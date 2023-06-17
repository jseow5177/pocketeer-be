package presenter

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/account"
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
	return &account.CreateAccountRequest{
		UserID:      goutil.String(userID),
		AccountName: m.AccountName,
		Balance:     m.Balance,
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
