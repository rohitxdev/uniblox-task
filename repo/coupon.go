package repo

import (
	"context"
	"errors"
)

var (
	ErrCouponNotFound = errors.New("coupon not found")
)

type Coupon struct {
	ID              int    `json:"id"`
	UserID          int    `json:"userId"`
	Code            string `json:"code"`
	DiscountPercent int    `json:"discountPercent"`
	IsUsed          bool   `json:"isUsed"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

func (r *Repo) CreateCoupon(ctx context.Context, userID int, code string, discountPercent int) (*Coupon, error) {
	var coupon Coupon
	err := r.db.QueryRowContext(ctx, `INSERT INTO coupons(user_id, code, discount_percent) 
        VALUES($1, $2, $3) 
        RETURNING id, user_id, code, discount_percent, created_at, updated_at;`,
		userID, code, discountPercent).Scan(&coupon.ID, &coupon.UserID, &coupon.Code, &coupon.DiscountPercent, &coupon.CreatedAt, &coupon.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *Repo) GetAllCoupons(ctx context.Context, page int, pageSize int) ([]Coupon, error) {
	var coupons []Coupon
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, code, discount_percent, is_used, created_at, updated_at FROM coupons ORDER BY created_at DESC LIMIT $1 OFFSET $2;`, pageSize, page*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var coupon Coupon
		err = rows.Scan(&coupon.ID, &coupon.UserID, &coupon.Code, &coupon.DiscountPercent, &coupon.IsUsed, &coupon.CreatedAt, &coupon.UpdatedAt)
		if err != nil {
			return nil, err
		}
		coupons = append(coupons, coupon)
	}
	return coupons, nil
}

func (r *Repo) GetAvailableCoupons(ctx context.Context, userID int) ([]Coupon, error) {
	var coupons []Coupon
	rows, err := r.db.QueryContext(ctx, `SELECT id, code, discount_percent, is_used, created_at, updated_at FROM coupons WHERE user_id=$1 AND is_used=FALSE;`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var coupon Coupon
		err = rows.Scan(&coupon.ID, &coupon.Code, &coupon.DiscountPercent, &coupon.IsUsed, &coupon.CreatedAt, &coupon.UpdatedAt)
		if err != nil {
			return nil, err
		}
		coupon.UserID = userID
		coupons = append(coupons, coupon)
	}
	return coupons, nil
}
