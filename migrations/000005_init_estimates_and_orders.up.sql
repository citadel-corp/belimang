CREATE TABLE IF NOT EXISTS
calculated_estimates (
    id VARCHAR(16) PRIMARY KEY,
    user_id CHAR(16) NOT NULL,
    total_price INT NOT NULL,
    user_location_lat float NOT NULL,
    user_location_lng float NOT NULL,
    estimated_delivery_time INT NOT NULL,
    ordered BOOLEAN NOT NULL DEFAULT FALSE,
    merchants JSONB NOT NULL, -- array of merchant_id
    items JSONB NOT NULL, -- array of object, object of item_id and quantity
    created_at TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE calculated_estimates ADD CONSTRAINT fk_estimate_user_id
    FOREIGN KEY (user_id)
			REFERENCES users(uid)
			ON DELETE CASCADE
			ON UPDATE NO ACTION;

CREATE INDEX IF NOT EXISTS calculated_estimates_created_at_desc
	ON calculated_estimates (created_at DESC);
CREATE INDEX IF NOT EXISTS calculated_estimates_items_created_at_asc
	ON calculated_estimates (created_at ASC);

CREATE TABLE IF NOT EXISTS
orders (
    id CHAR(16) PRIMARY KEY,
    calculated_estimate_id CHAR(16) NOT NULL,
    user_id CHAR(16) NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE orders ADD CONSTRAINT fk_orders_calculated_estimate_id
    FOREIGN KEY (calculated_estimate_id)
			REFERENCES calculated_estimates(id)
			ON DELETE CASCADE
			ON UPDATE NO ACTION;

ALTER TABLE orders ADD CONSTRAINT fk_orders_user_id
	FOREIGN KEY (user_id)
		REFERENCES users(uid)
		ON DELETE CASCADE
		ON UPDATE NO ACTION;

ALTER TABLE orders ADD CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
    REFERENCES users(uid)
    ON DELETE CASCADE
    ON UPDATE NO ACTION;

ALTER TABLE orders ADD CONSTRAINT fk_calculated_estimate_id
    FOREIGN KEY (calculated_estimate_id)
    REFERENCES calculated_estimates(id)
    ON DELETE CASCADE
    ON UPDATE NO ACTION;

CREATE INDEX IF NOT EXISTS orders_calculated_estimate_id
	ON orders USING HASH(calculated_estimate_id);
CREATE INDEX IF NOT EXISTS orders_user_id
	ON orders USING HASH(user_id);
CREATE INDEX IF NOT EXISTS orders_created_at_desc
	ON orders (created_at DESC);
CREATE INDEX IF NOT EXISTS orders_created_at_asc
	ON orders (created_at ASC);

CREATE TABLE IF NOT EXISTS
order_items(
    id CHAR(16) PRIMARY KEY,
    order_id CHAR(16) NOT NULL,
    merchant_id CHAR(16) NOT NULL,
    items JSONB NOT NULL, -- array of object, object of item_id and quantity
    created_at TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE order_items ADD CONSTRAINT fk_order_items_order_id
    FOREIGN KEY (order_id)
			REFERENCES orders(id)
			ON DELETE CASCADE
			ON UPDATE NO ACTION;

CREATE INDEX IF NOT EXISTS order_items_order_id
	ON order_items USING HASH(order_id);
CREATE INDEX IF NOT EXISTS order_items_created_at_desc
	ON order_items(created_at DESC);
CREATE INDEX IF NOT EXISTS order_items_created_at_asc
	ON order_items(created_at ASC);