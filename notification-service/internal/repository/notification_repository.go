package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Notification struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification Notification) error
	GetNotificationByID(ctx context.Context, id int) (Notification, error)
}

type notificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) CreateNotification(ctx context.Context, notification Notification) error {
	query := "INSERT INTO notifications (email, subject, body, created_at) VALUES (?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, notification.Email, notification.Subject, notification.Body, time.Now())
	return err
}

func (r *notificationRepository) GetNotificationByID(ctx context.Context, id int) (Notification, error) {
	var notification Notification
	query := "SELECT id, email, subject, body, created_at FROM notifications WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&notification.ID, &notification.Email, &notification.Subject, &notification.Body, &notification.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Notification{}, nil
		}
		return Notification{}, err
	}
	return notification, nil
}