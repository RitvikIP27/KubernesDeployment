-- SkillPulse PostgreSQL Database Schema

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT,
    oauth_provider TEXT,
    oauth_id TEXT,
    name TEXT,
    avatar_url TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS skills (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) DEFAULT '',
    target_hours INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS learning_logs (
    id SERIAL PRIMARY KEY,
    skill_id INT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    hours DECIMAL(4,1) NOT NULL,
    notes TEXT,
    log_date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS settings (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key VARCHAR(50) NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY (user_id, key)
);

-- Seed initial user and demo data
INSERT INTO users (email, password_hash, name)
VALUES ('admin@example.com', '$2b$12$XAxdZVaBMFdWah3WnS92l.hSsl9Nion996q8YX4eAx2ogiUOfH1be', 'Admin User')
ON CONFLICT (email) DO NOTHING;

-- Use the seeded admin user as owner for demo rows
INSERT INTO skills (user_id, id, name, category, target_hours) VALUES
    (1, 1, 'Docker', 'DevOps', 40),
    (1, 2, 'Kubernetes', 'DevOps', 60),
    (1, 3, 'Go', 'Programming', 50),
    (1, 4, 'Azure DevOps', 'Cloud', 30),
    (1, 5, 'Terraform', 'DevOps', 35)
ON CONFLICT (id) DO NOTHING;

SELECT setval('skills_id_seq', COALESCE((SELECT MAX(id)+1 FROM skills), 1), false);

INSERT INTO learning_logs (skill_id, user_id, hours, notes, log_date) VALUES
    (1, 1, 2.0, 'Learned Docker basics - images, containers, volumes', '2026-03-10'),
    (1, 1, 1.5, 'Built multi-stage Dockerfile for Go app', '2026-03-12'),
    (1, 1, 3.0, 'Docker Compose with multiple services', '2026-03-14'),
    (2, 1, 1.0, 'Kubernetes architecture overview', '2026-03-11'),
    (2, 1, 2.0, 'Deployed first pod and service', '2026-03-13'),
    (3, 1, 2.5, 'Go basics - structs, interfaces, goroutines', '2026-03-10'),
    (3, 1, 1.5, 'Built REST API with Gin framework', '2026-03-15'),
    (4, 1, 1.0, 'Created Azure DevOps org and project', '2026-03-16'),
    (5, 1, 1.5, 'Terraform basics - providers, resources, state', '2026-03-17')
ON CONFLICT DO NOTHING;

INSERT INTO settings (user_id, key, value) VALUES 
    (1, 'target_role', 'Site Reliability Engineer (SRE)')
ON CONFLICT (user_id, key) DO NOTHING;
