package productcontrollers

import (
        "curd-web-go/config"
        "curd-web-go/entities"
        "curd-web-go/models/categorymodels"
        "curd-web-go/models/productmodels"
        "html/template"
        "net/http"
        "strconv"
        "time"

        "github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
        ctx := r.Context()
        Products := productmodels.GetAll(ctx)

        data := map[string]any{
                "Products": Products,
        }

        temp, err := template.ParseFiles("views/products/index.html")
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
        ctx := r.Context()

        if r.Method == "GET" {
                categories := categorymodels.GetAll(ctx)
                data := map[string]any{
                        "categories": categories,
                }
                temp, err := template.ParseFiles("views/products/create.html")
                if err != nil {
                        panic(err)
                }
                temp.Execute(w, data)
                return
        }

        if r.Method == "POST" {
                var product entities.Product

                // Ambil data dari form
                product.Name = r.FormValue("name")
                product.Description = r.FormValue("description")

                // Convert stock dan price ke int & float
                stock, err := strconv.Atoi(r.FormValue("stock"))
                if err != nil {
                        http.Error(w, "Invalid stock value", http.StatusBadRequest)
                        return
                }
                product.Stock = stock

                price, err := strconv.ParseFloat(r.FormValue("price"), 64)
                if err != nil {
                        http.Error(w, "Invalid price value", http.StatusBadRequest)
                        return
                }
                product.Price = price

                categoryId, err := strconv.Atoi(r.FormValue("category_id"))
                if err != nil {
                        http.Error(w, "Invalid category ID", http.StatusBadRequest)
                        return
                }

                product.CategoryID = uint(categoryId)

                product.CreatedAt = time.Now()
                product.UpdatedAt = time.Now()

                if ok := productmodels.Create(ctx, product); !ok {
                        // gagal insert
                        temp, _ := template.ParseFiles("views/products/create.html")
                        temp.Execute(w, map[string]any{"error": "Gagal menyimpan produk"})
                        return
                }

                http.Redirect(w, r, "/product", http.StatusSeeOther)
        }
}

func Edit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
        ctx := r.Context()

        if r.Method == "GET" {
                temp, err := template.ParseFiles("views/products/edit.html")
                if err != nil {
                        panic(err)
                }

                IdString := r.URL.Query().Get("id")
                id, err := strconv.Atoi(IdString)
                if err != nil {
                        panic(err)
                }

                var product entities.Product
                err = config.DB.QueryRow(`SELECT id, name, price, stock, description, category_id FROM products WHERE id = $1`, id).Scan(
                        &product.ID, &product.Name, &product.Price, &product.Stock, &product.Description, &product.CategoryID,
                )
                if err != nil {
                        http.Error(w, "Product not found", http.StatusNotFound)
                        return
                }

                categories := categorymodels.GetAll(ctx)

                data := map[string]any{
                        "categories": categories,
                        "product":    product,
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

                price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
                stock, _ := strconv.Atoi(r.FormValue("stock"))
                categoryID, _ := strconv.Atoi(r.FormValue("category_id"))

                var product entities.Product
                product.Name = r.FormValue("name")
                product.Price = price
                product.Stock = stock
                product.CategoryID = uint(categoryID)
                product.Description = r.FormValue("description")
                product.UpdatedAt = time.Now()

                if ok := productmodels.Update(ctx,id, product); !ok {
                        http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
                        return
                }

                http.Redirect(w, r, "/product", http.StatusSeeOther) // ke halaman setelah update berhasil
        }
}

func Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
        ctx := r.Context()

        IdString := r.URL.Query().Get("id")
        id, err := strconv.Atoi(IdString)
        if err != nil {
                panic(err)
        }

        if err := productmodels.Delete(ctx, id); err != nil {
                panic(err)
        }

        http.Redirect(w, r, "/product", http.StatusSeeOther)
}
