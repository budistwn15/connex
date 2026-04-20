package connex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// FromJSON decodes central config JSON into PoolConfig.
// Accepted payload keys:
// max_open, max_idle, conn_max_lifetime_sec, conn_max_idle_time_sec, source, version.
func FromJSON(data []byte) (PoolConfig, error) {
	if len(data) == 0 {
		return PoolConfig{}, nil
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return PoolConfig{}, fmt.Errorf("decode central pool config json: %w", err)
	}
	return FromMap(raw)
}

// FromMap decodes central config map into PoolConfig.
// Accepted keys:
// max_open, max_idle, conn_max_lifetime_sec, conn_max_idle_time_sec, source, version.
func FromMap(raw map[string]any) (PoolConfig, error) {
	if len(raw) == 0 {
		return PoolConfig{}, nil
	}

	patch := poolConfigPatch{}

	if v, ok := lookupAny(raw, "max_open"); ok {
		parsed, err := toInt(v)
		if err != nil {
			return PoolConfig{}, fmt.Errorf("max_open: %w", err)
		}
		patch.MaxOpen = intPtr(parsed)
	}
	if v, ok := lookupAny(raw, "max_idle"); ok {
		parsed, err := toInt(v)
		if err != nil {
			return PoolConfig{}, fmt.Errorf("max_idle: %w", err)
		}
		patch.MaxIdle = intPtr(parsed)
	}
	if v, ok := lookupAny(raw, "conn_max_lifetime_sec"); ok {
		parsed, err := toInt(v)
		if err != nil {
			return PoolConfig{}, fmt.Errorf("conn_max_lifetime_sec: %w", err)
		}
		patch.ConnMaxLifetimeSec = intPtr(parsed)
	}
	if v, ok := lookupAny(raw, "conn_max_idle_time_sec"); ok {
		parsed, err := toInt(v)
		if err != nil {
			return PoolConfig{}, fmt.Errorf("conn_max_idle_time_sec: %w", err)
		}
		patch.ConnMaxIdleTimeSec = intPtr(parsed)
	}
	if v, ok := lookupAny(raw, "source"); ok {
		s, err := toString(v)
		if err != nil {
			return PoolConfig{}, fmt.Errorf("source: %w", err)
		}
		patch.Source = stringPtr(s)
	}
	if v, ok := lookupAny(raw, "version"); ok {
		s, err := toString(v)
		if err != nil {
			return PoolConfig{}, fmt.Errorf("version: %w", err)
		}
		patch.Version = stringPtr(s)
	}

	return patch.toPoolConfig(), nil
}

func lookupAny(m map[string]any, key string) (any, bool) {
	if v, ok := m[key]; ok {
		return v, true
	}
	camel := snakeToCamel(key)
	if v, ok := m[camel]; ok {
		return v, true
	}
	return nil, false
}

func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		if parts[i] == "" {
			continue
		}
		out += strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return out
}

func toInt(v any) (int, error) {
	switch x := v.(type) {
	case int:
		return x, nil
	case int8:
		return int(x), nil
	case int16:
		return int(x), nil
	case int32:
		return int(x), nil
	case int64:
		return int(x), nil
	case uint:
		return int(x), nil
	case uint8:
		return int(x), nil
	case uint16:
		return int(x), nil
	case uint32:
		return int(x), nil
	case uint64:
		return int(x), nil
	case float32:
		if float32(int(x)) != x {
			return 0, fmt.Errorf("must be integer, got %v", x)
		}
		return int(x), nil
	case float64:
		if float64(int(x)) != x {
			return 0, fmt.Errorf("must be integer, got %v", x)
		}
		return int(x), nil
	case json.Number:
		i, err := x.Int64()
		if err != nil {
			return 0, fmt.Errorf("must be integer, got %q", x)
		}
		return int(i), nil
	case string:
		t := strings.TrimSpace(x)
		if t == "" {
			return 0, fmt.Errorf("must be integer, got empty string")
		}
		i, err := strconv.Atoi(t)
		if err != nil {
			return 0, fmt.Errorf("must be integer, got %q", x)
		}
		return i, nil
	default:
		return 0, fmt.Errorf("must be integer, got %T", v)
	}
}

func toString(v any) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("must be string, got %T", v)
	}
	return strings.TrimSpace(s), nil
}
