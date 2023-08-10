package brevo

import (
	"context"
	"fmt"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/mailer"

	brevoGo "github.com/getbrevo/brevo-go/lib"
)

type BrevoMgr struct {
	client *brevoGo.APIClient
}

func NewBrevoMgr(cfg *config.Brevo) (*BrevoMgr, error) {
	brevoCfg := brevoGo.NewConfiguration()
	brevoCfg.AddDefaultHeader("api-key", cfg.APIKey)

	return &BrevoMgr{
		client: brevoGo.NewAPIClient(brevoCfg),
	}, nil
}

func (mgr *BrevoMgr) SendEmail(ctx context.Context, templateID mailer.Template, req *mailer.SendEmailRequest) error {
	var params interface{} = req.Params

	_, res, err := mgr.client.TransactionalEmailsApi.SendTransacEmail(ctx, brevoGo.SendSmtpEmail{
		TemplateId: int64(templateID),
		To: []brevoGo.SendSmtpEmailTo{{
			Email: req.To,
		}},
		Params: &params,
	})
	if err != nil {
		return fmt.Errorf("fail to send email, code: %v, err: %v", res.StatusCode, err)
	}
	return nil
}
