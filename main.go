package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/jobs"
	"github.com/BetOnz-Company/betonz-go/internal/routes"
	"github.com/BetOnz-Company/betonz-go/internal/routes/admin"
	"github.com/BetOnz-Company/betonz-go/internal/routes/admin/players"
	"github.com/BetOnz-Company/betonz-go/internal/routes/admin/report"
	"github.com/BetOnz-Company/betonz-go/internal/routes/producttype"
	"github.com/BetOnz-Company/betonz-go/internal/routes/profile"
	"github.com/BetOnz-Company/betonz-go/internal/routes/profile/bankingdetails"
	"github.com/BetOnz-Company/betonz-go/internal/routes/profile/transfer"
	"github.com/BetOnz-Company/betonz-go/internal/routes/resetpassword"
	"github.com/BetOnz-Company/betonz-go/internal/routes/verifyemail"

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
	r.Get("/login", routes.GetLogin(app))
	r.Post("/login", routes.PostLogin(app))
	r.Post("/logout", routes.PostLogout(app))
	r.Post("/register", routes.PostRegister(app))
	r.Route("/reset-password", func(r chi.Router) {
		r.Post("/", resetpassword.PostPasswordReset(app))
		r.Get("/{token}", resetpassword.GetPasswordResetToken(app))
		r.Post("/{token}", resetpassword.PostPasswordResetToken(app))
	})
	r.Get("/verify-email/{token}", verifyemail.GetVerifyEmailToken(app))
	r.Get("/leaderboard/{productType}", routes.GetLeaderboard(app))
	r.Route("/{productType}", func(r chi.Router) {
		r.Get("/", producttype.GetProducts(app))
		r.Get("/{product}", producttype.GetProduct(app))
		r.Post("/{product}", producttype.PostProduct(app))
	})
	r.Route("/profile", func(r chi.Router) {
		r.Get("/", profile.GetProfile(app))
		r.Post("/", profile.PostProfile(app))
		r.Post("/avatar", profile.PostAvatar(app))
		r.Get("/deposit", profile.GetDeposit(app))
		r.Post("/deposit", profile.PostDeposit(app))
		r.Get("/withdraw", profile.GetWithdraw(app))
		r.Post("/withdraw", profile.PostWithdraw(app))
		r.Route("/transfer", func(r chi.Router) {
			r.Get("/", transfer.GetTransfer(app))
			r.Post("/", transfer.PostTransfer(app))
			r.Get("/products", transfer.GetProducts(app))
		})
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
		r.Post("/account-settings", profile.PostAccountSettings(app))
	})
	r.Get("/verify-identity", routes.GetVerifyIdentity(app))
	r.Post("/verify-identity", routes.PostVerifyIdentity(app))
	r.Get("/sse", routes.GetSse(app))

	r.Route("/admin", func(r chi.Router) {
		r.Get("/", admin.GetIndex(app))
		r.Get("/transaction-request", admin.GetTransactionRequest(app))
		r.Post("/transaction-request", admin.PostTransactionRequest(app))
		r.Get("/report", report.GetReport(app))
		r.Get("/report/overview", report.GetOverview(app))
		r.Get("/players", players.GetPlayers(app))
		r.Get("/players/{id}", players.GetPlayersById(app))
		r.Post("/players", players.PostPlayers(app))
		r.Get("/banks", admin.GetBanks(app))
		r.Post("/banks", admin.PostBanks(app))
		r.Patch("/banks", admin.PatchBanks(app))
		r.Get("/bet-history", admin.GetBetHistory(app))
		r.Get("/activity-log", admin.GetActivityLog(app))
		r.Get("/identity-verification-request", admin.GetIdentityVerificationRequest(app))
		r.Post("/identity-verification-request", admin.PostIdentityVerificationRequest(app))
		r.Get("/file/{filename}", admin.GetFile(app))
		r.Get("/deposit", admin.GetDeposit(app))
		r.Post("/deposit", admin.PostDeposit(app))
		r.Get("/withdraw", admin.GetWithdraw(app))
		r.Post("/withdraw", admin.PostWithdraw(app))
		r.Get("/maintenance", admin.GetMaintenance(app))
		r.Post("/maintenance", admin.PostMaintenance(app))
		r.Get("/admins", admin.GetAdmin(app))
		r.Post("/admins", admin.PostAdmin(app))
	})

	minutes, _ := strconv.Atoi(os.Getenv("FETCH_BETS_INTERVAL_MINUTES"))
	if minutes > 0 {
		log.Printf("Starting bet fetch ticker every %d minutes", minutes)
		ticker := time.NewTicker(time.Duration(minutes) * time.Minute)
		defer ticker.Stop()

		go func() {
			jobs.FetchBets(app, 0)
			for {
				select {
				case <-ticker.C:
					jobs.FetchBets(app, 0)
				}
			}
		}()
	}

	log.Println("🥖 Server started at port 8080!")
	http.ListenAndServe(":8080", app.Scs.LoadAndSave(r))
}
