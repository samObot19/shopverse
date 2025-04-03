package services

import (
	"errors"
	"github.com/samObot19/shopverse/notification-service/internal/email"
	"github.com/samObot19/shopverse/notification-service/internal/models"
	"github.com/samObot19/shopverse/notification-service/internal/repository"
)

type NotificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (ns *NotificationService) ProcessOrderMessage(orderMessage models.OrderMessage) error {
	if err := validateOrderMessage(orderMessage); err != nil {
		return err
	}

	recipient := orderMessage.UserEmail
	subject := "Order Confirmation"
	body := generateEmailBody(orderMessage)

	if err := email.SendEmail(recipient, subject, body); err != nil {
		return errors.New("failed to send email: " + err.Error())
	}

	return ns.repo.SaveNotification(orderMessage)
}

func validateOrderMessage(orderMessage models.OrderMessage) error {
	if orderMessage.UserEmail == "" {
		return errors.New("user email is required")
	}
	// Additional validation logic can be added here
	return nil
}

func generateEmailBody(orderMessage models.OrderMessage) string {
	return "Dear User,\n\nThank you for your order!\n\nOrder ID: " + orderMessage.OrderID + "\n\nBest Regards,\nYour Company"
}