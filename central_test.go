package connex

import "testing"

func TestFromMap_SnakeCase(t *testing.T) {
	raw := map[string]any{
		"max_open":               80,
		"max_idle":               "10",
		"conn_max_lifetime_sec":  float64(1200),
		"conn_max_idle_time_sec": 300,
		"source":                 "central",
		"version":                "v1",
	}

	cfg, err := FromMap(raw)
	if err != nil {
		t.Fatalf("from map: %v", err)
	}

	if cfg.MaxOpen != 80 || cfg.MaxIdle != 10 || cfg.ConnMaxLifetimeSec != 1200 || cfg.ConnMaxIdleTimeSec != 300 || cfg.Source != "central" || cfg.Version != "v1" {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
}

func TestFromMap_CamelCase(t *testing.T) {
	raw := map[string]any{
		"maxOpen":            50,
		"maxIdle":            5,
		"connMaxLifetimeSec": 3600,
		"connMaxIdleTimeSec": 600,
	}

	cfg, err := FromMap(raw)
	if err != nil {
		t.Fatalf("from map: %v", err)
	}

	if cfg.MaxOpen != 50 || cfg.MaxIdle != 5 || cfg.ConnMaxLifetimeSec != 3600 || cfg.ConnMaxIdleTimeSec != 600 {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
}

func TestFromMap_InvalidType(t *testing.T) {
	raw := map[string]any{"max_open": "abc"}
	_, err := FromMap(raw)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFromJSON_Valid(t *testing.T) {
	cfg, err := FromJSON([]byte(`{"max_open": 70, "max_idle": 7, "source": "remote", "version": "2026-04-20"}`))
	if err != nil {
		t.Fatalf("from json: %v", err)
	}
	if cfg.MaxOpen != 70 || cfg.MaxIdle != 7 || cfg.Source != "remote" || cfg.Version != "2026-04-20" {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
}

func TestFromJSON_Invalid(t *testing.T) {
	_, err := FromJSON([]byte(`{"max_open":`))
	if err == nil {
		t.Fatal("expected json error")
	}
}

func TestFromJSON_Empty(t *testing.T) {
	cfg, err := FromJSON(nil)
	if err != nil {
		t.Fatalf("from json: %v", err)
	}
	if cfg.MaxOpen != 0 || cfg.MaxIdle != 0 || cfg.ConnMaxLifetimeSec != 0 || cfg.ConnMaxIdleTimeSec != 0 || cfg.Source != "" || cfg.Version != "" {
		t.Fatalf("expected zero config, got %+v", cfg)
	}
}
