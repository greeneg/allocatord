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
	"encoding/json"
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
//	@version		0.0.1
//	@description	An API for managing OS imaging

//	@contact.name	Gary Greene
//	@contact.url	https://github.com/greeneg/allocatord

//	@securityDefinitions.basic	BasicAuth

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:5000
//	@BasePath	/api/v1

// @schemas	http https
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
	AllocatorD := new(controllers.Allocator)
	AllocatorD.AppPath = appdir
	AllocatorD.ConfigPath = configDir
	AllocatorD.ConfStruct = config

	err = model.ConnectDatabase(AllocatorD.ConfStruct.DbPath)
	helpers.FatalCheckError(err)

	// set up our static assets
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*.html")

	// some defaults for using session support
	r.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))
	// frontend
	fePublic := r.Group("/")
	routes.FePublicRoutes(fePublic, AllocatorD)

	fePrivate := r.Group("/")
	fePrivate.Use(middleware.AuthCheck)
	routes.FePrivateRoutes(fePrivate, AllocatorD)

	// API
	public := r.Group("/api/v1")
	routes.PublicRoutes(public, AllocatorD)

	private := r.Group("/api/v1")
	private.Use(middleware.AuthCheck)
	routes.PrivateRoutes(private, AllocatorD)

	// swagger doc
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	tcpPort := strconv.Itoa(AllocatorD.ConfStruct.TcpPort)
	tlsTcpPort := strconv.Itoa(AllocatorD.ConfStruct.TLSTcpPort)
	tlsPemFile := AllocatorD.ConfStruct.TLSPemFile
	tlsKeyFile := AllocatorD.ConfStruct.TLSKeyFile
	if AllocatorD.ConfStruct.UseTLS {
		r.RunTLS(":"+tlsTcpPort, tlsPemFile, tlsKeyFile)
	} else {
		r.Run(":" + tcpPort)
	}
}
