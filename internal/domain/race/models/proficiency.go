package models

import "github.com/google/uuid"

type Proficiency struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey" swaggertype:"string" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
}
