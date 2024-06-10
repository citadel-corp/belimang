package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/citadel-corp/belimang/internal/common/db"
	"github.com/citadel-corp/belimang/internal/common/middleware"
	"github.com/citadel-corp/belimang/internal/image"
	merchantitems "github.com/citadel-corp/belimang/internal/merchant_items"
	"github.com/citadel-corp/belimang/internal/merchants"
	"github.com/citadel-corp/belimang/internal/order"
	"github.com/citadel-corp/belimang/internal/user"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Connect to database
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_PARAMS"))
	db, err := db.Connect(connStr)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Cannot connect to database: %v", err))
		os.Exit(1)
	}

	// Create migrations
	// err = db.UpMigration()
	// if err != nil {
	// 	log.Error().Msg(fmt.Sprintf("Up migration failed: %v", err))
	// 	os.Exit(1)
	// }

	// initialize user domain
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	// initialize merchants domain
	merchantRepository := merchants.NewRepository(db)
	merchantService := merchants.NewService(merchantRepository)
	merchantHandler := merchants.NewHandler(merchantService)

	// initialize merchant items domain
	merchantItemRepository := merchantitems.NewRepository(db)
	merchantItemService := merchantitems.NewService(merchantItemRepository, merchantRepository)
	merchantItemHandler := merchantitems.NewHandler(merchantItemService)

	// initialize order domain
	orderRepository := order.NewRepository(db)
	orderService := order.NewService(orderRepository, merchantRepository, merchantItemRepository)
	orderHandler := order.NewHandler(orderService)

	// initialize image domain
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Cannot create AWS session: %v", err))
		os.Exit(1)
	}
	imageService := image.NewService(sess)
	imageHandler := image.NewHandler(imageService)

	r := mux.NewRouter()
	r.Use(middleware.Logging)
	r.Use(middleware.PanicRecoverer)
	// v1 := r.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Service ready v2")
	})

	//
	r.HandleFunc("/merchants/{lat},{long}", middleware.AuthorizeRole(merchantHandler.ListByDistance, string(user.User))).Methods(http.MethodGet)

	// admin routes
	ar := r.PathPrefix("/admin").Subrouter()
	ar.HandleFunc("/register", userHandler.CreateAdmin).Methods(http.MethodPost)
	ar.HandleFunc("/login", userHandler.LoginUser).Methods(http.MethodPost)
	ar.HandleFunc("/merchants", middleware.AuthorizeRole(merchantHandler.Create, string(user.Admin))).Methods(http.MethodPost)
	ar.HandleFunc("/merchants", middleware.AuthorizeRole(merchantHandler.List, string(user.Admin))).Methods(http.MethodGet)
	ar.HandleFunc("/merchants/{merchantId}/items", middleware.AuthorizeRole(merchantItemHandler.Create, string(user.Admin))).Methods(http.MethodPost)
	ar.HandleFunc("/merchants/{merchantId}/items", middleware.AuthorizeRole(merchantItemHandler.List, string(user.Admin))).Methods(http.MethodGet)

	ur := r.PathPrefix("/users").Subrouter()
	ur.HandleFunc("/register", userHandler.CreateNonAdmin).Methods(http.MethodPost)
	ur.HandleFunc("/login", userHandler.LoginUser).Methods(http.MethodPost)

	ur.HandleFunc("/estimate", middleware.AuthorizeRole(orderHandler.CalculateEstimate, string(user.User))).Methods(http.MethodPost)
	ur.HandleFunc("/orders", middleware.AuthorizeRole(orderHandler.CreateOrder, string(user.User))).Methods(http.MethodPost)
	ur.HandleFunc("/orders", middleware.AuthorizeRole(orderHandler.SearchOrders, string(user.User))).Methods(http.MethodGet)

	// image routes
	ir := r.PathPrefix("/image").Subrouter()
	ir.HandleFunc("", middleware.Authorized(imageHandler.UploadToS3)).Methods(http.MethodPost)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Info().Msg(fmt.Sprintf("HTTP server listening on %s", httpServer.Addr))
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Error().Msg(fmt.Sprintf("HTTP server error: %v", err))
		}
		log.Info().Msg("Stopped serving new connections.")
	}()

	// Listen for the termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until termination signal received
	<-stop
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	log.Info().Msg(fmt.Sprintf("Shutting down HTTP server listening on %s", httpServer.Addr))
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Msg(fmt.Sprintf("HTTP server shutdown error: %v", err))
	}
	log.Info().Msg("Shutdown complete.")
}
