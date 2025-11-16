package handler

func RegisterRoutes(service Service) {
	service.GinFramework.POST("/v1/post", service.CreatePostHandler)
	// Home page (user-facing)
	service.GinFramework.GET("/", service.HomeHandler)

	// Create page - supports GET (render form) and POST (form submit)
	service.GinFramework.GET("/create", service.CreatePostHandler)
	service.GinFramework.POST("/create", service.CreatePostHandler)

	// Writer profile (writers can create/edit their posts)
	service.GinFramework.GET("/writer", service.WriterDashboardHandler)
	service.GinFramework.GET("/writer/edit/:id", service.WriterEditHandler)
	service.GinFramework.POST("/writer/edit/:id", service.WriterEditHandler)

	// Admin profile (management)
	service.GinFramework.GET("/admin", service.AdminHandler)
}
