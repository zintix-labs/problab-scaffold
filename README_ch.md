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

## Optimal 工作流程（Problab v0.8.0）

Scaffold 将普通文件目录 `internal/optimal` 作为唯一的 Optimal 结果库根目录。
运行时入口（`make run`、`make dev` 与 `make svr`）调用 `engine.New()`，并可从该
目录读取已发布的 bundle。`make opt` 则调用 `engine.NewForOptimizer()`：它使用同一套
由项目注入的 PRNG Factory、embed 配置与 Logic Registry，但在收集原始游戏结果时
不会预加载已经发布的 bundle。

Problab v0.8.0 支持两种互斥的运行时格式：

- **Artifact v1（推荐）：**由 manifest、二进制 probability、alias 与 Seed Bank
  文件组成。Manifest 会记录文件哈希与 PRNG snapshot 格式。
- **Legacy gacha：**每个 BetMode 各有一个 `gacha_<mode>.json.zst` 与一个
  `seed_bank_<mode>.bin`。既有集成仍可继续使用这种格式。

### Optimizer 由配置文件驱动，不再使用命令行参数

`cmd/opt` 是 Problab v0.8.0 的 **LP / 数学意图 Optimizer（v2）**。它**不接受任何
命令行参数**——`make opt` 直接运行，传入 flag 会被拒绝。每一次运行都由 embed 进
二进制的 `cmd/opt/opt_cfg.yaml`（`//go:embed`）唯一定义，这是运行的唯一可审计
来源。要改变运行行为，就编辑该文件并重新编译。

`opt_cfg.yaml` 分为三段，所有权刻意分离：

| 段落 | 归属 | 内容 |
| --- | --- | --- |
| `plans:` | 工程 | 路由与执行：目标 `game` + `bet_modes`、`seed`、收集参数 `workers` / `batch_size` / `max_spins`、输出 `format` 与 `directory`。按声明顺序执行；每个 BetMode 一个 plan，或在 `bet_modes` 中列多个。 |
| `intents:` | 设计师 | 数学合约：class 权重、闭区间 `win_range`、收集 `tags`、精确 `exp` 与 `median`、分桶的 `main_experience` group 及软性 `prefer` 分布、整体 `cv` 范围、可选的碰撞 `risk` 策略。 |
| `engine_options:` | 工程 | 仅数值控制：feasibility / optimality / quantile 容差与 bisection 迭代预算。infeasible 时不会自动放宽。 |

`intents[].collect.tags` 引用的收集 tag 由 Scaffold 在
`internal/logic/optimizer_tags/` 中拥有。`GameTags` 将每个 `spec.GID` 绑定到其
tag 判定函数集合；每个游戏新增一个文件并在此注册。

### 以 `game_id: 1901` 为例

1. 开发与收集阶段关闭 Optimal：

   ```yaml
   optimal_setting:
     use_optimal: false
   ```

2. 在 `internal/logic/optimizer_tags/` 中注册该游戏的收集 tag
   （复制 `demo_0_tags.go`，把新的 `spec.GID` 加入 `GameTags`）。

3. 在 `cmd/opt/opt_cfg.yaml` 中为每个 BetMode 添加一个 plan
   （设定 `target.game: 1901`、`target.bet_modes`、`intent`、`seed`、收集预算），
   然后运行：

   ```bash
   make opt
   ```

4. 检查生成结果：

   ```text
   build/optimizer/
   |-- artifact_v1/game_1901/
   |-- gacha/game_1901/
   `-- distribution_mode_<mode>.csv
   ```

   未完成的运行会保留在 `.pending` 暂存目录中。在所有 BetMode 齐全并完成原子发布
   之前，运行时不会看到这些暂存结果。

5. 按照下方其中一种方式发布 Artifact v1 或 legacy gacha。同一游戏不可同时配置
   两种格式。

### 推荐方式：发布 Artifact v1

复制整个审核通过的目录，不要只挑选部分文件：

```text
build/optimizer/artifact_v1/game_1901/
    -> internal/optimal/artifact_v1/game_1901/
```

安装后的目录包含一个 manifest，以及所有 BetMode 对应的文件：

```text
internal/optimal/artifact_v1/game_1901/
|-- manifest.json
|-- mode_0.json
|-- prob_0.bin
|-- aliases_0.bin
|-- seed_bank_0.bin
`-- ... 其他 BetMode 的文件
```

在游戏配置中启用最终结果：

```yaml
optimal_setting:
  use_optimal: true
  artifact: artifact_v1/game_1901/manifest.json
```

Manifest 会自动解析 probability、aliases 与 SeedBank，游戏配置不需要分别填写
这些文件路径。

### 既有集成：继续发布 legacy gacha

既有用户可以继续使用 `.json.zst` 加 Seed Bank 的目录格式。复制整个审核通过的
legacy 目录：

```text
build/optimizer/gacha/game_1901/
    -> internal/optimal/gacha/game_1901/
```

只有一个 BetMode 时，安装后的文件如下：

```text
internal/optimal/gacha/game_1901/
|-- gacha_0.json.zst
`-- seed_bank_0.bin
```

在配置中明确引用两个文件：

```yaml
optimal_setting:
  use_optimal: true
  gachas:
    - gacha/game_1901/gacha_0.json.zst
  seed_bank:
    - gacha/game_1901/seed_bank_0.bin
```

如果有多个 BetMode，请按照 BetMode 顺序分别在两个列表中加入路径。两个列表的长度
都必须与 `bet_units` 相同：index `0` 对应 `bet_units[0]`，index `1` 对应
`bet_units[1]`，以此类推。配置 `gachas` 与 `seed_bank` 时，不可同时设置
`artifact`。

### PRNG 与 snapshot 相容性

Optimal Seed Bank 保存的是序列化后的 PRNG snapshot，因此 Optimizer 与运行时必须
使用相容的 PRNG 实现与 snapshot 长度。v0.8.0 默认使用 ChaCha20。既有 PCG64
bundle 只有在项目明确注入相容的 PCG64 Factory 时才能继续运行；修改目录或 YAML
路径不会把 PCG64 bundle 转换为 ChaCha20 bundle。更换 PRNG 实现后，无论使用
Artifact v1 或 legacy gacha，都必须重新生成输出。

### 私有 Artifact 与部署

Git 默认忽略 `internal/optimal/artifact_v1/**` 与
`internal/optimal/gacha/**`。Optimal 输出可能包含具有商业价值的数学数据，而 Seed
Bank 也可能大到不适合普通 Git 托管。正式环境建议将审核通过的 bundle 存放在私有
Artifact 服务，并在部署时复制到 `internal/optimal`，或以只读方式挂载到该目录。

如果组织明确要在私有仓库中管理 bundle，请从 `.gitignore` 移除对应的
`internal/optimal` 规则，并使用 Git LFS 或其他大文件存储方案。GitHub 私有仓库仍有
一般的单文件大小限制。Legacy `.zst` 还会被全局 `*.zst` 规则忽略；要追踪 legacy
输出，也必须移除或覆盖该规则。切勿将具有商业价值的 Optimal 数据提交到公开仓库。

---

## 架构说明

- 配置文件通过 `internal/configs/` 进行 embed
- Config FS 采用 **扁平结构**
  - 仅支持目录内的 `*.yaml`
  - 不支持子目录
- 游戏逻辑通过 `internal/logic/` 中的 `init()` 自动注册
- `pkg/engine/problab.go` 拥有共用的 PRNG、Config 与 Logic 接线。运行时命令调用
  `engine.New()`；`make opt` 调用 `engine.NewForOptimizer()`，在不加载已发布 Artifact
  的情况下完成收集。
- 最终 Artifact v1 bundle 存放在
  `internal/optimal/artifact_v1/game_<gid>/`；legacy bundle 存放在
  `internal/optimal/gacha/game_<gid>/`。
- Artifact 使用普通文件而非 embed，使 Artifact v1 在支持的平台上可以只读 mmap。
- Optimizer 运行只由 embed 的 `cmd/opt/opt_cfg.yaml` 定义，二进制不接受任何
  flag。收集 tag 位于 `internal/logic/optimizer_tags/`，`GameTags` 将每个
  `spec.GID` 映射到其 tag 判定函数。
- 同一个 Problab 对每份 Artifact 只加载／映射一次，并由全部 Machine、Pool
  与 Simulator worker 共享。
- 各命令会在结束时关闭 Problab，正确释放 mmap。

这些限制是**刻意设计的约束**，  
用于保持系统行为可预测、结构清晰、易于维护。

---

## 常用命令

- `make run`：运行模拟器
- `make dev`：启动 Dev Web 面板
- `make svr`：启动 HTTP Server
- `make help`：查看全部命令

---

## 环境要求

- Go 1.25 或以上
- Problab v0.8.0
- `make`（非必须，但推荐）

---

## License

Apache-2.0
