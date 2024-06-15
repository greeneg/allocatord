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

func CreateMachineRole(m MachineRole, id int) (bool, error) {
	log.Println("INFO: Machine Role creation requested: " + m.MachineRoleName)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	q, err := t.Prepare("INSERT INTO MachineRoles (MachineRoleName, Description, CreatorId) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(m.MachineRoleName, id)
	if err != nil {
		log.Println("ERROR: Cannot create machine role '" + m.MachineRoleName + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: Machine Role '" + m.MachineRoleName + "' created")
	return true, nil
}

func DeleteMachineRole(machineRoleId int) (bool, error) {
	log.Println("INFO: Machine Role deletion requested: " + strconv.Itoa(machineRoleId))
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	q, err := DB.Prepare("DELETE FROM MachineRoles WHERE Id IS ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(machineRoleId)
	if err != nil {
		log.Println("ERROR: Cannot delete machine role '" + strconv.Itoa(machineRoleId) + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: Machine Role with Id '" + strconv.Itoa(machineRoleId) + "' has been deleted")
	return true, nil
}

func GetMachineRoles() ([]MachineRole, error) {
	log.Println("INFO: List of machine role objects requested")
	rows, err := DB.Query("SELECT * FROM MachineRoles")
	if err != nil {
		log.Println("ERROR: Could not run the DB query!" + string(err.Error()))
		return nil, err
	}

	machineRoles := make([]MachineRole, 0)
	for rows.Next() {
		machineRole := MachineRole{}
		err = rows.Scan(
			&machineRole.Id,
			&machineRole.MachineRoleName,
			&machineRole.Description,
			&machineRole.CreatorId,
			&machineRole.CreationDate,
		)
		if err != nil {
			log.Println("ERROR: Cannot marshal the machine role objects!" + string(err.Error()))
			return nil, err
		}

		machineRole.CreationDate = ConvertSqliteTimestamp(machineRole.CreationDate)

		machineRoles = append(machineRoles, machineRole)
	}

	log.Println("INFO: List of all machine roles retrieved")
	return machineRoles, nil
}

func GetMachineRoleById(id int) (MachineRole, error) {
	log.Println("INFO: Machine role by Id requested: " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM MachineRoles WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return MachineRole{}, err
	}

	machineRole := MachineRole{}
	err = rec.QueryRow(id).Scan(
		&machineRole.Id,
		&machineRole.MachineRoleName,
		&machineRole.Description,
		&machineRole.CreatorId,
		&machineRole.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such machine role found in DB: " + string(err.Error()))
			return MachineRole{}, nil
		}
		log.Println("ERROR: Cannot retrieve machine role from DB: " + string(err.Error()))
		return MachineRole{}, err
	}

	machineRole.CreationDate = ConvertSqliteTimestamp(machineRole.CreationDate)

	return machineRole, nil
}

func GetMachineRoleByName(machineRoleName string) (MachineRole, error) {
	log.Println("INFO: Machine Role by Name requested: " + machineRoleName)
	rec, err := DB.Prepare("SELECT * FROM MachineRoles WHERE MachineRoleName = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return MachineRole{}, err
	}

	machineRole := MachineRole{}
	err = rec.QueryRow(machineRoleName).Scan(
		&machineRole.Id,
		&machineRole.MachineRoleName,
		&machineRole.Description,
		&machineRole.CreatorId,
		&machineRole.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such machine role found in DB: " + string(err.Error()))
			return MachineRole{}, nil
		}
		log.Println("ERROR: Cannot retrieve machine role from DB: " + string(err.Error()))
		return MachineRole{}, err
	}

	machineRole.CreationDate = ConvertSqliteTimestamp(machineRole.CreationDate)

	return machineRole, nil
}

func UpdateMachineRoleById(machineRoleId int, m MachineRole) (bool, error) {
	log.Println("INFO: Update machine role by Id requested: " + strconv.Itoa(machineRoleId))
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := t.Prepare("UPDATE MachineRoles SET MachineRoleName = ?, Description =? WHERE Id = ?")
	if err != nil {
		return false, err
	}

	machineRole, err := json.Marshal(m)
	if err != nil {
		return false, err
	}
	_, err = q.Exec(machineRole, machineRoleId)
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}
