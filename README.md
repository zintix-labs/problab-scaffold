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

## Quick architecture notes

- Configs are embedded from `internal/configs/`.
- The config filesystem is **flat**: use `*.yaml` files in that folder (no subfolders).
- Logic is registered via `init()` in `internal/logic/` to the global registry.

## Commands

- `make run` : Run simulator (default `game=0`)  
- `make svr` : Run HTTP server  
- `make dev` : Run Dev web panel  
- `make help` : Show all targets and args

## Requirements

- Go 1.25+ (see `go.mod`)  
- `make` (optional but recommended)

## License

Apache-2.0
