package homecontrollers

import (
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

func Homecontrollers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	temp, err := template.ParseFiles("views/index.html")
	if err != nil {
		panic(err)
	}

	temp.Execute(w, nil)
}
