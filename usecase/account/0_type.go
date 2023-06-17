package account

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	GetAccount(ctx context.Context, req *GetAccountRequest) (*GetAccountResponse, error)

	CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error)
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
	return &repo.AccountFilter{
		UserID:    m.UserID,
		AccountID: m.AccountID,
	}
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

type CreateAccountRequest struct {
	UserID      *string
	AccountName *string
	Balance     *string
	Note        *string
	AccountType *uint32
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

func (m *CreateAccountRequest) ToAccountEntity() *entity.Account {
	ac := &entity.Account{
		UserID:      m.UserID,
		AccountName: m.AccountName,
		Note:        m.Note,
		AccountType: m.AccountType,
	}

	ac.SetBalance(m.GetBalance())

	return ac
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
