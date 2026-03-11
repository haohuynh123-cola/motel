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
	"tro-go/internal/adapter/repository"
	"tro-go/internal/usecase"
	"tro-go/pkg/config"
)

func runMigrations(dbURL string) {
	log.Println("Running database migrations...")

	var m *migrate.Migrate
	var err error

	// Retry logic (Thử lại 5 lần, mỗi lần cách nhau 2 giây)
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
	customerRepo := repository.NewCustomerRepository(dbPool)
	contractRepo := repository.NewContractRepository(dbPool)

	// 6. Instantiate UseCases
	houseUseCase := usecase.NewHouseUseCase(houseRepo)
	roomUseCase := usecase.NewRoomUseCase(roomRepo)
	userUseCase := usecase.NewUserUseCase(userRepo, cfg.JwtSecret)
	customerUseCase := usecase.NewCustomerUseCase(customerRepo)
	contractUseCase := usecase.NewContractUseCase(contractRepo, roomRepo, customerRepo)

	// 7. Setup Echo HTTP Server
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health Check Route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	e.GET("/", func(c echo.Context) error {
		htmlContent := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Tro-Go API</title>
			<style>
				body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; margin: 0; }
				.card { background: white; padding: 2rem; border-radius: 10px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; }
				h1 { color: #2563eb; }
				p { color: #4b5563; }
			</style>
		</head>
		<body>
			<div class="card">
				<h1>🏠 Tro-Go API</h1>
				<p>Hệ thống Backend đã khởi động thành công!</p>
				<p><i>Trạng thái: 🟢 Đang hoạt động</i></p>
			</div>
		</body>
		</html>
		`
		return c.HTML(http.StatusOK, htmlContent)
	})

	// 8. Register API Handlers
	v1 := e.Group("/api/v1")

	// Public routes (Không cần đăng nhập)
	handler.NewUserHandler(v1, userUseCase)

	// Protected routes (Bắt buộc phải có token JWT)
	houseGroup := v1.Group("") // Group này dùng chung tiền tố /api/v1 từ v1
	houseGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.JwtSecret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.MapClaims)
		},
	}))

	handler.NewHouseHandler(houseGroup, houseUseCase)
	handler.NewRoomHandler(houseGroup, roomUseCase)
	handler.NewProtectedUserHandler(houseGroup, userUseCase)
	handler.NewCustomerHandler(houseGroup, customerUseCase)
	handler.NewContractHandler(houseGroup, contractUseCase)

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

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := e.Shutdown(ctxShutdown); err != nil {
		e.Logger.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
