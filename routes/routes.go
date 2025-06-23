package routes

import (
	"curd-web-go/controllers/authcontrollers"
	"curd-web-go/controllers/categorycontrollers"
	"curd-web-go/controllers/homecontrollers"
	"curd-web-go/controllers/productcontrollers"
	"curd-web-go/middleware"

	"github.com/julienschmidt/httprouter"
)

func Routes() *httprouter.Router {
	router := httprouter.New()

	//Home
	router.Handler("GET", "/", middleware.Wrap(homecontrollers.Homecontrollers, "Home"))

	//JWT
	router.Handler("POST", "/auth/register", middleware.Wrap(authcontrollers.Register, "Register"))
	router.Handler("POST", "/auth/login", middleware.Wrap(authcontrollers.Login, "Login"))
	router.Handler("GET", "/auth/login", middleware.Wrap(authcontrollers.Login, "LoginGet"))
	router.Handler("POST", "/auth/logout", middleware.Wrap(authcontrollers.Logout, "Logout"))

	//Categories
	router.Handler("GET", "/categories", middleware.Wrap(categorycontrollers.Index, "CategoriesIndex"))
	router.Handler("GET", "/categories/add", middleware.Wrap(categorycontrollers.Add, "CategoriesAddGet"))
	router.Handler("POST", "/categories/add", middleware.Wrap(categorycontrollers.Add, "CategoriesAddPost"))
	router.Handler("GET", "/categories/edit", middleware.Wrap(categorycontrollers.Edit, "CategoriesEditGet"))
	router.Handler("POST", "/categories/edit", middleware.Wrap(categorycontrollers.Edit, "CategoriesEditPost"))
	router.Handler("GET", "/categories/delete", middleware.Wrap(categorycontrollers.Delete, "CategoriesDelete"))

	//Product
	router.Handler("GET", "/product", middleware.Wrap(productcontrollers.Index, "ProductIndex"))
	router.Handler("GET", "/product/add", middleware.Wrap(productcontrollers.Add, "ProductAddGet"))
	router.Handler("POST", "/product/add", middleware.Wrap(productcontrollers.Add, "ProductAddPost"))
	router.Handler("GET", "/product/edit", middleware.Wrap(productcontrollers.Edit, "ProductEditGet"))
	router.Handler("POST", "/product/edit", middleware.Wrap(productcontrollers.Edit, "ProductEditPost"))
	router.Handler("GET", "/product/delete", middleware.Wrap(productcontrollers.Delete, "ProductDelete"))

	return router
}
