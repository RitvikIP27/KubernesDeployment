-- HelixaCore Phase 2-5 Schema Migration
-- Phase 2: Analytics Dashboard
-- Phase 3: Career Readiness Score
-- Phase 4: Job Matching
-- Phase 5: Resume & Certificates

-- Phase 3: Career Readiness Score Tables
CREATE TABLE IF NOT EXISTS career_readiness (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track VARCHAR(50) NOT NULL, -- "Frontend", "Backend", "DevOps"
    score INT DEFAULT 0,
    last_updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, track)
);

-- Phase 4: Job Preferences & Job Matches
CREATE TABLE IF NOT EXISTS job_preferences (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(100) NOT NULL, -- "DevOps Engineer", "Backend Engineer", "Cloud Engineer"
    interest_level INT DEFAULT 1, -- 1-5 scale
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role)
);

CREATE TABLE IF NOT EXISTS job_matches (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(100) NOT NULL,
    readiness_score INT DEFAULT 0,
    matched_skills TEXT[], -- Array of matched skill names
    missing_skills TEXT[], -- Array of missing skill names
    learning_path TEXT[], -- Recommended skills to learn next
    recommended_next VARCHAR(100),
    calculated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role)
);

-- Phase 5: Projects
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    technologies TEXT[], -- Array of technologies used
    link TEXT, -- GitHub repo or project link
    duration_months INT,
    completion_date DATE,
    impact TEXT, -- Business impact or key achievements
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Phase 5: Certificates & Credentials
CREATE TABLE IF NOT EXISTS certificates (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    issuer VARCHAR(255), -- "AWS", "Google Cloud", "Microsoft Azure", etc.
    credential_id VARCHAR(255),
    credential_url TEXT,
    issue_date DATE NOT NULL,
    expiry_date DATE, -- NULL if no expiry
    skills_covered TEXT[], -- Array of skills this cert covers
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, credential_id)
);

-- Phase 5: Resume/Profile
CREATE TABLE IF NOT EXISTS user_profile (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    bio TEXT,
    headline VARCHAR(255), -- Professional headline
    location VARCHAR(100),
    website TEXT,
    resume_url TEXT, -- URL to uploaded resume
    github_url TEXT,
    linkedin_url TEXT,
    twitter_url TEXT,
    professional_profile TEXT, -- Generated professional summary
    visibility VARCHAR(20) DEFAULT 'private', -- 'private', 'public', 'link'
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_career_readiness_user ON career_readiness(user_id);
CREATE INDEX IF NOT EXISTS idx_job_preferences_user ON job_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_job_matches_user ON job_matches(user_id);
CREATE INDEX IF NOT EXISTS idx_projects_user ON projects(user_id);
CREATE INDEX IF NOT EXISTS idx_certificates_user ON certificates(user_id);
CREATE INDEX IF NOT EXISTS idx_learning_logs_date ON learning_logs(log_date);
