package connex

// Logger is intentionally minimal so it can adapt to many logging frameworks.
type Logger interface {
	Printf(format string, v ...any)
}

// LogApplied logs normalized pool config and current runtime stats.
func LogApplied(logger Logger, cfg PoolConfig, warnings []string, stats PoolStats) {
	if logger == nil {
		return
	}

	logger.Printf(
		"connex: applied pool config source=%q version=%q max_open=%d max_idle=%d conn_max_lifetime_sec=%d conn_max_idle_time_sec=%d",
		cfg.Source,
		cfg.Version,
		cfg.MaxOpen,
		cfg.MaxIdle,
		cfg.ConnMaxLifetimeSec,
		cfg.ConnMaxIdleTimeSec,
	)

	for _, warning := range warnings {
		logger.Printf("connex: warning: %s", warning)
	}

	logger.Printf(
		"connex: pool stats max_open=%d open=%d in_use=%d idle=%d wait_count=%d wait_duration=%s",
		stats.MaxOpen,
		stats.Open,
		stats.InUse,
		stats.Idle,
		stats.WaitCount,
		stats.WaitDuration,
	)
}
