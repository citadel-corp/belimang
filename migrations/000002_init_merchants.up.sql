DROP TYPE IF EXISTS merchant_category;
CREATE TYPE merchant_category AS ENUM(
    'SmallRestaurant', 
    'MediumRestaurant', 
    'LargeRestaurant', 
    'MerchandiseRestaurant',
    'BoothKiosk',
    'ConvenienceStore'
);

CREATE TABLE IF NOT EXISTS
merchants (
    id SERIAL PRIMARY KEY,
    uid CHAR(16) UNIQUE,
    name VARCHAR(30) NOT NULL,
    merchant_category merchant_category NOT NULL,
    image_url VARCHAR NOT NULL,
    location_lat FLOAT NOT NULL,
    location_lng FLOAT NOT NULL,
    --geom GEOMETRY(Point, 4326)
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS merchants_uid
	ON merchants USING HASH(uid);
CREATE INDEX IF NOT EXISTS merchants_category_idx
	ON merchants USING HASH (merchant_category);
CREATE INDEX IF NOT EXISTS merchants_name
	ON merchants USING HASH (name);
CREATE INDEX IF NOT EXISTS merchants_created_at_desc
	ON merchants(created_at DESC);
CREATE INDEX IF NOT EXISTS merchants_created_at_asc
	ON merchants(created_at ASC);
CREATE INDEX IF NOT EXISTS merchants_lat_lng 
	ON merchants (location_lat, location_lng);
    