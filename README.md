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

## Optimal workflow (Problab v0.6.0)

The scaffold uses the regular-file directory `internal/optimal` as its single
Optimal root. Problab read-only mmaps finalized Artifact bundles from this
directory on supported platforms. `make run`, `make dev`, `make svr`, and
`make opt` all continue to construct the engine through `engine.New()`.

For a new game such as `game_id: 1901`:

1. Develop and collect with Optimal disabled:

   ```yaml
   optimal_setting:
     use_optimal: false
   ```

2. Run the Optimizer for every bet mode:

   ```bash
   make opt game=1901 betmode=0
   ```

3. Review the generated bundle in `build/optimizer/game_1901/`.
4. Copy the complete approved directory to `internal/optimal/game_1901/`.
5. Enable the finalized result:

   ```yaml
   optimal_setting:
     use_optimal: true
     artifact: game_1901/manifest.json
   ```

The manifest resolves the probability, alias, and SeedBank files. Game configs
do not list those paths separately. The historical exchange files are still
written directly under `build/optimizer/` for existing external consumers.

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
- Problab v0.6.0
- `make` (optional but recommended)

## License

Apache-2.0
