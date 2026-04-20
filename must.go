package connex

// MustFromEnv loads config from env and panics on error.
func MustFromEnv(prefix string) PoolConfig {
	cfg, err := FromEnv(prefix)
	if err != nil {
		panic(err)
	}
	return cfg
}

// MustFromMap loads config from map and panics on error.
func MustFromMap(raw map[string]any) PoolConfig {
	cfg, err := FromMap(raw)
	if err != nil {
		panic(err)
	}
	return cfg
}

// MustFromJSON loads config from JSON and panics on error.
func MustFromJSON(data []byte) PoolConfig {
	cfg, err := FromJSON(data)
	if err != nil {
		panic(err)
	}
	return cfg
}
