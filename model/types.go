package model

/*

  Copyright 2024, YggdrasilSoft, LLC.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

// primary object structs
type Architecture struct {
	Id           int    `json:"Id"`
	ISEName      string `json:"iseName"`
	RegisterSize int    `json:"registerSize"`
	CreatorId    int    `json:"creatorId"`
	CreationDate string `json:"creationDate"`
}

type ArchitectureList struct {
	Data []Architecture `json:"data"`
}

type Building struct {
	Id           int    `json:"Id"`
	BuildingName string `json:"buildingName"`
	ShortName    string `json:"shortName"`
	City         string `json:"city"`
	Region       string `json:"region"`
	CreatorId    int    `json:"creatorId"`
	CreationDate string `json:"creationDate"`
}

type BuildingList struct {
	Data []Building `json:"data"`
}

type MachineRole struct {
	Id              int    `json:"Id"`
	MachineRoleName string `json:"machineRoleName"`
	Description     string `json:"description"`
	CreatorId       int    `json:"creatorId"`
	CreationDate    string `json:"creationDate"`
}

type MachineRoleList struct {
	Data []MachineRole `json:"data"`
}

type NetworkInterface struct {
	Id           int    `json:"Id"`
	DeviceModel  string `json:"deviceModel"`
	DeviceId     string `json:"deviceId"`
	MACAddress   string `json:"macAddress"`
	SystemId     int    `json:"systemId"`
	IpAddress    string `json:"ipAddress"`
	Bitmask      int    `json:"bitmask"`
	Gateway      string `json:"gateway"`
	CreatorId    int    `json:"creatorId"`
	CreationDate string `json:"creationDate"`
}

// Note that this is not stored in the DB, rather it is synthesized off the runtime data
type NetworkInterfaces struct {
	Interfaces []NetworkInterface `json:"interfaces"`
}

type OperatingSystemFamily struct {
	Id           int    `json:"Id"`
	OSFamilyName string `json:"osFamilyName"`
	CreatorId    int    `json:"creatorId"`
	CreationDate string `json:"creationDate"`
}

type OperatingSystemFamilyList struct {
	Data []OperatingSystemFamily `json:"data"`
}

type OperatingSystem struct {
	Id               int    `json:"Id"`
	OSName           string `json:"osName"`
	OSFamilyId       int    `json:"osFamilyId"`
	VendorId         int    `json:"vendorId"`
	OSImageUrl       string `json:"osImageUrl"`
	ImageUriProtocol string `json:"imageUriProtocol"`
	CreatorId        int    `json:"creatorId"`
	CreationDate     string `json:"creationDate"`
}

type OperatingSystemList struct {
	Data []OperatingSystem `json:"data"`
}

type OrgUnit struct {
	Id           int    `json:"Id"`
	OUName       string `json:"ouName"`
	Description  string `json:"description"`
	CreatorId    int    `json:"creatorId"`
	CreationDate string `json:"creationDate"`
}

type OrgUnitList struct {
	Data []OrgUnit `json:"data"`
}

type Role struct {
	Id           int    `json:"Id"`
	RoleName     string `json:"roleName"`
	Description  string `json:"description"`
	CreationDate string `json:"creationDate"`
}

type RolesList struct {
	Data []Role `json:"data"`
}

type StorageVolume struct {
	Id           int    `json:"Id"`
	VolumeName   string `json:"volumeName"`
	StorageType  string `json:"storageType"`
	DeviceModel  string `json:"deviceModel"`
	DeviceId     string `json:"deviceId"`
	MountPoint   string `json:"mountPoint"`
	VolumeSize   int    `json:"volumeSize"`
	VolumeFormat string `json:"volumeFormat"`
	VolumeLabel  string `json:"volumeLabel"`
	SystemId     int    `json:"systemId"`
	CreatorId    int    `json:"creatorId"`
	CreationDate string `json:"creationDate"`
}

// Note that this is not stored in the DB, it's synthesized from the data
type StorageVolumes struct {
	Volumes []StorageVolume `json:"volumes"`
}

type System struct {
	Id                int    `json:"Id"`
	SerialNumber      string `json:"serialNumber"`
	ModelId           int    `json:"modelId"`
	OperatingSystemId int    `json:"osId"`
	Reimage           bool   `json:"reimage"`
	HostVars          string `json:"HostVars"`
	BilledToOrgUnitId int    `json:"billedToOrgUnitId"`
	VendorId          int    `json:"vendorId"`
	ArchitectureId    int    `json:"architectureId"`
	RAM               int    `json:"ram"`
	CpuCores          int    `json:"cpuCores"`
	CreatorId         int    `json:"creatorId"`
	CreationDate      string `json:"creationDate"`
}

type SystemList struct {
	Data []System `json:"data"`
}

type PasswordChange struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ProposedUser struct {
	Id        int    `json:"Id"`
	UserName  string `json:"userName"`
	FullName  string `json:"fullName"`
	Status    string `json:"status" enum:"enabled,disabled"`
	OrgUnitId int    `json:"orgUnitId"`
	RoleId    int    `json:"roleId"`
	Password  string `json:"password"`
}

type User struct {
	Id                      int    `json:"Id"`
	UserName                string `json:"userName"`
	FullName                string `json:"fullName"`
	Status                  string `json:"status"`
	OrgUnitId               int    `json:"orgUnitId"`
	RoleId                  int    `json:"roleId"`
	PasswordHash            string `json:"passwordHash"`
	CreationDate            string `json:"creationDate"`
	LastPasswordChangedDate string `json:"lastPasswordChangedDate"`
}

type UsersList struct {
	Data []User `json:"data"`
}

type UserStatus struct {
	Status string `json:"status" enum:"enabled,disabled"`
}

type UserStatusMsg struct {
	Message    string `json:"message"`
	UserStatus string `json:"userStatus" enum:"enabled,disabled"`
}

type UserOrgUnitId struct {
	OrgUnitId int `json:"orgUnitId"`
}

type UserOrgUnitIdMsg struct {
	Message       string `json:"message"`
	UserOrgUnitId int    `json:"orgUnitId"`
}

type UserRoleId struct {
	RoleId int `json:"roleId"`
}

type UserRoleIdMsg struct {
	Message    string `json:"message"`
	UserRoleId int    `json:"roleId"`
}

type Vendor struct {
	Id           int    `json:"Id"`
	VendorName   string `json:"vendorName"`
	CreatorId    int    `json:"creatorId"`
	CreationDate string `json:"creationDate"`
}

type VendorList struct {
	Data []Vendor `json:"data"`
}

type FailureMsg struct {
	Error string `json:"error"`
}

type SuccessMsg struct {
	Message string `json:"message"`
}
