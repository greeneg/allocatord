--
-- File generated with SQLiteStudio v3.4.4 on Sun Jun 9 13:32:53 2024
--
-- Text encoding used: UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: Architectures
DROP TABLE IF EXISTS Architectures;

CREATE TABLE IF NOT EXISTS Architectures (
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


-- Table: Audit
DROP TABLE IF EXISTS Audit;

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


-- Table: Buildings
DROP TABLE IF EXISTS Buildings;

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


-- Table: MachineRoles
DROP TABLE IF EXISTS MachineRoles;

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


-- Table: NetworkInterfaces
DROP TABLE IF EXISTS NetworkInterfaces;

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


-- Table: OperatingSystemFamilies
DROP TABLE IF EXISTS OperatingSystemFamilies;

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


-- Table: OperatingSystems
DROP TABLE IF EXISTS OperatingSystems;

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


-- Table: OrganizationalUnits
DROP TABLE IF EXISTS OrganizationalUnits;

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

INSERT INTO OrganizationalUnits (
                                    Id,
                                    OUName,
                                    Description,
                                    CreatorId,
                                    CreationDate
                                )
                                VALUES (
                                    1,
                                    'Unassigned',
                                    'The OU used as a place holder when a system changes hands',
                                    1,
                                    '2024-06-01 15:38:42'
                                );


-- Table: Roles
DROP TABLE IF EXISTS Roles;

CREATE TABLE IF NOT EXISTS Roles (
    Id           INTEGER  PRIMARY KEY AUTOINCREMENT,
    RoleName     STRING   UNIQUE
                          NOT NULL,
    Description  STRING   NOT NULL,
    CreationDate DATETIME NOT NULL
                          DEFAULT (CURRENT_TIMESTAMP) 
);

INSERT INTO Roles (
                      Id,
                      RoleName,
                      Description,
                      CreationDate
                  )
                  VALUES (
                      1,
                      'SYSTEM',
                      'Built-in system role',
                      '2024-06-01 14:57:41'
                  );


-- Table: StorageVolumes
DROP TABLE IF EXISTS StorageVolumes;

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


-- Table: SystemModels
DROP TABLE IF EXISTS SystemModels;

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


-- Table: Systems
DROP TABLE IF EXISTS Systems;

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


-- Table: Users
DROP TABLE IF EXISTS Users;

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

INSERT INTO Users (
                      Id,
                      UserName,
                      FullName,
                      Status,
                      OrgUnitId,
                      RoleId,
                      PasswordHash,
                      CreationDate,
                      LastPasswordChangedDate
                  )
                  VALUES (
                      1,
                      'SYSTEM',
                      'Allocator System',
                      'enabled',
                      1,
                      1,
                      '!',
                      '2024-06-01 14:58:36',
                      '2024-06-01 14:58:36'
                  );


-- Table: Vendors
DROP TABLE IF EXISTS Vendors;

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


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
