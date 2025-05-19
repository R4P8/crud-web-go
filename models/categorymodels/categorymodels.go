package categorymodels

import (
    "context"
    "curd-web-go/config"
    "curd-web-go/entities"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer = otel.Tracer("category-model")

func GetAll(ctx context.Context) []entities.Category {
    ctx, span := tracer.Start(ctx, "GetALLCategories")
        defer span.End()

    rows, err := config.DB.QueryContext(ctx, "SELECT * FROM categories")
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
func Create(ctx context.Context, category entities.Category) bool {
    ctx, span := tracer.Start(ctx, "CreateCategory")
    defer span.End()

    err := config.DB.QueryRowContext(ctx,`
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

func Detail(ctx context.Context, id int) entities.Category {
    ctx, span := tracer.Start(ctx, "GetCategoryDetail")
    defer span.End()

    row := config.DB.QueryRowContext(ctx, `SELECT id, name FROM categories WHERE id =$1`, id)

    var categories entities.Category
        if err := row.Scan(&categories.ID, &categories.Name); err != nil {
            panic(err.Error())
        }

        return categories
}

func Update(ctx context.Context, id int, categories entities.Category) bool {
    ctx, span := tracer.Start(ctx, "UpdateCategory")
    defer span.End()

    result, err := config.DB.ExecContext(ctx,`
        UPDATE categories SET name = $1, updated_at = $2 WHERE id = $3
        `, categories.Name, categories.UpdatedAt, id)
        if err != nil {
            panic(err)
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
            panic(err)
        }

        return rowsAffected > 0
}

func Delete(ctx context.Context, id int) error {
    ctx, span := tracer.Start(ctx, "DeleteCategory")
    defer span.End()

    _, err := config.DB.ExecContext(ctx, "DELETE FROM categories WHERE id = $1", id)
    return err
}
