package connex

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const defaultEnvPrefix = "DB_POOL_"

// FromEnv loads partial pool config from environment variables.
// Empty prefix uses default prefix DB_POOL_.
// Supported keys: MAX_OPEN, MAX_IDLE, CONN_MAX_LIFETIME_SEC, CONN_MAX_IDLE_TIME_SEC, SOURCE, VERSION.
func FromEnv(prefix string) (PoolConfig, error) {
	if strings.TrimSpace(prefix) == "" {
		prefix = defaultEnvPrefix
	}

	patch := NewPatch()
	var err error
	var found bool
	var parsed int

	if parsed, found, err = lookupInt(prefix + "MAX_OPEN"); err != nil {
		return PoolConfig{}, err
	}
	if found {
		patch.MaxOpen = Int(parsed)
	}
	if parsed, found, err = lookupInt(prefix + "MAX_IDLE"); err != nil {
		return PoolConfig{}, err
	}
	if found {
		patch.MaxIdle = Int(parsed)
	}
	if parsed, found, err = lookupInt(prefix + "CONN_MAX_LIFETIME_SEC"); err != nil {
		return PoolConfig{}, err
	}
	if found {
		patch.ConnMaxLifetimeSec = Int(parsed)
	}
	if parsed, found, err = lookupInt(prefix + "CONN_MAX_IDLE_TIME_SEC"); err != nil {
		return PoolConfig{}, err
	}
	if found {
		patch.ConnMaxIdleTimeSec = Int(parsed)
	}

	if v, ok := os.LookupEnv(prefix + "SOURCE"); ok {
		patch.Source = String(strings.TrimSpace(v))
	}
	if v, ok := os.LookupEnv(prefix + "VERSION"); ok {
		patch.Version = String(strings.TrimSpace(v))
	}

	return patch.Config(), nil
}

func lookupInt(key string) (int, bool, error) {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return 0, false, nil
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, true, fmt.Errorf("%s is set but empty", key)
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return 0, true, fmt.Errorf("%s must be integer: %w", key, err)
	}
	return parsed, true, nil
}
