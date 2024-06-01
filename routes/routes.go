package routes

/*

  Copyright 2024, JAFAX, Inc.

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

func PublicRoutes(g *gin.RouterGroup, i *controllers.Allocator) {
	// Architectures
	g.GET("/architectures")    // get all architectures
	g.GET("/architecture/:id") // get architecture by Id
	// Buildings
	g.GET("/buildings")                    // get all buildings
	g.GET("/building/byId/:id")            // get building by Id
	g.GET("/building/byShortName/:abbrev") // get building by abbreviation
	// Machine Roles
	g.GET("/machineRoles")         // get all machine roles
	g.GET("/machineRole/byId/:id") // get a machine role by Id
	// Organizational Units
	g.GET("/organizationalUnits")           // get all organizational units
	g.GET("/organizationalUnit/byId/:ouId") // get organizational unit by Id
	// Roles
	g.GET("/roles")             // get all roles
	g.GET("/role/byId/:roleId") // get role by Id
	// Systems
	g.GET("/systems")                                // get all systems
	g.GET("/systems/byVendorId/:vendorid")           // get systems by vendor Id
	g.GET("/systems/byCpuCores/:coreCount")          // get systems by number of CPU Cores
	g.GET("/systems/byRAM/:memoryCount")             // get systems by amount of installed RAM
	g.GET("/systems/byMachineRoleId/:machineRoleId") // get systems by the machine's role Id
	g.GET("/systems/byOuId/:ouId")                   // get systems by organizational unit Id
	g.GET("/system/byId/:id")                        // get system by Id
	// Users
	g.GET("/users")                  // get all users
	g.GET("/users/byOuId/:ouId")     // get all users by organizational unit Id
	g.GET("/users/byRoleId/:roleId") // get all users by role Id
	g.GET("/user/byId/:id")          // get a user by Id
	// vendors
	g.GET("/vendors")         // get all vendors
	g.GET("/vendor/byId/:id") // get a vendor by Id
	// service related routes
	g.OPTIONS("/")   // API options
	g.GET("/health") // service health API
}
