package handlers

import (
	"database/sql"
	"fmt"
	"mac-cluster-monitor/server/middleware"
	"mac-cluster-monitor/server/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	DB  *sql.DB
}

func NewDeviceHandler(db *sql.DB) *DeviceHandler {
	return &DeviceHandler{DB: db}
}

// Heartbeat 设备心跳上报
func (h *DeviceHandler) Heartbeat(c *gin.Context) {
	var req models.HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查设备是否已存在
	var deviceID int64
	err := h.DB.QueryRow("SELECT id FROM devices WHERE mac_address = ?", req.MacAddress).Scan(&deviceID)

	if err == sql.ErrNoRows {
		// 新设备，插入记录
		result, err := h.DB.Exec(
			"INSERT INTO devices (mac_address, hostname, ip_address, total_memory, cpu_cores, status, last_seen_at) VALUES (?, ?, ?, ?, ?, 1, NOW())",
			req.MacAddress, req.Hostname, req.IPAddress, req.MemoryTotal, req.CpuCores,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册设备失败"})
			return
		}
		deviceID, _ = result.LastInsertId()
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	} else {
		// 更新设备信息
		_, err = h.DB.Exec(
			"UPDATE devices SET hostname = ?, ip_address = ?, total_memory = ?, cpu_cores = ?, status = 1, last_seen_at = NOW() WHERE id = ?",
			req.Hostname, req.IPAddress, req.MemoryTotal, req.CpuCores, deviceID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新设备失败"})
			return
		}
	}

	// 记录心跳历史
	_, _ = h.DB.Exec(
		"INSERT INTO device_heartbeats (device_id, ip_address, memory_used, memory_total, cpu_usage) VALUES (?, ?, ?, ?, ?)",
		deviceID, req.IPAddress, req.MemoryUsed, req.MemoryTotal, req.CpuUsage,
	)

	c.JSON(http.StatusOK, gin.H{"message": "心跳上报成功", "device_id": deviceID})
}

// GetDevices 获取设备列表 (admin看全部, 普通用户只看自己的)
func (h *DeviceHandler) GetDevices(c *gin.Context) {
	role := middleware.GetUserRole(c)
	userID := middleware.GetUserID(c)

	// 构建查询
	query := `
		SELECT d.id, d.mac_address, d.hostname, d.ip_address, d.user_id, d.total_memory, d.cpu_cores,
		       d.status, d.last_seen_at, d.created_at,
		       dh.memory_used, dh.memory_total, dh.cpu_usage,
		       u.username
		FROM devices d
		LEFT JOIN (
			SELECT device_id, memory_used, memory_total, cpu_usage
			FROM device_heartbeats
			WHERE id IN (
				SELECT MAX(id) FROM device_heartbeats GROUP BY device_id
			)
		) dh ON d.id = dh.device_id
		LEFT JOIN users u ON d.user_id = u.id
		WHERE 1=1
	`
	var args []interface{}

	// 普通用户只能看自己的设备
	if role != "admin" {
		query += " AND d.user_id = ?"
		args = append(args, userID)
	}

	query += " ORDER BY d.status DESC, d.last_seen_at DESC"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询设备失败"})
		return
	}
	defer rows.Close()

	var devices []models.DeviceResponse
	for rows.Next() {
		var d models.DeviceResponse
		var lastSeenAt, createdAt sql.NullTime
		var ipAddress, username sql.NullString
		var userIDVal sql.NullInt64
		var totalMemory sql.NullInt64
		var cpuCores sql.NullInt32
		var memoryUsed, memoryTotal sql.NullInt64
		var cpuUsage sql.NullFloat64

		err := rows.Scan(
			&d.ID, &d.MacAddress, &d.Hostname, &ipAddress, &userIDVal, &totalMemory, &cpuCores,
			&d.Status, &lastSeenAt, &createdAt,
			&memoryUsed, &memoryTotal, &cpuUsage,
			&username,
		)
		if err != nil {
			continue
		}

		d.IPAddress = ipAddress.String
		if userIDVal.Valid {
			d.UserID = userIDVal.Int64
		}
		if username.Valid {
			d.Username = username.String
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

		devices = append(devices, d)
	}

	c.JSON(http.StatusOK, gin.H{"devices": devices, "total": len(devices)})
}

// GetDeviceStats 获取设备统计 (admin看全部, 普通用户只看自己的)
func (h *DeviceHandler) GetDeviceStats(c *gin.Context) {
	role := middleware.GetUserRole(c)
	userID := middleware.GetUserID(c)

	var whereClause string
	var args []interface{}

	if role != "admin" {
		whereClause = "WHERE user_id = ?"
		args = append(args, userID)
	}

	// 查询总数
	var total int
	err := h.DB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM devices %s", whereClause), args...).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "统计失败"})
		return
	}

	// 查询在线数
	var online int
	onlineQuery := whereClause
	if whereClause != "" {
		onlineQuery += " AND status = 1"
	} else {
		onlineQuery = "WHERE status = 1"
	}
	err = h.DB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM devices %s", onlineQuery), args...).Scan(&online)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "统计失败"})
		return
	}

	offline := total - online
	onlineRate := float64(0)
	if total > 0 {
		onlineRate = float64(online) / float64(total) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"total":       total,
		"online":      online,
		"offline":     offline,
		"online_rate": onlineRate,
	})
}

// DeleteDevice 删除设备 (仅admin)
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id := c.Param("id")

	_, err := h.DB.Exec("DELETE FROM devices WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除设备失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设备删除成功"})
}

// CheckOfflineDevices 检查离线设备 (定时任务调用)
func (h *DeviceHandler) CheckOfflineDevices() {
	timeout := 2 * time.Minute
	_, err := h.DB.Exec(
		"UPDATE devices SET status = 0 WHERE status = 1 AND last_seen_at < ?",
		time.Now().Add(-timeout),
	)
	if err != nil {
		// 记录日志
	}
}

// GetUnboundDevices 获取未绑定用户的设备列表 (admin专用, 用于分配设备)
func (h *DeviceHandler) GetUnboundDevices(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT d.id, d.mac_address, d.hostname, d.ip_address, d.status, d.last_seen_at
		FROM devices d
		WHERE d.user_id IS NULL
		ORDER BY d.status DESC, d.last_seen_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询设备失败"})
		return
	}
	defer rows.Close()

	type SimpleDevice struct {
		ID         int64  `json:"id"`
		MacAddress string `json:"mac_address"`
		Hostname   string `json:"hostname"`
		IPAddress  string `json:"ip_address"`
		Status     int    `json:"status"`
		IsOnline   bool   `json:"is_online"`
	}

	var devices []SimpleDevice
	for rows.Next() {
		var d SimpleDevice
		var ipAddress sql.NullString
		var lastSeenAt sql.NullTime

		err := rows.Scan(&d.ID, &d.MacAddress, &d.Hostname, &ipAddress, &d.Status, &lastSeenAt)
		if err != nil {
			continue
		}

		d.IPAddress = ipAddress.String
		d.IsOnline = d.Status == 1
		devices = append(devices, d)
	}

	c.JSON(http.StatusOK, gin.H{"devices": devices})
}

// UpdateProfit 定时更新所有用户的盈利 (由定时任务调用)
func (h *DeviceHandler) UpdateProfit(increment float64) {
	// 每日重置: 如果 last_reset_date 不是今天, 则重置 today_profit
	_, _ = h.DB.Exec(`
		UPDATE user_profits 
		SET today_profit = 0, last_reset_date = CURDATE() 
		WHERE last_reset_date < CURDATE()
	`)

	// 增加盈利
	_, _ = h.DB.Exec(`
		UPDATE user_profits 
		SET today_profit = today_profit + ?, total_profit = total_profit + ?
	`, increment, increment)
}
