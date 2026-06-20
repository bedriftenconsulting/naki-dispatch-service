CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS dispatch_logs (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id   UUID NOT NULL,
    nurse_id     UUID,
    status       VARCHAR(50) NOT NULL DEFAULT 'pending',
    reason       TEXT DEFAULT '',
    distance     DECIMAL(10,2) DEFAULT 0,
    match_score  DECIMAL(10,4) DEFAULT 0,
    booking_type VARCHAR(50) DEFAULT 'scheduled',
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_dispatch_booking_id ON dispatch_logs(booking_id);
CREATE INDEX IF NOT EXISTS idx_dispatch_nurse_id ON dispatch_logs(nurse_id);
CREATE INDEX IF NOT EXISTS idx_dispatch_status ON dispatch_logs(status);
