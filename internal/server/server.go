package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/maeshinshin/go-multiapi/internal/database"
	"github.com/maeshinshin/go-multiapi/router"
)

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	router := router.NewRouter(database.New())

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
