package utils

import (
    "regexp"
)

// ValidateEmail checks if the provided email address is valid.
func ValidateEmail(email string) bool {
    // Regular expression for validating an Email
    regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    re := regexp.MustCompile(regex)
    return re.MatchString(email)
}

// ValidateOrderMessage checks if the order message contains valid data.
func ValidateOrderMessage(orderMessage OrderMessage) bool {
    return ValidateEmail(orderMessage.Email) && orderMessage.OrderID != ""
}