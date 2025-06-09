package main

import (
	"net/http"
	"time"

	app "dessert-ordering-go-system/internal/app"
	routes "dessert-ordering-go-system/routes"
)

func main() {
	a := app.NewApplication()
	appRoutes := routes.NewRoutes(a)

	defer a.DB.Close() // Ensure the connection is closed when main exits
	defer a.RedisPool.Close()

	// Configure server with timeouts
	s := &http.Server{
		Addr:         ":8080",
		Handler:      appRoutes,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	a.Loggers.Info.Println("Server Started On PORT :8080")

	a.Loggers.Error.Fatal(s.ListenAndServe())
}
