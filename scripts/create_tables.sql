-- Table for storing device tokens (for FCM push notifications)
CREATE TABLE IF NOT EXISTS device_tokens (
    id SERIAL PRIMARY KEY,
    car_id VARCHAR(50) NOT NULL,
    fcm_token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
