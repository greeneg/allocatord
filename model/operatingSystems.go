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

func CreateOperatingSystem(os OperatingSystem, id int) (bool, error) {
	log.Println("INFO: Operating System creation requested: " + os.OSName)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	q, err := t.Prepare("INSERT INTO OperatingSystems (OSName, OSFamilyId, OSImageUrl, ImageUriProtocol, CreatorId) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(os.OSName, os.OSFamilyId, os.OSImageUrl, os.ImageUriProtocol, id)
	if err != nil {
		log.Println("ERROR: Cannot create Operating System record for '" + os.OSName + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: Operating System record '" + os.OSName + "' created")
	return true, nil
}

func DeleteOperatingSystem(osId int) (bool, error) {
	log.Println("INFO: Operating System deletion requested: " + strconv.Itoa(osId))
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	q, err := DB.Prepare("DELETE FROM OperatingSystems WHERE Id IS ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(osId)
	if err != nil {
		log.Println("ERROR: Cannot delete Operating System with Id '" + strconv.Itoa(osId) + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: Operating System with Id '" + strconv.Itoa(osId) + "' has been deleted")
	return true, nil
}

func GetOperatingSystems() ([]OperatingSystem, error) {
	log.Println("INFO: List of Operating System objects requested")
	rows, err := DB.Query("SELECT * FROM OperatingSystems")
	if err != nil {
		log.Println("ERROR: Could not run the DB query!" + string(err.Error()))
		return nil, err
	}

	operatingSystems := make([]OperatingSystem, 0)
	for rows.Next() {
		operatingSystem := OperatingSystem{}
		err = rows.Scan(
			&operatingSystem.Id,
			&operatingSystem.OSName,
			&operatingSystem.OSFamilyId,
			&operatingSystem.VendorId,
			&operatingSystem.OSImageUrl,
			&operatingSystem.ImageUriProtocol,
			&operatingSystem.CreatorId,
			&operatingSystem.CreationDate,
		)
		if err != nil {
			log.Println("ERROR: Cannot marshal the Operating System objects!" + string(err.Error()))
			return nil, err
		}

		operatingSystem.CreationDate = ConvertSqliteTimestamp(operatingSystem.CreationDate)

		operatingSystems = append(operatingSystems, operatingSystem)
	}

	log.Println("INFO: List of all Operating System records retrieved")
	return operatingSystems, nil
}

func GetOperatingSystemById(id int) (OperatingSystem, error) {
	log.Println("INFO: Operating System by Id requested: " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM OperatingSystems WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return OperatingSystem{}, err
	}

	os := OperatingSystem{}
	err = rec.QueryRow(id).Scan(
		&os.Id,
		&os.OSName,
		&os.OSFamilyId,
		&os.VendorId,
		&os.OSImageUrl,
		&os.ImageUriProtocol,
		&os.CreatorId,
		&os.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such Operating System record found in DB: " + string(err.Error()))
			return OperatingSystem{}, nil
		}
		log.Println("ERROR: Cannot retrieve Operating System record from DB: " + string(err.Error()))
		return OperatingSystem{}, err
	}

	os.CreationDate = ConvertSqliteTimestamp(os.CreationDate)

	return os, nil
}

func GetOperatingSystemsByFamilyId(osFamilyId int) ([]OperatingSystem, error) {
	log.Println("INFO: Operating Systems by Name requested: " + strconv.Itoa(osFamilyId))
	rec, err := DB.Prepare("SELECT * FROM OperatingSystems WHERE OSFamilyId = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return nil, err
	}

	rows, err := rec.Query(osFamilyId)
	if err != nil {
		log.Println("ERROR: Could not query DB: " + string(err.Error()))
		return nil, err
	}

	operatingSystems := make([]OperatingSystem, 0)
	for rows.Next() {
		os := OperatingSystem{}
		err = rec.QueryRow(osFamilyId).Scan(
			&os.Id,
			&os.OSName,
			&os.OSFamilyId,
			&os.VendorId,
			&os.OSImageUrl,
			&os.ImageUriProtocol,
			&os.CreatorId,
			&os.CreationDate,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("ERROR: No such Operating System with family Id '" + strconv.Itoa(osFamilyId) + "' found in DB: " + string(err.Error()))
				return nil, err
			}
			log.Println("ERROR: Cannot retrieve Operating System with family Id '" + strconv.Itoa(osFamilyId) + "' from DB: " + string(err.Error()))
			return nil, err
		}

		os.CreationDate = ConvertSqliteTimestamp(os.CreationDate)

		operatingSystems = append(operatingSystems, os)
	}

	return operatingSystems, nil
}

func GetOperatingSystemsByVendorId(osVendorId int) ([]OperatingSystem, error) {
	log.Println("INFO: Operating Systems by Vendor Id requested: " + strconv.Itoa(osVendorId))
	rec, err := DB.Prepare("SELECT * FROM OperatingSystems WHERE VendorId = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return nil, err
	}

	rows, err := rec.Query(osVendorId)
	if err != nil {
		log.Println("ERROR: Could not query DB: " + string(err.Error()))
		return nil, err
	}

	operatingSystems := make([]OperatingSystem, 0)
	for rows.Next() {
		os := OperatingSystem{}
		err = rec.QueryRow(osVendorId).Scan(
			&os.Id,
			&os.OSName,
			&os.OSFamilyId,
			&os.VendorId,
			&os.OSImageUrl,
			&os.ImageUriProtocol,
			&os.CreatorId,
			&os.CreationDate,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("ERROR: No such Operating System with family Id '" + strconv.Itoa(osVendorId) + "' found in DB: " + string(err.Error()))
				return nil, err
			}
			log.Println("ERROR: Cannot retrieve Operating System with family Id '" + strconv.Itoa(osVendorId) + "' from DB: " + string(err.Error()))
			return nil, err
		}

		os.CreationDate = ConvertSqliteTimestamp(os.CreationDate)

		operatingSystems = append(operatingSystems, os)
	}

	return operatingSystems, nil
}

func UpdateOperatingSystemById(osId int, os OperatingSystem) (bool, error) {
	log.Println("INFO: Update Operating System by Id requested: " + strconv.Itoa(osId))
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := t.Prepare("UPDATE OperatingSystems SET OSName = ?, OSFamilyId = ?, OSImageUrl = ?, ImageUriProtocol = ? WHERE Id = ?")
	if err != nil {
		return false, err
	}

	operatingSystem, err := json.Marshal(os)
	if err != nil {
		return false, err
	}
	_, err = q.Exec(operatingSystem, osId)
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}
