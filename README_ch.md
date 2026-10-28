🌐 语言：中文 | [En](README.md)

---

# Problab Scaffold

<sub><b>Zintix Labs</b> — <a href="https://github.com/nextso">@nextso</a></sub>

这是一个**可直接投入生产环境**的 Problab Scaffold，专为 Slot Game 数学与后端研发设计。  
目标是 **让你尽快跑起来**，并将开发精力 **集中在 configs 与 logic 本身**。

本项目基于 Problab 引擎构建。

---

## 1 分钟检查清单

1. 获取项目  
   `git clone https://github.com/zintix-labs/problab-scaffold.git`

2. 运行  
   `make run`（或 `go run ./cmd/run`）

3. 私有化  
   - 修改 `go.mod` 中的 module 名称  
   - 全局替换 import 路径  
   - （建议）重新初始化 Git 仓库

4. 完成  
   从 `internal/configs/` 与 `internal/logic/` 开始新增游戏。

---

## 为什么选择这个 Scaffold

### A. 极速创建 Problab 项目

目标：**1 分钟运行成功，3 分钟进入商业开发状态**

### 步骤 1：获取项目

```bash
git clone https://github.com/zintix-labs/problab-scaffold.git
cd problab-scaffold
```

### 步骤 2：运行

支持 Makefile（推荐）：

```bash
make run
make dev
make svr
```

不使用 Makefile：

```bash
go run ./cmd/run
go run ./cmd/dev
go run ./cmd/svr
```

---

## 步骤 3：私有化（强烈建议）

### 1. 修改 Go module 名称

编辑 `go.mod`，将 module 名称修改为你自己的项目名：

```go
module <your project module name>
```

> module 名称 **不建议** 对应 GitHub 仓库地址。  
> 实际商业项目通常不会直接暴露在 GitHub。

---

### 2. 全局替换 import 路径

**VS Code（推荐）**

搜索：
```
github.com/zintix-labs/problab-scaffold
```

替换为：
```
<your project module name>
```

**命令行方式**

```bash
grep -rl "github.com/zintix-labs/problab-scaffold" . \
  | xargs sed -i '' 's|github.com/zintix-labs/problab-scaffold|<your project module name>|g'
```

---

### 3. 重新初始化 Git 历史（推荐）

```bash
rm -rf .git
git init
```

完成后，该项目将成为一个**完全私有、可直接用于商业开发的代码基线**。

---

## 新增一个游戏

新增游戏只需要 **两个步骤**：

1. 在 `internal/configs/` 中新增一个配置文件  
   （复制示例 YAML 并修改即可）

2. 在 `internal/logic/` 中新增一个逻辑实现  
   （复制示例游戏逻辑并调整规则）

完成后，新游戏即可立即运行，并具备 **生产可用性**。

---

## 选择游戏（分段游戏）

`demo_2`（`demo_choice`）与 `demo_3`（`demo_entry`）示范会停在选择点的游戏：不是整局
结束的段会在 `ChooseNext` 列出入口，下一段执行所选位置上的入口。

- `internal/logic/demo_2.go` 通过已注册的 checkpoint（`dto.RegisterCheckpoint`）把进入
  状态交给下一段。
- `internal/logic/demo_3.go` 不读 checkpoint：每一段只取决于入口值与 seed，因此每个入口
  都可以配置自己的 Optimal 结果对。

模拟这类游戏需要指定选择方式。`make run` 固定选 `first`；其他方式请直接运行 CLI：

```bash
make run game=2
go run ./cmd/run -game 2 -choose last
```

---

## 玩家体验分析

`make exp` 以 embed 的 `cmd/exp/exp_cfg.yaml` 分析有限玩家体验；`config=<path>`
改用其他 YAML：

```bash
make exp
make exp config=cmd/exp/examples/game_3.yaml
```

策略是按游戏注册在 `internal/strategies/` 的 Factory（`game_<gid>.go`，并列在
`reg.go`）。内置的 `builtin.Monotone`（`sdk/strategy/builtin`）固定押注；它的
`Choose` 在选择点挑选位置（`ChooseFirst`、`ChooseLast` 或 `ChooseIter`）。YAML 可以用
`player.max_segments` 限制一局的段数（默认 10）。

---

## Optimal 工作流程（Problab v0.17.0）

Scaffold 将普通文件目录 `internal/optimal` 作为唯一的 Optimal 结果库根目录。
运行时入口（`make run`、`make dev` 与 `make svr`）调用 `engine.New()`，并可从该
目录读取已发布的 gacha 结果对。`make opt` 则调用 `engine.NewForOptimizer()`：它使用
同一套由项目注入的 PRNG Factory、embed 配置与 Logic Registry，但在收集原始游戏结果
时不会预加载已经发布的结果对。

每个 BetMode 发布为一组 gacha 结果对：一个 `gacha.json.zst` 与一个 `seed_bank.bin`。

### Optimizer 由配置文件驱动，不再使用命令行参数

`cmd/opt` 是 Problab 的 **LP / 数学意图 Optimizer（v2）**。它**不接受任何
命令行参数**——`make opt` 直接运行，传入 flag 会被拒绝。每一次运行都由 embed 进
二进制的 `cmd/opt/opt_cfg.yaml`（`//go:embed`）唯一定义，这是运行的唯一可审计
来源。要改变运行行为，就编辑该文件并重新编译。第一次 Ctrl+C 会在当前求解结束后
停止；第二次 Ctrl+C 立即退出。

`opt_cfg.yaml` 分为三段，所有权刻意分离：

| 段落 | 归属 | 内容 |
| --- | --- | --- |
| `plans:` | 工程 | 路由与执行：目标 `game` + `bet_modes`（以及续段可选的 `input`）、`seed`、收集参数 `workers` / `batch_size` / `max_spins`、输出 `format` 与 `directory`。按声明顺序执行；每个 plan 只负责一个 BetMode。 |
| `intents:` | 设计师 | 数学合约：class 权重、闭区间 `win_range`、收集 `tags`、精确 `exp` 与 `median`、分桶的 `main_experience` group 及软性 `prefer` 分布、整体 `cv` 范围、可选的碰撞 `risk` 策略。 |
| `engine_options:` | 工程 | 仅数值控制：feasibility / optimality / quantile 容差与 bisection 迭代预算。infeasible 时不会自动放宽。 |

`intents[].collect.tags` 引用的收集 tag 由 Scaffold 在
`internal/tags/` 中拥有。`Catalog` 将每个 `spec.GID` 绑定到其
tag 判定函数集合；每个游戏新增一个文件并在此注册。

RGS 输出（`rgs-collected`、`rgs-optimized`）以 `cmd/opt/main.go` 的 `useConverter`
编码每一笔记录。默认输出完整 DTO，包含 PRNG 状态；`internal/converter` 另外提供会去掉
这些状态的 `WithoutSnapshotConverter`。

### 以 `game_id: 1901` 为例

1. 开发与收集阶段关闭 Optimal：

   ```yaml
   optimal_setting:
     use_optimal: false
   ```

2. 在 `internal/tags/` 注册该游戏的收集 tag
   （复制 `demo_0_tags.go`，并把新的 `spec.GID` 加入 `Catalog`）。

3. 在 `cmd/opt/opt_cfg.yaml` 为每个 BetMode 新增一个 plan（设置 `target.game: 1901`、
   `target.bet_modes`、`intent`、`seed`、收集预算），然后运行：

   ```bash
   make opt
   ```

4. 检查生成结果：

   ```text
   build/optimizer/
   |-- gacha/game_1901/mode_<mode>/<id>/
   |   |-- gacha.json.zst
   |   `-- seed_bank.bin
   `-- distribution_mode_<mode>.csv
   ```

   每个 BetMode 独立发布。结果对先写入暂存目录，经正式加载器检查后再改名就位；
   `<id>` 由内容决定，已发布的目录不会被改写。

5. 复制审核通过的目录，并附上产生它的 `opt_cfg.yaml`：

   ```text
   build/optimizer/gacha/game_1901/mode_0/<id>/
       -> internal/optimal/gacha/game_1901/mode_0/<id>/
   ```

6. 在游戏配置中引用这组结果对：

   ```yaml
   optimal_setting:
     use_optimal: true
     gachas:
       - gacha/game_1901/mode_0/<id>/gacha.json.zst
     seed_bank:
       - gacha/game_1901/mode_0/<id>/seed_bank.bin
   ```

   多个 BetMode 时，按 BetMode 顺序在两个列表各加一条路径。两个列表长度都必须与
   `bet_units` 相同：索引 `0` 对应 `bet_units[0]`，索引 `1` 对应 `bet_units[1]`，以此类推。

`gachas`/`seed_bank` 是保留的相容写法。任何游戏也都可以改在 `optimal_setting.modes`
按 BetMode 列出结果对：`start` 为初始段；选择游戏另以 `after` 按入口值列出续段（续段以
`plans[].target.input` 收集）。详见 `internal/optimal/README.md`。

### PRNG 与 snapshot 相容性

Optimal Seed Bank 保存的是序列化后的 PRNG snapshot，因此 Optimizer 与运行时必须
使用相容的 PRNG 实现与 snapshot 长度。默认使用 ChaCha20。既有 PCG64 结果对只有在
项目明确注入相容的 PCG64 Factory 时才能继续运行；修改目录或 YAML 路径不会把 PCG64
结果对转换为 ChaCha20 结果对。更换 PRNG 实现后必须重新生成输出。

### 私有输出与部署

Git 默认忽略 `internal/optimal/**` 下生成的文件。Optimal 输出可能包含具有商业价值的
数学数据，而 Seed Bank 也可能大到不适合普通 Git 托管。正式环境建议将审核通过的结果对
存放在私有 Artifact 服务，并在部署时复制到 `internal/optimal`，或以只读方式挂载到该目录。

如果组织明确要在私有仓库中管理结果对，请从 `.gitignore` 移除对应的
`internal/optimal` 规则，并使用 Git LFS 或其他大文件存储方案。GitHub 私有仓库仍有
一般的单文件大小限制。`.zst` 还会被全局 `*.zst` 规则忽略；要追踪这些输出，也必须
移除或覆盖该规则。切勿将具有商业价值的 Optimal 数据提交到公开仓库。

---

## 架构说明

- 配置文件通过 `internal/configs/` 进行 embed
- Config FS 采用 **扁平结构**
  - 仅支持目录内的 `*.yaml`
  - 不支持子目录
- 游戏逻辑通过 `internal/logic/` 中的 `init()` 自动注册
- `pkg/engine/problab.go` 拥有共用的 PRNG、Config 与 Logic 接线。运行时命令调用
  `engine.New()`；`make opt` 调用 `engine.NewForOptimizer()`，在不加载已发布结果对
  的情况下完成收集。
- 已发布的 gacha 结果对存放在
  `internal/optimal/gacha/game_<gid>/mode_<mode>/<id>/`。
- Optimal 文件使用普通文件而非 embed，使 Seed Bank 在支持的平台上可以只读 mmap。
- Optimizer 运行只由 embed 的 `cmd/opt/opt_cfg.yaml` 定义，二进制不接受任何
  flag。收集 tag 位于 `internal/tags/`，`Catalog` 将每个
  `spec.GID` 映射到其 tag 判定函数。
- 玩家策略位于 `internal/strategies/`，`Catalog` 将每个 `spec.GID` 映射到其策略
  Factory。
- 同一个 Problab 对每组结果对只加载／映射一次，并由全部 Machine、Pool
  与 Simulator worker 共享。
- 各命令会在结束时关闭 Problab，正确释放 mmap。

这些限制是**刻意设计的约束**，  
用于保持系统行为可预测、结构清晰、易于维护。

---

## 常用命令

- `make run`：运行模拟器（默认 `game=0`；选择游戏固定选 `first`）
- `make exp`：分析有限玩家体验（`config=<path>` 改用其他 YAML）
- `make opt`：运行 Optimizer（embed 的 `cmd/opt/opt_cfg.yaml`）
- `make dev`：启动 Dev Web 面板
- `make svr`：启动 HTTP Server
- `make help`：查看全部命令

---

## 环境要求

- Go 1.27.3 或以上
- Problab v0.17.0
- `make`（非必须，但推荐）

---

## License

Apache-2.0
