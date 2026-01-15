package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// NotificationService handles async notifications
type NotificationService struct {
	logger *zap.Logger
}

// NewNotificationService creates a new notification service
func NewNotificationService(logger *zap.Logger) *NotificationService {
	return &NotificationService{
		logger: logger,
	}
}

// SendBookingConfirmationAsync sends booking confirmation notification asynchronously
func (s *NotificationService) SendBookingConfirmationAsync(ctx context.Context, userEmail string, bookingID int, cinemaName string, seatNumbers []string, bookingTime string) {
	// Run notification in goroutine for async execution
	go func() {
		// Use background context to prevent cancellation affecting notification
		notifCtx := context.Background()

		s.logger.Info("Sending booking confirmation notification",
			zap.String("email", userEmail),
			zap.Int("booking_id", bookingID),
			zap.String("cinema", cinemaName),
		)

		// Simulate notification sending (email/SMS/push notification)
		time.Sleep(100 * time.Millisecond) // Simulate network delay

		// In production, this would:
		// - Send email via SMTP or email service (SendGrid, AWS SES)
		// - Send SMS via Twilio or similar
		// - Send push notification
		// - Store notification in database

		message := fmt.Sprintf(
			"Booking Confirmed! Cinema: %s, Seats: %v, Time: %s. Booking ID: #%d",
			cinemaName, seatNumbers, bookingTime, bookingID,
		)

		s.logger.Info("Booking notification sent successfully",
			zap.String("email", userEmail),
			zap.Int("booking_id", bookingID),
			zap.String("message", message),
		)

		// Log to file asynchronously
		s.logNotificationToFile(notifCtx, "booking_confirmation", userEmail, bookingID, message)
	}()
}

// SendPaymentConfirmationAsync sends payment confirmation notification asynchronously
func (s *NotificationService) SendPaymentConfirmationAsync(ctx context.Context, userEmail string, paymentID int, bookingID int, amount float64, paymentMethod string) {
	// Run notification in goroutine for async execution
	go func() {
		notifCtx := context.Background()

		s.logger.Info("Sending payment confirmation notification",
			zap.String("email", userEmail),
			zap.Int("payment_id", paymentID),
			zap.Int("booking_id", bookingID),
		)

		time.Sleep(100 * time.Millisecond) // Simulate network delay

		message := fmt.Sprintf(
			"Payment Successful! Amount: Rp %.2f, Method: %s, Booking ID: #%d, Payment ID: #%d",
			amount, paymentMethod, bookingID, paymentID,
		)

		s.logger.Info("Payment notification sent successfully",
			zap.String("email", userEmail),
			zap.Int("payment_id", paymentID),
			zap.String("message", message),
		)

		s.logNotificationToFile(notifCtx, "payment_confirmation", userEmail, paymentID, message)
	}()
}

// SendBookingReminderAsync sends booking reminder notification asynchronously
func (s *NotificationService) SendBookingReminderAsync(ctx context.Context, userEmail string, bookingID int, cinemaName string, showTime time.Time) {
	go func() {
		notifCtx := context.Background()

		// Calculate time until show
		timeUntilShow := time.Until(showTime)

		// Send reminder 2 hours before show
		if timeUntilShow > 2*time.Hour {
			time.Sleep(timeUntilShow - 2*time.Hour)
		}

		s.logger.Info("Sending booking reminder",
			zap.String("email", userEmail),
			zap.Int("booking_id", bookingID),
			zap.Time("show_time", showTime),
		)

		message := fmt.Sprintf(
			"Reminder: Your movie at %s starts at %s. Booking ID: #%d",
			cinemaName, showTime.Format("15:04"), bookingID,
		)

		s.logger.Info("Booking reminder sent",
			zap.String("email", userEmail),
			zap.Int("booking_id", bookingID),
			zap.String("message", message),
		)

		s.logNotificationToFile(notifCtx, "booking_reminder", userEmail, bookingID, message)
	}()
}

// logNotificationToFile logs notification to file asynchronously
func (s *NotificationService) logNotificationToFile(ctx context.Context, notificationType string, recipient string, entityID int, message string) {
	// This runs in goroutine already, so it's async
	// In production, this would write to a separate log file or database
	s.logger.Info("Notification logged",
		zap.String("type", notificationType),
		zap.String("recipient", recipient),
		zap.Int("entity_id", entityID),
		zap.String("message", message),
		zap.Time("timestamp", time.Now()),
	)
}

// ProcessBulkNotificationsAsync processes multiple notifications in parallel
func (s *NotificationService) ProcessBulkNotificationsAsync(ctx context.Context, notifications []NotificationTask) {
	// Use worker pool pattern with goroutines
	workers := 5
	tasks := make(chan NotificationTask, len(notifications))

	// Start workers
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			for task := range tasks {
				s.logger.Info("Worker processing notification",
					zap.Int("worker_id", workerID),
					zap.String("type", task.Type),
					zap.String("recipient", task.Recipient),
				)

				// Process notification
				time.Sleep(50 * time.Millisecond) // Simulate processing

				s.logger.Info("Worker completed notification",
					zap.Int("worker_id", workerID),
					zap.String("recipient", task.Recipient),
				)
			}
		}(i)
	}

	// Send tasks to workers
	for _, notification := range notifications {
		tasks <- notification
	}
	close(tasks)
}

// NotificationTask represents a notification task
type NotificationTask struct {
	Type      string
	Recipient string
	Message   string
	EntityID  int
}
