CREATE TABLE users (
    user_id uuid PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    user_role VARCHAR(50) CHECK (user_role IN ('client', 'moderator')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tokens (
    token_id SERIAL PRIMARY KEY,
    user_id uuid NOT NULL,
    token VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE houses (
    house_id SERIAL PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    year INT CHECK (year > 0) NOT NULL,
    developer VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE flats (
    flat_id SERIAL PRIMARY KEY,
    -- flat_number VARCHAR(100) NOT NULL,
    price INT NOT NULL,
    rooms INT NOT NULL,
    house_id INT NOT NULL,
    moderation_status VARCHAR(20) CHECK (moderation_status IN ('created', 'approved', 'declined', 'on moderation')) DEFAULT 'created' NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (house_id) REFERENCES houses(house_id) ON DELETE CASCADE
);


CREATE TABLE subscribers (
    user_id uuid,
    house_id INT,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (house_id) REFERENCES houses(house_id),
    PRIMARY KEY(user_id, house_id)
)

BEGIN; 
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE READ ONLY;
COMMIT;

