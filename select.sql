-- Fetch tables.
SELECT * FROM information_schema.tables
WHERE table_catalog = $1;

-- Fetch columns.
SELECT * FROM information_schema.columns
WHERE table_name = $1;