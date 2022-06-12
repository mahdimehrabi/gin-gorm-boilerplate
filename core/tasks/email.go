package tasks

import (
	"boilerplate/apps/authApp/services"
	"boilerplate/apps/userApp/repositories"
	"boilerplate/core/infrastructure"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"strconv"
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
	Payload        emailPayload
	logger         infrastructure.Logger
	authService    services.AuthService
	userRepository repositories.UserRepository
	email          infrastructure.Email
	env            infrastructure.Env
}

func NewEmailTask(logger infrastructure.Logger,
	authService services.AuthService,
	userRepository repositories.UserRepository,
	email infrastructure.Email,
	env infrastructure.Env,
) EmailTask {
	return EmailTask{
		logger:         logger,
		authService:    authService,
		userRepository: userRepository,
		email:          email,
		env:            env,
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
	return et.SendRegisterationEmail(et.Payload.UserID)
}

func (et EmailTask) SendRegisterationEmail(userID uint) error {
	user, err := et.userRepository.FindByField("id", strconv.Itoa(int(userID)))
	if err != nil {
		et.logger.Zap.Error("failed to find user:", err)
		return err
	}
	ch := make(chan error)
	htmlFile := et.env.BasePath + "/vendors/templates/mail/auth/register.tmpl"

	data := map[string]string{
		"name": user.FirstName,
		"link": et.env.SiteUrl + "/verify-email?token=" + user.VerifyEmailToken,
	}
	go et.email.SendEmail(ch, user.Email, "Verify email", htmlFile, data)
	err = <-ch
	if err != nil {
		et.logger.Zap.Error(err)
		return err
	}
	err = et.userRepository.UpdateColumn(&user, "last_verify_email_date", time.Now())
	if err != nil {
		et.logger.Zap.Error(err)
		return err
	}
	return nil
}

func (et *EmailTask) NewForgotEmailTask(userID uint) (*asynq.Task, error) {
	et.Payload.UserID = userID
	payload, err := json.Marshal(et.Payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		TypeSendForgotEmail,
		payload,
		asynq.Timeout(80*time.Second),
		asynq.MaxRetry(2)), nil
}

func (et EmailTask) HandleForgotEmailTask(ctx context.Context, t *asynq.Task) error {
	if err := json.Unmarshal(t.Payload(), &et.Payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	return et.sendForgotPassowrdEmail(et.Payload.UserID)
}

func (et EmailTask) sendForgotPassowrdEmail(userID uint) error {
	user, err := et.userRepository.FindByField("id", strconv.Itoa(int(userID)))
	if err != nil {
		et.logger.Zap.Error("failed to find user:", err)
		return err
	}

	ch := make(chan error)
	htmlFile := et.env.BasePath + "/vendors/templates/mail/auth/forgot.tmpl"

	data := map[string]string{
		"name": user.FirstName,
		"link": et.env.SiteUrl + "/forgot-password?token=" + user.VerifyEmailToken,
	}
	go et.email.SendEmail(ch, user.Email, "Recover password", htmlFile, data)
	err = <-ch
	if err != nil {
		et.logger.Zap.Error(err)
		return err
	}
	err = et.userRepository.UpdateColumn(&user, "last_forgot_email_date", time.Now())
	if err != nil {
		et.logger.Zap.Error(err)
		return err
	}
	return nil
}
