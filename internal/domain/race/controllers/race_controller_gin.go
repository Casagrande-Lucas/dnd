package controllers

import (
	"net/http"
	"strconv"

	"github.com/Casagrande-Lucas/dnd/internal/domain/race/entities"
	"github.com/Casagrande-Lucas/dnd/internal/domain/race/services"
	"github.com/gin-gonic/gin"
)

// raceControllerGin is a concrete implementation of RaceController using the Gin framework.
type raceControllerGin struct {
	service services.RaceService
}

// NewRaceControllerGin creates a new instance of raceControllerGin.
func NewRaceControllerGin(service services.RaceService) RaceController {
	return &raceControllerGin{
		service: service,
	}
}

// GetAllRaces handles GET /races to retrieve all races.
func (c *raceControllerGin) GetAllRaces(ctx *gin.Context) {
	races, err := c.service.ListRaces()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, races)
}

// GetRaceByID handles GET /races/:id to retrieve a race by its ID.
func (c *raceControllerGin) GetRaceByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid race ID"})
		return
	}

	race, err := c.service.GetRaceDetails(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, race)
}

// CreateRace handles POST /races to create a new race.
func (c *raceControllerGin) CreateRace(ctx *gin.Context) {
	var race entities.Race
	if err := ctx.ShouldBindJSON(&race); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if err := c.service.RegisterRace(&race); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, race)
}

// UpdateRace handles PUT /races/:id to update an existing race.
func (c *raceControllerGin) UpdateRace(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid race ID"})
		return
	}

	var race entities.Race
	if err := ctx.ShouldBindJSON(&race); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if err := c.service.UpdateRaceInfo(uint(id), &race); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, race)
}

// DeleteRace handles DELETE /races/:id to remove a race.
func (c *raceControllerGin) DeleteRace(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid race ID"})
		return
	}

	if err := c.service.RemoveRace(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// AddSubrace handles POST /races/:id/subraces to add a subrace to a race.
func (c *raceControllerGin) AddSubrace(ctx *gin.Context) {
	idStr := ctx.Param("id")
	raceID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid race ID"})
		return
	}

	var subrace entities.Subrace
	if err := ctx.ShouldBindJSON(&subrace); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if err := c.service.AddSubraceToRace(uint(raceID), &subrace); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, subrace)
}

// RemoveSubrace handles DELETE /races/:id/subraces/:subraceID to remove a subrace from a race.
func (c *raceControllerGin) RemoveSubrace(ctx *gin.Context) {
	idStr := ctx.Param("id")
	subraceIDStr := ctx.Param("subraceID")

	raceID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid race ID"})
		return
	}
	subraceID, err := strconv.Atoi(subraceIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid subrace ID"})
		return
	}

	if err := c.service.DetachSubraceFromRace(uint(raceID), uint(subraceID)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// AddTrait handles POST /races/:id/traits/:traitID to assign a trait to a race.
func (c *raceControllerGin) AddTrait(ctx *gin.Context) {
	idStr := ctx.Param("id")
	traitIDStr := ctx.Param("traitID")

	raceID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid race ID"})
		return
	}
	traitID, err := strconv.Atoi(traitIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trait ID"})
		return
	}

	if err := c.service.AssignTraitToRace(uint(raceID), uint(traitID)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusCreated)
}

// RemoveTrait handles DELETE /races/:id/traits/:traitID to unassign a trait from a race.
func (c *raceControllerGin) RemoveTrait(ctx *gin.Context) {
	idStr := ctx.Param("id")
	traitIDStr := ctx.Param("traitID")

	raceID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid race ID"})
		return
	}
	traitID, err := strconv.Atoi(traitIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trait ID"})
		return
	}

	if err := c.service.UnassignTraitFromRace(uint(raceID), uint(traitID)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// SearchRaces handles GET /races/search?key=value to search for races based on criteria.
func (c *raceControllerGin) SearchRaces(ctx *gin.Context) {
	criteria := make(map[string]string)
	queryParams := ctx.Request.URL.Query()
	for key, values := range queryParams {
		if len(values) > 0 {
			criteria[key] = values[0]
		}
	}

	races, err := c.service.FindRaces(criteria)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, races)
}
