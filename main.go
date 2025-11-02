package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-journaling/config"
	"backend-journaling/internal/database"
	"backend-journaling/internal/handlers"
	"backend-journaling/internal/middleware"
	"backend-journaling/internal/repository"
	"backend-journaling/internal/service"
	"backend-journaling/pkg/email"
	"backend-journaling/pkg/jwt"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewPostgres(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	mongoDB, err := database.NewMongoDB(cfg.MongoDB.URI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Disconnect(context.Background())

	mongoDatabase := mongoDB.Database(cfg.MongoDB.Database)

	jwtManager, err := jwt.NewManager(cfg.JWT.PrivateKeyPath, cfg.JWT.PublicKeyPath, cfg.JWT.AccessTokenDuration)
	if err != nil {
		log.Fatalf("Failed to initialize JWT manager: %v", err)
	}

	emailSender := email.NewSMTPSender(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
		cfg.SMTP.FromEmail,
		cfg.SMTP.FromName,
	)

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	authEventRepo := repository.NewAuthEventRepository(db)
	profileRepo := repository.NewProfileRepository(db)

	noteRepo := repository.NewNoteRepository(mongoDatabase)
	todoRepo := repository.NewTodoRepository(mongoDatabase)
	taskRepo := repository.NewTaskRepository(mongoDatabase)
	noteGroupRepo := repository.NewNoteGroupRepository(mongoDatabase)

	authService := service.NewAuthService(
		userRepo,
		otpRepo,
		refreshTokenRepo,
		authEventRepo,
		jwtManager,
		emailSender,
		cfg,
	)

	noteService := service.NewNoteService(noteRepo)
	todoService := service.NewTodoService(todoRepo)
	taskService := service.NewTaskService(taskRepo)
	noteGroupService := service.NewNoteGroupService(noteGroupRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo, profileRepo, authService)
	noteHandler := handlers.NewNoteHandler(noteService)
	todoHandler := handlers.NewTodoHandler(todoService)
	taskHandler := handlers.NewTaskHandler(taskService)
	noteGroupHandler := handlers.NewNoteGroupHandler(noteGroupService)

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(middleware.CORS(
		cfg.CORS.AllowedOrigins,
		cfg.CORS.AllowedMethods,
		cfg.CORS.AllowedHeaders,
		cfg.CORS.MaxAge,
	))
	r.Use(middleware.SecurityHeaders())

	// Health check endpoint (outside /api/v1 for monitoring)
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Authentication endpoints
		r.Route("/auth", func(r chi.Router) {
			r.Use(middleware.RateLimit(100, time.Minute))
			r.Post("/register", authHandler.Register)
			r.Post("/verify-otp", authHandler.VerifyOTP)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.Refresh)
			r.Post("/logout", authHandler.Logout)
			r.Post("/forgot-password", authHandler.ForgotPassword)
			r.Post("/reset-password", authHandler.ResetPassword)
			r.Post("/request-otp", authHandler.RequestOTP)
		})

		// Profile endpoints (authenticated)
		r.Route("/profile", func(r chi.Router) {
			r.Use(middleware.Authenticate(jwtManager))
			r.Get("/", userHandler.GetProfile)
			r.Post("/", userHandler.CreateProfile)
			r.Put("/", userHandler.UpdateProfile)
			r.Put("/avatar", userHandler.UpdateAvatar)
			r.Put("/change-password", userHandler.ChangePassword)
		})

		// User management endpoints (admin only)
		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.Authenticate(jwtManager))
			r.Use(middleware.RequireRole("admin"))
			r.Get("/{id}", userHandler.GetUserByID)
		})

		// Notes endpoints (authenticated)
		r.Route("/notes", func(r chi.Router) {
			r.Use(middleware.Authenticate(jwtManager))
			r.Get("/", noteHandler.GetNotes)
			r.Post("/", noteHandler.CreateNote)
			r.Get("/{id}", noteHandler.GetNote)
			r.Patch("/{id}", noteHandler.UpdateNote)
			r.Delete("/{id}", noteHandler.DeleteNote)
			r.Post("/{id}/blocks", noteHandler.AddBlock)
			r.Patch("/{id}/blocks/{blockId}", noteHandler.UpdateBlock)
			r.Delete("/{id}/blocks/{blockId}", noteHandler.DeleteBlock)
			r.Patch("/{id}/blocks/order", noteHandler.ReorderBlocks)
		})

		// Todos endpoints (authenticated)
		r.Route("/todos", func(r chi.Router) {
			r.Use(middleware.Authenticate(jwtManager))
			r.Get("/", todoHandler.GetTodos)
			r.Post("/", todoHandler.CreateTodo)
			r.Patch("/{id}", todoHandler.UpdateTodo)
			r.Delete("/{id}", todoHandler.DeleteTodo)
		})

		// Tasks endpoints (authenticated)
		r.Route("/tasks", func(r chi.Router) {
			r.Use(middleware.Authenticate(jwtManager))
			r.Get("/", taskHandler.GetTasks)
			r.Post("/", taskHandler.CreateTask)
			r.Get("/{id}", taskHandler.GetTask)
			r.Patch("/{id}", taskHandler.UpdateTask)
			r.Delete("/{id}", taskHandler.DeleteTask)
		})

		// Note Groups endpoints (authenticated)
		r.Route("/note-groups", func(r chi.Router) {
			r.Use(middleware.Authenticate(jwtManager))
			r.Get("/", noteGroupHandler.GetGroups)
			r.Post("/", noteGroupHandler.CreateGroup)
			r.Get("/{id}", noteGroupHandler.GetGroup)
			r.Patch("/{id}", noteGroupHandler.UpdateGroup)
			r.Delete("/{id}", noteGroupHandler.DeleteGroup)
			r.Patch("/{id}/pin", noteGroupHandler.PinGroup)
			r.Patch("/{id}/archive", noteGroupHandler.ArchiveGroup)
			r.Post("/{id}/notes", noteGroupHandler.AddNoteToGroup)
			r.Delete("/notes/{noteId}", noteGroupHandler.RemoveNoteFromGroup)
			r.Post("/{id}/move-notes", noteGroupHandler.MoveNotesToGroup)
			r.Get("/{id}/notes", noteGroupHandler.GetNotesInGroup)
		})
	})

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
