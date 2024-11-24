package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidCoupon = errors.New("invalid coupon")
)

type Order struct {
	ID               int    `json:"id"`
	UserID           int    `json:"userId"`
	Status           string `json:"status"`
	TotalAmount      int    `json:"totalAmount"`
	DiscountedAmount int    `json:"discountedAmount"`
	CouponID         int    `json:"couponId,omitempty"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
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

func (r *Repo) CreateOrder(ctx context.Context, cart []CartItem, userID int, preTotalAmount int, validCoupon *Coupon) (*Order, error) {
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

	// First check if all products have enough quantity
	for _, item := range orderItems {
		var quantityLeft int
		err = tx.QueryRowContext(ctx,
			`SELECT quantity_left FROM products WHERE id = $1 FOR UPDATE`,
			item.ProductID,
		).Scan(&quantityLeft)

		if err != nil {
			return nil, fmt.Errorf("failed to check product quantity: %w", err)
		}

		if quantityLeft < item.Quantity {
			return nil, fmt.Errorf("insufficient quantity for product %d", item.ProductID)
		}
	}

	totalAmount := preTotalAmount
	if validCoupon != nil {
		// Check if coupon is already used
		var isUsed bool
		err = tx.QueryRowContext(ctx,
			`SELECT is_used FROM coupons WHERE id = $1 FOR UPDATE`,
			validCoupon.ID,
		).Scan(&isUsed)

		if err != nil {
			return nil, fmt.Errorf("failed to check coupon: %w", err)
		}

		if isUsed {
			return nil, errors.New("coupon already used")
		}

		totalAmount -= totalAmount * validCoupon.DiscountPercent / 100

		row := tx.QueryRowContext(ctx,
			`INSERT INTO orders(user_id, status, total_amount, discounted_amount, coupon_id) 
			 VALUES($1, $2, $3, $4, $5) 
			 RETURNING id, user_id, status, total_amount, discounted_amount, coupon_id, created_at, updated_at`,
			userID, "completed", totalAmount, preTotalAmount-totalAmount, validCoupon.ID)

		err = row.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalAmount,
			&order.DiscountedAmount,
			&order.CouponID,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create order: %w", err)
		}

		// Mark coupon as used
		_, err = tx.ExecContext(ctx,
			`UPDATE coupons SET is_used = TRUE WHERE id = $1`,
			validCoupon.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update coupon: %w", err)
		}
	} else {
		row := tx.QueryRowContext(ctx,
			`INSERT INTO orders(user_id, status, total_amount, discounted_amount) 
			 VALUES($1, $2, $3, $4) 
			 RETURNING id, user_id, status, total_amount, discounted_amount, created_at, updated_at`,
			userID, "completed", totalAmount, preTotalAmount-totalAmount)

		err = row.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalAmount,
			&order.DiscountedAmount,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create order: %w", err)
		}
	}

	// Insert order items
	for _, orderItem := range orderItems {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO order_items(order_id, product_id, quantity) 
			 VALUES($1, $2, $3)`,
			order.ID, orderItem.ProductID, orderItem.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// Update product quantities
	for _, item := range orderItems {
		_, err = tx.ExecContext(ctx,
			`UPDATE products 
			 SET quantity_left = quantity_left - $1 
			 WHERE id = $2`,
			item.Quantity, item.ProductID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update product quantity: %w", err)
		}
	}

	return &order, nil
}

func (r *Repo) GetAllOrders(ctx context.Context, page int, pageSize int) ([]Order, error) {
	var orders []Order
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, status, total_amount, discounted_amount, created_at, updated_at FROM orders ORDER BY created_at DESC LIMIT $1 OFFSET $2;`, pageSize, page*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		err = rows.Scan(&order.ID, &order.UserID, &order.Status, &order.TotalAmount, &order.DiscountedAmount, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *Repo) GetOrdersCountForUser(ctx context.Context, userID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM orders WHERE user_id=$1;`, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
