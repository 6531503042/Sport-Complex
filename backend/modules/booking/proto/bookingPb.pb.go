// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v4.25.2
// source: modules/booking/proto/bookingPb.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Slot struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	StartTime string `protobuf:"bytes,2,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime   string `protobuf:"bytes,3,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
}

func (x *Slot) Reset() {
	*x = Slot{}
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Slot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Slot) ProtoMessage() {}

func (x *Slot) ProtoReflect() protoreflect.Message {
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Slot.ProtoReflect.Descriptor instead.
func (*Slot) Descriptor() ([]byte, []int) {
	return file_modules_booking_proto_bookingPb_proto_rawDescGZIP(), []int{0}
}

func (x *Slot) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Slot) GetStartTime() string {
	if x != nil {
		return x.StartTime
	}
	return ""
}

func (x *Slot) GetEndTime() string {
	if x != nil {
		return x.EndTime
	}
	return ""
}

type Booking struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId    string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	SlotId    string                 `protobuf:"bytes,3,opt,name=slot_id,json=slotId,proto3" json:"slot_id,omitempty"`
	Status    string                 `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Booking) Reset() {
	*x = Booking{}
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Booking) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Booking) ProtoMessage() {}

func (x *Booking) ProtoReflect() protoreflect.Message {
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Booking.ProtoReflect.Descriptor instead.
func (*Booking) Descriptor() ([]byte, []int) {
	return file_modules_booking_proto_bookingPb_proto_rawDescGZIP(), []int{1}
}

func (x *Booking) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Booking) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Booking) GetSlotId() string {
	if x != nil {
		return x.SlotId
	}
	return ""
}

func (x *Booking) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Booking) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Booking) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type InsertBookingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId       string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	FacilityName string `protobuf:"bytes,2,opt,name=facilityName,proto3" json:"facilityName,omitempty"`
	SlotId       string `protobuf:"bytes,3,opt,name=slot_id,json=slotId,proto3" json:"slot_id,omitempty"`
}

func (x *InsertBookingRequest) Reset() {
	*x = InsertBookingRequest{}
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InsertBookingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InsertBookingRequest) ProtoMessage() {}

func (x *InsertBookingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InsertBookingRequest.ProtoReflect.Descriptor instead.
func (*InsertBookingRequest) Descriptor() ([]byte, []int) {
	return file_modules_booking_proto_bookingPb_proto_rawDescGZIP(), []int{2}
}

func (x *InsertBookingRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *InsertBookingRequest) GetFacilityName() string {
	if x != nil {
		return x.FacilityName
	}
	return ""
}

func (x *InsertBookingRequest) GetSlotId() string {
	if x != nil {
		return x.SlotId
	}
	return ""
}

type InsertBookingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BookingId string `protobuf:"bytes,1,opt,name=bookingId,proto3" json:"bookingId,omitempty"`
	Status    string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *InsertBookingResponse) Reset() {
	*x = InsertBookingResponse{}
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InsertBookingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InsertBookingResponse) ProtoMessage() {}

func (x *InsertBookingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InsertBookingResponse.ProtoReflect.Descriptor instead.
func (*InsertBookingResponse) Descriptor() ([]byte, []int) {
	return file_modules_booking_proto_bookingPb_proto_rawDescGZIP(), []int{3}
}

func (x *InsertBookingResponse) GetBookingId() string {
	if x != nil {
		return x.BookingId
	}
	return ""
}

func (x *InsertBookingResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type UpdateBookingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BookingId    string `protobuf:"bytes,1,opt,name=booking_id,json=bookingId,proto3" json:"booking_id,omitempty"`
	FacilityName string `protobuf:"bytes,2,opt,name=facilityName,proto3" json:"facilityName,omitempty"`
}

func (x *UpdateBookingRequest) Reset() {
	*x = UpdateBookingRequest{}
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateBookingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateBookingRequest) ProtoMessage() {}

func (x *UpdateBookingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateBookingRequest.ProtoReflect.Descriptor instead.
func (*UpdateBookingRequest) Descriptor() ([]byte, []int) {
	return file_modules_booking_proto_bookingPb_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateBookingRequest) GetBookingId() string {
	if x != nil {
		return x.BookingId
	}
	return ""
}

func (x *UpdateBookingRequest) GetFacilityName() string {
	if x != nil {
		return x.FacilityName
	}
	return ""
}

type UpdateBookingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *UpdateBookingResponse) Reset() {
	*x = UpdateBookingResponse{}
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateBookingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateBookingResponse) ProtoMessage() {}

func (x *UpdateBookingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateBookingResponse.ProtoReflect.Descriptor instead.
func (*UpdateBookingResponse) Descriptor() ([]byte, []int) {
	return file_modules_booking_proto_bookingPb_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateBookingResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type FindBookingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BookingId string `protobuf:"bytes,1,opt,name=bookingId,proto3" json:"bookingId,omitempty"`
}

func (x *FindBookingRequest) Reset() {
	*x = FindBookingRequest{}
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FindBookingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindBookingRequest) ProtoMessage() {}

func (x *FindBookingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindBookingRequest.ProtoReflect.Descriptor instead.
func (*FindBookingRequest) Descriptor() ([]byte, []int) {
	return file_modules_booking_proto_bookingPb_proto_rawDescGZIP(), []int{6}
}

func (x *FindBookingRequest) GetBookingId() string {
	if x != nil {
		return x.BookingId
	}
	return ""
}

type FindBookingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BookingId string `protobuf:"bytes,1,opt,name=bookingId,proto3" json:"bookingId,omitempty"`
	Status    string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *FindBookingResponse) Reset() {
	*x = FindBookingResponse{}
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FindBookingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindBookingResponse) ProtoMessage() {}

func (x *FindBookingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_modules_booking_proto_bookingPb_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindBookingResponse.ProtoReflect.Descriptor instead.
func (*FindBookingResponse) Descriptor() ([]byte, []int) {
	return file_modules_booking_proto_bookingPb_proto_rawDescGZIP(), []int{7}
}

func (x *FindBookingResponse) GetBookingId() string {
	if x != nil {
		return x.BookingId
	}
	return ""
}

func (x *FindBookingResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_modules_booking_proto_bookingPb_proto protoreflect.FileDescriptor

var file_modules_booking_proto_bookingPb_proto_rawDesc = []byte{
	0x0a, 0x25, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x2f, 0x62, 0x6f, 0x6f, 0x6b, 0x69, 0x6e,
	0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x50,
	0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x50, 0x0a, 0x04, 0x53, 0x6c, 0x6f, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x19, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0xd9, 0x01, 0x0a, 0x07, 0x42,
	0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x17, 0x0a, 0x07, 0x73, 0x6c, 0x6f, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x73, 0x6c, 0x6f, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x6c, 0x0a, 0x14, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74,
	0x42, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x66, 0x61, 0x63, 0x69, 0x6c,
	0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x66,
	0x61, 0x63, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x73,
	0x6c, 0x6f, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6c,
	0x6f, 0x74, 0x49, 0x64, 0x22, 0x4d, 0x0a, 0x15, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x42, 0x6f,
	0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x62, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x62, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x59, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x6f,
	0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x62,
	0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x62, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x66, 0x61,
	0x63, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x66, 0x61, 0x63, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x2f,
	0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22,
	0x32, 0x0a, 0x12, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x6f, 0x6f, 0x6b, 0x69, 0x6e,
	0x67, 0x49, 0x64, 0x22, 0x4b, 0x0a, 0x13, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x6f, 0x6f, 0x6b, 0x69,
	0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x6f,
	0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62,
	0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x32, 0xca, 0x01, 0x0a, 0x0e, 0x42, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x3e, 0x0a, 0x0d, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x42, 0x6f, 0x6f,
	0x6b, 0x69, 0x6e, 0x67, 0x12, 0x15, 0x2e, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x42, 0x6f, 0x6f,
	0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x49, 0x6e,
	0x73, 0x65, 0x72, 0x74, 0x42, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x6f,
	0x6b, 0x69, 0x6e, 0x67, 0x12, 0x15, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x6f,
	0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0b, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x6f, 0x6f, 0x6b, 0x69,
	0x6e, 0x67, 0x12, 0x13, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x6f,
	0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3b, 0x5a,
	0x39, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x36, 0x35, 0x33, 0x31,
	0x35, 0x30, 0x33, 0x30, 0x34, 0x32, 0x2f, 0x53, 0x70, 0x6f, 0x72, 0x74, 0x2d, 0x43, 0x6f, 0x6d,
	0x70, 0x6c, 0x65, 0x78, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x2f, 0x62, 0x6f, 0x6f,
	0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_modules_booking_proto_bookingPb_proto_rawDescOnce sync.Once
	file_modules_booking_proto_bookingPb_proto_rawDescData = file_modules_booking_proto_bookingPb_proto_rawDesc
)

func file_modules_booking_proto_bookingPb_proto_rawDescGZIP() []byte {
	file_modules_booking_proto_bookingPb_proto_rawDescOnce.Do(func() {
		file_modules_booking_proto_bookingPb_proto_rawDescData = protoimpl.X.CompressGZIP(file_modules_booking_proto_bookingPb_proto_rawDescData)
	})
	return file_modules_booking_proto_bookingPb_proto_rawDescData
}

var file_modules_booking_proto_bookingPb_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_modules_booking_proto_bookingPb_proto_goTypes = []any{
	(*Slot)(nil),                  // 0: Slot
	(*Booking)(nil),               // 1: Booking
	(*InsertBookingRequest)(nil),  // 2: InsertBookingRequest
	(*InsertBookingResponse)(nil), // 3: InsertBookingResponse
	(*UpdateBookingRequest)(nil),  // 4: UpdateBookingRequest
	(*UpdateBookingResponse)(nil), // 5: UpdateBookingResponse
	(*FindBookingRequest)(nil),    // 6: FindBookingRequest
	(*FindBookingResponse)(nil),   // 7: FindBookingResponse
	(*timestamppb.Timestamp)(nil), // 8: google.protobuf.Timestamp
}
var file_modules_booking_proto_bookingPb_proto_depIdxs = []int32{
	8, // 0: Booking.created_at:type_name -> google.protobuf.Timestamp
	8, // 1: Booking.updated_at:type_name -> google.protobuf.Timestamp
	2, // 2: BookingService.InsertBooking:input_type -> InsertBookingRequest
	4, // 3: BookingService.UpdateBooking:input_type -> UpdateBookingRequest
	6, // 4: BookingService.FindBooking:input_type -> FindBookingRequest
	3, // 5: BookingService.InsertBooking:output_type -> InsertBookingResponse
	5, // 6: BookingService.UpdateBooking:output_type -> UpdateBookingResponse
	7, // 7: BookingService.FindBooking:output_type -> FindBookingResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_modules_booking_proto_bookingPb_proto_init() }
func file_modules_booking_proto_bookingPb_proto_init() {
	if File_modules_booking_proto_bookingPb_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_modules_booking_proto_bookingPb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_modules_booking_proto_bookingPb_proto_goTypes,
		DependencyIndexes: file_modules_booking_proto_bookingPb_proto_depIdxs,
		MessageInfos:      file_modules_booking_proto_bookingPb_proto_msgTypes,
	}.Build()
	File_modules_booking_proto_bookingPb_proto = out.File
	file_modules_booking_proto_bookingPb_proto_rawDesc = nil
	file_modules_booking_proto_bookingPb_proto_goTypes = nil
	file_modules_booking_proto_bookingPb_proto_depIdxs = nil
}
