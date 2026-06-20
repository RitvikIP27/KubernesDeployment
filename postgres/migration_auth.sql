-- SaaS auth migration for SkillPulse

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT,
    oauth_provider TEXT,
    oauth_id TEXT,
    name TEXT,
    avatar_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE IF EXISTS skills
ADD COLUMN IF NOT EXISTS user_id INT NOT NULL DEFAULT 1 REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE IF EXISTS learning_logs
ADD COLUMN IF NOT EXISTS user_id INT NOT NULL DEFAULT 1 REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE IF EXISTS settings
ADD COLUMN IF NOT EXISTS user_id INT NOT NULL DEFAULT 1 REFERENCES users(id) ON DELETE CASCADE;

INSERT INTO users (email, password_hash, name)
VALUES ('admin@example.com', '$2b$12$XAxdZVaBMFdWah3WnS92l.hSsl9Nion996q8YX4eAx2ogiUOfH1be', 'Admin User')
ON CONFLICT (email) DO NOTHING;
