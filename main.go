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
	resources := rice.MustFindBox("view/dist")
	model.MainPage = resources.MustString("app.html")
	router.GET("/", controller.Index)
	router.StaticFS("/view/dist", resources.HTTPBox())
	router.StaticFS("/uploads", web.Dir("uploads"))

	// Deploy
	router.UseController("/install", new(controller.Install), middleware.MustNotInstalled)

	// Login/out
	router.UseController("/login", new(controller.Login), middleware.MustInstalled)
	router.GET("/logout", controller.Logout)

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
