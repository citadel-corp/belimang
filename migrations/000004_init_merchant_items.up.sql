DROP TYPE IF EXISTS item_category;
CREATE TYPE item_category AS ENUM(
'Beverage', 
'Food', 
'Snack', 
'Condiments',
'Additions');

CREATE TABLE IF NOT EXISTS
merchant_items (
    id SERIAL PRIMARY KEY,
    uid CHAR(16) NOT NULL,
    merchant_id BIGINT NOT NULL,
    name VARCHAR(30) NOT NULL,
    item_category item_category NOT NULL,
    price INT NOT NULL,
    image_url VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE merchant_items ADD CONSTRAINT fk_merchant_id
    FOREIGN KEY (merchant_id)
    REFERENCES merchants(id)
    ON DELETE CASCADE
    ON UPDATE NO ACTION;

CREATE INDEX IF NOT EXISTS merchant_items_uid
	ON merchant_items USING HASH(uid);
CREATE INDEX IF NOT EXISTS merchant_items_merchant_id
	ON merchant_items USING HASH(merchant_id);
CREATE INDEX IF NOT EXISTS merchant_items_name
	ON merchant_items USING HASH(lower(name));
CREATE INDEX IF NOT EXISTS merchant_items_category
	ON merchant_items USING HASH(item_category);
CREATE INDEX IF NOT EXISTS merchant_items_created_at_desc
	ON merchant_items(created_at DESC);
CREATE INDEX IF NOT EXISTS merchant_items_created_at_asc
	ON merchant_items(created_at ASC);
CREATE INDEX IF NOT EXISTS merchant_items_price_desc
	ON merchant_items(price DESC);
CREATE INDEX IF NOT EXISTS merchant_items_price_asc
	ON merchant_items(price ASC);
