package logstore

import (
	"context"
	pb "logger/model/proto"
)


type Writer interface {
	//Write logs to storage
	WriteLogs(ctx context.Context,ld []pb.ResourceLogs)(error)
}