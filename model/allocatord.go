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

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDatabase(dbPath string) error {
	db, err := sql.Open("sqlite3", "file:"+dbPath+"?_foreign_keys=on&cache=shared&mode=rwc&_journal_mode=WAL&_tx_locking_mode=EXCLUSIVE&_busy_timeout=5000&_temp_store=MEMORY&_auto_vacuum=FULL&_synchronous=NORMAL")
	if err != nil {
		return err
	}

	// Set the appropriate pragmas for SQLite
	log.Println("NOTICE: Setting foreign keys to ON")
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return err
	}
	log.Println("NOTICE: Setting journal mode to WAL")
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return err
	}
	log.Println("NOTICE: Settting tx_locking mode to EXCLUSIVE")
	log.Println("NOTICE: Setting busy timeout to 5000ms")
	_, err = db.Exec("PRAGMA busy_timeout = 5000")
	if err != nil {
		return err
	}
	log.Println("NOTICE: Setting temp store to MEMORY")
	_, err = db.Exec("PRAGMA temp_store = MEMORY")
	if err != nil {
		return err
	}
	log.Println("NOTICE: Setting auto_vacuum to FULL")
	_, err = db.Exec("PRAGMA auto_vacuum = FULL")
	if err != nil {
		return err
	}
	log.Println("NOTICE: Setting synchronous to NORMAL")
	_, err = db.Exec("PRAGMA synchronous = NORMAL")
	if err != nil {
		return err
	}

	DB = db
	return nil
}
