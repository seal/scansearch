package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/seal/scansearch/pkg/controllers"
	"github.com/seal/scansearch/pkg/middleware"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(router chi.Router) {
	/*
	   This "leaks" the middleware, which means the fileserver tries to use user authentication for say, a favicon.ico request.
	   	router.Use(middleware.DeserializeUser)
	   	router.Get("/me", uc.userController.GetMe)
	*/

	userGroup := router.Group(nil)
	userGroup.Use(middleware.DeserializeUser)
	userGroup.Get("/user", uc.userController.GetMe)
	userGroup.Put("/user", uc.userController.PutUser)
	userGroup.Delete("/user", uc.userController.DeleteUser)
	userGroup.Get("/wardrobe", uc.userController.GetWardrobe)
	userGroup.Post("/wardrobe", uc.userController.AddWardrobe)
	userGroup.Put("/wardrobe", uc.userController.PutWardrobe)
	userGroup.Delete("/wardrobe", uc.userController.DeleteWardrobe)
}
