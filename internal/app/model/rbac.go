package model

import (
	"time"
)

type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null;unique" json:"name" binding:"required,min=2,max=100"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Resource struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null;unique" json:"name" binding:"required,min=2,max=100"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Action struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;not null;unique" json:"name" binding:"required,min=2,max=50"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Policy struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RoleID     uint      `gorm:"not null;uniqueIndex:unique_role_resource_action" json:"role_id" binding:"required"`
	Role       *Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	ResourceID uint      `gorm:"not null;uniqueIndex:unique_role_resource_action" json:"resource_id" binding:"required"`
	Resource   *Resource `gorm:"foreignKey:ResourceID" json:"resource,omitempty"`
	ActionID   uint      `gorm:"not null;uniqueIndex:unique_role_resource_action" json:"action_id" binding:"required"`
	Action     *Action   `gorm:"foreignKey:ActionID" json:"action,omitempty"`
	Effect     string    `gorm:"size:10;default:'allow';not null" json:"effect" binding:"required,oneof=allow deny"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserRole struct {
	UserID uint `gorm:"primaryKey" json:"user_id"`
	RoleID uint `gorm:"primaryKey" json:"role_id"`
}

// TableName overrides GORM's default plural naming for UserRole join table
func (UserRole) TableName() string {
	return "user_roles"
}
