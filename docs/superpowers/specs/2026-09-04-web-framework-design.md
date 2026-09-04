# Web 框架搭建设计（fx + gin + viper + zap + gorm）

日期：2026-09-04
状态：已获用户批准（方案 B · 轻量分层）

## 1. 背景与目标

现有 `general-agent` 是半成品骨架：`cmd/main.go` 为空（fx 从未启动）、无 viper/zap、
分层断层（Service 无方法、DAO 有 bug）、路由注册为死代码。本次在保留既有依赖选型的
前提下重搭框架，以 user 模块全链路贯穿验证。

**成功标准：**
1. `make run` 启动服务，5 个 user REST 接口可用（连真库）
2. `make test` 全绿，且不依赖任何外部服务
3. 新增一个业务模块 = 新建 `app/<domain>/` 目录 + `cmd/main.go` 加一行
4. API 响应不泄漏 entity 中的敏感字段（密码等）

## 2. 已确认的决策

| 维度 | 决定 |
|---|---|
| 分包方式 | 按业务域（feature-first），模块内含全部分层 |
| 抽象深度 | 仅 repository 接口化；service/controller 为具体类型 |
| 路由注册 | 写在各模块 `module.go` 的 `fx.Invoke` 中 |
| 配置组织 | `config/config.yaml`（公共）+ `config/{env}.yaml`（覆盖），`-env` 默认 dev |
| 错误处理 | 统一错误码 + 全局 Recovery + errs 辅助函数 |
| 代码生成器 | 保留，迁至 `tools/gen`，DSN 改读主配置 |
| 验证方式 | 单测（sqlite 内存库/手写 fake/httptest）+ 真库冒烟 |

## 3. 目标目录结构

```
general-agent/
├── cmd/main.go              # 唯一入口：fx.New(各模块...).Run()
├── config/
│   ├── config.yaml          # 公共配置
│   └── dev.yaml             # 环境覆盖（prod.yaml 同理）
├── app/
│   └── user/                # 业务模块（见第 6 节）
├── internal/
│   ├── config/              # viper 加载（见 4.1）
│   ├── logger/              # zap（见 4.2）
│   ├── database/            # gorm（见 4.3）
│   ├── server/              # gin + fx lifecycle（见 4.4）
│   └── errs/                # 错误码 + 响应 + Recovery（见 5）
├── tools/gen/               # 实体生成器（见 8）
├── docs/                    # swagger（保留，swag init 重新生成）
├── Makefile
└── go.mod                   # module general-agent
```

**删除：** `common/`（ResultResp 移入 errs）、旧 `config/`、`app/dao|services|controller|entity/`、
`cmd/generate/`（迁至 tools/gen，删除重复生成物 `cmd/generate/app/entity/user.gen.go`）。

## 4. 基础设施层（internal/）

### 4.1 config —— viper

- `Config` 结构体：`HTTP{Addr, BasePath}`、`MySQL{Host, Port, Username, Password, Database, MaxIdleConns, MaxOpenConns, ConnMaxLifetime}`、`Log{Level, Encoding}`。
- 加载：`flag -env`（默认 `dev`）→ viper 读 `config/config.yaml` → `MergeInConfig` 读 `config/{env}.yaml` 覆盖 → `Unmarshal` 到 `Config`。
- `fx.Provide` 单例；文件缺失或反序列化失败返回 error（启动即失败，不静默降级）。

### 4.2 logger —— zap

- `NewLogger(cfg)`：按 `log.level` 与 `log.encoding`（console/json）创建，输出 stdout。
- fx lifecycle `OnStop` 时 `Sync()`。
- 全应用共享一个 `*zap.Logger`（gin 中间件、gorm、业务层复用）。

### 4.3 database —— gorm

- `NewDB(cfg, log) (*gorm.DB, error)`：MySQL DSN（`charset=utf8mb4&loc=Asia%2FShanghai&parseTime=true`）。
- 连接池：`SetMaxIdleConns` / `SetMaxOpenConns` / `SetConnMaxLifetime`（修正原代码 `SetConnMaxIdleTime` 与 `ConnMaxLifetime` 语义混用的 bug）。
- 实现 `gorm logger.Interface` 桥接 zap：慢 SQL（阈值 1s）warn、错误 error，替代原标准库 log。
- 保留 `TranslateError: true`。

### 4.4 server —— gin

- 提供 `*gin.Engine`：`gin.New()` + zap 访问日志中间件 + errs.Recovery + swagger 路由（`/swagger/*any`）。
- 提供 `*gin.RouterGroup`：`Group(cfg.HTTP.BasePath)`，供业务模块注册路由。
- fx lifecycle：`OnStart` 异步 `ListenAndServe`（启动失败须令 app 退出，不能只打日志）；`OnStop` 优雅关停（`Shutdown` + 5s 超时）。

## 5. 统一错误处理（internal/errs）

- `BizError{Code int, Msg string}` 实现 `error`；预定义：`ErrInternal`（50000）、`ErrNotFound`（10001）、`ErrUserExists`（10002）等，可按模块扩展。
- `Response{Code int, Msg string, Data any}` 替代旧 `ResultResp`。
- 辅助函数（controller 调用）：
  - `errs.OK(ctx, data)` → HTTP 200 + `{code:0, msg:"success", data}`
  - `errs.Fail(ctx, err)` →
    - `*BizError`：HTTP 200 + `{code, msg}`（国内惯例，前端按 code 分支）
    - 其他 error：zap 记录 + HTTP 500 + `{code:50000, msg:"内部错误"}`（不泄漏细节）
- Recovery 中间件：panic → zap 记录堆栈 → 同上 500 响应。

## 6. user 业务模块（全链路模板）

```
app/user/
├── user.gen.go     # gorm gen 生成实体（由现有 app/entity/user.go 迁移，禁止手改）
├── repository.go   # Repository 接口 + impl
├── service.go      # Service 具体类型
├── controller.go   # Controller 具体类型 + swag 注释
├── dto.go          # 请求/响应 DTO
└── module.go       # fx.Module("user", Provide(NewRepository, NewService, NewController), Invoke(registerRoutes))
```

依赖链：`controller → *Service → Repository 接口 → *gorm.DB`。

### 6.1 API（挂在 basePath `/api/v1/general-agent` 下）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | /user/:id | 按 ID 查询 |
| GET | /user?pageNum=&pageSize= | 分页列表 |
| POST | /user | 创建 |
| PUT | /user/:id | 更新 |
| DELETE | /user/:id | 删除 |

### 6.2 DTO 与实体分离（硬需求）

现有 entity 含 `password`/`password_sm3`/`mds_pass` 等 6 个敏感字段且 JSON tag 全导出，
直接返回会泄漏密码哈希：
- 入参：`CreateUserReq`、`UpdateUserReq`、`ListUserReq`（pageNum/pageSize，默认 1/10）
- 出参：`UserResp`（仅 id、username、name、roleID、status、createTime）；
  列表接口返回 `PageResp{List []UserResp, Total int64}`
- controller 内完成 entity ↔ DTO 转换

### 6.3 各层职责

- **repository**：纯数据访问，`GetByID` / `List` / `Create` / `Update` / `Delete`；
  查询用 `First` 并原样返回 error（修正原 `Find(nil)` bug）。
- **service**：业务规则与错误转换——`errors.Is(err, gorm.ErrRecordNotFound)` → `errs.ErrNotFound`；
  用户名重复 → `errs.ErrUserExists`。错误语义的翻译只发生在 service 层。
- **controller**：参数绑定校验（gin binding）、DTO 转换、调 `errs.OK`/`errs.Fail`。

### 6.4 数据流（GetByID 示例）

```
GET /api/v1/general-agent/user/:id
→ controller.GetById（ShouldBindUri）
→ service.GetByID（ErrRecordNotFound → errs.ErrNotFound）
→ repository.GetByID（db.First）
→ MySQL
→ controller: errs.OK(ctx, toUserResp(user))
```

## 7. 组装（cmd/main.go）

```go
fx.New(
    config.Module,
    logger.Module,
    database.Module,
    server.Module,
    user.Module,
).Run()
```

依赖顺序由 fx 自动解析；`internal/` 可被同 module 的 `app/`、`tools/` 引用。

## 8. tools/gen 实体生成器

- `targets := map[string]string{"user": "./app/user"}`（表名 → 输出目录）。
- DSN 通过 `internal/config` 加载（dev 环境，删除硬编码 DSN）。
- 每个表生成 `<table>.gen.go`（包名 = 目录名）；新增表 = 加一行映射。

## 9. 测试策略

| 层 | 方式 | 外部依赖 |
|---|---|---|
| repository | sqlite 内存库（`github.com/glebarez/sqlite` 纯 Go 驱动）+ AutoMigrate 建表，覆盖 5 个方法 | 无 |
| service | 手写 fake 实现 Repository 接口，测错误转换与业务规则 | 无 |
| 全链路 | httptest 手动组装 repo→service→controller→gin，HTTP 断言 JSON | 无 |
| 冒烟 | `make run` + curl 连真库（dev.yaml 由用户填 DSN） | MySQL |

不引 gomock（接口少，手写 fake 更直观）。

## 10. 新增依赖

- 直接依赖转为：`github.com/spf13/viper`（新增）、`go.uber.org/fx`、`go.uber.org/zap`、`github.com/gin-gonic/gin`、`github.com/swaggo/gin-swagger`
- 测试依赖：`github.com/glebarez/sqlite`
- 保留：`gorm.io/gorm`、`gorm.io/driver/mysql`、`gorm.io/gen`（tools 用）

## 11. 非目标（本次不做，留作扩展点）

- 环境变量覆盖配置（Spring relaxed binding 风格）——viper 已支持，需要时加 `AutomaticEnv`
- 认证/鉴权中间件、限流、熔断
- 事务管理器（`db.Transaction` 在 service 内按需直接用）
- 多数据源、读写分离
- Redis/消息队列等额外基础设施

## 12. 现有问题的处置清单

| 现有问题 | 处置 |
|---|---|
| `cmd/main.go` 空 | 重写为 fx 组装入口 |
| 无 viper/zap | 按第 4 节补齐 |
| `dao.GetById` `Find(nil)` bug | repository 重写为 `First` + 显式错误判断 |
| `registryRouter` 死代码 | 路由统一注册进 module.go 的 fx.Invoke |
| `common/dao.go` 无用基类 | 删除 |
| `TernaryOperation` 硬编码 | 删除（日志级别走配置） |
| 实体两份重复 | 只保留 `app/user/user.gen.go` |
| 包名不一致（services/controller/dao） | 统一为 `app/<domain>/` 单包 |
