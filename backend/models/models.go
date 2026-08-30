package models

import "time"

// User is a per-password account. The username is deterministically derived
// from the user's password, so the same complex password always maps to the
// same account.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:64;uniqueIndex" json:"username"`
	PasswordHash string    `gorm:"size:255" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Countdown belongs to one user and holds an exact target time (second precision).
type Countdown struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index:idx_user_target,priority:1" json:"user_id"`
	Title      string    `gorm:"size:100;not null" json:"title"`
	TargetTime time.Time `gorm:"not null;index:idx_user_target,priority:2" json:"target_time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
