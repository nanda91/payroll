package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy *uint          `json:"created_by,omitempty"`
	UpdatedBy *uint          `json:"updated_by,omitempty"`
	IPAddress string         `json:"ip_address,omitempty"`
	RequestID string         `json:"request_id,omitempty"`
}
