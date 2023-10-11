package savesnapshot

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/api/finnhub"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mem"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"

	acuc "github.com/jseow5177/pockteer-be/usecase/account"
)

const DefaultSnapshotType = entity.SnapshotTypeAccount

type JobConfig struct {
	SnapshotType uint64
}

type SaveSnapshot struct {
	cfg JobConfig

	mongo *mongo.Mongo

	accountUseCase acuc.UseCase

	txMgr        repo.TxMgr
	userRepo     repo.UserRepo
	snapshotRepo repo.SnapshotRepo
}

func (c *SaveSnapshot) initFlags() error {
	flagSet := flag.NewFlagSet(fmt.Sprintf("%s %s", filepath.Base(os.Args[0]), os.Args[1]), flag.ExitOnError)

	flagSet.Uint64Var(&c.cfg.SnapshotType, "snapshotType", uint64(DefaultSnapshotType), "type of snapshot")

	if err := flagSet.Parse(os.Args[2:]); err != nil {
		return err
	}

	if err := entity.CheckSnapshotType(uint32(c.cfg.SnapshotType)); err != nil {
		return err
	}

	return nil
}

func (c *SaveSnapshot) Init(ctx context.Context, cfg *config.Config) error {
	var err error

	if err = c.initFlags(); err != nil {
		return err
	}

	// init mongo
	c.mongo, err = mongo.NewMongo(ctx, cfg.Mongo)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to init mongo client, err: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			_ = c.mongo.Close(ctx)
		}
	}()

	c.txMgr = c.mongo

	securityAPI := finnhub.NewFinnHubMgr(cfg.FinnHub)
	quoteRepo, err := mem.NewQuoteMemCache(cfg.QuoteMemCache, securityAPI)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to init quote repo, err: %v", err)
		return err
	}

	// init repos
	c.userRepo = mongo.NewUserMongo(c.mongo)
	c.snapshotRepo = mongo.NewSnapshotMongo(c.mongo)

	exchangeRateRepo, err := mongo.NewExchangeRateMongo(ctx, c.mongo)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to init exchange rate repo, err: %v", err)
		return err
	}

	// init use cases
	c.accountUseCase = acuc.NewAccountUseCase(
		c.mongo, mongo.NewAccountMongo(c.mongo), mongo.NewTransactionMongo(c.mongo), mongo.NewHoldingMongo(c.mongo),
		mongo.NewLotMongo(c.mongo), quoteRepo, mongo.NewSecurityMongo(c.mongo), exchangeRateRepo, c.snapshotRepo,
	)

	return nil
}

func (c *SaveSnapshot) Run(ctx context.Context) error {
	var (
		page  = 1
		limit = 1000
		now   = uint64(time.Now().UnixMilli())
	)

	p := &repo.Paging{
		Limit: goutil.Uint32(uint32(limit)),
		Page:  goutil.Uint32(uint32(page)),
	}
	sps := make([]*entity.Snapshot, 0)

	for {
		us, err := c.userRepo.GetMany(ctx, repo.NewUserFilter(
			repo.WithUserPaging(p),
		))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get users from repo, err: %v", err)
			return err
		}

		for _, u := range us {
			var r interface{}

			switch c.cfg.SnapshotType {
			case uint64(entity.SnapshotTypeAccount):
				ctxWithUser := entity.SetUserToCtx(ctx, u)
				r, err = c.accountUseCase.GetAccounts(ctxWithUser, &acuc.GetAccountsRequest{
					UserID: u.UserID,
				})
				if err != nil {
					log.Ctx(ctx).Error().Msgf("fail to get user account snapshot, userID: %v, err: %v", u.GetUserID(), err)
					return err
				}
			}

			if r != nil {
				record, err := json.Marshal(r)
				if err != nil {
					log.Ctx(ctx).Error().Msgf("fail to get marshal snapshot, type: %v, userID: %v, err: %v",
						c.cfg.SnapshotType, u.GetUserID(), err)
					return err
				}
				sps = append(sps, entity.NewSnapshot(
					u.GetUserID(),
					uint32(c.cfg.SnapshotType),
					entity.WithSnapshotCreateTime(goutil.Uint64(now)),
					entity.WithSnapshotTimestamp(goutil.Uint64(now)),
					entity.WithSnapshotRecord(goutil.String(string(record))),
				))
			}
		}

		if len(us) < limit {
			break
		}

		page++
		p.Page = goutil.Uint32(uint32(page))
	}

	if len(sps) == 0 {
		return nil
	}

	batchCount := (len(sps) + limit - 1) / limit
	if err := c.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		for i := 0; i < batchCount; i++ {
			var (
				start = i * limit
				end   = (i + 1) * limit
			)
			if end > len(sps) {
				end = len(sps)
			}

			batch := sps[start:end]
			if _, err := c.snapshotRepo.CreateMany(txCtx, batch); err != nil {
				log.Ctx(txCtx).Error().Msgf("fail to create snapshots in repo: %v", err)
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c *SaveSnapshot) Clean(ctx context.Context) error {
	return c.mongo.Close(ctx)
}
