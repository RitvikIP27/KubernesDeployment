package models

import "time"

// Phase 3: Career Readiness Score
type CareerReadiness struct {
	ID          int       `json:"id"`
	Track       string    `json:"track"` // "Frontend", "Backend", "DevOps"
	Score       int       `json:"score"`
	LastUpdated time.Time `json:"last_updated"`
}

type CareerReadinessDetail struct {
	Track           string   `json:"track"`
	Score           int      `json:"score"`
	MatchedSkills   []string `json:"matched_skills"`
	MissingSkills   []string `json:"missing_skills"`
	ProgressPercent int      `json:"progress_percent"`
}

// Phase 4: Job Preferences & Matching
type JobPreference struct {
	ID            int       `json:"id"`
	Role          string    `json:"role"`
	InterestLevel int       `json:"interest_level"` // 1-5 scale
	CreatedAt     time.Time `json:"created_at"`
}

type CreateJobPreferenceRequest struct {
	Role          string `json:"role" binding:"required"`
	InterestLevel int    `json:"interest_level"`
}

type JobMatch struct {
	ID              int       `json:"id"`
	Role            string    `json:"role"`
	ReadinessScore  int       `json:"readiness_score"`
	MatchedSkills   []string  `json:"matched_skills"`
	MissingSkills   []string  `json:"missing_skills"`
	LearningPath    []string  `json:"learning_path"`
	RecommendedNext string    `json:"recommended_next"`
	CalculatedAt    time.Time `json:"calculated_at"`
}

// Phase 5: Projects
type Project struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Technologies   []string  `json:"technologies"`
	Link           string    `json:"link"`
	DurationMonths int       `json:"duration_months"`
	CompletionDate string    `json:"completion_date"`
	Impact         string    `json:"impact"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateProjectRequest struct {
	Title          string   `json:"title" binding:"required"`
	Description    string   `json:"description"`
	Technologies   []string `json:"technologies"`
	Link           string   `json:"link"`
	DurationMonths int      `json:"duration_months"`
	CompletionDate string   `json:"completion_date"`
	Impact         string   `json:"impact"`
}

type UpdateProjectRequest struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Technologies   []string `json:"technologies"`
	Link           string   `json:"link"`
	DurationMonths int      `json:"duration_months"`
	CompletionDate string   `json:"completion_date"`
	Impact         string   `json:"impact"`
}

// Phase 5: Certificates
type Certificate struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Issuer        string    `json:"issuer"`
	CredentialID  string    `json:"credential_id"`
	CredentialURL string    `json:"credential_url"`
	IssueDate     string    `json:"issue_date"`
	ExpiryDate    string    `json:"expiry_date,omitempty"`
	SkillsCovered []string  `json:"skills_covered"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateCertificateRequest struct {
	Name          string   `json:"name" binding:"required"`
	Issuer        string   `json:"issuer"`
	CredentialID  string   `json:"credential_id" binding:"required"`
	CredentialURL string   `json:"credential_url"`
	IssueDate     string   `json:"issue_date" binding:"required"`
	ExpiryDate    string   `json:"expiry_date"`
	SkillsCovered []string `json:"skills_covered"`
}

// Phase 5: User Profile / Professional Profile
type UserProfile struct {
	ID                  int       `json:"id"`
	Bio                 string    `json:"bio"`
	Headline            string    `json:"headline"`
	Location            string    `json:"location"`
	Website             string    `json:"website"`
	ResumeURL           string    `json:"resume_url"`
	GithubURL           string    `json:"github_url"`
	LinkedinURL         string    `json:"linkedin_url"`
	TwitterURL          string    `json:"twitter_url"`
	ProfessionalProfile string    `json:"professional_profile"`
	Visibility          string    `json:"visibility"` // 'private', 'public', 'link'
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type UpdateProfileRequest struct {
	Bio         string `json:"bio"`
	Headline    string `json:"headline"`
	Location    string `json:"location"`
	Website     string `json:"website"`
	GithubURL   string `json:"github_url"`
	LinkedinURL string `json:"linkedin_url"`
	TwitterURL  string `json:"twitter_url"`
	Visibility  string `json:"visibility"`
}

// Professional profile generation response
type ProfessionalProfile struct {
	Summary         string         `json:"summary"`
	StrengthAreas   []string       `json:"strength_areas"`
	GrowthAreas     []string       `json:"growth_areas"`
	Recommendations []string       `json:"recommendations"`
	ReadinessScores map[string]int `json:"readiness_scores"`
	TopProjects     []string       `json:"top_projects"`
	Certifications  []string       `json:"certifications"`
}

// Analytics response structures
type AnalyticsSummary struct {
	TotalLearningHours  float64                 `json:"total_learning_hours"`
	ActiveSkills        int                     `json:"active_skills"`
	CurrentStreak       int                     `json:"current_streak"`
	LongestStreak       int                     `json:"longest_streak"`
	TopSkills           []SkillProgress         `json:"top_skills"`
	ReadinessScores     []CareerReadinessDetail `json:"readiness_scores"`
	RecommendedJobRoles []string                `json:"recommended_job_roles"`
}

type SkillProgress struct {
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Hours       float64 `json:"hours"`
	Progress    int     `json:"progress"` // percentage
	TargetHours int     `json:"target_hours"`
}

type HeatmapData struct {
	Date  string  `json:"date"`
	Hours float64 `json:"hours"`
	Level int     `json:"level"` // 0-4 intensity level
}

type ActivityCalendarResponse struct {
	Year  int           `json:"year"`
	Month int           `json:"month"`
	Days  []HeatmapData `json:"days"`
}
