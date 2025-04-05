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

func CreateArchitecture(a Architecture, id int) (bool, error) {
	log.Println("INFO: Architecture creation requested: " + a.ISEName)
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

	q, err := t.Prepare("INSERT INTO Architectures (ISEName, RegisterSize, CreatorId) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(a.ISEName, a.RegisterSize, id)
	if err != nil {
		log.Println("ERROR: Cannot create architecture '" + a.ISEName + "': " + string(err.Error()))
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Could not commit the DB transaction!" + string(err.Error()))
		return false, err
	}

	log.Println("INFO: Architecture '" + a.ISEName + "' created")
	return true, nil
}

func DeleteArchitecture(architectureId int) (bool, error) {
	log.Println("INFO: Architecture deletion requested: " + strconv.Itoa(architectureId))
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

	q, err := DB.Prepare("DELETE FROM Architectures WHERE Id IS ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(architectureId)
	if err != nil {
		log.Println("ERROR: Cannot delete architecture '" + strconv.Itoa(architectureId) + "': " + string(err.Error()))
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Could not commit the DB transaction!" + string(err.Error()))
		return false, err
	}

	log.Println("INFO: Architecture with Id '" + strconv.Itoa(architectureId) + "' has been deleted")
	return true, nil
}

func GetArchitectures() ([]Architecture, error) {
	log.Println("INFO: List of architecture objects requested")
	rows, err := DB.Query("SELECT * FROM Architectures")
	if err != nil {
		log.Println("ERROR: Could not run the DB query!" + string(err.Error()))
		return nil, err
	}
	defer rows.Close()

	architectures := make([]Architecture, 0)
	for rows.Next() {
		architecture := Architecture{}
		err = rows.Scan(
			&architecture.Id,
			&architecture.ISEName,
			&architecture.RegisterSize,
			&architecture.CreatorId,
			&architecture.CreationDate,
		)
		if err != nil {
			log.Println("ERROR: Cannot marshal the architecture objects!" + string(err.Error()))
			return nil, err
		}

		architecture.CreationDate = ConvertSqliteTimestamp(architecture.CreationDate)

		architectures = append(architectures, architecture)
	}

	log.Println("INFO: List of all architectures retrieved")
	return architectures, nil
}

func GetArchitectureById(id int) (Architecture, error) {
	log.Println("INFO: Architecture by Id requested: " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM Architectures WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return Architecture{}, err
	}
	defer rec.Close()

	architecture := Architecture{}

	r, err := rec.Query(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such architecture found in DB: " + string(err.Error()))
			return Architecture{}, nil
		}
		log.Println("ERROR: Cannot retrieve architecture from DB: " + string(err.Error()))
		return Architecture{}, err
	}
	defer r.Close()

	err = r.Scan(
		&architecture.Id,
		&architecture.ISEName,
		&architecture.RegisterSize,
		&architecture.CreatorId,
		&architecture.CreationDate,
	)
	if err != nil {
		log.Println("ERROR: Cannot scan the architecture object!" + string(err.Error()))
		return Architecture{}, err
	}

	architecture.CreationDate = ConvertSqliteTimestamp(architecture.CreationDate)

	log.Println("INFO: Architecture with Id '" + strconv.Itoa(id) + "' has been retrieved")
	return architecture, nil
}

func GetArchitectureByName(architectureName string) (Architecture, error) {
	log.Println("INFO: Architecture by Name requested: " + architectureName)
	rec, err := DB.Prepare("SELECT * FROM Architectures WHERE ISEName = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return Architecture{}, err
	}
	defer rec.Close()

	architecture := Architecture{}

	r, err := rec.Query(architectureName)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such architecture found in DB: " + string(err.Error()))
			return Architecture{}, nil
		}
		log.Println("ERROR: Cannot retrieve architecture from DB: " + string(err.Error()))
		return Architecture{}, err
	}
	defer r.Close()

	err = r.Scan(
		&architecture.Id,
		&architecture.ISEName,
		&architecture.RegisterSize,
		&architecture.CreatorId,
		&architecture.CreationDate,
	)
	if err != nil {
		log.Println("ERROR: Cannot scan the architecture object!" + string(err.Error()))
		return Architecture{}, err
	}

	architecture.CreationDate = ConvertSqliteTimestamp(architecture.CreationDate)

	log.Println("INFO: Architecture with Name '" + architectureName + "' has been retrieved")
	return architecture, nil
}
