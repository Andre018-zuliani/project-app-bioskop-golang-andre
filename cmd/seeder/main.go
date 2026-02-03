package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/config"
	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/andre/project-app-bioskop-golang/internal/repositories"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Load configuration
	cfg := config.LoadConfig()

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
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlDB.Close()

	// Ping database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = sqlDB.PingContext(ctx)
	cancel()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connected successfully")

	// Create PGX connection pool for queries
	pgxCfg, err := pgx.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Failed to parse database config: %v", err)
	}

	conn, err := pgx.ConnectConfig(context.Background(), pgxCfg)
	if err != nil {
		log.Fatalf("Failed to connect to database with pgx: %v", err)
	}
	defer conn.Close(context.Background())

	// Initialize repositories
	cinemaRepo := repositories.NewCinemaRepository(conn)
	seatRepo := repositories.NewSeatRepository(conn)

	ctx = context.Background()

	// Seed cinemas
	log.Println("Seeding cinemas...")
	cinemas := []*models.Cinema{
		{
			Name:       "CGV Cinemas - Jakarta",
			Location:   "Blok M Plaza",
			City:       "Jakarta",
			Address:    "Jl. Melawai No. 1, Blok M, Jakarta Selatan",
			TotalSeats: 150,
			ImageURL:   "https://via.placeholder.com/300x200?text=CGV+Jakarta",
		},
		{
			Name:       "Cinemaxx - Surabaya",
			Location:   "Pakuwon Indah",
			City:       "Surabaya",
			Address:    "Jl. Raya Pakuwon Indah, Surabaya",
			TotalSeats: 200,
			ImageURL:   "https://via.placeholder.com/300x200?text=Cinemaxx+Surabaya",
		},
		{
			Name:       "Premiere Cinema - Bandung",
			Location:   "Bandung Indah Plaza",
			City:       "Bandung",
			Address:    "Jl. Ir. H. Juanda No. 1, Bandung",
			TotalSeats: 120,
			ImageURL:   "https://via.placeholder.com/300x200?text=Premiere+Bandung",
		},
		{
			Name:       "TheScreen Cinemas - Medan",
			Location:   "Medan Fair",
			City:       "Medan",
			Address:    "Jl. Jend. Gatot Subroto No. 1, Medan",
			TotalSeats: 180,
			ImageURL:   "https://via.placeholder.com/300x200?text=TheScreen+Medan",
		},
		{
			Name:       "Studio 21 - Bali",
			Location:   "Denpasar",
			City:       "Bali",
			Address:    "Jl. Raya Puputan No. 1, Denpasar",
			TotalSeats: 160,
			ImageURL:   "https://via.placeholder.com/300x200?text=Studio21+Bali",
		},
	}

	var createdCinemas []*models.Cinema
	for _, cinema := range cinemas {
		err := cinemaRepo.CreateCinema(ctx, cinema)
		if err != nil {
			log.Printf("Warning: Could not create cinema %s: %v\n", cinema.Name, err)
			continue
		}
		log.Printf("Cinema created: %s (ID: %d)\n", cinema.Name, cinema.ID)
		createdCinemas = append(createdCinemas, cinema)
	}

	// Seed seats
	log.Println("Seeding seats...")
	for _, cinema := range createdCinemas {
		// Create seats: 5 rows, 30 seats per row
		seatCounter := 0
		for row := 1; row <= 5; row++ {
			for seatNum := 1; seatNum <= 30; seatNum++ {
				seatType := "standard"
				price := 50000.0

				// Premium seats (rows 3-4)
				if row >= 3 && row <= 4 {
					seatType = "premium"
					price = 70000.0
				}

				// VIP seats (row 5)
				if row == 5 {
					seatType = "vip"
					price = 100000.0
				}

				seatLetter := string(rune('A' + seatNum - 1))
				seat := &models.Seat{
					CinemaID:   cinema.ID,
					SeatNumber: fmt.Sprintf("%d%s", row, seatLetter),
					RowNumber:  row,
					SeatType:   seatType,
					Price:      price,
				}

				err := seatRepo.CreateSeat(ctx, seat)
				if err != nil {
					log.Printf("Warning: Could not create seat for cinema %s: %v\n", cinema.Name, err)
					continue
				}
				seatCounter++
			}
		}
		log.Printf("Created %d seats for cinema: %s\n", seatCounter, cinema.Name)

		// Create seat availability for next 10 days
		for i := 0; i < 10; i++ {
			date := time.Now().AddDate(0, 0, i)
			times := []string{"10:00", "13:00", "16:00", "19:00", "21:00"}
			for _, timeStr := range times {
				err := seatRepo.CreateSeatAvailability(ctx, cinema.ID, date, timeStr)
				if err != nil {
					log.Printf("Warning: Could not create seat availability: %v\n", err)
				}
			}
		}
		log.Printf("Created seat availability for cinema: %s\n", cinema.Name)
	}

	log.Println("Seeding completed successfully!")
}
