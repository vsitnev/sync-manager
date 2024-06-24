package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vsitnev/sync-manager/internal/dto"
	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/internal/repository"
	"github.com/vsitnev/sync-manager/internal/repository/pgdb"
	"github.com/vsitnev/sync-manager/pkg/amqpclient"
	"net/http"
	"strings"
)

type MessageService struct {
	repo          repository.Message
	sourceService *SourceService
	amqp          *amqpclient.Amqp
}

func NewMessageService(repo repository.Message, sourceService *SourceService, amqp *amqpclient.Amqp) *MessageService {
	return &MessageService{
		repo:          repo,
		sourceService: sourceService,
		amqp:          amqp,
	}
}

/*
Отправка через http
*/
/**
1. Find source
	not found -> send err
2. Get receive method:
	amqp -> publish
	http -> send

если мы гарантируем соблюдение контракта:
	http:
		400: return {save: "fail", sent: "fail", error: "...."}
		500: return {save: "success", sent: "fail", error: ""}
	amqp:
		{save: "success", sent: "success"}

	Далее делаем подписку на dead-letter-exchange -> из него обновляем dead = true
*/

type SendMessageResponse struct {
	Save  bool  `json:"save"`
	Sent  bool  `json:"sent"`
	Error error `json:"error"`
}

func sendMessageResponse(save, sent bool, err error) (SendMessageResponse, error) {
	return SendMessageResponse{
		Save:  save,
		Sent:  sent,
		Error: err,
	}, err
}

// handleMessageFromHttp
// handleMessageFromAmqp

func (s *MessageService) SendMessage(ctx context.Context, message model.Message) (SendMessageResponse, error) {
	source, err := s.sourceService.GetSourceByCode(ctx, message.Message.Source)
	if err != nil {
		return sendMessageResponse(false, false, err)
	}

	tx, err := s.repo.StartTx(ctx)
	if err != nil {
		return sendMessageResponse(false, false, err)
		//return msg, fmt.Errorf("MessageService.SendMessage - s.repo.StartTx: %v", err)
	}
	defer func() { tx.Rollback(ctx) }()

	data, err := s.repo.SaveMessage(ctx, message, tx)
	if err != nil {
		return sendMessageResponse(false, false, err)
		//return msg, fmt.Errorf("MessageService.SendMessage - s.repo.SaveMessage: %v", err)
	}
	msg := data.ToDto()

	switch source.ReceiveMethod {
	case "http":
		statusCode, err := s.sendByHttp(source, msg)
		if err != nil {
			return sendMessageResponse(false, false, err)
			//fmt.Errorf("MessageService.SendMessage - s.sendByHttp: %v", err),
		}
		if statusCode >= 400 && statusCode < 600 {
			if statusCode >= 400 && statusCode < 500 {
				return sendMessageResponse(false, false, calledServiceErrorResponse(statusCode))
			}
			dead := true
			err = s.repo.UpdateMessage(ctx, msg.ID, pgdb.UpdateMessageInput{Dead: &dead}, tx)
			if err != nil {
				return sendMessageResponse(false, false, calledServiceErrorResponse(statusCode))
			}
			err = tx.Commit(ctx)
			if err != nil {
				return sendMessageResponse(false, false, calledServiceErrorResponse(statusCode))
			}
			return sendMessageResponse(true, false, calledServiceErrorResponse(statusCode))
		}
	case "amqp":
		err = s.sendByAmqp(ctx, msg)
		if err != nil {
			return sendMessageResponse(false, false, err)
		}
	default:
		fmt.Println("There is no requested method of sending")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return sendMessageResponse(false, true, err)
		//return msg, fmt.Errorf("MessageService.SendMessage - tx.Commit: %v", err)
	}

	return sendMessageResponse(true, true, nil)
}

func (s *MessageService) getRouteUrlByName(routes []dto.Route, name string) (string, error) {
	for _, item := range routes {
		if item.Name == name {
			return item.Url, nil
		}
	}
	return "", errors.New("route URL not found")
}

func (s *MessageService) sendByHttp(source dto.Source, msg dto.Message) (int, error) {
	url, err := s.getRouteUrlByName(source.Routes, msg.Routing)
	if err != nil {
		return 0, fmt.Errorf("MessageService.sendByHttp - s.getRouteUrlByName: %v", err)
	}

	httpMsg, err := json.Marshal(msg.Message)
	if err != nil {
		return 0, fmt.Errorf("MessageService.sendByHttp - json.Marshal: %v", err)
	}

	reader := bytes.NewReader(httpMsg)
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		return 0, fmt.Errorf("MessageService.sendByHttp - http.Post: %v", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func (s *MessageService) sendByAmqp(ctx context.Context, msg dto.Message) error {
	amqpMsg, err := json.Marshal(msg.Message)
	if err != nil {
		return fmt.Errorf("MessageService.sendByAmqp - json.Marshal: %v", err)
	}
	exchange := fmt.Sprintf("%s.exc", strings.Split(msg.Routing, ".")[0])
	err = s.amqp.Publish(ctx, exchange, msg.Routing, amqpMsg)
	if err != nil {
		return fmt.Errorf("MessageService.sendByAmqp - s.amqp.Publish: %v", err)
	}
	return nil
}

func (s *MessageService) SaveMessages(ctx context.Context, messages []model.Message) error {
	return s.repo.SaveMessages(ctx, messages)
}

func (s *MessageService) GetMessages(ctx context.Context, input MessageInput) ([]dto.Message, error) {
	var msg []dto.Message
	data, err := s.repo.GetMessagesPagination(ctx, input.Source, input.Routing, input.SortType, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	for _, item := range data {
		msg = append(msg, item.ToDto())
	}
	return msg, nil
}

func (s *MessageService) GetMessageByID(ctx context.Context, ID int) (dto.Message, error) {
	var msg dto.Message
	data, err := s.repo.GetMessageByID(ctx, ID)
	if err != nil {
		return msg, err
	}
	return data.ToDto(), nil
}
