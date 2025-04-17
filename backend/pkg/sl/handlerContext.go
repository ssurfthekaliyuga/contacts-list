package sl

import (
	"context"
	"log/slog"
)

type handlerContext struct {
	slog.Handler
}

//todo посмотреть на replace attrs в конструкторе и тут и как отправлять и писать time
//todo мне не нравиться как выглядит поле source оно должно быть покороче

func (h *handlerContext) Handle(ctx context.Context, record slog.Record) error {
	attrs := extractAttrs(ctx)

	for _, attr := range attrs {
		record.AddAttrs(attr)
	}

	return h.Handler.Handle(ctx, record)
}
