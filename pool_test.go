package connex

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"log"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNormalize_Valid(t *testing.T) {
	cfg := PoolConfig{MaxOpen: 10, MaxIdle: 5, ConnMaxLifetimeSec: 120, ConnMaxIdleTimeSec: 30}
	normalized, warnings, err := Normalize(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", warnings)
	}
	if normalized != cfg {
		t.Fatalf("expected normalized to match cfg, got %+v", normalized)
	}
}

func TestNormalize_ClampMaxIdle(t *testing.T) {
	cfg := PoolConfig{MaxOpen: 4, MaxIdle: 9, ConnMaxLifetimeSec: 120, ConnMaxIdleTimeSec: 30}
	normalized, warnings, err := Normalize(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if normalized.MaxIdle != 4 {
		t.Fatalf("expected MaxIdle clamped to 4, got %d", normalized.MaxIdle)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "clamped") {
		t.Fatalf("expected one clamp warning, got %v", warnings)
	}
}

func TestNormalize_ErrorCases(t *testing.T) {
	tests := []PoolConfig{
		{MaxOpen: 0, MaxIdle: 0},
		{MaxOpen: 1, MaxIdle: -1},
		{MaxOpen: 1, MaxIdle: 0, ConnMaxLifetimeSec: -1},
		{MaxOpen: 1, MaxIdle: 0, ConnMaxIdleTimeSec: -1},
	}
	for _, tc := range tests {
		_, _, err := Normalize(tc)
		if err == nil {
			t.Fatalf("expected error for cfg %+v", tc)
		}
	}
}

func TestMerge_Precedence(t *testing.T) {
	defaultCfg := PoolConfig{MaxOpen: 50, MaxIdle: 10, ConnMaxLifetimeSec: 3600, ConnMaxIdleTimeSec: 600, Source: "default", Version: "v1"}
	envCfg := PoolConfig{MaxOpen: 40, MaxIdle: 0, ConnMaxLifetimeSec: 1800, Source: "env"}
	centralCfg := PoolConfig{MaxOpen: 25, Version: "v2", Source: "central"}

	merged := Merge(defaultCfg, envCfg, centralCfg)

	want := PoolConfig{MaxOpen: 25, MaxIdle: 0, ConnMaxLifetimeSec: 1800, ConnMaxIdleTimeSec: 600, Source: "central", Version: "v2"}
	if merged != want {
		t.Fatalf("unexpected merged config\nwant: %+v\n got: %+v", want, merged)
	}
}

func TestApplyAndStats(t *testing.T) {
	registerTestDriver(t)
	db, err := sql.Open("connex_test_driver", "")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	cfg := PoolConfig{MaxOpen: 8, MaxIdle: 3, ConnMaxLifetimeSec: 90, ConnMaxIdleTimeSec: 15}
	if err := Apply(db, cfg); err != nil {
		t.Fatalf("apply: %v", err)
	}

	stats, err := Stats(db)
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if stats.MaxOpen != 8 {
		t.Fatalf("expected MaxOpen=8, got %d", stats.MaxOpen)
	}
}

func TestApply_NilDB(t *testing.T) {
	if err := Apply(nil, PoolConfig{MaxOpen: 1}); err == nil {
		t.Fatal("expected error for nil db")
	}
}

func TestStats_NilDB(t *testing.T) {
	if _, err := Stats(nil); err == nil {
		t.Fatal("expected error for nil db")
	}
}

func TestLogApplied(t *testing.T) {
	var sink strings.Builder
	logger := log.New(&sink, "", 0)

	cfg := PoolConfig{MaxOpen: 1, MaxIdle: 0, ConnMaxLifetimeSec: 60, ConnMaxIdleTimeSec: 30, Source: "test", Version: "v1"}
	warnings := []string{"warn-a", "warn-b"}
	stats := PoolStats{MaxOpen: 1, Open: 1, InUse: 0, Idle: 1}

	LogApplied(logger, cfg, warnings, stats)
	out := sink.String()
	for _, token := range []string{"connex: applied pool config", "warn-a", "warn-b", "pool stats"} {
		if !strings.Contains(out, token) {
			t.Fatalf("expected log output to contain %q, got %q", token, out)
		}
	}
}

type testDriver struct{}
type testConn struct{}
type testStmt struct{}
type testTx struct{}
type testRows struct{}
type testResult struct{}

func (d testDriver) Open(name string) (driver.Conn, error)   { return testConn{}, nil }
func (c testConn) Prepare(query string) (driver.Stmt, error) { return testStmt{}, nil }
func (c testConn) Close() error                              { return nil }
func (c testConn) Begin() (driver.Tx, error)                 { return testTx{}, nil }
func (s testStmt) Close() error                              { return nil }
func (s testStmt) NumInput() int                             { return -1 }
func (s testStmt) Exec(args []driver.Value) (driver.Result, error) {
	return testResult{}, nil
}
func (s testStmt) Query(args []driver.Value) (driver.Rows, error) { return testRows{}, nil }
func (tx testTx) Commit() error                                   { return nil }
func (tx testTx) Rollback() error                                 { return nil }
func (r testRows) Columns() []string                              { return []string{"col"} }
func (r testRows) Close() error                                   { return nil }
func (r testRows) Next(dest []driver.Value) error                 { return io.EOF }
func (r testResult) LastInsertId() (int64, error)                 { return 0, nil }
func (r testResult) RowsAffected() (int64, error)                 { return 0, nil }

var registerOnce sync.Once

func registerTestDriver(t *testing.T) {
	t.Helper()
	registerOnce.Do(func() {
		sql.Register("connex_test_driver", testDriver{})
	})
	// warm-up open/close to ensure driver is usable
	db, err := sql.Open("connex_test_driver", "")
	if err != nil {
		t.Fatalf("register test driver open failed: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("register test driver ping failed: %v", err)
	}
	db.SetConnMaxLifetime(time.Second)
	db.SetConnMaxIdleTime(time.Second)
	_ = db.Close()
}
