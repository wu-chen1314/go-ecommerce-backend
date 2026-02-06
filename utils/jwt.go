package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

// 定义密钥 (生产环境请读取环境变量，不要写死在代码里)
var jwtSecret = []byte("my_super_secret_key_2026")

// 定义 Claims (载荷)，也就是令牌里存的数据
type Claims struct {
    UserID uint   `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims // 包含过期时间、签发人等标准字段
}

// GenerateToken 生成 Token
func GenerateToken(userID uint, username string) (string, error) {
    // 1. 设置有效期 (比如 24 小时)
    nowTime := time.Now()
    expireTime := nowTime.Add(24 * time.Hour)

    // 2. 创建 Claims
    claims := Claims{
        UserID: userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expireTime), // 过期时间
            Issuer:    "go-ecommerce",                // 签发人
        },
    }

    // 3. 使用 HS256 算法生成 Token 对象
    tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // 4. 使用密钥签名并生成字符串
    token, err := tokenClaims.SignedString(jwtSecret)
    return token, err
}
// 在 utils/jwt.go 里添加

// ParseToken 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
    // 解析 token
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    // 验证 token 是否有效
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, err
}