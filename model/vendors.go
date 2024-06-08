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

func CreateVendor(v Vendor, id int) (bool, error) {
	log.Println("INFO: Vendor creation requested: " + v.VendorName)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	q, err := t.Prepare("INSERT INTO Vendors (VendorName, CreatorId) VALUES (?, ?)")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(v.VendorName, id)
	if err != nil {
		log.Println("ERROR: Cannot create vendor '" + v.VendorName + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: Vendor '" + v.VendorName + "' created")
	return true, nil
}

func DeleteVendor(vendorId int) (bool, error) {
	log.Println("INFO: Vendor deletion requested: " + strconv.Itoa(vendorId))
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	q, err := DB.Prepare("DELETE FROM Vendors WHERE Id IS ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(vendorId)
	if err != nil {
		log.Println("ERROR: Cannot delete vendor '" + strconv.Itoa(vendorId) + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: Vendor with Id '" + strconv.Itoa(vendorId) + "' has been deleted")
	return true, nil
}

func GetVendors() ([]Vendor, error) {
	log.Println("INFO: List of vendor object requested")
	rows, err := DB.Query("SELECT * FROM Vendors")
	if err != nil {
		log.Println("ERROR: Could not run the DB query!" + string(err.Error()))
		return nil, err
	}

	vendors := make([]Vendor, 0)
	for rows.Next() {
		vendor := Vendor{}
		err = rows.Scan(
			&vendor.Id,
			&vendor.VendorName,
			&vendor.CreatorId,
			&vendor.CreationDate,
		)
		if err != nil {
			log.Println("ERROR: Cannot marshal the vendor objects!" + string(err.Error()))
			return nil, err
		}

		vendor.CreationDate = ConvertSqliteTimestamp(vendor.CreationDate)

		vendors = append(vendors, vendor)
	}

	log.Println("INFO: List of all vendors retrieved")
	return vendors, nil
}

func GetVendorById(id int) (Vendor, error) {
	log.Println("INFO: Vendor by Id requested: " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM Vendors WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return Vendor{}, err
	}

	vendor := Vendor{}
	err = rec.QueryRow(id).Scan(
		&vendor.Id,
		&vendor.VendorName,
		&vendor.CreatorId,
		&vendor.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such vendor found in DB: " + string(err.Error()))
			return Vendor{}, nil
		}
		log.Println("ERROR: Cannot retrieve vendor from DB: " + string(err.Error()))
		return Vendor{}, err
	}

	vendor.CreationDate = ConvertSqliteTimestamp(vendor.CreationDate)

	return vendor, nil
}

func GetVendorByName(vendorName string) (Vendor, error) {
	log.Println("INFO: Vendor by Name requested: " + vendorName)
	rec, err := DB.Prepare("SELECT * FROM Vendors WHERE VendorName = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return Vendor{}, err
	}

	vendor := Vendor{}
	err = rec.QueryRow(vendorName).Scan(
		&vendor.Id,
		&vendor.VendorName,
		&vendor.CreatorId,
		&vendor.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such user found in DB: " + string(err.Error()))
			return Vendor{}, nil
		}
		log.Println("ERROR: Cannot retrieve user from DB: " + string(err.Error()))
		return Vendor{}, err
	}

	vendor.CreationDate = ConvertSqliteTimestamp(vendor.CreationDate)

	return vendor, nil
}
