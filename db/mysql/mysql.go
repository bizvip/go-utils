package mysql

import (
	"fmt"
	"log"
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
	DbMaxIdleConn int    `toml:"dbMaxIdleConn"` // 最大空闲连接数
	DbMaxOpenConn int    `toml:"dbMaxOpenConn"` // 最大打开连接数
	DbMaxIdleTime int    `toml:"dbMaxIdleTime"` // 最大空闲时间（秒）
}

type myConf struct {
	Mysql DbConfig `toml:"mysql"`
}

func readConfig(configFile string) (*myConf, error) {
	var conf myConf
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := toml.NewDecoder(file)
	if _, err = decoder.Decode(&conf); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &conf, nil
}

func InitConnection(configFile string) {
	once.Do(func() {
		conf, err := readConfig(configFile)
		if err != nil {
			log.Fatalf("Failed to read configuration: %v", err)
		}

		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			conf.Mysql.DbUsername,
			conf.Mysql.DbPassword,
			conf.Mysql.DbHost,
			conf.Mysql.DbPort,
			conf.Mysql.DbDatabase,
			conf.Mysql.DbCharset,
		)

		newLogger := gormlogger.Default.LogMode(gormlogger.Silent)
		ormInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		sqlDB, err := ormInstance.DB()
		if err != nil {
			log.Fatalf("Failed to get database instance: %v", err)
		}

		// 设置连接池参数
		sqlDB.SetMaxIdleConns(conf.Mysql.DbMaxIdleConn)
		sqlDB.SetMaxOpenConns(conf.Mysql.DbMaxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Duration(conf.Mysql.DbMaxIdleTime) * time.Second)

		log.Println("Database connection successfully initialized")
	})
}

func GetOrmInstance() *gorm.DB {
	return ormInstance
}
