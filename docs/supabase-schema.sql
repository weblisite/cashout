-- Cashout Database Schema
-- Supabase PostgreSQL Schema

-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create custom types
CREATE TYPE kyc_status AS ENUM ('pending', 'verified', 'rejected');
CREATE TYPE transaction_type AS ENUM ('p2p', 'cash_in', 'cash_out', 'business');
CREATE TYPE transaction_status AS ENUM ('pending', 'completed', 'failed');
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended');
CREATE TYPE agent_status AS ENUM ('active', 'inactive', 'suspended');
CREATE TYPE business_status AS ENUM ('active', 'inactive', 'suspended');

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone_number VARCHAR(15) UNIQUE NOT NULL,
    hashed_id VARCHAR(14) UNIQUE NOT NULL,
    kyc_status kyc_status DEFAULT 'pending',
    wallet_balance DECIMAL(15,2) DEFAULT 0.00,
    status user_status DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Agents table
CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone_number VARCHAR(15) UNIQUE NOT NULL,
    float_balance DECIMAL(15,2) DEFAULT 0.00,
    commission_balance DECIMAL(15,2) DEFAULT 0.00,
    status agent_status DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Businesses table
CREATE TABLE businesses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    wallet_balance DECIMAL(15,2) DEFAULT 0.00,
    status business_status DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    recipient_id UUID REFERENCES users(id),
    business_id UUID REFERENCES businesses(id),
    agent_id UUID REFERENCES agents(id),
    amount DECIMAL(15,2) NOT NULL,
    fee DECIMAL(10,4) NOT NULL,
    agent_commission DECIMAL(10,4) DEFAULT 0.00,
    platform_margin DECIMAL(10,4) DEFAULT 0.00,
    type transaction_type NOT NULL,
    status transaction_status DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Fee structure table (for reference)
CREATE TABLE fee_structure (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    min_amount DECIMAL(15,2) NOT NULL,
    max_amount DECIMAL(15,2) NOT NULL,
    p2p_fee DECIMAL(10,4) NOT NULL,
    cash_out_fee DECIMAL(10,4) NOT NULL,
    business_fee DECIMAL(10,4) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert fee structure data with rounded values
INSERT INTO fee_structure (min_amount, max_amount, p2p_fee, cash_out_fee, business_fee) VALUES
(50.00, 100.00, 8.0, 8.0, 8.0),
(101.00, 500.00, 22.0, 22.0, 22.0),
(501.00, 1000.00, 22.0, 22.0, 22.0),
(1001.00, 1500.00, 22.0, 22.0, 22.0),
(1501.00, 2500.00, 22.0, 22.0, 22.0),
(2501.00, 3500.00, 39.0, 39.0, 39.0),
(3501.00, 5000.00, 52.0, 52.0, 52.0),
(5001.00, 7500.00, 65.0, 65.0, 65.0),
(7501.00, 10000.00, 86.0, 86.0, 86.0),
(10001.00, 15000.00, 125.0, 125.0, 125.0),
(15001.00, 20000.00, 139.0, 139.0, 139.0),
(20001.00, 35000.00, 148.0, 148.0, 148.0),
(35001.00, 50000.00, 209.0, 209.0, 209.0),
(50001.00, 250000.00, 232.0, 232.0, 232.0),
(250001.00, 500000.00, 513.0, 513.0, 513.0),
(500001.00, 1000000.00, 1076.0, 1076.0, 1076.0);

-- Create indexes for performance
CREATE INDEX idx_users_phone_number ON users(phone_number);
CREATE INDEX idx_users_hashed_id ON users(hashed_id);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_recipient_id ON transactions(recipient_id);
CREATE INDEX idx_transactions_business_id ON transactions(business_id);
CREATE INDEX idx_transactions_agent_id ON transactions(agent_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_agents_phone_number ON agents(phone_number);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_agents_updated_at BEFORE UPDATE ON agents FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_businesses_updated_at BEFORE UPDATE ON businesses FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Row Level Security (RLS) Policies

-- Enable RLS on all tables
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE agents ENABLE ROW LEVEL SECURITY;
ALTER TABLE businesses ENABLE ROW LEVEL SECURITY;
ALTER TABLE transactions ENABLE ROW LEVEL SECURITY;
ALTER TABLE fee_structure ENABLE ROW LEVEL SECURITY;

-- Users can only see their own data
CREATE POLICY "Users can view own profile" ON users FOR SELECT USING (auth.uid()::text = id::text);
CREATE POLICY "Users can update own profile" ON users FOR UPDATE USING (auth.uid()::text = id::text);

-- Users can see their own transactions
CREATE POLICY "Users can view own transactions" ON transactions FOR SELECT USING (auth.uid()::text = user_id::text OR auth.uid()::text = recipient_id::text);

-- Agents can see their own data and transactions
CREATE POLICY "Agents can view own profile" ON agents FOR SELECT USING (auth.uid()::text = id::text);
CREATE POLICY "Agents can view own transactions" ON transactions FOR SELECT USING (auth.uid()::text = agent_id::text);

-- Businesses can see their own data and transactions
CREATE POLICY "Businesses can view own profile" ON businesses FOR SELECT USING (auth.uid()::text = id::text);
CREATE POLICY "Businesses can view own transactions" ON transactions FOR SELECT USING (auth.uid()::text = business_id::text);

-- Fee structure is public (read-only)
CREATE POLICY "Fee structure is public" ON fee_structure FOR SELECT USING (true);

-- Functions for common operations

-- Function to calculate fee based on amount and transaction type
CREATE OR REPLACE FUNCTION calculate_fee(amount DECIMAL, transaction_type transaction_type)
RETURNS DECIMAL AS $$
DECLARE
    fee DECIMAL;
BEGIN
    SELECT 
        CASE 
            WHEN transaction_type = 'p2p' THEN p2p_fee
            WHEN transaction_type = 'cash_out' THEN cash_out_fee
            WHEN transaction_type = 'business' THEN business_fee
            ELSE 0
        END INTO fee
    FROM fee_structure 
    WHERE amount >= min_amount AND amount <= max_amount
    LIMIT 1;
    
    RETURN COALESCE(fee, 0);
END;
$$ LANGUAGE plpgsql;

-- Function to generate hashed ID from phone number
CREATE OR REPLACE FUNCTION generate_hashed_id(phone VARCHAR)
RETURNS VARCHAR AS $$
DECLARE
    hashed VARCHAR;
    salt VARCHAR;
BEGIN
    -- Generate hash from phone number
    hashed := encode(sha256(phone::bytea), 'hex');
    -- Take last 10 digits
    hashed := substring(hashed from length(hashed) - 9);
    -- Generate 4-digit random salt
    salt := lpad(floor(random() * 10000)::text, 4, '0');
    -- Return hashed ID with salt
    RETURN hashed || '-' || salt;
END;
$$ LANGUAGE plpgsql;

-- Function to process P2P transaction
CREATE OR REPLACE FUNCTION process_p2p_transaction(
    sender_id UUID,
    recipient_id UUID,
    amount DECIMAL
)
RETURNS UUID AS $$
DECLARE
    transaction_id UUID;
    fee DECIMAL;
BEGIN
    -- Calculate fee
    fee := calculate_fee(amount, 'p2p');
    
    -- Check sender balance
    IF (SELECT wallet_balance FROM users WHERE id = sender_id) < (amount + fee) THEN
        RAISE EXCEPTION 'Insufficient balance';
    END IF;
    
    -- Create transaction
    INSERT INTO transactions (user_id, recipient_id, amount, fee, platform_margin, type, status)
    VALUES (sender_id, recipient_id, amount, fee, fee, 'p2p', 'completed')
    RETURNING id INTO transaction_id;
    
    -- Update balances
    UPDATE users SET wallet_balance = wallet_balance - (amount + fee) WHERE id = sender_id;
    UPDATE users SET wallet_balance = wallet_balance + amount WHERE id = recipient_id;
    
    RETURN transaction_id;
END;
$$ LANGUAGE plpgsql;

-- Function to process cash-in transaction
CREATE OR REPLACE FUNCTION process_cash_in_transaction(
    user_id UUID,
    agent_id UUID,
    amount DECIMAL
)
RETURNS UUID AS $$
DECLARE
    transaction_id UUID;
BEGIN
    -- Check agent float balance
    IF (SELECT float_balance FROM agents WHERE id = agent_id) < amount THEN
        RAISE EXCEPTION 'Insufficient agent float balance';
    END IF;
    
    -- Create transaction (no fee for cash-in)
    INSERT INTO transactions (user_id, agent_id, amount, fee, type, status)
    VALUES (user_id, agent_id, amount, 0, 'cash_in', 'completed')
    RETURNING id INTO transaction_id;
    
    -- Update balances
    UPDATE users SET wallet_balance = wallet_balance + amount WHERE id = user_id;
    UPDATE agents SET float_balance = float_balance - amount WHERE id = agent_id;
    
    RETURN transaction_id;
END;
$$ LANGUAGE plpgsql;

-- Function to process cash-out transaction
CREATE OR REPLACE FUNCTION process_cash_out_transaction(
    user_id UUID,
    agent_id UUID,
    amount DECIMAL
)
RETURNS UUID AS $$
DECLARE
    transaction_id UUID;
    fee DECIMAL;
    agent_commission DECIMAL;
    platform_margin DECIMAL;
BEGIN
    -- Calculate fee
    fee := calculate_fee(amount, 'cash_out');
    agent_commission := fee * 0.25; -- 25% agent commission
    platform_margin := fee * 0.75;  -- 75% platform margin
    
    -- Check user balance
    IF (SELECT wallet_balance FROM users WHERE id = user_id) < (amount + fee) THEN
        RAISE EXCEPTION 'Insufficient user balance';
    END IF;
    
    -- Create transaction
    INSERT INTO transactions (user_id, agent_id, amount, fee, agent_commission, platform_margin, type, status)
    VALUES (user_id, agent_id, amount, fee, agent_commission, platform_margin, 'cash_out', 'completed')
    RETURNING id INTO transaction_id;
    
    -- Update balances
    UPDATE users SET wallet_balance = wallet_balance - (amount + fee) WHERE id = user_id;
    UPDATE agents SET 
        float_balance = float_balance + (amount + fee),
        commission_balance = commission_balance + agent_commission
    WHERE id = agent_id;
    
    RETURN transaction_id;
END;
$$ LANGUAGE plpgsql;

-- Function to process business transaction
CREATE OR REPLACE FUNCTION process_business_transaction(
    user_id UUID,
    business_id UUID,
    amount DECIMAL,
    is_user_initiated BOOLEAN
)
RETURNS UUID AS $$
DECLARE
    transaction_id UUID;
    fee DECIMAL;
    user_fee DECIMAL;
    business_fee DECIMAL;
BEGIN
    -- Calculate fee
    fee := calculate_fee(amount, 'business');
    user_fee := fee * 0.5;    -- 50% user fee
    business_fee := fee * 0.5; -- 50% business fee
    
    -- Check balances based on who initiates
    IF is_user_initiated THEN
        -- User pays business
        IF (SELECT wallet_balance FROM users WHERE id = user_id) < (amount + user_fee) THEN
            RAISE EXCEPTION 'Insufficient user balance';
        END IF;
        
        -- Create transaction
        INSERT INTO transactions (user_id, business_id, amount, fee, platform_margin, type, status)
        VALUES (user_id, business_id, amount, fee, fee, 'business', 'completed')
        RETURNING id INTO transaction_id;
        
        -- Update balances
        UPDATE users SET wallet_balance = wallet_balance - (amount + user_fee) WHERE id = user_id;
        UPDATE businesses SET wallet_balance = wallet_balance + (amount - business_fee) WHERE id = business_id;
    ELSE
        -- Business pays user
        IF (SELECT wallet_balance FROM businesses WHERE id = business_id) < (amount + business_fee) THEN
            RAISE EXCEPTION 'Insufficient business balance';
        END IF;
        
        -- Create transaction
        INSERT INTO transactions (user_id, business_id, amount, fee, platform_margin, type, status)
        VALUES (user_id, business_id, amount, fee, fee, 'business', 'completed')
        RETURNING id INTO transaction_id;
        
        -- Update balances
        UPDATE businesses SET wallet_balance = wallet_balance - (amount + business_fee) WHERE id = business_id;
        UPDATE users SET wallet_balance = wallet_balance + (amount - user_fee) WHERE id = user_id;
    END IF;
    
    RETURN transaction_id;
END;
$$ LANGUAGE plpgsql; 