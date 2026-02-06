package routers

import (
    "github.com/gin-gonic/gin"
    
    // 这两个包必须在下面用到，否则就会报你截图里的错
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"

    "go-web-server/controllers"
    "go-web-server/middlewares"

    // 记得这里要是你自己的 module 名
    _ "go-web-server/docs" 
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // ==========================================
    // 💡 加上这行！这就是在使用上面导入的包
    // ==========================================
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    
    // 公开组
    public := r.Group("/")
    {
        public.POST("/user/register", controllers.Register)
        public.POST("/user/login", controllers.Login)
        public.GET("/products", controllers.GetProducts)
    }

    // 私有组
    authorized := r.Group("/")
    authorized.Use(middlewares.JWTAuthMiddleware())
    {
        authorized.GET("/user/me", controllers.GetUserProfile)
        authorized.POST("/product/add", controllers.AddProduct)
        authorized.PUT("/product/:id", controllers.UpdateProduct)
        authorized.DELETE("/product/:id", controllers.DeleteProduct)
    }

    return r
}