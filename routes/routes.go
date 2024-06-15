package routes

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

import (
	"github.com/gin-gonic/gin"

	"github.com/greeneg/allocatord/controllers"
)

func PrivateRoutes(g *gin.RouterGroup, a *controllers.Allocator) {
	// Architectures
	g.GET("/architectures", a.GetArchitectures)                              // get all architectures
	g.GET("/architecture/byId/:architectureId", a.GetArchitectureById)       // get architecture by Id
	g.GET("/architecture/byName/:architectureName", a.GetArchitectureByName) // get architectures by name
	g.POST("/architecture", a.CreateArchitecture)                            // create a new architecture record
	g.DELETE("/architecture/:architectureId", a.DeleteArchitecture)          // delete an architecture by Id
	// Buildings
	g.GET("/buildings", a.GetBuildings)                              // get all buildings
	g.GET("/building/byId/:id", a.GetBuildingById)                   // get building by Id
	g.GET("/building/byShortName/:abbrev", a.GetBuildingByShortName) // get building by abbreviation
	g.POST("/building", a.CreateBuilding)                            // create a new building
	g.PATCH("/building/:buildingId", a.UpdateBuildingById)           // update a building by its Id
	g.DELETE("/building/:buildingId", a.DeleteBuilding)              // delete a building by its Id
	// Machine Roles
	g.GET("/machineRoles", a.GetMachineRoles)                       // get all machine roles
	g.GET("/machineRole/byId/:id", a.GetMachineRoleById)            // get a machine role by Id
	g.POST("/machineRole", a.CreateMachineRole)                     // create a new machine role
	g.PATCH("/machineRole/:machineRoleId", a.UpdateMachineRoleById) // update a machine role by Id
	g.DELETE("/machineRole/:machineRoleId", a.DeleteMachineRole)    // delete a machine role by Id
	// Network Interfaces
	g.GET("/networkInterfaces", a.GetNetworkInterfaces)                                    // get all network interfaces
	g.GET("/networkInterfaces/:systemId", a.GetNetworkInterfacesBySystemId)                // get all network interfaces by system
	g.GET("/networkInterface/byId/:networkInterfaceId", a.GetNetworkInterfaceById)         // get network interface by Id
	g.GET("/networkInterface/byIpAddress/:ipAddress", a.GetNetworkInterfaceByIpAddress)    // get network interface by IP address
	g.GET("/networkInterface/byMACAddress/:macAddress", a.GetNetworkInterfaceByMACAddress) // get network interface by MAC address
	g.POST("/networkInterface", a.CreateNetworkInterface)                                  // create a new network interface
	g.PATCH("/networkInterface/:networkInterfaceId", a.UpdateNetworkInterface)             // update a network interface
	g.DELETE("/networkInterface/:networkInterfaceId", a.DeleteNetworkInterface)            // delete a network interface
	// Operating System Families
	g.GET("/osFamilies", a.GetOSFamilies)                        // get all operating system families
	g.GET("/osFamily/byId/:osFamilyId", a.GetOSFamilyById)       // get operating system family by Id
	g.GET("/osFamily/byName/:osFamilyName", a.GetOSFamilyByName) // get operating system family by name
	g.POST("/osFamily", a.CreateOSFamily)                        // create a new operating system family
	g.DELETE("/osFamily/:osFamilyId", a.DeleteOSFamily)          // delete an operating system by Id
	// Operating Systems
	g.GET("/operatingSystems", a.GetOperatingSystems)                                  // get all operating systems
	g.GET("/operatingSystems/byFamilyId/:osFamilyId", a.GetOperatingSystemsByFamilyId) // get operating sytems by OS Family Id
	g.GET("/operatingSystem/byId/:osId", a.GetOperatingSystemById)                     // get operating systems by Id
	g.POST("/operatingSystem", a.CreateOperatingSystem)                                // create operating system
	g.PATCH("/operatingSystem/:osId", a.UpdateOperatingSystemById)                     // update an operating system by Id
	g.DELETE("/operatingSystem/:osId", a.DeleteOperatingSystem)                        // delete an operating system
	// Operating System Versions
	g.GET("/osVersions")
	g.GET("/osVersion/byId/:osVersionId")
	g.GET("/osVersion/byOSId/:osId")
	g.POST("/osVersion")
	g.DELETE("/osVersion/:osVersionId")
	// Organizational Units
	g.GET("/organizationalUnits", a.GetOUs)              // get all organizational units
	g.GET("/organizationalUnit/byId/:ouId", a.GetOUById) // get organizational unit by Id
	g.POST("/organizationalUnit", a.CreateOU)            // create a new organizational unit
	g.DELETE("/organizationalUnit/:ouId", a.DeleteOU)    // delete an organizational unit by Id
	// Roles
	g.GET("/roles", a.GetRoles)                      // get all roles
	g.GET("/role/byId/:roleId", a.GetRoleById)       // get role by Id
	g.GET("/role/byName/:roleName", a.GetRoleByName) // get role by name
	g.POST("/role", a.CreateRole)                    // create new role
	g.DELETE("/role/:roleId", a.DeleteRole)          // delete a role by Id
	// Storage Volumes
	g.GET("/storageVolumes", a.GetStorageVolumes)                                            // get all storage volumes
	g.GET("/storageVolumes/:systemId", a.GetStorageVolumesBySystemId)                        // get storage volumes by system Id
	g.GET("/storageVolume/byId/:storageVolumeId", a.GetStorageVolumeById)                    // get a storage volume by Id
	g.GET("/storageVolume/:systemId/byLabel/:storageVolumeLabel", a.GetStorageVolumeByLabel) // get a storage volume by label
	g.POST("/storageVolume", a.CreateStorageVolume)                                          // create a new storage volume
	g.PATCH("/storageVolume/:storageVolumeId", a.UpdateStorageVolume)                        // update a storage volume
	g.DELETE("/storageVolume/:storageVolumeId", a.DeleteStorageVolume)                       // delete a storage volume
	// user related routes
	g.GET("/users", a.GetUsers)                          // get all users
	g.GET("/users/byOuId/:ouId", a.GetUsersByOuId)       // get all users by organizational unit Id
	g.GET("/users/byRoleId/:roleId", a.GetUsersByRoleId) // get all users by role Id
	g.GET("/user/:name/status", a.GetUserStatus)         // get whether a user is locked or not
	g.POST("/user", a.CreateUser)                        // create new user
	g.PATCH("/user/:name/status", a.SetUserStatus)       // lock a user
	g.PATCH("/user/:name/ouId", a.SetUserOuId)           // set a user's organizational unit Id
	g.PATCH("/user/:name/roleId", a.SetUserRoleId)       // set a user's role Id
	g.DELETE("/user/:name", a.DeleteUser)                // trash a user
	// Vendors
	g.GET("/vendors", a.GetVendors)               // get all vendors
	g.GET("/vendor/byId/:id", a.GetVendorById)    // get a vendor by Id
	g.POST("/vendor", a.CreateVendor)             // create new vendor record
	g.DELETE("/vendor/:vendorId", a.DeleteVendor) // delete a vendor record by Id
}

func PublicRoutes(g *gin.RouterGroup, a *controllers.Allocator) {
	// User related routes
	g.GET("/user/byId/:id", a.GetUserById)          // get a user by Id
	g.GET("/user/:name", a.GetUserByUserName)       // get a user by username
	g.PATCH("/user/:name", a.ChangeAccountPassword) // update a user password
	// Systems
	g.GET("/systems")                                // get all systems
	g.GET("/systems/byVendorId/:vendorid")           // get systems by vendor Id
	g.GET("/systems/byCpuCores/:coreCount")          // get systems by number of CPU Cores
	g.GET("/systems/byRAM/:memoryCount")             // get systems by amount of installed RAM
	g.GET("/systems/byMachineRoleId/:machineRoleId") // get systems by the machine's role Id
	g.GET("/systems/byOuId/:ouId")                   // get systems by organizational unit Id
	g.GET("/system/byId/:id")                        // get system by Id
	// service related routes
	g.OPTIONS("/")   // API options
	g.GET("/health") // service health API
}
