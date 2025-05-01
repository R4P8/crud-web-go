package routes

import (
	"curd-web-go/controllers/categorycontrollers"
	"curd-web-go/controllers/homecontrollers"
	"curd-web-go/controllers/productcontrollers"

	"github.com/julienschmidt/httprouter"
)

func Routes() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", homecontrollers.Homecontrollers)

	router.GET("/categories", categorycontrollers.Index)
	router.GET("/categories/add", categorycontrollers.Add)
	router.POST("/categories/add", categorycontrollers.Add)
	router.GET("/categories/edit", categorycontrollers.Edit)
	router.POST("/categories/edit", categorycontrollers.Edit)
	router.GET("/categories/delete", categorycontrollers.Delete)

	router.GET("/product", productcontrollers.Index)
	router.GET("/product/add", productcontrollers.Add)
	router.POST("/product/add", productcontrollers.Add)
	router.GET("/product/edit", productcontrollers.Edit)
	router.POST("/product/edit", productcontrollers.Edit)
	router.GET("/product/delete", productcontrollers.Delete)

	return router
}
