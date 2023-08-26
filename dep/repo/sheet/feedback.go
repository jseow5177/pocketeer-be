package sheet

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/sheet/model"
	"github.com/jseow5177/pockteer-be/entity"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type sheetCred struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	Type        string `json:"type"`
}

func toSheetCred(cfg *config.GoogleSheet) *sheetCred {
	pkTmpl := "-----BEGIN PRIVATE KEY-----\n%s\n-----END PRIVATE KEY-----\n"

	return &sheetCred{
		ClientEmail: cfg.ClientEmail,
		PrivateKey:  fmt.Sprintf(pkTmpl, cfg.PrivateKey),
		Type:        cfg.Type,
	}
}

type feedbackSheet struct {
	client     *sheets.Service
	sheetID    string
	writeRange string
}

func NewFeedbackSheet(ctx context.Context, cfg *config.GoogleSheet) (repo.FeedbackRepo, error) {
	b, err := json.Marshal(toSheetCred(cfg))
	if err != nil {
		return nil, err
	}
	client, err := sheets.NewService(ctx, option.WithCredentialsJSON(b), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		return nil, fmt.Errorf("fail to start sheet service, err: %v", err)
	}

	return &feedbackSheet{
		client:     client,
		sheetID:    cfg.SheetID,
		writeRange: cfg.WriteRange,
	}, nil
}

func (fs *feedbackSheet) Create(ctx context.Context, f *entity.Feedback) error {
	var (
		fm = model.ToFeedbackModelFromEntity(f)
		vs = [][]interface{}{{fm.GetUserID(), fm.GetScore(), fm.GetText()}}
	)

	vr := &sheets.ValueRange{
		Values: vs,
	}

	_, err := fs.client.Spreadsheets.Values.
		Append(fs.sheetID, fs.writeRange, vr).
		ValueInputOption("RAW").
		Do()

	return err
}
