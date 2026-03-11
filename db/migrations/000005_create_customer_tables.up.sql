-- Bảng Khách thuê (Customers)
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    identity_number VARCHAR(20) UNIQUE NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    address TEXT,
    birthday DATE,
    gender VARCHAR(10),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Bảng Hợp đồng (Contracts)
CREATE TABLE IF NOT EXISTS contracts (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id),
    room_id INTEGER REFERENCES rooms(id),
    start_date DATE NOT NULL,
    end_date DATE,
    deposit DECIMAL(15, 2) DEFAULT 0,
    monthly_rent DECIMAL(15, 2) NOT NULL,
    payment_day INTEGER CHECK (payment_day >= 1 AND payment_day <= 31),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Cấu hình điện nước (Utility Configs) cho Nhà trọ
CREATE TABLE IF NOT EXISTS utility_configs (
    id SERIAL PRIMARY KEY,
    house_id INTEGER UNIQUE REFERENCES houses(id),
    electricity_price DECIMAL(10, 2) NOT NULL,
    water_price DECIMAL(10, 2) NOT NULL,
    trash_price DECIMAL(10, 2) DEFAULT 0,
    internet_price DECIMAL(10, 2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Chốt số điện nước hàng tháng (Utility Usages)
CREATE TABLE IF NOT EXISTS utility_usages (
    id SERIAL PRIMARY KEY,
    room_id INTEGER REFERENCES rooms(id),
    month INTEGER NOT NULL,
    year INTEGER NOT NULL,
    electricity_begin DECIMAL(15, 2) NOT NULL,
    electricity_end DECIMAL(15, 2) NOT NULL,
    water_begin DECIMAL(15, 2) NOT NULL,
    water_end DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(room_id, month, year)
);
