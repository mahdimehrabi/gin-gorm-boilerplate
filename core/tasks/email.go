package tasks

import (
	"boilerplate/apps/authApp/services"
	"boilerplate/core/infrastructure"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"time"
)

const (
	TypeSendVerifyEmail = "sendEmail:verify"
	TypeSendForgotEmail = "sendEmail:forgot"
)

type emailPayload struct {
	UserID uint
}

type EmailTask struct {
	Payload     emailPayload
	logger      infrastructure.Logger
	authService services.AuthService
}

func NewEmailTask(logger infrastructure.Logger, authService services.AuthService) EmailTask {
	return EmailTask{
		logger:      logger,
		authService: authService,
	}
}

func (et *EmailTask) NewVerifyEmailTask(userID uint) (*asynq.Task, error) {
	et.Payload.UserID = userID
	payload, err := json.Marshal(et.Payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		TypeSendVerifyEmail,
		payload,
		asynq.Timeout(80*time.Second),
		asynq.MaxRetry(2)), nil
}

func (et EmailTask) HandleVerifyEmailTask(ctx context.Context, t *asynq.Task) error {
	if err := json.Unmarshal(t.Payload(), &et.Payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	et.authService.SendRegisterationEmail(et.Payload.UserID)
	return nil
}
