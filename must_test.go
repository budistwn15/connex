package connex

import "testing"

func TestMustFromEnv_OK(t *testing.T) {
	t.Setenv("DB_POOL_MAX_OPEN", "11")
	cfg := MustFromEnv("")
	if cfg.MaxOpen != 11 {
		t.Fatalf("expected MaxOpen=11, got %d", cfg.MaxOpen)
	}
}

func TestMustFromEnv_Panic(t *testing.T) {
	t.Setenv("DB_POOL_MAX_OPEN", "bad")
	mustPanic(t, func() {
		_ = MustFromEnv("")
	})
}

func TestMustFromMap_OK(t *testing.T) {
	cfg := MustFromMap(map[string]any{"max_open": 22})
	if cfg.MaxOpen != 22 {
		t.Fatalf("expected MaxOpen=22, got %d", cfg.MaxOpen)
	}
}

func TestMustFromMap_Panic(t *testing.T) {
	mustPanic(t, func() {
		_ = MustFromMap(map[string]any{"max_open": "bad"})
	})
}

func TestMustFromJSON_OK(t *testing.T) {
	cfg := MustFromJSON([]byte(`{"max_open":33}`))
	if cfg.MaxOpen != 33 {
		t.Fatalf("expected MaxOpen=33, got %d", cfg.MaxOpen)
	}
}

func TestMustFromJSON_Panic(t *testing.T) {
	mustPanic(t, func() {
		_ = MustFromJSON([]byte(`{"max_open":`))
	})
}

func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}
