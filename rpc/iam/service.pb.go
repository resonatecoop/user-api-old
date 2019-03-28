// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpc/iam/service.proto

/*
Package iam is a generated protocol buffer package.

It is generated from these files:
	rpc/iam/service.proto

It has these top-level messages:
	AuthReq
	AuthResp
	RefreshReq
	RefreshResp
*/
package iam

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/mwitkow/go-proto-validators"
import _ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Password authentication request
type AuthReq struct {
	// Required
	Auth string `protobuf:"bytes,1,opt,name=auth" json:"auth,omitempty"`
	// Required
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *AuthReq) Reset()                    { *m = AuthReq{} }
func (m *AuthReq) String() string            { return proto.CompactTextString(m) }
func (*AuthReq) ProtoMessage()               {}
func (*AuthReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *AuthReq) GetAuth() string {
	if m != nil {
		return m.Auth
	}
	return ""
}

func (m *AuthReq) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

// Password authentication response
type AuthResp struct {
	// Access token
	Token string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
	// Refresh token
	RefreshToken string `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken" json:"refresh_token,omitempty"`
}

func (m *AuthResp) Reset()                    { *m = AuthResp{} }
func (m *AuthResp) String() string            { return proto.CompactTextString(m) }
func (*AuthResp) ProtoMessage()               {}
func (*AuthResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *AuthResp) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *AuthResp) GetRefreshToken() string {
	if m != nil {
		return m.RefreshToken
	}
	return ""
}

// Refresh token request
type RefreshReq struct {
	// Required
	Token string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
}

func (m *RefreshReq) Reset()                    { *m = RefreshReq{} }
func (m *RefreshReq) String() string            { return proto.CompactTextString(m) }
func (*RefreshReq) ProtoMessage()               {}
func (*RefreshReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *RefreshReq) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

// Refresh token response
type RefreshResp struct {
	// Access token
	Token string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
}

func (m *RefreshResp) Reset()                    { *m = RefreshResp{} }
func (m *RefreshResp) String() string            { return proto.CompactTextString(m) }
func (*RefreshResp) ProtoMessage()               {}
func (*RefreshResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RefreshResp) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func init() {
	proto.RegisterType((*AuthReq)(nil), "twisk.iam.AuthReq")
	proto.RegisterType((*AuthResp)(nil), "twisk.iam.AuthResp")
	proto.RegisterType((*RefreshReq)(nil), "twisk.iam.RefreshReq")
	proto.RegisterType((*RefreshResp)(nil), "twisk.iam.RefreshResp")
}

func init() { proto.RegisterFile("rpc/iam/service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 383 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x52, 0x4d, 0x4f, 0xe3, 0x30,
	0x14, 0x54, 0x9a, 0x6e, 0x3f, 0xbc, 0xbb, 0x17, 0x6f, 0xbb, 0x5b, 0x45, 0x5a, 0x35, 0x4a, 0x0f,
	0xbb, 0x02, 0x12, 0x43, 0x91, 0x10, 0xea, 0x2d, 0x15, 0x1c, 0x7a, 0xe8, 0x25, 0xf4, 0x80, 0xb8,
	0x20, 0x37, 0x35, 0x89, 0x55, 0x12, 0x1b, 0xdb, 0x69, 0xc4, 0x8d, 0xdf, 0xc0, 0x1f, 0xac, 0xd4,
	0x5f, 0x82, 0x12, 0x47, 0xa1, 0x15, 0x70, 0xf2, 0x7b, 0x33, 0xa3, 0x79, 0xcf, 0x1e, 0x83, 0xbe,
	0xe0, 0x21, 0xa2, 0x38, 0x41, 0x92, 0x88, 0x0d, 0x0d, 0x89, 0xc7, 0x05, 0x53, 0x0c, 0x76, 0x55,
	0x4e, 0xe5, 0xda, 0xa3, 0x38, 0xb1, 0x2e, 0x22, 0xaa, 0xe2, 0x6c, 0xe9, 0x85, 0x2c, 0x41, 0x49,
	0x4e, 0xd5, 0x9a, 0xe5, 0x28, 0x62, 0x6e, 0xa9, 0x73, 0x37, 0xf8, 0x91, 0xae, 0xb0, 0x62, 0x42,
	0xa2, 0xba, 0xd4, 0x16, 0xd6, 0x49, 0x79, 0x84, 0x6e, 0x44, 0x52, 0x57, 0xe6, 0x38, 0x8a, 0x88,
	0x40, 0x8c, 0x2b, 0xca, 0x52, 0x89, 0x70, 0x9a, 0x32, 0x85, 0xcb, 0x5a, 0xab, 0x9d, 0x19, 0x68,
	0xfb, 0x99, 0x8a, 0x03, 0xf2, 0x04, 0x2d, 0xd0, 0xc4, 0x99, 0x8a, 0x07, 0x86, 0x6d, 0xfc, 0xef,
	0x4e, 0x5b, 0xbb, 0xed, 0xb0, 0x71, 0x6b, 0x04, 0x25, 0x06, 0x1d, 0xd0, 0xe1, 0x58, 0xca, 0x9c,
	0x89, 0xd5, 0xa0, 0x71, 0xc0, 0xd7, 0xb8, 0x73, 0x0d, 0x3a, 0xda, 0x4a, 0x72, 0xd8, 0x03, 0xdf,
	0x14, 0x5b, 0x93, 0x54, 0x9b, 0x05, 0xba, 0x81, 0x23, 0xf0, 0x53, 0x90, 0x07, 0x41, 0x64, 0x7c,
	0xaf, 0xd9, 0xd2, 0x2a, 0xf8, 0x51, 0x81, 0x8b, 0x02, 0x73, 0x8e, 0x01, 0x08, 0x74, 0x5f, 0x2c,
	0xf5, 0xf7, 0xc0, 0x68, 0xda, 0xde, 0x6d, 0x87, 0xe6, 0x8b, 0xd1, 0xab, 0x1c, 0x9d, 0x11, 0xf8,
	0x5e, 0x8b, 0xbf, 0x1a, 0x3b, 0xe6, 0xc0, 0x9c, 0xf9, 0x73, 0x88, 0x40, 0xb3, 0xd8, 0x0f, 0x42,
	0xaf, 0x7e, 0x64, 0xaf, 0xba, 0xbb, 0xf5, 0xeb, 0x03, 0x26, 0x39, 0xbc, 0x04, 0xed, 0xca, 0x1c,
	0xf6, 0xf7, 0xf8, 0xf7, 0xed, 0xac, 0xdf, 0x9f, 0xc1, 0x92, 0x4f, 0xf9, 0xab, 0x1f, 0xc3, 0x7f,
	0xc0, 0x5e, 0x14, 0xac, 0x3d, 0xf3, 0xe7, 0xf6, 0x8d, 0x8e, 0xd9, 0xb6, 0xaf, 0x58, 0x98, 0x25,
	0x24, 0xd5, 0x29, 0x8c, 0xcd, 0x33, 0xef, 0xf4, 0xc8, 0x68, 0x88, 0x09, 0xf8, 0xa3, 0xb5, 0x3a,
	0x74, 0x5b, 0x10, 0xce, 0x24, 0x55, 0x4c, 0x3c, 0xc3, 0x61, 0xac, 0x14, 0x97, 0x13, 0x84, 0xf6,
	0xfe, 0x83, 0xa0, 0x4b, 0x1a, 0x12, 0x54, 0x0e, 0xbf, 0x33, 0x29, 0x4e, 0x96, 0xad, 0x32, 0xce,
	0xf3, 0xb7, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf8, 0x5a, 0x82, 0x2a, 0x58, 0x02, 0x00, 0x00,
}