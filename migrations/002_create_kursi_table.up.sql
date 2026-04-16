-- File: migrations/002_create_kursi_table.up.sql
CREATE TABLE kursi (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    baris INTEGER NOT NULL,
    nomor_kursi INTEGER NOT NULL,
    tribun_id UUID NOT NULL REFERENCES tribun(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(baris, nomor_kursi, tribun_id)
);

CREATE INDEX idx_kursi_tribun_id ON kursi(tribun_id);
CREATE INDEX idx_kursi_baris_nomor ON kursi(baris, nomor_kursi);