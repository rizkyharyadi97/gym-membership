package model

import "time"

type Booking struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	GymClassID  uint      `gorm:"not null;index" json:"gym_class_id"`
	BookingDate time.Time `gorm:"not null;index" json:"booking_date"`
	Status      string    `gorm:"size:20;not null;default:booked;index" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	GymClass GymClass `gorm:"foreignKey:GymClassID" json:"gym_class,omitempty"`
}
