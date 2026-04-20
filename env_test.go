package connex

import "testing"

func TestFromEnv_DefaultPrefix(t *testing.T) {
	t.Setenv("DB_POOL_MAX_OPEN", "100")
	t.Setenv("DB_POOL_MAX_IDLE", "0")
	t.Setenv("DB_POOL_CONN_MAX_LIFETIME_SEC", "1200")
	t.Setenv("DB_POOL_CONN_MAX_IDLE_TIME_SEC", "300")
	t.Setenv("DB_POOL_SOURCE", "env")
	t.Setenv("DB_POOL_VERSION", "v2026")

	cfg, err := FromEnv("")
	if err != nil {
		t.Fatalf("from env: %v", err)
	}

	if cfg.MaxOpen != 100 || cfg.MaxIdle != 0 || cfg.ConnMaxLifetimeSec != 1200 || cfg.ConnMaxIdleTimeSec != 300 || cfg.Source != "env" || cfg.Version != "v2026" {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
}

func TestFromEnv_CustomPrefix(t *testing.T) {
	t.Setenv("MY_POOL_MAX_OPEN", "30")

	cfg, err := FromEnv("MY_POOL_")
	if err != nil {
		t.Fatalf("from env: %v", err)
	}
	if cfg.MaxOpen != 30 {
		t.Fatalf("expected MaxOpen=30, got %d", cfg.MaxOpen)
	}
}

func TestFromEnv_InvalidInt(t *testing.T) {
	t.Setenv("DB_POOL_MAX_OPEN", "abc")
	_, err := FromEnv("")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestFromEnv_EmptyValue(t *testing.T) {
	t.Setenv("DB_POOL_MAX_OPEN", "")
	_, err := FromEnv("")
	if err == nil {
		t.Fatal("expected empty value error")
	}
}
