package controllers

import (
	"go-web-server/conf"
	"go-web-server/models"
	"go-web-server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register 用户注册
func Register(c *gin.Context) {
    // 1. 接收参数
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 2. 校验用户名是否已存在
    var temp models.User
    // Try to find a user with this username
    if err := conf.DB.Where("username = ?", user.Username).First(&temp).Error; err == nil {
        // 如果 err == nil 说明找到了，也就是用户名已存在
        c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
        return
    }

    // 3. 密码加密 (核心步骤！)
    // GenerateFromPassword 会自动加盐
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "加密失败"})
        return
    }
    // 将加密后的密码赋值回去
    user.Password = string(hashedPassword)

    // 4. 创建用户
    conf.DB.Create(&user)

    c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}
//登录
// Login 用户登录
func Login(c *gin.Context) {
    // 1. 定义一个结构体接收前端传来的数据
    // 这里也可以直接复用 models.User，但为了代码清晰，我们定义一个临时的
    var inputUser struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    // 绑定参数
    if err := c.ShouldBindJSON(&inputUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 2. 根据用户名去数据库查找用户
    var dbUser models.User
    if err := conf.DB.Where("username = ?", inputUser.Username).First(&dbUser).Error; err != nil {
        // 如果查询不到，说明用户不存在
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        return
    }

    // 3. 验证密码 (核心步骤)
    // 参数1: 数据库里存的加密哈希值 (dbUser.Password)
    // 参数2: 前端传来的明文密码 (inputUser.Password)
    err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(inputUser.Password))
    
    if err != nil {
        // 如果 err 不为空，说明密码比对失败
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        return
    }

    // 4. 登录成功
    // (下一章我们会在这里生成 JWT Token 返回给前端)
   token, err := utils.GenerateToken(dbUser.ID, dbUser.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
        return
    }

    // 5. 返回 Token 给前端
    c.JSON(http.StatusOK, gin.H{
        "message": "登录成功",
        "token":   token, //这就是你的通行证！
    })
}
// GetUserProfile 获取当前登录用户的个人信息
func GetUserProfile(c *gin.Context) {
    // 1. 从上下文取出中间件存进去的 userID
    // c.Get 返回的是 interface{} (空接口)，需要用断言转成 uint
    userID, exists := c.Get("userID")
    if !exists {
        // 理论上经过中间件不应该进这里，但为了稳健
        c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
        return
    }

    // 2. 去数据库查这个人的详细信息
    var user models.User
    if err := conf.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
        return
    }

    // 3. 返回信息 (注意不要把加密后的密码也返回去，很危险！)
    c.JSON(http.StatusOK, gin.H{
        "id":       user.ID,
        "username": user.Username,
        "role":     "普通用户", // 以后可以扩展角色字段
    })
}