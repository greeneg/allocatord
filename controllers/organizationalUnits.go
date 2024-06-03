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

// CreateOU Register a new organizational unit
//
//	@Summary		Register organizational unit
//	@Description	Add a new organizational unit
//	@Tags			orgs
//	@Accept			json
//	@Produce		json
//	@Param			orgUnit	body	model.OrgUnit	true	"Organizational Unit data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/organizationalUnit [post]
func (a *Allocator) CreateOU(c *gin.Context) {
	userObject, authed := a.GetUserId(c)
	if authed {
		var json model.OrgUnit
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateOU(json, userObject.Id)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Orgaizational Unit '" + json.OUName + "' has been added to system"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// DeleteOU Remove an organizational unit
//
//	@Summary		Delete OU
//	@Description	Delete an organizational unit
//	@Tags			orgs
//	@Accept			json
//	@Produce		json
//	@Param			ouId	path	int	true	"Organizational Unit Id"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/organizationalUnit/{ouId} [delete]
func (a *Allocator) DeleteOU(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		ouId, _ := strconv.Atoi(c.Param("ouId"))
		status, err := model.DeleteOU(ouId)
		if err != nil {
			log.Println("ERROR: Cannot delete organizational unit: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove organizational unit! " + string(err.Error())})
			return
		}

		if status {
			ouIdStr := strconv.Itoa(ouId)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Role Id " + ouIdStr + " has been removed from system"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove user!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetOUs Retrieve list of all organizational units
//
//	@Summary		Retrieve list of all organizational units
//	@Description	Retrieve list of all organizational units
//	@Tags			orgs
//	@Produce		json
//	@Success		200	{object}	model.OrgUnitList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/organizationalUnits [get]
func (a *Allocator) GetOUs(c *gin.Context) {
	ouList, err := model.GetOUs()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
		return
	}

	if ouList == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": ouList})
	}
}

// GetOUById Retrieve a role by its Id
//
//	@Summary		Retrieve an organizational unit by its Id
//	@Description	Retrieve an organizational unit by its Id
//	@Tags			orgs
//	@Produce		json
//	@Param			ouId	path int true "Organizationl Unit ID"
//	@Success		200	{object}	model.OrgUnit
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/organizationalUnit/byId/{ouId} [get]
func (a *Allocator) GetOUById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("ouId"))
	ou, err := model.GetOUById(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
		return
	}

	if ou.OUName == "" {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with organizational unit id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, ou)
	}
}
