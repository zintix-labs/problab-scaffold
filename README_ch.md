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

## Optimal 工作流程（Problab v0.7.0）

Scaffold 将普通文件目录 `internal/optimal` 作为唯一的 Optimal 结果库根目录。
在支持的平台上，Problab 会对其中确认完成的 Artifact 进行只读 mmap。
`make run`、`make dev`、`make svr` 与 `make opt` 仍统一通过 `engine.New()`
建立引擎，不需要分别修改命令入口。

### Optimizer 由配置文件驱动，不再使用命令行参数

`cmd/opt` 是 Problab v0.7.0 的 **LP / 数学意图 Optimizer（v2）**。它**不接受任何
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

4. 在 `build/optimizer/artifact_v1/game_1901/` 检查生成的 Artifact bundle；
   旧版 gacha 交换格式写在 `build/optimizer/gacha/game_1901/`，每个 mode 的分布
   CSV 会写在输出旁边。未完成的运行保留在 `.pending` 暂存目录中，对运行时不可见，
   直到所有 BetMode 齐全并原子发布。

5. 将确认完成的整个 `artifact_v1/game_1901/` 目录复制到
   `internal/optimal/game_1901/`。

6. 开启最终结果：

   ```yaml
   optimal_setting:
     use_optimal: true
     artifact: game_1901/manifest.json
   ```

Manifest 会自动解析 probability、aliases 与 SeedBank，游戏配置不需要分别填写
这些文件路径。

Docker image 会将 `internal/optimal` 复制到 `/app/internal/optimal`，让本地与
容器使用相同目录结构。正式部署也可以将确认后的结果库以只读方式挂载到该位置。

---

## 架构说明

- 配置文件通过 `internal/configs/` 进行 embed
- Config FS 采用 **扁平结构**
  - 仅支持目录内的 `*.yaml`
  - 不支持子目录
- 游戏逻辑通过 `internal/logic/` 中的 `init()` 自动注册
- `pkg/engine/problab.go` 只有一个私有的
  `WithOptimalDir("internal/optimal")` 接线点；各命令仍只调用 `engine.New()`。
- 确认后的 Optimal 存放于 `internal/optimal/game_<gid>/`，使用普通文件而非
  embed，使运行时可以只读 mmap。
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
- Problab v0.7.0
- `make`（非必须，但推荐）

---

## License

Apache-2.0
