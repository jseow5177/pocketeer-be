package presenter

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/account"
	"github.com/jseow5177/pockteer-be/util"
)

type Account struct {
	AccountID     *string `json:"account_id,omitempty"`
	AccountName   *string `json:"account_name,omitempty"`
	Balance       *string `json:"balance,omitempty"`
	AccountType   *uint32 `json:"account_type,omitempty"`
	AccountStatus *uint32 `json:"account_status,omitempty"`
	Note          *string `json:"note,omitempty"`
	CreateTime    *uint64 `json:"create_time,omitempty"`
	UpdateTime    *uint64 `json:"update_time,omitempty"`

	TotalCost *string    `json:"total_cost,omitempty"`
	Holdings  []*Holding `json:"holdings,omitempty"`
}

func (ac *Account) GetAccountID() string {
	if ac != nil && ac.AccountID != nil {
		return *ac.AccountID
	}
	return ""
}

func (ac *Account) GetAccountName() string {
	if ac != nil && ac.AccountName != nil {
		return *ac.AccountName
	}
	return ""
}

func (ac *Account) GetBalance() string {
	if ac != nil && ac.Balance != nil {
		return *ac.Balance
	}
	return ""
}

func (ac *Account) GetAccountStatus() uint32 {
	if ac != nil && ac.AccountStatus != nil {
		return *ac.AccountStatus
	}
	return 0
}

func (ac *Account) GetAccountType() uint32 {
	if ac != nil && ac.AccountType != nil {
		return *ac.AccountType
	}
	return 0
}

func (ac *Account) GetNote() string {
	if ac != nil && ac.Note != nil {
		return *ac.Note
	}
	return ""
}

func (ac *Account) GetCreateTime() uint64 {
	if ac != nil && ac.CreateTime != nil {
		return *ac.CreateTime
	}
	return 0
}

func (ac *Account) GetUpdateTime() uint64 {
	if ac != nil && ac.UpdateTime != nil {
		return *ac.UpdateTime
	}
	return 0
}

func (ac *Account) GetTotalCost() string {
	if ac != nil && ac.TotalCost != nil {
		return *ac.TotalCost
	}
	return ""
}

func (ac *Account) GetHoldings() []*Holding {
	if ac != nil && ac.Holdings != nil {
		return ac.Holdings
	}
	return nil
}

type CreateAccountRequest struct {
	AccountName *string `json:"account_name,omitempty"`
	Balance     *string `json:"balance,omitempty"`
	Note        *string `json:"note,omitempty"`
	AccountType *uint32 `json:"account_type,omitempty"`
}

func (m *CreateAccountRequest) GetAccountName() string {
	if m != nil && m.AccountName != nil {
		return *m.AccountName
	}
	return ""
}

func (m *CreateAccountRequest) GetBalance() string {
	if m != nil && m.Balance != nil {
		return *m.Balance
	}
	return ""
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

func (m *CreateAccountRequest) ToUseCaseReq(userID string) *account.CreateAccountRequest {
	var balance *float64
	if m.Balance != nil {
		b, _ := util.MonetaryStrToFloat(m.GetBalance())
		balance = goutil.Float64(b)
	}
	return &account.CreateAccountRequest{
		UserID:      goutil.String(userID),
		AccountName: m.AccountName,
		Balance:     balance,
		AccountType: m.AccountType,
		Note:        m.Note,
	}
}

type CreateAccountResponse struct {
	Account *Account `json:"account,omitempty"`
}

func (m *CreateAccountResponse) GetAccount() *Account {
	if m != nil && m.Account != nil {
		return m.Account
	}
	return nil
}

func (m *CreateAccountResponse) Set(useCaseRes *account.CreateAccountResponse) {
	m.Account = toAccount(useCaseRes.Account)
}

type GetAccountRequest struct {
	AccountID *string `json:"account_id,omitempty"`
}

func (m *GetAccountRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (m *GetAccountRequest) ToUseCaseReq(userID string) *account.GetAccountRequest {
	return &account.GetAccountRequest{
		UserID:    goutil.String(userID),
		AccountID: m.AccountID,
	}
}

type GetAccountResponse struct {
	Account *Account `json:"account,omitempty"`
}

func (m *GetAccountResponse) GetAccount() *Account {
	if m != nil && m.Account != nil {
		return m.Account
	}
	return nil
}

func (m *GetAccountResponse) Set(useCaseRes *account.GetAccountResponse) {
	m.Account = toAccount(useCaseRes.Account)
}

type GetAccountsRequest struct {
	AccountType *uint32 `json:"account_type,omitempty"`
}

func (m *GetAccountsRequest) GetAccountType() uint32 {
	if m != nil && m.AccountType != nil {
		return *m.AccountType
	}
	return 0
}

func (m *GetAccountsRequest) ToUseCaseReq(userID string) *account.GetAccountsRequest {
	return &account.GetAccountsRequest{
		UserID:      goutil.String(userID),
		AccountType: m.AccountType,
	}
}

type GetAccountsResponse struct {
	Accounts []*Account `json:"accounts,omitempty"`
}

func (m *GetAccountsResponse) GetAccounts() []*Account {
	if m != nil && m.Accounts != nil {
		return m.Accounts
	}
	return nil
}

func (m *GetAccountsResponse) Set(useCaseRes *account.GetAccountsResponse) {
	m.Accounts = toAccounts(useCaseRes.Accounts)
}

type UpdateAccountRequest struct {
	AccountID   *string `json:"account_id,omitempty"`
	AccountName *string `json:"account_name,omitempty"`
	Balance     *string `json:"balance,omitempty"`
	Note        *string `json:"note,omitempty"`
	UpdateMode  *uint32 `json:"update_mode,omitempty"`
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

func (m *UpdateAccountRequest) GetBalance() string {
	if m != nil && m.Balance != nil {
		return *m.Balance
	}
	return ""
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

func (m *UpdateAccountRequest) ToUseCaseReq(userID string) *account.UpdateAccountRequest {
	var balance *float64
	if m.Balance != nil {
		b, _ := util.MonetaryStrToFloat(m.GetBalance())
		balance = goutil.Float64(b)
	}
	return &account.UpdateAccountRequest{
		UserID:      goutil.String(userID),
		AccountID:   m.AccountID,
		AccountName: m.AccountName,
		Balance:     balance,
		Note:        m.Note,
		UpdateMode:  m.UpdateMode,
	}
}

type UpdateAccountResponse struct {
	Account *Account `json:"account,omitempty"`
}

func (m *UpdateAccountResponse) GetAccount() *Account {
	if m != nil && m.Account != nil {
		return m.Account
	}
	return nil
}

func (m *UpdateAccountResponse) Set(useCaseRes *account.UpdateAccountResponse) {
	m.Account = toAccount(useCaseRes.Account)
}
