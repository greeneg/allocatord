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

// CreateArchitecture Register a new architecture
//
//	@Summary		Register architecture
//	@Description	Add a new architecture
//	@Tags			architectures
//	@Accept			json
//	@Produce		json
//	@Param			architecture	body	model.Architecture	true	"Architecture data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/architecture [post]
func (a *Allocator) CreateArchitecture(c *gin.Context) {
	userObject, authed := a.GetUserId(c)
	if authed {
		var json model.Architecture
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateArchitecture(json, userObject.Id)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Machine Role '" + json.ISEName + "' has been added to system"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// DeleteArchitecture Remove an architecture
//
//	@Summary		Delete architecture
//	@Description	Delete an architecture by Id
//	@Tags			architectures
//	@Accept			json
//	@Produce		json
//	@Param			architectureId	path	int	true	"Architecture Id"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/architecture/{architectureId} [delete]
func (a *Allocator) DeleteArchitecture(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		architectureId, _ := strconv.Atoi(c.Param("architectureId"))
		status, err := model.DeleteArchitecture(architectureId)
		if err != nil {
			log.Println("ERROR: Cannot delete architecture record: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove architecture! " + string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Architecture Id " + strconv.Itoa(architectureId) + " has been removed from system"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove architecture!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetArchitectures Retrieve list of all architectures
//
//	@Summary		Retrieve list of all architectures
//	@Description	Retrieve list of all architectures
//	@Tags			architectures
//	@Produce		json
//	@Security		BasicAuth
//	@Success		200	{object}	model.ArchitectureList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/architectures [get]
func (a *Allocator) GetArchitectures(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		architectures, err := model.GetArchitectures()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if architectures == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No records found!"})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"data": architectures})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetArchitectureById Retrieve an architecture by its Id
//
//	@Summary		Retrieve an architecture by its Id
//	@Description	Retrieve an architecture by its Id
//	@Tags			architectures
//	@Produce		json
//	@Param			architectureId	path int true "Architecture ID"
//	@Security		BasicAuth
//	@Success		200	{object}	model.Architecture
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/architecture/byId/{architectureId} [get]
func (a *Allocator) GetArchitectureById(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("architectureId"))
		architecture, err := model.GetArchitectureById(id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if architecture.ISEName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with architecture Id " + strconv.Itoa(id)})
		} else {
			c.IndentedJSON(http.StatusOK, architecture)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetArchitectureByName Retrieve an architecture by its name
//
//	@Summary		Retrieve an architecture by its name
//	@Description	Retrieve an architecture by its name
//	@Tags			architectures
//	@Produce		json
//	@Param			architectureName	path int true "Architecture Name"
//	@Security		BasicAuth
//	@Success		200	{object}	model.Architecture
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/architecture/byName/{architectureName} [get]
func (a *Allocator) GetArchitectureByName(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		name := c.Param("architectureName")
		architecture, err := model.GetArchitectureByName(name)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if architecture.ISEName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with architecture name " + name})
		} else {
			c.IndentedJSON(http.StatusOK, architecture)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}
