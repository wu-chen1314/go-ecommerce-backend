package conf

import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// 全局 DB 对象
var DB *gorm.DB

func InitDB() {
    // 这里的配置还是用你刚才成功的配置
    dsn := "root:12345678@tcp(host.docker.internal:3306)/go_ecommerce?charset=utf8mb4&parseTime=True&loc=Local"
    
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("数据库连接失败: " + err.Error())
    }
    fmt.Println("MySQL 连接成功！")
}