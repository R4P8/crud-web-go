package categorymodels

import (
	"curd-web-go/config"
	"curd-web-go/entities"
)

func GetAll() []entities.Category {
	rows, err := config.DB.Query("SELECT * FROM categories")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var categories []entities.Category

	for rows.Next() {
		var category entities.Category
		err := rows.Scan(&category.ID, &category.Name, &category.UpdatedAt, &category.CreatedAt)
		if err != nil {
			panic(err)
		}
		categories = append(categories, category)
	}

	return categories
}

// Dummy function biar gak error compile
func Create(category entities.Category) bool {
	err := config.DB.QueryRow(`
		INSERT INTO categories (name, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`,
		category.Name,
		category.CreatedAt,
		category.UpdatedAt,
	).Scan(&category.ID)

	if err != nil {
		panic(err)
	}

	return category.ID > 0
}

func Detail(id int) entities.Category {
	row := config.DB.QueryRow(`SELECT id, name FROM categories WHERE id =$1`, id)

	var categories entities.Category
	if err := row.Scan(&categories.ID, &categories.Name); err != nil {
		panic(err.Error())
	}

	return categories
}

func Update(id int, categories entities.Category) bool {
	query, err := config.DB.Exec(`
		UPDATE categories SET name = $1, updated_at = $2 WHERE id = $3
	`, categories.Name, categories.UpdatedAt, id)
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
	_, err := config.DB.Exec("DELETE FROM categories WHERE id = $1", id)
	return err
}
