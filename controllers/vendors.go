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

// CreateVendor Register a new vendor
//
//	@Summary		Register vendor
//	@Description	Add a new vendor
//	@Tags			vendors
//	@Accept			json
//	@Produce		json
//	@Param			vendor	body	model.Vendor	true	"Vendor data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/vendor [post]
func (a *Allocator) CreateVendor(c *gin.Context) {
	userObject, authed := a.GetUserId(c)
	if authed {
		var json model.Vendor
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateVendor(json, userObject.Id)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Vendor '" + json.VendorName + "' has been added to system"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// DeleteVendor Remove an organizational unit
//
//	@Summary		Delete Vendor
//	@Description	Delete a vendor
//	@Tags			vendors
//	@Accept			json
//	@Produce		json
//	@Param			vendorId	path	int	true	"Vendor Id"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/vendor/{vendorId} [delete]
func (a *Allocator) DeleteVendor(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		ouId, _ := strconv.Atoi(c.Param("ouId"))
		status, err := model.DeleteVendor(ouId)
		if err != nil {
			log.Println("ERROR: Cannot delete vendor record: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove vendor! " + string(err.Error())})
			return
		}

		if status {
			ouIdStr := strconv.Itoa(ouId)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Vendor Id " + ouIdStr + " has been removed from system"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove vendor!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetVendors Retrieve list of all vendor records
//
//	@Summary		Retrieve list of all vendors
//	@Description	Retrieve list of all vendors
//	@Tags			vendors
//	@Produce		json
//	@Success		200	{object}	model.VendorList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/vendorss [get]
func (a *Allocator) GetVendors(c *gin.Context) {
	vendorList, err := model.GetVendors()
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

// GetVendorById Retrieve a vendor by its Id
//
//	@Summary		Retrieve a vendor by its Id
//	@Description	Retrieve a vendor by its Id
//	@Tags			vendors
//	@Produce		json
//	@Param			vendorId	path int true "Vendor ID"
//	@Success		200	{object}	model.Vendor
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/vendor/byId/{ouId} [get]
func (a *Allocator) GetVendorById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("vendorId"))
	vendor, err := model.GetVendorById(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
		return
	}

	if vendor.VendorName == "" {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with vendor id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, vendor)
	}
}
