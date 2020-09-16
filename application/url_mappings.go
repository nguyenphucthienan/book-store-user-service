package application

import (
	"github.com/nguyenphucthienan/book-store-user-service/controller/user"
)

const (
	apiPrefix = "/api"
)

func mapUrls() {
	router.POST(apiPrefix+"/users", user.Create)
	router.GET(apiPrefix+"/users/:user_id", user.Get)
	router.PUT(apiPrefix+"/users/:user_id", user.Update)
	router.PATCH(apiPrefix+"/users/:user_id", user.Update)
	router.DELETE(apiPrefix+"/users/:user_id", user.Delete)
	router.POST(apiPrefix+"/users/login", user.Login)
	router.GET(apiPrefix+"/internal/users/search", user.Search)
}
