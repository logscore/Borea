-- Connect to the pg_borea database
\c pg_borea;

-- Create unique_users table
CREATE TABLE IF NOT EXISTS unique_users (
    userId UUID PRIMARY KEY UNIQUE,
    lastActivityTime TIMESTAMP -- Store last activity time as a timestamp with timezone
);

-- Create admin_users table
CREATE TABLE IF NOT EXISTS admin_users (
    id SERIAL PRIMARY KEY,              -- Auto-incrementing ID for each admin user
    username TEXT NOT NULL UNIQUE,      -- Username must be unique
    password_hash TEXT NOT NULL         -- Password hash for the admin user
);

-- Create sessions table with ordered columns
CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    session_id UUID NOT NULL,                -- Matches sessionId
    last_activity_time TIMESTAMP DEFAULT NOW(), -- Matches lastActivityTime
    user_id UUID DEFAULT NULL,               -- Matches userId
    session_duration INTEGER DEFAULT NULL,   -- Matches sessionDuration
    user_agent TEXT,                         -- Matches userAgent
    referrer TEXT,                           -- Matches referrer
    token TEXT DEFAULT NULL,                 -- Matches token
    start_time TIMESTAMP DEFAULT NOW(),      -- Matches startTime
    language TEXT                           -- Matches language
    -- FOREIGN KEY (user_id) REFERENCES unique_users(userId) ON DELETE SET NULL -- Reference to unique_users table
);

-- Grant privileges to the user 'borea'
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO borea;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO borea;
