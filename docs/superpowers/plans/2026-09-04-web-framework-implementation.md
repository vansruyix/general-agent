# Web 框架搭建（fx + gin + viper + zap + gorm）实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 基于 fx/gin/viper/zap/gorm 搭建类 Spring Boot 分层 web 框架，以 user 模块全链路贯穿验证。

**Architecture:** 业务按 feature 分包（app/user/），基础设施收拢 internal/ 下。仅 repository 接口化，service/controller 为具体类型。路由注册在 module.go 的 fx.Invoke 中。自底向上装配：config → logger → database → errs → server → user → main。

**Tech Stack:** Go 1.26, go.uber.org/fx, github.com/gin-gonic/gin, github.com/spf13/viper, go.uber.org/zap, gorm.io/gorm, gorm.io/driver/mysql, github.com/glebarez/sqlite（测试）, github.com/gin-contrib/sse, github.com/swaggo/gin-swagger

**Spec:** `docs/superpowers/specs/2026-09-04-web-framework-design.md`

## Global Constraints

- module path: `general-agent`
- Go 1.26.8
- 配置：公共配置 + 环境覆盖，-env 默认 dev
- 错误处理：统一错误码，BizError 返回 HTTP 200，未知错误返回 HTTP 500
- API 响应不泄漏 entity 中的敏感字段（密码等）
- 仅 repository 接口化；service/controller 为具体类型
- 测试：单测不依赖外部服务（sqlite 内存库/手写 fake/httptest）
- 不做 git commit（用户自行 review 后提交）

---

### Task 1: 清理旧代码与安装依赖

**Files:**
- Delete: `common/`, `config/`, `app/dao/`, `app/services/`, `app/controller/`, `app/entity/`, `cmd/generate/`
- Modify: `go.mod`

**Interfaces:**
- Consumes: 无
- Produces: 干净的项目目录，所有新依赖就绪

- [ ] **Step 1: 删除旧目录**

```bash
rm -rf common/ config/ app/dao/ app/services/ app/controller/ app/entity/ cmd/generate/
```

- [ ] **Step 2: 安装新依赖**

```bash
cd /venus/agent/general-agent
go get github.com/spf13/viper@latest
go get github.com/glebarez/sqlite@latest
go mod tidy
```

- [ ] **Step 3: 验证依赖项**

```bash
go build ./...
```
Expected: 应该失败（因为 cmd/main.go 仍然引用了已删除的包），但下载/编译依赖本身不应报错。如果报 "package not found" 是预期的——下一步会重写 main.go。

---

### Task 2: internal/config —— viper 配置加载

**Files:**
- Create: `internal/config/config.go`
- Create: `config/config.yaml`
- Create: `config/dev.yaml`

**Interfaces:**
- Consumes: 无
- Produces: `*config.Config`, `config.Module` (fx.Option)
- Types: `Config{HTTP{Addr, BasePath}, MySQL{Host, Port, Username, Password, Database, MaxIdleConns, MaxOpenConns, ConnMaxLifetime}, Log{Level, Encoding}}`

- [ ] **Step 1: 写配置文件**

`config/config.yaml`:
```yaml
http:
  addr: ":8080"
  basePath: "/api/v1/general-agent"

log:
  level: "info"
  encoding: "console"

mysql:
  host: "127.0.0.1"
  port: 3306
  username: "root"
  password: ""
  database: "general_agent"
  maxIdleConns: 10
  maxOpenConns: 100
  connMaxLifetime: 3600
```

`config/dev.yaml`:
```yaml
log:
  level: "debug"

mysql:
  host: "127.0.0.1"
  port: 3306
  username: "root"
  password: "root123456"
  database: "general_agent"
```

- [ ] **Step 2: 写 Config 结构体与 viper 加载**

`internal/config/config.go`:
```go
package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Config struct {
	HTTP  HTTP  `mapstructure:"http"`
	MySQL MySQL `mapstructure:"mysql"`
	Log   Log   `mapstructure:"log"`
}

type HTTP struct {
	Addr     string `mapstructure:"addr"`
	BasePath string `mapstructure:"basePath"`
}

type MySQL struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
}

type Log struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}

func New() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config.yaml: %w", err)
	}

	env := "dev"
	v.SetConfigName(env)
	if err := v.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("merge %s.yaml: %w", env, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}

var Module = fx.Module("config", fx.Provide(New))
```

- [ ] **Step 3: 验证编译**

```bash
cd /venus/agent/general-agent && go build ./internal/config/...
```
Expected: 编译成功。

---

### Task 3: internal/logger —— zap 日志

**Files:**
- Create: `internal/logger/logger.go`

**Interfaces:**
- Consumes: `*config.Config`
- Produces: `*zap.Logger`, `logger.Module` (fx.Option)

- [ ] **Step 1: 写 logger 实现**

`internal/logger/logger.go`:
```go
package logger

import (
	"general-agent/internal/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg *config.Config) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if cfg.Log.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(zapcore.Lock(zapcore.NewAtomicLevelAt(level))), level)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}

var Module = fx.Module("logger",
	fx.Provide(New),
	fx.Invoke(func(lc fx.Lifecycle, log *zap.Logger) {
		lc.Append(fx.Hook{OnStop: func(ctx context.Context) error { return log.Sync() }})
	}),
)
```

- [ ] **Step 2: 验证编译**

```bash
cd /venus/agent/general-agent && go build ./internal/logger/...
```
Expected: 编译成功。

---

### Task 4: internal/database —— gorm + zap 桥接

**Files:**
- Create: `internal/database/database.go`

**Interfaces:**
- Consumes: `*config.Config`, `*zap.Logger`
- Produces: `*gorm.DB`, `database.Module` (fx.Option)

- [ ] **Step 1: 写 database 实现（含 zap-bridge gorm logger）**

`internal/database/database.go`:
```go
package database

import (
	"context"
	"fmt"
	"time"

	"general-agent/internal/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type zapGormLogger struct {
	log *zap.Logger
}

func (l *zapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface { return l }

func (l *zapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.log.Info(fmt.Sprintf(msg, data...))
}

func (l *zapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.log.Warn(fmt.Sprintf(msg, data...))
}

func (l *zapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log.Error(fmt.Sprintf(msg, data...))
}

func (l *zapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		l.log.Error("gorm query error",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
		)
	} else if elapsed > time.Second {
		l.log.Warn("slow query",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
		)
	}
}

func NewDB(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=Asia%%2FShanghai&parseTime=true",
		cfg.MySQL.Username, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         &zapGormLogger{log: log},
	})
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MySQL.ConnMaxLifetime) * time.Second)

	return db, nil
}

var Module = fx.Module("database", fx.Provide(NewDB))
```

- [ ] **Step 2: 验证编译**

```bash
cd /venus/agent/general-agent && go build ./internal/database/...
```
Expected: 编译成功。

---

### Task 5: internal/errs —— 统一错误处理与响应

**Files:**
- Create: `internal/errs/errs.go`

**Interfaces:**
- Consumes: `*zap.Logger`
- Produces: `*BizError`, `Response`, `OK()`, `Fail()`, `Recovery()` (gin.HandlerFunc), `errs.Module` (fx.Option)

- [ ] **Step 1: 写 errs 实现**

`internal/errs/errs.go`:
```go
package errs

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var globalLog *zap.Logger

func init() { globalLog = zap.NewNop() }

type BizError struct {
	Code int
	Msg  string
}

func (e *BizError) Error() string { return e.Msg }

func NewBizError(code int, msg string) *BizError { return &BizError{Code: code, Msg: msg} }

var (
	ErrInternal   = &BizError{Code: 50000, Msg: "内部错误"}
	ErrNotFound   = &BizError{Code: 10001, Msg: "记录不存在"}
	ErrUserExists = &BizError{Code: 10002, Msg: "用户名已存在"}
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func OK(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{Code: 0, Msg: "success", Data: data})
}

func Fail(ctx *gin.Context, err error) {
	if bizErr, ok := err.(*BizError); ok {
		ctx.JSON(http.StatusOK, Response{Code: bizErr.Code, Msg: bizErr.Msg})
		return
	}
	globalLog.Error("unexpected error", zap.Error(err))
	ctx.JSON(http.StatusInternalServerError, Response{Code: ErrInternal.Code, Msg: ErrInternal.Msg})
}

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recovered",
					zap.Any("panic", r),
					zap.String("stack", string(debug.Stack())),
				)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError,
					Response{Code: ErrInternal.Code, Msg: ErrInternal.Msg})
			}
		}()
		ctx.Next()
	}
}

var Module = fx.Module("errs", fx.Provide(NewErrHandler))

func NewErrHandler(log *zap.Logger) gin.HandlerFunc {
	globalLog = log
	return Recovery(log)
}
```

- [ ] **Step 2: 验证编译**

```bash
cd /venus/agent/general-agent && go build ./internal/errs/...
```
Expected: 编译成功。

---

### Task 6: internal/server —— gin 引擎 + fx lifecycle

**Files:**
- Create: `internal/server/server.go`

**Interfaces:**
- Consumes: `*config.Config`, `*zap.Logger`, `gin.HandlerFunc` (recovery)
- Produces: `*gin.Engine`, `*gin.RouterGroup`, `server.Module` (fx.Option)

- [ ] **Step 1: 写 server 实现**

`internal/server/server.go`:
```go
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"general-agent/internal/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewEngine(cfg *config.Config, log *zap.Logger, recovery gin.HandlerFunc) (*gin.Engine, *gin.RouterGroup) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/health"))
	engine.Use(recovery)
	engine.Use(zapAccessLog(log))

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rg := engine.Group(cfg.HTTP.BasePath)
	return engine, rg
}

func zapAccessLog(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		log.Info("access",
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.Int("status", ctx.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}

var Module = fx.Module("server",
	fx.Provide(NewEngine),
	fx.Invoke(func(lc fx.Lifecycle, cfg *config.Config, engine *gin.Engine, log *zap.Logger) {
		srv := &http.Server{Addr: cfg.HTTP.Addr, Handler: engine}
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						log.Fatal("server start failed", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := srv.Shutdown(shutdownCtx); err != nil {
					return fmt.Errorf("server shutdown: %w", err)
				}
				return nil
			},
		})
	}),
)
```

- [ ] **Step 2: 验证编译**

```bash
cd /venus/agent/general-agent && go build ./internal/server/...
```
Expected: 编译成功。

---

### Task 7: app/user —— 实体迁移与 DTO 定义

**Files:**
- Create: `app/user/user.gen.go`（从现有 `app/entity/user.go` 迁移，勿改内容）
- Create: `app/user/dto.go`

**Interfaces:**
- Consumes: 无（实体自包含）
- Produces: `User` (entity), `CreateUserReq`, `UpdateUserReq`, `ListUserReq`, `UserResp`, `PageResp`

- [ ] **Step 1: 迁移实体文件**

直接将现有 `app/entity/user.go` 的内容拷入 `app/user/user.gen.go`，包名改为 `user`，保留 `// Code generated by gorm.io/gen. DO NOT EDIT.` 注释。

```bash
mkdir -p /venus/agent/general-agent/app/user
```

`app/user/user.gen.go`:
```go
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package user

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID                             string `gorm:"column:id;primaryKey" json:"id"`
	Username                       string `gorm:"column:username;comment:用户名" json:"username"`
	Password                       string `gorm:"column:password;comment:密码，MD5加密" json:"password"`
	PasswordSm3                    string `gorm:"column:password_sm3;comment:密码，SM3加密" json:"password_sm3"`
	Name                           string `gorm:"column:name" json:"name"`
	CreateTime                     int64  `gorm:"column:create_time;comment:创建时间" json:"create_time"`
	LatestLoginTime                int64  `gorm:"column:latest_login_time;comment:最近一次登录时间" json:"latest_login_time"`
	PasswordChangeTime             int64  `gorm:"column:password_change_time;comment:修改密码时间" json:"password_change_time"`
	Status                         int32  `gorm:"column:status;not null;default:1;comment:0:正常|1:第一次登录|2:密码过期" json:"status"`
	RoleID                         int32  `gorm:"column:role_id;comment:角色id" json:"role_id"`
	Classification                 int32  `gorm:"column:classification;not null;default:1;comment:用户分类 0:系统内置用户|1:自定义用户" json:"classification"`
	GatherType                     int32  `gorm:"column:gather_type;default:1;comment:告警视角，默认按照事件名称聚合" json:"gather_type"`
	PlatformSha256Pass             string `gorm:"column:platform_sha256_pass;comment:平台共用的sha256加密后的密码" json:"platform_sha256_pass"`
	IsLock                         int32  `gorm:"column:is_lock;not null;comment:当前用户是否被锁定" json:"is_lock"`
	OtpSn                          string `gorm:"column:otp_sn;comment:otp令牌的序列号" json:"otp_sn"`
	MdsPass                        string `gorm:"column:mds_pass;comment:MDS平台密码" json:"mds_pass"`
	CloudPass                      string `gorm:"column:cloud_pass;comment:云平台规范格式的密码" json:"cloud_pass"`
	RadiusServerID                 int32  `gorm:"column:radius_server_id;comment:radius认证服务器ID" json:"radius_server_id"`
	AllowOperationModuleIds        string `gorm:"column:allow_operation_module_ids;comment:允许操作的功能模块id，多个逗号分割" json:"allow_operation_module_ids"`
	AllowOperationModuleExpireTime int64  `gorm:"column:allow_operation_module_expire_time;not null;comment:允许操作的功能模块限期过期时间，为0则代表永不过期" json:"allow_operation_module_expire_time"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
```

- [ ] **Step 2: 写 DTO**

`app/user/dto.go`:
```go
package user

type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	RoleID   int32  `json:"role_id"`
}

type UpdateUserReq struct {
	Name   string `json:"name"`
	RoleID *int32 `json:"role_id"`
}

type ListUserReq struct {
	PageNum  int `form:"pageNum"`
	PageSize int `form:"pageSize"`
}

type UserResp struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	RoleID     int32  `json:"role_id"`
	Status     int32  `json:"status"`
	CreateTime int64  `json:"create_time"`
}

type PageResp struct {
	List  []UserResp `json:"list"`
	Total int64      `json:"total"`
}

func toUserResp(u *User) UserResp {
	return UserResp{
		ID:         u.ID,
		Username:   u.Username,
		Name:       u.Name,
		RoleID:     u.RoleID,
		Status:     u.Status,
		CreateTime: u.CreateTime,
	}
}
```

- [ ] **Step 3: 验证编译**

```bash
cd /venus/agent/general-agent && go build ./app/user/...
```
Expected: 编译成功。

---

### Task 8: app/user —— Repository 接口与实现

**Files:**
- Create: `app/user/repository.go`

**Interfaces:**
- Consumes: `*gorm.DB`
- Produces: `Repository` (interface), `NewRepository(*gorm.DB) *repositoryImpl`

- [ ] **Step 1: 写 Repository 接口与实现**

`app/user/repository.go`:
```go
package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetByID(id string) (*User, error)
	List(offset, limit int) ([]User, int64, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id string) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) GetByID(id string) (*User, error) {
	var user User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repositoryImpl) List(offset, limit int) ([]User, int64, error) {
	var users []User
	var total int64
	if err := r.db.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *repositoryImpl) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *repositoryImpl) Update(user *User) error {
	return r.db.Where("id = ?", user.ID).Updates(user).Error
}

func (r *repositoryImpl) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&User{}).Error
}
```

- [ ] **Step 2: 验证编译**

```bash
cd /venus/agent/general-agent && go build ./app/user/...
```
Expected: 编译成功。

---

### Task 9: app/user —— Repository 单元测试（sqlite 内存库）

**Files:**
- Create: `app/user/repository_test.go`

**Interfaces:**
- Consumes: `Repository`, `User`
- Produces: 测试覆盖 GetByID / List / Create / Update / Delete

- [ ] **Step 1: 写 Repository 测试**

`app/user/repository_test.go`:
```go
package user

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func seedUser(t *testing.T, db *gorm.DB, u *User) {
	t.Helper()
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
}

func TestRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	seedUser(t, db, &User{ID: "1", Username: "alice", Name: "Alice"})

	got, err := repo.GetByID("1")
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Username != "alice" {
		t.Errorf("expected username=alice, got %s", got.Username)
	}

	_, err = repo.GetByID("999")
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	seedUser(t, db, &User{ID: "1", Username: "alice"})
	seedUser(t, db, &User{ID: "2", Username: "bob"})

	users, total, err := repo.List(0, 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total=2, got %d", total)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	err := repo.Create(&User{ID: "1", Username: "alice"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	var count int64
	db.Model(&User{}).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 user, got %d", count)
	}
}

func TestRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	seedUser(t, db, &User{ID: "1", Username: "alice", Name: "Alice"})

	err := repo.Update(&User{ID: "1", Name: "Alice Updated"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, _ := repo.GetByID("1")
	if got.Name != "Alice Updated" {
		t.Errorf("expected Name='Alice Updated', got %s", got.Name)
	}
}

func TestRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	seedUser(t, db, &User{ID: "1", Username: "alice"})

	err := repo.Delete("1")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err = repo.GetByID("1")
	if err == nil {
		t.Error("expected error after delete")
	}
}
```

- [ ] **Step 2: 运行测试验证**

```bash
cd /venus/agent/general-agent && go test ./app/user/ -run "TestRepository" -v
```
Expected: 5 tests PASS。

---

### Task 10: app/user —— Service

**Files:**
- Create: `app/user/service.go`

**Interfaces:**
- Consumes: `Repository` (interface)
- Produces: `*Service`, `NewService(Repository) *Service`

- [ ] **Step 1: 写 Service**

`app/user/service.go`:
```go
package user

import (
	"errors"
	"time"

	"general-agent/internal/errs"

	"gorm.io/gorm"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(id string) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) List(pageNum, pageSize int) ([]User, int64, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize
	return s.repo.List(offset, pageSize)
}

func (s *Service) Create(req *CreateUserReq) (*User, error) {
	user := &User{
		ID:         generateID(),
		Username:   req.Username,
		Password:   req.Password,
		Name:       req.Name,
		RoleID:     req.RoleID,
		Status:     1,
		CreateTime: time.Now().Unix(),
	}
	if err := s.repo.Create(user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errs.ErrUserExists
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) Update(id string, req *UpdateUserReq) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.RoleID != nil {
		user.RoleID = *req.RoleID
	}
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}

func generateID() string {
	return time.Now().Format("20060102150405") + "000000"
}
```

- [ ] **Step 2: 验证编译**

```bash
cd /venus/agent/general-agent && go build ./app/user/...
```
Expected: 编译成功。

---

### Task 11: app/user —— Service 单元测试（手写 fake Repository）

**Files:**
- Create: `app/user/service_test.go`

**Interfaces:**
- Consumes: `*Service`, `Repository` (fake)
- Produces: 测试覆盖错误转换（NotFound → ErrNotFound）、分页默认值、Create 正向

- [ ] **Step 1: 写 Service 测试（含 fake Repository）**

`app/user/service_test.go`:
```go
package user

import (
	"errors"
	"testing"

	"general-agent/internal/errs"

	"gorm.io/gorm"
)

type fakeRepo struct {
	users map[string]*User
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{users: make(map[string]*User)}
}

func (f *fakeRepo) GetByID(id string) (*User, error) {
	u, ok := f.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (f *fakeRepo) List(offset, limit int) ([]User, int64, error) {
	users := make([]User, 0, len(f.users))
	for _, u := range f.users {
		users = append(users, *u)
	}
	total := int64(len(users))
	if offset > len(users) {
		return nil, total, nil
	}
	end := offset + limit
	if end > len(users) {
		end = len(users)
	}
	return users[offset:end], total, nil
}

func (f *fakeRepo) Create(user *User) error {
	if _, ok := f.users[user.Username]; ok {
		return gorm.ErrDuplicatedKey
	}
	f.users[user.ID] = user
	return nil
}

func (f *fakeRepo) Update(user *User) error {
	f.users[user.ID] = user
	return nil
}

func (f *fakeRepo) Delete(id string) error {
	delete(f.users, id)
	return nil
}

func TestService_GetByID_NotFound(t *testing.T) {
	svc := NewService(newFakeRepo())
	_, err := svc.GetByID("999")
	if !errors.Is(err, errs.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestService_GetByID_Success(t *testing.T) {
	repo := newFakeRepo()
	repo.users["1"] = &User{ID: "1", Username: "alice"}
	svc := NewService(repo)

	got, err := svc.GetByID("1")
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Username != "alice" {
		t.Errorf("expected username=alice, got %s", got.Username)
	}
}

func TestService_List_Defaults(t *testing.T) {
	svc := NewService(newFakeRepo())
	users, total, err := svc.List(0, 0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if total != 0 {
		t.Errorf("expected total=0, got %d", total)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

func TestService_Create_Success(t *testing.T) {
	svc := NewService(newFakeRepo())
	user, err := svc.Create(&CreateUserReq{Username: "alice", Password: "pass", Name: "Alice"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if user.Username != "alice" {
		t.Errorf("expected username=alice, got %s", user.Username)
	}
}
```

- [ ] **Step 2: 运行测试验证**

```bash
cd /venus/agent/general-agent && go test ./app/user/ -run "TestService" -v
```
Expected: 4 tests PASS。

---

### Task 12: app/user —— Controller

**Files:**
- Create: `app/user/controller.go`

**Interfaces:**
- Consumes: `*Service`
- Produces: `*Controller`, `NewController(*Service) *Controller`

- [ ] **Step 1: 写 Controller**

`app/user/controller.go`:
```go
package user

import (
	"general-agent/internal/errs"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	svc *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{svc: svc}
}

// GetByID godoc
// @Summary 获取用户信息
// @Description 根据ID查询单个用户
// @Tags 用户
// @Param id path string true "用户ID"
// @Success 200 {object} errs.Response{data=UserResp}
// @Router /user/{id} [get]
func (ctrl *Controller) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := ctrl.svc.GetByID(id)
	if err != nil {
		errs.Fail(ctx, err)
		return
	}
	errs.OK(ctx, toUserResp(user))
}

// List godoc
// @Summary 用户列表
// @Description 分页查询用户列表
// @Tags 用户
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页大小"
// @Success 200 {object} errs.Response{data=PageResp}
// @Router /user [get]
func (ctrl *Controller) List(ctx *gin.Context) {
	var req ListUserReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errs.Fail(ctx, err)
		return
	}
	users, total, err := ctrl.svc.List(req.PageNum, req.PageSize)
	if err != nil {
		errs.Fail(ctx, err)
		return
	}
	resps := make([]UserResp, len(users))
	for i := range users {
		resps[i] = toUserResp(&users[i])
	}
	errs.OK(ctx, PageResp{List: resps, Total: total})
}

// Create godoc
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户
// @Param req body CreateUserReq true "创建请求"
// @Success 200 {object} errs.Response{data=UserResp}
// @Router /user [post]
func (ctrl *Controller) Create(ctx *gin.Context) {
	var req CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errs.Fail(ctx, err)
		return
	}
	user, err := ctrl.svc.Create(&req)
	if err != nil {
		errs.Fail(ctx, err)
		return
	}
	errs.OK(ctx, toUserResp(user))
}

// Update godoc
// @Summary 更新用户
// @Description 根据ID更新用户信息
// @Tags 用户
// @Param id path string true "用户ID"
// @Param req body UpdateUserReq true "更新请求"
// @Success 200 {object} errs.Response{data=UserResp}
// @Router /user/{id} [put]
func (ctrl *Controller) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errs.Fail(ctx, err)
		return
	}
	user, err := ctrl.svc.Update(id, &req)
	if err != nil {
		errs.Fail(ctx, err)
		return
	}
	errs.OK(ctx, toUserResp(user))
}

// Delete godoc
// @Summary 删除用户
// @Description 根据ID删除用户
// @Tags 用户
// @Param id path string true "用户ID"
// @Success 200 {object} errs.Response
// @Router /user/{id} [delete]
func (ctrl *Controller) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := ctrl.svc.Delete(id); err != nil {
		errs.Fail(ctx, err)
		return
	}
	errs.OK(ctx, nil)
}
```

- [ ] **Step 2: 验证编译**

```bash
cd /venus/agent/general-agent && go build ./app/user/...
```
Expected: 编译成功。

---

### Task 13: app/user —— Controller 集成测试（httptest）

**Files:**
- Create: `app/user/controller_test.go`

**Interfaces:**
- Consumes: `*Controller`, `*Service`, fake `Repository`
- Produces: 测试 5 个 HTTP 端点

- [ ] **Step 1: 写 Controller 集成测试**

`app/user/controller_test.go`:
```go
package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	repo := newFakeRepo()
	repo.users["1"] = &User{ID: "1", Username: "alice", Name: "Alice", RoleID: 1, Status: 1, CreateTime: 1000}
	svc := NewService(repo)
	ctrl := NewController(svc)

	r := gin.New()
	rg := r.Group("/api/v1/general-agent")
	rg.GET("/user/:id", ctrl.GetByID)
	rg.GET("/user", ctrl.List)
	rg.POST("/user", ctrl.Create)
	rg.PUT("/user/:id", ctrl.Update)
	rg.DELETE("/user/:id", ctrl.Delete)
	return r
}

func TestController_GetByID(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/general-agent/user/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Code int      `json:"code"`
		Msg  string   `json:"msg"`
		Data UserResp `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Username != "alice" {
		t.Errorf("expected username=alice, got %s", resp.Data.Username)
	}
}

func TestController_GetByID_NotFound(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/general-agent/user/999", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 10001 {
		t.Errorf("expected code=10001, got %d", resp.Code)
	}
}

func TestController_List(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/general-agent/user?pageNum=1&pageSize=10", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Code int      `json:"code"`
		Data PageResp `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Total != 1 {
		t.Errorf("expected total=1, got %d", resp.Data.Total)
	}
}

func TestController_Create(t *testing.T) {
	r := setupRouter()
	body := `{"username":"bob","password":"pass","name":"Bob"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/general-agent/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Code int      `json:"code"`
		Data UserResp `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Username != "bob" {
		t.Errorf("expected username=bob, got %s", resp.Data.Username)
	}
}

func TestController_Update(t *testing.T) {
	r := setupRouter()
	body := `{"name":"Alice Updated"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/general-agent/user/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Code int      `json:"code"`
		Data UserResp `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Name != "Alice Updated" {
		t.Errorf("expected Name='Alice Updated', got %s", resp.Data.Name)
	}
}

func TestController_Delete(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/general-agent/user/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Code int `json:"code"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("expected code=0, got %d", resp.Code)
	}
}
```

- [ ] **Step 2: 运行全部测试**

```bash
cd /venus/agent/general-agent && go test ./app/user/ -v
```
Expected: 全部 14 tests PASS（5 repository + 4 service + 5 controller）。

---

### Task 14: app/user —— fx.Module（路由注册）

**Files:**
- Create: `app/user/module.go`

**Interfaces:**
- Consumes: `*gin.RouterGroup`, `*Controller`
- Produces: `user.Module` (fx.Option)

- [ ] **Step 1: 写 module**

`app/user/module.go`:
```go
package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Module("user",
	fx.Provide(NewRepository, NewService, NewController),
	fx.Invoke(registerRoutes),
)

func registerRoutes(rg *gin.RouterGroup, ctrl *Controller) {
	rg.GET("/user/:id", ctrl.GetByID)
	rg.GET("/user", ctrl.List)
	rg.POST("/user", ctrl.Create)
	rg.PUT("/user/:id", ctrl.Update)
	rg.DELETE("/user/:id", ctrl.Delete)
}
```

- [ ] **Step 2: 验证编译**

```bash
cd /venus/agent/general-agent && go build ./app/user/...
```
Expected: 编译成功。

---

### Task 15: cmd/main.go —— fx 组装入口

**Files:**
- Rewrite: `cmd/main.go`

**Interfaces:**
- Consumes: 所有 infrastructure module + user.Module
- Produces: 可运行的应用

- [ ] **Step 1: 重写 main.go**

`cmd/main.go`:
```go
package main

import (
	"general-agent/app/user"
	"general-agent/internal/config"
	"general-agent/internal/database"
	"general-agent/internal/errs"
	"general-agent/internal/logger"
	"general-agent/internal/server"

	"go.uber.org/fx"
)

// @title           General Agent
// @version         1.0
// @description     这是一个通用Agent智能体.
// @host            localhost:8080
// @BasePath        /api/v1/general-agent
func main() {
	fx.New(
		config.Module,
		logger.Module,
		database.Module,
		errs.Module,
		server.Module,
		user.Module,
	).Run()
}
```

- [ ] **Step 2: 验证编译**

```bash
cd /venus/agent/general-agent && go build ./cmd/...
```
Expected: 编译成功。

---

### Task 16: tools/gen —— 实体生成器改造

**Files:**
- Create: `tools/gen/main.go`

**Interfaces:**
- Consumes: `internal/config`
- Produces: 可执行生成器

- [ ] **Step 1: 写生成器**

`tools/gen/main.go`:
```go
package main

import (
	"fmt"
	"log"

	"general-agent/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var targets = map[string]string{
	"user": "./app/user",
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.Username, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatalf("connect mysql: %v", err)
	}

	g := gen.NewGenerator(gen.Config{
		ModelPkgPath: "./app/user",
	})
	g.UseDB(db)

	for tableName := range targets {
		g.ApplyBasic(g.GenerateModel(tableName))
	}
	g.Execute()
}
```

- [ ] **Step 2: 验证编译**

```bash
cd /venus/agent/general-agent && go build ./tools/gen/...
```
Expected: 编译成功。

---

### Task 17: Makefile

**Files:**
- Create: `Makefile`

- [ ] **Step 1: 写 Makefile**

```makefile
.PHONY: run test gen swagger

run:
	go run ./cmd

test:
	go test ./... -v

gen:
	go run ./tools/gen

swagger:
	swag init -g cmd/main.go -o docs
```

- [ ] **Step 2: 验证 make test**

```bash
cd /venus/agent/general-agent && make test
```
Expected: 全部 14 tests PASS。

---

### Task 18: 最终验证

- [ ] **Step 1: 确认全量编译通过**

```bash
cd /venus/agent/general-agent && go build ./...
```
Expected: 无错误。

- [ ] **Step 2: 确认全量测试通过**

```bash
cd /venus/agent/general-agent && go test ./... -v
```
Expected: 全部 PASS。

- [ ] **Step 3: 确认最终目录结构**

```bash
cd /venus/agent/general-agent && find . -type f -not -path '*/.git/*' -not -path '*/go.sum' | sort
```
Expected: 与 spec 第 3 节目录结构一致。

- [ ] **Step 4: 报告给用户**

列出所有创建/删除/修改的文件清单，以及测试结果摘要，提醒用户自行 review 并 git commit。