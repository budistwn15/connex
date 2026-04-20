package connex

type poolConfigPatch struct {
	MaxOpen            *int
	MaxIdle            *int
	ConnMaxLifetimeSec *int
	ConnMaxIdleTimeSec *int
	Source             *string
	Version            *string
}

func (p poolConfigPatch) toPoolConfig() PoolConfig {
	cfg := PoolConfig{}

	if p.MaxOpen != nil {
		cfg.MaxOpen = *p.MaxOpen
		cfg.setMask |= fieldMaxOpen
	}
	if p.MaxIdle != nil {
		cfg.MaxIdle = *p.MaxIdle
		cfg.setMask |= fieldMaxIdle
	}
	if p.ConnMaxLifetimeSec != nil {
		cfg.ConnMaxLifetimeSec = *p.ConnMaxLifetimeSec
		cfg.setMask |= fieldConnMaxLifetimeSec
	}
	if p.ConnMaxIdleTimeSec != nil {
		cfg.ConnMaxIdleTimeSec = *p.ConnMaxIdleTimeSec
		cfg.setMask |= fieldConnMaxIdleTimeSec
	}
	if p.Source != nil {
		cfg.Source = *p.Source
		cfg.setMask |= fieldSource
	}
	if p.Version != nil {
		cfg.Version = *p.Version
		cfg.setMask |= fieldVersion
	}

	return cfg
}

func intPtr(v int) *int { return &v }

func stringPtr(v string) *string { return &v }
