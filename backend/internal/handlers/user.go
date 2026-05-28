package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/uit/vinhomes-management/internal/database"
	"github.com/uit/vinhomes-management/internal/models"
)

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin manager staff resident"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Role     string `json:"role" binding:"omitempty,oneof=admin manager staff resident"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive locked"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

func toUserResponse(u models.User) UserResponse {
	return UserResponse{
		ID: u.ID, Username: u.Username, Email: u.Email,
		Role: u.Role, Status: u.Status, CreatedAt: u.CreatedAt,
	}
}

// GetUsers — UC 2.7.7 Admin quản lý tài khoản
func GetUsers(c *gin.Context) {
	var users []models.User
	q := database.DB
	if role := c.Query("role"); role != "" {
		q = q.Where("role = ?", role)
	}
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	resp := make([]UserResponse, len(users))
	for i, u := range users {
		resp[i] = toUserResponse(u)
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existing models.User
	if err := database.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username hoặc email đã tồn tại"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user := models.User{
		ID: uuid.New().String(), Username: req.Username, Email: req.Email,
		Password: string(hashed), Role: req.Role, Status: "active",
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	createRoleRecord(user.ID, user.Role)
	c.JSON(http.StatusCreated, gin.H{"data": toUserResponse(user)})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	updates := map[string]interface{}{}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updates["password"] = string(hashed)
	}
	if len(updates) > 0 {
		database.DB.Model(&user).Updates(updates)
		database.DB.First(&user, "id = ?", id)
	}
	c.JSON(http.StatusOK, gin.H{"data": toUserResponse(user)})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	if id == userID.(string) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể xóa tài khoản của chính mình"})
		return
	}
	if err := database.DB.Model(&models.User{}).Where("id = ?", id).Update("status", "inactive").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Đã vô hiệu hóa tài khoản"})
}

func createRoleRecord(userID, role string) {
	switch role {
	case "admin":
		database.DB.Create(&models.Admin{ID: uuid.New().String(), UserID: userID})
	case "manager":
		database.DB.Create(&models.Manager{ID: uuid.New().String(), UserID: userID})
	case "staff":
		database.DB.Create(&models.Staff{ID: uuid.New().String(), UserID: userID})
	case "resident":
		database.DB.Create(&models.Resident{ID: uuid.New().String(), UserID: userID})
	}
}
