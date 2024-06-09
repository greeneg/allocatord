package model

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
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
)

func CreateNetworkInterface(n NetworkInterface, id int) (bool, error) {
	log.Println("INFO: Network Interface creation requested: " + n.DeviceModel)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	q, err := t.Prepare("INSERT INTO NetworkInterfaces (DeviceModel, DeviceId, MACAddress, SystemId, IpAddress, Bitmask, Gateway, CreatorId) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(n.DeviceModel, n.DeviceId, n.MACAddress, n.SystemId, n.IpAddress, n.Bitmask, n.Gateway, id)
	if err != nil {
		log.Println("ERROR: Cannot create network interface '" + n.DeviceModel + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: Network Interface '" + n.DeviceModel + "' created")
	return true, nil
}

func DeleteNetworkInterface(networkInterfaceId int) (bool, error) {
	log.Println("INFO: Network Interface deletion requested: " + strconv.Itoa(networkInterfaceId))
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	q, err := DB.Prepare("DELETE FROM NetworkInterfaces WHERE Id IS ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(networkInterfaceId)
	if err != nil {
		log.Println("ERROR: Cannot delete network interface '" + strconv.Itoa(networkInterfaceId) + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: Network Interface with Id '" + strconv.Itoa(networkInterfaceId) + "' has been deleted")
	return true, nil
}

func GetNetworkInterfaces() ([]NetworkInterface, error) {
	log.Println("INFO: List of network interface objects requested")
	rows, err := DB.Query("SELECT * FROM NetworkInterfaces")
	if err != nil {
		log.Println("ERROR: Could not run the DB query!" + string(err.Error()))
		return nil, err
	}

	networkInterfaces := make([]NetworkInterface, 0)
	for rows.Next() {
		networkInterface := NetworkInterface{}
		err = rows.Scan(
			&networkInterface.Id,
			&networkInterface.DeviceModel,
			&networkInterface.DeviceId,
			&networkInterface.MACAddress,
			&networkInterface.SystemId,
			&networkInterface.IpAddress,
			&networkInterface.Bitmask,
			&networkInterface.Gateway,
			&networkInterface.CreatorId,
			&networkInterface.CreationDate,
		)
		if err != nil {
			log.Println("ERROR: Cannot marshal the machine role objects!" + string(err.Error()))
			return nil, err
		}

		networkInterface.CreationDate = ConvertSqliteTimestamp(networkInterface.CreationDate)

		networkInterfaces = append(networkInterfaces, networkInterface)
	}

	log.Println("INFO: List of all machine roles retrieved")
	return networkInterfaces, nil
}

func GetNetworkInterfaceById(id int) (NetworkInterface, error) {
	log.Println("INFO: Network Interface by Id requested: " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM NetworkInterfaces WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return NetworkInterface{}, err
	}

	networkInterface := NetworkInterface{}
	err = rec.QueryRow(id).Scan(
		&networkInterface.Id,
		&networkInterface.DeviceModel,
		&networkInterface.DeviceId,
		&networkInterface.MACAddress,
		&networkInterface.SystemId,
		&networkInterface.IpAddress,
		&networkInterface.Bitmask,
		&networkInterface.Gateway,
		&networkInterface.CreatorId,
		&networkInterface.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such network interface found in DB: " + string(err.Error()))
			return NetworkInterface{}, nil
		}
		log.Println("ERROR: Cannot retrieve network interface from DB: " + string(err.Error()))
		return NetworkInterface{}, err
	}

	networkInterface.CreationDate = ConvertSqliteTimestamp(networkInterface.CreationDate)

	return networkInterface, nil
}

func GetNetworkInterfaceByIpAddress(ipAddr string) (NetworkInterface, error) {
	log.Println("INFO: Network Interface by IP address requested: " + ipAddr)
	rec, err := DB.Prepare("SELECT * FROM NetworkInterfaces WHERE IpAddress = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return NetworkInterface{}, err
	}

	networkInterface := NetworkInterface{}
	err = rec.QueryRow(ipAddr).Scan(
		&networkInterface.Id,
		&networkInterface.DeviceModel,
		&networkInterface.DeviceId,
		&networkInterface.MACAddress,
		&networkInterface.SystemId,
		&networkInterface.IpAddress,
		&networkInterface.Bitmask,
		&networkInterface.Gateway,
		&networkInterface.CreatorId,
		&networkInterface.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such network interface found in DB: " + string(err.Error()))
			return NetworkInterface{}, nil
		}
		log.Println("ERROR: Cannot retrieve network interface from DB: " + string(err.Error()))
		return NetworkInterface{}, err
	}

	networkInterface.CreationDate = ConvertSqliteTimestamp(networkInterface.CreationDate)

	return networkInterface, nil
}

func GetNetworkInterfaceByMACAddress(macAddress string) (NetworkInterface, error) {
	log.Println("INFO: Network Interface by MAC address requested: " + macAddress)
	rec, err := DB.Prepare("SELECT * FROM NetworkInterfaces WHERE MACAddress = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return NetworkInterface{}, err
	}

	networkInterface := NetworkInterface{}
	err = rec.QueryRow(macAddress).Scan(
		&networkInterface.Id,
		&networkInterface.DeviceModel,
		&networkInterface.DeviceId,
		&networkInterface.MACAddress,
		&networkInterface.SystemId,
		&networkInterface.IpAddress,
		&networkInterface.Bitmask,
		&networkInterface.Gateway,
		&networkInterface.CreatorId,
		&networkInterface.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such network interface found in DB: " + string(err.Error()))
			return NetworkInterface{}, nil
		}
		log.Println("ERROR: Cannot retrieve network interface from DB: " + string(err.Error()))
		return NetworkInterface{}, err
	}

	networkInterface.CreationDate = ConvertSqliteTimestamp(networkInterface.CreationDate)

	return networkInterface, nil
}

func GetNetworkInterfacesBySystemId(systemId int) ([]NetworkInterface, error) {
	log.Println("INFO: Network Interfaces by System Id requested: " + strconv.Itoa(systemId))
	rec, err := DB.Prepare("SELECT * FROM NetworkInterfaces WHERE SystemId = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return nil, err
	}

	rows, err := rec.Query(systemId)
	if err != nil {
		log.Println("ERROR: Could not query DB: " + string(err.Error()))
		return nil, err
	}

	networkInterfaces := make([]NetworkInterface, 0)
	for rows.Next() {
		networkInterface := NetworkInterface{}
		err = rows.Scan(
			&networkInterface.Id,
			&networkInterface.DeviceModel,
			&networkInterface.DeviceId,
			&networkInterface.MACAddress,
			&networkInterface.SystemId,
			&networkInterface.IpAddress,
			&networkInterface.Bitmask,
			&networkInterface.Gateway,
			&networkInterface.CreatorId,
			&networkInterface.CreationDate,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("ERROR: No such network interface found in DB: " + string(err.Error()))
				return nil, nil
			}
			log.Println("ERROR: Cannot retrieve network interface from DB: " + string(err.Error()))
			return nil, err
		}

		networkInterface.CreationDate = ConvertSqliteTimestamp(networkInterface.CreationDate)

		networkInterfaces = append(networkInterfaces, networkInterface)
	}

	return networkInterfaces, nil
}

func UpdateNetworkInterface(networkInterfaceId int, n NetworkInterface) (bool, error) {
	log.Println("INFO: Update network interface by Id requested: " + strconv.Itoa(networkInterfaceId))
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := t.Prepare("UPDATE NetworkInterfaces SET DeviceModel = ?, DeviceId = ?, MACAddress = ?, SystemId = ?, IpAddress = ?, Bitmask = ?, Gateway = ? WHERE Id = ?")
	if err != nil {
		return false, err
	}

	networkInterface, err := json.Marshal(n)
	if err != nil {
		return false, err
	}
	_, err = q.Exec(networkInterface, networkInterfaceId)
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}
