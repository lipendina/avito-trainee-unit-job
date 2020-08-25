CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS balance (id UUID DEFAULT uuid_generate_v4() PRIMARY KEY, user_id UUID NOT NULL, amount BIGINT NOT NULL CHECK (amount >= 0), UNIQUE(user_id));
CREATE TABLE IF NOT EXISTS "transaction" (id UUID DEFAULT uuid_generate_v4() PRIMARY KEY, user_id UUID REFERENCES balance(user_id) NOT NULL, change_balance BIGINT NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE INDEX balance_user_id_idx ON balance (user_id);