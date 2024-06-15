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

// CreateOperatingSystem Register a new operating system
//
//	@Summary		Register operating system
//	@Description	Add a new operating system
//	@Tags			operating-systems
//	@Accept			json
//	@Produce		json
//	@Param			operatingSystem	body	model.OperatingSystem	true	"Operating System data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/operatingSystem [post]
func (a *Allocator) CreateOperatingSystem(c *gin.Context) {
	userObject, authed := a.GetUserId(c)
	if authed {
		var json model.OperatingSystem
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateOperatingSystem(json, userObject.Id)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Operating System '" + json.OSName + "' has been added to system"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// DeleteOperatingSystem Remove an operating system
//
//	@Summary		Delete operating system
//	@Description	Delete an operating system by Id
//	@Tags			operating-systems
//	@Accept			json
//	@Produce		json
//	@Param			osId	path	int	true	"Operating System Id"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/operatingSystem/{osId} [delete]
func (a *Allocator) DeleteOperatingSystem(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		osId, _ := strconv.Atoi(c.Param("osId"))
		status, err := model.DeleteOperatingSystem(osId)
		if err != nil {
			log.Println("ERROR: Cannot delete Operating System record: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove Operating System! " + string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Operating System with Id '" + strconv.Itoa(osId) + "' has been removed from system"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove Operating System!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetOperatingSystems Retrieve list of all operating systems
//
//	@Summary		Retrieve list of all operating systems
//	@Description	Retrieve list of all operating systems
//	@Tags			operating-systems
//	@Produce		json
//	@Security		BasicAuth
//	@Success		200	{object}	model.OperatingSystemList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/operatingSystems [get]
func (a *Allocator) GetOperatingSystems(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		operatingSystems, err := model.GetOperatingSystems()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if operatingSystems == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No records found!"})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"data": operatingSystems})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetOperatingSystemById Retrieve an operating system by its Id
//
//	@Summary		Retrieve an operating system by its Id
//	@Description	Retrieve an operating system by its Id
//	@Tags			operating-systems
//	@Produce		json
//	@Param			osId	path int true "Operating System ID"
//	@Security		BasicAuth
//	@Success		200	{object}	model.OperatingSystem
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/operatingSystem/byId/{osId} [get]
func (a *Allocator) GetOperatingSystemById(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("osId"))
		operatingSystem, err := model.GetOperatingSystemById(id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if operatingSystem.OSName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with Operating System id " + strconv.Itoa(id)})
		} else {
			c.IndentedJSON(http.StatusOK, operatingSystem)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetOperatingSystemByFamilyId Retrieve an operating system by its family Id
//
//	@Summary		Retrieve an operating system by its family Id
//	@Description	Retrieve an operating system by its family Id
//	@Tags			operating-systems
//	@Produce		json
//	@Param			osFamilyId	path int true "Operating System Family ID"
//	@Security		BasicAuth
//	@Success		200	{object}	model.OperatingSystem
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/operatingSystem/byFamilyId/{osFamilyId} [get]
func (a *Allocator) GetOperatingSystemsByFamilyId(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("osFamilyId"))
		operatingSystem, err := model.GetOperatingSystemsByFamilyId(id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if operatingSystem == nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with Operating System name " + strconv.Itoa(id)})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"data": operatingSystem})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// UpdateOperatingSystemById Update an operating system by its Id
//
//	@Summary		Update an operating system by its Id
//	@Description	Update an operating system by its Id
//	@Tags			operating-systems
//	@Produce		json
//	@Param			operatingSystemId	path int true "Operating System ID"
//	@Param			operatingSystemData	body model.OperatingSystem	true	"Operating System data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/operatingSystem/{osId} [patch]
func (a *Allocator) UpdateOperatingSystemById(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		osId, _ := strconv.Atoi(c.Param("machineRoleId"))
		var json model.OperatingSystem
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		osIdStr := strconv.Itoa(osId)
		status, err := model.UpdateOperatingSystemById(osId, json)
		if err != nil {
			log.Println("ERROR: Cannot update Operating System with Id '" + osIdStr + "': " + string(err.Error()))
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to update Operating System record: " + string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Operating System with Id '" + osIdStr + "' has been updated"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to update Operating System record with Id '" + osIdStr + "'"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}
