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
	UpdateProfile                                Permission = "UpdateProfile"
	DepositToOwnWallet                           Permission = "DepositToOwnWallet"
	WithdrawFromOwnWallet                        Permission = "WithdrawFromOwnWallet"
	TransferBetweenWallets                       Permission = "TransferBetweenWallets"
	ViewOwnTransactionHistory                    Permission = "ViewOwnTransactionHistory"
	ManageOwnBankingDetails                      Permission = "ManageOwnBankingDetails"
	ViewNotifications                            Permission = "ViewNotifications"
	ViewReports                                  Permission = "ViewReports"
	ManageTransactionRequests                    Permission = "ManageTransactionRequests"
	ManageIdentityVerificationRequests           Permission = "ManageIdentityVerificationRequests"
	OverruleApprovedIdentityVerificationRequests Permission = "OverruleApprovedIdentityVerificationRequests"
	AdminDeposit                                 Permission = "AdminDeposit"
	AdminWithdraw                                Permission = "AdminWithdraw"
	ManagePlayers                                Permission = "ManagePlayers"
	ManageAdmins                                 Permission = "ManageAdmins"
	ToggleSystemBanks                            Permission = "ToggleSystemBanks"
	ManageSystemBanks                            Permission = "ManageSystemBanks"
	ViewBetHistory                               Permission = "ViewBetHistory"
	ViewActivityLog                              Permission = "ViewActivityLog"
	ViewSuperadminActivityLog                    Permission = "ViewSuperadminActivityLog"
	ViewFiles                                    Permission = "ViewFiles"
)

var Acl = map[db.Role][]Permission{
	db.RolePLAYER: {
		UpdateProfile,
		DepositToOwnWallet,
		WithdrawFromOwnWallet,
		TransferBetweenWallets,
		ViewOwnTransactionHistory,
		ManageOwnBankingDetails,
		ViewNotifications,
	},
	db.RoleADMIN: {
		ViewReports,
		ManageTransactionRequests,
		ManageIdentityVerificationRequests,
		ManagePlayers,
		ToggleSystemBanks,
		ViewBetHistory,
		ViewActivityLog,
		ViewFiles,
	},
	db.RoleSUPERADMIN: {
		UpdateProfile,
		DepositToOwnWallet,
		WithdrawFromOwnWallet,
		TransferBetweenWallets,
		ViewOwnTransactionHistory,
		ManageOwnBankingDetails,
		ViewNotifications,
		ViewReports,
		ManageTransactionRequests,
		ManageIdentityVerificationRequests,
		OverruleApprovedIdentityVerificationRequests,
		AdminDeposit,
		AdminWithdraw,
		ManagePlayers,
		ManageAdmins,
		ToggleSystemBanks,
		ManageSystemBanks,
		ViewBetHistory,
		ViewActivityLog,
		ViewSuperadminActivityLog,
		ViewFiles,
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
