-- DDL
-- Users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    deposit_amount DECIMAL(12, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT users_deposit_amount_check
        CHECK (deposit_amount >= 0)
);

-- Memberships
CREATE TABLE memberships (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    costs DECIMAL(12, 2) NOT NULL,
    duration_days INT NOT NULL,
    availability BOOLEAN NOT NULL DEFAULT TRUE,
    category VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT memberships_costs_check
        CHECK (costs >= 0),
    CONSTRAINT memberships_duration_days_check
        CHECK (duration_days > 0)
);

-- Trainers
CREATE TABLE trainers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    specialization VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Gym class
CREATE TABLE gym_classes (
    id SERIAL PRIMARY KEY,
    trainer_id INT NOT NULL,
    class_name VARCHAR(255) NOT NULL,
    schedule TIMESTAMP NOT NULL,
    quota INT NOT NULL,
    availability BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT gym_classes_trainer_fk 
        FOREIGN KEY (trainer_id)
        REFERENCES trainers(id)
        ON DELETE RESTRICT,

    CONSTRAINT gym_classes_quota_check
        CHECK (quota > 0)
);

-- Transaction
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    membership_id INT NOT NULL,
    total_amount NUMERIC(12, 2) NOT NULL,
    deposit_payment NUMERIC(12, 2) NOT NULL DEFAULT 0,
    status VARCHAR(255) NOT NULL DEFAULT 'pending',
    transaction_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT transactions_user_fk 
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT,

    CONSTRAINT transactions_membership_fk 
        FOREIGN KEY (membership_id)
        REFERENCES memberships(id)
        ON DELETE RESTRICT,

    CONSTRAINT transactions_total_amount_check
        CHECK (total_amount >= 0),

    CONSTRAINT transactions_deposit_payment_check
        CHECK (deposit_payment >= 0),

    CONSTRAINT transactions_status_check
        CHECK (status IN ('pending', 'paid', 'failed', 'cancelled'))
);

-- Booking
CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    gym_class_id INT NOT NULL,
    booking_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(255) NOT NULL DEFAULT 'booked',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT bookings_user_fk 
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT,

    CONSTRAINT bookings_gym_class_fk 
        FOREIGN KEY (gym_class_id)
        REFERENCES gym_classes(id)
        ON DELETE RESTRICT,

    CONSTRAINT bookings_status_check
        CHECK (status IN ('booked', 'cancelled', 'attended')),

    CONSTRAINT bookings_user_class_unique
        UNIQUE (user_id, gym_class_id)
);