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

## 架构说明

- 配置文件通过 `internal/configs/` 进行 embed
- Config FS 采用 **扁平结构**
  - 仅支持目录内的 `*.yaml`
  - 不支持子目录
- 游戏逻辑通过 `internal/logic/` 中的 `init()` 自动注册

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
- `make`（非必须，但推荐）

---

## License

Apache-2.0