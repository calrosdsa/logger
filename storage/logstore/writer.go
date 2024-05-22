package logstore

import (
	"context"
	"logger/model"
)


type Writer interface {
	//Write logs to storage
	// WriteLogs(ctx context.Context,ld []*model.LogRecord)(error)
	WriteLog(ctx context.Context,ld *model.LogRecord)(error)
	// WriteLogs(ctx context.Context,ld []pb.ResourceLogs)(error)
}