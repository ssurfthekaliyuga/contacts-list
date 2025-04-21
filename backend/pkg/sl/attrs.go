package sl

import (
	"errors"
	"log/slog"
	"reflect"
)

func Error(err error) slog.Attr {
	if err == nil {
		err = errors.New("sl.Error: got nil error")
	}
	return slog.String("error", err.Error())
}

func Panic(p any) slog.Attr {
	return slog.Any("panic", p)
}

func Struct(s any) slog.Attr {
	name := reflect.TypeOf(s).String()
	return slog.Any(name, s)
}
