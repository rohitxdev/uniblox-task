-- Create case-insensitive text column
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    role TEXT NOT NULL CHECK (role IN ('user', 'admin')) DEFAULT 'user',
    email CITEXT NOT NULL UNIQUE CHECK (LENGTH(email) <= 64),
    full_name TEXT CHECK (LENGTH(full_name) <= 64),
    date_of_birth DATE CHECK (date_of_birth >= '1900-01-01'),
    gender TEXT CHECK (gender IN ('male', 'female', 'other')),
    phone_number TEXT CHECK (LENGTH(phone_number) <= 16),
    account_status TEXT NOT NULL CHECK (
        account_status IN ('active', 'suspended', 'deleted')
    ) DEFAULT 'active',
    image_url TEXT,
    is_verified BOOL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp
);

-- Create a function to update the updated_at column
CREATE
OR REPLACE FUNCTION set_updated_at_column () RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = current_timestamp; 
    RETURN NEW; 
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_user_updated_at BEFORE
UPDATE ON users FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column ();


CREATE TABLE products (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL CHECK (LENGTH(name) <= 128),
    price BIGINT NOT NULL CHECK (price > 0),
    quantity_left BIGINT NOT NULL CHECK (quantity_left >= 0),
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp
);

CREATE TRIGGER set_products_updated_at BEFORE
UPDATE ON products FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column ();


CREATE TABLE orders (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id),
    status TEXT NOT NULL CHECK (
        status IN ('pending', 'processing', 'completed', 'cancelled')
    ),
    total_amount BIGINT NOT NULL CHECK (total_amount > 0),
    discounted_amount BIGINT NOT NULL DEFAULT 0 CHECK (discounted_amount >= 0),
    coupon_id BIGINT,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp
);

CREATE TRIGGER set_orders_updated_at BEFORE
UPDATE ON products FOR EACH ROW EXECUTE FUNCTION set_updated_at_column ();

CREATE TABLE order_items (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders (id),
    product_id BIGINT NOT NULL REFERENCES products (id),
    quantity BIGINT NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp
    UNIQUE (order_id, product_id)
);

CREATE TRIGGER set_order_items_updated_at BEFORE
UPDATE ON order_items FOR EACH ROW EXECUTE FUNCTION set_updated_at_column ();


CREATE TABLE cart_items (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id),
    product_id BIGINT NOT NULL REFERENCES products (id),
    quantity BIGINT NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp
    UNIQUE (user_id, product_id)
);

CREATE TRIGGER set_cart_items_updated_at BEFORE
UPDATE ON cart_items FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column ();


CREATE TABLE coupons (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id),
    code TEXT NOT NULL CHECK (LENGTH(code) <= 128),
    discount_percent INT NOT NULL CHECK (
        discount_percent >= 0
        AND discount_percent <= 100
    ),
    is_used BOOL NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp
);

CREATE TRIGGER set_coupons_updated_at BEFORE
UPDATE ON coupons FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column ();