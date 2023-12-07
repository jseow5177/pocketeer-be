package presenter

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/account"
	"github.com/jseow5177/pockteer-be/usecase/category"
	"github.com/jseow5177/pockteer-be/usecase/user"
)

type UserMeta struct {
	Currency *string `json:"currency,omitempty"`
	HideInfo *bool   `json:"hide_info,omitempty"`
}

func (um *UserMeta) GetHideInfo() bool {
	if um != nil && um.HideInfo != nil {
		return *um.HideInfo
	}
	return false
}

func (um *UserMeta) GetCurrency() string {
	if um != nil && um.Currency != nil {
		return *um.Currency
	}
	return ""
}

type User struct {
	UserID     *string   `json:"user_id,omitempty"`
	Email      *string   `json:"email,omitempty"`
	Username   *string   `json:"username,omitempty"`
	UserFlag   *uint32   `json:"user_flag,omitempty"`
	UserStatus *uint32   `json:"user_status,omitempty"`
	CreateTime *uint64   `json:"create_time,omitempty"`
	UpdateTime *uint64   `json:"update_time,omitempty"`
	Meta       *UserMeta `json:"meta,omitempty"`
}

func (u *User) GetUserID() string {
	if u != nil && u.UserID != nil {
		return *u.UserID
	}
	return ""
}

func (u *User) GetUsername() string {
	if u != nil && u.Username != nil {
		return *u.Username
	}
	return ""
}

func (u *User) GetEmail() string {
	if u != nil && u.Email != nil {
		return *u.Email
	}
	return ""
}

func (u *User) GetUserFlag() uint32 {
	if u != nil && u.UserFlag != nil {
		return *u.UserFlag
	}
	return 0
}

func (u *User) GetUserStatus() uint32 {
	if u != nil && u.UserStatus != nil {
		return *u.UserStatus
	}
	return 0
}

func (u *User) GetCreateTime() uint64 {
	if u != nil && u.CreateTime != nil {
		return *u.CreateTime
	}
	return 0
}

func (u *User) GetUpdateTime() uint64 {
	if u != nil && u.UpdateTime != nil {
		return *u.UpdateTime
	}
	return 0
}

func (u *User) GetMeta() *UserMeta {
	if u != nil && u.Meta != nil {
		return u.Meta
	}
	return nil
}

type GetUserRequest struct{}

func (m *GetUserRequest) ToUseCaseReq(userID string) *user.GetUserRequest {
	return &user.GetUserRequest{
		UserID: goutil.String(userID),
	}
}

type GetUserResponse struct {
	User *User `json:"user,omitempty"`
}

func (m *GetUserResponse) GetUser() *User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

func (m *GetUserResponse) Set(useCaseRes *user.GetUserResponse) {
	m.User = toUser(useCaseRes.User)
}

type SignUpRequest struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
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

func (m *SignUpRequest) ToUseCaseReq() *user.SignUpRequest {
	return &user.SignUpRequest{
		Email:    m.Email,
		Password: m.Password,
	}
}

type SignUpResponse struct {
	AccessToken *string `json:"access_token,omitempty"`
	User        *User   `json:"user,omitempty"`
}

func (m *SignUpResponse) GetAccessToken() string {
	if m != nil && m.AccessToken != nil {
		return *m.AccessToken
	}
	return ""
}

func (m *SignUpResponse) GetUser() *User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

func (m *SignUpResponse) Set(useCaseRes *user.SignUpResponse) {
	m.AccessToken = useCaseRes.AccessToken
	m.User = toUser(useCaseRes.User)
}

type LogInRequest struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
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

func (m *LogInRequest) ToUseCaseReq() *user.LogInRequest {
	return &user.LogInRequest{
		Email:    m.Email,
		Password: m.Password,
	}
}

func (m *LogInResponse) Set(useCaseRes *user.LogInResponse) {
	m.AccessToken = useCaseRes.AccessToken
	m.User = toUser(useCaseRes.User)
}

type LogInResponse struct {
	AccessToken *string `json:"access_token,omitempty"`
	User        *User   `json:"user,omitempty"`
}

func (m *LogInResponse) GetAccessToken() string {
	if m != nil && m.AccessToken != nil {
		return *m.AccessToken
	}
	return ""
}

func (m *LogInResponse) GetUser() *User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

type VerifyEmailRequest struct {
	Email *string `json:"email,omitempty"`
	Code  *string `json:"code,omitempty"`
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

func (m *VerifyEmailRequest) ToUseCaseReq() *user.VerifyEmailRequest {
	return &user.VerifyEmailRequest{
		Email: m.Email,
		Code:  m.Code,
	}
}

type VerifyEmailResponse struct {
	AccessToken *string `json:"access_token,omitempty"`
	User        *User   `json:"user,omitempty"`
}

func (m *VerifyEmailResponse) Set(useCaseRes *user.VerifyEmailResponse) {
	m.AccessToken = useCaseRes.AccessToken
	m.User = toUser(useCaseRes.User)
}

type InitUserRequest struct {
	Currency   *string                  `json:"currency,omitempty"`
	Accounts   []*CreateAccountRequest  `json:"accounts,omitempty"`
	Categories []*CreateCategoryRequest `json:"categories,omitempty"`
}

func (m *InitUserRequest) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
	}
	return ""
}

func (m *InitUserRequest) ToUseCaseReq(userID string) *user.InitUserRequest {
	acs := make([]*account.CreateAccountRequest, 0)
	for _, r := range m.Accounts {
		r.Currency = m.Currency // TODO: Support currency on account creation
		acs = append(acs, r.ToUseCaseReq(userID))
	}

	cs := make([]*category.CreateCategoryRequest, 0)
	for _, r := range m.Categories {
		if r.Budget != nil {
			// TODO: Support currency on budget creation
			r.Budget.Currency = m.Currency

			// Default to all time
			r.Budget.BudgetRepeat = goutil.Uint32(uint32(entity.BudgetRepeatAllTime))
		}
		cs = append(cs, r.ToUseCaseReq(userID))
	}

	return &user.InitUserRequest{
		UserID:     goutil.String(userID),
		Currency:   m.Currency,
		Accounts:   acs,
		Categories: cs,
	}
}

type InitUserResponse struct{}

func (m *InitUserResponse) Set(useCaseRes *user.InitUserResponse) {}

type SendOTPRequest struct {
	Email *string `json:"email,omitempty"`
}

func (m *SendOTPRequest) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *SendOTPRequest) ToUseCaseReq() *user.SendOTPRequest {
	return &user.SendOTPRequest{
		Email: m.Email,
	}
}

type SendOTPResponse struct{}

func (m *SendOTPResponse) Set(useCaseRes *user.SendOTPResponse) {}

type UpdateUserMetaRequest struct {
	Currency *string `json:"currency,omitempty"`
	HideInfo *bool   `json:"hide_info,omitempty"`
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

func (m *UpdateUserMetaRequest) ToUseCaseReq(userID string) *user.UpdateUserMetaRequest {
	return &user.UpdateUserMetaRequest{
		UserID:   goutil.String(userID),
		Currency: m.Currency,
		HideInfo: m.HideInfo,
	}
}

type UpdateUserMetaResponse struct {
	User *User `json:"user,omitempty"`
}

func (m *UpdateUserMetaResponse) GetUser() *User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

func (m *UpdateUserMetaResponse) Set(useCaseRes *user.UpdateUserMetaResponse) {
	m.User = toUser(useCaseRes.User)
}
