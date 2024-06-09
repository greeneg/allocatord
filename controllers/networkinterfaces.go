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

// CreateNetworkInterface Register a new network interface
//
//	@Summary		Register network interface
//	@Description	Add a new network interface
//	@Tags			network-interfaces
//	@Accept			json
//	@Produce		json
//	@Param			networkInterface	body	model.NetworkInterface	true	"Network Interface data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/networkInterface [post]
func (a *Allocator) CreateNetworkInterface(c *gin.Context) {
	userObject, authed := a.GetUserId(c)
	if authed {
		var json model.NetworkInterface
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateNetworkInterface(json, userObject.Id)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Network Interface '" + strconv.Itoa(json.Id) + "' has been added to system"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// DeleteNetworkInterface Remove a network interface
//
//	@Summary		Delete network interface
//	@Description	Delete a network interface by Id
//	@Tags			network-interfaces
//	@Accept			json
//	@Produce		json
//	@Param			networkInterfaceId	path	int	true	"Network Interface Id"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/networkInterface/{networkInterfaceId} [delete]
func (a *Allocator) DeleteNetworkInterface(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		networkInterfaceId, _ := strconv.Atoi(c.Param("networkInterfaceId"))
		status, err := model.DeleteNetworkInterface(networkInterfaceId)
		if err != nil {
			log.Println("ERROR: Cannot delete network interface record: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to network interface role! " + string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Network Interface Id " + strconv.Itoa(networkInterfaceId) + " has been removed from system"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove network interface!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetNetworkInterfaces Retrieve list of all network interfaces
//
//	@Summary		Retrieve list of all network interfaces
//	@Description	Retrieve list of all network interfaces
//	@Tags			network-interfaces
//	@Produce		json
//	@Security		BasicAuth
//	@Success		200	{object}	model.NetworkInterfaces
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/networkInterfaces [get]
func (a *Allocator) GetNetworkInterfaces(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		networkInterfaces, err := model.GetNetworkInterfaces()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if networkInterfaces == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"interfaces": networkInterfaces})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetNetworkInterfaceById Retrieve a network interface by its Id
//
//	@Summary		Retrieve a network interface by its Id
//	@Description	Retrieve a network interface by its Id
//	@Tags			network-interfaces
//	@Produce		json
//	@Param			networkInterfaceId	path int true "Network Interface ID"
//	@Security		BasicAuth
//	@Success		200	{object}	model.NetworkInterface
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/networkInterface/byId/{networkInterfaceId} [get]
func (a *Allocator) GetNetworkInterfaceById(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("networkInterfaceId"))
		networkInterface, err := model.GetNetworkInterfaceById(id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if networkInterface.DeviceModel == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with network interface id " + strconv.Itoa(id)})
		} else {
			c.IndentedJSON(http.StatusOK, networkInterface)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetNetworkInterfaceByIpAddress Retrieve a network interface by its IP Address
//
//	@Summary		Retrieve a network interface by its IP Address
//	@Description	Retrieve a network interface by its IP Address
//	@Tags			network-interfaces
//	@Produce		json
//	@Param			networkInterfaceIpAddress	path int true "Network Interface IP Address"
//	@Security		BasicAuth
//	@Success		200	{object}	model.NetworkInterface
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/networkInterface/byIpAddress/{networkInterfaceIpAddress} [get]
func (a *Allocator) GetNetworkInterfaceByIpAddress(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		ipAddr := c.Param("networkInterfaceIpAddress")
		networkInterface, err := model.GetNetworkInterfaceByIpAddress(ipAddr)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if networkInterface.DeviceModel == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with IP Address " + ipAddr})
		} else {
			c.IndentedJSON(http.StatusOK, networkInterface)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetNetworkInterfaceByMACAddress Retrieve a network interface by its MAC Address
//
//	@Summary		Retrieve a network interface by its MAC Address
//	@Description	Retrieve a network interface by its MAC Address
//	@Tags			network-interfaces
//	@Produce		json
//	@Param			networkInterfaceMACAddress	path int true "Network Interface MAC Address"
//	@Security		BasicAuth
//	@Success		200	{object}	model.NetworkInterface
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/networkInterface/byMACAddress/{networkInterfaceMACAddress} [get]
func (a *Allocator) GetNetworkInterfaceByMACAddress(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		macAddress := c.Param("networkInterfaceId")
		networkInterface, err := model.GetNetworkInterfaceByMACAddress(macAddress)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if networkInterface.DeviceModel == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with MAC Address " + macAddress})
		} else {
			c.IndentedJSON(http.StatusOK, networkInterface)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetNetworkInterfacesBySystemId Retrieve the list of network interfaces for a system's Id
//
//	@Summary		Retrieve the list of network interfaces for a system's Id
//	@Description	Retrieve the list of network interfaces for a system's Id
//	@Tags			network-interfaces
//	@Produce		json
//	@Param			systemId	path int true "System ID"
//	@Security		BasicAuth
//	@Success		200	{object}	model.NetworkInterfaces
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/networkInterfaces/{systemId} [get]
func (a *Allocator) GetNetworkInterfacesBySystemId(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("systemId"))
		networkInterfaces, err := model.GetNetworkInterfacesBySystemId(id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if networkInterfaces.DeviceModel == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with system id " + strconv.Itoa(id)})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"interfaces": networkInterfaces})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// UpdateNetworkInterfaceById Update a network interface by its Id
//
//	@Summary		Update a network interface by its Id
//	@Description	Update a network interface by its Id
//	@Tags			network-interfaces
//	@Produce		json
//	@Param			networkInterfaceId	path int true "Network Interface ID"
//	@Param			networkInterfaceData	body model.MachineRole	true	"Network Interface data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/networkInterface/{networkInterfaceId} [patch]
func (a *Allocator) UpdateNetworkInterface(c *gin.Context) {
	_, authed := a.GetUserId(c)
	if authed {
		networkInterfaceId := c.Param("networkInterfaceId")
		id, _ := strconv.Atoi(networkInterfaceId)
		var json model.NetworkInterface
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		status, err := model.UpdateNetworkInterface(id, json)
		if err != nil {
			log.Println("ERROR: Cannot update network interface with Id '" + networkInterfaceId + "': " + string(err.Error()))
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to update network interface: " + string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, "Network interface with Id '"+networkInterfaceId+"' has been updated")
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to update network interface with Id '" + networkInterfaceId + "'"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}
