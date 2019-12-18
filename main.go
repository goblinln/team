package main

import (
	"team/controller"
	"team/middleware"
	"team/model"
	"team/web"

	rice "github.com/GeertJohan/go.rice"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Open database connections.
	model.Environment.Prepare()

	// Create router for httpd service
	router := web.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.PanicAsError)

	// Resources.
	router.SetPage("/", model.Page)
	router.StaticFS("/view/dist", rice.MustFindBox("view/dist").HTTPBox())
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
	router.Start(model.Environment.AppPort)
}
