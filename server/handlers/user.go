package handlers

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"mac-cluster-monitor/server/middleware"
	"mac-cluster-monitor/server/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// UserResponse 用户响应结构
type UserResponse struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Role        string `json:"role"`
	LicenseKey  string `json:"license_key"`
	DeviceCount int    `json:"device_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// GetUsers 获取所有用户列表 (admin专用)
func (h *UserHandler) GetUsers(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT u.id, u.username, u.email, u.phone, u.role, u.license_key, u.created_at, u.updated_at,
			   (SELECT COUNT(*) FROM devices WHERE user_id = u.id) as device_count
		FROM users u
		WHERE u.role != 'admin'
		ORDER BY u.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}
	defer rows.Close()

	var users []UserResponse
	for rows.Next() {
		var id int64
		var username, role, licenseKey string
		var email, phone sql.NullString
		var createdAt, updatedAt sql.NullTime
		var deviceCount int

		err := rows.Scan(&id, &username, &email, &phone, &role, &licenseKey, &createdAt, &updatedAt, &deviceCount)
		if err != nil {
			continue
		}

		user := UserResponse{
			ID:          id,
			Username:    username,
			Role:        role,
			LicenseKey:  licenseKey,
			DeviceCount: deviceCount,
		}

		if email.Valid {
			user.Email = email.String
		}
		if phone.Valid {
			user.Phone = phone.String
		}
		if createdAt.Valid {
			user.CreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
		}
		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time.Format("2006-01-02 15:04:05")
		}

		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{"users": users, "total": len(users)})
}

// CreateUser 创建用户 (admin专用, 卡密由手机号MD5生成)
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查手机号是否已存在
	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone = ?)", req.Phone).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "手机号已存在"})
		return
	}

	// 用手机号MD5生成卡密（取前8位大写）
	hash := md5.Sum([]byte(req.Phone))
	licenseKey := strings.ToUpper(hex.EncodeToString(hash[:])[:8])

	// 创建用户
	result, err := h.DB.Exec(
		"INSERT INTO users (username, password, phone, role, license_key) VALUES (?, ?, ?, ?, ?)",
		req.Username, licenseKey, req.Phone, req.Role, licenseKey,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败: " + err.Error()})
		return
	}

	userID, _ := result.LastInsertId()

	// 为普通用户初始化盈利记录
	if req.Role == "user" {
		_, _ = h.DB.Exec(
			"INSERT INTO user_profits (user_id, today_profit, total_profit, last_reset_date) VALUES (?, 0, 0, CURDATE())",
			userID,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户创建成功",
		"user": UserResponse{
			ID:         userID,
			Username:   req.Username,
			Phone:      req.Phone,
			Role:       req.Role,
			LicenseKey: licenseKey,
		},
	})
}

// UpdateUser 更新用户角色 (admin专用)
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.DB.Exec("UPDATE users SET role = ? WHERE id = ?", req.Role, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "角色更新成功"})
}

// BindDevice 绑定设备到用户 (admin专用)
func (h *UserHandler) BindDevice(c *gin.Context) {
	deviceID := c.Param("id")
	var req models.BindDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.DB.Exec("UPDATE devices SET user_id = ? WHERE id = ?", req.UserID, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "绑定设备失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设备绑定成功"})
}

// UnbindDevice 解绑设备 (admin专用)
func (h *UserHandler) UnbindDevice(c *gin.Context) {
	deviceID := c.Param("id")

	_, err := h.DB.Exec("UPDATE devices SET user_id = NULL WHERE id = ?", deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解绑设备失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设备解绑成功"})
}

// DeleteUser 删除用户 (admin专用, 不能删除自己)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := middleware.GetUserID(c)

	if userID == string(rune(currentUserID)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除当前登录用户"})
		return
	}

	_, err := h.DB.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// GetUserDashboard 获取普通用户Dashboard数据
func (h *UserHandler) GetUserDashboard(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取用户设备
	rows, err := h.DB.Query(`
		SELECT d.id, d.mac_address, d.hostname, d.ip_address, d.user_id, d.total_memory, d.cpu_cores,
		       d.status, d.last_seen_at, d.created_at,
		       dh.memory_used, dh.memory_total, dh.cpu_usage
		FROM devices d
		LEFT JOIN (
			SELECT device_id, memory_used, memory_total, cpu_usage
			FROM device_heartbeats
			WHERE id IN (
				SELECT MAX(id) FROM device_heartbeats GROUP BY device_id
			)
		) dh ON d.id = dh.device_id
		WHERE d.user_id = ?
		ORDER BY d.status DESC, d.last_seen_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询设备失败"})
		return
	}
	defer rows.Close()

	var devices []models.DeviceResponse
	onlineCount := 0
	offlineCount := 0

	for rows.Next() {
		var d models.DeviceResponse
		var lastSeenAt, createdAt sql.NullTime
		var ipAddress sql.NullString
		var userIDVal sql.NullInt64
		var totalMemory sql.NullInt64
		var cpuCores sql.NullInt32
		var memoryUsed, memoryTotal sql.NullInt64
		var cpuUsage sql.NullFloat64

		err := rows.Scan(
			&d.ID, &d.MacAddress, &d.Hostname, &ipAddress, &userIDVal, &totalMemory, &cpuCores,
			&d.Status, &lastSeenAt, &createdAt,
			&memoryUsed, &memoryTotal, &cpuUsage,
		)
		if err != nil {
			continue
		}

		d.IPAddress = ipAddress.String
		if userIDVal.Valid {
			d.UserID = userIDVal.Int64
		}
		d.TotalMemory = totalMemory.Int64
		d.CpuCores = cpuCores.Int32
		d.IsOnline = d.Status == 1
		if lastSeenAt.Valid {
			d.LastSeenAt = lastSeenAt.Time
		}
		if createdAt.Valid {
			d.CreatedAt = createdAt.Time
		}
		d.MemoryUsed = memoryUsed.Int64
		d.MemoryTotal = memoryTotal.Int64
		d.CpuUsage = cpuUsage.Float64

		if d.IsOnline {
			onlineCount++
		} else {
			offlineCount++
		}

		devices = append(devices, d)
	}

	// 获取盈利数据
	var todayProfit, totalProfit float64
	err = h.DB.QueryRow(
		"SELECT COALESCE(today_profit, 0), COALESCE(total_profit, 0) FROM user_profits WHERE user_id = ?",
		userID,
	).Scan(&todayProfit, &totalProfit)
	if err != nil {
		todayProfit = 0
		totalProfit = 0
	}

	c.JSON(http.StatusOK, models.UserDashboardResponse{
		DeviceCount:  len(devices),
		OnlineCount:  onlineCount,
		OfflineCount: offlineCount,
		TodayProfit:  todayProfit,
		TotalProfit:  totalProfit,
		Devices:      devices,
	})
}
