-- 1. Products Table
CREATE TABLE products (
    product_uuid UUID PRIMARY KEY,
    name TEXT NOT NULL,
    -- NUMERIC is critical for financial calculations, avoiding floating-point rounding errors
    price NUMERIC(10, 2) NOT NULL,
    weight REAL NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Locations Table
CREATE TABLE locations (
    location_uuid UUID PRIMARY KEY,
    name TEXT NOT NULL,
    city TEXT NOT NULL,
    address TEXT NOT NULL,

    coordinates GEOGRAPHY(Point, 4326) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Inventory Table (Mapping Products to Locations with Quantity)
CREATE TABLE inventory (
    product_uuid UUID NOT NULL REFERENCES products(product_uuid) ON DELETE RESTRICT,
    location_uuid UUID NOT NULL REFERENCES locations(location_uuid) ON DELETE RESTRICT,

    -- Statuses become columns
    qty_available INT NOT NULL DEFAULT 0,
    qty_reserved  INT NOT NULL DEFAULT 0,
    qty_damaged   INT NOT NULL DEFAULT 0,

    -- Constraints ensure no bucket ever goes negative
    CONSTRAINT chk_avail_non_negative CHECK (qty_available >= 0),
    CONSTRAINT chk_rsvd_non_negative  CHECK (qty_reserved >= 0),
    CONSTRAINT chk_dmgd_non_negative  CHECK (qty_damaged >= 0),

    -- The original primary key works perfectly here
    PRIMARY KEY (product_uuid, location_uuid)
);

CREATE INDEX idx_locations_coordinates ON locations USING GIST (coordinates);

---- create above / drop below ----

DROP TABLE inventory;
DROP TABLE locations;
DROP TABLE products;
DROP EXTENSION IF EXISTS postgis;