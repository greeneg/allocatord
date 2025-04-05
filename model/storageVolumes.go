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

func CreateStorageVolume(s StorageVolume, id int) (bool, error) {
	log.Println("INFO: Storage Volume creation requested: " + s.VolumeName)
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

	q, err := t.Prepare("INSERT INTO StorageVolumes (VolumeName, StorageType, DeviceModel, DeviceId, MountPoint, VolumeSize, VolumeFormat, VolumeLabel, SystemId, CreatorId) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(s.VolumeName, s.StorageType, s.DeviceModel, s.DeviceId, s.MountPoint, s.VolumeSize, s.VolumeFormat, s.VolumeLabel, s.SystemId, id)
	if err != nil {
		log.Println("ERROR: Cannot create storage volume '" + s.VolumeName + "': " + string(err.Error()))
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Could not commit the DB transaction!" + string(err.Error()))
		return false, err
	}

	log.Println("INFO: Storage Volume '" + s.VolumeName + "' created")
	return true, nil
}

func DeleteStorageVolume(storageVolumeId int) (bool, error) {
	log.Println("INFO: Storage Volume deletion requested: " + strconv.Itoa(storageVolumeId))
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

	q, err := DB.Prepare("DELETE FROM StorageVolumes WHERE Id IS ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(storageVolumeId)
	if err != nil {
		log.Println("ERROR: Cannot delete storage volume with ID '" + strconv.Itoa(storageVolumeId) + "': " + string(err.Error()))
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Could not commit the DB transaction!" + string(err.Error()))
		return false, err
	}

	log.Println("INFO: Storage Volume with Id '" + strconv.Itoa(storageVolumeId) + "' has been deleted")
	return true, nil
}

func GetStorageVolumes() ([]StorageVolume, error) {
	log.Println("INFO: List of storage volume object requested")
	rows, err := DB.Query("SELECT * FROM StorageVolumes")
	if err != nil {
		log.Println("ERROR: Could not run the DB query!" + string(err.Error()))
		return nil, err
	}
	defer rows.Close()

	volumes := make([]StorageVolume, 0)
	for rows.Next() {
		volume := StorageVolume{}
		err = rows.Scan(
			&volume.Id,
			&volume.VolumeName,
			&volume.StorageType,
			&volume.DeviceModel,
			&volume.DeviceId,
			&volume.MountPoint,
			&volume.VolumeSize,
			&volume.VolumeFormat,
			&volume.VolumeLabel,
			&volume.SystemId,
			&volume.CreatorId,
			&volume.CreationDate,
		)
		if err != nil {
			log.Println("ERROR: Cannot marshal the storage volume objects!" + string(err.Error()))
			return nil, err
		}

		volume.CreationDate = ConvertSqliteTimestamp(volume.CreationDate)

		volumes = append(volumes, volume)
	}

	log.Println("INFO: List of all storage volumes retrieved")
	return volumes, nil
}

func GetStorageVolumeById(id int) (StorageVolume, error) {
	log.Println("INFO: Storage Volume by Id requested: " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM StorageVolumes WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return StorageVolume{}, err
	}
	defer rec.Close()

	volume := StorageVolume{}

	r, err := rec.Query(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such storage volume found in DB: " + string(err.Error()))
			return StorageVolume{}, nil
		}
		log.Println("ERROR: Cannot retrieve storage volume from DB: " + string(err.Error()))
		return StorageVolume{}, err
	}
	defer r.Close()

	err = r.Scan(
		&volume.Id,
		&volume.VolumeName,
		&volume.StorageType,
		&volume.DeviceModel,
		&volume.DeviceId,
		&volume.MountPoint,
		&volume.VolumeSize,
		&volume.VolumeFormat,
		&volume.VolumeLabel,
		&volume.SystemId,
		&volume.CreatorId,
		&volume.CreationDate,
	)
	if err != nil {
		log.Println("ERROR: Cannot scan the storage volume object!" + string(err.Error()))
		return StorageVolume{}, err
	}

	volume.CreationDate = ConvertSqliteTimestamp(volume.CreationDate)

	log.Println("INFO: Storage Volume with Id '" + strconv.Itoa(volume.Id) + "' has been retrieved")
	return volume, nil
}

func GetStorageVolumeByLabel(label string, id int) (StorageVolume, error) {
	log.Println("INFO: Storage Volume by label requested: " + label)
	rec, err := DB.Prepare("SELECT * FROM StorageVolumes WHERE Id = ? AND VolumeLabel = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return StorageVolume{}, err
	}
	defer rec.Close()

	volume := StorageVolume{}

	r, err := rec.Query(id, label)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such storage volume found in DB: " + string(err.Error()))
			return StorageVolume{}, nil
		}
		log.Println("ERROR: Cannot retrieve storage volume from DB: " + string(err.Error()))
		return StorageVolume{}, err
	}
	defer r.Close()

	err = r.Scan(
		&volume.Id,
		&volume.VolumeName,
		&volume.StorageType,
		&volume.DeviceModel,
		&volume.DeviceId,
		&volume.MountPoint,
		&volume.VolumeSize,
		&volume.VolumeFormat,
		&volume.VolumeLabel,
		&volume.SystemId,
		&volume.CreatorId,
		&volume.CreationDate,
	)
	if err != nil {
		log.Println("ERROR: Cannot scan the storage volume object!" + string(err.Error()))
		return StorageVolume{}, err
	}

	volume.CreationDate = ConvertSqliteTimestamp(volume.CreationDate)

	log.Println("INFO: Storage Volume with Id '" + strconv.Itoa(volume.Id) + "' has been retrieved")
	return volume, nil
}

func GetStorageVolumesBySystemId(systemId int) ([]StorageVolume, error) {
	log.Println("INFO: Storage Volumes by System Id requested: " + strconv.Itoa(systemId))
	rec, err := DB.Prepare("SELECT * FROM StorageVolumes WHERE SystemId = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return nil, err
	}
	defer rec.Close()

	rows, err := rec.Query(systemId)
	if err != nil {
		log.Println("ERROR: Could not query DB: " + string(err.Error()))
		return nil, err
	}
	defer rows.Close()

	storageVolumes := make([]StorageVolume, 0)

	for rows.Next() {
		storageVolume := StorageVolume{}
		err = rows.Scan(
			&storageVolume.Id,
			&storageVolume.VolumeName,
			&storageVolume.StorageType,
			&storageVolume.DeviceModel,
			&storageVolume.DeviceId,
			&storageVolume.MountPoint,
			&storageVolume.VolumeSize,
			&storageVolume.VolumeFormat,
			&storageVolume.SystemId,
			&storageVolume.CreatorId,
			&storageVolume.CreationDate,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("ERROR: No such storage volumes found in DB: " + string(err.Error()))
				return nil, nil
			}
			log.Println("ERROR: Cannot retrieve storage volumes from DB: " + string(err.Error()))
			return nil, err
		}

		storageVolume.CreationDate = ConvertSqliteTimestamp(storageVolume.CreationDate)

		storageVolumes = append(storageVolumes, storageVolume)
	}

	log.Println("INFO: Storage Volumes with System Id '" + strconv.Itoa(systemId) + "' have been retrieved")
	return storageVolumes, nil
}

func UpdateStorageVolume(id int, s StorageVolume) (bool, error) {
	log.Println("INFO: Update storage volume by Id requested: " + strconv.Itoa(id))
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

	q, err := t.Prepare("UPDATE StorageVolumes SET VolumeName = ?, StorageType = ?, DeviceModel = ?, DeviceId = ?, MountPoint = ?, VolumeSize = ?, VolumeFormat = ?, VolumeLabel = ?, SystemId = ? WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	storageVolume, err := json.Marshal(s)
	if err != nil {
		log.Println("ERROR: Cannot marshal the storage volume object!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(storageVolume, id)
	if err != nil {
		log.Println("ERROR: Cannot update storage volume with ID '" + strconv.Itoa(id) + "': " + string(err.Error()))
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Could not commit the DB transaction!" + string(err.Error()))
		return false, err
	}

	log.Println("INFO: Storage Volume with Id '" + strconv.Itoa(id) + "' has been updated")
	return true, nil
}
