package postgres

import (
	"fmt"
	"general-agent/config"
	"general-agent/extension/logz"
	"general-agent/extension/objects"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	MASTER = "MASTER"              // 默认主数据源
	RDBs   = map[string]*gorm.DB{} // 初始化时加载数据源到集合
)

type PostgresDb struct {
	*gorm.DB
}

func NewPostgresDB(conf *config.Config) *PostgresDb {
	// 初始化主数据源
	primary, err := NewPostgresClient(conf, conf.Postgres.Database)
	if err != nil {
		panic(err)
	}
	// 初始化其他数据源
	// for _, database := range conf.Postgres.OtherDB {
	// 	datasource, err := NewPostgresClient(conf, database)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	RDBs[database] = datasource
	// }

	return &PostgresDb{primary}
}

func (m *PostgresDb) WithTX(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return m.DB
}

// Use 获取数据源,并控制是否开启事务
func (m *PostgresDb) Use(ds ...string) *gorm.DB {
	dsName := MASTER
	if len(ds) > 0 && len(ds[0]) > 0 {
		dsName = ds[0]
	}
	return RDBs[dsName]
}

// NewPostgresClient 初始化数据库实例
func NewPostgresClient(cf *config.Config, database string) (*gorm.DB, error) {
	option := cf.Postgres
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		option.Host, option.Username, option.Password, database, option.Port)
	options := gorm.Config{
		TranslateError: true,
	}
	if true {
		options.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      objects.If(false, logger.Info, logger.Silent),
				Colorful:      true,
			},
		)
	}
	db, err := gorm.Open(postgres.Open(dsn), &options)
	if err != nil {
		logz.ErrorNoCtx("PostgreSQL init error", "err", err)
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		logz.ErrorNoCtx("PostgreSQL init error", "err", err)
		return nil, err
	}
	// 连接池
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(3600)
	return db, nil
}
