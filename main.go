package main

/*

  Giron-Service - Golang-based web service for managing panel events

  Author:  Gary L. Greene, Jr.
  License: Apache v2.0

  Copyright 2024, YggdrasilSoft, LLC

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
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/greeneg/allocatord/controllers"
	_ "github.com/greeneg/allocatord/docs"
	"github.com/greeneg/allocatord/globals"
	"github.com/greeneg/allocatord/helpers"
	"github.com/greeneg/allocatord/middleware"
	"github.com/greeneg/allocatord/model"
	"github.com/greeneg/allocatord/routes"
)

//	@title			Allocator Daemon
//	@version		0.1.1
//	@description	An API for managing OS imaging

//	@contact.name	Gary Greene
//	@contact.url	https://github.com/greeneg/allocatord

//	@securityDefinitions.basic	BasicAuth

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:5000
//	@BasePath	/api/v1

// @schemas	http https

func createDB(dbName string) (bool, error) {
	log.Println("INFO: DB doesn't exist. Attempt to create it")
	const schema string = `CREATE TABLE IF NOT EXISTS Architectures (
		Id           INTEGER  PRIMARY KEY AUTOINCREMENT
							  NOT NULL
							  UNIQUE,
		ISEName      STRING   UNIQUE
							  NOT NULL,
		CreatorId    INTEGER  NOT NULL
							  REFERENCES Users (Id),
		CreationDate DATETIME NOT NULL
							  DEFAULT (CURRENT_TIMESTAMP) 
	);
	CREATE TABLE IF NOT EXISTS Audit (
		Id           INTEGER  PRIMARY KEY AUTOINCREMENT
							  NOT NULL
							  UNIQUE,
		ChangedById  INTEGER  REFERENCES Users (Id) 
							  NOT NULL,
		TableChanged STRING   NOT NULL,
		ChangeClass  STRING   NOT NULL,
		ChangeDate   DATETIME NOT NULL
							  DEFAULT (CURRENT_TIMESTAMP) 
	);
	CREATE TABLE IF NOT EXISTS Buildings (
		Id           INTEGER  PRIMARY KEY AUTOINCREMENT
							  UNIQUE
							  NOT NULL,
		BuildingName STRING   NOT NULL
							  UNIQUE,
		ShortName    STRING   NOT NULL
							  UNIQUE,
		City         STRING   NOT NULL,
		Region       STRING   NOT NULL,
		CreatorId    INTEGER  REFERENCES Users (Id) 
							  NOT NULL,
		CreationDate DATETIME NOT NULL
							  DEFAULT (CURRENT_TIMESTAMP) 
	);
	CREATE TABLE IF NOT EXISTS MachineRoles (
		Id              INTEGER  PRIMARY KEY AUTOINCREMENT
								 UNIQUE
								 NOT NULL,
		MachineRoleName STRING   UNIQUE
								 NOT NULL,
		Description     STRING   NOT NULL,
		CreatorId       INTEGER  REFERENCES Users (Id) 
								 NOT NULL,
		CreationDate    DATETIME NOT NULL
								 DEFAULT (CURRENT_TIMESTAMP) 
	);
	CREATE TABLE IF NOT EXISTS NetworkInterfaces (
		Id           INTEGER  PRIMARY KEY AUTOINCREMENT
							  NOT NULL
							  UNIQUE,
		DeviceModel  STRING   NOT NULL,
		DeviceId     STRING   NOT NULL,
		MACAddress   STRING   NOT NULL
							  UNIQUE,
		SystemId     INTEGER  REFERENCES Systems (Id) 
							  NOT NULL,
		IpAddress    STRING   NOT NULL,
		Bitmask      INTEGER  NOT NULL,
		Gateway      STRING   NOT NULL,
		CreatorId    INTEGER  REFERENCES Users (Id) 
							  NOT NULL,
		CreationDate DATETIME NOT NULL
							  DEFAULT (CURRENT_TIMESTAMP) 
	);
	CREATE TABLE IF NOT EXISTS OperatingSystemFamilies (
		Id           INTEGER  PRIMARY KEY AUTOINCREMENT
							  UNIQUE
							  NOT NULL,
		OSFamilyName STRING   UNIQUE
							  NOT NULL,
		CreatorId    INTEGER  REFERENCES Users (Id) 
							  NOT NULL,
		CreationDate DATETIME NOT NULL
							  DEFAULT (CURRENT_TIMESTAMP) 
	);
	CREATE TABLE IF NOT EXISTS OperatingSystems (
		Id           INTEGER  PRIMARY KEY AUTOINCREMENT
							  UNIQUE
							  NOT NULL,
		OSName       STRING   UNIQUE
							  NOT NULL,
		OSFamilyId   INTEGER  REFERENCES OperatingSystemFamilies (Id) 
							  NOT NULL,
		OSImageUrl   STRING   UNIQUE
							  NOT NULL,
		CreatorId    INTEGER  REFERENCES Users (Id) 
							  NOT NULL,
		CreationDate DATETIME NOT NULL
							  DEFAULT (CURRENT_TIMESTAMP) 
	);
	CREATE TABLE IF NOT EXISTS OrganizationalUnits (
		Id           INTEGER  PRIMARY KEY AUTOINCREMENT
							  NOT NULL
							  UNIQUE,
		OUName       STRING   UNIQUE
							  NOT NULL,
		Description  STRING   NOT NULL,
		CreatorId    INTEGER  REFERENCES Users (Id) 
							  NOT NULL,
		CreationDate DATETIME NOT NULL
							  DEFAULT (CURRENT_TIMESTAMP) 
	);

	INSERT INTO OrganizationalUnits (Id, OUName, Description, CreatorId, CreationDate)
		VALUES ( 1, 'Unassigned', 'The OU used as a place holder when a system changes hands', 1, '2024-06-01 15:38:42' );

	CREATE TABLE IF NOT EXISTS Roles (
		Id           INTEGER  PRIMARY KEY AUTOINCREMENT,
		RoleName     STRING   UNIQUE
							  NOT NULL,
		Description  STRING   NOT NULL,
		CreationDate DATETIME NOT NULL
							  DEFAULT (CURRENT_TIMESTAMP) 
	);

	INSERT INTO Roles (Id, RoleName, Description, CreationDate)
		VALUES ( 1, 'SYSTEM', 'Built-in system role', '2024-06-01 14:57:41' );

	CREATE TABLE IF NOT EXISTS StorageVolumes (
		Id           INTEGER  PRIMARY KEY AUTOINCREMENT
							  UNIQUE
							  NOT NULL,
		StorageType  STRING   NOT NULL,
		DeviceModel  STRING   NOT NULL,
		DeviceId     STRING   NOT NULL,
		MountPoint   STRING   NOT NULL,
		VolumeSize   INTEGER  NOT NULL,
		VolumeFormat STRING   NOT NULL,
		VolumeLabel  STRING   NOT NULL,
		SystemId     INTEGER  REFERENCES Systems (Id) 
							  NOT NULL,
		CreatorId    INTEGER  REFERENCES Users (Id) 
							  NOT NULL,
		CreationDate DATETIME NOT NULL
							  DEFAULT (CURRENT_TIMESTAMP) 
	);
	CREATE TABLE IF NOT EXISTS SystemModels (
		Id           INTEGER  PRIMARY KEY AUTOINCREMENT
							  UNIQUE
							  NOT NULL,
		ModelName    STRING   NOT NULL
							  UNIQUE,
		CreatorId    INTEGER  REFERENCES Users (Id) 
							  NOT NULL,
		CreationDate DATETIME NOT NULL
							  DEFAULT (CURRENT_TIMESTAMP) 
	);
	CREATE TABLE IF NOT EXISTS Systems (
		Id                INTEGER  PRIMARY KEY AUTOINCREMENT
								   UNIQUE
								   NOT NULL,
		SerialNumber      STRING   NOT NULL
								   UNIQUE,
		ModelId           INTEGER  REFERENCES SystemModels (Id) 
								   NOT NULL,
		OperatingSystemId INTEGER  NOT NULL
								   REFERENCES OperatingSystems (Id),
		Reimage           BOOL     NOT NULL
								   DEFAULT (FALSE),
		HostVars          STRING   NOT NULL,
		BilledToOrgUnitId INTEGER  REFERENCES OrganizationalUnits (Id) 
								   NOT NULL,
		MachineRoleId     INTEGER  NOT NULL
								   REFERENCES MachineRoles (Id),
		BuildingId        INTEGER  REFERENCES Buildings (Id) 
								   NOT NULL,
		VendorId          INTEGER  NOT NULL
								   REFERENCES Vendors (Id),
		ArchitectureId    INTEGER  REFERENCES Architectures (Id) 
								   NOT NULL,
		RAM               INTEGER  NOT NULL,
		CPUCores          INTEGER  NOT NULL,
		CreatorId         INTEGER  REFERENCES Users (Id) 
								   NOT NULL,
		CreationDate      DATETIME NOT NULL
								   DEFAULT (CURRENT_TIMESTAMP) 
	);
	CREATE TABLE IF NOT EXISTS Users (
		Id                      INTEGER  PRIMARY KEY AUTOINCREMENT
										 UNIQUE
										 NOT NULL,
		UserName                STRING   NOT NULL
										 UNIQUE,
		FullName                STRING   NOT NULL,
		Status                  STRING   NOT NULL
										 DEFAULT enabled,
		OrgUnitId               INTEGER  REFERENCES OrganizationalUnits (Id) 
										 NOT NULL,
		RoleId                  INTEGER  REFERENCES Roles (Id) 
										 NOT NULL,
		PasswordHash            STRING   NOT NULL,
		CreationDate            DATETIME NOT NULL
										 DEFAULT (CURRENT_TIMESTAMP),
		LastPasswordChangedDate DATETIME NOT NULL
										 DEFAULT (CURRENT_TIMESTAMP) 
	);

	INSERT INTO Users (Id, UserName, FullName, Status, OrgUnitId, RoleId, PasswordHash, CreationDate, LastPasswordChangedDate)
		VALUES ( 1, 'SYSTEM', 'Allocator System', 'enabled', 1, 1, '!', '2024-06-01 14:58:36', '2024-06-01 14:58:36' );

	CREATE TABLE IF NOT EXISTS Vendors (
		Id           INTEGER  PRIMARY KEY AUTOINCREMENT
							  NOT NULL
							  UNIQUE,
		VendorName   STRING   UNIQUE
							  NOT NULL,
		CreatorId    INTEGER  REFERENCES Users (Id),
		CreationDate DATETIME NOT NULL
							  DEFAULT (CURRENT_TIMESTAMP) 
	);
	`

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		helpers.FatalCheckError(err)
	}
	if _, err := db.Exec(schema); err != nil {
		helpers.FatalCheckError(err)
	}
	return true, err
}

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// lets get our working directory
	appdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	helpers.FatalCheckError(err)

	// config path is derived from app working directory
	configDir := filepath.Join(appdir, "config")

	// now that we have our appdir and configDir, lets read in our app config
	// and marshall it to the Config struct
	config := globals.Config{}
	jsonContent, err := os.ReadFile(filepath.Join(configDir, "config.json"))
	helpers.FatalCheckError(err)
	err = json.Unmarshal(jsonContent, &config)
	helpers.FatalCheckError(err)

	// create an app object that contains our routes and the configuration
	Allocator := new(controllers.Allocator)
	Allocator.AppPath = appdir
	Allocator.ConfigPath = configDir
	Allocator.ConfStruct = config

	if _, err := os.Stat(Allocator.ConfStruct.DbPath); errors.Is(err, os.ErrNotExist) {
		_, err := createDB(Allocator.ConfStruct.DbPath)
		if err != nil {
			helpers.FatalCheckError(err)
		}
	}

	err = model.ConnectDatabase(Allocator.ConfStruct.DbPath)
	helpers.FatalCheckError(err)

	// set up our static assets
	// r.Static("/assets", "./assets")
	// r.LoadHTMLGlob("templates/*.html")

	// some defaults for using session support
	r.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))
	// frontend
	// fePublic := r.Group("/")
	// routes.FePublicRoutes(fePublic, AllocatorD)

	// fePrivate := r.Group("/")
	// fePrivate.Use(middleware.AuthCheck)
	// routes.FePrivateRoutes(fePrivate, AllocatorD)

	// API
	public := r.Group("/api/v1")
	routes.PublicRoutes(public, Allocator)

	private := r.Group("/api/v1")
	private.Use(middleware.AuthCheck)
	routes.PrivateRoutes(private, Allocator)

	// swagger doc
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	tcpPort := strconv.Itoa(Allocator.ConfStruct.TcpPort)
	tlsTcpPort := strconv.Itoa(Allocator.ConfStruct.TLSTcpPort)
	tlsPemFile := Allocator.ConfStruct.TLSPemFile
	tlsKeyFile := Allocator.ConfStruct.TLSKeyFile
	if Allocator.ConfStruct.UseTLS {
		r.RunTLS(":"+tlsTcpPort, tlsPemFile, tlsKeyFile)
	} else {
		r.Run(":" + tcpPort)
	}
}
