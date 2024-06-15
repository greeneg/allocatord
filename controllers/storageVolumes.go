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

// CreateStorageVolume Register a storage volume
//
//	@Summary		Register storage volume
//	@Description	Add a new storage volume
//	@Tags			storage-volumes
//	@Accept			json
//	@Produce		json
//	@Param			storageVolume	body	model.Role	true	"Storage Volume data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/storageVolume [post]
func (a *Allocator) CreateStorageVolume(c *gin.Context) {
	userObject, authed := a.GetUserId(c)
	if authed {
		var json model.StorageVolume
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateStorageVolume(json, userObject.Id)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Storage Volume with ID '" + strconv.Itoa(json.Id) + "' has been added to system"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// DeleteStorageVolume Remove a storage volume
//
//	@Summary		Delete storage volume
//	@Description	Delete a storage volume
//	@Tags			storage-volumes
//	@Accept			json
//	@Produce		json
//	@Param			storageVolumeId	path	int	true	"Storage Volume Id"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/storageVolume/{storageVolumeId} [delete]
func (a *Allocator) DeleteStorageVolume(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		storageVolumeId, _ := strconv.Atoi(c.Param("storageVolumeId"))
		status, err := model.DeleteStorageVolume(storageVolumeId)
		if err != nil {
			log.Println("ERROR: Cannot delete storage volume: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove storage volume! " + string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Storage Volume with ID " + strconv.Itoa(storageVolumeId) + " has been removed from system"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove storage volume!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetStorageVolumes Retrieve list of all storage volumes
//
//	@Summary		Retrieve list of all storage volumes
//	@Description	Retrieve list of all storage volumes
//	@Tags			storage-volumes
//	@Produce		json
//	@Security		BasicAuth
//	@Success		200	{object}	model.StorageVolumes
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/storageVolumes [get]
func (a *Allocator) GetStorageVolumes(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		volumes, err := model.GetStorageVolumes()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if volumes == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No records found!"})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"data": volumes})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetStorageVolumesById Retrieve a storage volume by its Id
//
//	@Summary		Retrieve a storage volume by its Id
//	@Description	Retrieve a storage volume by its Id
//	@Tags			storage-volumes
//	@Produce		json
//	@Param			storageVolumeId	path int true "Storage Volume ID"
//	@Security		BasicAuth
//	@Success		200	{object}	model.StorageVolume
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/storageVolume/byId/{storageVolumeId} [get]
func (a *Allocator) GetStorageVolumeById(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("storageVolumeId"))
		volume, err := model.GetStorageVolumeById(id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if volume.VolumeName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with storage volume ID " + strconv.Itoa(id)})
		} else {
			c.IndentedJSON(http.StatusOK, volume)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetStorageVolumesByLabel Retrieve a storage volume by its volume label
//
//	@Summary		Retrieve a storage volume by its volume label
//	@Description	Retrieve a storage volume by its volume label
//	@Tags			storage-volumes
//	@Produce		json
//	@Param			systemId	path	string	true	"System ID"
//	@Param			storageVolumeLabel	path	string	true	"Storage Volume Label"
//	@Security		BasicAuth
//	@Success		200	{object}	model.StorageVolume
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/storageVolume/{systemId}/byLabel/{storageVolumeLabel} [get]
func (a *Allocator) GetStorageVolumeByLabel(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		systemId, _ := strconv.Atoi(c.Param("systemId"))
		label := c.Param("storageVolumeLabel")
		volume, err := model.GetStorageVolumeByLabel(label, systemId)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if volume.VolumeName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with storage volume having label " + label})
		} else {
			c.IndentedJSON(http.StatusOK, volume)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetStorageVolumesBySystemId Retrieve storage volumes by system Id
//
//	@Summary		Retrieve storage volumes by system Id
//	@Description	Retrieve storage volumes by system Id
//	@Tags			storage-volumes
//	@Produce		json
//	@Param			systemId	path int true "System ID"
//	@Security		BasicAuth
//	@Success		200	{object}	model.StorageVolumes
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/storageVolumes/{systemId} [get]
func (a *Allocator) GetStorageVolumesBySystemId(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		systemId, _ := strconv.Atoi(c.Param("systemId"))
		volumes, err := model.GetStorageVolumesBySystemId(systemId)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if volumes == nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with storage volumes for system ID " + strconv.Itoa(systemId)})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"volumes": volumes})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// UpdateStorageVolume Update a storage volume by its Id
//
//	@Summary		Update a storage volume by its Id
//	@Description	Update a storage volume by its Id
//	@Tags			storage-volumes
//	@Produce		json
//	@Param			storageVolumeId	path int true "Storage Volume ID"
//	@Param			storageVolumeData	body model.StorageVolume	true	"Storage Volume data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/storageVolume/{storageVolumeId} [patch]
func (a *Allocator) UpdateStorageVolume(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		storageVolumeId := c.Param("storageVolumeId")
		id, _ := strconv.Atoi(storageVolumeId)
		var json model.StorageVolume
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		status, err := model.UpdateStorageVolume(id, json)
		if err != nil {
			log.Println("ERROR: Cannot update storage volume with Id '" + storageVolumeId + "': " + string(err.Error()))
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to update storage volume: " + string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Machine role with Id '" + strconv.Itoa(id) + "' has been updated"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to update machine role with Id '" + strconv.Itoa(id) + "'"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}
