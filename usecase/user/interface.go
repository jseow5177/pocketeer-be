package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/account"
	"github.com/jseow5177/pockteer-be/usecase/category"
	"github.com/jseow5177/pockteer-be/usecase/token"
	"github.com/jseow5177/pockteer-be/util"
)

type UseCase interface {
	GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)
	IsAuthenticated(ctx context.Context, req *IsAuthenticatedRequest) (*IsAuthenticatedResponse, error)
	VerifyEmail(ctx context.Context, req *VerifyEmailRequest) (*VerifyEmailResponse, error)
	InitUser(ctx context.Context, req *InitUserRequest) (*InitUserResponse, error)
	SendOTP(ctx context.Context, req *SendOTPRequest) (*SendOTPResponse, error)
	SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error)
	LogIn(ctx context.Context, req *LogInRequest) (*LogInResponse, error)
	UpdateUserMeta(ctx context.Context, req *UpdateUserMetaRequest) (*UpdateUserMetaResponse, error)
}

type GetUserRequest struct {
	UserID *string
}

func (m *GetUserRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetUserRequest) ToUserFilter() *repo.UserFilter {
	return repo.NewUserFilter(
		repo.WithUserID(m.UserID),
	)
}

type GetUserResponse struct {
	*entity.User
}

func (m *GetUserResponse) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

type SignUpRequest struct {
	Email    *string
	Password *string
}

func (m *SignUpRequest) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *SignUpRequest) GetPassword() string {
	if m != nil && m.Password != nil {
		return *m.Password
	}
	return ""
}

func (m *SignUpRequest) ToUserFilter() *repo.UserFilter {
	return repo.NewUserFilter(
		repo.WithUserEmail(m.Email),
	)
}

func (m *SignUpRequest) ToOTPFilter() *repo.OTPFilter {
	return &repo.OTPFilter{
		Email: m.Email,
	}
}

func (m *SignUpRequest) ToUserEntity() (*entity.User, error) {
	username := util.GetEmailPrefix(m.GetEmail())
	return entity.NewUser(
		m.GetEmail(),
		entity.WithUserPassword(m.Password),
		entity.WithUsername(goutil.String(username)),
	)
}

type SignUpResponse struct {
	User *entity.User
}

func (m *SignUpResponse) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

type LogInRequest struct {
	Email    *string
	Password *string
}

func (m *LogInRequest) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *LogInRequest) GetPassword() string {
	if m != nil && m.Password != nil {
		return *m.Password
	}
	return ""
}

func (m *LogInRequest) ToUserFilter() *repo.UserFilter {
	return repo.NewUserFilter(
		repo.WithUserEmail(m.Email),
	)
}

type LogInResponse struct {
	AccessToken *string
	User        *entity.User
}

func (m *LogInResponse) GetAccessToken() string {
	if m != nil && m.AccessToken != nil {
		return *m.AccessToken
	}
	return ""
}

func (m *LogInResponse) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

type VerifyEmailRequest struct {
	Email *string
	Code  *string
}

func (m *VerifyEmailRequest) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *VerifyEmailRequest) GetCode() string {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return ""
}

func (m *VerifyEmailRequest) ToOTPFilter() *repo.OTPFilter {
	return &repo.OTPFilter{
		Email: m.Email,
	}
}

func (m *VerifyEmailRequest) ToUserFilter(email string) *repo.UserFilter {
	return repo.NewUserFilter(
		repo.WithUserEmail(m.Email),
		repo.WithUserStatus(goutil.Uint32(uint32(entity.UserStatusPending))),
	)
}

type VerifyEmailResponse struct {
	AccessToken *string
	User        *entity.User
}

func (m *VerifyEmailResponse) GetAccessToken() string {
	if m != nil && m.AccessToken != nil {
		return *m.AccessToken
	}
	return ""
}

func (m *VerifyEmailResponse) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

type IsAuthenticatedRequest struct {
	AccessToken *string
}

func (m *IsAuthenticatedRequest) GetAccessToken() string {
	if m != nil && m.AccessToken != nil {
		return *m.AccessToken
	}
	return ""
}

func (m *IsAuthenticatedRequest) ToValidateTokenRequest() *token.ValidateTokenRequest {
	return &token.ValidateTokenRequest{
		TokenType: goutil.Uint32(uint32(entity.TokenTypeAccess)),
		Token:     m.AccessToken,
	}
}

func (m *IsAuthenticatedRequest) ToUserFilter(userID string) *repo.UserFilter {
	return repo.NewUserFilter(
		repo.WithUserID(goutil.String(userID)),
	)
}

type IsAuthenticatedResponse struct {
	User *entity.User
}

func (m *IsAuthenticatedResponse) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

type InitCategoryRequest struct {
}

type InitLotRequest struct {
	Shares       *float64
	CostPerShare *float64
	TradeDate    *uint64
}

func (m *InitLotRequest) GetShares() float64 {
	if m != nil && m.Shares != nil {
		return *m.Shares
	}
	return 0
}

func (m *InitLotRequest) GetCostPerShare() float64 {
	if m != nil && m.CostPerShare != nil {
		return *m.CostPerShare
	}
	return 0
}

func (m *InitLotRequest) GetTradeDate() uint64 {
	if m != nil && m.TradeDate != nil {
		return *m.TradeDate
	}
	return 0
}

func (m *InitLotRequest) ToLotEntity(userID, holdingID, currency string) *entity.Lot {
	return entity.NewLot(
		userID,
		holdingID,
		entity.WithLotCostPerShare(m.CostPerShare),
		entity.WithLotShares(m.Shares),
		entity.WithLotTradeDate(m.TradeDate),
		entity.WithLotCurrency(goutil.String(currency)),
	)
}

type InitHoldingRequest struct {
	Symbol      *string
	HoldingType *uint32
	TotalCost   *float64
	LatestValue *float64
	Lots        []*InitLotRequest
}

func (m *InitHoldingRequest) GetSymbol() string {
	if m != nil && m.Symbol != nil {
		return *m.Symbol
	}
	return ""
}

func (m *InitHoldingRequest) GetHoldingType() uint32 {
	if m != nil && m.HoldingType != nil {
		return *m.HoldingType
	}
	return 0
}

func (m *InitHoldingRequest) GetTotalCost() float64 {
	if m != nil && m.TotalCost != nil {
		return *m.TotalCost
	}
	return 0
}

func (m *InitHoldingRequest) GetLatestValue() float64 {
	if m != nil && m.LatestValue != nil {
		return *m.LatestValue
	}
	return 0
}

func (m *InitHoldingRequest) GetLots() []*InitLotRequest {
	if m != nil && m.Lots != nil {
		return m.Lots
	}
	return nil
}

func (m *InitHoldingRequest) ToSecurityFilter() *repo.SecurityFilter {
	return repo.NewSecurityFilter(
		repo.WithSecuritySymbol(m.Symbol),
	)
}

func (m *InitHoldingRequest) ToHoldingEntity(userID, accountID, currency string) (*entity.Holding, error) {
	return entity.NewHolding(
		userID,
		accountID,
		m.GetSymbol(),
		entity.WithHoldingType(m.HoldingType),
		entity.WithHoldingTotalCost(m.TotalCost),
		entity.WithHoldingLatestValue(m.LatestValue),
	)
}

func (m *InitHoldingRequest) ToLotEntities(userID, holdingID, currency string) []*entity.Lot {
	ls := make([]*entity.Lot, 0)
	for _, r := range m.Lots {
		ls = append(ls, r.ToLotEntity(userID, holdingID, currency))
	}
	return ls
}

type InitAccountRequest struct {
	AccountName *string
	Balance     *float64
	Note        *string
	AccountType *uint32
	Holdings    []*InitHoldingRequest
}

func (m *InitAccountRequest) GetAccountName() string {
	if m != nil && m.AccountName != nil {
		return *m.AccountName
	}
	return ""
}

func (m *InitAccountRequest) GetBalance() float64 {
	if m != nil && m.Balance != nil {
		return *m.Balance
	}
	return 0
}

func (m *InitAccountRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *InitAccountRequest) GetAccountType() uint32 {
	if m != nil && m.AccountType != nil {
		return *m.AccountType
	}
	return 0
}

func (m *InitAccountRequest) GetHoldings() []*InitHoldingRequest {
	if m != nil && m.Holdings != nil {
		return m.Holdings
	}
	return nil
}

func (m *InitAccountRequest) ToAccountEntity(userID, currency string) (*entity.Account, error) {
	return entity.NewAccount(
		userID,
		entity.WithAccountName(m.AccountName),
		entity.WithAccountBalance(m.Balance),
		entity.WithAccountType(m.AccountType),
		entity.WithAccountNote(m.Note),
		entity.WithAccountCurrency(goutil.String(currency)),
	)
}

type InitUserRequest struct {
	UserID     *string
	Currency   *string
	Accounts   []*account.CreateAccountRequest
	Categories []*category.CreateCategoryRequest
}

func (m *InitUserRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *InitUserRequest) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
	}
	return ""
}

func (m *InitUserRequest) GetAccounts() []*account.CreateAccountRequest {
	if m != nil && m.Accounts != nil {
		return m.Accounts
	}
	return nil
}

func (m *InitUserRequest) GetCategories() []*category.CreateCategoryRequest {
	if m != nil && m.Categories != nil {
		return m.Categories
	}
	return nil
}

func (m *InitUserRequest) ToAccountEntities() ([]*entity.Account, error) {
	acs := make([]*entity.Account, 0)
	for _, r := range m.Accounts {
		ac, err := r.ToAccountEntity()
		if err != nil {
			return nil, err
		}
		acs = append(acs, ac)
	}

	return acs, nil
}

func (m *InitUserRequest) ToCategoryEntities() ([]*entity.Category, error) {
	cs := make([]*entity.Category, 0)
	for _, r := range m.Categories {
		c, err := r.ToCategoryEntity()
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

func (m *InitUserRequest) ToUserFilter() *repo.UserFilter {
	return repo.NewUserFilter(
		repo.WithUserID(m.UserID),
	)
}

type InitUserResponse struct{}

type SendOTPRequest struct {
	Email *string
}

func (m *SendOTPRequest) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *SendOTPRequest) ToUserFilter() *repo.UserFilter {
	return repo.NewUserFilter(
		repo.WithUserEmail(m.Email),
	)
}

type SendOTPResponse struct {
	Email *string
}

type UpdateUserMetaRequest struct {
	UserID   *string
	Currency *string
	HideInfo *bool
}

func (m *UpdateUserMetaRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *UpdateUserMetaRequest) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
	}
	return ""
}

func (m *UpdateUserMetaRequest) GetHideInfo() bool {
	if m != nil && m.HideInfo != nil {
		return *m.HideInfo
	}
	return false
}

func (m *UpdateUserMetaRequest) ToUserFilter() *repo.UserFilter {
	return repo.NewUserFilter(
		repo.WithUserID(m.UserID),
	)
}

type UpdateUserMetaResponse struct {
	User *entity.User
}

func (m *UpdateUserMetaResponse) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}
