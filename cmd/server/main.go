package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/handler"
	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/middleware"
	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/repository/mysql"
	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/usecase"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "root:@tcp(127.0.0.1:3306)/bank_db?parseTime=true"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-jwt-key"
	}

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("Warning: Database ping failed (%v). Server will start, but database requests may fail until MySQL is up.", err)
	}

	customerRepo := mysql.NewCustomerRepository(db)
	transactionRepo := mysql.NewTransactionRepository(db)

	customerUC := usecase.NewCustomerUseCase(customerRepo, jwtSecret)
	transactionUC := usecase.NewTransactionUseCase(transactionRepo)

	authHandler := handler.NewAuthHandler(customerUC)
	txHandler := handler.NewTransactionHandler(transactionUC)

	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","message":"REST API BANK backend is running"}`))
	})

	r.Post("/customers/register", authHandler.Register)
	r.Post("/customers/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTMiddleware(jwtSecret))

		r.Post("/transactions", txHandler.Create)
		r.Get("/transactions", txHandler.GetList)
		r.Get("/transactions/{transactionId}", txHandler.GetByID)
		r.Patch("/transactions/{transactionId}", txHandler.Update)
		r.Delete("/transactions/{transactionId}", txHandler.Delete)
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
