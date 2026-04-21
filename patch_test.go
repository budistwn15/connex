package connex

import "testing"

func TestPatch_ConfigAndMerge_ExplicitZeroDuration(t *testing.T) {
	defaultCfg := DefaultConfig()
	manual := Patch{
		ConnMaxLifetimeSec: Int(0),
		ConnMaxIdleTimeSec: Ptr(0),
	}

	merged := Merge(defaultCfg, manual.Config(), PoolConfig{})
	if merged.ConnMaxLifetimeSec != 0 || merged.ConnMaxIdleTimeSec != 0 {
		t.Fatalf("expected explicit zero duration override, got lifetime=%d idle_time=%d", merged.ConnMaxLifetimeSec, merged.ConnMaxIdleTimeSec)
	}
}

func TestPatch_ConfigAndMerge_SourceOnly_DoesNotOverrideMaxIdle(t *testing.T) {
	defaultCfg := DefaultConfig()
	manual := Patch{Source: String("manual")}

	merged := Merge(defaultCfg, manual.Config(), PoolConfig{})
	if merged.MaxIdle != defaultCfg.MaxIdle {
		t.Fatalf("expected MaxIdle=%d, got %d", defaultCfg.MaxIdle, merged.MaxIdle)
	}
	if merged.Source != "manual" {
		t.Fatalf("expected Source=manual, got %q", merged.Source)
	}
}

func TestNewPatch(t *testing.T) {
	p := NewPatch()
	cfg := p.Config()
	if cfg.MaxOpen != 0 || cfg.MaxIdle != 0 || cfg.ConnMaxLifetimeSec != 0 || cfg.ConnMaxIdleTimeSec != 0 || cfg.Source != "" || cfg.Version != "" {
		t.Fatalf("expected zero config, got %+v", cfg)
	}
}
