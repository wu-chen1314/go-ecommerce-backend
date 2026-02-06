package middlewares

import (
    "net/http"
    "strings"
    "go-web-server/utils" // 注意你的 module 名
    "github.com/gin-gonic/gin"
)


// JWTAuthMiddleware 中间件
func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取 Authorization Header
        // 约定格式: Bearer <token>
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "未携带令牌，禁止访问"})
            c.Abort() // 阻止请求继续往下走
            return
        }

        // 2. 解析 Token (去掉 "Bearer " 前缀)
        parts := strings.SplitN(authHeader, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
             c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌格式错误"})
             c.Abort()
             return
        }
        
        tokenString := parts[1]

        // 3. 解析并验证
        // 这里需要把 utils/jwt.go 里的 ParseToken 补上（稍后给你代码）
        // 我们先假设有一个 ParseToken 函数
        claims, err := utils.ParseToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌无效或已过期"})
            c.Abort()
            return
        }

        // 4. 将当前用户 ID 存入上下文，方便后面的 Controller 使用
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
		

        // 5. 放行
        c.Next()
    }
}