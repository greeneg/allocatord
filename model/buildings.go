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

func CreateBuilding(b Building, id int) (bool, error) {
	log.Println("INFO: Building creation requested: " + b.BuildingName)
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

	q, err := t.Prepare("INSERT INTO Buildings (BuildingName, ShortName, City, Region, CreatorId) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(b.BuildingName, b.ShortName, b.City, b.Region, id)
	if err != nil {
		log.Println("ERROR: Cannot create building '" + b.BuildingName + "': " + string(err.Error()))
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Could not commit the DB transaction!" + string(err.Error()))
		return false, err
	}

	log.Println("INFO: Building '" + b.BuildingName + "' created")
	return true, nil
}

func DeleteBuilding(buildingId int) (bool, error) {
	log.Println("INFO: Building deletion requested: " + strconv.Itoa(buildingId))
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

	q, err := DB.Prepare("DELETE FROM Buildings WHERE Id IS ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(buildingId)
	if err != nil {
		log.Println("ERROR: Cannot delete building with Id '" + strconv.Itoa(buildingId) + "': " + string(err.Error()))
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Could not commit the DB transaction!" + string(err.Error()))
		return false, err
	}

	log.Println("INFO: Building with Id '" + strconv.Itoa(buildingId) + "' has been deleted")
	return true, nil
}

func GetBuildings() ([]Building, error) {
	log.Println("INFO: List of building object requested")
	rows, err := DB.Query("SELECT * FROM Buildings")
	if err != nil {
		log.Println("ERROR: Could not run the DB query!" + string(err.Error()))
		return nil, err
	}
	defer rows.Close()

	buildings := make([]Building, 0)
	for rows.Next() {
		building := Building{}
		err = rows.Scan(
			&building.Id,
			&building.BuildingName,
			&building.ShortName,
			&building.City,
			&building.Region,
			&building.CreatorId,
			&building.CreationDate,
		)
		if err != nil {
			log.Println("ERROR: Cannot marshal the building objects!" + string(err.Error()))
			return nil, err
		}

		building.CreationDate = ConvertSqliteTimestamp(building.CreationDate)

		buildings = append(buildings, building)
	}

	log.Println("INFO: List of all buildings retrieved")
	return buildings, nil
}

func GetBuildingById(id int) (Building, error) {
	log.Println("INFO: Building by Id requested: " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM Buildings WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return Building{}, err
	}
	defer rec.Close()

	building := Building{}
	r, err := rec.Query(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such building found in DB: " + string(err.Error()))
			return Building{}, nil
		}
		log.Println("ERROR: Cannot retrieve building from DB: " + string(err.Error()))
		return Building{}, err
	}
	defer r.Close()

	r.Scan(
		&building.Id,
		&building.BuildingName,
		&building.ShortName,
		&building.City,
		&building.Region,
		&building.CreatorId,
		&building.CreationDate,
	)

	building.CreationDate = ConvertSqliteTimestamp(building.CreationDate)

	return building, nil
}

func GetBuildingByShortName(buildingShortName string) (Building, error) {
	log.Println("INFO: Building by Short Name requested: " + buildingShortName)
	rec, err := DB.Prepare("SELECT * FROM Buildings WHERE ShortName = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return Building{}, err
	}
	defer rec.Close()

	building := Building{}

	r, err := rec.Query(buildingShortName)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such building found in DB: " + string(err.Error()))
			return Building{}, nil
		}
		log.Println("ERROR: Cannot retrieve building from DB: " + string(err.Error()))
		return Building{}, err
	}
	defer r.Close()

	r.Scan(
		&building.Id,
		&building.BuildingName,
		&building.ShortName,
		&building.City,
		&building.Region,
		&building.CreatorId,
		&building.CreationDate,
	)

	building.CreationDate = ConvertSqliteTimestamp(building.CreationDate)

	return building, nil
}

func UpdateBuildingById(buildingId int, b Building) (bool, error) {
	log.Println("INFO: Update building by Id requested: " + strconv.Itoa(buildingId))
	t, err := DB.Begin()
	if err != nil {
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

	q, err := t.Prepare("UPDATE Buildings SET BuildingName = ?, ShortName = ?, City = ?, Region = ? WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	building, err := json.Marshal(b)
	if err != nil {
		log.Println("ERROR: Cannot marshal the building object!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(building, buildingId)
	if err != nil {
		log.Println("ERROR: Cannot update building with Id '" + strconv.Itoa(buildingId) + "': " + string(err.Error()))
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Could not commit the DB transaction!" + string(err.Error()))
		return false, err
	}

	log.Println("INFO: Building with Id '" + strconv.Itoa(buildingId) + "' has been updated")
	return true, nil
}
