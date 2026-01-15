package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/config"
	"github.com/andre/project-app-bioskop-golang/internal/handlers"
	"github.com/andre/project-app-bioskop-golang/internal/middleware"
	"github.com/andre/project-app-bioskop-golang/internal/repositories"
	"github.com/andre/project-app-bioskop-golang/internal/services"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg := config.LoadConfig()
	logger.Info("Configuration loaded",
		zap.String("server_port", cfg.Server.Port),
		zap.String("server_env", cfg.Server.Env),
	)

	// Connect to database
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer sqlDB.Close()

	// Ping database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = sqlDB.PingContext(ctx)
	cancel()
	if err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}
	logger.Info("Database connected successfully")

	// Create PGX connection pool for queries
	pgxCfg, err := pgx.ParseConfig(dbURL)
	if err != nil {
		logger.Fatal("Failed to parse database config", zap.Error(err))
	}

	conn, err := pgx.ConnectConfig(context.Background(), pgxCfg)
	if err != nil {
		logger.Fatal("Failed to connect to database with pgx", zap.Error(err))
	}
	defer conn.Close(context.Background())

	// Initialize validator
	validate := validator.New()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(conn)
	cinemaRepo := repositories.NewCinemaRepository(conn)
	seatRepo := repositories.NewSeatRepository(conn)
	bookingRepo := repositories.NewBookingRepository(conn)
	paymentRepo := repositories.NewPaymentRepository(conn)
	emailRepo := repositories.NewEmailVerificationRepository(conn)

	// Initialize services
	emailService := services.NewEmailService(emailRepo, logger, cfg.Email.APIURL, cfg.Email.APIKey)
	userService := services.NewUserService(userRepo, emailService, cfg.JWT.Secret)
	cinemaService := services.NewCinemaService(cinemaRepo)
	seatService := services.NewSeatService(seatRepo)
	bookingService := services.NewBookingService(bookingRepo, seatRepo, cinemaRepo)
	paymentService := services.NewPaymentService(paymentRepo, bookingRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, validate, logger)
	cinemaHandler := handlers.NewCinemaHandler(cinemaService, validate, logger)
	seatHandler := handlers.NewSeatHandler(seatService, validate, logger)
	bookingHandler := handlers.NewBookingHandler(bookingService, validate, logger)
	paymentHandler := handlers.NewPaymentHandler(paymentService, validate, logger)
	emailHandler := handlers.NewEmailHandler(emailService, validate, logger)

	// Setup router
	router := chi.NewRouter()

	// Global middleware
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Recoverer)
	router.Use(middleware.LoggingMiddleware(logger))

	// Public routes
	router.Post("/api/register", userHandler.Register)
	router.Post("/api/login", userHandler.Login)

	// Email verification routes (public)
	router.Post("/api/verify-email", emailHandler.VerifyEmail)
	router.Post("/api/resend-otp", emailHandler.ResendOTP)

	// Cinema routes (public)
	router.Get("/api/cinemas", cinemaHandler.GetAllCinemas)
	router.Get("/api/cinemas/{cinemaId}", cinemaHandler.GetCinemaByID)

	// Seat routes (public)
	router.Get("/api/cinemas/{cinemaId}/seats", seatHandler.GetSeatAvailability)

	// Payment methods (public)
	router.Get("/api/payment-methods", paymentHandler.GetPaymentMethods)

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(userService))

		// User routes
		r.Post("/api/logout", userHandler.Logout)
		r.Get("/api/user/profile", userHandler.GetProfile)

		// Booking routes
		r.Post("/api/booking", bookingHandler.CreateBooking)
		r.Get("/api/user/bookings", bookingHandler.GetUserBookings)

		// Payment routes
		r.Post("/api/pay", paymentHandler.ProcessPayment)
	})

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", zap.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	logger.Info("Shutting down server...")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	} else {
		logger.Info("Server shut down successfully")
	}
}
