// admin_extensions.go — hand-written extension types and the AccountingAdminService gRPC service.
// Uses legacy proto v1 struct-tag encoding compatible with grpc v1.60+.
package accountingv1

import (
	"context"
	"fmt"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// ─── AccountingService extension types ──────────────────────────────────────

// DayCutHistoryEntry holds the aggregated day-cut status for a single (cut_date, run_id) across all shards.
type DayCutHistoryEntry struct {
	CutDate     string `protobuf:"bytes,1,opt,name=cut_date,json=cutDate,proto3"           json:"cut_date,omitempty"`
	RunId       int32  `protobuf:"varint,2,opt,name=run_id,json=runId,proto3"              json:"run_id,omitempty"`
	TotalShards int32  `protobuf:"varint,3,opt,name=total_shards,json=totalShards,proto3"  json:"total_shards,omitempty"`
	Pending     int32  `protobuf:"varint,4,opt,name=pending,proto3"                        json:"pending,omitempty"`
	Processing  int32  `protobuf:"varint,5,opt,name=processing,proto3"                     json:"processing,omitempty"`
	Completed   int32  `protobuf:"varint,6,opt,name=completed,proto3"                      json:"completed,omitempty"`
	Failed      int32  `protobuf:"varint,7,opt,name=failed,proto3"                         json:"failed,omitempty"`
}

func (x *DayCutHistoryEntry) Reset()        { *x = DayCutHistoryEntry{} }
func (x *DayCutHistoryEntry) String() string { return fmt.Sprintf("%+v", *x) }
func (*DayCutHistoryEntry) ProtoMessage()    {}

// ListDayCutHistoryRequest is the request for the ListDayCutHistory RPC.
type ListDayCutHistoryRequest struct{}

func (x *ListDayCutHistoryRequest) Reset()        { *x = ListDayCutHistoryRequest{} }
func (x *ListDayCutHistoryRequest) String() string { return fmt.Sprintf("%+v", *x) }
func (*ListDayCutHistoryRequest) ProtoMessage()    {}

// ListDayCutHistoryResponse is the response from the ListDayCutHistory RPC.
type ListDayCutHistoryResponse struct {
	Code    int32                `protobuf:"varint,1,opt,name=code,proto3"   json:"code,omitempty"`
	Message string               `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Entries []*DayCutHistoryEntry `protobuf:"bytes,3,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *ListDayCutHistoryResponse) Reset()        { *x = ListDayCutHistoryResponse{} }
func (x *ListDayCutHistoryResponse) String() string { return fmt.Sprintf("%+v", *x) }
func (*ListDayCutHistoryResponse) ProtoMessage()    {}

// ListSnapshotDatesRequest is the request for the ListSnapshotDates RPC.
type ListSnapshotDatesRequest struct{}

func (x *ListSnapshotDatesRequest) Reset()        { *x = ListSnapshotDatesRequest{} }
func (x *ListSnapshotDatesRequest) String() string { return fmt.Sprintf("%+v", *x) }
func (*ListSnapshotDatesRequest) ProtoMessage()    {}

// ListSnapshotDatesResponse is the response from the ListSnapshotDates RPC.
type ListSnapshotDatesResponse struct {
	Code    int32    `protobuf:"varint,1,opt,name=code,proto3"   json:"code,omitempty"`
	Message string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Dates   []string `protobuf:"bytes,3,rep,name=dates,proto3"   json:"dates,omitempty"`
}

func (x *ListSnapshotDatesResponse) Reset()        { *x = ListSnapshotDatesResponse{} }
func (x *ListSnapshotDatesResponse) String() string { return fmt.Sprintf("%+v", *x) }
func (*ListSnapshotDatesResponse) ProtoMessage()    {}

// ─── AccountingAdminService — maintenance / batch-task RPCs ─────────────────

// ProcessAsyncTasksRequest triggers a batch run of pending async tasks.
type ProcessAsyncTasksRequest struct {
	// BatchSize is the max number of tasks to process per shard per call (default 100).
	BatchSize int32 `protobuf:"varint,1,opt,name=batch_size,json=batchSize,proto3" json:"batch_size,omitempty"`
}

func (x *ProcessAsyncTasksRequest) Reset()        { *x = ProcessAsyncTasksRequest{} }
func (x *ProcessAsyncTasksRequest) String() string { return fmt.Sprintf("%+v", *x) }
func (*ProcessAsyncTasksRequest) ProtoMessage()    {}

// ProcessAsyncTasksResponse reports how many tasks were processed.
type ProcessAsyncTasksResponse struct {
	Code      int32  `protobuf:"varint,1,opt,name=code,proto3"      json:"code,omitempty"`
	Message   string `protobuf:"bytes,2,opt,name=message,proto3"    json:"message,omitempty"`
	Processed int32  `protobuf:"varint,3,opt,name=processed,proto3" json:"processed,omitempty"`
}

func (x *ProcessAsyncTasksResponse) Reset()        { *x = ProcessAsyncTasksResponse{} }
func (x *ProcessAsyncTasksResponse) String() string { return fmt.Sprintf("%+v", *x) }
func (*ProcessAsyncTasksResponse) ProtoMessage()    {}

// DayCutWatchdogRequest triggers a scan for stuck day-cut shards.
type DayCutWatchdogRequest struct {
	// StuckThresholdSeconds is how many seconds a shard must be unchanged to be considered stuck (default 300).
	StuckThresholdSeconds int32 `protobuf:"varint,1,opt,name=stuck_threshold_seconds,json=stuckThresholdSeconds,proto3" json:"stuck_threshold_seconds,omitempty"`
}

func (x *DayCutWatchdogRequest) Reset()        { *x = DayCutWatchdogRequest{} }
func (x *DayCutWatchdogRequest) String() string { return fmt.Sprintf("%+v", *x) }
func (*DayCutWatchdogRequest) ProtoMessage()    {}

// DayCutWatchdogResponse reports the watchdog result.
type DayCutWatchdogResponse struct {
	Code    int32  `protobuf:"varint,1,opt,name=code,proto3"    json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3"  json:"message,omitempty"`
}

func (x *DayCutWatchdogResponse) Reset()        { *x = DayCutWatchdogResponse{} }
func (x *DayCutWatchdogResponse) String() string { return fmt.Sprintf("%+v", *x) }
func (*DayCutWatchdogResponse) ProtoMessage()    {}

// ManualTask describes a single task that needs human intervention.
type ManualTask struct {
	TaskId       string `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3"             json:"task_id,omitempty"`
	TaskType     string `protobuf:"bytes,2,opt,name=task_type,json=taskType,proto3"         json:"task_type,omitempty"`
	BusinessNo   string `protobuf:"bytes,3,opt,name=business_no,json=businessNo,proto3"     json:"business_no,omitempty"`
	ErrorMessage string `protobuf:"bytes,4,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	RetryCount   int32  `protobuf:"varint,5,opt,name=retry_count,json=retryCount,proto3"    json:"retry_count,omitempty"`
}

func (x *ManualTask) Reset()        { *x = ManualTask{} }
func (x *ManualTask) String() string { return fmt.Sprintf("%+v", *x) }
func (*ManualTask) ProtoMessage()    {}

// ListManualTasksRequest requests the list of tasks pending manual processing.
type ListManualTasksRequest struct {
	// Limit is the maximum number of tasks to return (default 50).
	Limit int32 `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *ListManualTasksRequest) Reset()        { *x = ListManualTasksRequest{} }
func (x *ListManualTasksRequest) String() string { return fmt.Sprintf("%+v", *x) }
func (*ListManualTasksRequest) ProtoMessage()    {}

// ListManualTasksResponse returns tasks awaiting manual processing.
type ListManualTasksResponse struct {
	Code    int32         `protobuf:"varint,1,opt,name=code,proto3"   json:"code,omitempty"`
	Message string        `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Tasks   []*ManualTask `protobuf:"bytes,3,rep,name=tasks,proto3"   json:"tasks,omitempty"`
}

func (x *ListManualTasksResponse) Reset()        { *x = ListManualTasksResponse{} }
func (x *ListManualTasksResponse) String() string { return fmt.Sprintf("%+v", *x) }
func (*ListManualTasksResponse) ProtoMessage()    {}

// ─── AccountingAdminServiceClient ───────────────────────────────────────────

// AccountingAdminServiceClient is the client API for AccountingAdminService.
type AccountingAdminServiceClient interface {
	ProcessAsyncTasks(ctx context.Context, in *ProcessAsyncTasksRequest, opts ...grpc.CallOption) (*ProcessAsyncTasksResponse, error)
	DayCutWatchdog(ctx context.Context, in *DayCutWatchdogRequest, opts ...grpc.CallOption) (*DayCutWatchdogResponse, error)
	ListManualTasks(ctx context.Context, in *ListManualTasksRequest, opts ...grpc.CallOption) (*ListManualTasksResponse, error)
}

type accountingAdminServiceClient struct {
	cc grpc.ClientConnInterface
}

// NewAccountingAdminServiceClient creates a new AccountingAdminServiceClient.
func NewAccountingAdminServiceClient(cc grpc.ClientConnInterface) AccountingAdminServiceClient {
	return &accountingAdminServiceClient{cc}
}

func (c *accountingAdminServiceClient) ProcessAsyncTasks(ctx context.Context, in *ProcessAsyncTasksRequest, opts ...grpc.CallOption) (*ProcessAsyncTasksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProcessAsyncTasksResponse)
	if err := c.cc.Invoke(ctx, "/accounting.v1.AccountingAdminService/ProcessAsyncTasks", in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountingAdminServiceClient) DayCutWatchdog(ctx context.Context, in *DayCutWatchdogRequest, opts ...grpc.CallOption) (*DayCutWatchdogResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DayCutWatchdogResponse)
	if err := c.cc.Invoke(ctx, "/accounting.v1.AccountingAdminService/DayCutWatchdog", in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountingAdminServiceClient) ListManualTasks(ctx context.Context, in *ListManualTasksRequest, opts ...grpc.CallOption) (*ListManualTasksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListManualTasksResponse)
	if err := c.cc.Invoke(ctx, "/accounting.v1.AccountingAdminService/ListManualTasks", in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── AccountingAdminServiceServer ───────────────────────────────────────────

// AccountingAdminServiceServer is the server API for AccountingAdminService.
type AccountingAdminServiceServer interface {
	ProcessAsyncTasks(context.Context, *ProcessAsyncTasksRequest) (*ProcessAsyncTasksResponse, error)
	DayCutWatchdog(context.Context, *DayCutWatchdogRequest) (*DayCutWatchdogResponse, error)
	ListManualTasks(context.Context, *ListManualTasksRequest) (*ListManualTasksResponse, error)
	mustEmbedUnimplementedAccountingAdminServiceServer()
}

// UnimplementedAccountingAdminServiceServer provides default implementations.
type UnimplementedAccountingAdminServiceServer struct{}

func (UnimplementedAccountingAdminServiceServer) ProcessAsyncTasks(context.Context, *ProcessAsyncTasksRequest) (*ProcessAsyncTasksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessAsyncTasks not implemented")
}
func (UnimplementedAccountingAdminServiceServer) DayCutWatchdog(context.Context, *DayCutWatchdogRequest) (*DayCutWatchdogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DayCutWatchdog not implemented")
}
func (UnimplementedAccountingAdminServiceServer) ListManualTasks(context.Context, *ListManualTasksRequest) (*ListManualTasksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListManualTasks not implemented")
}
func (UnimplementedAccountingAdminServiceServer) mustEmbedUnimplementedAccountingAdminServiceServer() {
}

// UnsafeAccountingAdminServiceServer may be embedded to opt out of forward compatibility.
type UnsafeAccountingAdminServiceServer interface {
	mustEmbedUnimplementedAccountingAdminServiceServer()
}

// ─── gRPC handler functions ──────────────────────────────────────────────────

func _AccountingAdminService_ProcessAsyncTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessAsyncTasksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountingAdminServiceServer).ProcessAsyncTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accounting.v1.AccountingAdminService/ProcessAsyncTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountingAdminServiceServer).ProcessAsyncTasks(ctx, req.(*ProcessAsyncTasksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountingAdminService_DayCutWatchdog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DayCutWatchdogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountingAdminServiceServer).DayCutWatchdog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accounting.v1.AccountingAdminService/DayCutWatchdog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountingAdminServiceServer).DayCutWatchdog(ctx, req.(*DayCutWatchdogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountingAdminService_ListManualTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListManualTasksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountingAdminServiceServer).ListManualTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accounting.v1.AccountingAdminService/ListManualTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountingAdminServiceServer).ListManualTasks(ctx, req.(*ListManualTasksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountingAdminService_ServiceDesc is the grpc.ServiceDesc for AccountingAdminService.
var AccountingAdminService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "accounting.v1.AccountingAdminService",
	HandlerType: (*AccountingAdminServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcessAsyncTasks",
			Handler:    _AccountingAdminService_ProcessAsyncTasks_Handler,
		},
		{
			MethodName: "DayCutWatchdog",
			Handler:    _AccountingAdminService_DayCutWatchdog_Handler,
		},
		{
			MethodName: "ListManualTasks",
			Handler:    _AccountingAdminService_ListManualTasks_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

// RegisterAccountingAdminServiceServer registers the AccountingAdminService implementation with the gRPC server.
func RegisterAccountingAdminServiceServer(s grpc.ServiceRegistrar, srv AccountingAdminServiceServer) {
	s.RegisterService(&AccountingAdminService_ServiceDesc, srv)
}
