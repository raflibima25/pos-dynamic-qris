-- Create payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    method VARCHAR(50) NOT NULL CHECK (method IN ('qris')),
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'success', 'failed', 'expired', 'cancelled')),
    external_id VARCHAR(255), -- Midtrans transaction ID
    external_response TEXT,   -- Midtrans response JSON
    paid_at TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create qris_codes table
CREATE TABLE qris_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    qr_code TEXT NOT NULL,
    qr_image TEXT, -- Base64 encoded image
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_payments_transaction_id ON payments(transaction_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_external_id ON payments(external_id);
CREATE INDEX idx_qris_codes_transaction_id ON qris_codes(transaction_id);
CREATE INDEX idx_qris_codes_payment_id ON qris_codes(payment_id);

-- Create triggers to update updated_at timestamp
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();