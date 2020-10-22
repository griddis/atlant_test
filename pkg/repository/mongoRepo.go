package repository

import (
	"context"
	"errors"
	"strconv"

	"github.com/griddis/atlant_test/configs"
	"github.com/griddis/atlant_test/tools/logging"
)

func NewRepo(ctx context.Context, logger *logging.Logger, cfg *configs.Database) (Repository, error) {
	ctx = logging.WithContext(ctx, logger)
	if cfg.Driver != "mongodb" {
		logger = logger.With("pkg", "inmemory")
		logger.Debug("init", "inmemory")
		return NewInMemoryStore(ctx)
	} else {
		logger = logger.With("pkg", "mongodb")
		logger.Debug("init", "mongodb")
		return NewMongoStore(ctx, "mongodb://"+cfg.User+":"+cfg.Password+"@"+cfg.Host+":"+strconv.Itoa(cfg.Port)+"/"+cfg.DatabaseName+cfg.Args, cfg.DatabaseName, "price")
	}
	return nil, errors.New("no select repository driver")
}
