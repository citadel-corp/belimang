CREATE INDEX IF NOT EXISTS merchants_name_lower
	ON merchants USING HASH (LOWER(name));
