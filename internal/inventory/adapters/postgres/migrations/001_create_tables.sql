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
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Inventory Table (Mapping Products to Locations with Quantity)
CREATE TABLE inventory (
    product_uuid UUID NOT NULL REFERENCES products(product_uuid) ON DELETE RESTRICT,
    location_uuid UUID NOT NULL REFERENCES locations(location_uuid) ON DELETE RESTRICT,

    quantity INT NOT NULL DEFAULT 0,
    status TEXT NOT NULL CHECK (status IN ('available', 'reserved', 'damaged')),

    -- Prevent negative inventory at the database level
    CONSTRAINT chk_quantity_non_negative CHECK (quantity >= 0),

    -- A product can only have one inventory record per location
    PRIMARY KEY (product_uuid, location_uuid)
);


---- create above / drop below ----

DROP TABLE inventory;
DROP TABLE locations;
DROP TABLE products;