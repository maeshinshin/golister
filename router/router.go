package router

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/maeshinshin/go-multiapi/cmd/web"
	"github.com/maeshinshin/go-multiapi/handlers"
	"github.com/maeshinshin/go-multiapi/internal/database"
)

type Router struct {
	handlers *handlers.Handlers
}

func NewRouter(db database.Service) *Router {
	return &Router{handlers.NewHandlers(db)}
}

func (r *Router) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/web", echo.WrapHandler(templ.Handler(web.HelloForm())))
	e.POST("/hello", echo.WrapHandler(http.HandlerFunc(web.HelloWebHandler)))

	e.GET("/", r.handlers.HelloWorldHandler)

	e.GET("/health", r.handlers.HealthHandler)

	return e
}
