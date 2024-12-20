package entities

type Proficiency struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"unique;not null"`
	Description string `json:"description"`
}
