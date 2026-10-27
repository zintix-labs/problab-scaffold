# `internal/optimal/`

Read-only root for Optimal gacha pairs. `pkg/engine/problab.go` wires this single directory for the whole process:

```go
const optimalDir = "internal/optimal"
optimalOption problab.ProblabOption = problab.WithOptimalDir(optimalDir)
```

Each game opts in per-config, under `optimal_setting` in its `internal/configs/*.yaml`, with paths relative to this root. The compatible shorthand lists one pair per bet mode, in bet-mode order:

```yaml
optimal_setting:
  use_optimal: true
  gachas:
    - gacha/game_0/mode_0/<id>/gacha.json.zst
  seed_bank:
    - gacha/game_0/mode_0/<id>/seed_bank.bin
```

Any game can instead list its pairs under `modes`, one entry per bet mode (`modes` excludes `gachas` and `seed_bank`). `start` is the initial segment's pair, and is all a single-segment game needs. A choice game adds `after`, keyed by entry value (`ChooseNext[Choice]`); every entry the game offers must be listed. `{dist: null, bank: null}` runs that segment on the machine RNG:

```yaml
optimal_setting:
  use_optimal: true
  modes:
    - start: {dist: gacha/game_3/mode_0/<id>/gacha.json.zst, bank: gacha/game_3/mode_0/<id>/seed_bank.bin}
      after:
        3: {dist: null, bank: null}
        7: {dist: gacha/game_3/mode_0/<id7>/gacha.json.zst, bank: gacha/game_3/mode_0/<id7>/seed_bank.bin}
```

## Directory layout

`cmd/opt` publishes every bet mode on its own, in an immutable directory named by its content id:

```
optimal/
  gacha/
    game_0/
      mode_0/
        <id>/
          gacha.json.zst      # zstd-compressed probability/alias description
          seed_bank.bin       # deterministic replay seeds, read-only mmap
          opt_cfg.yaml        # the exact opt_cfg.yaml that produced the pair
```

Copy the published `<id>/` directory as-is. A new run with different content gets a new `<id>`, so older pairs stay in place until no config references them.

## Reproducibility: keep the config that built it

Every published pair ships together with the exact `opt_cfg.yaml` that `cmd/opt` ran to produce it, copied into the same `<id>/` directory. This is not optional. Without it, nobody can:

- rebuild the pair if the generated files are lost, corrupted, or need to move to a new machine;
- check a later game-logic version for drift against what was verified at build time;
- tell a reviewer, teammate, or auditor which config produced which `<id>`.

Treat this file as the single source of truth for "what generated this."

## What's tracked in git, what isn't

`.gitignore` scopes the generated files out of this root:

```
/internal/optimal/**/*.json
/internal/optimal/**/*.zst
/internal/optimal/**/*.bin
```

So `gacha.json.zst` and `seed_bank.bin` are **not** committed — they're large (a full seed bank can run into the hundreds of MB), fully regenerable from the retained yaml, and `seed_bank.bin` in particular carries replayable PRNG snapshot data that doesn't need to sit in version-control history. Only each retained `opt_cfg.yaml` and this `README.md` are tracked.

## Rebuilding / publishing a new pair

1. Put the desired config in place as `cmd/opt/opt_cfg.yaml` (start from a retained `opt_cfg.yaml` if you're reproducing or extending an existing pair), then run `make opt` (or `go run ./cmd/opt`).
2. `cmd/opt` publishes each successful bet mode to `build/optimizer/gacha/game_<GID>/mode_<N>/<id>/`. The pair is staged, checked with the production loader, and renamed into place; rerunning with the same content reuses the directory.
3. Copy `build/optimizer/gacha/game_<GID>/mode_<N>/<id>/` into `internal/optimal/gacha/game_<GID>/mode_<N>/<id>/`, then copy the `opt_cfg.yaml` that produced it alongside.
4. Point the game's config under `internal/configs/` at it: set `use_optimal: true` and list the pair under `gachas`/`seed_bank` (or under `modes` for a choice game).
