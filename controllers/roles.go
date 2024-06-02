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
	"github.com/greeneg/allocatord/helpers"
	"github.com/greeneg/allocatord/model"
)

// CreateUser Register a role for user rights assignment
//
//	@Summary		Register role
//	@Description	Add a new role
//	@Tags			role
//	@Accept			json
//	@Produce		json
//	@Param			role	body	model.Role	true	"Role data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/role [post]
func (a *Allocator) CreateRole(c *gin.Context) {
	var json model.Role
	if err := c.ShouldBindJSON(&json); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s, err := model.CreateRole(json)
	if s {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Role '" + json.RoleName + "' has been added to system"})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// DeleteRole Remove a role
//
//	@Summary		Delete role
//	@Description	Delete a role
//	@Tags			role
//	@Accept			json
//	@Produce		json
//	@Param			roleId	path	int	true	"Role Id"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/role/{roleId} [delete]
func (a *Allocator) DeleteRole(c *gin.Context) {
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	status, err := model.DeleteRole(roleId)
	if err != nil {
		log.Println("ERROR: Cannot delete role: " + string(err.Error()))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove role! " + string(err.Error())})
		return
	}

	if status {
		roleIdStr := strconv.Itoa(roleId)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Role Id " + roleIdStr + " has been removed from system"})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove user!"})
	}
}

// GetRoles Retrieve list of all roles
//
//	@Summary		Retrieve list of all roles
//	@Description	Retrieve list of all roles
//	@Tags			role
//	@Produce		json
//	@Success		200	{object}	model.RolesList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/roles [get]
func (a *Allocator) GetRoles(c *gin.Context) {
	roles, err := model.GetRoles()
	helpers.FatalCheckError(err)

	if roles == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": roles})
	}
}

// GetRoleById Retrieve a role by its Id
//
//	@Summary		Retrieve a role by its Id
//	@Description	Retrieve a role by its Id
//	@Tags			role
//	@Produce		json
//	@Param			roleId	path int true "Role ID"
//	@Success		200	{object}	model.Role
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/id/{id} [get]
func (a *Allocator) GetRoleById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("roleId"))
	role, err := model.GetRoleById(id)
	helpers.FatalCheckError(err)

	if err != nil {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with role id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, role)
	}
}
