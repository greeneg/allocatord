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
	g.GET("/architectures")    // get all architectures
	g.GET("/architecture/:id") // get architecture by Id
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
	g.GET("/networkInterfaces", a.GetNetworkInterfaces)
	g.GET("/networkInterfaces/:systemId", a.GetNetworkInterfacesBySystemId)
	g.GET("/networkInterface/byId/:networkInterfaceId", a.GetNetworkInterfaceById)
	g.GET("/networkInterface/byIpAddress/:ipAddress", a.GetNetworkInterfaceByIpAddress)
	g.GET("/networkInterface/byMACAddress/:macAddress", a.GetNetworkInterfaceByMACAddress)
	g.POST("/networkInterface", a.CreateNetworkInterface)
	g.PATCH("/networkInterface/:networkInterfaceId", a.UpdateNetworkInterface)
	g.DELETE("/networkInterface/:networkInterfaceId", a.DeleteNetworkInterface)
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
	g.GET("/storageVolumes")
	g.GET("/storageVolumes/:systemId")
	g.GET("/storageVolume/byId/:storageVolumeId")
	g.GET("/storageVolume/byLabel/:storageVolumeLabel")
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
