package controllers

import (
	"github.com/gin-gonic/gin"
)

type RaceController interface {
	GetAllRaces(ctx *gin.Context)
	GetRaceByID(ctx *gin.Context)
	CreateRace(ctx *gin.Context)
	UpdateRace(ctx *gin.Context)
	DeleteRace(ctx *gin.Context)
	AddSubrace(ctx *gin.Context)
	RemoveSubrace(ctx *gin.Context)
	AddTrait(ctx *gin.Context)
	RemoveTrait(ctx *gin.Context)
	SearchRaces(ctx *gin.Context)
}
