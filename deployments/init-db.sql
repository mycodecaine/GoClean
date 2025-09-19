-- Create database for Keycloak if it doesn't exist
CREATE DATABASE IF NOT EXISTS keycloak;

-- Create additional extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- You can add additional database initialization scripts here
-- For example, create additional users, schemas, etc.