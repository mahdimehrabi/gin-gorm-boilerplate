package tasks

import (
	"encoding/json"
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
	Payload emailPayload
}

func NewEmailTask(userID uint) EmailTask {
	return EmailTask{Payload: emailPayload{UserID: userID}}
}

func (et EmailTask) NewVerifyEmailTask(userID uint) (*asynq.Task, error) {
	payload, err := json.Marshal(et.Payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSendVerifyEmail, payload, asynq.Timeout(80*time.Second)), nil
}
