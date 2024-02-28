package admin

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/product"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/sliceutils"
	"github.com/doorman2137/betonz-go/internal/utils/timeutils"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateMaintenanceItemForm struct {
	Product product.Product `form:"product" validate:"required"`
	GMT     string          `form:"gmt" validate:"required"`
	From    string          `form:"from" validate:"required,datetime=2006/01/02 15:04:05|datetime=2006/01/02 3:04:05 PM"`
	To      string          `form:"to" validate:"required,datetime=2006/01/02 15:04:05|datetime=2006/01/02 3:04:05 PM"`
}

type DeleteMaintenanceItemForm struct {
	Id int32 `form:"id" validate:"required"`
}

type UpdateMaintenanceItemForm struct {
	Id   int32  `form:"id" validate:"required"`
	From string `form:"from" validate:"required,datetime=2006/01/02 15:04:05|datetime=2006/01/02 3:04:05 PM"`
	To   string `form:"to" validate:"required,datetime=2006/01/02 15:04:05|datetime=2006/01/02 3:04:05 PM"`
	GMT  string `form:"gmt" validate:"required"`
}

type MaintenanceItem struct {
	ID            int32              `json:"id"`
	Product       string             `json:"product"`
	StartTime     pgtype.Timestamptz `json:"startTime"`
	EndTime       pgtype.Timestamptz `json:"endTime"`
	GmtOffsetSecs int32              `json:"gmtOffsetSecs"`
}

type MaintenanceGetResponse struct {
	Products                 map[product.Product]string `json:"products"`
	ProductsUnderMaintenance []MaintenanceItem          `json:"productsUnderMaintenance"`
}

type MaintenancePostResponse struct {
	Action string `json:"action"`
}

func GetMaintenance(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageProductMaintenance) != nil {
			return
		}

		productsUnderMaintenance, err := app.DB.GetMaintenanceList(r.Context())
		if err != nil {
			log.Panicln("Error fetching maintained products: ", err.Error())
		}

		productNames := make(map[product.Product]string)

		for _, p := range product.AllProducts {
			productNames[p] = p.String()
		}

		jsonutils.Write(w, MaintenanceGetResponse{
			Products: productNames,
			ProductsUnderMaintenance: sliceutils.Map(productsUnderMaintenance, func(r db.Maintenance) MaintenanceItem {
				productsUnderMaintenance := product.Product(int(r.ProductCode)).String()

				return MaintenanceItem{
					ID:            r.ID,
					Product:       productsUnderMaintenance,
					StartTime:     pgtype.Timestamptz(r.MaintenancePeriod.Lower),
					EndTime:       pgtype.Timestamptz(r.MaintenancePeriod.Upper),
					GmtOffsetSecs: r.GmtOffsetSecs,
				}
			}),
		}, http.StatusOK)
	}
}

func PostMaintenance(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageProductMaintenance) != nil {
			return
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		if r.URL.Query().Has("/create") {
			var createMaintenance CreateMaintenanceItemForm
			if formutils.ParseDecodeValidate(app, w, r, &createMaintenance) != nil {
				return
			}

			gmtInt, _ := strconv.Atoi(createMaintenance.GMT)
			loc := time.FixedZone("GMT", gmtInt)

			from, err := timeutils.ParseDateTimeInLocation(createMaintenance.From, *loc)
			if err != nil {
				http.Error(w, "Invalid date", http.StatusBadRequest)
				return
			}
			to, err := timeutils.ParseDateTimeInLocation(createMaintenance.To, *loc)
			if err != nil {
				http.Error(w, "Invalid date", http.StatusBadRequest)
				return
			}

			if to.Before(from) {
				http.Error(w, "maintenance.dateRange.error", http.StatusBadRequest)
				return
			}

			err = qtx.CreateMaintenanceItem(r.Context(), db.CreateMaintenanceItemParams{
				ProductCode:       int32(createMaintenance.Product),
				MaintenancePeriod: pgtype.Range[pgtype.Timestamptz]{Lower: pgtype.Timestamptz{Time: from, Valid: true}, Upper: pgtype.Timestamptz{Time: to, Valid: true}, LowerType: pgtype.Inclusive, UpperType: pgtype.Inclusive, Valid: true},
				GmtOffsetSecs:     int32(gmtInt),
			})
			if err != nil {
				log.Panicln("Cannot create maintenceProduct", err.Error())
			}

			err = utils.LogEvent(qtx, r, user.ID, db.EventTypeMAINTENANCEADD, db.EventResultSUCCESS, "", map[string]any{
				"productCode": createMaintenance.Product,
				"startTime":   createMaintenance.From,
				"endTime":     createMaintenance.To,
				"timezone":    createMaintenance.GMT})
			if err != nil {
				log.Panicln("Error creating event: ", err.Error())
			}

			jsonutils.Write(w, MaintenancePostResponse{
				Action: "create",
			}, http.StatusCreated)
		} else if r.URL.Query().Has("/delete") {
			var deleteMaintenanceItemForm DeleteMaintenanceItemForm
			if formutils.ParseDecodeValidate(app, w, r, &deleteMaintenanceItemForm) != nil {
				return
			}

			err = qtx.DeleteMaintenanceItem(r.Context(), deleteMaintenanceItemForm.Id)
			if err != nil {
				log.Panicln("Error deleting maintained product: ", err.Error())
			}

			err = utils.LogEvent(qtx, r, user.ID, db.EventTypeMAINTENANCEDELETE, db.EventResultSUCCESS, "", map[string]any{
				"id": deleteMaintenanceItemForm.Id})
			if err != nil {
				log.Panicln("Error creating event: ", err.Error())
			}

			jsonutils.Write(w, MaintenancePostResponse{
				Action: "delete",
			}, http.StatusOK)
		} else if r.URL.Query().Has("/update") {
			var updateMaintenanceItemForm UpdateMaintenanceItemForm
			if formutils.ParseDecodeValidate(app, w, r, &updateMaintenanceItemForm) != nil {
				return
			}

			gmtInt, _ := strconv.Atoi(updateMaintenanceItemForm.GMT)
			loc := time.FixedZone("GMT", gmtInt)

			from, err := timeutils.ParseDateTimeInLocation(updateMaintenanceItemForm.From, *loc)
			if err != nil {
				http.Error(w, "Invalid date", http.StatusBadRequest)
				return
			}
			to, err := timeutils.ParseDateTimeInLocation(updateMaintenanceItemForm.To, *loc)
			if err != nil {
				http.Error(w, "Invalid date", http.StatusBadRequest)
				return
			}
			if to.Before(from) {
				http.Error(w, "maintenance.dateRange.error", http.StatusBadRequest)
				return
			}

			err = qtx.UpdateMaintenanceItem(r.Context(), db.UpdateMaintenanceItemParams{
				ID:                updateMaintenanceItemForm.Id,
				MaintenancePeriod: pgtype.Range[pgtype.Timestamptz]{Lower: pgtype.Timestamptz{Time: from, Valid: true}, Upper: pgtype.Timestamptz{Time: to, Valid: true}, LowerType: pgtype.Inclusive, UpperType: pgtype.Inclusive, Valid: true},
				GmtOffsetSecs:     int32(gmtInt),
			})
			if err != nil {
				log.Panicln("Error updating maintained product: ", err.Error())
			}

			err = utils.LogEvent(qtx, r, user.ID, db.EventTypeMAINTENANCEUPDATE, db.EventResultSUCCESS, "", map[string]any{
				"id":        updateMaintenanceItemForm.Id,
				"startTime": updateMaintenanceItemForm.From,
				"endTime":   updateMaintenanceItemForm.To,
				"timeZone":  updateMaintenanceItemForm.GMT,
			})
			if err != nil {
				log.Panicln("Error creating event: ", err.Error())
			}

			jsonutils.Write(w, MaintenancePostResponse{
				Action: "update",
			}, http.StatusOK)
		}
		tx.Commit(r.Context())
	}
}
