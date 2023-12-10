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
)

var acl = map[db.Role][]Permission{
	db.RoleADMIN: {
		ManageTransactionRequests,
	},
	db.RoleSUPERADMIN: {
		ManageTransactionRequests,
	},
}

// Returns a bool indicating whether a role is authorized with a given permission
func IsAuthorized(role db.Role, permission Permission) bool {
	return slices.Contains(acl[role], permission)
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
