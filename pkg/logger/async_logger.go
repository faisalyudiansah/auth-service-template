package logger

import (
	"context"
	"fmt"
	"sync"
)

type queryLogEvent struct {
	ctx context.Context
	msg string
}

var (
	queryLogCh chan queryLogEvent
	onceLogger sync.Once
)

func InitQueryLogger(buffer int) {
	onceLogger.Do(func() {
		queryLogCh = make(chan queryLogEvent, buffer)

		go func() {
			for e := range queryLogCh {
				logger := loggerFromContextSafe(e.ctx)
				logger.Info(e.msg)
			}
		}()
	})
}

func PushQueryAsyncLog(ctx context.Context, msg string) {
	fmt.Println()
	if queryLogCh == nil {
		loggerFromContextSafe(ctx).Info(msg)
		return
	}
	select {
	case queryLogCh <- queryLogEvent{
		ctx: ctx,
		msg: msg,
	}:
	default:
		// DROP log if channel is full
	}
}

func loggerFromContextSafe(ctx context.Context) Logger {
	if ctx == nil {
		return FromContext(context.Background())
	}
	return FromContext(ctx)
}
