package v1

import (
	"errors"
	"github.com/vsitnev/sync-manager/internal/dto"
	"github.com/vsitnev/sync-manager/internal/repository/repoerr"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/internal/service"
)

type MessageRoutes struct {
	service service.Message
}

func newMessageRoutes(handler *gin.RouterGroup, messageService service.Message) {
	r := &MessageRoutes{
		service: messageService,
	}
	handler.POST("", r.sendMessage)
	handler.GET("", r.getMessages)
	handler.GET("/:id", r.getMessageByID)
}

type messageCreateRequestDto struct {
	Routing string            `json:"routing" binding:"required"`
	Message model.AmqpMessage `json:"message" binding:"required"`
}

// @Summary Создание элемента "Сообщение"
// @Description Message
// @Tags Messages / Сообщения
// @Accept json
// @Produce json
// @Param input body messageCreateRequestDto true "messages"
// @Success 201 {object} service.SendMessageResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/messages [post]
func (r *MessageRoutes) sendMessage(c *gin.Context) {
	var input messageCreateRequestDto

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := r.service.SendMessage(c.Request.Context(), model.Message{
		Routing: input.Routing,
		Message: input.Message,
	})
	if err != nil {
		if errors.Is(err, service.ErrRouteUrlNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

type messageResponse struct {
	Data []dto.Message `json:"data"`
}

type messageInput struct {
	Source   string `json:"source,omitempty" form:"source"`
	Routing  string `json:"routing,omitempty" form:"routing"`
	SortType string `json:"sort_type,omitempty" form:"sort_type"`
	Offset   int    `json:"offset,omitempty" form:"offset"`
	Limit    int    `json:"limit,omitempty" form:"limit"`
}

// @Summary Список элементов "Сообщение"
// @Description Message
// @Tags Messages / Сообщения
// @Accept json
// @Produce json
// @Param input query messageInput false "input"
// @Success 200 {object} messageResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/messages [get]
func (r *MessageRoutes) getMessages(c *gin.Context) {
	var input messageInput
	if err := c.BindQuery(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	data, err := r.service.GetMessages(c.Request.Context(), service.MessageInput{
		Source:  input.Source,
		Routing: input.Routing,
		PaginationFilter: service.PaginationFilter{
			SortType: input.SortType,
			Limit:    input.Limit,
			Offset:   input.Offset,
		},
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(data) == 0 {
		data = []dto.Message{}
	}
	c.JSON(http.StatusOK, messageResponse{
		Data: data,
	})
}

// @Summary Получение элемента "Сообщение" по id
// @Description Сообщение
// @Tags Messages / Сообщения
// @Accept json
// @Produce json
// @Param id path int true "Message ID"
// @Success 200 {object} dto.Message
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/messages/{id} [get]
func (r *MessageRoutes) getMessageByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}

	message, err := r.service.GetMessageByID(c, id)
	if err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, http.StatusText(500))
		return
	}
	c.JSON(http.StatusOK, message)
}
