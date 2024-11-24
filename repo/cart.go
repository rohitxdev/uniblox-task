package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrCartNotFound             = errors.New("cart not found")
	ErrCartItemNotFound         = errors.New("cart item not found")
	ErrCartItemQuantityExceeded = errors.New("cart item quantity exceeded")
)

type CartItem struct {
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	ProductID int    `json:"productId"`
	Quantity  int    `json:"quantity"`
}

func (r *Repo) GetCart(ctx context.Context, userID int) ([]CartItem, error) {
	cartItems := make([]CartItem, 0)
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, product_id, quantity, created_at, updated_at FROM cart_items WHERE user_id=$1 ORDER BY created_at DESC;`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cartItem CartItem
		err = rows.Scan(&cartItem.ID, &cartItem.UserID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CreatedAt, &cartItem.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cartItems = append(cartItems, cartItem)
	}
	return cartItems, nil
}

func (r *Repo) GetCartItem(ctx context.Context, userID int, productID int) (*CartItem, error) {
	var cartItem CartItem
	err := r.db.QueryRowContext(ctx, `SELECT id, user_id, product_id, quantity, created_at, updated_at FROM cart_items WHERE user_id=$1 AND product_id=$2 LIMIT 1;`, userID, productID).Scan(&cartItem.ID, &cartItem.UserID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CreatedAt, &cartItem.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCartItemNotFound
		}
		return nil, err
	}
	return &cartItem, nil
}

func (r *Repo) AddToCart(ctx context.Context, userID int, productID int) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO cart_items(user_id, product_id, quantity) VALUES($1, $2, $3);`, userID, productID, 1)
	return err
}

func (r *Repo) UpdateCartItemQuantity(ctx context.Context, userID int, productID int, quantity int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var productQuantityLeft int
	err = r.db.QueryRowContext(ctx, `SELECT quantity_left FROM products WHERE id=$1 FOR UPDATE;`, productID).Scan(&productQuantityLeft)
	if err != nil {
		return err
	}
	if quantity > productQuantityLeft {
		return ErrCartItemQuantityExceeded
	}
	_, err = r.db.ExecContext(ctx, `UPDATE cart_items SET quantity=$1 WHERE user_id=$2 AND product_id=$3;`, quantity, userID, productID)
	return err
}

func (r *Repo) DiscardCart(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM cart_items WHERE user_id=$1;`, userID)
	return err
}

func (r *Repo) GetCartWithProducts(ctx context.Context, userID int) ([]CartItem, []Product, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	query := `
		SELECT 
			ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at, ci.updated_at,
			p.id, p.name, p.price, p.quantity_left, p.created_at, p.updated_at
		FROM cart_items ci
		LEFT JOIN products p ON p.id = ci.product_id
		WHERE ci.user_id = $1;`

	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get cart: %w", err)
	}
	defer rows.Close()

	cartItems := make([]CartItem, 0)
	products := make([]Product, 0)

	for rows.Next() {
		var cartItem CartItem
		var product Product
		err := rows.Scan(
			&cartItem.ID,
			&cartItem.UserID,
			&cartItem.ProductID,
			&cartItem.Quantity,
			&cartItem.CreatedAt,
			&cartItem.UpdatedAt,
			&product.ID,
			&product.Name,
			&product.Price,
			&product.QuantityLeft,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan cart item: %w", err)
		}
		cartItems = append(cartItems, cartItem)
		products = append(products, product)
	}

	return cartItems, products, nil
}

func (r *Repo) DeleteCartItem(ctx context.Context, userID int, productID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM cart_items WHERE user_id=$1 AND product_id=$2;`, userID, productID)
	return err
}
