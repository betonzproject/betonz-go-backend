package main

import (
	"log"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/routes"
	"github.com/doorman2137/betonz-go/internal/routes/admin"
	"github.com/doorman2137/betonz-go/internal/routes/admin/players"
	"github.com/doorman2137/betonz-go/internal/routes/admin/report"
	"github.com/doorman2137/betonz-go/internal/routes/producttype"
	"github.com/doorman2137/betonz-go/internal/routes/profile"
	"github.com/doorman2137/betonz-go/internal/routes/profile/bankingdetails"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := app.NewApp()

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Get("/", routes.GetIndex(app))
	r.Post("/login", routes.PostLogin(app))
	r.Post("/logout", routes.PostLogout(app))
	r.Get("/leaderboard/{productType}", routes.GetLeaderboard(app))
	r.Route("/{productType}", func(r chi.Router) {
		r.Get("/", producttype.GetProducts)
		r.Get("/{product}", producttype.GetProduct(app))
		r.Post("/{product}", producttype.PostProduct(app))
	})
	r.Route("/profile", func(r chi.Router) {
		r.Post("/", profile.PostProfile(app))
		r.Post("/avatar", profile.PostAvatar(app))
		r.Get("/deposit", profile.GetDeposit(app))
		r.Post("/deposit", profile.PostDeposit(app))
		r.Get("/withdraw", profile.GetWithdraw(app))
		r.Post("/withdraw", profile.PostWithdraw(app))
		r.Get("/history", profile.GetHistory(app))
		r.Route("/banking-details", func(r chi.Router) {
			r.Get("/", bankingdetails.GetBanks(app))
			r.Post("/", bankingdetails.DeleteBank(app))
			r.Post("/add-bank", bankingdetails.AddBank(app))
			r.Get("/{bankId}", bankingdetails.GetBankById(app))
			r.Patch("/{bankId}", bankingdetails.PatchBankById(app))
		})
		r.Get("/notifications", profile.GetNotifications(app))
		r.Post("/notifications", profile.PostNotifications(app))
	})

	r.Route("/admin", func(r chi.Router) {
		r.Get("/transaction-request", admin.GetTransactionRequest(app))
		r.Get("/transaction-log", admin.GetTransactionLog(app))
		r.Get("/report", report.GetReport(app))
		r.Get("/report/overview", report.GetOverview(app))
		r.Get("/players", players.GetPlayers(app))
		r.Get("/players/{id}", players.GetPlayersById(app))
		r.Post("/players", players.PostPlayers(app))
		r.Get("/banks", admin.GetBanks(app))
		r.Post("/banks", admin.PostBanks(app))
		r.Patch("/banks", admin.PatchBanks(app))
		r.Get("/activity-log", admin.GetActivityLog(app))
	})

	log.Println("🥖 Server started at port 8080!")
	http.ListenAndServe(":8080", app.Scs.LoadAndSave(r))
}
