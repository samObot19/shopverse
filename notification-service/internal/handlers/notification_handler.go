package handlers

import (
    "net/http"
    "encoding/json"
    "github.com/samObot19/shopverse/notification-service/internal/models"
    "github.com/samObot19/shopverse/notification-service/internal/services"
    "github.com/samObot19/shopverse/notification-service/pkg/logger"
)

type NotificationHandler struct {
    NotificationService services.NotificationService
}

func NewNotificationHandler(ns services.NotificationService) *NotificationHandler {
    return &NotificationHandler{
        NotificationService: ns,
    }
}

func (h *NotificationHandler) HandleNotification(w http.ResponseWriter, r *http.Request) {
    var orderMessage models.OrderMessage

    if err := json.NewDecoder(r.Body).Decode(&orderMessage); err != nil {
        logger.Error("Failed to decode order message: %v", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := h.NotificationService.ProcessOrderMessage(orderMessage); err != nil {
        logger.Error("Failed to process order message: %v", err)
        http.Error(w, "Failed to process notification", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(map[string]string{"status": "notification sent"})
}