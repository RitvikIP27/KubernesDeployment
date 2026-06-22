# HelixaCore Database Migration - Quick Start

## Step 1: Access Supabase SQL Editor

1. Go to [Supabase Dashboard](https://supabase.com/dashboard)
2. Select your project
3. Go to **SQL Editor** → **New Query**

## Step 2: Copy Migration SQL

The migration file is located at: `postgres/migration_phase2-5.sql`

Here's the complete migration (ready to copy-paste):

```sql
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
```

## Step 3: Execute the Migration

1. **Paste** the SQL code into the Supabase SQL Editor
2. Click **Run** button (or press Ctrl+Enter)
3. You should see: "Query executed successfully" messages for each table

## Step 4: Verify Tables Created

Run this verification query:

```sql
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
ORDER BY table_name;
```

You should see these new tables:
- ✅ career_readiness
- ✅ job_preferences
- ✅ job_matches
- ✅ projects
- ✅ certificates
- ✅ user_profile

## Step 5: Rebuild & Deploy

After migration is complete:

```bash
# 1. Rebuild backend Docker image
cd /home/ritvik-kant/KubernesDeployment
docker compose build --no-cache backend

# 2. Restart services
docker compose down
docker compose up -d

# 3. Verify services are running
docker compose ps
```

## Step 6: Test the Features

1. Open http://localhost in your browser
2. Login with your account
3. Try the new tabs:
   - **📊 Analytics Dashboard** - View charts and heatmaps
   - **🎯 Career Readiness** - See your readiness scores
   - **💼 Job Matching** - Match against job roles
   - **📁 Portfolio** - Add projects and certificates
   - **👤 Profile** - Build your professional profile

## Troubleshooting

### Issue: "Table already exists"
- This is fine! The migration uses `IF NOT EXISTS`
- The tables won't be duplicated

### Issue: "Foreign key constraint failed"
- Make sure the `users` table exists first
- Run: `SELECT * FROM users;` to verify

### Issue: Tables not appearing in Supabase
- Refresh the Supabase dashboard
- Check the function list on the left sidebar
- Look under **Database** → **Tables**

### Issue: Backend returns 404 errors for new endpoints
- Verify backend is rebuilt with latest code
- Check that all handler files are in `backend/handlers/`
- Ensure `main.go` has all the new routes added

## Files Modified

### Backend
- ✅ `backend/main.go` - Added new routes
- ✅ `backend/models/analytics.go` - New model types
- ✅ `backend/handlers/projects_certificates.go` - New handlers
- ✅ `backend/handlers/career_profile.go` - New handlers

### Frontend
- ✅ `frontend/index.html` - Added new tabs and modals
- ✅ `frontend/css/style.css` - Added new styles
- ✅ `frontend/js/app.js` - Added new functions

### Database
- ✅ `postgres/migration_phase2-5.sql` - Migration file

## What's New

### Phase 2: Analytics Dashboard 🚀
- Learning hours visualization
- Weekly/monthly progress charts
- 84-day activity heatmap
- Learning streaks tracking
- Top 5 skills display

### Phase 3: Career Readiness Score 🎯
- Frontend readiness tracking
- Backend readiness tracking
- DevOps readiness tracking
- Skill matching analysis
- Progress indicators

### Phase 4: Job Matching 💼
- DevOps Engineer matching
- Backend Engineer matching
- Cloud Engineer matching
- Missing skills identification
- Learning path recommendations

### Phase 5: Portfolio + Profile 📁👤
- Project portfolio management
- Certificate/credential tracking
- Professional profile builder
- Auto-generated career summary
- Social links integration

---

**Ready to transform HelixaCore?** 🚀

Follow the steps above and you'll be up and running in minutes!

