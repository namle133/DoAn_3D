package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uit/vinhomes-management/config"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/middleware"
	"github.com/uit/vinhomes-management/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UserID       string `json:"user_id"`
	Role         string `json:"role"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Name         string `json:"name"`
}

type RegisterRequest struct {
	Username         string `json:"username" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=8"`
	Role             string `json:"role" binding:"required,oneof=admin manager staff resident"`
	FullName         string `json:"full_name"`
	IDCard           string `json:"id_card"`
	DateOfBirth      string `json:"date_of_birth"`
	Gender           string `json:"gender"`
	PhoneNumber      string `json:"phone_number"`
	PermanentAddress string `json:"permanent_address"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func Login(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		loginID := req.Username
		if loginID == "" {
			loginID = req.Email
		}
		if loginID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email hoặc username là bắt buộc"})
			return
		}

		var user models.User
		if err := database.DB.Where("username = ? OR email = ?", loginID, loginID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email hoặc mật khẩu không đúng"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email hoặc mật khẩu không đúng"})
			return
		}

		if user.Status != "active" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Tài khoản đã bị khóa hoặc vô hiệu hóa"})
			return
		}

		token, err := middleware.GenerateToken(&user, cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": LoginResponse{
				Token:    token,
				UserID:   user.ID,
				Role:     user.Role,
				Username: user.Username,
				Email:    user.Email,
				Name:     user.Username,
			},
		})
	}
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Debug: log received data
	fmt.Printf("[DEBUG] Received registration request:\n")
	fmt.Printf("  Username: %q\n", req.Username)
	fmt.Printf("  Email: %q\n", req.Email)
	fmt.Printf("  Role: %q\n", req.Role)
	fmt.Printf("  FullName: %q\n", req.FullName)
	fmt.Printf("  IDCard: %q\n", req.IDCard)
	fmt.Printf("  DateOfBirth: %q\n", req.DateOfBirth)
	fmt.Printf("  Gender: %q\n", req.Gender)
	fmt.Printf("  PhoneNumber: %q\n", req.PhoneNumber)
	fmt.Printf("  PermanentAddress: %q\n", req.PermanentAddress)

	// For resident role, validate required fields
	if req.Role == "resident" {
		if req.FullName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Họ tên là bắt buộc"})
			return
		}
		if len(req.FullName) < 2 || len(req.FullName) > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Họ tên phải từ 2-100 ký tự"})
			return
		}
		if req.IDCard == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CCCD là bắt buộc"})
			return
		}
		if req.DateOfBirth == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ngày sinh là bắt buộc"})
			return
		}
		if req.Gender == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Giới tính là bắt buộc"})
			return
		}
		if req.Gender != "M" && req.Gender != "F" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Giới tính phải là M hoặc F"})
			return
		}
		if req.PhoneNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Số điện thoại là bắt buộc"})
			return
		}
		if req.PermanentAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Địa chỉ thường trú là bắt buộc"})
			return
		}
		if len(req.PermanentAddress) < 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Địa chỉ phải ít nhất 5 ký tự"})
			return
		}
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Người dùng đã tồn tại"})
		return
	}

	// Check if ID card already exists (for resident role)
	if req.Role == "resident" {
		var existingResident models.Resident
		if err := database.DB.Where("id_card = ?", req.IDCard).First(&existingResident).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "CCCD đã được đăng ký"})
			return
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mã hóa mật khẩu"})
		return
	}

	// Create user
	user := models.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		Status:   "active",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo người dùng"})
		return
	}

	// Create role-specific record
	createRoleRecord(user.ID, req.Role)

	// If resident role, create Resident record
	if req.Role == "resident" {
		dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			// Rollback user creation
			database.DB.Delete(&user)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng ngày sinh không hợp lệ (YYYY-MM-DD)"})
			return
		}

		resident := models.Resident{
			ID:               uuid.New().String(),
			UserID:           user.ID,
			FullName:         req.FullName,
			IDCard:           req.IDCard,
			DateOfBirth:      dateOfBirth,
			Gender:           req.Gender,
			PhoneNumber:      req.PhoneNumber,
			Email:            req.Email,
			PermanentAddress: req.PermanentAddress,
		}

		if err := database.DB.Create(&resident).Error; err != nil {
			// Rollback user and role record
			database.DB.Delete(&user)
			database.DB.Where("user_id = ?", user.ID).Delete(&models.Resident{})
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo hồ sơ cư dân"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Đăng ký thành công",
		"user_id": user.ID,
	})
}

// ForgotPassword — UC 2.7.4: tạo mã reset (demo: trả về token trong response)
func ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Không tiết lộ email có tồn tại hay không
		c.JSON(http.StatusOK, gin.H{"message": "Nếu email tồn tại, mã đặt lại mật khẩu đã được gửi"})
		return
	}

	tokenBytes := make([]byte, 16)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)
	expiry := time.Now().Add(30 * time.Minute)

	database.DB.Model(&user).Updates(map[string]interface{}{
		"reset_token":        token,
		"reset_token_expiry": expiry,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Mã đặt lại mật khẩu đã được tạo (demo). Vui lòng dùng mã bên dưới.",
		"data": gin.H{
			"reset_token": token,
			"expires_in":  "30 phút",
		},
	})
}

// ResetPassword — UC 2.7.4
func ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ? AND reset_token = ?", req.Email, req.Token).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã reset không hợp lệ hoặc đã hết hạn"})
		return
	}

	if user.ResetTokenExpiry == nil || user.ResetTokenExpiry.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã reset đã hết hạn"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	database.DB.Model(&user).Updates(map[string]interface{}{
		"password":           string(hashed),
		"reset_token":        "",
		"reset_token_expiry": nil,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Đặt lại mật khẩu thành công"})
}

func GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User

	if err := database.DB.Preload("Resident").Preload("Manager").Preload("Staff").Preload("Admin").
		First(&user, "id = ?", userID.(string)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateProfile — UC 2.7.5/2.7.6
func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		Email       string `json:"email"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID.(string)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if req.Email != "" {
		database.DB.Model(&user).Update("email", req.Email)
	}

	if req.NewPassword != "" {
		if req.OldPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cần nhập mật khẩu cũ"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Mật khẩu cũ không đúng"})
			return
		}
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		database.DB.Model(&user).Update("password", string(hashed))
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật hồ sơ thành công"})
}
