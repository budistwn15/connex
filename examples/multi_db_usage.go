package examples

import (
	"database/sql"

	"github.com/budistwn15/connex"
)

func setupPools(mainDB, readDB, analyticsDB *sql.DB) error {
	defaultCfg := connex.DefaultConfig()

	mainEnv, err := connex.FromEnv("DB_POOL_")
	if err != nil {
		return err
	}
	readEnv, err := connex.FromEnv("READ_DB_POOL_")
	if err != nil {
		return err
	}
	analyticsEnv, err := connex.FromEnv("ANALYTICS_DB_POOL_")
	if err != nil {
		return err
	}

	mainCfg, _, err := connex.Normalize(connex.Merge(defaultCfg, mainEnv, connex.PoolConfig{Source: "main"}))
	if err != nil {
		return err
	}
	if err := connex.Apply(mainDB, mainCfg); err != nil {
		return err
	}

	readCfg, _, err := connex.Normalize(connex.Merge(defaultCfg, readEnv, connex.PoolConfig{Source: "read-replica"}))
	if err != nil {
		return err
	}
	if err := connex.Apply(readDB, readCfg); err != nil {
		return err
	}

	analyticsCfg, _, err := connex.Normalize(connex.Merge(defaultCfg, analyticsEnv, connex.PoolConfig{Source: "analytics"}))
	if err != nil {
		return err
	}
	if err := connex.Apply(analyticsDB, analyticsCfg); err != nil {
		return err
	}

	return nil
}
