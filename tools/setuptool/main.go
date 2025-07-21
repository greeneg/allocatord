package main

import (
	"os"
	"path"
	"strings"

	"github.com/pborman/getopt/v2"
)

var VERSION string = "0.1.0"

var appFullPath, _ = os.Executable()
var app = path.Base(appFullPath)

// setup the global flags

var (
	dbFile             string
	account            string
	fullName           string
	orgUnitName        string
	orgUnitDescription string
	role               string
	roleDescription    string
	optHelp            = getopt.BoolLong("help", 'h', "This help message")
	optVersion         = getopt.BoolLong("version", 'v', "Show the version")
)

func showHelp() {
	println(app + " - Setup tool for Allocator Daemon")
	dividerLine := strings.Repeat("=", 43)
	println(dividerLine)
	println("Add and configure roles or accounts for the Allocator Daemon\n")
	println("OPTIONS:")
	println("   -d|--database-file FILENAME_PATH       REQUIRED: The full or relative path")
	println("                                          to the database file")
	println("   -a|--account ACCOUNT_NAME              OPTIONAL: The account to create")
	println("   -r|--role ROLE_NAME                    OPTIONAL: The role to create")
	println("   -o|--org-unit ORGANIZATIONAL_UNIT_NAME OPTIONAL: The organizational unit to")
	println("                                          create")
	println("   -f|--fullname QUOTED_FULLNAME          CONDITIONALLY OPTIONAL: If the")
	println("                                          account flag is set, this is required.")
	println("                                          This should be the full name or")
	println("                                          description for the account to be")
	println("                                          registered with the system.")
	println("   -D|--role-description ROLE_DESCRIPTION CONDITIONALLY OPTIONAL: If the")
	println("                                          role flag is set, this is required.")
	println("                                          This should be the description for")
	println("                                          the role to be registered with the")
	println("                                          system.")
	println("   -O|--org-unit-description DESCRIPTION  CONDITIONALLY OPTIONAL: If the")
	println("                                          org unit flag is set, this is required.")
	println("                                          This describes the organizational unit")
	println("                                          for the org-unit to be registered with")
	println("                                          the system.")
	println("")
	println("Author: Gary L. Greene, Jr. <greeneg@tolharadys.net>")
	println("License: Apache Public License, v2")
	showVersion()
}

func showVersion() {
	println("version: " + VERSION)
}

func init() {
	getopt.FlagLong(&dbFile, "database-file", 'd', "The full path to the database file")
	getopt.FlagLong(&account, "account", 'a', "The account to add to the system")
	getopt.FlagLong(&fullName, "fullname", 'f', "The full name to associate with the account")
	getopt.FlagLong(&orgUnitName, "org-unit", 'o', "The name of the organizational unit")
	getopt.FlagLong(&orgUnitDescription, "org-unit-description", 'O', "The description of the organizational unit to process")
	getopt.FlagLong(&role, "role", 'r', "The role to add to the system")
	getopt.FlagLong(&roleDescription, "role-description", 'D', "The description of the role to process")
}

func processFlags() {
	if *optHelp {
		showHelp()
		os.Exit(0)
	}

	if *optVersion {
		showVersion()
		os.Exit(0)
	}

	println("Starting setuptool... \n")

	// first, setup our DB connection
	if dbFile != "" {
		println("Database file: " + dbFile)
		ConnectDatabase(dbFile)
		infoPrintln("Database connection completed")
	} else {
		errPrintln("Database file must be defined.")
		showHelp()
		os.Exit(1)
	}

	// do we need to process an account the user passed in?
	if account != "" {
		println("Account: " + account)
		if fullName != "" {
			println("Account fullname: " + fullName)
		} else {
			errPrintln("Account must have a full name")
			showHelp()
			os.Exit(1)
		}
	}

	// next, do we need to process an organizational unit?
	if orgUnitName != "" {
		println("Organizational Unit: " + orgUnitName)
		if orgUnitDescription != "" {
			println("Organizational Unit Description: " + orgUnitDescription)
		} else {
			errPrintln("Organizational unit must have a description")
			showHelp()
			os.Exit(1)
		}
	}

	// How about processing roles?
	if role != "" {
		println("Role: " + role)
		if roleDescription != "" {
			println("Role description: " + roleDescription)
		} else {
			errPrintln("Role must have a role description")
			showHelp()
			os.Exit(1)
		}
	}
}

func main() {
	getopt.Parse()
	processFlags()

	creator, err := getAccountByName("SYSTEM")
	if err != nil {
		errPrintln("Encountered error when looking up the 'SYSTEM' account")
		os.Exit(1)
	}

	ouRecord, err := createInformationTechnologyOu(creator)
	if err != nil {
		errPrintln("Encountered error when creating the InformationTechnology organizational unit: " + string(err.Error()))
		os.Exit(1)
	}

	roleRecord, err := createAdministratorsRole(creator)
	if err != nil {
		errPrintln("Encountered error when creating the Administrators role: " + string(err.Error()))
		os.Exit(1)
	}

	// now handle our admin user
	_, err = createAdminAccount(creator, ouRecord, roleRecord)
	if err != nil {
		errPrintln("Encountered error when creating the admin account: " + string(err.Error()))
		os.Exit(1)
	}
}
