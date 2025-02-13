package dal

import (
	"fmt"

	"blockSBOM/backend/internal/config"
	"blockSBOM/backend/internal/dal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) error {
	var err error

	// 先连接MySQL（不指定数据库）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 创建数据库（如果不存在）
	createDB := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", cfg.Database.DBName)
	if err := db.Exec(createDB).Error; err != nil {
		return fmt.Errorf("创建数据库失败: %v", err)
	}

	// 重新连接（指定数据库）
	dsn = cfg.GetDSN()
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接指定数据库失败: %v", err)
	}

	// 自动迁移表结构
	if err := DB.AutoMigrate(&model.User{}); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("获取数据库实例失败: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("关闭数据库连接失败: %v", err)
		}
	}
	return nil
}
