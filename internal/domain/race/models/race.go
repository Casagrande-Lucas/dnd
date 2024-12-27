package models

type Race struct {
	ID                  uint                `json:"id" gorm:"primaryKey"`
	Name                string              `json:"name" gorm:"unique;not null"`
	Description         string              `json:"description"`
	AbilityScoreBonuses AbilityScoreBonuses `json:"ability_score_bonuses" gorm:"embedded"`
	Age                 Age                 `json:"age" gorm:"foreignKey:RaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
