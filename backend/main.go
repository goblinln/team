package main

import (
	"team/controller"
	"team/middleware"
	"team/model"
	"team/web"
	"team/web/orm"
)

func main() {
	// Open database connections.
	if model.Environment.Installed {
		mysql := model.Environment.MySQL
		err := orm.ConnectDB(
			mysql.Host,
			mysql.User,
			mysql.Password,
			mysql.Database,
			mysql.MaxConns)
		if err != nil {
			web.Logger.Fatal("Failed to prepare database : %s", err.Error())
		}
	}

	// Create router for httpd service
	router := web.New()
	router.Use(middleware.Logger)

	// Resources.
	router.GET("/", controller.Index)
	router.StaticFS(`/www/[\s\S]+`, "www", "/www/")
	router.StaticFS(`/uploads/[\s\S]+`, "uploads", "/uploads/")

	// Deploy
	router.UseController("/install", new(controller.Install), middleware.MustNotInstalled)

	// Login/out
	router.POST("/login", controller.Login, middleware.MustInstalled)
	router.POST("/logout", controller.Logout, middleware.MustInstalled)

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
	web.Logger.Info("Service started at %s", model.Environment.AppPort)
	web.Start(model.Environment.AppPort, router)
}
