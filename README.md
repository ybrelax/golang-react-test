
一个基于Go和SQLite的图书馆管理系统后端服务，提供用户认证、图书管理和借阅功能。

## 功能特点

- 用户认证：注册、登录和JWT令牌管理
- 角色权限：管理员和普通用户角色区分
- 图书管理：图书的增删改查和多副本管理
- 借阅系统：图书借阅、归还和状态跟踪
- API限流：防止过度请求保护系统

## 技术栈

- **后端框架**：Gin
- **ORM**：GORM
- **数据库**：SQLite
- **认证**：JWT
- **限流**：golang.org/x/time/rate

## 项目结构

```
library-api/
├── config/         # 配置管理
├── controllers/    # 控制器
├── database/       # 数据库连接
├── middleware/     # 中间件
├── models/         # 数据模型
├── routes/         # 路由定义
├── go.mod          # 依赖管理
└── main.go         # 应用入口
```

## 安装与运行

### 前提条件

- Go 1.22+ 
- Git

### 安装步骤

1. 克隆仓库
```bash
git clone <repository-url>
cd library-api
```

2. 安装依赖
```bash
go mod tidy
```

3. 创建环境变量文件（可选）
```bash
touch .env
# 添加以下内容
JWT_SECRET=your_jwt_secret
DB_PATH=./library.db
SERVER_PORT=8080
```

4. 运行应用
```bash
go run main.go
```

服务器将在 http://localhost:8080 启动

## API 文档

### 认证接口

#### 注册用户
- **URL**: `/auth/register`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }
  ```
- **响应**: 201 Created

#### 用户登录
- **URL**: `/auth/login`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "email": "john@example.com",
    "password": "password123"
  }
  ```
- **响应**: 200 OK (返回JWT令牌)

### 用户接口

#### 获取我的借阅
- **URL**: `/api/user/borrows`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**: 200 OK (借阅列表)

### 图书接口

#### 获取图书列表
- **URL**: `/api/books`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**: 200 OK (图书列表)

#### 获取图书详情
- **URL**: `/api/books/:id`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**: 200 OK (图书详情)

#### 借阅图书
- **URL**: `/api/books/borrow`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "book_id": 1,
    "days": 14
  }
  ```
- **响应**: 200 OK (借阅信息)

#### 归还图书
- **URL**: `/api/books/return`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "borrow_id": 1
  }
  ```
- **响应**: 200 OK (归还确认)

### 管理员接口

#### 添加图书
- **URL**: `/api/books`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "title": "Go Programming",
    "isbn": "9781234567890",
    "description": "Learn Go programming",
    "publisher": "Tech Press",
    "publication_date": "2023-01-15",
    "author_ids": [1, 2]
  }
  ```
- **响应**: 201 Created (新图书信息)

#### 更新图书
- **URL**: `/api/books/:id`
- **方法**: `PUT`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**: (与添加图书类似)
- **响应**: 200 OK (更新后的图书信息)

#### 删除图书
- **URL**: `/api/books/:id`
- **方法**: `DELETE`
- **请求头**: `Authorization: Bearer {token}`
- **响应**: 200 OK (删除确认)

#### 添加图书副本
- **URL**: `/api/books/:id/copies`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "copies_count": 5
  }
  ```
- **响应**: 201 Created (副本添加结果)

## 配置说明

通过环境变量或.env文件配置以下参数：

- `JWT_SECRET`: JWT签名密钥（必填）
- `DB_PATH`: SQLite数据库文件路径（默认：./library.db）
- `SERVER_PORT`: 服务器端口（默认：8080）
- `JWT_EXPIRY_HOURS`: JWT过期时间（小时，默认：72）
- `RATE_LIMIT_RPS`: API限流（每秒请求数，默认：10）
- `RATE_LIMIT_BURST`: 限流突发容量（默认：20）

## 开发说明

### 数据库迁移

应用启动时会自动执行数据库迁移，创建所需表结构。

### 依赖管理

依赖项在go.mod中定义，使用以下命令更新依赖：
```bash
go get -u
```

### 测试

添加测试用例后，使用以下命令运行测试：
```bash
go test ./...
```

## 许可证

[MIT](LICENSE)