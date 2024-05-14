/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package mysql

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	ormInstance *gorm.DB
	once        sync.Once
)

type DbConfig struct {
	DbHost        string `toml:"dbHost"`
	DbPort        int    `toml:"dbPort"`
	DbUsername    string `toml:"dbUsername"`
	DbPassword    string `toml:"dbPassword"`
	DbDatabase    string `toml:"dbDatabase"`
	DbCharset     string `toml:"dbCharset"`
	DbCollation   string `toml:"dbCollation"`
	DbPrefix      string `toml:"dbPrefix"`
	DbMaxIdleConn int    `toml:"dbMaxIdleConn"` // 添加最大空闲连接数设置
	DbMaxOpenConn int    `toml:"dbMaxOpenConn"` // 添加最大打开连接数设置
	DbMaxIdleTime int    `toml:"dbMaxIdleTime"` // 最大空闲时间（秒）
}

func InitConnection(configFile string) {
	once.Do(func() {
		var err error

		file, err := os.Open(configFile)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to open TOML file: %v", err))
		}
		defer func(file *os.File) { _ = file.Close() }(file)

		var dbConfig DbConfig
		decoder := toml.NewDecoder(file)
		if _, err = decoder.Decode(&dbConfig); err != nil {
			fmt.Println(fmt.Errorf("failed to decode TOML file: %v", err))
		}

		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			dbConfig.DbUsername,
			dbConfig.DbPassword,
			dbConfig.DbHost,
			dbConfig.DbPort,
			dbConfig.DbDatabase,
			dbConfig.DbCharset,
		)

		newLogger := gormlogger.Default.LogMode(gormlogger.Silent)
		ormInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
		if err != nil {
			panic(fmt.Errorf("failed to connect to MySQL database: %v", err))
		}

		sqlDB, err := ormInstance.DB()
		if err != nil {
			panic(fmt.Errorf("failed to get database object: %v", err))
		}

		// 设置连接池参数
		sqlDB.SetMaxIdleConns(dbConfig.DbMaxIdleConn)
		sqlDB.SetMaxOpenConns(dbConfig.DbMaxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.DbMaxIdleTime) * time.Second)
	})
}

func GetOrmInstance() *gorm.DB {
	return ormInstance
}
