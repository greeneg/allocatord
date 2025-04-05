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

func CreateOSFamily(osFamily OperatingSystemFamily, id int) (bool, error) {
	log.Println("INFO: Operating System Family creation requested: " + osFamily.OSFamilyName)
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

	q, err := t.Prepare("INSERT INTO OperatingSystemFamilies (OSFamilyName, CreatorId) VALUES (?, ?)")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(osFamily.OSFamilyName, id)
	if err != nil {
		log.Println("ERROR: Cannot create operating system family '" + osFamily.OSFamilyName + "': " + string(err.Error()))
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Could not commit the DB transaction!" + string(err.Error()))
		return false, err
	}

	log.Println("INFO: Operating System Family '" + osFamily.OSFamilyName + "' created")
	return true, nil
}

func DeleteOSFamily(osFamilyIdId int) (bool, error) {
	log.Println("INFO: Operating System Family deletion requested: " + strconv.Itoa(osFamilyIdId))
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

	q, err := DB.Prepare("DELETE FROM OperatingSystemFamilies WHERE Id IS ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(osFamilyIdId)
	if err != nil {
		log.Println("ERROR: Cannot delete operating system family with ID '" + strconv.Itoa(osFamilyIdId) + "': " + string(err.Error()))
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Could not commit the DB transaction!" + string(err.Error()))
		return false, err
	}

	log.Println("INFO: Operating System Family with Id '" + strconv.Itoa(osFamilyIdId) + "' has been deleted")
	return true, nil
}

func GetOSFamilies() ([]OperatingSystemFamily, error) {
	log.Println("INFO: List of Operating System Family objects requested")
	rows, err := DB.Query("SELECT * FROM OperatingSystemFamilies")
	if err != nil {
		log.Println("ERROR: Could not run the DB query!" + string(err.Error()))
		return nil, err
	}
	defer rows.Close()

	osFamilies := make([]OperatingSystemFamily, 0)
	for rows.Next() {
		osFamily := OperatingSystemFamily{}
		err = rows.Scan(
			&osFamily.Id,
			&osFamily.OSFamilyName,
			&osFamily.CreatorId,
			&osFamily.CreationDate,
		)
		if err != nil {
			log.Println("ERROR: Cannot marshal the storage volume objects!" + string(err.Error()))
			return nil, err
		}

		osFamily.CreationDate = ConvertSqliteTimestamp(osFamily.CreationDate)

		osFamilies = append(osFamilies, osFamily)
	}

	log.Println("INFO: List of all Operating System Families retrieved")
	return osFamilies, nil
}

func GetOSFamilyById(id int) (OperatingSystemFamily, error) {
	log.Println("INFO: Operating System Family by Id requested: " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM OperatingSystemFamilies WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return OperatingSystemFamily{}, err
	}
	defer rec.Close()

	osFamily := OperatingSystemFamily{}

	r, err := rec.Query(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such Operating System Family found in DB: " + string(err.Error()))
			return OperatingSystemFamily{}, nil
		}
		log.Println("ERROR: Cannot retrieve Operating System Family from DB: " + string(err.Error()))
		return OperatingSystemFamily{}, err
	}
	defer r.Close()

	r.Scan(
		&osFamily.Id,
		&osFamily.OSFamilyName,
		&osFamily.CreatorId,
		&osFamily.CreationDate,
	)

	osFamily.CreationDate = ConvertSqliteTimestamp(osFamily.CreationDate)

	log.Println("INFO: Operating System Family with Id '" + strconv.Itoa(id) + "' has been retrieved")
	return osFamily, nil
}

func GetOSFamilyByName(name string) (OperatingSystemFamily, error) {
	log.Println("INFO: Operating System Family by name requested: " + name)
	rec, err := DB.Prepare("SELECT * FROM OperatingSystemFamilies WHERE OSFamilyName = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return OperatingSystemFamily{}, err
	}
	defer rec.Close()

	osFamily := OperatingSystemFamily{}

	r, err := rec.Query(name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such Operating System Family found in DB: " + string(err.Error()))
			return OperatingSystemFamily{}, nil
		}
		log.Println("ERROR: Cannot retrieve Operating System Family from DB: " + string(err.Error()))
		return OperatingSystemFamily{}, err
	}
	defer r.Close()

	r.Scan(
		&osFamily.Id,
		&osFamily.OSFamilyName,
		&osFamily.CreatorId,
		&osFamily.CreationDate,
	)

	osFamily.CreationDate = ConvertSqliteTimestamp(osFamily.CreationDate)

	log.Println("INFO: Operating System Family with name '" + name + "' has been retrieved")
	return osFamily, nil
}
