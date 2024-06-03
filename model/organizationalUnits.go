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

func CreateOU(o OrgUnit, id int) (bool, error) {
	log.Println("INFO: Organizational Unit creation requested: " + o.OUName)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	q, err := t.Prepare("INSERT INTO Roles (OUName, Description, CreatorId) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(o.OUName, o.Description, id)
	if err != nil {
		log.Println("ERROR: Cannot create organizational unit '" + o.OUName + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: Organizational Unit '" + o.OUName + "' created")
	return true, nil
}

func DeleteOU(ouId int) (bool, error) {
	log.Println("INFO: Organizational Unit deletion requested: " + strconv.Itoa(ouId))
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	q, err := DB.Prepare("DELETE FROM OrganizationalUnits WHERE Id IS ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(ouId)
	if err != nil {
		log.Println("ERROR: Cannot delete organizational unit Id '" + strconv.Itoa(ouId) + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: Organizational Unit with Id '" + strconv.Itoa(ouId) + "' has been deleted")
	return true, nil
}

func GetOUs() ([]OrgUnit, error) {
	log.Println("INFO: List of organizational unit object requested")
	rows, err := DB.Query("SELECT * FROM OrganizationalUnits")
	if err != nil {
		log.Println("ERROR: Could not run the DB query!" + string(err.Error()))
		return nil, err
	}

	units := make([]OrgUnit, 0)
	for rows.Next() {
		unit := OrgUnit{}
		err = rows.Scan(
			&unit.Id,
			&unit.OUName,
			&unit.Description,
			&unit.CreatorId,
			&unit.CreationDate,
		)
		if err != nil {
			log.Println("ERROR: Cannot marshal the organizational unit objects!" + string(err.Error()))
			return nil, err
		}

		unit.CreationDate = ConvertSqliteTimestamp(unit.CreationDate)

		units = append(units, unit)
	}

	log.Println("INFO: List of all organizational units retrieved")
	return units, nil
}

func GetOUById(id int) (OrgUnit, error) {
	log.Println("INFO: Organizational Unit by Id requested: " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM OrganizationalUnits WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return OrgUnit{}, err
	}

	ou := OrgUnit{}
	err = rec.QueryRow(id).Scan(
		&ou.Id,
		&ou.OUName,
		&ou.Description,
		&ou.CreatorId,
		&ou.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such user found in DB: " + string(err.Error()))
			return OrgUnit{}, nil
		}
		log.Println("ERROR: Cannot retrieve user from DB: " + string(err.Error()))
		return OrgUnit{}, err
	}

	ou.CreationDate = ConvertSqliteTimestamp(ou.CreationDate)

	return ou, nil
}
