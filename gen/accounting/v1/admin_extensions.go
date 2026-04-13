// admin_extensions.go — ListDayCutHistory / ListSnapshotDates extension types.
// Hand-written structs using legacy proto v1 struct-tag encoding,
// compatible with both grpc v1.60+ (via protoadapt) and grpc v1.64+ (via legacyWrapMessage).
package accountingv1

import "fmt"

// DayCutHistoryEntry holds the aggregated day-cut status for a single (cut_date, run_id) across all shards.
type DayCutHistoryEntry struct {
	CutDate     string `protobuf:"bytes,1,opt,name=cut_date,json=cutDate,proto3"           json:"cut_date,omitempty"`
	RunId       int32  `protobuf:"varint,2,opt,name=run_id,json=runId,proto3"              json:"run_id,omitempty"`
	TotalShards int32  `protobuf:"varint,3,opt,name=total_shards,json=totalShards,proto3"  json:"total_shards,omitempty"`
	Pending     int32  `protobuf:"varint,4,opt,name=pending,proto3"                         json:"pending,omitempty"`
	Processing  int32  `protobuf:"varint,5,opt,name=processing,proto3"                      json:"processing,omitempty"`
	Completed   int32  `protobuf:"varint,6,opt,name=completed,proto3"                       json:"completed,omitempty"`
	Failed      int32  `protobuf:"varint,7,opt,name=failed,proto3"                          json:"failed,omitempty"`
}

func (x *DayCutHistoryEntry) Reset()         { *x = DayCutHistoryEntry{} }
func (x *DayCutHistoryEntry) String() string  { return fmt.Sprintf("%+v", *x) }
func (*DayCutHistoryEntry) ProtoMessage()     {}

// ListDayCutHistoryRequest is the request for the ListDayCutHistory RPC.
type ListDayCutHistoryRequest struct{}

func (x *ListDayCutHistoryRequest) Reset()         { *x = ListDayCutHistoryRequest{} }
func (x *ListDayCutHistoryRequest) String() string  { return fmt.Sprintf("%+v", *x) }
func (*ListDayCutHistoryRequest) ProtoMessage()     {}

// ListDayCutHistoryResponse is the response from the ListDayCutHistory RPC.
type ListDayCutHistoryResponse struct {
	Code    int32                `protobuf:"varint,1,opt,name=code,proto3"    json:"code,omitempty"`
	Message string               `protobuf:"bytes,2,opt,name=message,proto3"  json:"message,omitempty"`
	Entries []*DayCutHistoryEntry `protobuf:"bytes,3,rep,name=entries,proto3"  json:"entries,omitempty"`
}

func (x *ListDayCutHistoryResponse) Reset()         { *x = ListDayCutHistoryResponse{} }
func (x *ListDayCutHistoryResponse) String() string  { return fmt.Sprintf("%+v", *x) }
func (*ListDayCutHistoryResponse) ProtoMessage()     {}

// ListSnapshotDatesRequest is the request for the ListSnapshotDates RPC.
type ListSnapshotDatesRequest struct{}

func (x *ListSnapshotDatesRequest) Reset()         { *x = ListSnapshotDatesRequest{} }
func (x *ListSnapshotDatesRequest) String() string  { return fmt.Sprintf("%+v", *x) }
func (*ListSnapshotDatesRequest) ProtoMessage()     {}

// ListSnapshotDatesResponse is the response from the ListSnapshotDates RPC.
type ListSnapshotDatesResponse struct {
	Code    int32    `protobuf:"varint,1,opt,name=code,proto3"    json:"code,omitempty"`
	Message string   `protobuf:"bytes,2,opt,name=message,proto3"  json:"message,omitempty"`
	Dates   []string `protobuf:"bytes,3,rep,name=dates,proto3"    json:"dates,omitempty"`
}

func (x *ListSnapshotDatesResponse) Reset()         { *x = ListSnapshotDatesResponse{} }
func (x *ListSnapshotDatesResponse) String() string  { return fmt.Sprintf("%+v", *x) }
func (*ListSnapshotDatesResponse) ProtoMessage()     {}
