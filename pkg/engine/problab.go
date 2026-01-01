// Copyright 2025 Zintix Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package engine provides a thin, opinionated wiring layer for building a runnable
// Problab instance in this scaffold repository.
//
// This is the main integration boundary between your game project (configs + logic)
// and the upstream Problab engine.
//
// Customize points:
//   - PRNG / core factory: swap `core.Default()` with your own deterministic PRNGFactory.
//   - Config FS: you may mount multiple fs.FS sources (e.g., embed + os.DirFS), but
//     keeping a single source is recommended for simplicity.
//   - Logic registry: you may register multiple logic sets, but keeping one registry
//     is recommended to reduce operational complexity.
//
// The goal is to keep app code and cmd/* entrypoints clean: they can call `engine.New()`
// (or `engine.MustNew()` for CLI/dev) and focus on running simulation/server/dev.
package engine

import (
	"io/fs"

	"github.com/zintix-labs/problab"
	"github.com/zintix-labs/problab-scaffold/internal/configs"
	"github.com/zintix-labs/problab-scaffold/internal/logic"
	"github.com/zintix-labs/problab/sdk/core"
	"github.com/zintix-labs/problab/sdk/slot"
)

// NOTE: This scaffold intentionally does NOT expose runtime dependency injection
// (no Options, setters, env/flag-driven wiring, or external overrides).
//
// Why?
//   - The engine wiring (PRNG, config FS, logic registry) is a build-time decision owned by
//     the game/math side, not something application code should be able to toggle at runtime.
//   - Keeping these choices in code makes changes explicit, reviewable, and versioned.
//   - It prevents accidental or unsafe production deployments caused by misconfigured injections.
//   - It reduces the debugging surface area ("which PRNG/configs/logics is this process using?")
//
// If you need to customize wiring, edit the private variables below and ship a new version.
var (
	// PRNGFactory: replace this with your own deterministic PRNG implementation if needed.
	// The engine only depends on the PRNG interface/factory, not a specific algorithm.
	// See package `github.com/zintix-labs/problab/sdk/core` for the `PRNG` and `PRNGFactory` interface definitions.
	// (On GitHub, the source lives under `github.com/zintix-labs/problab/blob/main/sdk/core/core.go`.)
	pRNGFactory core.PRNGFactory = core.Default()
	// Config FS: provide game settings/spec files (usually embedded via `embed`).
	//
	// Default expectation (recommended): a flat FS layout.
	// Your provided fs.FS should map directly to the config folder itself:
	//   - files are addressed by filename only (no nested paths)
	//   - no subdirectories are used for config lookup
	// This keeps config loading deterministic and avoids runtime path dependencies.
	//
	// Advanced: you can mount multiple fs.FS sources via catalog/multi-FS patterns,
	// but start with a single FS to keep maintenance simple.
	cfgs []fs.FS = problab.Configs(configs.FS)
	// Logic registry: register your game logic builders/handlers.
	// You can merge multiple registries, but a single registry is easiest to reason about.
	logics []*slot.LogicRegistry = problab.Logics(logic.Logics)
)

// New constructs a Problab instance using the scaffold's embedded configs and logic registry.
//
// It returns an error so production callers can decide how to report/handle engine failures.
// Typical usage:
//
//	pb, err := engine.New()
func New() (*problab.Problab, error) {
	pb, err := problab.NewAuto(pRNGFactory, cfgs, logics)
	if err != nil {
		return nil, err
	}
	return pb, nil
}

// MustNew is a convenience helper for CLI/dev entrypoints.
// It panics on error (instead of exiting the process), keeping this package usable as a library.
func MustNew() *problab.Problab {
	pb, err := New()
	if err != nil {
		panic(err)
	}
	return pb
}
