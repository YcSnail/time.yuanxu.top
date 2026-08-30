package handlers

import (
	"net/http"
	"time"

	"github.com/YcSnail/time.yuanxu.top/backend/models"
	"github.com/YcSnail/time.yuanxu.top/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler handles password-based enter (register-or-login).
type AuthHandler struct {
	DB     *gorm.DB
	Secret string
	Expire int // token lifetime in days
}

type enterRequest struct {
	Password string `json:"password" binding:"required"`
}

// Enter implements: same password => same account. First use creates the
// account (password complexity enforced), later uses log in.
func (h *AuthHandler) Enter(c *gin.Context) {
	var req enterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入密码"})
		return
	}

	username := utils.UsernameFromPassword(req.Password)
	var user models.User
	err := h.DB.Where("username = ?", username).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		// New account: enforce the password policy before creating.
		if ok, reason := utils.ValidatePassword(req.Password); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": reason})
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "系统繁忙,请稍后再试"})
			return
		}
		user = models.User{Username: username, PasswordHash: string(hash)}
		if err := h.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建账号失败"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统繁忙,请稍后再试"})
		return
	} else if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	token, err := h.issueToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统繁忙,请稍后再试"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// Me returns the current logged-in user.
func (h *AuthHandler) Me(c *gin.Context) {
	uid := c.GetUint("uid")
	var user models.User
	if err := h.DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) issueToken(uid uint) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"uid": uid,
		"iat": now.Unix(),
		"exp": now.Add(time.Duration(h.Expire) * 24 * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.Secret))
}
