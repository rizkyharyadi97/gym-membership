package model

import "time"

type GymClass struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TrainerID    uint      `gorm:"not null;index" json:"trainer_id"`
	ClassName    string    `gorm:"size:100;not null" json:"class_name"`
	Schedule     time.Time `gorm:"not null;index" json:"schedule"`
	Quota        int       `gorm:"not null" json:"quota"`
	Availability bool      `gorm:"not null;default:true" json:"availability"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	Trainer  Trainer   `gorm:"foreignKey:TrainerID" json:"trainer,omitempty"`
	Bookings []Booking `gorm:"foreignKey:GymClassID" json:"bookings,omitempty"`
}
