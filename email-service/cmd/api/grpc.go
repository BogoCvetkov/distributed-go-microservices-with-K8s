package main

import (
	"context"
	"email-service/cmd/api/email_proto"
	mailer "email-service/cmd/api/service"
	"fmt"
)

type EmailServer struct {
	// UnimplementedEmailServiceServer
	email_proto.UnimplementedEmailServiceServer
	Mailer *mailer.Mail
}

func (s *EmailServer) SendEmail(ctx context.Context, data *email_proto.EmailRequest) (*email_proto.EmailResponse, error) {

	msg := mailer.Message{
		To:      data.Data.To,
		Subject: data.Data.Subject,
		Data:    data.Data.Message,
	}

	err := s.Mailer.SendSMTPMessage(msg)

	if err != nil {
		return nil, err
	}

	response := email_proto.EmailResponse{
		Message: fmt.Sprintf("GRPC: Message send to --> %s", data.Data.To),
	}

	return &response, nil
}
