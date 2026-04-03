--Delete trigger
DELETE TRIGGER IF EXISTS set_user_update_at ON users;
--Delete trigger function
DELETE FUNCTION IF EXISTS update_user_updated_at;
DELETE FUNCTION IF EXISTS random_string;
--Drop table
DROP TABLE IF EXISTS users;
