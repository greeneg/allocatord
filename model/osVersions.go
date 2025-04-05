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
	"log"
	"strconv"
)

func CreateOSVersion(osVersion OperatingSystemVersion, id int) (bool, error) {
	log.Println("INFO: Operating System Version record creation requested: " + osVersion.VersionNumber)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			t.Rollback()
			log.Println("ERROR: Transaction rolled back due to panic: " + string(r.(error).Error()))
		}
		if err != nil {
			t.Rollback()
			log.Println("ERROR: Transaction rolled back due to error: " + string(err.Error()))
		}
	}()

	q, err := t.Prepare("INSERT INTO OperatingSystemVersions (OperatingSystemId, VersionNumber, CreatorId) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(osVersion.OperatingSystemId, osVersion.VersionNumber, id)
	if err != nil {
		log.Println("ERROR: Cannot create Operating System Version record for version number '" + osVersion.VersionNumber + "': " + string(err.Error()))
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Could not commit the DB transaction!" + string(err.Error()))
		return false, err
	}

	log.Println("INFO: Operating System Version record '" + osVersion.VersionNumber + "' created")
	return true, nil
}

func DeleteOSVersion(osVersionId int) (bool, error) {
	log.Println("INFO: Operating System Version deletion requested: " + strconv.Itoa(osVersionId))
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			t.Rollback()
			log.Println("ERROR: Transaction rolled back due to panic: " + string(r.(error).Error()))
		}
		if err != nil {
			t.Rollback()
			log.Println("ERROR: Transaction rolled back due to error: " + string(err.Error()))
		}
	}()

	q, err := DB.Prepare("DELETE FROM OperatingSystemVersions WHERE Id IS ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(osVersionId)
	if err != nil {
		log.Println("ERROR: Cannot delete Operating System Version with Id '" + strconv.Itoa(osVersionId) + "': " + string(err.Error()))
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Could not commit the DB transaction!" + string(err.Error()))
		return false, err
	}

	log.Println("INFO: Operating System Version with Id '" + strconv.Itoa(osVersionId) + "' has been deleted")
	return true, nil
}

func GetOSVersions() ([]OperatingSystemVersion, error) {
	log.Println("INFO: List of Operating System objects requested")
	rows, err := DB.Query("SELECT * FROM OperatingSystems")
	if err != nil {
		log.Println("ERROR: Could not run the DB query!" + string(err.Error()))
		return nil, err
	}
	defer rows.Close()

	osVersions := make([]OperatingSystemVersion, 0)
	for rows.Next() {
		osVersion := OperatingSystemVersion{}
		err = rows.Scan(
			&osVersion.Id,
			&osVersion.OperatingSystemId,
			&osVersion.VersionNumber,
			&osVersion.CreatorId,
			&osVersion.CreationDate,
		)
		if err != nil {
			log.Println("ERROR: Cannot marshal the Operating System Version objects!" + string(err.Error()))
			return nil, err
		}

		osVersion.CreationDate = ConvertSqliteTimestamp(osVersion.CreationDate)

		osVersions = append(osVersions, osVersion)
	}

	log.Println("INFO: List of all Operating System Version records retrieved")
	return osVersions, nil
}

func GetOSVersionById(id int) (OperatingSystemVersion, error) {
	log.Println("INFO: Operating System Version by Id requested: " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM OperatingSystemVersions WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return OperatingSystemVersion{}, err
	}
	defer rec.Close()

	osVersion := OperatingSystemVersion{}

	r, err := rec.Query(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such Operating System Version record found in DB: " + string(err.Error()))
			return OperatingSystemVersion{}, nil
		}
		log.Println("ERROR: Cannot retrieve Operating System Version record from DB: " + string(err.Error()))
		return OperatingSystemVersion{}, err
	}
	defer r.Close()

	r.Scan(
		&osVersion.Id,
		&osVersion.OperatingSystemId,
		&osVersion.VersionNumber,
		&osVersion.CreatorId,
		&osVersion.CreationDate,
	)

	osVersion.CreationDate = ConvertSqliteTimestamp(osVersion.CreationDate)

	log.Println("INFO: Operating System Version with Id '" + strconv.Itoa(id) + "' has been retrieved")
	return osVersion, nil
}

func GetOSVersionsByOSId(osId int) ([]OperatingSystemVersion, error) {
	log.Println("INFO: Operating System Versions by OS Id requested: " + strconv.Itoa(osId))
	rec, err := DB.Prepare("SELECT * FROM OperatingSystemVersions WHERE OperatingSystemId = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return nil, err
	}
	defer rec.Close()

	rows, err := rec.Query(osId)
	if err != nil {
		log.Println("ERROR: Could not query DB: " + string(err.Error()))
		return nil, err
	}
	defer rows.Close()

	versions := make([]OperatingSystemVersion, 0)
	for rows.Next() {
		osVersion := OperatingSystemVersion{}
		err = rows.Scan(
			&osVersion.Id,
			&osVersion.OperatingSystemId,
			&osVersion.VersionNumber,
			&osVersion.CreatorId,
			&osVersion.CreationDate,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("ERROR: No such Operating System Version with OS Id '" + strconv.Itoa(osId) + "' found in DB: " + string(err.Error()))
				return nil, err
			}
			log.Println("ERROR: Cannot retrieve Operating System with family Id '" + strconv.Itoa(osId) + "' from DB: " + string(err.Error()))
			return nil, err
		}

		osVersion.CreationDate = ConvertSqliteTimestamp(osVersion.CreationDate)

		versions = append(versions, osVersion)
	}

	log.Println("INFO: List of Operating System Versions with OS Id '" + strconv.Itoa(osId) + "' retrieved")
	return versions, nil
}
