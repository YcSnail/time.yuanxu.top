package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/YcSnail/time.yuanxu.top/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CountdownHandler manages a user's countdown items.
type CountdownHandler struct {
	DB *gorm.DB
}

// List returns the current user's countdowns, nearest target time first.
func (h *CountdownHandler) List(c *gin.Context) {
	uid := c.GetUint("uid")
	var items []models.Countdown
	h.DB.Where("user_id = ?", uid).Order("target_time ASC").Find(&items)
	c.JSON(http.StatusOK, gin.H{"items": items})
}

type createRequest struct {
	Title      string `json:"title" binding:"required"`
	TargetTime string `json:"target_time" binding:"required"`
}

// Create adds a new countdown. TargetTime accepts RFC3339 or "2006-01-02 15:04:05".
func (h *CountdownHandler) Create(c *gin.Context) {
	uid := c.GetUint("uid")

	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写标题和目标时间"})
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题不能为空"})
		return
	}
	if utf8.RuneCountInString(title) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题不能超过 100 个字符"})
		return
	}

	target, err := parseTargetTime(req.TargetTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "时间格式错误,请选择有效时间"})
		return
	}
	if !target.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "目标时间必须晚于当前时间"})
		return
	}

	item := models.Countdown{UserID: uid, Title: title, TargetTime: target}
	if err := h.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败,请稍后再试"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": item})
}

// Delete removes a countdown owned by the current user.
func (h *CountdownHandler) Delete(c *gin.Context) {
	uid := c.GetUint("uid")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的倒计时"})
		return
	}

	res := h.DB.Where("id = ? AND user_id = ?", id, uid).Delete(&models.Countdown{})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败,请稍后再试"})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "倒计时不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func parseTargetTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
}
