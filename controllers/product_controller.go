package controllers

import (
    "encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	
	"go-web-server/conf"   // 引入配置包
	"go-web-server/models" // 引入模型包
)
// AddProduct 添加商品
// @Summary      添加新商品
// @Description  创建一个新的商品记录 (需要管理员权限)
// @Tags         商品管理
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string        true  "Bearer Token"
// @Param        product        body      models.Product true "商品信息"
// @Success      200            {object}  models.Product "成功返回商品信息"
// @Router       /product/add [post]
func AddProduct(c *gin.Context) {
    var p models.Product
    if err := c.ShouldBindJSON(&p); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 使用全局 DB
    if err := conf.DB.Create(&p).Error; err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库写入错误"})
         return
    }

    c.JSON(http.StatusOK, gin.H{"message": "添加成功", "product": p})
}

// GetProducts 获取所有商品
// GetProducts 获取所有商品 (带 Redis 缓存版)
// @Summary      获取商品列表
// @Description  优先查缓存，缓存没有查数据库
// @Tags         商品管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Product
// @Router       /products [get]
func GetProducts(c *gin.Context) {
	// 定义缓存的 Key
	cacheKey := "products_list"

	// ================= STEP 1: 先查 Redis =================
	// RDB.Get 需要传入 Context，我们复用 conf.Ctx
	val, err := conf.RDB.Get(conf.Ctx, cacheKey).Result()
	
	if err == nil {
		// --- 情况 A: 缓存命中 (Hit) ---
		fmt.Println("🚀 命中 Redis 缓存，直接返回！")
		
		var products []models.Product
		// Redis 里存的是 JSON 字符串，取出来要反序列化变回 Struct
		json.Unmarshal([]byte(val), &products)
		
		c.JSON(http.StatusOK, gin.H{"data": products, "source": "redis_cache"})
		return
	} else if err != redis.Nil {
		// 如果报错不是因为“没找到”，而是 Redis 挂了等其他原因，打印一下但不中断，继续查库
		fmt.Println("Redis 异常:", err)
	}

	// ================= STEP 2: 缓存没命中，查 MySQL =================
	fmt.Println("🐢 缓存未命中，正在查询 MySQL...")
	
	var products []models.Product
	result := conf.DB.Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	}

	// ================= STEP 3: 写入 Redis (回填) =================
	// 也就是“下次我就记得了”
	
	// 1. 将数据序列化为 JSON 字符串
	data, _ := json.Marshal(products)
	
	// 2. 存入 Redis
	// 参数: Context, Key, Value, 过期时间
	// 我们设置 10 秒过期，方便测试（生产环境可能设置 1 小时）
	err = conf.RDB.Set(conf.Ctx, cacheKey, data, 10*time.Second).Err()
	if err != nil {
		fmt.Println("写入缓存失败:", err)
	} else {
		fmt.Println("✅ 数据已回填至 Redis")
	}

	c.JSON(http.StatusOK, gin.H{"data": products, "source": "mysql_db"})
}

// UpdateProduct 更新商品 (你也把之前的逻辑搬过来)
func UpdateProduct(c *gin.Context) {
    id := c.Param("id")
    var p models.Product
    if err := conf.DB.First(&p, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
        return
    }
    // 假设更新价格
    conf.DB.Model(&p).Update("Price", 88888)
    c.JSON(http.StatusOK, gin.H{"message": "更新成功", "product": p})
}

// DeleteProduct 删除商品
func DeleteProduct(c *gin.Context) {
    id := c.Param("id")
    conf.DB.Delete(&models.Product{}, id)
    c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}