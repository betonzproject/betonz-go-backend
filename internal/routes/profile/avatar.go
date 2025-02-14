package profile

import (
	"log"
	"net/http"

	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/utils/formutils"

	"github.com/jackc/pgx/v5/pgtype"
)

type AvatarForm struct {
	ProfileImage string `form:"profileImage" validate:"omitempty,oneof=ant bike cake car deer desert fish forest gentleman hiking lunchbox nature night pet roof surf"`
}

func PostAvatar(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		var avatarForm AvatarForm
		if formutils.ParseDecodeValidateMultipart(app, w, r, &avatarForm) != nil {
			return
		}

		err = app.DB.UpdateUserProfileImage(r.Context(), db.UpdateUserProfileImageParams{
			ID:           user.ID,
			ProfileImage: pgtype.Text{String: avatarForm.ProfileImage, Valid: avatarForm.ProfileImage != ""},
		})
		if err != nil {
			log.Panicln("Can't update user: " + err.Error())
		}

		w.WriteHeader(http.StatusOK)
	}
}
