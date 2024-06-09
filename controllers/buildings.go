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

// CreateBuilding Register a new building
//
//	@Summary		Register building
//	@Description	Add a new building
//	@Tags			buildings
//	@Accept			json
//	@Produce		json
//	@Param			building	body	model.Building	true	"Building data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/building [post]
func (a *Allocator) CreateBuilding(c *gin.Context) {
	userObject, authed := a.GetUserId(c)
	if authed {
		var json model.Building
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateBuilding(json, userObject.Id)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Building '" + json.BuildingName + "' has been added to system"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// DeleteBuilding Remove a building
//
//	@Summary		Delete building
//	@Description	Delete a building by Id
//	@Tags			buildings
//	@Accept			json
//	@Produce		json
//	@Param			buildingId	path	int	true	"Building Id"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/building/{buildingId} [delete]
func (a *Allocator) DeleteBuilding(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		buildingId, _ := strconv.Atoi(c.Param("buildingId"))
		status, err := model.DeleteBuilding(buildingId)
		if err != nil {
			log.Println("ERROR: Cannot delete building record: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove building: " + string(err.Error())})
			return
		}

		if status {
			buildingIdStr := strconv.Itoa(buildingId)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Building with Id '" + buildingIdStr + "' has been removed from system"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove building!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetBuildings Retrieve list of all building objects
//
//	@Summary		Retrieve list of all building objects
//	@Description	Retrieve list of all building objects
//	@Tags			buildings
//	@Produce		json
//	@Security		BasicAuth
//	@Success		200	{object}	model.BuildingList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/buildings [get]
func (a *Allocator) GetBuildings(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		buildingList, err := model.GetBuildings()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if buildingList == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"data": buildingList})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetBuildingById Retrieve a building by its Id
//
//	@Summary		Retrieve a building by its Id
//	@Description	Retrieve a building by its Id
//	@Tags			buildings
//	@Produce		json
//	@Param			buildingId	path int true "Building ID"
//	@Security		BasicAuth
//	@Success		200	{object}	model.Building
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/building/byId/{buildingId} [get]
func (a *Allocator) GetBuildingById(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("buildingId"))
		building, err := model.GetBuildingById(id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if building.BuildingName == "" {
			strId := strconv.Itoa(id)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with building id " + strId})
		} else {
			c.IndentedJSON(http.StatusOK, building)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetBuildingByShortName Retrieve a building by its abbreviated name
//
//	@Summary		Retrieve a building by its abbreviated name
//	@Description	Retrieve a building by its abbreviated name
//	@Tags			buildings
//	@Produce		json
//	@Param			buildingShortName	path int true "Building abbreviation"
//	@Security		BasicAuth
//	@Success		200	{object}	model.Building
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/building/byShortName/{buildingShortName} [get]
func (a *Allocator) GetBuildingByShortName(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		buildingShortName := c.Param("buildingShortName")
		building, err := model.GetBuildingByShortName(buildingShortName)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if building.BuildingName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with building abbreviation " + buildingShortName})
		} else {
			c.IndentedJSON(http.StatusOK, building)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// UpdateBuildingById Update a building by its Id
//
//	@Summary		Update a building by its Id
//	@Description	Update a building by its Id
//	@Tags			buildings
//	@Produce		json
//	@Param			buildingId	path int true "Building ID"
//	@Param			buildingData	body model.Building	true	"Building data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/building/{buildingId} [patch]
func (a *Allocator) UpdateBuildingById(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		buildingId := c.Param("buildingId")
		id, _ := strconv.Atoi(buildingId)
		var json model.Building
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		status, err := model.UpdateBuildingById(id, json)
		if err != nil {
			log.Println("ERROR: Cannot update building with Id '" + buildingId + "': " + string(err.Error()))
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to update building: " + string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, "machine role with Id '"+buildingId+"' has been updated")
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to update building with Id '" + buildingId + "'"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}
