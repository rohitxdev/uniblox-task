package repo

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Price is in the smallest unit of the currency
	Price        int    `json:"price"`
	QuantityLeft int    `json:"quantityLeft"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

func (r *Repo) GetProducts(ctx context.Context) ([]Product, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, price, quantity_left, created_at, updated_at FROM products;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Price, &p.QuantityLeft, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *Repo) GetProduct(ctx context.Context, id int) (*Product, error) {
	var p Product
	err := r.db.QueryRowContext(ctx, `SELECT id, name, price, quantity_left, created_at, updated_at FROM products WHERE id=$1 LIMIT 1;`, id).Scan(&p.ID, &p.Name, &p.Price, &p.QuantityLeft, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &p, nil
}
