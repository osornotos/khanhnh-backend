CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX products_description_idx ON products USING GIN (description gin_trgm_ops);
CREATE INDEX products_name_idx ON products USING GIN (name gin_trgm_ops);
CREATE INDEX products_created_at_idx ON products (created_at);
CREATE INDEX products_category_id_idx ON products (category_id);
CREATE INDEX products_status_idx ON products (status);
