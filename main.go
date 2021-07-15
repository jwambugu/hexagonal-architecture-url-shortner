package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jwambugu/hexagonal-architecture-url-shortener/api"
	mongoRepo "github.com/jwambugu/hexagonal-architecture-url-shortener/repository/mongodb"
	redisRepo "github.com/jwambugu/hexagonal-architecture-url-shortener/repository/redis"
	"github.com/jwambugu/hexagonal-architecture-url-shortener/shortener"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func chooseRepo(repo string) shortener.RedirectRepository {
	switch repo {
	case "redis":
		repository, err := redisRepo.NewRedisRepository("redis://localhost:6379")

		if err != nil {
			log.Fatal(err)
		}
		return repository
	case "mongo":
		url := ""
		repository, err := mongoRepo.NewMongoRepository(url, "redirects", 30)

		if err != nil {
			log.Fatal(err)
		}
		return repository
	default:
		return nil
	}
}

func main() {
	repo := chooseRepo("redis")
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
