package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vsitnev/sync-manager/internal/dto"
	"github.com/vsitnev/sync-manager/internal/repository/pgdb"
	"github.com/vsitnev/sync-manager/internal/repository/repoerr"
	"github.com/vsitnev/sync-manager/internal/service"
	"net/http"
	"strconv"
)

type SourceRoutes struct {
	service service.Source
}

func newSourceRoutes(handler *gin.RouterGroup, sourceService service.Source) {
	r := &SourceRoutes{
		service: sourceService,
	}
	handler.POST("", r.saveSource)
	handler.GET("", r.getSources)
	handler.GET("/:id", r.getSourceByID)
	handler.PATCH("/:id", r.updateSource)
}

type sourceCreateRequestDto struct {
	Name          string      `json:"name" binding:"required"`
	Description   string      `json:"description" binding:"required"`
	Code          string      `json:"code" binding:"required"`
	ReceiveMethod string      `json:"receive_method" binding:"required"`
	Routes        []dto.Route `json:"routes" binding:"required,dive"`
}
type sourceRouteCreateRequestDto struct {
	Name string `json:"name" db:"name"`
	Url  string `json:"url" db:"url"`
}

type sourceCreateResponseDto struct {
	ID int
}

// @Summary Создание элемента "Источник"
// @Description Source
// @Tags Sources / Источники
// @Accept json
// @Produce json
// @Param input body sourceCreateRequestDto true "sources"
// @Success 201 {object} dto.Source
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/sources [post]
func (r *SourceRoutes) saveSource(c *gin.Context) {
	var input sourceCreateRequestDto

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := r.service.SaveSource(c.Request.Context(), dto.Source{
		Name:          input.Name,
		Description:   input.Description,
		Code:          input.Code,
		ReceiveMethod: input.ReceiveMethod,
		Routes:        input.Routes,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, sourceCreateResponseDto{
		ID: id,
	})
}

type sourceResponse struct {
	Data []dto.Source `json:"data"`
}

type sourceInput struct {
	SortType string `json:"sort_type,omitempty" form:"sort_type"`
	Offset   int    `json:"offset,omitempty" form:"offset"`
	Limit    int    `json:"limit,omitempty" form:"limit"`
}

// @Summary Список элементов "Источник"
// @Description Source
// @Tags Sources / Источники
// @Accept json
// @Produce json
// @Param input query sourceInput false "input"
// @Success 200 {object} sourceResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/sources [get]
func (r *SourceRoutes) getSources(c *gin.Context) {
	var input sourceInput
	if err := c.BindQuery(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	data, err := r.service.GetSources(c.Request.Context(), service.PaginationFilter{
		SortType: input.SortType,
		Limit:    input.Limit,
		Offset:   input.Offset,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(data) == 0 {
		data = []dto.Source{}
	}
	c.JSON(http.StatusOK, sourceResponse{
		Data: data,
	})
}

// @Summary Получение элемента "Источник" по id
// @Description Source
// @Tags Sources / Источники
// @Accept json
// @Produce json
// @Param id path int true "Source ID"
// @Success 200 {object} dto.Source
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/sources/{id} [get]
func (r *SourceRoutes) getSourceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}

	data, err := r.service.GetSourceByID(c, id)
	if err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, http.StatusText(500))
		return
	}
	c.JSON(http.StatusOK, data)
}

type updateSourceInput struct {
	Name          *string      `json:"name,omitempty"`
	Description   *string      `json:"description,omitempty"`
	Code          *string      `json:"code,omitempty"`
	ReceiveMethod *string      `json:"receive_method,omitempty"`
	Routes        *[]dto.Route `json:"routes,omitempty"`
}
type updateSourceResponse struct {
	Success bool `json:"success"`
}

// @Summary Изменение элемента "Источник" по id
// @Description Source
// @Tags Sources / Источники
// @Accept json
// @Produce json
// @Param input body updateSourceInput true "source"
// @Param id path int true "Source ID"
// @Success 200 {object} updateSourceInput
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/sources/{id} [patch]
func (r *SourceRoutes) updateSource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}
	var input updateSourceInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = r.service.UpdateSource(c, id, pgdb.UpdateSourceInput{
		Name:          input.Name,
		Description:   input.Description,
		Code:          input.Code,
		ReceiveMethod: input.ReceiveMethod,
		Routes:        input.Routes,
	})
	if err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, updateSourceResponse{
		Success: true,
	})
}
