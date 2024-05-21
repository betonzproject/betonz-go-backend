package routes

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
)

func GetSse(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetExtendedUser(app, w, r)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		randomBytes := make([]byte, 16)
		rand.Read(randomBytes)

		connection := app.EventServer.Subscribe([16]byte(randomBytes), user)
		log.Printf("%x (%s) connected to SSE channel", randomBytes, user.Username)

		rc := http.NewResponseController(w)
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		fmt.Fprintf(w, "data: %s\n\n", "connect")
		rc.Flush()
		for {
			select {
			case <-r.Context().Done():
				log.Printf("%x (%s) disconnected from SSE channel", randomBytes, user.Username)
				app.EventServer.Unsubscribe([16]byte(randomBytes))
				return
			case <-time.After(time.Duration(30 * time.Second)):
				fmt.Fprint(w, "data: keepalive\n\n")
				rc.Flush()
			case message := <-connection.MessageChannel:
				log.Printf("Sending message %s to %x (%s)\n", message, randomBytes, user.Username)
				fmt.Fprintf(w, "data: %s\n\n", message)
				rc.Flush()
			}
		}
	}
}
