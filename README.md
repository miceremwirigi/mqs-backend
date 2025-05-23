# mqs
MQS is a medical Equipment Servicing management app that helps service engineers schedule upcoming tasks, and notifies them when they are due.

git clone git@github.com:miceremwirigi/mqs-backend.git
go mod init
go mod tidy

# Configure database
Install postgresql, then if on linux, do -> sudo -u postgres psql 
CREATE DATABASE mqs_db;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE USER mqs_user WITH SUPERUSER CREATEDB CREATEROLE LOGIN PASSWORD 'mqs_pass';
ALTER ROLE mqs_user SET client_encoding TO 'utf8';
ALTER ROLE mqs_user SET default_transaction_isolation TO 'read committed';
ALTER ROLE mqs_user SET timezone TO 'UTC';
GRANT ALL PRIVILEGES ON DATABASE mqs_db TO mqs_user;
CREATE DATABASE mqs_test;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
GRANT ALL PRIVILEGES ON DATABASE mqs_test TO mqs_user;
\q

# Load env
source env.sh

