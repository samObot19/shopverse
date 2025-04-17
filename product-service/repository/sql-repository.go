package repository

import (
    "context"
    "database/sql"
	"fmt"
	"strings"
    "github.com/samObot19/shopverse/product-service/models"
)


type MySQLProductRepository struct {
    DB *sql.DB
}

// NewMySQLProductRepository creates a new instance of MySQLProductRepository
func NewMySQLProductRepository(db *sql.DB) *MySQLProductRepository {
    return &MySQLProductRepository{DB: db}
}

func (r *MySQLProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
    tx, err := r.DB.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.ExecContext(ctx, `
        INSERT INTO products (id, title, description, price, stock, category, ratings, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
        product.ID, product.Title, product.Description, product.Price,
        product.Stock, product.Category, product.Ratings, product.CreatedAt,
    )
    if err != nil {
        return err
    }

    _, err = tx.ExecContext(ctx, `
        INSERT INTO product_attributes (product_id, color)
        VALUES (?, ?)`, product.ID, product.Attributes.Color)
    if err != nil {
        return err
    }

    for _, size := range product.Attributes.Size {
        _, err = tx.ExecContext(ctx, `
            INSERT INTO product_sizes (product_id, size)
            VALUES (?, ?)`, product.ID, size)
        if err != nil {
            return err
        }
    }

    for _, img := range product.Images {
        _, err = tx.ExecContext(ctx, `
            INSERT INTO product_images (product_id, image_url)
            VALUES (?, ?)`, product.ID, img)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (r *MySQLProductRepository) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
    var product models.Product

    err := r.DB.QueryRowContext(ctx, `
        SELECT id, title, description, price, stock, category, ratings, created_at
        FROM products WHERE id = ?`, id).
        Scan(&product.ID, &product.Title, &product.Description, &product.Price,
            &product.Stock, &product.Category, &product.Ratings, &product.CreatedAt)
    if err != nil {
        return nil, err
    }

    err = r.DB.QueryRowContext(ctx, `
        SELECT color FROM product_attributes WHERE product_id = ?`, id).
        Scan(&product.Attributes.Color)
    if err != nil {
        return nil, err
    }

    sizeRows, err := r.DB.QueryContext(ctx, `
        SELECT size FROM product_sizes WHERE product_id = ?`, id)
    if err != nil {
        return nil, err
    }
    defer sizeRows.Close()

    for sizeRows.Next() {
        var size string
        if err := sizeRows.Scan(&size); err != nil {
            return nil, err
        }
        product.Attributes.Size = append(product.Attributes.Size, size)
    }

    imgRows, err := r.DB.QueryContext(ctx, `
        SELECT image_url FROM product_images WHERE product_id = ?`, id)
    if err != nil {
        return nil, err
    }
    defer imgRows.Close()

    for imgRows.Next() {
        var img string
        if err := imgRows.Scan(&img); err != nil {
            return nil, err
        }
        product.Images = append(product.Images, img)
    }

    return &product, nil
}

func (r *MySQLProductRepository) GetAllProducts(ctx context.Context, filters map[string]interface{}) ([]*models.Product, error) {
    query := "SELECT id FROM products WHERE 1=1"
    args := []interface{}{}

    for key, value := range filters {
        query += fmt.Sprintf(" AND %s = ?", key)
        args = append(args, value)
    }

    rows, err := r.DB.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []*models.Product
    for rows.Next() {
        var id string
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }
        product, err := r.GetProductByID(ctx, id)
        if err != nil {
            return nil, err
        }
        products = append(products, product)
    }

    return products, nil
}

func (r *MySQLProductRepository) UpdateProduct(ctx context.Context, id string, updated *models.Product) error {
    _, err := r.DB.ExecContext(ctx, `
        UPDATE products SET title = ?, description = ?, price = ?, stock = ?, category = ?, ratings = ?
        WHERE id = ?`,
        updated.Title, updated.Description, updated.Price,
        updated.Stock, updated.Category, updated.Ratings, id,
    )
    return err
}

func (r *MySQLProductRepository) DeleteProduct(ctx context.Context, id string) error {
    tx, err := r.DB.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, _ = tx.ExecContext(ctx, `DELETE FROM product_attributes WHERE product_id = ?`, id)
    _, _ = tx.ExecContext(ctx, `DELETE FROM product_sizes WHERE product_id = ?`, id)
    _, _ = tx.ExecContext(ctx, `DELETE FROM product_images WHERE product_id = ?`, id)

    _, err = tx.ExecContext(ctx, `DELETE FROM products WHERE id = ?`, id)
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (r *MySQLProductRepository) UpdateStock(ctx context.Context, id string, quantity int) error {
    _, err := r.DB.ExecContext(ctx, `
        UPDATE products SET stock = ? WHERE id = ?`, quantity, id)
    return err
}

func (r *MySQLProductRepository) GetProductsByCategory(ctx context.Context, category string) ([]*models.Product, error) {
    rows, err := r.DB.QueryContext(ctx, `SELECT id FROM products WHERE category = ?`, category)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []*models.Product
    for rows.Next() {
        var id string
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }
        product, err := r.GetProductByID(ctx, id)
        if err != nil {
            return nil, err
        }
        products = append(products, product)
    }

    return products, nil
}

func (r *MySQLProductRepository) SearchProducts(ctx context.Context, query string) ([]*models.Product, error) {
    search := "%" + strings.ToLower(query) + "%"
    rows, err := r.DB.QueryContext(ctx, `
        SELECT id FROM products
        WHERE LOWER(title) LIKE ? OR LOWER(description) LIKE ?`, search, search)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []*models.Product
    for rows.Next() {
        var id string
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }
        product, err := r.GetProductByID(ctx, id)
        if err != nil {
            return nil, err
        }
        products = append(products, product)
    }

    return products, nil
}
