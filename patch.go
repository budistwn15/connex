package connex

// Patch represents field-set aware overrides.
// Use pointers so explicit zero values are distinguishable from "not set".
type Patch struct {
	MaxOpen            *int
	MaxIdle            *int
	ConnMaxLifetimeSec *int
	ConnMaxIdleTimeSec *int
	Source             *string
	Version            *string
}

// NewPatch creates an empty patch for fluent/manual assignment.
func NewPatch() Patch {
	return Patch{}
}

// Ptr returns a pointer to v. Useful for explicit overrides, including zero values.
func Ptr[T any](v T) *T {
	return &v
}

// Int returns pointer to int.
func Int(v int) *int {
	return &v
}

// String returns pointer to string.
func String(v string) *string {
	return &v
}

// Config materializes a patch into PoolConfig with field-set metadata.
func (p Patch) Config() PoolConfig {
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
