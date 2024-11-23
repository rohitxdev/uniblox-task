package repo

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type Order struct {
	ID          int    `json:"id"`
	UserID      int    `json:"userId"`
	Status      string `json:"status"`
	TotalAmount int    `json:"totalAmount"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type OrderItem struct {
	ID        int    `json:"id"`
	OrderID   int    `json:"orderId"`
	ProductID int    `json:"productId"`
	Quantity  int    `json:"quantity"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (r *Repo) GetOrders(ctx context.Context) ([]Order, error) {
	var orders []Order
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, status, total_amount, created_at, updated_at FROM orders;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		err = rows.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *Repo) GetOrder(ctx context.Context, id int) (*Order, error) {
	var order Order
	err := r.db.QueryRowContext(ctx, `SELECT id, user_id, status, total_amount, created_at, updated_at FROM orders WHERE id=$1 LIMIT 1;`, id).Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return nil, nil
}

func (r *Repo) CreateOrder(ctx context.Context, cart []CartItem, userID int, totalAmount int) (*Order, error) {
	var order Order
	var orderItems []OrderItem

	for _, cartItem := range cart {
		orderItem := OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
		}
		orderItems = append(orderItems, orderItem)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	row := tx.QueryRowContext(ctx, `INSERT INTO orders(user_id, status, total_amount) VALUES($1, $2, $3) RETURNING (id, user_id, status, total_amount, created_at, updated_at);`, userID, "completed", totalAmount)
	if err = row.Scan(&order.ID); err != nil {
		return nil, err
	}

	for _, orderItem := range orderItems {
		_, err = tx.ExecContext(ctx, `INSERT INTO order_items(order_id, product_id, quantity) VALUES($1, $2, $3);`, order.ID, orderItem.ProductID, orderItem.Quantity)
		if err != nil {
			return nil, err
		}
	}

	return &order, nil
}
