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

type mysqlConf struct {
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

type conf struct {
	MySQL mysqlConf `toml:"mysql"`
}

func InitConnection() {
	once.Do(func() {
		config := conf{}
		file, err := os.Open("config.toml")
		if err != nil {
			fmt.Println(fmt.Errorf("failed to open TOML file: %v", err))
		}
		defer func(file *os.File) { _ = file.Close() }(file)

		decoder := toml.NewDecoder(file)
		if _, err = decoder.Decode(&config); err != nil {
			fmt.Println(fmt.Errorf("failed to decode TOML file: %v", err))
		}

		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			config.MySQL.DbUsername,
			config.MySQL.DbPassword,
			config.MySQL.DbHost,
			config.MySQL.DbPort,
			config.MySQL.DbDatabase,
			config.MySQL.DbCharset,
		)

		newLogger := gormlogger.Default.LogMode(gormlogger.Silent)
		ormInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
		if err != nil {
			fmt.Println(fmt.Errorf("failed to connect to MySQL database: %v", err))
		}

		sqlDB, err := ormInstance.DB()
		if err != nil {
			fmt.Println(fmt.Errorf("failed to get database object: %v", err))
		}

		// 设置连接池参数
		sqlDB.SetMaxIdleConns(config.MySQL.DbMaxIdleConn)
		sqlDB.SetMaxOpenConns(config.MySQL.DbMaxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Duration(config.MySQL.DbMaxIdleTime) * time.Second)
	})
}
func GetOrmInstance() *gorm.DB {
	return ormInstance
}
