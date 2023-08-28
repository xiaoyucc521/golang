package dao

import (
	"context"
	"db/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

var _db *gorm.DB

// InitDB 配置
func InitDB() {
	DBConfig := config.Conf.Mysql
	host := DBConfig.Host
	port := DBConfig.Port
	database := DBConfig.Database
	username := DBConfig.UserName
	password := DBConfig.Password
	charset := DBConfig.Charset
	dsn := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=" + charset + "&parseTime=true"}, "")
	err := Database(dsn)
	if err != nil {
		panic(err)
	}
}

// Database 初始化链接
func Database(dsn string) error {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	// 设置空闲连接池中的最大连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置连接可重用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)
	_db = db

	// 迁移
	migration()

	return err
}

// NewDBClient 从连接池中取一个链接
func NewDBClient(context context.Context) *gorm.DB {
	db := _db
	return db.WithContext(context)
}
