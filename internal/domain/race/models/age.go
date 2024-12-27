package models

type Age struct {
	RaceID          uint   `json:"race_id" gorm:"primaryKey"`
	AverageLifespan string `json:"average_lifespan"`
	MinimumAge      int    `json:"minimum_age"`
	MaximumAge      int    `json:"maximum_age"`
}
