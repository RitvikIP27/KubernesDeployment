-- SkillPulse PostgreSQL Database Schema

CREATE TABLE IF NOT EXISTS skills (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) DEFAULT '',
    target_hours INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS learning_logs (
    id SERIAL PRIMARY KEY,
    skill_id INT NOT NULL,
    hours DECIMAL(4,1) NOT NULL,
    notes TEXT,
    log_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_skill FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS settings (
    key VARCHAR(50) PRIMARY KEY,
    value TEXT NOT NULL
);

-- Seed skills data
INSERT INTO skills (id, name, category, target_hours) VALUES
    (1, 'Docker', 'DevOps', 40),
    (2, 'Kubernetes', 'DevOps', 60),
    (3, 'Go', 'Programming', 50),
    (4, 'Azure DevOps', 'Cloud', 30),
    (5, 'Terraform', 'DevOps', 35)
ON CONFLICT (id) DO NOTHING;

-- Reset serial sequence to start after the seeded values
SELECT setval('skills_id_seq', COALESCE((SELECT MAX(id)+1 FROM skills), 1), false);

-- Seed learning logs data
INSERT INTO learning_logs (skill_id, hours, notes, log_date) VALUES
    (1, 2.0, 'Learned Docker basics - images, containers, volumes', '2026-03-10'),
    (1, 1.5, 'Built multi-stage Dockerfile for Go app', '2026-03-12'),
    (1, 3.0, 'Docker Compose with multiple services', '2026-03-14'),
    (2, 1.0, 'Kubernetes architecture overview', '2026-03-11'),
    (2, 2.0, 'Deployed first pod and service', '2026-03-13'),
    (3, 2.5, 'Go basics - structs, interfaces, goroutines', '2026-03-10'),
    (3, 1.5, 'Built REST API with Gin framework', '2026-03-15'),
    (4, 1.0, 'Created Azure DevOps org and project', '2026-03-16'),
    (5, 1.5, 'Terraform basics - providers, resources, state', '2026-03-17')
ON CONFLICT DO NOTHING;

-- Seed default settings
INSERT INTO settings (key, value) VALUES 
    ('target_role', 'Site Reliability Engineer (SRE)')
ON CONFLICT (key) DO NOTHING;
