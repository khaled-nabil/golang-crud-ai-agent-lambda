CREATE TABLE domain (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(70) NOT NULL,
    instructions text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
);

CREATE INDEX IF NOT EXISTS idx_domain_id ON domain (id);