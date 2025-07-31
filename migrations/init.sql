CREATE DATABASE subscriptions;
CREATE USER subscription WITH PASSWORD 'subscription';
ALTER DATABASE subscriptions OWNER TO subscription;
-- GRANT ALL PRIVILEGES ON DATABASE "subscriptions" to subscription;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";