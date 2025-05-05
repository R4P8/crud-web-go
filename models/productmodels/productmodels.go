package productmodels

import (
	"curd-web-go/config"
	"curd-web-go/entities"
)

func GetAll() []entities.Product {
	rows, err := config.DB.Query(`
		SELECT
			products.id,
			products.name,
			products.price,
			products.stock,
			products.description,
			products.category_id,
			products.created_at,
			products.updated_at,
			categories.name
		FROM products 
		JOIN categories ON products.category_id = categories.id
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var products []entities.Product

	for rows.Next() {
		var product entities.Product
		var categoryName string
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Stock,
			&product.Description,
			&product.CategoryID,
			&product.CreatedAt,
			&product.UpdatedAt,
			&categoryName,
		); err != nil {
			panic(err)
		}

		product.Category.Name = categoryName

		products = append(products, product)
	}

	return products
}

func Create(product entities.Product) bool {
	err := config.DB.QueryRow(`
		INSERT INTO products (name, price,  category_id, stock, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`,
		product.Name,
		product.Price,
		product.CategoryID,
		product.Stock,
		product.Description,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID)

	if err != nil {
		panic(err)
	}

	return product.ID > 0
}

func Update(id int, product entities.Product) bool {
	query, err := config.DB.Exec(`
		UPDATE products SET 
			name = $1,
			price = $2,
			category_id = $3,
			stock = $4,
			description = $5,
			updated_at = $6 
			WHERE id = $7`,

		product.Name,
		product.Price,
		product.CategoryID,
		product.Stock,
		product.Description,
		product.UpdatedAt,
		id,
	)
	if err != nil {
		panic(err)
	}

	result, err := query.RowsAffected()
	if err != nil {
		panic(err)
	}

	return result > 0
}

func Delete(id int) error {
	_, err := config.DB.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}
