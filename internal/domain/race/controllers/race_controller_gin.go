package controllers

import (
	"net/http"

	"github.com/Casagrande-Lucas/dnd/internal/domain/race/models"
	"github.com/Casagrande-Lucas/dnd/internal/domain/race/services"
	"github.com/Casagrande-Lucas/dnd/pkg/httperror"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// GetAllRaces godoc
// @Summary      List all races
// @Description  Return all registered races
// @Tags         Races
// @Accept       json
// @Produce      json
// @Success      200 {array}  models.Race
// @Failure      500 {object} httperror.ErrorResponse
// @Router       /races [get]
func (c *raceControllerGin) GetAllRaces(ctx *gin.Context) {
	races, err := c.service.ListRaces()
	if err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}
	ctx.JSON(http.StatusOK, races)
}

// GetRaceByID godoc
// @Summary      Get race by ID
// @Description  Retrieve a race using the provided ID
// @Tags         Races
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Race ID (UUID)"
// @Success      200  {object}  models.Race
// @Failure      400  {object}  httperror.ErrorResponse
// @Failure      404  {object}  httperror.ErrorResponse
// @Router       /races/{id} [get]
func (c *raceControllerGin) GetRaceByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}

	race, err := c.service.GetRaceDetails(id)
	if err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}
	ctx.JSON(http.StatusOK, race)
}

// CreateRace godoc
// @Summary      Create race
// @Description  Create a new race
// @Tags         Races
// @Accept       json
// @Produce      json
// @Param        race  body      models.Race  true  "Race info"
// @Success      201   {object}  models.Race
// @Failure      400   {object}  httperror.ErrorResponse
// @Failure      500   {object}  httperror.ErrorResponse
// @Router       /races [post]
func (c *raceControllerGin) CreateRace(ctx *gin.Context) {
	var race models.Race
	if err := ctx.ShouldBindJSON(&race); err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}

	if err := c.service.RegisterRace(&race); err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}
	ctx.JSON(http.StatusCreated, race)
}

// UpdateRace godoc
// @Summary      Update race
// @Description  Update an existing race
// @Tags         Races
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "Race ID (UUID)"
// @Param        race  body      models.Race  true  "Race info"
// @Success      200   {object}  models.Race
// @Failure      400   {object}  httperror.ErrorResponse
// @Failure      404   {object}  httperror.ErrorResponse
// @Router       /races/{id} [put]
func (c *raceControllerGin) UpdateRace(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}

	var race models.Race
	if err := ctx.ShouldBindJSON(&race); err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}

	if err := c.service.UpdateRaceInfo(id, &race); err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}
	ctx.JSON(http.StatusOK, race)
}

// DeleteRace godoc
// @Summary      Delete race
// @Description  Delete an existing race
// @Tags         Races
// @Param        id   path      string  true  "Race ID (UUID)"
// @Success      204
// @Failure      400 {object} httperror.ErrorResponse
// @Failure      404 {object} httperror.ErrorResponse
// @Router       /races/{id} [delete]
func (c *raceControllerGin) DeleteRace(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}

	if err := c.service.RemoveRace(id); err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// AddSubrace godoc
// @Summary      Add subrace
// @Description  Add a new subrace to an existing race
// @Tags         Races
// @Accept       json
// @Produce      json
// @Param        id       path      string          true  "Race ID (UUID)"
// @Param        subrace  body      models.Subrace  true  "Subrace info"
// @Success      201 {object} models.Subrace
// @Failure      400 {object} httperror.ErrorResponse
// @Failure      500 {object} httperror.ErrorResponse
// @Router       /races/{id}/subraces [post]
func (c *raceControllerGin) AddSubrace(ctx *gin.Context) {
	idStr := ctx.Param("id")
	raceID, err := uuid.Parse(idStr)
	if err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}

	var subrace models.Subrace
	if err := ctx.ShouldBindJSON(&subrace); err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}

	if err := c.service.AddSubraceToRace(raceID, &subrace); err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}
	ctx.JSON(http.StatusCreated, subrace)
}

// RemoveSubrace godoc
// @Summary      Remove subrace
// @Description  Remove an existing subrace from a race
// @Tags         Races
// @Accept       json
// @Produce      json
// @Param        id         path  string true "Race ID (UUID)"
// @Param        subraceID  path  string true "Subrace ID (UUID)"
// @Success      204
// @Failure      400 {object} httperror.ErrorResponse
// @Failure      404 {object} httperror.ErrorResponse
// @Router       /races/{id}/subraces/{subraceID} [delete]
func (c *raceControllerGin) RemoveSubrace(ctx *gin.Context) {
	idStr := ctx.Param("id")
	subraceIDStr := ctx.Param("subraceID")

	raceID, err := uuid.Parse(idStr)
	if err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}
	subraceID, err := uuid.Parse(subraceIDStr)
	if err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}

	if err := c.service.DetachSubraceFromRace(raceID, subraceID); err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// AddTrait godoc
// @Summary      Add trait to race
// @Description  Add a new trait to an existing race
// @Tags         Races
// @Accept       json
// @Produce      json
// @Param        id       path  string  true  "Race ID (UUID)"
// @Param        traitID  path  string  true  "Trait ID (UUID)"
// @Success      201
// @Failure      400 {object} httperror.ErrorResponse
// @Failure      404 {object} httperror.ErrorResponse
// @Router       /races/{id}/traits/{traitID} [post]
func (c *raceControllerGin) AddTrait(ctx *gin.Context) {
	idStr := ctx.Param("id")
	traitIDStr := ctx.Param("traitID")

	raceID, err := uuid.Parse(idStr)
	if err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}
	traitID, err := uuid.Parse(traitIDStr)
	if err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}

	if err := c.service.AssignTraitToRace(raceID, traitID); err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}
	ctx.Status(http.StatusCreated)
}

// RemoveTrait godoc
// @Summary      Remove trait from race
// @Description  Remove an existing trait from a race
// @Tags         Races
// @Accept       json
// @Produce      json
// @Param        id       path  string  true  "Race ID (UUID)"
// @Param        traitID  path  string  true  "Trait ID (UUID)"
// @Success      204
// @Failure      400 {object} httperror.ErrorResponse
// @Failure      404 {object} httperror.ErrorResponse
// @Router       /races/{id}/traits/{traitID} [delete]
func (c *raceControllerGin) RemoveTrait(ctx *gin.Context) {
	idStr := ctx.Param("id")
	traitIDStr := ctx.Param("traitID")

	raceID, err := uuid.Parse(idStr)
	if err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}
	traitID, err := uuid.Parse(traitIDStr)
	if err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}

	if err := c.service.UnassignTraitFromRace(raceID, traitID); err != nil {
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// SearchRaces godoc
// @Summary      Search races
// @Description  Search for races based on query parameters
// @Tags         Races
// @Accept       json
// @Produce      json
// @Param        key    query  string false "Key to filter"
// @Param        value  query  string false "Value to filter"
// @Success      200 {array}  models.Race
// @Failure      400 {object} httperror.ErrorResponse
// @Failure      500 {object} httperror.ErrorResponse
// @Router       /races/search [get]
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
		apiError := httperror.FormError(err)
		ctx.JSON(apiError.StatusCode, apiError.ObjectErr)
		return
	}
	ctx.JSON(http.StatusOK, races)
}
