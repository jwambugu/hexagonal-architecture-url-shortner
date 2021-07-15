package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jwambugu/hexagonal-architecture-url-shortener/api"
	"github.com/jwambugu/hexagonal-architecture-url-shortener/config"
	mongoRepo "github.com/jwambugu/hexagonal-architecture-url-shortener/repository/mongodb"
	redisRepo "github.com/jwambugu/hexagonal-architecture-url-shortener/repository/redis"
	"github.com/jwambugu/hexagonal-architecture-url-shortener/shortener"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func chooseRepo() shortener.RedirectRepository {
	switch os.Getenv("MS_REPOSITORY") {
	case "redis":
		url := os.Getenv("REDIS_URL")

		client, err := config.NewRedisConfig(url).RedisClient()

		if err != nil {
			log.Fatal(err)
		}

		repo := redisRepo.NewRedisRepository(client)
		redirectRepository, err := repo.RedirectRepository()

		if err != nil {
			log.Fatal(err)
		}

		return redirectRepository
	case "mongo":
		url := os.Getenv("MONGO_URL")
		database := os.Getenv("MONGO_DB")

		repository, err := mongoRepo.NewMongoRepository(url, database, 30)

		if err != nil {
			log.Fatal(err)
		}
		return repository
	default:
		return nil
	}
}

func main() {
	repo := chooseRepo()

	if repo == nil {
		log.Fatal("Invalid repository provided.")
	}
	service := shortener.NewRedirectService(repo)
	handler := api.NewHandler(service)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/{code}", handler.Get)
	router.Post("/", handler.Post)

	errorsChan := make(chan error, 2)

	go func() {
		fmt.Println("Listening on port :8000")
		errorsChan <- http.ListenAndServe(":8000", router)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)

		errorsChan <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated: %s", <-errorsChan)
}
