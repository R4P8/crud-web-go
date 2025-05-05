package categorycontrollers

import (
	"curd-web-go/entities"
	"curd-web-go/models/categorymodels"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	categories := categorymodels.GetAll()
	data := map[string]any{
		"categories": categories,
	}

	temp, err := template.ParseFiles("views/category/index.html")
	if err != nil {
		panic(err)
	}

	err = temp.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "GET" {
		temp, err := template.ParseFiles("views/category/create.html")
		if err != nil {
			panic(err)
		}
		temp.Execute(w, nil)
	}

	if r.Method == "POST" {
		var categories entities.Category

		categories.Name = r.FormValue("name")
		categories.UpdatedAt = time.Now()
		categories.CreatedAt = time.Now()

		if ok := categorymodels.Create(categories); !ok {
			// kalau gagal insert, render lagi create page
			temp, _ := template.ParseFiles("views/category/create.html")
			temp.Execute(w, nil)
		}

		// Kalau sukses, redirect ke /categories
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
	}
}

func Edit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "GET" {
		temp, err := template.ParseFiles("views/category/edit.html")
		if err != nil {
			panic(err)
		}

		IdString := r.URL.Query().Get("id")
		id, err := strconv.Atoi(IdString)
		if err != nil {
			panic(err)
		}

		categories := categorymodels.Detail(id)

		data := map[string]any{
			"categories": categories,
		}

		temp.Execute(w, data)

	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		idString := r.FormValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			panic(err)
		}

		var categories entities.Category
		categories.Name = r.FormValue("name")
		categories.UpdatedAt = time.Now()

		if ok := categorymodels.Update(id, categories); !ok {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/categories", http.StatusSeeOther) // ke halaman setelah update berhasil
	}
}

func Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	IdString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(IdString)
	if err != nil {
		panic(err)
	}

	if err := categorymodels.Delete(id); err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}
