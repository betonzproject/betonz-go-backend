package acl

import (
	"errors"
	"net/http"
	"slices"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/db"
)

type Permission string

const (
	ManageTransactionRequests Permission = "ManageTransactionRequests"
	ViewTransactionLogs       Permission = "ViewTransactionLogs"
	ManagePlayers             Permission = "ManagePlayers"
	ManageAdmins              Permission = "ManageAdmins"
	ToggleSystemBanks         Permission = "ToggleSystemBanks"
	ManageSystemBanks         Permission = "ManageSystemBanks"
	ViewActivityLog           Permission = "ViewActivityLog"
	ViewSuperadminActivityLog Permission = "ViewSuperadminActivityLog"
)

var Acl = map[db.Role][]Permission{
	db.RoleADMIN: {
		ManageTransactionRequests,
		ViewTransactionLogs,
		ManagePlayers,
		ToggleSystemBanks,
		ViewActivityLog,
	},
	db.RoleSUPERADMIN: {
		ManageTransactionRequests,
		ViewTransactionLogs,
		ManagePlayers,
		ManageAdmins,
		ToggleSystemBanks,
		ManageSystemBanks,
		ViewActivityLog,
		ViewSuperadminActivityLog,
	},
}

// Returns a bool indicating whether a role is authorized with a given permission
func IsAuthorized(role db.Role, permission Permission) bool {
	return slices.Contains(Acl[role], permission)
}

// Checks that the user is authorized with the given permission. If not, return an error and show 404 to the
// user
func Authorize(app *app.App, w http.ResponseWriter, r *http.Request, role db.Role, permission Permission) error {
	if !IsAuthorized(role, permission) {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return errors.New("Unauthorized")
	}
	return nil
}
