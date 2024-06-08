package controllers

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
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/allocatord/model"
)

// CreateMachineRole Register a new machine role
//
//	@Summary		Register machine role
//	@Description	Add a new machine role
//	@Tags			machine-roles
//	@Accept			json
//	@Produce		json
//	@Param			machineRole	body	model.Vendor	true	"Machine Role data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/machineRole [post]
func (a *Allocator) CreateMachineRole(c *gin.Context) {
	userObject, authed := a.GetUserId(c)
	if authed {
		var json model.MachineRole
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateMachineRole(json, userObject.Id)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Machine Role '" + json.MachineRoleName + "' has been added to system"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// DeleteMachineRole Remove a machine role
//
//	@Summary		Delete machine role
//	@Description	Delete a machine role by Id
//	@Tags			machine-roles
//	@Accept			json
//	@Produce		json
//	@Param			machineRoleId	path	int	true	"Machine Role Id"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/machineRole/{machineRoleId} [delete]
func (a *Allocator) DeleteMachineRole(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		machineRoleId, _ := strconv.Atoi(c.Param("machineRoleId"))
		status, err := model.DeleteMachineRole(machineRoleId)
		if err != nil {
			log.Println("ERROR: Cannot delete machine role record: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove machine role! " + string(err.Error())})
			return
		}

		if status {
			machineRoleIdStr := strconv.Itoa(machineRoleId)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Machine Role Id " + machineRoleIdStr + " has been removed from system"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove machine role!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetMachineRoles Retrieve list of all machine roles
//
//	@Summary		Retrieve list of all machine roles
//	@Description	Retrieve list of all machine roles
//	@Tags			machine-roles
//	@Produce		json
//	@Success		200	{object}	model.MachineRoleList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/machineRoles [get]
func (a *Allocator) GetMachineRoles(c *gin.Context) {
	vendorList, err := model.GetMachineRoles()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
		return
	}

	if vendorList == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": vendorList})
	}
}

// GetMachineRolesById Retrieve a machine role by its Id
//
//	@Summary		Retrieve a machine role by its Id
//	@Description	Retrieve a machine role by its Id
//	@Tags			machine-roles
//	@Produce		json
//	@Param			machineRoleId	path int true "Machine Role ID"
//	@Success		200	{object}	model.MachineRole
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/machineRole/byId/{ouId} [get]
func (a *Allocator) GetMachineRoleById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("machineRoleId"))
	machineRole, err := model.GetMachineRoleById(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
		return
	}

	if machineRole.MachineRoleName == "" {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with machine role id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, machineRole)
	}
}

// UpdateMachineRoleById Retrieve a machine role by its Id
//
//	@Summary		Update a machine role by its Id
//	@Description	Update a machine role by its Id
//	@Tags			machine-roles
//	@Produce		json
//	@Param			machineRoleId	path int true "Machine Role ID"
//	@Param			machineRoleData	body model.MachineRole	true	"Machine Role data"
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/machineRole/{machineRoleId} [patch]
func (a *Allocator) UpdateMachineRoleById(c *gin.Context) {
	machineRoleId := c.Param("machineRoleId")
	id, _ := strconv.Atoi(machineRoleId)
	var json model.MachineRole
	if err := c.ShouldBindJSON(&json); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := model.UpdateMachineRoleById(id, json)
	if err != nil {
		log.Println("ERROR: Cannot update machine role with Id '" + machineRoleId + "': " + string(err.Error()))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to update machine role: " + string(err.Error())})
		return
	}

	if status {
		c.IndentedJSON(http.StatusOK, "machine role with Id '"+machineRoleId+"' has been updated")
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to update machine role with Id '" + machineRoleId + "'"})
	}
}
