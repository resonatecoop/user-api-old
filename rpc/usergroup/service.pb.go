// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpc/usergroup/service.proto

package usergroup

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import user "user-api/rpc/user"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type UserGroup struct {
	Id                   string              `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	DisplayName          string              `protobuf:"bytes,2,opt,name=display_name,json=displayName" json:"display_name,omitempty"`
	Description          string              `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
	ShortBio             string              `protobuf:"bytes,4,opt,name=short_bio,json=shortBio" json:"short_bio,omitempty"`
	Avatar               []byte              `protobuf:"bytes,5,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Banner               []byte              `protobuf:"bytes,6,opt,name=banner,proto3" json:"banner,omitempty"`
	OwnerId              string              `protobuf:"bytes,7,opt,name=owner_id,json=ownerId" json:"owner_id,omitempty"`
	Type                 *GroupTaxonomy      `protobuf:"bytes,8,opt,name=type" json:"type,omitempty"`
	Followers            []*user.User        `protobuf:"bytes,9,rep,name=followers" json:"followers,omitempty"`
	Members              []*user.User        `protobuf:"bytes,10,rep,name=members" json:"members,omitempty"`
	SubGroups            []*UserGroup        `protobuf:"bytes,11,rep,name=sub_groups,json=subGroups" json:"sub_groups,omitempty"`
	Links                []*Link             `protobuf:"bytes,12,rep,name=links" json:"links,omitempty"`
	Tags                 []*Tag              `protobuf:"bytes,13,rep,name=tags" json:"tags,omitempty"`
	Address              *user.StreetAddress `protobuf:"bytes,14,opt,name=address" json:"address,omitempty"`
	Privacy              *Privacy            `protobuf:"bytes,16,opt,name=privacy" json:"privacy,omitempty"`
	RecommendedArtists   []*UserGroup        `protobuf:"bytes,17,rep,name=recommended_artists,json=recommendedArtists" json:"recommended_artists,omitempty"`
	Labels               []*UserGroup        `protobuf:"bytes,18,rep,name=labels" json:"labels,omitempty"`
	HighlightedTracks    []string            `protobuf:"bytes,19,rep,name=highlighted_tracks,json=highlightedTracks" json:"highlighted_tracks,omitempty"`
	FeaturedTrackGroup   string              `protobuf:"bytes,20,opt,name=featured_track_group,json=featuredTrackGroup" json:"featured_track_group,omitempty"`
	GroupEmailAddress    string              `protobuf:"bytes,21,opt,name=group_email_address,json=groupEmailAddress" json:"group_email_address,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *UserGroup) Reset()         { *m = UserGroup{} }
func (m *UserGroup) String() string { return proto.CompactTextString(m) }
func (*UserGroup) ProtoMessage()    {}
func (*UserGroup) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_d11390e3dbee1cbf, []int{0}
}
func (m *UserGroup) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserGroup.Unmarshal(m, b)
}
func (m *UserGroup) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserGroup.Marshal(b, m, deterministic)
}
func (dst *UserGroup) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserGroup.Merge(dst, src)
}
func (m *UserGroup) XXX_Size() int {
	return xxx_messageInfo_UserGroup.Size(m)
}
func (m *UserGroup) XXX_DiscardUnknown() {
	xxx_messageInfo_UserGroup.DiscardUnknown(m)
}

var xxx_messageInfo_UserGroup proto.InternalMessageInfo

func (m *UserGroup) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UserGroup) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *UserGroup) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *UserGroup) GetShortBio() string {
	if m != nil {
		return m.ShortBio
	}
	return ""
}

func (m *UserGroup) GetAvatar() []byte {
	if m != nil {
		return m.Avatar
	}
	return nil
}

func (m *UserGroup) GetBanner() []byte {
	if m != nil {
		return m.Banner
	}
	return nil
}

func (m *UserGroup) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (m *UserGroup) GetType() *GroupTaxonomy {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *UserGroup) GetFollowers() []*user.User {
	if m != nil {
		return m.Followers
	}
	return nil
}

func (m *UserGroup) GetMembers() []*user.User {
	if m != nil {
		return m.Members
	}
	return nil
}

func (m *UserGroup) GetSubGroups() []*UserGroup {
	if m != nil {
		return m.SubGroups
	}
	return nil
}

func (m *UserGroup) GetLinks() []*Link {
	if m != nil {
		return m.Links
	}
	return nil
}

func (m *UserGroup) GetTags() []*Tag {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *UserGroup) GetAddress() *user.StreetAddress {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *UserGroup) GetPrivacy() *Privacy {
	if m != nil {
		return m.Privacy
	}
	return nil
}

func (m *UserGroup) GetRecommendedArtists() []*UserGroup {
	if m != nil {
		return m.RecommendedArtists
	}
	return nil
}

func (m *UserGroup) GetLabels() []*UserGroup {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *UserGroup) GetHighlightedTracks() []string {
	if m != nil {
		return m.HighlightedTracks
	}
	return nil
}

func (m *UserGroup) GetFeaturedTrackGroup() string {
	if m != nil {
		return m.FeaturedTrackGroup
	}
	return ""
}

func (m *UserGroup) GetGroupEmailAddress() string {
	if m != nil {
		return m.GroupEmailAddress
	}
	return ""
}

type UserGroupToUserGroups struct {
	UserGroupId          string       `protobuf:"bytes,1,opt,name=user_group_id,json=userGroupId" json:"user_group_id,omitempty"`
	UserGroups           []*UserGroup `protobuf:"bytes,2,rep,name=user_groups,json=userGroups" json:"user_groups,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *UserGroupToUserGroups) Reset()         { *m = UserGroupToUserGroups{} }
func (m *UserGroupToUserGroups) String() string { return proto.CompactTextString(m) }
func (*UserGroupToUserGroups) ProtoMessage()    {}
func (*UserGroupToUserGroups) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_d11390e3dbee1cbf, []int{1}
}
func (m *UserGroupToUserGroups) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserGroupToUserGroups.Unmarshal(m, b)
}
func (m *UserGroupToUserGroups) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserGroupToUserGroups.Marshal(b, m, deterministic)
}
func (dst *UserGroupToUserGroups) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserGroupToUserGroups.Merge(dst, src)
}
func (m *UserGroupToUserGroups) XXX_Size() int {
	return xxx_messageInfo_UserGroupToUserGroups.Size(m)
}
func (m *UserGroupToUserGroups) XXX_DiscardUnknown() {
	xxx_messageInfo_UserGroupToUserGroups.DiscardUnknown(m)
}

var xxx_messageInfo_UserGroupToUserGroups proto.InternalMessageInfo

func (m *UserGroupToUserGroups) GetUserGroupId() string {
	if m != nil {
		return m.UserGroupId
	}
	return ""
}

func (m *UserGroupToUserGroups) GetUserGroups() []*UserGroup {
	if m != nil {
		return m.UserGroups
	}
	return nil
}

type UserGroupMember struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	DisplayName          string   `protobuf:"bytes,2,opt,name=display_name,json=displayName" json:"display_name,omitempty"`
	Avatar               []byte   `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Tags                 []*Tag   `protobuf:"bytes,4,rep,name=tags" json:"tags,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserGroupMember) Reset()         { *m = UserGroupMember{} }
func (m *UserGroupMember) String() string { return proto.CompactTextString(m) }
func (*UserGroupMember) ProtoMessage()    {}
func (*UserGroupMember) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_d11390e3dbee1cbf, []int{2}
}
func (m *UserGroupMember) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserGroupMember.Unmarshal(m, b)
}
func (m *UserGroupMember) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserGroupMember.Marshal(b, m, deterministic)
}
func (dst *UserGroupMember) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserGroupMember.Merge(dst, src)
}
func (m *UserGroupMember) XXX_Size() int {
	return xxx_messageInfo_UserGroupMember.Size(m)
}
func (m *UserGroupMember) XXX_DiscardUnknown() {
	xxx_messageInfo_UserGroupMember.DiscardUnknown(m)
}

var xxx_messageInfo_UserGroupMember proto.InternalMessageInfo

func (m *UserGroupMember) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UserGroupMember) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *UserGroupMember) GetAvatar() []byte {
	if m != nil {
		return m.Avatar
	}
	return nil
}

func (m *UserGroupMember) GetTags() []*Tag {
	if m != nil {
		return m.Tags
	}
	return nil
}

type UserGroupMembers struct {
	UserGroupId          string             `protobuf:"bytes,1,opt,name=user_group_id,json=userGroupId" json:"user_group_id,omitempty"`
	Members              []*UserGroupMember `protobuf:"bytes,2,rep,name=members" json:"members,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *UserGroupMembers) Reset()         { *m = UserGroupMembers{} }
func (m *UserGroupMembers) String() string { return proto.CompactTextString(m) }
func (*UserGroupMembers) ProtoMessage()    {}
func (*UserGroupMembers) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_d11390e3dbee1cbf, []int{3}
}
func (m *UserGroupMembers) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserGroupMembers.Unmarshal(m, b)
}
func (m *UserGroupMembers) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserGroupMembers.Marshal(b, m, deterministic)
}
func (dst *UserGroupMembers) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserGroupMembers.Merge(dst, src)
}
func (m *UserGroupMembers) XXX_Size() int {
	return xxx_messageInfo_UserGroupMembers.Size(m)
}
func (m *UserGroupMembers) XXX_DiscardUnknown() {
	xxx_messageInfo_UserGroupMembers.DiscardUnknown(m)
}

var xxx_messageInfo_UserGroupMembers proto.InternalMessageInfo

func (m *UserGroupMembers) GetUserGroupId() string {
	if m != nil {
		return m.UserGroupId
	}
	return ""
}

func (m *UserGroupMembers) GetMembers() []*UserGroupMember {
	if m != nil {
		return m.Members
	}
	return nil
}

type GroupTaxonomy struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Type                 string   `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GroupTaxonomy) Reset()         { *m = GroupTaxonomy{} }
func (m *GroupTaxonomy) String() string { return proto.CompactTextString(m) }
func (*GroupTaxonomy) ProtoMessage()    {}
func (*GroupTaxonomy) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_d11390e3dbee1cbf, []int{4}
}
func (m *GroupTaxonomy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GroupTaxonomy.Unmarshal(m, b)
}
func (m *GroupTaxonomy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GroupTaxonomy.Marshal(b, m, deterministic)
}
func (dst *GroupTaxonomy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GroupTaxonomy.Merge(dst, src)
}
func (m *GroupTaxonomy) XXX_Size() int {
	return xxx_messageInfo_GroupTaxonomy.Size(m)
}
func (m *GroupTaxonomy) XXX_DiscardUnknown() {
	xxx_messageInfo_GroupTaxonomy.DiscardUnknown(m)
}

var xxx_messageInfo_GroupTaxonomy proto.InternalMessageInfo

func (m *GroupTaxonomy) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *GroupTaxonomy) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *GroupTaxonomy) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type GroupTaxonomies struct {
	Types                []*GroupTaxonomy `protobuf:"bytes,1,rep,name=types" json:"types,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GroupTaxonomies) Reset()         { *m = GroupTaxonomies{} }
func (m *GroupTaxonomies) String() string { return proto.CompactTextString(m) }
func (*GroupTaxonomies) ProtoMessage()    {}
func (*GroupTaxonomies) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_d11390e3dbee1cbf, []int{5}
}
func (m *GroupTaxonomies) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GroupTaxonomies.Unmarshal(m, b)
}
func (m *GroupTaxonomies) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GroupTaxonomies.Marshal(b, m, deterministic)
}
func (dst *GroupTaxonomies) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GroupTaxonomies.Merge(dst, src)
}
func (m *GroupTaxonomies) XXX_Size() int {
	return xxx_messageInfo_GroupTaxonomies.Size(m)
}
func (m *GroupTaxonomies) XXX_DiscardUnknown() {
	xxx_messageInfo_GroupTaxonomies.DiscardUnknown(m)
}

var xxx_messageInfo_GroupTaxonomies proto.InternalMessageInfo

func (m *GroupTaxonomies) GetTypes() []*GroupTaxonomy {
	if m != nil {
		return m.Types
	}
	return nil
}

type Link struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Platform             string   `protobuf:"bytes,2,opt,name=platform" json:"platform,omitempty"`
	Type                 string   `protobuf:"bytes,3,opt,name=type" json:"type,omitempty"`
	Uri                  string   `protobuf:"bytes,4,opt,name=uri" json:"uri,omitempty"`
	PersonalData         bool     `protobuf:"varint,5,opt,name=personal_data,json=personalData" json:"personal_data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Link) Reset()         { *m = Link{} }
func (m *Link) String() string { return proto.CompactTextString(m) }
func (*Link) ProtoMessage()    {}
func (*Link) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_d11390e3dbee1cbf, []int{6}
}
func (m *Link) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Link.Unmarshal(m, b)
}
func (m *Link) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Link.Marshal(b, m, deterministic)
}
func (dst *Link) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Link.Merge(dst, src)
}
func (m *Link) XXX_Size() int {
	return xxx_messageInfo_Link.Size(m)
}
func (m *Link) XXX_DiscardUnknown() {
	xxx_messageInfo_Link.DiscardUnknown(m)
}

var xxx_messageInfo_Link proto.InternalMessageInfo

func (m *Link) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Link) GetPlatform() string {
	if m != nil {
		return m.Platform
	}
	return ""
}

func (m *Link) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Link) GetUri() string {
	if m != nil {
		return m.Uri
	}
	return ""
}

func (m *Link) GetPersonalData() bool {
	if m != nil {
		return m.PersonalData
	}
	return false
}

type Tag struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Type                 string   `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tag) Reset()         { *m = Tag{} }
func (m *Tag) String() string { return proto.CompactTextString(m) }
func (*Tag) ProtoMessage()    {}
func (*Tag) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_d11390e3dbee1cbf, []int{7}
}
func (m *Tag) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tag.Unmarshal(m, b)
}
func (m *Tag) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tag.Marshal(b, m, deterministic)
}
func (dst *Tag) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tag.Merge(dst, src)
}
func (m *Tag) XXX_Size() int {
	return xxx_messageInfo_Tag.Size(m)
}
func (m *Tag) XXX_DiscardUnknown() {
	xxx_messageInfo_Tag.DiscardUnknown(m)
}

var xxx_messageInfo_Tag proto.InternalMessageInfo

func (m *Tag) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Tag) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Tag) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Privacy struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Private              bool     `protobuf:"varint,2,opt,name=private" json:"private,omitempty"`
	OwnedTracks          bool     `protobuf:"varint,3,opt,name=owned_tracks,json=ownedTracks" json:"owned_tracks,omitempty"`
	SupportedArtists     bool     `protobuf:"varint,4,opt,name=supported_artists,json=supportedArtists" json:"supported_artists,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Privacy) Reset()         { *m = Privacy{} }
func (m *Privacy) String() string { return proto.CompactTextString(m) }
func (*Privacy) ProtoMessage()    {}
func (*Privacy) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_d11390e3dbee1cbf, []int{8}
}
func (m *Privacy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Privacy.Unmarshal(m, b)
}
func (m *Privacy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Privacy.Marshal(b, m, deterministic)
}
func (dst *Privacy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Privacy.Merge(dst, src)
}
func (m *Privacy) XXX_Size() int {
	return xxx_messageInfo_Privacy.Size(m)
}
func (m *Privacy) XXX_DiscardUnknown() {
	xxx_messageInfo_Privacy.DiscardUnknown(m)
}

var xxx_messageInfo_Privacy proto.InternalMessageInfo

func (m *Privacy) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Privacy) GetPrivate() bool {
	if m != nil {
		return m.Private
	}
	return false
}

func (m *Privacy) GetOwnedTracks() bool {
	if m != nil {
		return m.OwnedTracks
	}
	return false
}

func (m *Privacy) GetSupportedArtists() bool {
	if m != nil {
		return m.SupportedArtists
	}
	return false
}

type GroupedUserGroups struct {
	Labels               []*UserGroup `protobuf:"bytes,1,rep,name=labels" json:"labels,omitempty"`
	Artists              []*UserGroup `protobuf:"bytes,2,rep,name=artists" json:"artists,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *GroupedUserGroups) Reset()         { *m = GroupedUserGroups{} }
func (m *GroupedUserGroups) String() string { return proto.CompactTextString(m) }
func (*GroupedUserGroups) ProtoMessage()    {}
func (*GroupedUserGroups) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_d11390e3dbee1cbf, []int{9}
}
func (m *GroupedUserGroups) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GroupedUserGroups.Unmarshal(m, b)
}
func (m *GroupedUserGroups) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GroupedUserGroups.Marshal(b, m, deterministic)
}
func (dst *GroupedUserGroups) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GroupedUserGroups.Merge(dst, src)
}
func (m *GroupedUserGroups) XXX_Size() int {
	return xxx_messageInfo_GroupedUserGroups.Size(m)
}
func (m *GroupedUserGroups) XXX_DiscardUnknown() {
	xxx_messageInfo_GroupedUserGroups.DiscardUnknown(m)
}

var xxx_messageInfo_GroupedUserGroups proto.InternalMessageInfo

func (m *GroupedUserGroups) GetLabels() []*UserGroup {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *GroupedUserGroups) GetArtists() []*UserGroup {
	if m != nil {
		return m.Artists
	}
	return nil
}

func init() {
	proto.RegisterType((*UserGroup)(nil), "resonate.api.user.UserGroup")
	proto.RegisterType((*UserGroupToUserGroups)(nil), "resonate.api.user.UserGroupToUserGroups")
	proto.RegisterType((*UserGroupMember)(nil), "resonate.api.user.UserGroupMember")
	proto.RegisterType((*UserGroupMembers)(nil), "resonate.api.user.UserGroupMembers")
	proto.RegisterType((*GroupTaxonomy)(nil), "resonate.api.user.GroupTaxonomy")
	proto.RegisterType((*GroupTaxonomies)(nil), "resonate.api.user.GroupTaxonomies")
	proto.RegisterType((*Link)(nil), "resonate.api.user.Link")
	proto.RegisterType((*Tag)(nil), "resonate.api.user.Tag")
	proto.RegisterType((*Privacy)(nil), "resonate.api.user.Privacy")
	proto.RegisterType((*GroupedUserGroups)(nil), "resonate.api.user.GroupedUserGroups")
}

func init() {
	proto.RegisterFile("rpc/usergroup/service.proto", fileDescriptor_service_d11390e3dbee1cbf)
}

var fileDescriptor_service_d11390e3dbee1cbf = []byte{
	// 911 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0xcf, 0x6f, 0xdb, 0x36,
	0x14, 0x86, 0x62, 0x27, 0xb6, 0x9f, 0xec, 0x26, 0x66, 0xda, 0x8e, 0x73, 0x07, 0xcc, 0x53, 0x77,
	0x30, 0x36, 0xc4, 0xd9, 0xb2, 0xac, 0x87, 0x6d, 0x3d, 0xa4, 0x6b, 0x11, 0x64, 0x6b, 0x8a, 0x4d,
	0x75, 0x2f, 0xbb, 0x08, 0xb4, 0xf5, 0xe2, 0x10, 0x91, 0x44, 0x81, 0xa4, 0xd3, 0x79, 0xc0, 0x0e,
	0xc3, 0x2e, 0xdd, 0x7f, 0x3d, 0x90, 0xfa, 0x61, 0xa5, 0x89, 0x1d, 0xb7, 0xc8, 0xc5, 0x10, 0xf9,
	0xbe, 0xef, 0xe9, 0xf1, 0xf1, 0x7d, 0x9f, 0x05, 0x8f, 0x64, 0x3a, 0xd9, 0x9f, 0x29, 0x94, 0x53,
	0x29, 0x66, 0xe9, 0xbe, 0x42, 0x79, 0xc9, 0x27, 0x38, 0x4c, 0xa5, 0xd0, 0x82, 0x74, 0x25, 0x2a,
	0x91, 0x30, 0x8d, 0x43, 0x96, 0xf2, 0xa1, 0x41, 0xf5, 0x3e, 0x37, 0xbf, 0x7b, 0x2c, 0xe5, 0xfb,
	0x05, 0xf1, 0x2a, 0xc7, 0x7b, 0xd7, 0x80, 0xd6, 0x1b, 0x85, 0xf2, 0xd8, 0xe4, 0x23, 0xf7, 0x60,
	0x83, 0x87, 0xd4, 0xe9, 0x3b, 0x83, 0x96, 0xbf, 0xc1, 0x43, 0xf2, 0x05, 0xb4, 0x43, 0xae, 0xd2,
	0x88, 0xcd, 0x83, 0x84, 0xc5, 0x48, 0x37, 0x6c, 0xc4, 0xcd, 0xf7, 0x5e, 0xb1, 0x18, 0x49, 0x1f,
	0xdc, 0x10, 0xd5, 0x44, 0xf2, 0x54, 0x73, 0x91, 0xd0, 0x5a, 0x8e, 0x58, 0x6c, 0x91, 0x47, 0xd0,
	0x52, 0xe7, 0x42, 0xea, 0x60, 0xcc, 0x05, 0xad, 0xdb, 0x78, 0xd3, 0x6e, 0x3c, 0xe3, 0x82, 0x3c,
	0x84, 0x2d, 0x76, 0xc9, 0x34, 0x93, 0x74, 0xb3, 0xef, 0x0c, 0xda, 0x7e, 0xbe, 0x32, 0xfb, 0x63,
	0x96, 0x24, 0x28, 0xe9, 0x56, 0xb6, 0x9f, 0xad, 0xc8, 0xa7, 0xd0, 0x14, 0x6f, 0x13, 0x94, 0x01,
	0x0f, 0x69, 0xc3, 0xe6, 0x6a, 0xd8, 0xf5, 0x49, 0x48, 0x0e, 0xa1, 0xae, 0xe7, 0x29, 0xd2, 0x66,
	0xdf, 0x19, 0xb8, 0x07, 0xfd, 0xe1, 0xb5, 0x6e, 0x0c, 0xed, 0x21, 0x47, 0xec, 0x4f, 0x91, 0x88,
	0x78, 0xee, 0x5b, 0x34, 0xf9, 0x1e, 0x5a, 0x67, 0x22, 0x8a, 0xc4, 0x5b, 0x94, 0x8a, 0xb6, 0xfa,
	0xb5, 0x81, 0x7b, 0xf0, 0xc9, 0x0d, 0x54, 0xd3, 0x23, 0x7f, 0x81, 0x24, 0xdf, 0x42, 0x23, 0xc6,
	0x78, 0x6c, 0x48, 0xb0, 0x9a, 0x54, 0xe0, 0xc8, 0x8f, 0x00, 0x6a, 0x36, 0x0e, 0xec, 0xcd, 0x29,
	0xea, 0x5a, 0xd6, 0x67, 0x4b, 0x58, 0xb6, 0x52, 0xbf, 0xa5, 0x66, 0x63, 0xfb, 0xa4, 0xc8, 0x1e,
	0x6c, 0x46, 0x3c, 0xb9, 0x50, 0xb4, 0xbd, 0xf4, 0x6d, 0x2f, 0x79, 0x72, 0xe1, 0x67, 0x28, 0xf2,
	0x15, 0xd4, 0x35, 0x9b, 0x2a, 0xda, 0xb1, 0xe8, 0x87, 0x37, 0xa0, 0x47, 0x6c, 0xea, 0x5b, 0x0c,
	0xf9, 0x01, 0x1a, 0x2c, 0x0c, 0x25, 0x2a, 0x45, 0xef, 0x2d, 0x6d, 0xdd, 0x6b, 0x2d, 0x11, 0xf5,
	0x51, 0x86, 0xf3, 0x0b, 0x02, 0x39, 0x84, 0x46, 0x2a, 0xf9, 0x25, 0x9b, 0xcc, 0xe9, 0x8e, 0xe5,
	0xf6, 0x6e, 0xe0, 0xfe, 0x96, 0x21, 0xfc, 0x02, 0x4a, 0x4e, 0x61, 0x57, 0xe2, 0x44, 0xc4, 0x31,
	0x26, 0x21, 0x86, 0x01, 0x93, 0x9a, 0x2b, 0xad, 0x68, 0x77, 0x8d, 0x96, 0x90, 0x0a, 0xf1, 0x28,
	0xe3, 0x91, 0x43, 0xd8, 0x8a, 0xd8, 0x18, 0x23, 0x45, 0xc9, 0x1a, 0x19, 0x72, 0x2c, 0xd9, 0x03,
	0x72, 0xce, 0xa7, 0xe7, 0x11, 0x9f, 0x9e, 0x6b, 0x0c, 0x03, 0x2d, 0xd9, 0xe4, 0x42, 0xd1, 0xdd,
	0x7e, 0x6d, 0xd0, 0xf2, 0xbb, 0x95, 0xc8, 0xc8, 0x06, 0xc8, 0x37, 0x70, 0xff, 0x0c, 0x99, 0x9e,
	0xc9, 0x02, 0x9b, 0x5d, 0x24, 0xbd, 0x6f, 0x87, 0x90, 0x14, 0x31, 0x8b, 0xce, 0xc4, 0x34, 0x84,
	0x5d, 0x0b, 0x09, 0x30, 0x66, 0x3c, 0x0a, 0x8a, 0x1e, 0x3f, 0xb0, 0x84, 0xae, 0x0d, 0xbd, 0x30,
	0x91, 0xbc, 0xa9, 0xde, 0x5f, 0xf0, 0xa0, 0xac, 0x72, 0x24, 0xca, 0x47, 0x45, 0x3c, 0xe8, 0x98,
	0x33, 0x64, 0x2f, 0x0c, 0x4a, 0x81, 0xba, 0xb3, 0x02, 0x72, 0x12, 0x92, 0xa7, 0xe0, 0x2e, 0x30,
	0x8a, 0x6e, 0xac, 0xd1, 0x08, 0x28, 0xf9, 0xca, 0x7b, 0xe7, 0xc0, 0x76, 0x19, 0x39, 0xb5, 0x03,
	0xfb, 0x31, 0x66, 0xb0, 0x50, 0x73, 0xed, 0x8a, 0x9a, 0x8b, 0x71, 0xac, 0xdf, 0x3e, 0x8e, 0x9e,
	0x86, 0x9d, 0xf7, 0x2a, 0x59, 0xaf, 0x03, 0x3f, 0x2d, 0x14, 0x99, 0x9d, 0xde, 0x5b, 0x75, 0xfa,
	0x2c, 0x73, 0x29, 0x4e, 0xef, 0x18, 0x3a, 0x57, 0xdc, 0xe1, 0xda, 0xe9, 0x49, 0xee, 0x2e, 0xd9,
	0xa9, 0x33, 0xef, 0x20, 0x50, 0xb7, 0x9d, 0xc8, 0x4c, 0xcf, 0x3e, 0x7b, 0x27, 0xb0, 0x5d, 0x4d,
	0xc4, 0x51, 0x91, 0x27, 0xb0, 0x69, 0xe0, 0x8a, 0x3a, 0xb6, 0xae, 0xdb, 0x9d, 0x29, 0x83, 0x7b,
	0x7f, 0x43, 0xdd, 0x68, 0xfa, 0x5a, 0x29, 0x3d, 0x68, 0xa6, 0x11, 0xd3, 0x67, 0x42, 0xc6, 0x79,
	0x39, 0xe5, 0xba, 0x2c, 0xb3, 0x56, 0x29, 0x73, 0x07, 0x6a, 0x33, 0xc9, 0x73, 0xeb, 0x35, 0x8f,
	0xe4, 0x31, 0x74, 0x52, 0x94, 0xa6, 0x88, 0x28, 0x08, 0x99, 0x66, 0xd6, 0x7c, 0x9b, 0x7e, 0xbb,
	0xd8, 0x7c, 0xce, 0x34, 0xf3, 0x9e, 0x42, 0x6d, 0xc4, 0xa6, 0x1f, 0xdd, 0x88, 0x7f, 0x1d, 0x68,
	0xe4, 0xca, 0xbf, 0x96, 0x83, 0xe6, 0xb6, 0xa1, 0xb3, 0x34, 0x4d, 0xbf, 0x58, 0x9a, 0x21, 0x33,
	0x7e, 0x5e, 0xea, 0xb1, 0x66, 0xc3, 0xae, 0xdd, 0xcb, 0x95, 0xf8, 0x35, 0x74, 0xd5, 0x2c, 0x4d,
	0x85, 0xd4, 0x15, 0xef, 0xa8, 0x5b, 0xdc, 0x4e, 0x19, 0xc8, 0xbd, 0xc1, 0xfb, 0xc7, 0x81, 0xae,
	0x6d, 0x2e, 0x86, 0x15, 0x45, 0x2d, 0x1c, 0xc3, 0xf9, 0x00, 0xc7, 0x78, 0x02, 0x8d, 0xe2, 0x75,
	0xeb, 0xe8, 0xab, 0x00, 0x1f, 0xfc, 0xb7, 0x59, 0x19, 0xe9, 0xd7, 0xd9, 0xdf, 0x2f, 0x39, 0x85,
	0xed, 0x9f, 0x25, 0x32, 0x8d, 0x8b, 0x7f, 0xdf, 0x95, 0xe9, 0x7a, 0x2b, 0xa3, 0xe4, 0x17, 0x68,
	0x1f, 0xa3, 0xbe, 0x9b, 0x5c, 0x27, 0xb0, 0xfd, 0x26, 0x0d, 0x3f, 0xa0, 0x34, 0x7a, 0x43, 0xf4,
	0x45, 0x9c, 0xea, 0xb9, 0x49, 0xf5, 0x1c, 0x23, 0xbc, 0x8b, 0x54, 0x23, 0x20, 0xc7, 0xa8, 0x5f,
	0x9a, 0xab, 0xa8, 0xdc, 0xe4, 0x52, 0x7c, 0xef, 0xcb, 0x65, 0x32, 0xbb, 0x32, 0x09, 0xbf, 0x43,
	0xb7, 0xda, 0xb7, 0x91, 0x11, 0xde, 0x8a, 0xa4, 0xde, 0x2d, 0xda, 0x35, 0x72, 0xff, 0x15, 0xe0,
	0x28, 0x0c, 0x0b, 0xeb, 0x7a, 0x7c, 0xbb, 0x0b, 0xa9, 0x15, 0xa7, 0x7e, 0x05, 0x1d, 0x1f, 0x63,
	0x71, 0x89, 0x77, 0x93, 0xef, 0x99, 0xfb, 0x47, 0xab, 0xfc, 0x7c, 0x1c, 0x6f, 0xd9, 0x6f, 0xc0,
	0xef, 0xfe, 0x0f, 0x00, 0x00, 0xff, 0xff, 0x4b, 0xde, 0xcb, 0x74, 0x56, 0x0a, 0x00, 0x00,
}
