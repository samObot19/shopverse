# Notification Microservice

This repository contains a notification microservice built with Go. The service is responsible for sending email notifications to users when an order is placed. It listens for order messages from the order microservice and triggers email notifications accordingly.

## Project Structure

```
notification-service
├── cmd
│   └── main.go                  # Entry point of the application
├── config
│   └── config.go                # Configuration settings
├── internal
│   ├── email
│   │   ├── email.go             # Email sending service
│   │   └── email_test.go        # Unit tests for email service
│   ├── handlers
│   │   └── notification_handler.go # HTTP handler for notifications
│   ├── models
│   │   └── order_message.go      # Order message structure
│   ├── repository
│   │   └── notification_repository.go # Database interaction for notifications
│   ├── services
│   │   └── notification_service.go # Business logic for notifications
│   └── utils
│       └── validation.go         # Utility functions for validation
├── pkg
│   └── logger
│       └── logger.go             # Logging utility
├── go.mod                        # Module definition and dependencies
├── go.sum                        # Dependency checksums
└── README.md                     # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.16 or later
- A working email service (e.g., SMTP server) for sending emails

### Installation

1. Clone the repository:

   ```
   git clone https://github.com/yourusername/notification-service.git
   cd notification-service
   ```

2. Install the dependencies:

   ```
   go mod tidy
   ```

### Configuration

Update the `config/config.go` file with your email service credentials and any other necessary configuration settings.

### Running the Service

To run the notification microservice, execute the following command:

```
go run cmd/main.go
```

### Testing

To run the tests for the email service, use the following command:

```
go test ./internal/email
```

## Usage

Once the service is running, it will listen for incoming order messages. Upon receiving an order message, it will send an email notification to the user associated with the order.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.