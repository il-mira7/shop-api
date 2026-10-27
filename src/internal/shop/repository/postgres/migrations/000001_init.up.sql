CREATE TABLE IF NOT EXISTS address (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    street VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image bytea NOT NULL
);

CREATE TABLE IF NOT EXISTS supplier (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    address_id UUID REFERENCES address(id) ON DELETE SET NULL,
    phone_number VARCHAR(20) 
);

CREATE TABLE IF NOT EXISTS client (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_name VARCHAR(100) NOT NULL,
    client_surname VARCHAR(100) NOT NULL,
    birthday DATE NOT NULL,
    gender VARCHAR(10) NOT NULL,
    registration_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    address_id UUID REFERENCES address(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS product (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    category VARCHAR(100) NOT NULL,
    price NUMERIC(15, 2) NOT NULL CHECK (price >= 0), 
    available_stock INT NOT NULL DEFAULT 0 CHECK (available_stock >= 0), 
    last_update_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    supplier_id UUID REFERENCES supplier(id) ON DELETE CASCADE,
    image_id  UUID REFERENCES images(id) ON DELETE SET NULL
);
