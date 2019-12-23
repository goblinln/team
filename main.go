package main

import (
	"team/config"
	"team/controller"
	"team/middleware"
	"team/orm"
	"team/web"

	rice "github.com/GeertJohan/go.rice"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration.
	config.Default.Load()

	// Open database
	if config.Default.Installed {
		if err := orm.OpenDB("mysql", config.Default.GetMySQLAddr()); err != nil {
			web.Logger.Fatal("Failed to connect to database: %s. Reason: %v", config.Default.GetMySQLAddr(), err)
		}
	}

	// Create router for httpd service
	router := web.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.PanicAsError)

	// Resources.
	router.SetPage("/", rice.MustFindBox("view/dist").MustString("app.html"))
	router.StaticFS("/assets", rice.MustFindBox("view/dist").HTTPBox())
	router.StaticFS("/uploads", web.Dir("uploads"))

	// Home/Login/Logout
	router.GET("/home", controller.Home)
	router.GET("/logout", controller.Logout)
	router.POST("/login", controller.Login)

	// Install
	router.UseController("/install", new(controller.Install), middleware.MustNotInstalled)

	// Normal API
	api := router.Group("/api")
	api.Use(middleware.MustInstalled)
	api.Use(middleware.AutoLogin)
	api.Use(middleware.MustLogined)
	api.UseController("/user", new(controller.User))
	api.UseController("/task", new(controller.Task))
	api.UseController("/project", new(controller.Project))
	api.UseController("/document", new(controller.Document))
	api.UseController("/file", new(controller.File))
	api.UseController("/notice", new(controller.Notice))

	// Admin API
	router.UseController(
		"/admin",
		new(controller.Admin),
		middleware.MustInstalled,
		middleware.AutoLogin,
		middleware.MustLoginedAsAdmin)

	// Start service.
	router.Start(config.Default.AppPort)
}
