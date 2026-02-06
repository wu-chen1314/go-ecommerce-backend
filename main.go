package main

import (
    "go-web-server/conf"
    "go-web-server/models"
    "go-web-server/routers"
)
// @title           Go电商后端 API
// @version         1.0
// @description     这是一个基于 Gin + GORM + JWT 的电商后端项目
// @host            localhost:8080
// @BasePath        /

func main() {
    // 1. 初始化数据库
    conf.InitDB()
    conf.InitRedis()
    // 2. 自动迁移表结构 (放在这里或者 conf 里都可以)
    conf.DB.AutoMigrate(&models.Product{}, &models.User{}) // 迁移商品表和用户表

    // 3. 加载路由
    r := routers.SetupRouter()

    // 4. 启动服务
    r.Run(":8080")
}