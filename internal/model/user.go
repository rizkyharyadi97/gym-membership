package model

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Username      string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email         string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password      string    `gorm:"size:255;not null" json:"-"`
	DepositAmount float64   `gorm:"type:numeric(12,2);not null;default:0" json:"deposit_amount"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
	Bookings     []Booking     `gorm:"foreignKey:UserID" json:"bookings,omitempty"`
}
