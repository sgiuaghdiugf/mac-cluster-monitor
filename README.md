# Mac Cluster Monitor

Mac Mini 集群设备监控系统

## 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                     公网 Linux 服务器                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Go Server  │  │   MySQL     │  │    Vue Web App      │  │
│  │  (API服务)   │  │   (数据)     │  │    (监控面板)        │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              ▲
                              │ HTTP API (心跳上报)
                              │
                    ┌─────────┴─────────┐
                    │     机房局域网      │
                    │  ┌─────────────┐   │
                    │  │  Mac Mini   │   │
                    │  │  + Go Agent │   │
                    │  │ (开机自运行)  │   │
                    │  └─────────────┘   │
                    └─────────────────────┘
```

## 功能特性

- **设备自动发现**: Mac Mini 开机后自动上报到监控系统
- **实时监控**: 显示设备在线状态、CPU使用率、内存使用情况
- **离线检测**: 10秒内未收到心跳自动标记为离线
- **用户管理**: 支持多用户登录注册
- **Web界面**: 美观的监控面板,支持实时刷新

## 监控指标

- 设备名称 (macOS主机名)
- MAC地址
- IP地址
- 内存使用 (已用/总量)
- CPU使用率
- 在线状态
- 最后在线时间

## 快速开始

### 1. 服务端部署 (Linux服务器)

#### 环境要求
- Linux 系统 (Ubuntu/CentOS等)
- MySQL 5.7+
- Go 1.21+

#### 部署步骤

```bash
# 1. 克隆代码
cd /opt
git clone <your-repo>
cd mac-cluster-monitor

# 2. 创建数据库
mysql -u root -p < database/schema.sql

# 3. 配置环境变量
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=mac_cluster_monitor
export JWT_SECRET=your-secret-key
export SERVER_PORT=8080

# 4. 编译并启动服务端
cd server
go mod tidy
go build -o server
./server
```

#### 使用 systemd 管理服务

```bash
# 复制服务文件
sudo cp scripts/mac-monitor-server.service /etc/systemd/system/

# 编辑配置
sudo vim /etc/systemd/system/mac-monitor-server.service
# 修改 WorkingDirectory 和 Environment 变量

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable mac-monitor-server
sudo systemctl start mac-monitor-server
```

### 2. 前端部署

```bash
cd web
npm install
npm run build

# 将 dist 目录部署到 Nginx 或其他Web服务器
```

#### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        root /opt/mac-cluster-monitor/web/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 3. Agent部署 (Mac Mini)

#### 方式一: 使用安装脚本

```bash
# 1. 在开发机上编译 Agent
# macOS ARM64 (M1/M2/M3/M4)
GOOS=darwin GOARCH=arm64 go build -o agent agent/main.go

# macOS AMD64 (Intel)
GOOS=darwin GOARCH=amd64 go build -o agent agent/main.go

# 2. 将 agent 和 install.sh 复制到U盘
# 3. 在Mac Mini上运行安装脚本
sudo ./install.sh http://your-server:8080
```

#### 方式二: 手动安装

```bash
# 1. 复制agent到系统目录
sudo cp agent /usr/local/bin/mac-cluster-monitor-agent
sudo chmod +x /usr/local/bin/mac-cluster-monitor-agent

# 2. 创建启动配置文件
sudo vim /Library/LaunchDaemons/com.macmonitor.agent.plist

# 3. 加载服务
sudo launchctl load /Library/LaunchDaemons/com.macmonitor.agent.plist
sudo launchctl start com.macmonitor.agent
```

## 环境变量配置

### Server 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 3306 |
| DB_USER | 数据库用户名 | root |
| DB_PASSWORD | 数据库密码 | - |
| DB_NAME | 数据库名 | mac_cluster_monitor |
| JWT_SECRET | JWT密钥 | - |
| SERVER_PORT | 服务端口 | 8080 |
| OFFLINE_THRESHOLD | 离线阈值(秒) | 10 |

### Agent 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| MONITOR_SERVER | 监控服务器地址 | http://localhost:8080 |
| MONITOR_INTERVAL | 上报间隔(秒) | 5 |

## API接口

### 公开接口

- `POST /api/register` - 用户注册
- `POST /api/login` - 用户登录
- `POST /api/heartbeat` - 设备心跳上报

### 需要认证的接口

- `GET /api/devices` - 获取设备列表
- `GET /api/devices/:id` - 获取单个设备
- `DELETE /api/devices/:id` - 删除设备

### 心跳上报数据格式

```json
{
  "mac_address": "aa:bb:cc:dd:ee:ff",
  "hostname": "mac-mini-01",
  "ip_address": "192.168.1.100",
  "memory_used": 8192,
  "memory_total": 16384,
  "cpu_usage": 45.5,
  "cpu_cores": 10
}
```

## 目录结构

```
mac-cluster-monitor/
├── agent/              # Mac Mini 监控客户端
│   ├── main.go
│   ├── install.sh      # 安装脚本
│   └── go.mod
├── server/             # 服务端
│   ├── main.go
│   ├── handlers/       # HTTP处理器
│   ├── middleware/     # 中间件
│   ├── models/         # 数据模型
│   ├── database/       # 数据库连接
│   ├── config/         # 配置
│   └── go.mod
├── web/                # Vue前端
│   ├── src/
│   ├── package.json
│   └── vite.config.js
├── database/           # 数据库脚本
│   └── schema.sql
├── scripts/            # 部署脚本
│   └── mac-monitor-server.service
└── README.md
```

## 开发计划

- [x] 基础心跳上报
- [x] 设备状态监控
- [x] Web管理界面
- [x] 用户认证
- [ ] 历史数据统计
- [ ] 告警通知
- [ ] 设备分组管理

## 许可证

MIT License
