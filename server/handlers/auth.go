package handlers

import (
	"database/sql"
	"log"
	"mac-cluster-monitor/server/config"
	"mac-cluster-monitor/server/middleware"
	"mac-cluster-monitor/server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	DB  *sql.DB
	Cfg *config.Config
}

func NewAuthHandler(db *sql.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{DB: db, Cfg: cfg}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户（用手机号登录）
	var user models.User
	var email, phone sql.NullString
	err := h.DB.QueryRow(
		"SELECT id, username, password, email, phone, role, license_key, created_at, updated_at FROM users WHERE phone = ?",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Password, &email, &phone, &user.Role, &user.LicenseKey, &user.CreatedAt, &user.UpdatedAt)
	user.Email = email
	user.Phone = phone

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		log.Printf("查询用户失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	// 验证密码（明文比对）
	if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role, h.Cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User: models.User{
			ID:         user.ID,
			Username:   user.Username,
			Email:      email,
			Phone:      phone,
			Role:       user.Role,
			LicenseKey: user.LicenseKey,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		},
	})
}
