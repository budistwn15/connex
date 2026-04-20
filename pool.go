package connex

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// PoolConfig defines standardized database pool settings.
type PoolConfig struct {
	MaxOpen            int
	MaxIdle            int
	ConnMaxLifetimeSec int
	ConnMaxIdleTimeSec int
	Source             string // Optional source identifier for observability.
	Version            string // Optional config version for observability.
}

// PoolStats is a lightweight view of sql.DB stats for observability.
type PoolStats struct {
	MaxOpen      int
	Open         int
	InUse        int
	Idle         int
	WaitCount    int64
	WaitDuration time.Duration
}

// DefaultConfig returns safe default values.
// Tune these values according to service workload and DB capacity.
func DefaultConfig() PoolConfig {
	return PoolConfig{
		MaxOpen:            50,
		MaxIdle:            10,
		ConnMaxLifetimeSec: 3600,
		ConnMaxIdleTimeSec: 600,
	}
}

// Normalize validates and normalizes the config.
// It returns warnings for non-fatal corrections (e.g., clamping MaxIdle).
func Normalize(cfg PoolConfig) (PoolConfig, []string, error) {
	warnings := make([]string, 0, 1)

	if cfg.MaxOpen < 1 {
		return PoolConfig{}, nil, errors.New("max_open must be >= 1")
	}
	if cfg.MaxIdle < 0 {
		return PoolConfig{}, nil, errors.New("max_idle must be >= 0")
	}
	if cfg.ConnMaxLifetimeSec < 0 {
		return PoolConfig{}, nil, errors.New("conn_max_lifetime_sec must be >= 0")
	}
	if cfg.ConnMaxIdleTimeSec < 0 {
		return PoolConfig{}, nil, errors.New("conn_max_idle_time_sec must be >= 0")
	}

	normalized := cfg
	if normalized.MaxIdle > normalized.MaxOpen {
		normalized.MaxIdle = normalized.MaxOpen
		warnings = append(warnings, fmt.Sprintf("max_idle (%d) > max_open (%d), clamped to %d", cfg.MaxIdle, cfg.MaxOpen, normalized.MaxIdle))
	}

	return normalized, warnings, nil
}

// Apply normalizes and applies pool settings to sql.DB.
func Apply(sqlDB *sql.DB, cfg PoolConfig) error {
	if sqlDB == nil {
		return errors.New("sqlDB is nil")
	}

	normalized, _, err := Normalize(cfg)
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(normalized.MaxOpen)
	sqlDB.SetMaxIdleConns(normalized.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(normalized.ConnMaxLifetimeSec) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(normalized.ConnMaxIdleTimeSec) * time.Second)
	return nil
}

// Merge combines config layers with precedence: central > env > default.
// Zero values mean "not set", except MaxIdle where 0 can be intentional.
func Merge(defaultCfg, envCfg, centralCfg PoolConfig) PoolConfig {
	result := defaultCfg
	applyLayer(&result, envCfg)
	applyLayer(&result, centralCfg)
	return result
}

func applyLayer(dst *PoolConfig, layer PoolConfig) {
	if layer.MaxOpen != 0 {
		dst.MaxOpen = layer.MaxOpen
	}
	if isMaxIdleSet(layer) {
		dst.MaxIdle = layer.MaxIdle
	}
	if layer.ConnMaxLifetimeSec != 0 {
		dst.ConnMaxLifetimeSec = layer.ConnMaxLifetimeSec
	}
	if layer.ConnMaxIdleTimeSec != 0 {
		dst.ConnMaxIdleTimeSec = layer.ConnMaxIdleTimeSec
	}
	if layer.Source != "" {
		dst.Source = layer.Source
	}
	if layer.Version != "" {
		dst.Version = layer.Version
	}
}

// MaxIdle can be 0 as valid value, so we treat it as explicitly set when
// layer carries at least one non-zero/non-empty signal, or MaxIdle itself is non-zero.
func isMaxIdleSet(layer PoolConfig) bool {
	if layer.MaxIdle != 0 {
		return true
	}
	return layer.MaxOpen != 0 ||
		layer.ConnMaxLifetimeSec != 0 ||
		layer.ConnMaxIdleTimeSec != 0 ||
		layer.Source != "" ||
		layer.Version != ""
}

// Stats returns a simplified snapshot of sql.DB runtime pool stats.
func Stats(sqlDB *sql.DB) (PoolStats, error) {
	if sqlDB == nil {
		return PoolStats{}, errors.New("sqlDB is nil")
	}

	raw := sqlDB.Stats()
	return PoolStats{
		MaxOpen:      raw.MaxOpenConnections,
		Open:         raw.OpenConnections,
		InUse:        raw.InUse,
		Idle:         raw.Idle,
		WaitCount:    raw.WaitCount,
		WaitDuration: raw.WaitDuration,
	}, nil
}
