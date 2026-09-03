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

## Optimal workflow (Problab v0.8.0)

The scaffold uses the regular-file directory `internal/optimal` as its single
Optimal root. Runtime entrypoints (`make run`, `make dev`, and `make svr`) call
`engine.New()` and may read finalized bundles from this root. `make opt` calls
`engine.NewForOptimizer()` instead: it uses the same project-owned PRNG factory,
embedded configs, and logic registry, but deliberately does not preload a
published bundle while collecting raw game results.

Problab v0.8.0 supports two mutually exclusive runtime formats:

- **Artifact v1 (recommended):** a manifest plus binary probability, alias, and
  Seed Bank files. The manifest records file hashes and the PRNG snapshot
  format.
- **Legacy gacha:** one `gacha_<mode>.json.zst` and one
  `seed_bank_<mode>.bin` for every bet mode. This remains supported for existing
  integrations.

### The optimizer is configured, not invoked with flags

`cmd/opt` is the Problab v0.8.0 **LP / math-intent optimizer (v2)**. It takes
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

4. Review the generated output:

   ```text
   build/optimizer/
   |-- artifact_v1/game_1901/
   |-- gacha/game_1901/
   `-- distribution_mode_<mode>.csv
   ```

   Incomplete runs stay in a `.pending` staging directory and are invisible to
   runtime consumers until all bet modes are present and published atomically.

5. Publish either Artifact v1 or legacy gacha by following one of the sections
   below. Do not configure both formats for the same game.

### Recommended: publish Artifact v1

Copy the complete approved directory without selecting individual files:

```text
build/optimizer/artifact_v1/game_1901/
    -> internal/optimal/artifact_v1/game_1901/
```

The installed directory contains one manifest and the files for every bet mode:

```text
internal/optimal/artifact_v1/game_1901/
|-- manifest.json
|-- mode_0.json
|-- prob_0.bin
|-- aliases_0.bin
|-- seed_bank_0.bin
`-- ... files for additional bet modes
```

Enable the finalized result in the game config:

```yaml
optimal_setting:
  use_optimal: true
  artifact: artifact_v1/game_1901/manifest.json
```

The manifest resolves the probability, alias, and SeedBank files. Game configs
do not list those paths separately.

### Existing integrations: publish legacy gacha

Legacy users may continue using the `.json.zst` plus Seed Bank layout. Copy the
complete approved legacy directory:

```text
build/optimizer/gacha/game_1901/
    -> internal/optimal/gacha/game_1901/
```

For one bet mode the installed files are:

```text
internal/optimal/gacha/game_1901/
|-- gacha_0.json.zst
`-- seed_bank_0.bin
```

Reference both files explicitly:

```yaml
optimal_setting:
  use_optimal: true
  gachas:
    - gacha/game_1901/gacha_0.json.zst
  seed_bank:
    - gacha/game_1901/seed_bank_0.bin
```

For multiple bet modes, add one path to each list in bet-mode order. Both lists
must have the same length as `bet_units`: index `0` maps to `bet_units[0]`, index
`1` maps to `bet_units[1]`, and so on. Do not set `artifact` when `gachas` and
`seed_bank` are configured.

### PRNG and snapshot compatibility

An Optimal Seed Bank stores serialized PRNG snapshots. The optimizer and
runtime must therefore use compatible PRNG implementations and snapshot sizes.
The v0.8.0 default is ChaCha20. Existing PCG64 bundles can continue to run only
when the project explicitly injects the compatible PCG64 factory; changing the
directory or YAML path does not convert a PCG64 bundle into a ChaCha20 bundle.
Regenerate either output format after changing PRNG implementation.

### Private artifacts and deployment

`internal/optimal/artifact_v1/**` and `internal/optimal/gacha/**` are ignored by
Git by default. Optimal outputs can contain proprietary math data, and Seed Bank
files can be too large for ordinary Git hosting. The recommended production
flow is to store approved bundles in a private artifact service and copy or
read-only mount the selected bundle into `internal/optimal` during deployment.

If your organization intentionally versions bundles in a private repository,
remove the corresponding `internal/optimal` rule from `.gitignore` and use Git
LFS or another large-file storage system. Private GitHub repositories retain the
normal per-file size limit. Legacy `.zst` files are also covered by the global
`*.zst` ignore rule, so that rule must be removed or overridden before tracking
legacy output. Never commit proprietary Optimal data to a public repository.

## Quick architecture notes

- Configs are embedded from `internal/configs/`.
- The config filesystem is **flat**: use `*.yaml` files in that folder (no subfolders).
- Logic is registered via `init()` in `internal/logic/` to the global registry.
- `pkg/engine/problab.go` owns the shared PRNG, config, and logic wiring.
  Runtime commands call `engine.New()`; `make opt` calls
  `engine.NewForOptimizer()` to collect without a published Artifact.
- Finalized Artifact v1 bundles live under
  `internal/optimal/artifact_v1/game_<gid>/`; legacy bundles live under
  `internal/optimal/gacha/game_<gid>/`.
- Artifact files are regular files rather than embedded assets, allowing
  read-only mmap for Artifact v1 on supported platforms.
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
- Problab v0.8.0
- `make` (optional but recommended)

## License

Apache-2.0
