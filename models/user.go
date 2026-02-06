package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Username string `json:"username" gorm:"unique"` // 用户名唯一
    Password string `json:"password"`               // 存加密后的哈希值
}