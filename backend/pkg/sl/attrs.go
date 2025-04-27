package sl

import (
	"errors"
	"golang.org/x/exp/constraints"
	"log/slog"
	"reflect"
	"time"
)

type Attr = slog.Attr

func Err(err error) Attr {
	if err == nil {
		err = errors.New("sl.Err: got nil error")
	}
	return slog.String("error", err.Error())
}

func Panic(p any) Attr {
	return slog.Any("panic", p)
}

func Struct(s any) Attr {
	name := reflect.TypeOf(s).String()
	return slog.Any(name, s)
}

func String(key string, value string) Attr {
	return slog.String(key, value)
}

func Int[T constraints.Integer](key string, value T) Attr {
	switch any(value).(type) {
	case int:
		return slog.Int(key, int(value))
	case int8, int16, int32, int64:
		return slog.Int64(key, int64(value))
	case uint, uint8, uint16, uint32, uint64:
		return slog.Uint64(key, uint64(value))
	default:
		return slog.String(key, "unknown type of integer value")
	}
}

func Float[T constraints.Float](key string, value T) Attr {
	return slog.Float64(key, float64(value))
}

func Bool(key string, value bool) Attr {
	return slog.Bool(key, value)
}

func Time(key string, value time.Time) Attr {
	return slog.Time(key, value)
}

func Duration(key string, value time.Duration) Attr {
	return slog.Duration(key, value)
}

func Any(key string, value any) Attr {
	return slog.Any(key, value)
}

func Group(key string, attrs Attr) Attr {
	return slog.Group(key, attrs)
}
