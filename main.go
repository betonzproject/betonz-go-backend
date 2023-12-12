package main

import (
	"log"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/routes"
	"github.com/doorman2137/betonz-go/internal/routes/admin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := app.NewApp()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", routes.GetIndex(app))
	r.Post("/login", routes.PostLogin(app))
	r.Post("/logout", routes.PostLogout(app))
	r.Route("/admin", func(r chi.Router) {
		r.Get("/transaction-request", admin.GetTransactionRequest(app))
		r.Get("/players", admin.GetPlayers(app))
		r.Post("/players", admin.PostPlayers(app))
	})

	log.Println("🥖 Server started at port 8080!")
	http.ListenAndServe(":8080", app.Scs.LoadAndSave(r))
}
