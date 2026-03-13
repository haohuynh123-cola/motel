CREATE TABLE IF NOT EXISTS appointments (
    id SERIAL PRIMARY KEY,
    room_id INTEGER NOT NULL REFERENCES rooms(id),
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(20) NOT NULL,
    appointment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    note TEXT,
    status VARCHAR(20) DEFAULT 'pending', -- pending, confirmed, cancelled
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_appointments_room_id ON appointments(room_id);
CREATE INDEX idx_appointments_customer_email ON appointments(customer_email);
