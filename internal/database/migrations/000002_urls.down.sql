--Delete trigger
DROP TRIGGER IF EXISTS set_url_update_at ON urls;
--Delete trigger function
DROP FUNCTION IF EXISTS update_url_updated_at();
--Drop table
DROP TABLE IF EXISTS urls CASCADE;
