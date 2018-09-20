// Code generated by protoc-gen-go. DO NOT EDIT.
// source: alerts.proto

package alerts

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import cap "github.com/alerting/alerts/cap"
import protobuf "github.com/alerting/alerts/protobuf"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// * Messages *
type Coordinate struct {
	Lat                  float64  `protobuf:"fixed64,1,opt,name=lat,proto3" json:"lat,omitempty"`
	Lon                  float64  `protobuf:"fixed64,2,opt,name=lon,proto3" json:"lon,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Coordinate) Reset()         { *m = Coordinate{} }
func (m *Coordinate) String() string { return proto.CompactTextString(m) }
func (*Coordinate) ProtoMessage()    {}
func (*Coordinate) Descriptor() ([]byte, []int) {
	return fileDescriptor_alerts_19aa2e8bfb50d8b5, []int{0}
}
func (m *Coordinate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Coordinate.Unmarshal(m, b)
}
func (m *Coordinate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Coordinate.Marshal(b, m, deterministic)
}
func (dst *Coordinate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Coordinate.Merge(dst, src)
}
func (m *Coordinate) XXX_Size() int {
	return xxx_messageInfo_Coordinate.Size(m)
}
func (m *Coordinate) XXX_DiscardUnknown() {
	xxx_messageInfo_Coordinate.DiscardUnknown(m)
}

var xxx_messageInfo_Coordinate proto.InternalMessageInfo

func (m *Coordinate) GetLat() float64 {
	if m != nil {
		return m.Lat
	}
	return 0
}

func (m *Coordinate) GetLon() float64 {
	if m != nil {
		return m.Lon
	}
	return 0
}

type TimeConditions struct {
	Gte                  *timestamp.Timestamp `protobuf:"bytes,1,opt,name=gte,proto3" json:"gte,omitempty"`
	Gt                   *timestamp.Timestamp `protobuf:"bytes,2,opt,name=gt,proto3" json:"gt,omitempty"`
	Lte                  *timestamp.Timestamp `protobuf:"bytes,3,opt,name=lte,proto3" json:"lte,omitempty"`
	Lt                   *timestamp.Timestamp `protobuf:"bytes,4,opt,name=lt,proto3" json:"lt,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *TimeConditions) Reset()         { *m = TimeConditions{} }
func (m *TimeConditions) String() string { return proto.CompactTextString(m) }
func (*TimeConditions) ProtoMessage()    {}
func (*TimeConditions) Descriptor() ([]byte, []int) {
	return fileDescriptor_alerts_19aa2e8bfb50d8b5, []int{1}
}
func (m *TimeConditions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TimeConditions.Unmarshal(m, b)
}
func (m *TimeConditions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TimeConditions.Marshal(b, m, deterministic)
}
func (dst *TimeConditions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TimeConditions.Merge(dst, src)
}
func (m *TimeConditions) XXX_Size() int {
	return xxx_messageInfo_TimeConditions.Size(m)
}
func (m *TimeConditions) XXX_DiscardUnknown() {
	xxx_messageInfo_TimeConditions.DiscardUnknown(m)
}

var xxx_messageInfo_TimeConditions proto.InternalMessageInfo

func (m *TimeConditions) GetGte() *timestamp.Timestamp {
	if m != nil {
		return m.Gte
	}
	return nil
}

func (m *TimeConditions) GetGt() *timestamp.Timestamp {
	if m != nil {
		return m.Gt
	}
	return nil
}

func (m *TimeConditions) GetLte() *timestamp.Timestamp {
	if m != nil {
		return m.Lte
	}
	return nil
}

func (m *TimeConditions) GetLt() *timestamp.Timestamp {
	if m != nil {
		return m.Lt
	}
	return nil
}

type FindCriteria struct {
	Start int32    `protobuf:"varint,1,opt,name=start,proto3" json:"start,omitempty"`
	Count int32    `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	Sort  []string `protobuf:"bytes,3,rep,name=sort,proto3" json:"sort,omitempty"`
	// cap.Alert
	Superseded    bool                  `protobuf:"varint,4,opt,name=superseded,proto3" json:"superseded,omitempty"`
	NotSuperseded bool                  `protobuf:"varint,5,opt,name=not_superseded,json=notSuperseded,proto3" json:"not_superseded,omitempty"`
	Status        cap.Alert_Status      `protobuf:"varint,6,opt,name=status,proto3,enum=cap.Alert_Status" json:"status,omitempty"`
	MessageType   cap.Alert_MessageType `protobuf:"varint,7,opt,name=message_type,json=messageType,proto3,enum=cap.Alert_MessageType" json:"message_type,omitempty"`
	Scope         cap.Alert_Scope       `protobuf:"varint,8,opt,name=scope,proto3,enum=cap.Alert_Scope" json:"scope,omitempty"`
	// cap.Info
	Language    string             `protobuf:"bytes,9,opt,name=language,proto3" json:"language,omitempty"`
	Certainty   cap.Info_Certainty `protobuf:"varint,10,opt,name=certainty,proto3,enum=cap.Info_Certainty" json:"certainty,omitempty"`
	Severity    cap.Info_Severity  `protobuf:"varint,11,opt,name=severity,proto3,enum=cap.Info_Severity" json:"severity,omitempty"`
	Urgency     cap.Info_Urgency   `protobuf:"varint,12,opt,name=urgency,proto3,enum=cap.Info_Urgency" json:"urgency,omitempty"`
	Headline    string             `protobuf:"bytes,13,opt,name=headline,proto3" json:"headline,omitempty"`
	Description string             `protobuf:"bytes,14,opt,name=description,proto3" json:"description,omitempty"`
	Instruction string             `protobuf:"bytes,15,opt,name=instruction,proto3" json:"instruction,omitempty"`
	Effective   *TimeConditions    `protobuf:"bytes,16,opt,name=effective,proto3" json:"effective,omitempty"`
	Onset       *TimeConditions    `protobuf:"bytes,17,opt,name=onset,proto3" json:"onset,omitempty"`
	Expires     *TimeConditions    `protobuf:"bytes,18,opt,name=expires,proto3" json:"expires,omitempty"`
	// cap.Area
	AreaDescription      string      `protobuf:"bytes,19,opt,name=area_description,json=areaDescription,proto3" json:"area_description,omitempty"`
	Point                *Coordinate `protobuf:"bytes,20,opt,name=point,proto3" json:"point,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *FindCriteria) Reset()         { *m = FindCriteria{} }
func (m *FindCriteria) String() string { return proto.CompactTextString(m) }
func (*FindCriteria) ProtoMessage()    {}
func (*FindCriteria) Descriptor() ([]byte, []int) {
	return fileDescriptor_alerts_19aa2e8bfb50d8b5, []int{2}
}
func (m *FindCriteria) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindCriteria.Unmarshal(m, b)
}
func (m *FindCriteria) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindCriteria.Marshal(b, m, deterministic)
}
func (dst *FindCriteria) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindCriteria.Merge(dst, src)
}
func (m *FindCriteria) XXX_Size() int {
	return xxx_messageInfo_FindCriteria.Size(m)
}
func (m *FindCriteria) XXX_DiscardUnknown() {
	xxx_messageInfo_FindCriteria.DiscardUnknown(m)
}

var xxx_messageInfo_FindCriteria proto.InternalMessageInfo

func (m *FindCriteria) GetStart() int32 {
	if m != nil {
		return m.Start
	}
	return 0
}

func (m *FindCriteria) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *FindCriteria) GetSort() []string {
	if m != nil {
		return m.Sort
	}
	return nil
}

func (m *FindCriteria) GetSuperseded() bool {
	if m != nil {
		return m.Superseded
	}
	return false
}

func (m *FindCriteria) GetNotSuperseded() bool {
	if m != nil {
		return m.NotSuperseded
	}
	return false
}

func (m *FindCriteria) GetStatus() cap.Alert_Status {
	if m != nil {
		return m.Status
	}
	return cap.Alert_STATUS_UNKNOWN
}

func (m *FindCriteria) GetMessageType() cap.Alert_MessageType {
	if m != nil {
		return m.MessageType
	}
	return cap.Alert_MESSAGE_TYPE_UNKNOWN
}

func (m *FindCriteria) GetScope() cap.Alert_Scope {
	if m != nil {
		return m.Scope
	}
	return cap.Alert_SCOPE_UNKNOWN
}

func (m *FindCriteria) GetLanguage() string {
	if m != nil {
		return m.Language
	}
	return ""
}

func (m *FindCriteria) GetCertainty() cap.Info_Certainty {
	if m != nil {
		return m.Certainty
	}
	return cap.Info_CERTAINTY_UNKNOWN
}

func (m *FindCriteria) GetSeverity() cap.Info_Severity {
	if m != nil {
		return m.Severity
	}
	return cap.Info_SEVERITY_UNKNOWN
}

func (m *FindCriteria) GetUrgency() cap.Info_Urgency {
	if m != nil {
		return m.Urgency
	}
	return cap.Info_URGENCY_UNKNOWN
}

func (m *FindCriteria) GetHeadline() string {
	if m != nil {
		return m.Headline
	}
	return ""
}

func (m *FindCriteria) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *FindCriteria) GetInstruction() string {
	if m != nil {
		return m.Instruction
	}
	return ""
}

func (m *FindCriteria) GetEffective() *TimeConditions {
	if m != nil {
		return m.Effective
	}
	return nil
}

func (m *FindCriteria) GetOnset() *TimeConditions {
	if m != nil {
		return m.Onset
	}
	return nil
}

func (m *FindCriteria) GetExpires() *TimeConditions {
	if m != nil {
		return m.Expires
	}
	return nil
}

func (m *FindCriteria) GetAreaDescription() string {
	if m != nil {
		return m.AreaDescription
	}
	return ""
}

func (m *FindCriteria) GetPoint() *Coordinate {
	if m != nil {
		return m.Point
	}
	return nil
}

type Hit struct {
	Id                   string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Alert                *cap.Alert `protobuf:"bytes,2,opt,name=alert,proto3" json:"alert,omitempty"`
	Info                 *cap.Info  `protobuf:"bytes,3,opt,name=info,proto3" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Hit) Reset()         { *m = Hit{} }
func (m *Hit) String() string { return proto.CompactTextString(m) }
func (*Hit) ProtoMessage()    {}
func (*Hit) Descriptor() ([]byte, []int) {
	return fileDescriptor_alerts_19aa2e8bfb50d8b5, []int{3}
}
func (m *Hit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Hit.Unmarshal(m, b)
}
func (m *Hit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Hit.Marshal(b, m, deterministic)
}
func (dst *Hit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hit.Merge(dst, src)
}
func (m *Hit) XXX_Size() int {
	return xxx_messageInfo_Hit.Size(m)
}
func (m *Hit) XXX_DiscardUnknown() {
	xxx_messageInfo_Hit.DiscardUnknown(m)
}

var xxx_messageInfo_Hit proto.InternalMessageInfo

func (m *Hit) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Hit) GetAlert() *cap.Alert {
	if m != nil {
		return m.Alert
	}
	return nil
}

func (m *Hit) GetInfo() *cap.Info {
	if m != nil {
		return m.Info
	}
	return nil
}

type FindResult struct {
	Total                int64    `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Hits                 []*Hit   `protobuf:"bytes,2,rep,name=hits,proto3" json:"hits,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FindResult) Reset()         { *m = FindResult{} }
func (m *FindResult) String() string { return proto.CompactTextString(m) }
func (*FindResult) ProtoMessage()    {}
func (*FindResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_alerts_19aa2e8bfb50d8b5, []int{4}
}
func (m *FindResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindResult.Unmarshal(m, b)
}
func (m *FindResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindResult.Marshal(b, m, deterministic)
}
func (dst *FindResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindResult.Merge(dst, src)
}
func (m *FindResult) XXX_Size() int {
	return xxx_messageInfo_FindResult.Size(m)
}
func (m *FindResult) XXX_DiscardUnknown() {
	xxx_messageInfo_FindResult.DiscardUnknown(m)
}

var xxx_messageInfo_FindResult proto.InternalMessageInfo

func (m *FindResult) GetTotal() int64 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *FindResult) GetHits() []*Hit {
	if m != nil {
		return m.Hits
	}
	return nil
}

func init() {
	proto.RegisterType((*Coordinate)(nil), "alerts.Coordinate")
	proto.RegisterType((*TimeConditions)(nil), "alerts.TimeConditions")
	proto.RegisterType((*FindCriteria)(nil), "alerts.FindCriteria")
	proto.RegisterType((*Hit)(nil), "alerts.Hit")
	proto.RegisterType((*FindResult)(nil), "alerts.FindResult")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AlertsServiceClient is the client API for AlertsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AlertsServiceClient interface {
	// Add a new alert.
	Add(ctx context.Context, in *cap.Alert, opts ...grpc.CallOption) (*cap.Alert, error)
	// Returns the alert that matches the provided reference.
	// NOTE: Either id or (identifier, sender, sent) must be provided.
	Get(ctx context.Context, in *cap.Reference, opts ...grpc.CallOption) (*cap.Alert, error)
	// Returns whether or not an alert matches the provided reference.
	// NOTE: Either id or (identifier, sender, sent) must be provided.
	Has(ctx context.Context, in *cap.Reference, opts ...grpc.CallOption) (*protobuf.BooleanResult, error)
	// Find alerts matching the provided criteria.
	Find(ctx context.Context, in *FindCriteria, opts ...grpc.CallOption) (*FindResult, error)
}

type alertsServiceClient struct {
	cc *grpc.ClientConn
}

func NewAlertsServiceClient(cc *grpc.ClientConn) AlertsServiceClient {
	return &alertsServiceClient{cc}
}

func (c *alertsServiceClient) Add(ctx context.Context, in *cap.Alert, opts ...grpc.CallOption) (*cap.Alert, error) {
	out := new(cap.Alert)
	err := c.cc.Invoke(ctx, "/alerts.AlertsService/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertsServiceClient) Get(ctx context.Context, in *cap.Reference, opts ...grpc.CallOption) (*cap.Alert, error) {
	out := new(cap.Alert)
	err := c.cc.Invoke(ctx, "/alerts.AlertsService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertsServiceClient) Has(ctx context.Context, in *cap.Reference, opts ...grpc.CallOption) (*protobuf.BooleanResult, error) {
	out := new(protobuf.BooleanResult)
	err := c.cc.Invoke(ctx, "/alerts.AlertsService/Has", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertsServiceClient) Find(ctx context.Context, in *FindCriteria, opts ...grpc.CallOption) (*FindResult, error) {
	out := new(FindResult)
	err := c.cc.Invoke(ctx, "/alerts.AlertsService/Find", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AlertsServiceServer is the server API for AlertsService service.
type AlertsServiceServer interface {
	// Add a new alert.
	Add(context.Context, *cap.Alert) (*cap.Alert, error)
	// Returns the alert that matches the provided reference.
	// NOTE: Either id or (identifier, sender, sent) must be provided.
	Get(context.Context, *cap.Reference) (*cap.Alert, error)
	// Returns whether or not an alert matches the provided reference.
	// NOTE: Either id or (identifier, sender, sent) must be provided.
	Has(context.Context, *cap.Reference) (*protobuf.BooleanResult, error)
	// Find alerts matching the provided criteria.
	Find(context.Context, *FindCriteria) (*FindResult, error)
}

func RegisterAlertsServiceServer(s *grpc.Server, srv AlertsServiceServer) {
	s.RegisterService(&_AlertsService_serviceDesc, srv)
}

func _AlertsService_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(cap.Alert)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertsServiceServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alerts.AlertsService/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertsServiceServer).Add(ctx, req.(*cap.Alert))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlertsService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(cap.Reference)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertsServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alerts.AlertsService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertsServiceServer).Get(ctx, req.(*cap.Reference))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlertsService_Has_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(cap.Reference)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertsServiceServer).Has(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alerts.AlertsService/Has",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertsServiceServer).Has(ctx, req.(*cap.Reference))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlertsService_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindCriteria)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertsServiceServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alerts.AlertsService/Find",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertsServiceServer).Find(ctx, req.(*FindCriteria))
	}
	return interceptor(ctx, in, info, handler)
}

var _AlertsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "alerts.AlertsService",
	HandlerType: (*AlertsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _AlertsService_Add_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _AlertsService_Get_Handler,
		},
		{
			MethodName: "Has",
			Handler:    _AlertsService_Has_Handler,
		},
		{
			MethodName: "Find",
			Handler:    _AlertsService_Find_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "alerts.proto",
}

func init() { proto.RegisterFile("alerts.proto", fileDescriptor_alerts_19aa2e8bfb50d8b5) }

var fileDescriptor_alerts_19aa2e8bfb50d8b5 = []byte{
	// 741 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x94, 0x51, 0x6f, 0xdb, 0x36,
	0x10, 0xc7, 0x61, 0xcb, 0x76, 0xe2, 0x73, 0xe2, 0xa6, 0xd7, 0xa0, 0x20, 0x0c, 0xb4, 0xf5, 0x0c,
	0x6c, 0x70, 0xb7, 0x42, 0x49, 0xb3, 0xed, 0x61, 0x8f, 0x9d, 0x87, 0x2d, 0x7b, 0xd8, 0x0b, 0x93,
	0xed, 0x35, 0x60, 0xa4, 0xb3, 0x42, 0x40, 0x26, 0x05, 0xf2, 0x14, 0xcc, 0x5f, 0x65, 0x1f, 0x65,
	0xcf, 0xfb, 0x60, 0x83, 0x48, 0x2b, 0x56, 0x51, 0xc0, 0x79, 0x30, 0xc0, 0xfb, 0xdf, 0xef, 0x8e,
	0x77, 0x16, 0xef, 0xe0, 0x44, 0x95, 0xe4, 0xd8, 0xa7, 0x95, 0xb3, 0x6c, 0x71, 0x14, 0xad, 0xd9,
	0xbb, 0xc2, 0xda, 0xa2, 0xa4, 0x8b, 0xa0, 0xde, 0xd7, 0xeb, 0x0b, 0xd6, 0x1b, 0xf2, 0xac, 0x36,
	0x55, 0x04, 0x67, 0xcb, 0x42, 0xf3, 0x43, 0x7d, 0x9f, 0x66, 0x76, 0x73, 0x11, 0x62, 0xb4, 0x29,
	0xe2, 0xc1, 0x5f, 0x64, 0xaa, 0x6a, 0x7e, 0x3b, 0xf2, 0xe3, 0x01, 0xf2, 0x29, 0xfd, 0xbd, 0xb5,
	0x25, 0x29, 0x13, 0x43, 0x16, 0x97, 0x00, 0x2b, 0x6b, 0x5d, 0xae, 0x8d, 0x62, 0xc2, 0x33, 0x48,
	0x4a, 0xc5, 0xa2, 0x37, 0xef, 0x2d, 0x7b, 0xb2, 0x39, 0x06, 0xc5, 0x1a, 0xd1, 0xdf, 0x29, 0xd6,
	0x2c, 0xfe, 0xeb, 0xc1, 0xf4, 0x56, 0x6f, 0x68, 0x65, 0x4d, 0xae, 0x59, 0x5b, 0xe3, 0xf1, 0x03,
	0x24, 0x05, 0x53, 0x08, 0x9b, 0x5c, 0xcd, 0xd2, 0xd8, 0x50, 0xda, 0xde, 0x98, 0xde, 0xb6, 0x0d,
	0xc9, 0x06, 0xc3, 0x6f, 0xa1, 0x5f, 0x70, 0xc8, 0x78, 0x18, 0xee, 0x17, 0xdc, 0x64, 0x2e, 0x99,
	0x44, 0xf2, 0x7c, 0xe6, 0x32, 0x66, 0x2e, 0x59, 0x0c, 0x9e, 0xcf, 0x5c, 0xf2, 0xe2, 0x9f, 0x11,
	0x9c, 0xfc, 0xaa, 0x4d, 0xbe, 0x72, 0x9a, 0xc9, 0x69, 0x85, 0xe7, 0x30, 0xf4, 0xac, 0x5c, 0xec,
	0x7e, 0x28, 0xa3, 0xd1, 0xa8, 0x99, 0xad, 0x4d, 0xac, 0x77, 0x28, 0xa3, 0x81, 0x08, 0x03, 0x6f,
	0x1d, 0x8b, 0x64, 0x9e, 0x2c, 0xc7, 0x32, 0x9c, 0xf1, 0x2d, 0x80, 0xaf, 0x2b, 0x72, 0x9e, 0x72,
	0xca, 0x43, 0x11, 0xc7, 0xb2, 0xa3, 0xe0, 0xd7, 0x30, 0x35, 0x96, 0xef, 0x3a, 0xcc, 0x30, 0x30,
	0xa7, 0xc6, 0xf2, 0xcd, 0x1e, 0x7b, 0x0f, 0x23, 0xcf, 0x8a, 0x6b, 0x2f, 0x46, 0xf3, 0xde, 0x72,
	0x7a, 0xf5, 0x32, 0x6d, 0xbe, 0xef, 0xa7, 0xe6, 0x23, 0xa6, 0x37, 0xc1, 0x21, 0x77, 0x00, 0xfe,
	0x04, 0x27, 0x1b, 0xf2, 0x5e, 0x15, 0x74, 0xc7, 0xdb, 0x8a, 0xc4, 0x51, 0x08, 0x78, 0xdd, 0x09,
	0xf8, 0x23, 0xba, 0x6f, 0xb7, 0x15, 0xc9, 0xc9, 0x66, 0x6f, 0xe0, 0x37, 0x30, 0xf4, 0x99, 0xad,
	0x48, 0x1c, 0x87, 0x98, 0xb3, 0xee, 0x25, 0x8d, 0x2e, 0xa3, 0x1b, 0x67, 0x70, 0x5c, 0x2a, 0x53,
	0xd4, 0xaa, 0x20, 0x31, 0x9e, 0xf7, 0x96, 0x63, 0xf9, 0x64, 0xe3, 0x47, 0x18, 0x67, 0xe4, 0x58,
	0x69, 0xc3, 0x5b, 0x01, 0x21, 0xcf, 0xab, 0x90, 0xe7, 0x77, 0xb3, 0xb6, 0xe9, 0xaa, 0x75, 0xc9,
	0x3d, 0x85, 0x29, 0x1c, 0x7b, 0x7a, 0x24, 0xa7, 0x79, 0x2b, 0x26, 0x21, 0x02, 0xf7, 0x11, 0x37,
	0x3b, 0x8f, 0x7c, 0x62, 0xf0, 0x3b, 0x38, 0xaa, 0x5d, 0x41, 0x26, 0xdb, 0x8a, 0x93, 0xce, 0xbf,
	0x11, 0xf0, 0x3f, 0xa3, 0x43, 0xb6, 0x44, 0x53, 0xeb, 0x03, 0xa9, 0xbc, 0xd4, 0x86, 0xc4, 0x69,
	0xac, 0xb5, 0xb5, 0x71, 0x0e, 0x93, 0x9c, 0x7c, 0xe6, 0x74, 0xd5, 0xbc, 0x58, 0x31, 0x0d, 0xee,
	0xae, 0xd4, 0x10, 0xda, 0x78, 0x76, 0x75, 0x16, 0x88, 0x17, 0x91, 0xe8, 0x48, 0xf8, 0x03, 0x8c,
	0x69, 0xbd, 0xa6, 0x8c, 0xf5, 0x23, 0x89, 0xb3, 0xf0, 0xc8, 0x5e, 0xa7, 0xbb, 0x91, 0xfe, 0x7c,
	0x20, 0xe4, 0x1e, 0xc4, 0x0f, 0x30, 0xb4, 0xc6, 0x13, 0x8b, 0x97, 0x07, 0x23, 0x22, 0x84, 0x97,
	0x70, 0x44, 0x7f, 0x57, 0xda, 0x91, 0x17, 0x78, 0x90, 0x6f, 0x31, 0x7c, 0x0f, 0x67, 0xca, 0x91,
	0xba, 0xeb, 0xb6, 0xf7, 0x2a, 0x14, 0xff, 0xa2, 0xd1, 0x7f, 0xe9, 0xb4, 0xb8, 0x84, 0x61, 0x65,
	0xb5, 0x61, 0x71, 0x1e, 0x52, 0x63, 0x9b, 0x7a, 0xbf, 0x00, 0x64, 0x04, 0x16, 0x7f, 0x41, 0x72,
	0xad, 0x19, 0xa7, 0xd0, 0xd7, 0x79, 0x98, 0x87, 0xb1, 0xec, 0xeb, 0x1c, 0xe7, 0x30, 0x0c, 0x21,
	0xbb, 0xe1, 0x85, 0xfd, 0xab, 0x91, 0xd1, 0x81, 0x6f, 0x60, 0xa0, 0xcd, 0xda, 0xee, 0x06, 0x76,
	0xfc, 0xf4, 0xb5, 0x64, 0x90, 0x17, 0x2b, 0x80, 0x66, 0xe6, 0x24, 0xf9, 0xba, 0x0c, 0xb3, 0xc5,
	0x96, 0x55, 0x19, 0x6e, 0x48, 0x64, 0x34, 0xf0, 0x1d, 0x0c, 0x1e, 0x34, 0x7b, 0xd1, 0x9f, 0x27,
	0xcb, 0xc9, 0xd5, 0xa4, 0x2d, 0xf2, 0x5a, 0xb3, 0x0c, 0x8e, 0xab, 0x7f, 0x7b, 0x70, 0x1a, 0x2e,
	0xf5, 0x37, 0xe4, 0x1e, 0x75, 0x46, 0xf8, 0x06, 0x92, 0x4f, 0x79, 0x8e, 0x9d, 0x7a, 0x66, 0x9d,
	0x33, 0x7e, 0x05, 0xc9, 0x6f, 0xc4, 0x38, 0x0d, 0x92, 0xa4, 0x35, 0x39, 0x32, 0x19, 0x7d, 0x86,
	0xfc, 0x08, 0xc9, 0xb5, 0xf2, 0x5f, 0x20, 0x6f, 0xd3, 0xee, 0xca, 0x6e, 0x96, 0xc8, 0xcf, 0x71,
	0x7b, 0xee, 0x3a, 0xb8, 0x84, 0x41, 0xd3, 0x0f, 0x9e, 0xb7, 0x5c, 0x77, 0xa3, 0xcc, 0xb0, 0xab,
	0xc6, 0x88, 0xfb, 0x51, 0xc8, 0xf4, 0xfd, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0xe9, 0xb4, 0xf3,
	0xd5, 0x0c, 0x06, 0x00, 0x00,
}
