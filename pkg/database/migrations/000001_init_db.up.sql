CREATE TYPE seats_status AS ENUM ('AVAILABLE', 'RESERVED', 'SOLD');

CREATE TYPE bookings_status AS ENUM ('PENDING', 'PAID', 'CANCELLED', 'FAILED');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50),
    hash_password TEXT
);

CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    event_date TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE seats (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT REFERENCES events(id),
    seat_number VARCHAR(10),
    price DECIMAL(10, 2),
    status seats_status DEFAULT 'AVAILABLE',
    version INT DEFAULT 1
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    event_id BIGINT REFERENCES events(id),
    total_amount DECIMAL(10, 2),
    status bookings_status DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMP
);

CREATE TABLE booking_items (
    booking_id UUID REFERENCES bookings(id),
    seat_id BIGINT REFERENCES seats(id)
);
