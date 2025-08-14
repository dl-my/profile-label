package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"profile-label/config"
	"profile-label/model"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		config.GlobalConfig.DB.Username,
		config.GlobalConfig.DB.Password,
		config.GlobalConfig.DB.Host,
		config.GlobalConfig.DB.Port,
		config.GlobalConfig.DB.Name,
		config.GlobalConfig.DB.Charset,
		config.GlobalConfig.DB.ParseTime,
		config.GlobalConfig.DB.Loc,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("数据库连接失败%v\n", err)
	}

	// 自动迁移
	err = DB.AutoMigrate(&model.Solscan{})
	if err != nil {
		log.Printf("自动迁移失败%v\n", err)
	}
	fmt.Println("数据库连接成功")
}
