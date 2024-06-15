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

// CreateOSVersion Register a new operating system version
//
//	@Summary		Register operating system version
//	@Description	Add a new operating system version
//	@Tags			operating-system-versions
//	@Accept			json
//	@Produce		json
//	@Param			osVersion	body	model.OperatingSystemVersion	true	"Operating System Version data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/osVersion [post]
func (a *Allocator) CreateOSVersion(c *gin.Context) {
	userObject, authed := a.GetUserId(c)
	if authed {
		var json model.OperatingSystemVersion
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateOSVersion(json, userObject.Id)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Operating System Version '" + json.VersionNumber + "' has been added to system"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// DeleteOSVersion Remove an operating system version
//
//	@Summary		Delete operating system version
//	@Description	Delete an operating system version by Id
//	@Tags			operating-system-versions
//	@Accept			json
//	@Produce		json
//	@Param			osVersionId	path	int	true	"Operating System Version Id"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/osVersion/{osVersionId} [delete]
func (a *Allocator) DeleteOSVersion(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		osVersionId, _ := strconv.Atoi(c.Param("osVersionId"))
		status, err := model.DeleteOSVersion(osVersionId)
		if err != nil {
			log.Println("ERROR: Cannot delete Operating System Version record: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove Operating System Version! " + string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Operating System Version with Id '" + strconv.Itoa(osVersionId) + "' has been removed from system"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove Operating System Version!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetOSVersions Retrieve list of all operating systems versions
//
//	@Summary		Retrieve list of all operating systems versions
//	@Description	Retrieve list of all operating systems versions
//	@Tags			operating-system-versions
//	@Produce		json
//	@Security		BasicAuth
//	@Success		200	{object}	model.OperatingSystemVersionList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/osVersions [get]
func (a *Allocator) GetOSVersions(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		versions, err := model.GetOSVersions()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if versions == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No records found!"})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"data": versions})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetOSVersionById Retrieve an operating system version by its Id
//
//	@Summary		Retrieve an operating system version by its Id
//	@Description	Retrieve an operating system version by its Id
//	@Tags			operating-system-versions
//	@Produce		json
//	@Param			osVersionId	path int true "Operating System Version ID"
//	@Security		BasicAuth
//	@Success		200	{object}	model.OperatingSystemVersion
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/osVersion/byId/{osVersionId} [get]
func (a *Allocator) GetOSVersionById(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("osVersionId"))
		osVersion, err := model.GetOSVersionById(id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if osVersion.VersionNumber == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with Operating System Version Id " + strconv.Itoa(id)})
		} else {
			c.IndentedJSON(http.StatusOK, osVersion)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetOSVersionsByOSId Retrieve operating systems versions by operating system Id
//
//	@Summary		Retrieve operating system versions by operating system Id
//	@Description	Retrieve operating system versions by operating system Id
//	@Tags			operating-system-versions
//	@Produce		json
//	@Param			osId	path int true "Operating System ID"
//	@Security		BasicAuth
//	@Success		200	{object}	model.OperatingSystem
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/osVersions/byOSId/{osId} [get]
func (a *Allocator) GetOSVersionsByOSId(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("osId"))
		osVersions, err := model.GetOSVersionsByOSId(id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if osVersions == nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with Operating System Version associated with Operating System Id " + strconv.Itoa(id)})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"data": osVersions})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}
