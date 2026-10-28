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

## Choice games (segmented games)

`demo_2` (`demo_choice`) and `demo_3` (`demo_entry`) show games that stop at a
choice point: a segment that is not the game's end lists its entries in
`ChooseNext`, and the next segment runs the entry at the chosen position.

- `internal/logic/demo_2.go` passes entry state to the next segment through a
  registered checkpoint (`dto.RegisterCheckpoint`).
- `internal/logic/demo_3.go` reads no checkpoint: each segment depends only on
  its entry value and seed, so each entry can have its own Optimal pair.

Simulations need a choose mode for these games. `make run` always picks
`first`; for another mode, run the CLI directly:

```bash
make run game=2
go run ./cmd/run -game 2 -choose last
```

## Player experience analysis

`make exp` analyzes finite player sessions with the embedded
`cmd/exp/exp_cfg.yaml`; `config=<path>` selects another YAML:

```bash
make exp
make exp config=cmd/exp/examples/game_3.yaml
```

Strategies are factories registered per game in `internal/strategies/`
(`game_<gid>.go`, listed in `reg.go`). The built-in `builtin.Monotone`
(`sdk/strategy/builtin`) bets a fixed wager; its `Choose` picks a position at
choice points (`ChooseFirst`, `ChooseLast` or `ChooseIter`). A YAML can cap the
segments of one game with `player.max_segments` (default 10).

## Optimal workflow (Problab v0.17.1)

The scaffold uses the regular-file directory `internal/optimal` as its single
Optimal root. Runtime entrypoints (`make run`, `make dev`, and `make svr`) call
`engine.New()` and may read published gacha pairs from this root. `make opt`
calls `engine.NewForOptimizer()` instead: it uses the same project-owned PRNG
factory, embedded configs, and logic registry, but deliberately does not
preload a published pair while collecting raw game results.

Each bet mode is published as a gacha pair: one `gacha.json.zst` and one
`seed_bank.bin`.

### The optimizer is configured, not invoked with flags

`cmd/opt` is the Problab **LP / math-intent optimizer (v2)**. It takes
**no command-line arguments** — `make opt` runs it as-is and it rejects any
flag. Every run is defined by the embedded file `cmd/opt/opt_cfg.yaml`
(`//go:embed`), which is the single auditable source for a run. Changing a run
means editing that file and rebuilding. The first Ctrl+C stops after the
current solve; a second Ctrl+C exits immediately.

`opt_cfg.yaml` has three sections with deliberately separated ownership:

| Section | Owner | Contents |
| --- | --- | --- |
| `plans:` | Engineering | Routing and execution: target `game` + `bet_modes` (and the optional `input` of a continuation segment), `seed`, collection `workers` / `batch_size` / `max_spins`, output `format` and `directory`. Plans run in declaration order; each plan owns exactly one bet mode. |
| `intents:` | Designer | The math contract: class weights, inclusive `win_range`, collection `tags`, exact `exp` and `median`, bucketed `main_experience` groups with a soft `prefer` profile, overall `cv` bounds, optional collision `risk` policy. |
| `engine_options:` | Engineering | Numerical controls only: feasibility / optimality / quantile tolerances and bisection iteration budgets. Never relaxed automatically after an infeasible result. |

Collection tags referenced by `intents[].collect.tags` are owned by the scaffold
in `internal/tags/`. `Catalog` binds each `spec.GID` to its tag
predicate set; add one file per game and register it there.

RGS outputs (`rgs-collected`, `rgs-optimized`) encode each record with
`useConverter` in `cmd/opt/main.go`. The default emits the whole DTO, including
PRNG state; `internal/converter` also provides `WithoutSnapshotConverter`, which
drops it.

### Steps for a new game such as `game_id: 1901`

1. Develop and collect with Optimal disabled:

   ```yaml
   optimal_setting:
     use_optimal: false
   ```

2. Register the game's collection tags in `internal/tags/`
   (copy `demo_0_tags.go`, add the new `spec.GID` to `Catalog`).

3. Add a plan per bet mode in `cmd/opt/opt_cfg.yaml` (set `target.game: 1901`,
   `target.bet_modes`, `intent`, `seed`, collection budget), then run:

   ```bash
   make opt
   ```

4. Review the generated output:

   ```text
   build/optimizer/
   |-- gacha/game_1901/mode_<mode>/<id>/
   |   |-- gacha.json.zst
   |   `-- seed_bank.bin
   `-- distribution_mode_<mode>.csv
   ```

   Each bet mode is published on its own. The pair is staged, checked with the
   production loader, and renamed into place; `<id>` comes from the content, so
   a published directory is never rewritten.

5. Copy the approved directory, together with the `opt_cfg.yaml` that produced
   it:

   ```text
   build/optimizer/gacha/game_1901/mode_0/<id>/
       -> internal/optimal/gacha/game_1901/mode_0/<id>/
   ```

6. Reference the pair in the game config:

   ```yaml
   optimal_setting:
     use_optimal: true
     gachas:
       - gacha/game_1901/mode_0/<id>/gacha.json.zst
     seed_bank:
       - gacha/game_1901/mode_0/<id>/seed_bank.bin
   ```

   For multiple bet modes, add one path to each list in bet-mode order. Both
   lists must have the same length as `bet_units`: index `0` maps to
   `bet_units[0]`, index `1` maps to `bet_units[1]`, and so on.

`gachas`/`seed_bank` is the compatible shorthand. Any game can instead list its
pairs under `optimal_setting.modes`, one entry per bet mode: `start` for the
initial segment and, for a choice game, `after` keyed by entry value (collected
with `plans[].target.input`). See `internal/optimal/README.md`.

### PRNG and snapshot compatibility

An Optimal Seed Bank stores serialized PRNG snapshots. The optimizer and
runtime must therefore use compatible PRNG implementations and snapshot sizes.
The default is ChaCha20. Existing PCG64 pairs can continue to run only when the
project explicitly injects the compatible PCG64 factory; changing the directory
or YAML path does not convert a PCG64 pair into a ChaCha20 pair. Regenerate the
output after changing PRNG implementation.

### Private outputs and deployment

`internal/optimal/**` generated files are ignored by Git by default. Optimal
outputs can contain proprietary math data, and Seed Bank files can be too large
for ordinary Git hosting. The recommended production flow is to store approved
pairs in a private artifact service and copy or read-only mount the selected
pairs into `internal/optimal` during deployment.

If your organization intentionally versions pairs in a private repository,
remove the corresponding `internal/optimal` rules from `.gitignore` and use Git
LFS or another large-file storage system. Private GitHub repositories retain the
normal per-file size limit. `.zst` files are also covered by the global `*.zst`
ignore rule, so that rule must be removed or overridden before tracking them.
Never commit proprietary Optimal data to a public repository.

## Quick architecture notes

- Configs are embedded from `internal/configs/`.
- The config filesystem is **flat**: use `*.yaml` files in that folder (no subfolders).
- Logic is registered via `init()` in `internal/logic/` to the global registry.
- `pkg/engine/problab.go` owns the shared PRNG, config, and logic wiring.
  Runtime commands call `engine.New()`; `make opt` calls
  `engine.NewForOptimizer()` to collect without a published pair.
- Published gacha pairs live under
  `internal/optimal/gacha/game_<gid>/mode_<mode>/<id>/`.
- Optimal files are regular files rather than embedded assets, allowing
  read-only mmap of seed banks on supported platforms.
- Optimizer runs are defined only by the embedded `cmd/opt/opt_cfg.yaml`; the
  binary takes no flags. Collection tags live in
  `internal/tags/`, where `Catalog` maps each `spec.GID` to its
  tag predicates.
- Player strategies live in `internal/strategies/`, where `Catalog` maps each
  `spec.GID` to its strategy factories.
- A single Problab instance loads/maps each referenced pair once and shares
  it across every Machine, pool, and Simulator worker.
- Application entrypoints close Problab on shutdown so mapped files are released.

## Commands

- `make run` : Run simulator (default `game=0`; choice games pick `first`)
- `make exp` : Analyze finite player experiences (`config=<path>` for another YAML)
- `make opt` : Run optimizer (embedded `cmd/opt/opt_cfg.yaml`)
- `make svr` : Run HTTP server
- `make dev` : Run Dev web panel
- `make help` : Show all targets and args

## Requirements

- Go 1.27.3+ (see `go.mod`)
- Problab v0.17.1
- `make` (optional but recommended)

## License

Apache-2.0
