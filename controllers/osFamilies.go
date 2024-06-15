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

// CreateOSFamily Register an operating system family
//
//	@Summary		Register operating system family
//	@Description	Add a new operating system family
//	@Tags			operating-system-families
//	@Accept			json
//	@Produce		json
//	@Param			osFamily	body	model.OperatingSystemFamily	true	"Operating System Family data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/osFamily [post]
func (a *Allocator) CreateOSFamily(c *gin.Context) {
	userObject, authed := a.GetUserId(c)
	if authed {
		var json model.OperatingSystemFamily
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateOSFamily(json, userObject.Id)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Operating System Family with ID '" + strconv.Itoa(json.Id) + "' has been added to system"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// DeleteOSFamily Remove an operating system family
//
//	@Summary		Delete operating system family
//	@Description	Delete an operating system family
//	@Tags			operating-system-families
//	@Accept			json
//	@Produce		json
//	@Param			osFamilyId	path	int	true	"Operating System Family Id"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/osFamily/{storageVolumeId} [delete]
func (a *Allocator) DeleteOSFamily(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		osFamilyId, _ := strconv.Atoi(c.Param("osFamilyId"))
		status, err := model.DeleteOSFamily(osFamilyId)
		if err != nil {
			log.Println("ERROR: Cannot delete operating system family: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove operating system family! " + string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Operating System Family with ID " + strconv.Itoa(osFamilyId) + " has been removed from system"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove operating system family!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetOSFamilies Retrieve list of all operating system families
//
//	@Summary		Retrieve list of all storage volumes
//	@Description	Retrieve list of all storage volumes
//	@Tags			operating-system-families
//	@Produce		json
//	@Security		BasicAuth
//	@Success		200	{object}	model.OperatingSystemFamilyList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/osFamilies [get]
func (a *Allocator) GetOSFamilies(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		osFamilies, err := model.GetOSFamilies()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if osFamilies == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No records found!"})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"data": osFamilies})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetOSFamilyById Retrieve an operating system family by its Id
//
//	@Summary		Retrieve an operating system family by its Id
//	@Description	Retrieve an operating system family by its Id
//	@Tags			operating-system-families
//	@Produce		json
//	@Param			osFamilyId	path int true "Operating System Family ID"
//	@Security		BasicAuth
//	@Success		200	{object}	model.OperatingSystemFamily
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/osFamily/byId/{osFamilyId} [get]
func (a *Allocator) GetOSFamilyById(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("osFamilyId"))
		osFamily, err := model.GetOSFamilyById(id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if osFamily.OSFamilyName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with operating system family ID " + strconv.Itoa(id)})
		} else {
			c.IndentedJSON(http.StatusOK, osFamily)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetOSFamilyByName Retrieve an operating system family by its name
//
//	@Summary		Retrieve an operating system family by its name
//	@Description	Retrieve an operating system family by its name
//	@Tags			operating-system-families
//	@Produce		json
//	@Param			osFamilyName	path	string	true	"Operating System Family Name"
//	@Security		BasicAuth
//	@Success		200	{object}	model.OperatingSystemFamily
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/osFamily/byName/{osFamilyName} [get]
func (a *Allocator) GetOSFamilyByName(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		osFamilyName := c.Param("osFamilyName")
		osFamily, err := model.GetOSFamilyByName(osFamilyName)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if osFamily.OSFamilyName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with storage volume having label " + osFamilyName})
		} else {
			c.IndentedJSON(http.StatusOK, osFamily)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}
