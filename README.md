🌐 Language: En | [中文](README_ch.md)

---

# Problab Scaffold

<sub>Maintained by <b>Zintix Labs</b> — <a href="https://github.com/nextso">@nextso</a></sub>

A production-ready Problab scaffold for slot game R&D. It is designed to get you running fast and keep game development focused on configs and logic only.

This repository uses the Problab architecture.

## 1‑Minute Checklist

1. Clone  
   `git clone https://github.com/zintix-labs/problab-scaffold.git`

2. Run  
   `make run` (or `go run ./cmd/run`)

3. Privatize  
   - Change `go.mod` module name  
   - Replace all import paths  
   - (Recommended) `git init`

4. Done  
   Start adding games via `internal/configs/` and `internal/logic/`.

---

## Why this scaffold

### A. Ultra-fast project bootstrap

Goal: **run successfully in 1 minute, and be ready for commercial development in 3 minutes**

**Step 1: Get the project**

```bash
git clone https://github.com/zintix-labs/problab-scaffold.git
cd problab-scaffold
```

**Step 2: Run**

With Makefile (recommended):

```bash
make run
make dev
make svr
```

Without Makefile:

```bash
go run ./cmd/run
go run ./cmd/dev
go run ./cmd/svr
```

### Step 3: Privatize (Recommended)

#### 1. Change the Go module name

Edit `go.mod` and replace the module name with your own project module:

```go
module <your project module name>
```

> The module name does not need to be hosted on GitHub.

#### 2. Replace all import paths

**VS Code (recommended)**  
Search:
```
github.com/zintix-labs/problab-scaffold
```
Replace with:
```
<your project module name>
```

**Command line**

```bash
grep -rl "github.com/zintix-labs/problab-scaffold" . \
  | xargs sed -i '' 's|github.com/zintix-labs/problab-scaffold|<your project module name>|g'
```

#### 3. Reinitialize Git history (recommended)

```bash
rm -rf .git
git init
```

The project is now fully private and ready for commercial development.

---

## Extremely simple development flow

To create a new game, you only need two changes:

1. Add a config in `internal/configs/` (copy a demo YAML and modify it).  
2. Add logic in `internal/logic/` (copy a demo game and modify the logic).

That is all. You now have a **production-ready** new game. Development has never been this straightforward.

## Optimal workflow (Problab v0.7.0)

The scaffold uses the regular-file directory `internal/optimal` as its single
Optimal root. Problab read-only mmaps finalized Artifact bundles from this
directory on supported platforms. `make run`, `make dev`, `make svr`, and
`make opt` all construct the engine through `engine.New()`.

### The optimizer is configured, not invoked with flags

`cmd/opt` is the Problab v0.7.0 **LP / math-intent optimizer (v2)**. It takes
**no command-line arguments** — `make opt` runs it as-is and it rejects any
flag. Every run is defined by the embedded file `cmd/opt/opt_cfg.yaml`
(`//go:embed`), which is the single auditable source for a run. Changing a run
means editing that file and rebuilding.

`opt_cfg.yaml` has three sections with deliberately separated ownership:

| Section | Owner | Contents |
| --- | --- | --- |
| `plans:` | Engineering | Routing and execution: target `game` + `bet_modes`, `seed`, collection `workers` / `batch_size` / `max_spins`, output `format` and `directory`. Plans run in declaration order; use one plan per bet mode, or list several in `bet_modes`. |
| `intents:` | Designer | The math contract: class weights, inclusive `win_range`, collection `tags`, exact `exp` and `median`, bucketed `main_experience` groups with a soft `prefer` profile, overall `cv` bounds, optional collision `risk` policy. |
| `engine_options:` | Engineering | Numerical controls only: feasibility / optimality / quantile tolerances and bisection iteration budgets. Never relaxed automatically after an infeasible result. |

Collection tags referenced by `intents[].collect.tags` are owned by the scaffold
in `internal/logic/optimizer_tags/`. `GameTags` binds each `spec.GID` to its tag
predicate set; add one file per game and register it there.

### Steps for a new game such as `game_id: 1901`

1. Develop and collect with Optimal disabled:

   ```yaml
   optimal_setting:
     use_optimal: false
   ```

2. Register the game's collection tags in `internal/logic/optimizer_tags/`
   (copy `demo_0_tags.go`, add the new `spec.GID` to `GameTags`).

3. Add a plan per bet mode in `cmd/opt/opt_cfg.yaml` (set `target.game: 1901`,
   `target.bet_modes`, `intent`, `seed`, collection budget), then run:

   ```bash
   make opt
   ```

4. Review the generated bundle in `build/optimizer/artifact_v1/game_1901/`; the
   legacy gacha exchange files are written under
   `build/optimizer/gacha/game_1901/`, and a per-mode distribution CSV is
   written next to the output. Incomplete runs stay in a `.pending` staging
   directory and are invisible to the runtime until all bet modes are present
   and published atomically.

5. Copy the complete approved `artifact_v1/game_1901/` directory to
   `internal/optimal/game_1901/`.

6. Enable the finalized result:

   ```yaml
   optimal_setting:
     use_optimal: true
     artifact: game_1901/manifest.json
   ```

The manifest resolves the probability, alias, and SeedBank files. Game configs
do not list those paths separately.

The Docker image copies `internal/optimal` to the same path under `/app`, so the
local and container layouts are identical. A deployment may instead mount a
finalized artifact directory read-only at `/app/internal/optimal`.

## Quick architecture notes

- Configs are embedded from `internal/configs/`.
- The config filesystem is **flat**: use `*.yaml` files in that folder (no subfolders).
- Logic is registered via `init()` in `internal/logic/` to the global registry.
- `pkg/engine/problab.go` owns one private `WithOptimalDir("internal/optimal")`
  wiring point, so commands keep calling `engine.New()` unchanged.
- Finalized Optimal bundles live under `internal/optimal/game_<gid>/` and are
  regular files rather than embedded assets, allowing read-only mmap.
- Optimizer runs are defined only by the embedded `cmd/opt/opt_cfg.yaml`; the
  binary takes no flags. Collection tags live in
  `internal/logic/optimizer_tags/`, where `GameTags` maps each `spec.GID` to its
  tag predicates.
- A single Problab instance loads/maps each referenced Artifact once and shares
  it across every Machine, pool, and Simulator worker.
- Application entrypoints close Problab on shutdown so mapped files are released.

## Commands

- `make run` : Run simulator (default `game=0`)  
- `make svr` : Run HTTP server  
- `make dev` : Run Dev web panel  
- `make help` : Show all targets and args

## Requirements

- Go 1.25+ (see `go.mod`)
- Problab v0.7.0
- `make` (optional but recommended)

## License

Apache-2.0
