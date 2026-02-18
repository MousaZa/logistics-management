-- Create the Orders Table
CREATE TABLE orders (
    order_uuid UUID PRIMARY KEY,

    -- Basic Info
    placed_by TEXT NOT NULL,
    destination TEXT NOT NULL,

    -- Money & Weight
    -- Using NUMERIC for money is safer than float, even if Go uses float32
    order_total NUMERIC(10, 2) NOT NULL,
    weight REAL NOT NULL,

    -- Status with Validation
    -- This ensures only your specific Go constants can be saved
    status TEXT NOT NULL CHECK (status IN ('cancelled', 'completed', 'confirmed', 'delivered', 'pending', 'shipped')),

    -- The Line Items (Value Objects)
    -- Storing this as JSONB is efficient and maps directly to your []LineItem slice
    line_items JSONB NOT NULL DEFAULT '[]'::jsonb,

    -- Timestamps
    ordered_date TIMESTAMPTZ NOT NULL,
    shipped_date TIMESTAMPTZ,   -- Nullable because it might no`t be shipped yet
    delivered_date TIMESTAMPTZ, -- Nullable
    completed_date TIMESTAMPTZ  -- Nullable
);

-- Create the Outbox Table (For the Transactional Outbox Pattern)
-- CREATE TABLE events (
--     uuid UUID PRIMARY KEY,
--     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     payload JSONB NOT NULL,
--     metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
--     topic TEXT NOT NULL
-- );



---- create above / drop below ----

DROP TABLE events;
DROP TABLE orders;