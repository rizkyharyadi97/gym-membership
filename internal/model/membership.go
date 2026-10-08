package model

import "time"

type Membership struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Costs        float64   `gorm:"type:numeric(12,2);not null" json:"costs"`
	DurationDays int       `gorm:"not null" json:"duration_days"`
	Availability bool      `gorm:"not null;default:true" json:"availability"`
	Category     string    `gorm:"size:50;not null" json:"category"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	Transactions []Transaction `gorm:"foreignKey:MembershipID" json:"transactions,omitempty"`
}
