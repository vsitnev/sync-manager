package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/vsitnev/sync-manager/internal/dto/request"
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
}


// @Summary Создание элемента "Сообщение"
// @Description Message
// @Tags Message / Сообщение
// @Accept json
// @Produce json
// @Param input body dto.MessageRequestDto true "messages"
// @Success 201 {object} v1.MessageRoutes.create.res
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/messages [post]
func (r *MessageRoutes) sendMessage(c *gin.Context) {
	var input dto.MessageRequestDto

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := r.service.SendMessage(c.Request.Context(), input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	type res struct {
		Id int `json:"id"`
	}
	c.JSON(http.StatusOK, res{
		Id: id,
	})
}

type messageResponse struct {
	Data []model.Message `json:"data"`
}

// @Summary Список элементов "Сообщение"
// @Description Message
// @Tags Message / Сообщение
// @Accept json
// @Produce json
// @Success 200 {object} messageResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/messages [get]
func (r *MessageRoutes) getMessages(c *gin.Context) {
	data, err := r.service.GetMessages(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("data: ", data)

	
	if len(data) == 0 {
		data = []model.Message{}
	}
	c.JSON(http.StatusOK, messageResponse{
		Data: data,
	})
}
