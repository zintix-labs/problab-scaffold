# `internal/optimal/`

Read-only mmap root for Optimal Artifacts. `pkg/engine/problab.go` wires this single directory for the whole process:

```go
const optimalDir = "internal/optimal"
optimalOption problab.ProblabOption = problab.WithOptimalDir(optimalDir)
```

Each game opts in per-config, under `optimal_setting` in its `internal/configs/*.yaml`, by pointing at a path relative to this root:

```yaml
optimal_setting:
  use_optimal: true
  artifact: artifact_v1/game_0/manifest.json     # new format — see "Choosing a format" below
  # -- or, mutually exclusive with `artifact` --
  # gachas:
  #   - gacha/game_0/gacha_0.json.zst
  # seed_bank:
  #   - gacha/game_0/seed_bank_0.bin
```

A game uses exactly one format, never both.

## Choosing a format

| | `artifact_v1/` (recommended) | `gacha/` (legacy) |
|---|---|---|
| On-disk shape | Split binary files: `prob_*.bin`, `aliases_*.bin`, `seed_bank_*.bin`, plus a `manifest.json` that indexes them | One zstd-compressed JSON blob (`gacha_*.json.zst`) plus a separate `seed_bank_*.bin` |
| Load-time integrity | Every file in `manifest.json` carries a `sha256`, verified against the file on disk at load time | No per-file hash — the loader trusts the JSON payload as-is |
| Read path | Alias table and seed bank are read directly via mmap, no decompression/parse step | The JSON blob must be decompressed and parsed before use |
| Status in `internal/optimalrt` | Actively developed binary format | Kept only so deployments that already shipped this format keep working |

**Default to `artifact_v1/` for any new game.** It's the actively maintained path, and the per-file SHA-256 check in `manifest.json` catches a corrupted or partially-copied artifact at load time instead of at some unpredictable point during play — `gacha/` has no equivalent check. Reach for `gacha/` only when you're maintaining or migrating a game that was already published against it; don't start a new game there.

## Directory layout

One subfolder per game id under each format, named `game_<GID>/`:

```
optimal/
  artifact_v1/
    game_0/
      manifest.json               # schema_version, artifact_id, snapshot_format, modes[] (SHA-256-verified file refs)
      mode_0.json                  # standalone copy of modes[0], one per configured bet mode
      prob_0.bin                    # alias-table probabilities for bet mode 0
      aliases_0.bin                 # alias-table indices for bet mode 0
      seed_bank_0.bin                 # N deterministic replay seeds (snapshot_format-encoded)
      game_0_mode_0_v1.yaml            # the exact opt_cfg.yaml that produced everything above
  gacha/
    game_0/
      gacha_0.json.zst              # zstd-compressed JSON probability/alias description for bet mode 0
      seed_bank_0.bin                 # same deterministic replay seeds, gacha-format layout
      game_0_mode_0_v1.yaml            # the exact opt_cfg.yaml that produced the two files above
```

Add one `mode_<N>`/`prob_<N>`/`aliases_<N>` set (or `gacha_<N>.json.zst`) per configured bet mode.

## Reproducibility: keep the config that built it, whichever format you use

Every published artifact, in either format, ships together with the exact `opt_cfg.yaml` that `cmd/opt` ran to produce it — copied into the same `game_<GID>/` folder and renamed `game_<GID>_mode_<N>_v1.yaml`. This is not optional. Without it, nobody can:

- rebuild the artifact bit-for-bit if the generated files are lost, corrupted, or need to move to a new machine;
- re-run the library's `VerifyRuntimeMaterialized` against a later game-logic version to check for drift between what was verified at build time and what current code would produce;
- tell a reviewer, teammate, or auditor which config produced which `artifact_id` / manifest hashes.

Treat this file as the single source of truth for "what generated this." If you regenerate an artifact from a changed config, replace this file in the same commit as the new binaries — never leave a stale yaml sitting next to binaries it no longer describes.

## What's tracked in git, what isn't

`.gitignore` scopes the generated files out of both folders:

```
/internal/optimal/artifact_v1/**/*.json
/internal/optimal/artifact_v1/**/*.bin
/internal/optimal/gacha/**/*.json
/internal/optimal/gacha/**/*.zst
/internal/optimal/gacha/**/*.bin
```

So `manifest.json`, `mode_*.json`, `prob_*.bin`, `aliases_*.bin`, `seed_bank_*.bin`, and `gacha_*.json.zst` are **not** committed — they're large (a full seed bank can run into the hundreds of MB), fully regenerable from the retained yaml, and `seed_bank_*.bin` in particular carries replayable PRNG snapshot data that doesn't need to sit in version-control history. Only each `game_<GID>_mode_<N>_v1.yaml` and this `README.md` are tracked.

## Rebuilding / publishing a new artifact

The flow is the same for both formats; only the target subfolder and config keys differ.

1. Put the desired config in place as `cmd/opt/opt_cfg.yaml` (start from a saved `game_<GID>_mode_<N>_v1.yaml` if you're reproducing or extending an existing artifact), then run `make opt` (or `go run ./cmd/opt`).
2. `cmd/opt` writes its output to `build/optimizer/<format>/game_<GID>/`, staging through a `game_<GID>.pending/` directory before the final files land — `<format>` is `artifact_v1` or `gacha` depending on which output mode you ran.
3. Once you're satisfied with the run, copy `build/optimizer/<format>/game_<GID>/*` into `internal/optimal/<format>/game_<GID>/`, then copy the `opt_cfg.yaml` that produced it alongside as `game_<GID>_mode_<N>_v1.yaml`.
4. Point the game's config under `internal/configs/` at it: set `use_optimal: true` and either `artifact: artifact_v1/game_<GID>/manifest.json`, or `gachas`/`seed_bank` under the `gacha/` path — never both.
