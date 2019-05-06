package gateway

import "log"
import "github.com/labstack/echo"
import "github.com/labstack/echo/middleware"

func StartEchoHandle() {
	// Echo instance

	e = echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	err := LoadRoutes()
	if err != nil {
		log.Println("Load routes error", err)
	}

	// Start server
	e.Logger.Fatal(e.Start(":" + conf.Base.Port))
}
