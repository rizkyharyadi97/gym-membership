package model

import "time"

type Trainer struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	Specialization string    `gorm:"size:100;not null" json:"specialization"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations
	GymClasses []GymClass `gorm:"foreignKey:TrainerID" json:"gym_classes,omitempty"`
}
