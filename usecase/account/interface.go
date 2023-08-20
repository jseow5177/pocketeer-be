package account

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/holding"
)

var (
	ErrInvalidUpdateMode = errors.New("invalid update mode")
)

type UseCase interface {
	GetAccount(ctx context.Context, req *GetAccountRequest) (*GetAccountResponse, error)
	GetAccounts(ctx context.Context, req *GetAccountsRequest) (*GetAccountsResponse, error)

	CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error)
	CreateAccounts(ctx context.Context, req *CreateAccountsRequest) (*CreateAccountsResponse, error)
	UpdateAccount(ctx context.Context, req *UpdateAccountRequest) (*UpdateAccountResponse, error)
}

type GetAccountRequest struct {
	UserID    *string
	AccountID *string
}

func (m *GetAccountRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetAccountRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (m *GetAccountRequest) ToAccountFilter() *repo.AccountFilter {
	return repo.NewAccountFilter(m.GetUserID(), repo.WithAccountID(m.AccountID))
}

type GetAccountResponse struct {
	Account *entity.Account
}

func (m *GetAccountResponse) GetAccount() *entity.Account {
	if m != nil && m.Account != nil {
		return m.Account
	}
	return nil
}

type GetAccountsRequest struct {
	UserID      *string
	AccountType *uint32
}

func (m *GetAccountsRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetAccountsRequest) GetAccountType() uint32 {
	if m != nil && m.AccountType != nil {
		return *m.AccountType
	}
	return 0
}

func (m *GetAccountsRequest) ToAccountFilter() *repo.AccountFilter {
	return repo.NewAccountFilter(
		m.GetUserID(),
		repo.WitAccountType(m.AccountType),
	)
}

type GetAccountsResponse struct {
	Accounts []*entity.Account
}

func (m *GetAccountsResponse) GetAccounts() []*entity.Account {
	if m != nil && m.Accounts != nil {
		return m.Accounts
	}
	return nil
}

type CreateAccountRequest struct {
	UserID      *string
	AccountName *string
	Balance     *float64
	Note        *string
	AccountType *uint32
	Holdings    []*holding.CreateHoldingRequest
}

func (m *CreateAccountRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *CreateAccountRequest) GetAccountName() string {
	if m != nil && m.AccountName != nil {
		return *m.AccountName
	}
	return ""
}

func (m *CreateAccountRequest) GetBalance() float64 {
	if m != nil && m.Balance != nil {
		return *m.Balance
	}
	return 0
}

func (m *CreateAccountRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *CreateAccountRequest) GetAccountType() uint32 {
	if m != nil && m.AccountType != nil {
		return *m.AccountType
	}
	return 0
}

func (m *CreateAccountRequest) ToAccountEntity() (*entity.Account, error) {
	return entity.NewAccount(
		m.GetUserID(),
		entity.WithAccountName(m.AccountName),
		entity.WithAccountBalance(m.Balance),
		entity.WithAccountType(m.AccountType),
		entity.WithAccountNote(m.Note),
		entity.WithAccountStatus(goutil.Uint32(uint32(entity.AccountStatusNormal))),
	)
}

func (m *CreateAccountRequest) ToCreateHoldingsRequest(accountID string) *holding.CreateHoldingsRequest {
	return &holding.CreateHoldingsRequest{
		UserID:    m.UserID,
		AccountID: goutil.String(accountID),
		Holdings:  m.Holdings,
	}
}

func (m *CreateAccountRequest) GetHoldings() []*holding.CreateHoldingRequest {
	if m != nil && m.Holdings != nil {
		return m.Holdings
	}
	return nil
}

type CreateAccountResponse struct {
	Account *entity.Account
}

func (m *CreateAccountResponse) GetAccount() *entity.Account {
	if m != nil && m.Account != nil {
		return m.Account
	}
	return nil
}

type CreateAccountsRequest struct {
	Accounts []*CreateAccountRequest
}

func (m *CreateAccountsRequest) GetAccounts() []*CreateAccountRequest {
	if m != nil && m.Accounts != nil {
		return m.Accounts
	}
	return nil
}

type CreateAccountsResponse struct {
	Accounts []*entity.Account
}

func (m *CreateAccountsResponse) GetAccounts() []*entity.Account {
	if m != nil && m.Accounts != nil {
		return m.Accounts
	}
	return nil
}

type UpdateMode uint32

const (
	UpdateModeDefault UpdateMode = iota
	UpdateModeOffsetTransaction
)

var UpdateModes = map[uint32]string{
	uint32(UpdateModeDefault):           "default",
	uint32(UpdateModeOffsetTransaction): "offset transaction",
}

func CheckUpdateMode(updateMode uint32) error {
	checkUpdateMode := updateMode
	for k := range UpdateModes {
		checkUpdateMode = checkUpdateMode & ^uint32(k)
	}
	if checkUpdateMode == 0 {
		return nil
	}

	return ErrInvalidUpdateMode
}

type UpdateAccountRequest struct {
	UserID      *string
	AccountID   *string
	AccountName *string
	Balance     *float64
	Note        *string
	UpdateMode  *uint32
}

func (m *UpdateAccountRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *UpdateAccountRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (m *UpdateAccountRequest) GetAccountName() string {
	if m != nil && m.AccountName != nil {
		return *m.AccountName
	}
	return ""
}

func (m *UpdateAccountRequest) GetBalance() float64 {
	if m != nil && m.Balance != nil {
		return *m.Balance
	}
	return 0
}

func (m *UpdateAccountRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *UpdateAccountRequest) GetUpdateMode() uint32 {
	if m != nil && m.UpdateMode != nil {
		return *m.UpdateMode
	}
	return 0
}

func (m *UpdateAccountRequest) ToAccountFilter() *repo.AccountFilter {
	return repo.NewAccountFilter(m.GetUserID(), repo.WithAccountID(m.AccountID))
}

func (m *UpdateAccountRequest) ToAccountUpdate() *entity.AccountUpdate {
	return entity.NewAccountUpdate(
		entity.WithUpdateAccountName(m.AccountName),
		entity.WithUpdateAccountBalance(m.Balance),
		entity.WithUpdateAccountNote(m.Note),
	)
}

func (m *UpdateAccountRequest) NeedOffsetTransaction() bool {
	return m.GetUpdateMode()&uint32(UpdateModeOffsetTransaction) > 0
}

type UpdateAccountResponse struct {
	Account *entity.Account
}

func (m *UpdateAccountResponse) GetAccount() *entity.Account {
	if m != nil && m.Account != nil {
		return m.Account
	}
	return nil
}
