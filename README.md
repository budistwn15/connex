# connex

`connex` adalah shared Go package untuk standardisasi database connection pooling lintas service.

## Instalasi

```bash
go get github.com/budistwn15/connex
```

## API Ringkas

- `DefaultConfig() PoolConfig`
- `FromEnv(prefix string) (PoolConfig, error)`
- `FromMap(raw map[string]any) (PoolConfig, error)`
- `FromJSON(data []byte) (PoolConfig, error)`
- `MustFromEnv(prefix string) PoolConfig`
- `MustFromMap(raw map[string]any) PoolConfig`
- `MustFromJSON(data []byte) PoolConfig`
- `Normalize(cfg PoolConfig) (PoolConfig, []string, error)`
- `Apply(sqlDB *sql.DB, cfg PoolConfig) error`
- `Merge(defaultCfg, envCfg, centralCfg PoolConfig) PoolConfig`
- `Stats(sqlDB *sql.DB) (PoolStats, error)`
- `LogApplied(logger Logger, cfg PoolConfig, warnings []string, stats PoolStats)`

## Contoh Pakai (dengan GORM)

```go
package main

import (
	"log"

	"github.com/budistwn15/connex"
	"gorm.io/gorm"
)

func setupPool(gdb *gorm.DB) error {
	sqlDB, err := gdb.DB()
	if err != nil {
		return err
	}

	defaultCfg := connex.DefaultConfig()
	envCfg, err := connex.FromEnv("") // default prefix: DB_POOL_
	if err != nil {
		return err
	}
	centralCfg, err := connex.FromJSON([]byte(`{
		"max_idle": 20,
		"source": "central-config",
		"version": "2026-04-20"
	}`))
	if err != nil {
		return err
	}

	merged := connex.Merge(defaultCfg, envCfg, centralCfg)
	normalized, warnings, err := connex.Normalize(merged)
	if err != nil {
		return err
	}
	if err := connex.Apply(sqlDB, normalized); err != nil {
		return err
	}

	stats, err := connex.Stats(sqlDB)
	if err != nil {
		return err
	}
	connex.LogApplied(log.Default(), normalized, warnings, stats)
	return nil
}
```

### Mapping Env Vars

Default prefix: `DB_POOL_` (bisa custom lewat parameter `prefix`).

- `DB_POOL_MAX_OPEN`
- `DB_POOL_MAX_IDLE`
- `DB_POOL_CONN_MAX_LIFETIME_SEC`
- `DB_POOL_CONN_MAX_IDLE_TIME_SEC`
- `DB_POOL_SOURCE`
- `DB_POOL_VERSION`

### Prefix Custom (Multi-DB)

Gunakan prefix custom saat 1 service punya lebih dari 1 koneksi DB agar env tidak tabrakan.

```go
mainEnv, err := connex.FromEnv("DB_POOL_")
readEnv, err := connex.FromEnv("READ_DB_POOL_")
analyticsEnv, err := connex.FromEnv("ANALYTICS_DB_POOL_")
```

Contoh lengkap:

- [`.env.multi-db.example`](./.env.multi-db.example)
- [`examples/multi_db_usage.go`](./examples/multi_db_usage.go)

### Central Config Adapter

`FromMap` menerima key snake_case atau camelCase:

- `max_open` / `maxOpen`
- `max_idle` / `maxIdle`
- `conn_max_lifetime_sec` / `connMaxLifetimeSec`
- `conn_max_idle_time_sec` / `connMaxIdleTimeSec`
- `source`
- `version`

`FromJSON` adalah wrapper tipis untuk decode JSON lalu delegasi ke `FromMap`.

`Must*` helpers cocok untuk bootstrap/fail-fast startup; fungsi ini panic jika parsing gagal.

## Validasi & Normalisasi

Rules bawaan:

- `MaxOpen >= 1`
- `MaxIdle >= 0`
- `MaxIdle <= MaxOpen` (jika lebih besar akan di-clamp + warning)
- `ConnMaxLifetimeSec >= 0`
- `ConnMaxIdleTimeSec >= 0`

## Merge Precedence

Urutan prioritas:

1. `centralCfg`
2. `envCfg`
3. `defaultCfg`

Catatan `zero-value`:

- Secara umum `0` dianggap `not set`.
- Khusus `MaxIdle`, nilai `0` valid. Implementasi `Merge` memakai strategy internal supaya `MaxIdle=0` tetap bisa override saat layer tersebut memang membawa konfigurasi.

## Migration Guide dari Hardcoded

Sebelum:

```go
sqlDB.SetMaxOpenConns(50)
sqlDB.SetMaxIdleConns(10)
sqlDB.SetConnMaxLifetime(time.Hour)
sqlDB.SetConnMaxIdleTime(10 * time.Minute)
```

Sesudah:

```go
cfg := connex.Merge(connex.DefaultConfig(), envCfg, centralCfg)
cfg, warnings, err := connex.Normalize(cfg)
if err != nil {
	return err
}
if err := connex.Apply(sqlDB, cfg); err != nil {
	return err
}
stats, _ := connex.Stats(sqlDB)
connex.LogApplied(logger, cfg, warnings, stats)
```

## Contoh Patch Integrasi Service

Lihat: [`examples/service_integration_patch.diff`](./examples/service_integration_patch.diff)
