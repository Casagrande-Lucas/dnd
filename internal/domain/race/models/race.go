package models

import "github.com/google/uuid"

type Race struct {
	ID                  uuid.UUID           `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey" swaggertype:"string" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name                string              `json:"name" gorm:"unique;not null"`
	Description         string              `json:"description"`
	AbilityScoreBonuses AbilityScoreBonuses `json:"ability_score_bonuses" gorm:"embedded"`
	Age                 Age                 `json:"age" gorm:"foreignKey:RaceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Size                string              `json:"size"`
	Speed               int8                `json:"speed"`
	Alignment           string              `json:"alignment,omitempty"`
	Proficiencies       []Proficiency       `json:"proficiencies,omitempty" gorm:"many2many:race_proficiencies;"`
	LanguagesKnown      []Language          `json:"languages_known,omitempty" gorm:"many2many:race_languages;"`
	Traits              []Trait             `json:"traits,omitempty" gorm:"many2many:race_traits;"`
	Subraces            []Subrace           `json:"subraces,omitempty" gorm:"foreignKey:RaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type AbilityScoreBonuses struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}
