package model

import "time"

// User 用户实体
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // 不序列化密码
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`   // 角色：admin, staff, store, customer
	Status    int       `json:"status"` // 状态：1-正常，0-禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
