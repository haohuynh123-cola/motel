-- Create Houses Table
CREATE TABLE IF NOT EXISTS houses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create Rooms Table
CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    house_id INTEGER NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    area DECIMAL(10, 2) NOT NULL,
    price DECIMAL(15, 2) NOT NULL,
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create Indexes for performance
CREATE INDEX IF NOT EXISTS idx_rooms_house_id ON rooms(house_id);
CREATE INDEX IF NOT EXISTS idx_rooms_is_available ON rooms(is_available);
