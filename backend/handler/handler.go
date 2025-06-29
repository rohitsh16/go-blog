package handler

func RegisterRoutes(service Service) {
	service.GinFramework.POST("/v1/post", service.CreatePostHandler)
}
