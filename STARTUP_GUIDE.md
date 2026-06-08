# Mac Cluster Monitor 启动指南

## 环境准备

### 1. 安装必要软件

- **Go 1.21+**: https://golang.org/dl/
- **MySQL 8.0**: https://dev.mysql.com/downloads/installer/
- **Node.js 20+**: https://nodejs.org/

### 2. 初始化数据库

```bash
# 登录MySQL
mysql -u root -p

# 创建数据库
CREATE DATABASE mac_cluster_monitor CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 导入表结构
mysql -u root -p mac_cluster_monitor < database/schema.sql
```

## 启动步骤

### 第一步：启动服务端

```bash
# 1. 进入服务端目录
cd server

# 2. 复制配置文件（首次运行）
copy ..\config\server.local.yaml config.yaml

# 3. 根据你的MySQL配置修改 config.yaml
# 主要修改 database.password 为你的MySQL密码

# 4. 下载依赖
go mod tidy

# 5. 编译并运行
go build -o server.exe
server.exe

# 服务端启动在 http://localhost:8080
```

### 第二步：启动前端（新终端窗口）

```bash
# 1. 进入前端目录
cd web

# 2. 安装依赖
npm install

# 3. 启动开发服务器
npm run dev

# 前端启动在 http://localhost:3000
```

### 第三步：测试Agent（新终端窗口）

```bash
# 1. 进入Agent目录
cd agent

copy ..\config\agent.local.yaml config.yaml
# 2. 编译Agent（Windows版本，用于测试）
go build -o agent.exe main.go

# 3. 设置环境变量并运行
set MONITOR_SERVER=http://localhost:8080
set MONITOR_INTERVAL=5
agent.exe
```

## 访问系统

1. 打开浏览器访问: http://localhost:3000
2. 使用默认账号登录:
   - 用户名: `admin`
   - 密码: `admin123`
3. 你应该能看到Agent上报的设备信息

## 项目结构

```
mac-cluster-monitor/
├── config/              # 配置文件模板
├── server/              # Go后端
│   ├── config.yaml      # 当前使用的配置
│   └── main.go
├── web/                 # Vue前端
│   └── src/
├── agent/               # Go Agent
│   └── main.go
└── database/            # 数据库脚本
    └── schema.sql
```

## 常用命令

### 服务端

```bash
cd server

# 开发模式运行
go run main.go

# 编译
 go build -o server.exe

# 指定配置文件运行
set CONFIG_PATH=config.yaml
server.exe
```

### 前端

```bash
cd web

# 安装依赖
npm install

# 开发模式
npm run dev

# 构建生产版本
npm run build
```

### Agent

```bash
cd agent

# Windows测试
go build -o agent.exe
set MONITOR_SERVER=http://localhost:8080
agent.exe

# 编译Mac版本（交叉编译）
# Apple Silicon Mac
set GOOS=darwin
set GOARCH=arm64
go build -o agent-darwin-arm64

# Intel Mac
set GOOS=darwin
set GOARCH=amd64
go build -o agent-darwin-amd64
```

## 配置说明

### 服务端配置 (server/config.yaml)

```yaml
server:
  port: "8080"              # 服务端口

database:
  host: "localhost"         # MySQL地址
  port: "3306"              # MySQL端口
  user: "root"              # MySQL用户名
  password: "your_password" # MySQL密码（务必修改）
  name: "mac_cluster_monitor"

jwt:
  secret: "your-secret"     # JWT密钥（务必修改）
```

### Agent环境变量

- `MONITOR_SERVER`: 监控服务器地址，如 `http://localhost:8080`
- `MONITOR_INTERVAL`: 上报间隔（秒），默认5秒

## 部署到生产环境

### 服务端部署到Linux服务器

1. 修改 `config/server.prod.yaml` 中的配置
2. 复制到服务器：`server/config.yaml`
3. 使用 `scripts/deploy-server.sh` 自动部署

### Agent部署到Mac Mini

1. 交叉编译Mac版本
2. 修改 `config/agent.prod.yaml` 中的服务器地址
3. 复制到U盘：`agent/config.yaml`
4. 使用 `agent/install.sh` 安装

## 故障排查

### 服务端无法启动

- 检查MySQL是否运行
- 检查config.yaml中的数据库密码是否正确
- 检查端口8080是否被占用

### 前端无法访问

- 检查服务端是否已启动
- 检查vite.config.js中的代理配置

### Agent无法上报

- 检查MONITOR_SERVER环境变量是否正确
- 检查服务端是否允许跨域（CORS）
- 检查防火墙设置

## 下一步

1. ✅ 本地开发和测试
2. 🔄 购买/准备Linux服务器
3. 🔄 部署服务端到公网服务器
4. 🔄 准备U盘，安装Agent到Mac Mini
