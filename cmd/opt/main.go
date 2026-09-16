//go:build !poc

// Copyright 2025 Zintix Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/zintix-labs/problab-scaffold/internal/converter"
	"github.com/zintix-labs/problab-scaffold/internal/logic/game_tags"
	"github.com/zintix-labs/problab-scaffold/pkg/engine"
	"github.com/zintix-labs/problab/dto"
	optimizerv2 "github.com/zintix-labs/problab/optimizer/v2"
)

const embeddedConfigName = "opt_cfg.yaml"

// These two package variables are the command's only application-owned
// injection points. A game project repoints them at its own implementation
// without editing runV2, problab's optimizer/v2, or the embedded YAML; the
// values below are the defaults that ship with this scaffold.
var (
	// gameTags is a game_tags.GameTagCatalog, i.e.
	// map[spec.GID]map[string]optimizer.IsTag, where each optimizer.IsTag is a
	// func(*buf.SpinResult) bool predicate. WithCollectionTags stores it on the
	// Collector and resolves it per plan against that plan's target game, so a
	// GID no plan targets is simply unused and a game absent from the catalog
	// collects with the built-in tags only. The inner map keys are the names
	// that classes[].collect.tags.matches / .mismatches reference in
	// opt_cfg.yaml; bg and fg are built in and must not be redefined here.
	//
	// The scaffold catalog lives in internal/logic/game_tags (tags_catalog.go
	// binds GID 0 to demo_0_tags.go). To add a game, write its predicates next
	// to demo_0_tags.go and register them in game_tags.GameTags:
	//
	//	// internal/logic/game_tags/mygame_tags.go
	//	func IsFreeSpins(sr *buf.SpinResult) bool { ... }
	//	var MyGame_7_Tags = map[string]optimizer.IsTag{"free_spins": IsFreeSpins}
	//
	//	// internal/logic/game_tags/tags_catalog.go
	//	var GameTags = GameTagCatalog{0: Demo_0_Tags, 7: MyGame_7_Tags}
	gameTags game_tags.GameTagCatalog = game_tags.GameTags

	// useConverter is a dto.ResultConverter, i.e.
	// func(dto.SpinResult) (json.RawMessage, error). WithResultConverter binds
	// exactly one per Tuner, and both RGS families (rgs-collected and
	// rgs-optimized) call it once per replayed record. The returned bytes are
	// compacted and written verbatim as one line of results.jsonl.zst, so each
	// call must return exactly one valid UTF-8 JSON value; a converter error or
	// invalid JSON stops the Run instead of skipping the record. It cannot
	// affect the Parquet distribution or the native artifact_v1/gacha output,
	// and it is never called when opt_cfg.yaml lists no rgs-* format.
	//
	// The scaffold default, converter.IdentityConverter in internal/converter,
	// emits the whole dto.SpinResult exactly like dto.IdentityConverter. That
	// payload carries every record's start/after PRNG snapshots and checkpoint —
	// the complete seed material — so do not hand its output to a party that
	// must not receive it.
	//
	// The run report's converter field is decided by function identity, not by
	// behavior: only nil or dto.IdentityConverter itself report
	// converter=identity. Every other function, including the local
	// converter.IdentityConverter copy, reports converter=custom, which states
	// provenance only, not that replay state was actually removed.
	//
	// To supply your own, implement it in internal/converter (see problab's
	// optimizer/v2/examples/rgs for one that drops the replay state) and assign
	// it here:
	//
	//	useConverter dto.ResultConverter = converter.MyPlatformConverter
	useConverter dto.ResultConverter = converter.IdentityConverter
)

// main is deliberately a thin composition root: it loads the command-owned
// embedded config, constructs Problab plus optimizer/v2, and runs every declared
// plan in order. Collection, LP rows, bisection, materialization, and status
// mapping remain in optimizer/v2 and are shared by CLI and programmatic callers.
func main() {
	exitCode, err := runV2(os.Args[1:], os.Stdout, os.Stderr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "cmd/opt: %v\n", err)
		exitCode = 1
	}
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

// runV2 owns command lifecycle and returns a stable process classification:
// zero for a verified OPTIMAL mode that was durably staged (and may also have
// completed the all-mode manifest) or an EXPORTED rgs-collected-only plan, two
// for an expected typed non-success Run status, and an error for operational
// failures.
func runV2(arguments []string, _ io.Writer, stderr io.Writer) (int, error) {
	if len(arguments) != 0 {
		return 1, fmt.Errorf(
			"cmd/opt does not accept command-line parameters; edit embedded %q instead (got %s)",
			embeddedConfigName,
			strings.Join(arguments, " "),
		)
	}

	reporter := newCLIProgressReporter(stderr)
	configSource := fmt.Sprintf("embedded %q", embeddedConfigName)
	reporter.Report(optimizerv2.StageEvent{Stage: "load-config", State: "started", BetMode: -1, Message: configSource})
	configStarted := time.Now()
	config, err := loadV2Config()
	if err != nil {
		duration := time.Since(configStarted)
		var invalidConfig *invalidV2ConfigError
		if !errors.As(err, &invalidConfig) {
			reporter.Report(optimizerv2.StageEvent{Stage: "load-config", State: "failed", BetMode: -1, Message: err.Error(), Duration: duration})
			return 1, err
		}
		failedStage := invalidV2ConfigStage(invalidConfig)
		if failedStage == "static-validation" {
			reporter.Report(optimizerv2.StageEvent{Stage: "load-config", State: "completed", BetMode: -1, Duration: duration})
			reporter.Report(optimizerv2.StageEvent{Stage: failedStage, State: "started", BetMode: -1})
			reporter.Report(optimizerv2.StageEvent{Stage: failedStage, State: "failed", BetMode: -1, Message: err.Error()})
		} else {
			reporter.Report(optimizerv2.StageEvent{Stage: failedStage, State: "failed", BetMode: -1, Message: err.Error(), Duration: duration})
		}
		// A document that was read successfully but cannot be decoded or
		// validated is an expected, actionable Run outcome. Keep its typed value
		// internally, but present only one human-readable result on stderr.
		reportV2Outcome(stderr, configInvalidRunResult(invalidConfig, failedStage, duration))
		return 2, nil
	}
	reporter.Report(optimizerv2.StageEvent{Stage: "load-config", State: "completed", BetMode: -1, Duration: time.Since(configStarted)})
	lab, err := engine.NewForOptimizer()
	if err != nil {
		return 1, fmt.Errorf("construct optimizer Problab: %w", err)
	}
	defer func() { _ = lab.Close() }()
	tuner, err := optimizerv2.NewTuner(config, lab,
		optimizerv2.WithReporter(reporter),
		optimizerv2.WithCollectionTags(gameTags),
		optimizerv2.WithResultConverter(useConverter),
	)
	if err != nil {
		return 1, fmt.Errorf("construct optimizer/v2 tuner: %w", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	exitCode := 0
	for _, plan := range config.Plans {
		result, err := tuner.Run(ctx, optimizerv2.RunRequest{PlanID: plan.ID})
		// Report before inspecting err: some failing stages (for example
		// collection) still return a partial Report whose diagnostics explain
		// the failure. Stages such as RGS export return a zero RunResult with
		// the error, in which case this prints only a generic "Failed ()" line
		// and the returned error carries the real cause.
		reportV2Outcome(stderr, result)
		if err != nil {
			return 1, err
		}
		if !result.Succeeded() {
			exitCode = 2
			continue
		}
		// The verified distribution is a file artifact, not terminal noise: write
		// it next to the published output and print only where it landed. An
		// EXPORTED (rgs-collected only) Run skips the LP and has no Modes, so
		// nothing is written for it.
		paths, err := writeModeDistributionCSVs(plan.Output.Directory, result.Report.Modes)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "[warn] distribution report: %v\n", err)
		}
		for _, path := range paths {
			_, _ = fmt.Fprintf(stderr, "[Result] distribution written: %s\n", path)
		}
	}
	return exitCode, nil
}

// invalidV2ConfigError marks only successfully read configuration bytes whose
// YAML, schema, or semantic content is invalid. Keeping this distinction local
// to the command lets runV2 map bad user input to a typed RunResult while real
// filesystem failures continue through the operational Go-error/exit-one path.
type invalidV2ConfigError struct {
	source string
	cause  error
}

// Error includes the embedded source so a ConfigInvalid diagnostic identifies
// the command-owned document. Unwrap preserves yaml.v3 and optimizer/v2
// ConfigError inspection.
func (e *invalidV2ConfigError) Error() string {
	return fmt.Sprintf("load v2 config %s: %v", e.source, e.cause)
}

// Unwrap preserves the parser or ConfigError cause for callers that need
// structured inspection while runV2 relies only on this wrapper's provenance.
func (e *invalidV2ConfigError) Unwrap() error { return e.cause }

// parseV2ConfigBytes is the strict content boundary for embedded YAML.
// ParseConfig enforces KnownFields, one document, and all v2 semantic
// invariants; every failure after bytes have been obtained is therefore safe to
// classify as INFEASIBLE_CONFIG rather than an operational command failure.
func parseV2ConfigBytes(raw []byte, source string) (optimizerv2.Config, error) {
	config, err := optimizerv2.ParseConfig(raw)
	if err != nil {
		return optimizerv2.Config{}, &invalidV2ConfigError{source: source, cause: err}
	}
	return config, nil
}

// loadV2Config reads the command-owned document. cmd/opt intentionally has no
// external config or field-level override path: embeddedConfigName is the one
// auditable source of every RunPlan value executed by this binary.
func loadV2Config() (optimizerv2.Config, error) {
	raw, err := optConfig.ReadFile(embeddedConfigName)
	if err != nil {
		return optimizerv2.Config{}, fmt.Errorf("read embedded v2 config %q: %w", embeddedConfigName, err)
	}
	return parseV2ConfigBytes(raw, fmt.Sprintf("embedded %q", embeddedConfigName))
}

// configInvalidRunResult converts a strict configuration parse/validation
// failure into the public status contract used by Tuner.Run. No partially
// decoded Config is copied into the report because rejected input has no valid
// resolved plan or trustworthy effective values.
func configInvalidRunResult(err error, failedStage string, duration time.Duration) optimizerv2.RunResult {
	return optimizerv2.RunResult{
		Status: optimizerv2.StatusInfeasibleConfig,
		Diagnostics: optimizerv2.Diagnostics{{
			Code:        optimizerv2.DiagnosticConfigInvalid,
			Status:      optimizerv2.StatusInfeasibleConfig,
			Message:     err.Error(),
			SourcePaths: []string{"config"},
		}},
		Report: optimizerv2.RunReport{Stages: []optimizerv2.StageDuration{{Stage: failedStage, Duration: duration}}},
	}
}

// invalidV2ConfigStage distinguishes semantic Config validation from YAML/schema
// loading. This preserves the user's mental pipeline even though ParseConfig is
// intentionally one strict API that performs decoding and validation together.
func invalidV2ConfigStage(err error) string {
	var configError *optimizerv2.ConfigError
	if errors.As(err, &configError) {
		return "static-validation"
	}
	return "load-config"
}
