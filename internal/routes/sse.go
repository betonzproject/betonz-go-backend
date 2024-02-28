package routes

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
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
		log.Printf("%s (%x) connected to SSE channel", user.Username, randomBytes)

		rc := http.NewResponseController(w)
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		fmt.Fprintf(w, "data: %s\n\n", "connect")
		rc.Flush()
		for {
			select {
			case <-r.Context().Done():
				log.Printf("%s (%x) disconnected from SSE channel", user.Username, randomBytes)
				app.EventServer.Unsubscribe([16]byte(randomBytes))
				return
			case message := <-connection.MessageChannel:
				log.Printf("Sending message %s to %s\n", message, user.Username)
				fmt.Fprintf(w, "data: %s\n\n", message)
				rc.Flush()
			}
		}
	}
}
