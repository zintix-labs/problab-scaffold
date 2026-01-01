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

// Package main provides the scaffold's server entrypoint.
//
// This command is intentionally thin: it only parses flags, wires a default
// Problab engine (configs + logic registry), and starts the HTTP server.
//
// The goal is to make `go run ./cmd/svr` (or `make svr`) work out-of-the-box
// for new adopters, while keeping all Problab engine code inside the upstream
// `github.com/zintix-labs/problab` module.

package main

import (
	"flag"
	"fmt"

	"github.com/zintix-labs/problab-scaffold/pkg/engine"
	"github.com/zintix-labs/problab/server"
	"github.com/zintix-labs/problab/server/logger"
	"github.com/zintix-labs/problab/server/svrcfg"
)

// main loads runtime configuration from flags and starts the Problab HTTP server.
//
// Any configuration/engineping error is treated as fatal, because a partially
// initialized server is almost always the wrong behavior for an example scaffold.
func main() {
	cfg, err := loadConfigFromFlags()
	if err != nil {
		fmt.Println(err)
		return
	}
	server.Run(cfg)
}

// config holds CLI flag values.
//
// NOTE: This is intentionally small and human-readable. The scaffold is meant to
// be copied and modified by end users.
type config struct {
	Log         string // Logger mode: dev|prod|discard
	Mode        string // Server run mode: dev|prod
	SlotBufSize int    // Number of machine instances per game (pool size)
}

// loadConfigFromFlags parses CLI flags and builds a `svrcfg.SvrCfg`.
//
// Flags:
//
//	-log   : dev|prod|discard
//	         dev     -> developer-friendly console logs
//	         prod    -> production-style logs
//	         discard -> silence all logs (useful for raw benchmarking)
//	-buf   : number of machine instances per game (pool size)
//	-mode  : dev|prod
//	         dev  -> enables development/debugging HTTP endpoints (unsafe for public exposure)
//	         prod -> exposes only production-safe HTTP surface
//
// Defaults are chosen to be safe and predictable:
//
//	-log  defaults to "dev" for local visibility.
//	-mode defaults to "prod" to avoid accidentally exposing dev endpoints.
func loadConfigFromFlags() (*svrcfg.SvrCfg, error) {
	cfg := new(config)
	flag.StringVar(&cfg.Log, "log", "dev", "log mode: dev|prod|discard")
	flag.StringVar(&cfg.Mode, "mode", "prod", "svr mode: dev|prod")
	flag.IntVar(&cfg.SlotBufSize, "buf", 3, "number of machine instances per game")

	flag.Parse()

	// Print the raw flag values exactly as provided by the user.
	fmt.Printf("[scaffold][flags] -log=%s -mode=%s -buf=%d\n", cfg.Log, cfg.Mode, cfg.SlotBufSize)

	// Create an async logger with a small internal buffer. Most users should keep this as-is.
	log := logger.NewDefaultAsyncLogger(cfg.normLog())

	// engine wires configs (FS) + logic registry into a ready-to-run Problab instance.
	// This is the main value of the scaffold: users add YAML configs and logic builders,
	// and the rest is assembled for them.
	pb := engine.MustNew()

	// Assemble the server configuration used by `problab/server`.
	sCfg := &svrcfg.SvrCfg{
		Log:         log,
		SlotBufSize: cfg.SlotBufSize,
		Problab:     pb,
		Mode:        cfg.normSvrMode(),
	}
	return sCfg, nil
}

// normLog converts the `-log` flag into the logger's enum.
//
// Unknown values fall back to dev logging for better local ergonomics.
func (cfg *config) normLog() logger.LogMode {
	switch cfg.Log {
	case "dev":
		return logger.ModeDev
	case "prod":
		return logger.ModeProd
	case "discard":
		return logger.ModeSilence
	default:
		return logger.ModeDev
	}
}

// normSvrMode converts the `-mode` flag into the server RunMode.
//
// Unknown values fall back to production mode to avoid accidentally exposing
// debugging/simulation endpoints.
func (cfg *config) normSvrMode() svrcfg.RunMode {
	switch cfg.Mode {
	case "dev":
		return svrcfg.ModeDev
	case "prod":
		return svrcfg.ModeProd
	default:
		return svrcfg.ModeProd
	}
}
