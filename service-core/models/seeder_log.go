package models

import "time"

type SeederLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SeederName string    `json:"seeder_name" gorm:"not null"`
	ExecutedAt time.Time `json:"executed_at" gorm:"default:CURRENT_TIMESTAMP"`
}
