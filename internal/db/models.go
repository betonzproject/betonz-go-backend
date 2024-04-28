// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type BankName string

const (
	BankNameAGD      BankName = "AGD"
	BankNameAYA      BankName = "AYA"
	BankNameCB       BankName = "CB"
	BankNameKBZ      BankName = "KBZ"
	BankNameKBZPAY   BankName = "KBZPAY"
	BankNameOKDOLLAR BankName = "OK_DOLLAR"
	BankNameWAVEPAY  BankName = "WAVE_PAY"
	BankNameYOMA     BankName = "YOMA"
)

func (e *BankName) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = BankName(s)
	case string:
		*e = BankName(s)
	default:
		return fmt.Errorf("unsupported scan type for BankName: %T", src)
	}
	return nil
}

type NullBankName struct {
	BankName BankName `json:"BankName"`
	Valid    bool     `json:"valid"` // Valid is true if BankName is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullBankName) Scan(value interface{}) error {
	if value == nil {
		ns.BankName, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.BankName.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullBankName) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.BankName), nil
}

type EventResult string

const (
	EventResultSUCCESS EventResult = "SUCCESS"
	EventResultFAIL    EventResult = "FAIL"
)

func (e *EventResult) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = EventResult(s)
	case string:
		*e = EventResult(s)
	default:
		return fmt.Errorf("unsupported scan type for EventResult: %T", src)
	}
	return nil
}

type NullEventResult struct {
	EventResult EventResult `json:"EventResult"`
	Valid       bool        `json:"valid"` // Valid is true if EventResult is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullEventResult) Scan(value interface{}) error {
	if value == nil {
		ns.EventResult, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.EventResult.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullEventResult) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.EventResult), nil
}

type EventType string

const (
	EventTypeLOGIN                          EventType = "LOGIN"
	EventTypeREGISTER                       EventType = "REGISTER"
	EventTypePASSWORDRESETREQUEST           EventType = "PASSWORD_RESET_REQUEST"
	EventTypePASSWORDRESETTOKENVERIFICATION EventType = "PASSWORD_RESET_TOKEN_VERIFICATION"
	EventTypePASSWORDRESET                  EventType = "PASSWORD_RESET"
	EventTypeAUTHENTICATION                 EventType = "AUTHENTICATION"
	EventTypeAUTHORIZATION                  EventType = "AUTHORIZATION"
	EventTypePROFILEUPDATE                  EventType = "PROFILE_UPDATE"
	EventTypeUSERNAMECHANGE                 EventType = "USERNAME_CHANGE"
	EventTypePASSWORDCHANGE                 EventType = "PASSWORD_CHANGE"
	EventTypeBANKADD                        EventType = "BANK_ADD"
	EventTypeBANKUPDATE                     EventType = "BANK_UPDATE"
	EventTypeBANKDELETE                     EventType = "BANK_DELETE"
	EventTypeCHANGEUSERSTATUS               EventType = "CHANGE_USER_STATUS"
	EventTypeEMAILVERIFICATION              EventType = "EMAIL_VERIFICATION"
	EventTypeACTIVE                         EventType = "ACTIVE"
	EventTypeTRANSFERWALLET                 EventType = "TRANSFER_WALLET"
	EventTypeRESTOREWALLET                  EventType = "RESTORE_WALLET"
	EventTypeTRANSACTION                    EventType = "TRANSACTION"
	EventTypeFLAG                           EventType = "FLAG"
	EventTypeSYSTEMBANKADD                  EventType = "SYSTEM_BANK_ADD"
	EventTypeSYSTEMBANKUPDATE               EventType = "SYSTEM_BANK_UPDATE"
	EventTypeSYSTEMBANKDELETE               EventType = "SYSTEM_BANK_DELETE"
	EventTypeMAINTENANCEADD                 EventType = "MAINTENANCE_ADD"
	EventTypeMAINTENANCEUPDATE              EventType = "MAINTENANCE_UPDATE"
	EventTypeMAINTENANCEDELETE              EventType = "MAINTENANCE_DELETE"
)

func (e *EventType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = EventType(s)
	case string:
		*e = EventType(s)
	default:
		return fmt.Errorf("unsupported scan type for EventType: %T", src)
	}
	return nil
}

type NullEventType struct {
	EventType EventType `json:"EventType"`
	Valid     bool      `json:"valid"` // Valid is true if EventType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullEventType) Scan(value interface{}) error {
	if value == nil {
		ns.EventType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.EventType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullEventType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.EventType), nil
}

type FlagStatus string

const (
	FlagStatusPENDING    FlagStatus = "PENDING"
	FlagStatusRESOLVED   FlagStatus = "RESOLVED"
	FlagStatusRESTRICTED FlagStatus = "RESTRICTED"
)

func (e *FlagStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FlagStatus(s)
	case string:
		*e = FlagStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for FlagStatus: %T", src)
	}
	return nil
}

type NullFlagStatus struct {
	FlagStatus FlagStatus `json:"FlagStatus"`
	Valid      bool       `json:"valid"` // Valid is true if FlagStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullFlagStatus) Scan(value interface{}) error {
	if value == nil {
		ns.FlagStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.FlagStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullFlagStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.FlagStatus), nil
}

type IdentityVerificationStatus string

const (
	IdentityVerificationStatusVERIFIED   IdentityVerificationStatus = "VERIFIED"
	IdentityVerificationStatusREJECTED   IdentityVerificationStatus = "REJECTED"
	IdentityVerificationStatusPENDING    IdentityVerificationStatus = "PENDING"
	IdentityVerificationStatusINCOMPLETE IdentityVerificationStatus = "INCOMPLETE"
)

func (e *IdentityVerificationStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = IdentityVerificationStatus(s)
	case string:
		*e = IdentityVerificationStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for IdentityVerificationStatus: %T", src)
	}
	return nil
}

type NullIdentityVerificationStatus struct {
	IdentityVerificationStatus IdentityVerificationStatus `json:"IdentityVerificationStatus"`
	Valid                      bool                       `json:"valid"` // Valid is true if IdentityVerificationStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullIdentityVerificationStatus) Scan(value interface{}) error {
	if value == nil {
		ns.IdentityVerificationStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.IdentityVerificationStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullIdentityVerificationStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.IdentityVerificationStatus), nil
}

type NotificationType string

const (
	NotificationTypeTRANSACTION          NotificationType = "TRANSACTION"
	NotificationTypeIDENTITYVERIFICATION NotificationType = "IDENTITY_VERIFICATION"
)

func (e *NotificationType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = NotificationType(s)
	case string:
		*e = NotificationType(s)
	default:
		return fmt.Errorf("unsupported scan type for NotificationType: %T", src)
	}
	return nil
}

type NullNotificationType struct {
	NotificationType NotificationType `json:"NotificationType"`
	Valid            bool             `json:"valid"` // Valid is true if NotificationType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullNotificationType) Scan(value interface{}) error {
	if value == nil {
		ns.NotificationType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.NotificationType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullNotificationType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.NotificationType), nil
}

type PromotionType string

const (
	PromotionTypeINACTIVEBONUS             PromotionType = "INACTIVE_BONUS"
	PromotionTypeFIVEPERCENTUNLIMITEDBONUS PromotionType = "FIVE_PERCENT_UNLIMITED_BONUS"
	PromotionTypeTENPERCENTUNLIMITEDBONUS  PromotionType = "TEN_PERCENT_UNLIMITED_BONUS"
)

func (e *PromotionType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PromotionType(s)
	case string:
		*e = PromotionType(s)
	default:
		return fmt.Errorf("unsupported scan type for PromotionType: %T", src)
	}
	return nil
}

type NullPromotionType struct {
	PromotionType PromotionType `json:"PromotionType"`
	Valid         bool          `json:"valid"` // Valid is true if PromotionType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPromotionType) Scan(value interface{}) error {
	if value == nil {
		ns.PromotionType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PromotionType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPromotionType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PromotionType), nil
}

type Role string

const (
	RolePLAYER     Role = "PLAYER"
	RoleADMIN      Role = "ADMIN"
	RoleSUPERADMIN Role = "SUPERADMIN"
	RoleSYSTEM     Role = "SYSTEM"
)

func (e *Role) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Role(s)
	case string:
		*e = Role(s)
	default:
		return fmt.Errorf("unsupported scan type for Role: %T", src)
	}
	return nil
}

type NullRole struct {
	Role  Role `json:"Role"`
	Valid bool `json:"valid"` // Valid is true if Role is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRole) Scan(value interface{}) error {
	if value == nil {
		ns.Role, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Role.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Role), nil
}

type TransactionStatus string

const (
	TransactionStatusPENDING  TransactionStatus = "PENDING"
	TransactionStatusAPPROVED TransactionStatus = "APPROVED"
	TransactionStatusDECLINED TransactionStatus = "DECLINED"
)

func (e *TransactionStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TransactionStatus(s)
	case string:
		*e = TransactionStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for TransactionStatus: %T", src)
	}
	return nil
}

type NullTransactionStatus struct {
	TransactionStatus TransactionStatus `json:"TransactionStatus"`
	Valid             bool              `json:"valid"` // Valid is true if TransactionStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTransactionStatus) Scan(value interface{}) error {
	if value == nil {
		ns.TransactionStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TransactionStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTransactionStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TransactionStatus), nil
}

type TransactionType string

const (
	TransactionTypeDEPOSIT  TransactionType = "DEPOSIT"
	TransactionTypeWITHDRAW TransactionType = "WITHDRAW"
)

func (e *TransactionType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TransactionType(s)
	case string:
		*e = TransactionType(s)
	default:
		return fmt.Errorf("unsupported scan type for TransactionType: %T", src)
	}
	return nil
}

type NullTransactionType struct {
	TransactionType TransactionType `json:"TransactionType"`
	Valid           bool            `json:"valid"` // Valid is true if TransactionType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTransactionType) Scan(value interface{}) error {
	if value == nil {
		ns.TransactionType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TransactionType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTransactionType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TransactionType), nil
}

type UserStatus string

const (
	UserStatusNORMAL     UserStatus = "NORMAL"
	UserStatusRESTRICTED UserStatus = "RESTRICTED"
)

func (e *UserStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserStatus(s)
	case string:
		*e = UserStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for UserStatus: %T", src)
	}
	return nil
}

type NullUserStatus struct {
	UserStatus UserStatus `json:"UserStatus"`
	Valid      bool       `json:"valid"` // Valid is true if UserStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserStatus) Scan(value interface{}) error {
	if value == nil {
		ns.UserStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserStatus), nil
}

type Bank struct {
	ID            pgtype.UUID        `json:"id"`
	UserId        pgtype.UUID        `json:"userId"`
	Name          BankName           `json:"name"`
	AccountName   string             `json:"accountName"`
	AccountNumber string             `json:"accountNumber"`
	CreatedAt     pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt     pgtype.Timestamptz `json:"updatedAt"`
	Disabled      bool               `json:"disabled"`
}

type Bet struct {
	ID               int32              `json:"id"`
	RefId            string             `json:"refId"`
	EtgUsername      string             `json:"etgUsername"`
	ProviderUsername string             `json:"providerUsername"`
	ProductCode      int32              `json:"productCode"`
	ProductType      int32              `json:"productType"`
	GameId           pgtype.Text        `json:"gameId"`
	Details          string             `json:"details"`
	Turnover         pgtype.Numeric     `json:"turnover"`
	Bet              pgtype.Numeric     `json:"bet"`
	Payout           pgtype.Numeric     `json:"payout"`
	Status           int32              `json:"status"`
	StartTime        pgtype.Timestamptz `json:"startTime"`
	MatchTime        pgtype.Timestamptz `json:"matchTime"`
	EndTime          pgtype.Timestamptz `json:"endTime"`
	SettleTime       pgtype.Timestamptz `json:"settleTime"`
	ProgShare        pgtype.Numeric     `json:"progShare"`
	ProgWin          pgtype.Numeric     `json:"progWin"`
	Commission       pgtype.Numeric     `json:"commission"`
	WinLoss          pgtype.Numeric     `json:"winLoss"`
}

type Event struct {
	ID          int32              `json:"id"`
	SourceIp    pgtype.Text        `json:"sourceIp"`
	UserId      pgtype.UUID        `json:"userId"`
	Type        EventType          `json:"type"`
	Result      EventResult        `json:"result"`
	Reason      pgtype.Text        `json:"reason"`
	Data        map[string]any     `json:"data"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
	HttpRequest HttpRequest        `json:"httpRequest"`
}

type Flag struct {
	UserId       pgtype.UUID        `json:"userId"`
	ModifiedById pgtype.UUID        `json:"modifiedById"`
	Reason       pgtype.Text        `json:"reason"`
	Remarks      pgtype.Text        `json:"remarks"`
	Status       FlagStatus         `json:"status"`
	CreatedAt    pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz `json:"updatedAt"`
}

type IdentityVerificationRequest struct {
	ID           int32                      `json:"id"`
	UserId       pgtype.UUID                `json:"userId"`
	ModifiedById pgtype.UUID                `json:"modifiedById"`
	Status       IdentityVerificationStatus `json:"status"`
	Remarks      pgtype.Text                `json:"remarks"`
	NricFront    string                     `json:"nricFront"`
	NricBack     string                     `json:"nricBack"`
	HolderFace   string                     `json:"holderFace"`
	NricName     string                     `json:"nricName"`
	Nric         string                     `json:"nric"`
	CreatedAt    pgtype.Timestamptz         `json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz         `json:"updatedAt"`
	Dob          pgtype.Date                `json:"dob"`
}

type Maintenance struct {
	ID                int32                            `json:"id"`
	ProductCode       int32                            `json:"productCode"`
	MaintenancePeriod pgtype.Range[pgtype.Timestamptz] `json:"maintenancePeriod"`
	GmtOffsetSecs     int32                            `json:"gmtOffsetSecs"`
	CreatedAt         pgtype.Timestamptz               `json:"createdAt"`
	UpdatedAt         pgtype.Timestamptz               `json:"updatedAt"`
}

type Notification struct {
	ID        int32              `json:"id"`
	UserId    pgtype.UUID        `json:"userId"`
	Type      NotificationType   `json:"type"`
	Message   pgtype.Text        `json:"message"`
	Variables map[string]any     `json:"variables"`
	Read      bool               `json:"read"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
}

type PasswordResetToken struct {
	TokenHash string             `json:"tokenHash"`
	UserId    pgtype.UUID        `json:"userId"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
}

type PrismaMigration struct {
	ID                string             `json:"id"`
	Checksum          string             `json:"checksum"`
	FinishedAt        pgtype.Timestamptz `json:"finished_at"`
	MigrationName     string             `json:"migration_name"`
	Logs              pgtype.Text        `json:"logs"`
	RolledBackAt      pgtype.Timestamptz `json:"rolled_back_at"`
	StartedAt         pgtype.Timestamptz `json:"started_at"`
	AppliedStepsCount int32              `json:"applied_steps_count"`
}

type TransactionRequest struct {
	ID                           int32              `json:"id"`
	UserId                       pgtype.UUID        `json:"userId"`
	ModifiedById                 pgtype.UUID        `json:"modifiedById"`
	BankName                     NullBankName       `json:"bankName"`
	BankAccountName              pgtype.Text        `json:"bankAccountName"`
	BankAccountNumber            pgtype.Text        `json:"bankAccountNumber"`
	BeneficiaryBankAccountName   pgtype.Text        `json:"beneficiaryBankAccountName"`
	BeneficiaryBankAccountNumber pgtype.Text        `json:"beneficiaryBankAccountNumber"`
	Amount                       pgtype.Numeric     `json:"amount"`
	Type                         TransactionType    `json:"type"`
	ReceiptPath                  pgtype.Text        `json:"receiptPath"`
	Status                       TransactionStatus  `json:"status"`
	Remarks                      pgtype.Text        `json:"remarks"`
	CreatedAt                    pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt                    pgtype.Timestamptz `json:"updatedAt"`
	Bonus                        pgtype.Numeric     `json:"bonus"`
	WithdrawBankFees             pgtype.Numeric     `json:"withdrawBankFees"`
	DepositToWallet              pgtype.Int4        `json:"depositToWallet"`
	Promotion                    NullPromotionType  `json:"promotion"`
}

type TurnoverTarget struct {
	ID                   int32              `json:"id"`
	Target               pgtype.Numeric     `json:"target"`
	TransactionRequestId int32              `json:"transactionRequestId"`
	CreatedAt            pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt            pgtype.Timestamptz `json:"updatedAt"`
}

type User struct {
	ID              pgtype.UUID        `json:"id"`
	Username        string             `json:"username"`
	Email           string             `json:"email"`
	PasswordHash    string             `json:"passwordHash"`
	DisplayName     pgtype.Text        `json:"displayName"`
	PhoneNumber     pgtype.Text        `json:"phoneNumber"`
	CreatedAt       pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt       pgtype.Timestamptz `json:"updatedAt"`
	EtgUsername     string             `json:"etgUsername"`
	Role            Role               `json:"role"`
	MainWallet      pgtype.Numeric     `json:"mainWallet"`
	LastUsedBankId  pgtype.UUID        `json:"lastUsedBankId"`
	ProfileImage    pgtype.Text        `json:"profileImage"`
	Status          UserStatus         `json:"status"`
	IsEmailVerified bool               `json:"isEmailVerified"`
	Dob             pgtype.Date        `json:"dob"`
	PendingEmail    pgtype.Text        `json:"pendingEmail"`
	ReferralCode    pgtype.Text        `json:"referralCode"`
	InvitedBy       pgtype.Text        `json:"invitedBy"`
}

type VerificationPin struct {
	Pin       string             `json:"pin"`
	UserId    pgtype.UUID        `json:"userId"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
}

type VerificationToken struct {
	TokenHash    string             `json:"tokenHash"`
	UserId       pgtype.UUID        `json:"userId"`
	CreatedAt    pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz `json:"updatedAt"`
	RegisterInfo *RegisterInfo      `json:"registerInfo"`
}
