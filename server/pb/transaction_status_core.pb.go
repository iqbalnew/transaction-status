// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.22.2
// source: transaction_status_core.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Direction int32

const (
	Direction_UNKNOWN_DIRECTION Direction = 0
	Direction_DESC              Direction = 1
	Direction_ASC               Direction = 2
)

// Enum value maps for Direction.
var (
	Direction_name = map[int32]string{
		0: "UNKNOWN_DIRECTION",
		1: "DESC",
		2: "ASC",
	}
	Direction_value = map[string]int32{
		"UNKNOWN_DIRECTION": 0,
		"DESC":              1,
		"ASC":               2,
	}
)

func (x Direction) Enum() *Direction {
	p := new(Direction)
	*p = x
	return p
}

func (x Direction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Direction) Descriptor() protoreflect.EnumDescriptor {
	return file_transaction_status_core_proto_enumTypes[0].Descriptor()
}

func (Direction) Type() protoreflect.EnumType {
	return &file_transaction_status_core_proto_enumTypes[0]
}

func (x Direction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Direction.Descriptor instead.
func (Direction) EnumDescriptor() ([]byte, []int) {
	return file_transaction_status_core_proto_rawDescGZIP(), []int{0}
}

type TaskStep int32

const (
	TaskStep_UNKNOWN_TASK_STEP TaskStep = 0
	TaskStep_MAKER             TaskStep = 1
	TaskStep_CHECKER           TaskStep = 2
	TaskStep_SIGNER            TaskStep = 3
	TaskStep_RELEASER          TaskStep = 4
	TaskStep_COMPLETE          TaskStep = 5
)

// Enum value maps for TaskStep.
var (
	TaskStep_name = map[int32]string{
		0: "UNKNOWN_TASK_STEP",
		1: "MAKER",
		2: "CHECKER",
		3: "SIGNER",
		4: "RELEASER",
		5: "COMPLETE",
	}
	TaskStep_value = map[string]int32{
		"UNKNOWN_TASK_STEP": 0,
		"MAKER":             1,
		"CHECKER":           2,
		"SIGNER":            3,
		"RELEASER":          4,
		"COMPLETE":          5,
	}
)

func (x TaskStep) Enum() *TaskStep {
	p := new(TaskStep)
	*p = x
	return p
}

func (x TaskStep) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TaskStep) Descriptor() protoreflect.EnumDescriptor {
	return file_transaction_status_core_proto_enumTypes[1].Descriptor()
}

func (TaskStep) Type() protoreflect.EnumType {
	return &file_transaction_status_core_proto_enumTypes[1]
}

func (x TaskStep) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TaskStep.Descriptor instead.
func (TaskStep) EnumDescriptor() ([]byte, []int) {
	return file_transaction_status_core_proto_rawDescGZIP(), []int{1}
}

type Status int32

const (
	Status_UNKNOWN_STATUS             Status = 0
	Status_PENDING_MAKER_CONFIRMATION Status = 1
	Status_ON_CHECKER                 Status = 2
	Status_ON_SIGNER                  Status = 3
	Status_ON_RELEASER                Status = 4
	Status_SUCCESS                    Status = 5
	Status_FAILED                     Status = 6
	Status_RETUR                      Status = 7
	Status_WAITING_SCHEDULE           Status = 8
	Status_REJECTED                   Status = 9
	Status_WAITING_PROCESS            Status = 10
	Status_ON_PROGRESS                Status = 11
	Status_PROCEED_NEXT_PROCESS       Status = 12
	Status_APPROVED                   Status = 13
	Status_ACTIVE                     Status = 14
	Status_INACTIVE                   Status = 15
	Status_EXPIRED                    Status = 16
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0:  "UNKNOWN_STATUS",
		1:  "PENDING_MAKER_CONFIRMATION",
		2:  "ON_CHECKER",
		3:  "ON_SIGNER",
		4:  "ON_RELEASER",
		5:  "SUCCESS",
		6:  "FAILED",
		7:  "RETUR",
		8:  "WAITING_SCHEDULE",
		9:  "REJECTED",
		10: "WAITING_PROCESS",
		11: "ON_PROGRESS",
		12: "PROCEED_NEXT_PROCESS",
		13: "APPROVED",
		14: "ACTIVE",
		15: "INACTIVE",
		16: "EXPIRED",
	}
	Status_value = map[string]int32{
		"UNKNOWN_STATUS":             0,
		"PENDING_MAKER_CONFIRMATION": 1,
		"ON_CHECKER":                 2,
		"ON_SIGNER":                  3,
		"ON_RELEASER":                4,
		"SUCCESS":                    5,
		"FAILED":                     6,
		"RETUR":                      7,
		"WAITING_SCHEDULE":           8,
		"REJECTED":                   9,
		"WAITING_PROCESS":            10,
		"ON_PROGRESS":                11,
		"PROCEED_NEXT_PROCESS":       12,
		"APPROVED":                   13,
		"ACTIVE":                     14,
		"INACTIVE":                   15,
		"EXPIRED":                    16,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_transaction_status_core_proto_enumTypes[2].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_transaction_status_core_proto_enumTypes[2]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_transaction_status_core_proto_rawDescGZIP(), []int{2}
}

var File_transaction_status_core_proto protoreflect.FileDescriptor

var file_transaction_status_core_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x1d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2a, 0x35,
	0x0a, 0x09, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x15, 0x0a, 0x11, 0x55,
	0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x44, 0x49, 0x52, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e,
	0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x45, 0x53, 0x43, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03,
	0x41, 0x53, 0x43, 0x10, 0x02, 0x2a, 0x61, 0x0a, 0x08, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x65,
	0x70, 0x12, 0x15, 0x0a, 0x11, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x54, 0x41, 0x53,
	0x4b, 0x5f, 0x53, 0x54, 0x45, 0x50, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x41, 0x4b, 0x45,
	0x52, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x48, 0x45, 0x43, 0x4b, 0x45, 0x52, 0x10, 0x02,
	0x12, 0x0a, 0x0a, 0x06, 0x53, 0x49, 0x47, 0x4e, 0x45, 0x52, 0x10, 0x03, 0x12, 0x0c, 0x0a, 0x08,
	0x52, 0x45, 0x4c, 0x45, 0x41, 0x53, 0x45, 0x52, 0x10, 0x04, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x4f,
	0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x05, 0x2a, 0xa9, 0x02, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x0e, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x53,
	0x54, 0x41, 0x54, 0x55, 0x53, 0x10, 0x00, 0x12, 0x1e, 0x0a, 0x1a, 0x50, 0x45, 0x4e, 0x44, 0x49,
	0x4e, 0x47, 0x5f, 0x4d, 0x41, 0x4b, 0x45, 0x52, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x52, 0x4d,
	0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x4f, 0x4e, 0x5f, 0x43, 0x48,
	0x45, 0x43, 0x4b, 0x45, 0x52, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x4f, 0x4e, 0x5f, 0x53, 0x49,
	0x47, 0x4e, 0x45, 0x52, 0x10, 0x03, 0x12, 0x0f, 0x0a, 0x0b, 0x4f, 0x4e, 0x5f, 0x52, 0x45, 0x4c,
	0x45, 0x41, 0x53, 0x45, 0x52, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45,
	0x53, 0x53, 0x10, 0x05, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x06,
	0x12, 0x09, 0x0a, 0x05, 0x52, 0x45, 0x54, 0x55, 0x52, 0x10, 0x07, 0x12, 0x14, 0x0a, 0x10, 0x57,
	0x41, 0x49, 0x54, 0x49, 0x4e, 0x47, 0x5f, 0x53, 0x43, 0x48, 0x45, 0x44, 0x55, 0x4c, 0x45, 0x10,
	0x08, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x4a, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x09, 0x12,
	0x13, 0x0a, 0x0f, 0x57, 0x41, 0x49, 0x54, 0x49, 0x4e, 0x47, 0x5f, 0x50, 0x52, 0x4f, 0x43, 0x45,
	0x53, 0x53, 0x10, 0x0a, 0x12, 0x0f, 0x0a, 0x0b, 0x4f, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x47, 0x52,
	0x45, 0x53, 0x53, 0x10, 0x0b, 0x12, 0x18, 0x0a, 0x14, 0x50, 0x52, 0x4f, 0x43, 0x45, 0x45, 0x44,
	0x5f, 0x4e, 0x45, 0x58, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x43, 0x45, 0x53, 0x53, 0x10, 0x0c, 0x12,
	0x0c, 0x0a, 0x08, 0x41, 0x50, 0x50, 0x52, 0x4f, 0x56, 0x45, 0x44, 0x10, 0x0d, 0x12, 0x0a, 0x0a,
	0x06, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x0e, 0x12, 0x0c, 0x0a, 0x08, 0x49, 0x4e, 0x41,
	0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x0f, 0x12, 0x0b, 0x0a, 0x07, 0x45, 0x58, 0x50, 0x49, 0x52,
	0x45, 0x44, 0x10, 0x10, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_transaction_status_core_proto_rawDescOnce sync.Once
	file_transaction_status_core_proto_rawDescData = file_transaction_status_core_proto_rawDesc
)

func file_transaction_status_core_proto_rawDescGZIP() []byte {
	file_transaction_status_core_proto_rawDescOnce.Do(func() {
		file_transaction_status_core_proto_rawDescData = protoimpl.X.CompressGZIP(file_transaction_status_core_proto_rawDescData)
	})
	return file_transaction_status_core_proto_rawDescData
}

var file_transaction_status_core_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_transaction_status_core_proto_goTypes = []interface{}{
	(Direction)(0), // 0: transaction_status.service.v1.Direction
	(TaskStep)(0),  // 1: transaction_status.service.v1.TaskStep
	(Status)(0),    // 2: transaction_status.service.v1.Status
}
var file_transaction_status_core_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_transaction_status_core_proto_init() }
func file_transaction_status_core_proto_init() {
	if File_transaction_status_core_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_transaction_status_core_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_transaction_status_core_proto_goTypes,
		DependencyIndexes: file_transaction_status_core_proto_depIdxs,
		EnumInfos:         file_transaction_status_core_proto_enumTypes,
	}.Build()
	File_transaction_status_core_proto = out.File
	file_transaction_status_core_proto_rawDesc = nil
	file_transaction_status_core_proto_goTypes = nil
	file_transaction_status_core_proto_depIdxs = nil
}