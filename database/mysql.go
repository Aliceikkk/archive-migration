package database

import (
    "datarp/config"
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "log"
    "os"
    "time"
)

var DB *gorm.DB

func InitDB() error {
    conf := config.GetConfig()
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        conf.MySQL.Username,
        conf.MySQL.Password,
        conf.MySQL.Host,
        conf.MySQL.Port,
        conf.MySQL.Database,
    )

    // 自定义GORM日志配置
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             time.Second,
            LogLevel:                  logger.Info,
            IgnoreRecordNotFoundError: true,
            Colorful:                  true,
        },
    )

    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: newLogger,
    })
    if err != nil {
        return fmt.Errorf("连接数据库失败: %v", err)
    }

    sqlDB, err := DB.DB()
    if err != nil {
        return fmt.Errorf("获取数据库实例失败: %v", err)
    }

    // 设置连接池参数
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return nil
} 