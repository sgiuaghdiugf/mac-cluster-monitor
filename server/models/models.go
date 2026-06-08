package models

import (
	"database/sql"
	"time"
)

// User 用户模型
type User struct {
	ID         int64          `json:"id"`
	Username   string         `json:"username"`
	Password   string         `json:"-"`
	Email      sql.NullString `json:"email"`
	Phone      sql.NullString `json:"phone"`
	Role       string         `json:"role"`
	LicenseKey string         `json:"license_key"`
	DeviceCount int            `json:"device_count,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// Device 设备模型
type Device struct {
	ID          int64          `json:"id"`
	MacAddress  string         `json:"mac_address"`
	Hostname    string         `json:"hostname"`
	IPAddress   sql.NullString `json:"ip_address"`
	UserID      sql.NullInt64  `json:"user_id"`
	TotalMemory sql.NullInt64  `json:"total_memory"`
	CpuCores    sql.NullInt32  `json:"cpu_cores"`
	Status      int            `json:"status"` // 0=离线, 1=在线
	LastSeenAt  sql.NullTime   `json:"last_seen_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// 非数据库字段,用于展示
	MemoryUsed  int64   `json:"memory_used,omitempty"`
	MemoryTotal int64   `json:"memory_total,omitempty"`
	CpuUsage    float64 `json:"cpu_usage,omitempty"`
	IsOnline    bool    `json:"is_online"`
}

// HeartbeatRequest 心跳请求
type HeartbeatRequest struct {
	MacAddress  string  `json:"mac_address" binding:"required"`
	Hostname    string  `json:"hostname" binding:"required"`
	IPAddress   string  `json:"ip_address"`
	MemoryUsed  int64   `json:"memory_used"`
	MemoryTotal int64   `json:"memory_total"`
	CpuUsage    float64 `json:"cpu_usage"`
	CpuCores    int     `json:"cpu_cores"`
}

// DeviceResponse 设备响应
type DeviceResponse struct {
	ID          int64     `json:"id"`
	MacAddress  string    `json:"mac_address"`
	Hostname    string    `json:"hostname"`
	IPAddress   string    `json:"ip_address"`
	UserID      int64     `json:"user_id"`
	Username    string    `json:"username,omitempty"`
	TotalMemory int64     `json:"total_memory"`
	CpuCores    int32     `json:"cpu_cores"`
	Status      int       `json:"status"`
	IsOnline    bool      `json:"is_online"`
	LastSeenAt  time.Time `json:"last_seen_at"`
	CreatedAt   time.Time `json:"created_at"`

	// 实时数据
	MemoryUsed  int64   `json:"memory_used"`
	MemoryTotal int64   `json:"memory_total"`
	CpuUsage    float64 `json:"cpu_usage"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// CreateUserRequest 创建用户请求 (admin专用)
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=1,max=50"`
	Phone    string `json:"phone" binding:"required"` // 手机号，用于登录和生成卡密
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Role string `json:"role" binding:"required,oneof=admin user"`
}

// BindDeviceRequest 绑定设备到用户请求
type BindDeviceRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

// UserProfit 用户盈利
type UserProfit struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	TodayProfit   float64   `json:"today_profit"`
	TotalProfit   float64   `json:"total_profit"`
	LastResetDate string    `json:"last_reset_date"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// UserDashboardResponse 普通用户Dashboard响应
type UserDashboardResponse struct {
	DeviceCount int          `json:"device_count"`
	OnlineCount int          `json:"online_count"`
	OfflineCount int         `json:"offline_count"`
	TodayProfit float64      `json:"today_profit"`
	TotalProfit float64      `json:"total_profit"`
	Devices     []DeviceResponse `json:"devices"`
}
