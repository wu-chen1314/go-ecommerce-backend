# Go 电商后端系统 (Go E-Commerce Backend)

![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)
![Gin](https://img.shields.io/badge/Framework-Gin-000000?style=flat&logo=go)
![GORM](https://img.shields.io/badge/ORM-GORM-red?style=flat)
![Redis](https://img.shields.io/badge/Cache-Redis-DC382D?style=flat&logo=redis)
![Docker](https://img.shields.io/badge/Deploy-Docker-2496ED?style=flat&logo=docker)

这是一个基于 Golang 开发的企业级电商后端 API 项目。采用 **Gin** 框架作为路由核心，结合 **GORM** 进行数据库操作，引入 **Redis** 实现高并发缓存策略，并使用 **Docker** 实现了容器化部署。

## 🛠 技术栈 (Tech Stack)

- **语言**: Golang 1.25+
- **Web 框架**: Gin
- **数据库**: MySQL 8.0 (GORM v2)
- **缓存**: Redis v9 (Cache-Aside Pattern)
- **鉴权**: JWT (JSON Web Token) + BCrypt 加密
- **文档**: Swagger (自动生成 API 文档)
- **部署**: Docker + Docker Compose (可选)

## ✨ 核心功能 (Features)

1.  **用户系统**: 注册、登录、JWT 身份认证、密码加密存储。
2.  **商品管理**: 商品的增删改查 (CRUD)，支持管理员权限控制。
3.  **性能优化**: 实现 Redis 缓存旁路策略，大幅提升热点数据查询速度。
4.  **工程化**: 标准 MVC 分层架构 (Controller / Service / Model)。
5.  **API 文档**: 集成 Swagger，访问 `/swagger/index.html` 即可在线调试。

## 🚀 快速开始 (How to Run)

### 方式一：使用 Docker (推荐)

```bash
# 1. 构建镜像
docker build -t go-mall .

# 2. 运行容器
docker run -p 8080:8080 --name my-app go-mall
