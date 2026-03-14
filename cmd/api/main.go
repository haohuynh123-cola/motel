package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"tro-go/internal/adapter/db/postgres"
	"tro-go/internal/adapter/handler"
	kafkaAdapter "tro-go/internal/adapter/kafka"
	"tro-go/internal/adapter/repository"
	"tro-go/internal/usecase"
	"tro-go/pkg/config"
	"tro-go/pkg/email"
	"tro-go/pkg/kafka"
)

func runMigrations(dbURL string) {
	log.Println("Running database migrations...")

	var m *migrate.Migrate
	var err error

	for i := 0; i < 5; i++ {
		m, err = migrate.New("file://db/migrations", dbURL)
		if err == nil {
			break
		}
		log.Printf("Migration: Waiting for database... (attempt %d/5)\n", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not create migration instance after retries: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not run up migrations: %v", err)
	}
	log.Println("Database migrations completed successfully.")
}

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Run Migrations Automatically
	runMigrations(cfg.DatabaseURL)

	// 3. Setup Context for Application Lifecycle
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 4. Setup Database Connection Pool
	dbPool, err := postgres.ConnectPool(ctx, cfg.DatabaseURL, cfg.MaxConns)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer dbPool.Close()

	// 5. Instantiate Repositories
	houseRepo := repository.NewHouseRepository(dbPool)
	roomRepo := repository.NewRoomRepository(dbPool)
	userRepo := repository.NewUserRepository(dbPool)
	chatRepo := repository.NewChatRepository(dbPool)
	appRepo := repository.NewAppointmentRepository(dbPool)

	// 5.5 Instantiate Email Sender (Dùng cho Worker)
	emailSender := email.NewSMTPEmailSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword)

	// 5.6 Initialize Kafka Components
	kafkaProducer := kafka.NewProducer(cfg.KafkaBrokers, "notifications")
	notificationProvider := kafkaAdapter.NewNotificationAdapter(kafkaProducer)

	kafkaConsumer := kafka.NewConsumer(cfg.KafkaBrokers, "notifications", "email-group")
	notificationWorker := usecase.NewNotificationWorker(kafkaConsumer, emailSender)

	// Start Worker in background
	go notificationWorker.Start(ctx)

	// 6. Instantiate UseCases
	houseUseCase := usecase.NewHouseUseCase(houseRepo)
	roomUseCase := usecase.NewRoomUseCase(roomRepo, appRepo, notificationProvider)
	userUseCase := usecase.NewUserUseCase(userRepo, cfg.JwtSecret)
	chatUseCase := usecase.NewChatUseCase(chatRepo)

	// 6.5 Initialize Chat Hub
	chatHub := handler.NewChatHub(chatUseCase)
	go chatHub.Run()

	// 7. Setup Echo HTTP Server
	e := echo.New()

	// Docs Handler (Register early to be outside of auth groups)
	handler.NewDocsHandler(e)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health Check Route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Tro-Go API is running")
	})

	// 8. Register API Handlers
	v1 := e.Group("/api/v1")

	// Public routes
	handler.NewUserHandler(v1, userUseCase)

	// Protected routes
	houseGroup := v1.Group("")
	houseGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(cfg.JwtSecret),
		TokenLookup: "header:Authorization:Bearer ,query:token",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.MapClaims)
		},
	}))

	handler.NewHouseHandler(houseGroup, houseUseCase)
	handler.NewRoomHandler(houseGroup, roomUseCase)

	handler.NewProtectedUserHandler(houseGroup, userUseCase)
	handler.NewChatHandler(houseGroup, chatHub, chatUseCase)

	// 9. Start Server with Graceful Shutdown
	go func() {
		if err := e.Start(":" + cfg.AppPort); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Gracefully shutting down server...")

	// Close Kafka Producer
	kafkaProducer.Close()
	kafkaConsumer.Close()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := e.Shutdown(ctxShutdown); err != nil {
		e.Logger.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
