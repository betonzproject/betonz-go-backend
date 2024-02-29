package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/timeutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type GetVerifyIdentityResponse struct {
	*db.IdentityVerificationStatus `json:"identityVerificationStatus"`
}

type IdentityVerificationFormStep1 struct {
	FullName   string `form:"fullName" validate:"required,max=64"`
	NricCode   int    `form:"nricCode" validate:"required,min=1,max=14"`
	NricRegion string `form:"nricRegion" validate:"required,alpha"`
	NricNumber string `form:"nricNumber" validate:"required,number"`
	NricOwner  string `form:"nricOwner" validate:"required,oneof=E N P S T Y"`
	Dob        string `form:"dob" validate:"required,datetime=2006/01/02" key:"verifyIdentity.dob"`
}

type IdentityVerificationFormStep2 struct {
	NricFront string `form:"nricFront"`
	NricBack  string `form:"nricBack"`
}

type IdentityVerificationFormStep3 struct {
	HolderFace string `form:"holderFace"`
}

func GetVerifyIdentity(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		request, err := app.DB.GetLatestIdentityVerificationRequestByUserId(r.Context(), user.ID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			log.Panicln("Can't get latest verification request: " + err.Error())
		}

		step, _ := strconv.Atoi(r.URL.Query().Get("step"))
		if step > 1 {
			// Validate step 1
			if errors.Is(err, pgx.ErrNoRows) {
				http.Redirect(w, r, "/verify-identity?step=1", http.StatusFound)
				return
			}
		}
		if step > 2 {
			// Validate step 2
			if request.NricFront == "" || request.NricBack == "" {
				http.Redirect(w, r, "/verify-identity?step=2", http.StatusFound)
				return
			}
		}

		jsonutils.Write(w, GetVerifyIdentityResponse{IdentityVerificationStatus: &request.Status}, http.StatusOK)
	}
}

func PostVerifyIdentity(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		request, err := app.DB.GetLatestIdentityVerificationRequestByUserId(r.Context(), user.ID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			log.Panicln("Can't get latest identity verification request: " + err.Error())
		}
		if request.Status == db.IdentityVerificationStatusPENDING || request.Status == db.IdentityVerificationStatusVERIFIED {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		step, _ := strconv.Atoi(r.URL.Query().Get("step"))
		switch step {
		case 1:
			var identityVerificationForm IdentityVerificationFormStep1
			err2 := formutils.ParseDecodeValidate(app, w, r, &identityVerificationForm)
			if err2 != nil {
				return
			}

			dob, _ := timeutils.ParseDate(identityVerificationForm.Dob)
			if errors.Is(err, pgx.ErrNoRows) {
				err2 := app.DB.CreateIdentityVerificationRequest(r.Context(), db.CreateIdentityVerificationRequestParams{
					UserId:   user.ID,
					NricName: identityVerificationForm.FullName,
					Nric:     fmt.Sprintf("%d/%s(%s)%s", identityVerificationForm.NricCode, identityVerificationForm.NricRegion, identityVerificationForm.NricOwner, identityVerificationForm.NricNumber),
					Dob:      pgtype.Date{Time: dob, Valid: true},
				})
				if err2 != nil {
					log.Panicln("Can't create identity verification request: " + err2.Error())
				}
			} else {
				err2 := app.DB.UpdateIdentityVerificationRequestById(r.Context(), db.UpdateIdentityVerificationRequestByIdParams{
					ID:       request.ID,
					NricName: pgtype.Text{String: identityVerificationForm.FullName, Valid: true},
					Nric:     pgtype.Text{String: fmt.Sprintf("%d/%s(%s)%s", identityVerificationForm.NricCode, identityVerificationForm.NricRegion, identityVerificationForm.NricOwner, identityVerificationForm.NricNumber), Valid: true},
					Dob:      pgtype.Date{Time: dob, Valid: true},
				})
				if err2 != nil {
					log.Panicln("Can't update identity verification request: " + err2.Error())
				}
			}

			http.Redirect(w, r, "/verify-identity?step=2", http.StatusFound)
			return
		case 2:
			var identityVerificationForm IdentityVerificationFormStep2
			err2 := formutils.ParseDecodeValidateMultipart(app, w, r, &identityVerificationForm)
			if err2 != nil {
				return
			}

			err2 = app.DB.UpdateIdentityVerificationRequestById(r.Context(), db.UpdateIdentityVerificationRequestByIdParams{
				ID:        request.ID,
				NricFront: pgtype.Text{String: identityVerificationForm.NricFront, Valid: true},
				NricBack:  pgtype.Text{String: identityVerificationForm.NricBack, Valid: true},
			})
			if err2 != nil {
				log.Panicln("Can't update identity verification request: " + err2.Error())
			}

			http.Redirect(w, r, "/verify-identity?step=3", http.StatusFound)
			return
		case 3:
			var identityVerificationForm IdentityVerificationFormStep3
			err2 := formutils.ParseDecodeValidateMultipart(app, w, r, &identityVerificationForm)
			if err2 != nil {
				return
			}

			err2 = app.DB.UpdateIdentityVerificationRequestById(r.Context(), db.UpdateIdentityVerificationRequestByIdParams{
				ID:         request.ID,
				HolderFace: pgtype.Text{String: identityVerificationForm.HolderFace, Valid: true},
				Status:     db.NullIdentityVerificationStatus{IdentityVerificationStatus: db.IdentityVerificationStatusPENDING, Valid: true},
			})
			if err2 != nil {
				log.Panicln("Can't update identity verification request: " + err2.Error())
			}

			app.EventServer.NotifyAdmins("request")

			w.WriteHeader(http.StatusCreated)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}
