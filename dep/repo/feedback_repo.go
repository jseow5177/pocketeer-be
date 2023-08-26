package repo

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type FeedbackRepo interface {
	Create(ctx context.Context, f *entity.Feedback) error
}
