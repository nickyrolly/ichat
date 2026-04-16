-- File: database/schema.sql
-- Database: ichat
-- Schema: public

-- Create database (run this manually in PostgreSQL)
-- CREATE DATABASE ichat;

-- Drop tables if they exist (for development)
DROP TABLE IF EXISTS kursi CASCADE;
DROP TABLE IF EXISTS tribun CASCADE;

-- Create tribun table
CREATE TABLE tribun (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create kursi table
CREATE TABLE kursi (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    baris INTEGER NOT NULL,
    nomor_kursi INTEGER NOT NULL,
    tribun_id UUID NOT NULL REFERENCES tribun(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(baris, nomor_kursi, tribun_id)
);

-- Create indexes for better performance
CREATE INDEX idx_kursi_tribun_id ON kursi(tribun_id);
CREATE INDEX idx_kursi_baris_nomor ON kursi(baris, nomor_kursi);

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_tribun_updated_at BEFORE UPDATE ON tribun
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_kursi_updated_at BEFORE UPDATE ON kursi
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert sample data
INSERT INTO tribun (nama) VALUES 
('Tribun Utara'),
('Tribun Selatan'),
('Tribun Timur');

-- Get tribun IDs for sample kursi data
INSERT INTO kursi (baris, nomor_kursi, tribun_id) 
SELECT 
    generate_series(1, 10) as baris,
    generate_series(1, 20) as nomor_kursi,
    id as tribun_id
FROM tribun;