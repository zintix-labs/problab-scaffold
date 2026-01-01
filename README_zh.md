# Problab Scaffold

*由 **Zintix Labs** 維護*

<p align="right">
  <a href="./README.md">English</a> | <b>中文</b>
</p>

這是一個可直接投入生產的 Problab scaffold，專為 slot game 研發而設計。目標是**讓你快速跑起來**，並且**把開發集中在 configs 與 logic**。

本專案使用 Problab 架構。

## 1 分鐘檢查表

1. 取得專案  
   `git clone https://github.com/zintix-labs/problab-scaffold.git`

2. 執行  
   `make run`（或 `go run ./cmd/run`）

3. 私有化  
   - 修改 `go.mod` module 名稱  
   - 全域替換 import paths  
   - （建議）`git init`

4. 完成  
   從 `internal/configs/` 與 `internal/logic/` 開始新增遊戲。

---

## 為什麼是這個 scaffold

### A. 超快速建立 Problab 專案

目標：**1 分鐘執行成功，3 分鐘商業開發就緒**

**步驟 1：取得專案**

```bash
git clone https://github.com/zintix-labs/problab-scaffold.git
cd problab-scaffold
```

**步驟 2：執行**

支援 Makefile（建議）：

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

### 步驟 3：私有化（建議）

#### 1. 修改 Go module 名稱

編輯 `go.mod`，將 module 名稱改為你自己的專案名稱：

```go
module <your project module name>
```

> module 名稱不需要對應 GitHub。

#### 2. 全域替換 import paths

**VS Code（建議）**  
搜尋：
```
github.com/zintix-labs/problab-scaffold
```
取代為：
```
<your project module name>
```

**命令列**

```bash
grep -rl "github.com/zintix-labs/problab-scaffold" . \
  | xargs sed -i '' 's|github.com/zintix-labs/problab-scaffold|<your project module name>|g'
```

#### 3. 重新初始化 Git 歷史（建議）

```bash
rm -rf .git
git init
```

完成後，專案即為完全私有、可商業開發的基底。

---

## 新增一個遊戲

新增遊戲只需要兩個步驟：

1. 在 `internal/configs/` 新增一個設定檔  
   （複製範例 YAML 後修改即可）

2. 在 `internal/logic/` 新增一個邏輯實作  
   （複製範例遊戲並調整規則）

完成後，新的遊戲即可直接執行，並具備生產可用性。

---

## 架構說明

- 設定檔由 `internal/configs/` embed
- Config FS 採 **平坦結構**  
  （僅支援資料夾內的 `*.yaml`，不使用子目錄）
- 遊戲邏輯透過 `internal/logic/` 中的 `init()` 註冊

這些限制是刻意設計，用來保持系統可預期與易於維護。

---

## 常用指令

- `make run`：執行模擬器
- `make dev`：啟動 Dev Web 面板
- `make svr`：啟動 HTTP Server
- `make help`：顯示所有指令

---

## 環境需求

- Go 1.25 以上
- `make`（非必要，但建議）

---

## License

Apache-2.0
