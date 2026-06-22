package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RitvikIP27/KubernesDeployment/database"
	"github.com/RitvikIP27/KubernesDeployment/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// Phase 3: Career Readiness Handlers

// GetCareerReadiness returns career readiness scores for all tracks
func GetCareerReadiness(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get user's skills
	skills, err := loadAnalyticsSkills(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load skills"})
		return
	}

	// Calculate readiness scores
	readinessScores := []models.CareerReadinessDetail{}
	readinessRequirements := map[string][]string{
		"Frontend": {"HTML & CSS", "JavaScript", "React", "Tailwind CSS", "Web Performance", "Git"},
		"Backend":  {"Go", "REST APIs", "SQL & PostgreSQL", "Redis & Caching", "Docker", "System Design"},
		"DevOps":   {"Linux Basics", "Docker", "Kubernetes", "AWS Cloud", "Terraform", "CI/CD"},
	}

	for track, requirements := range readinessRequirements {
		matched, missing, score := evaluateRequirements(skills, requirements)
		progressPercent := 0
		if len(requirements) > 0 {
			progressPercent = (len(matched) * 100) / len(requirements)
		}

		readinessScores = append(readinessScores, models.CareerReadinessDetail{
			Track:           track,
			Score:           score,
			MatchedSkills:   matched,
			MissingSkills:   missing,
			ProgressPercent: progressPercent,
		})

		// Update/Insert into database
		_, err := database.DB.Exec(`
			INSERT INTO career_readiness (user_id, track, score, last_updated)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id, track) DO UPDATE SET score = $3, last_updated = $4
		`, userID, track, score, time.Now())
		if err != nil {
			// Log but don't fail the response
			continue
		}
	}

	c.JSON(http.StatusOK, readinessScores)
}

// GetCareerReadinessByTrack returns readiness score for a specific track
func GetCareerReadinessByTrack(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	track := c.Param("track")
	if track == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "track parameter required"})
		return
	}

	skills, err := loadAnalyticsSkills(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load skills"})
		return
	}

	readinessRequirements := map[string][]string{
		"Frontend": {"HTML & CSS", "JavaScript", "React", "Tailwind CSS", "Web Performance", "Git"},
		"Backend":  {"Go", "REST APIs", "SQL & PostgreSQL", "Redis & Caching", "Docker", "System Design"},
		"DevOps":   {"Linux Basics", "Docker", "Kubernetes", "AWS Cloud", "Terraform", "CI/CD"},
	}

	requirements, exists := readinessRequirements[track]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track"})
		return
	}

	matched, missing, score := evaluateRequirements(skills, requirements)
	progressPercent := 0
	if len(requirements) > 0 {
		progressPercent = (len(matched) * 100) / len(requirements)
	}

	c.JSON(http.StatusOK, models.CareerReadinessDetail{
		Track:           track,
		Score:           score,
		MatchedSkills:   matched,
		MissingSkills:   missing,
		ProgressPercent: progressPercent,
	})
}

// Phase 4: Job Preferences Handlers

// GetJobPreferences returns user's job preferences
func GetJobPreferences(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, role, interest_level, created_at
		FROM job_preferences
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch preferences"})
		return
	}
	defer rows.Close()

	prefs := []models.JobPreference{}
	for rows.Next() {
		var pref models.JobPreference
		if err := rows.Scan(&pref.ID, &pref.Role, &pref.InterestLevel, &pref.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan preference"})
			return
		}
		prefs = append(prefs, pref)
	}

	if prefs == nil {
		prefs = []models.JobPreference{}
	}
	c.JSON(http.StatusOK, prefs)
}

// AddJobPreference adds a job preference
func AddJobPreference(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.CreateJobPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var prefID int
	err := database.DB.QueryRow(`
		INSERT INTO job_preferences (user_id, role, interest_level, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, role) DO UPDATE SET interest_level = $3
		RETURNING id
	`, userID, req.Role, req.InterestLevel, time.Now()).Scan(&prefID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add preference"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "preference added successfully", "id": prefID})
}

// RemoveJobPreference removes a job preference
func RemoveJobPreference(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role parameter required"})
		return
	}

	result, err := database.DB.Exec(`
		DELETE FROM job_preferences WHERE user_id = $1 AND role = $2
	`, userID, role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove preference"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "preference not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "preference removed successfully"})
}

// Phase 4: Job Matching Handlers

// GetJobMatches returns job matches for the user based on their skills
func GetJobMatches(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	skills, err := loadAnalyticsSkills(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load skills"})
		return
	}

	jobRequirements := map[string][]string{
		"DevOps Engineer":  {"Linux Basics", "Docker", "Kubernetes", "AWS Cloud", "Terraform", "CI/CD", "Python Scripting"},
		"Backend Engineer": {"Go", "REST APIs", "SQL & PostgreSQL", "Redis & Caching", "Docker", "System Design"},
		"Cloud Engineer":   {"Linux Basics", "Docker", "Kubernetes", "AWS Cloud", "Terraform", "CI/CD", "Security"},
	}

	matches := []models.JobMatch{}
	for role, requirements := range jobRequirements {
		matched, missing, score := evaluateRequirements(skills, requirements)
		recommendedNext := "Keep strengthening your project portfolio."
		if len(missing) > 0 {
			recommendedNext = "Start with " + missing[0] + "."
		}

		matches = append(matches, models.JobMatch{
			Role:            role,
			ReadinessScore:  score,
			MatchedSkills:   matched,
			MissingSkills:   missing,
			LearningPath:    missing,
			RecommendedNext: recommendedNext,
			CalculatedAt:    time.Now(),
		})

		// Save to database
		_, _ = database.DB.Exec(`
			INSERT INTO job_matches (user_id, role, readiness_score, matched_skills, missing_skills, learning_path, recommended_next, calculated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (user_id, role) DO UPDATE SET 
				readiness_score = $3, matched_skills = $4, missing_skills = $5, learning_path = $6, recommended_next = $7, calculated_at = $8
		`, userID, role, score, pq.Array(matched), pq.Array(missing), pq.Array(missing), recommendedNext, time.Now())
	}

	c.JSON(http.StatusOK, matches)
}

// GetJobMatch returns job match for a specific role
func GetJobMatch(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role parameter required"})
		return
	}

	skills, err := loadAnalyticsSkills(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load skills"})
		return
	}

	jobRequirements := map[string][]string{
		"DevOps Engineer":  {"Linux Basics", "Docker", "Kubernetes", "AWS Cloud", "Terraform", "CI/CD", "Python Scripting"},
		"Backend Engineer": {"Go", "REST APIs", "SQL & PostgreSQL", "Redis & Caching", "Docker", "System Design"},
		"Cloud Engineer":   {"Linux Basics", "Docker", "Kubernetes", "AWS Cloud", "Terraform", "CI/CD", "Security"},
	}

	requirements, exists := jobRequirements[role]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	matched, missing, score := evaluateRequirements(skills, requirements)
	recommendedNext := "Keep strengthening your project portfolio."
	if len(missing) > 0 {
		recommendedNext = "Start with " + missing[0] + "."
	}

	c.JSON(http.StatusOK, models.JobMatch{
		Role:            role,
		ReadinessScore:  score,
		MatchedSkills:   matched,
		MissingSkills:   missing,
		LearningPath:    missing,
		RecommendedNext: recommendedNext,
		CalculatedAt:    time.Now(),
	})
}

// Phase 5: User Profile Handlers

// GetUserProfile returns the user's profile
func GetUserProfile(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var profile models.UserProfile
	err := database.DB.QueryRow(`
		SELECT id, bio, headline, location, website, resume_url, github_url, linkedin_url, twitter_url, professional_profile, visibility, created_at, updated_at
		FROM user_profile
		WHERE user_id = $1
	`, userID).Scan(&profile.ID, &profile.Bio, &profile.Headline, &profile.Location, &profile.Website, &profile.ResumeURL, &profile.GithubURL, &profile.LinkedinURL, &profile.TwitterURL, &profile.ProfessionalProfile, &profile.Visibility, &profile.CreatedAt, &profile.UpdatedAt)

	if err == sql.ErrNoRows {
		// Create default profile
		err = database.DB.QueryRow(`
			INSERT INTO user_profile (user_id, visibility, created_at, updated_at)
			VALUES ($1, 'private', $2, $3)
			ON CONFLICT (user_id) DO NOTHING
			RETURNING id, bio, headline, location, website, resume_url, github_url, linkedin_url, twitter_url, professional_profile, visibility, created_at, updated_at
		`, userID, time.Now(), time.Now()).Scan(&profile.ID, &profile.Bio, &profile.Headline, &profile.Location, &profile.Website, &profile.ResumeURL, &profile.GithubURL, &profile.LinkedinURL, &profile.TwitterURL, &profile.ProfessionalProfile, &profile.Visibility, &profile.CreatedAt, &profile.UpdatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profile"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateUserProfile updates the user's profile
func UpdateUserProfile(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec(`
		INSERT INTO user_profile (user_id, bio, headline, location, website, github_url, linkedin_url, twitter_url, visibility, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (user_id) DO UPDATE SET 
			bio = $2, headline = $3, location = $4, website = $5, github_url = $6, linkedin_url = $7, twitter_url = $8, visibility = $9, updated_at = $11
	`, userID, req.Bio, req.Headline, req.Location, req.Website, req.GithubURL, req.LinkedinURL, req.TwitterURL, req.Visibility, time.Now(), time.Now())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}

// GenerateProfessionalProfile generates a comprehensive professional profile
func GenerateProfessionalProfile(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	skills, _ := loadAnalyticsSkills(userID)

	// Get top projects
	projects := getTopProjects(userID, 3)

	// Get certificates
	certs := getUserCertificates(userID)

	// Calculate readiness scores
	readinessScores := calculateAllReadinessScores(skills)

	// Identify strength and growth areas
	strengthAreas := []string{}
	growthAreas := []string{}

	readinessRequirements := map[string][]string{
		"Frontend": {"HTML & CSS", "JavaScript", "React", "Tailwind CSS", "Web Performance", "Git"},
		"Backend":  {"Go", "REST APIs", "SQL & PostgreSQL", "Redis & Caching", "Docker", "System Design"},
		"DevOps":   {"Linux Basics", "Docker", "Kubernetes", "AWS Cloud", "Terraform", "CI/CD"},
	}

	for track, requirements := range readinessRequirements {
		_, _, score := evaluateRequirements(skills, requirements)
		if score >= 70 {
			strengthAreas = append(strengthAreas, track)
		} else if score < 40 {
			growthAreas = append(growthAreas, track)
		}
	}

	recommendations := buildRecommendations(skills, readinessScores)

	profProfile := models.ProfessionalProfile{
		Summary:         generateSummary(skills, strengthAreas, growthAreas),
		StrengthAreas:   strengthAreas,
		GrowthAreas:     growthAreas,
		Recommendations: recommendations,
		ReadinessScores: readinessScores,
		TopProjects:     projects,
		Certifications:  certs,
	}

	c.JSON(http.StatusOK, profProfile)
}

// Helper functions

func getTopProjects(userID int, limit int) []string {
	rows, err := database.DB.Query(`
		SELECT title FROM projects WHERE user_id = $1 ORDER BY completion_date DESC LIMIT $2
	`, userID, limit)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	projects := []string{}
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			continue
		}
		projects = append(projects, title)
	}
	return projects
}

func getUserCertificates(userID int) []string {
	rows, err := database.DB.Query(`
		SELECT name FROM certificates WHERE user_id = $1 ORDER BY issue_date DESC
	`, userID)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	certs := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		certs = append(certs, name)
	}
	return certs
}

func calculateAllReadinessScores(skills []analyticsSkill) map[string]int {
	scores := map[string]int{}
	readinessRequirements := map[string][]string{
		"Frontend": {"HTML & CSS", "JavaScript", "React", "Tailwind CSS", "Web Performance", "Git"},
		"Backend":  {"Go", "REST APIs", "SQL & PostgreSQL", "Redis & Caching", "Docker", "System Design"},
		"DevOps":   {"Linux Basics", "Docker", "Kubernetes", "AWS Cloud", "Terraform", "CI/CD"},
	}

	for track, requirements := range readinessRequirements {
		_, _, score := evaluateRequirements(skills, requirements)
		scores[track] = score
	}
	return scores
}

func buildRecommendations(skills []analyticsSkill, scores map[string]int) []string {
	recommendations := []string{}

	skillCount := len(skills)
	if skillCount < 5 {
		recommendations = append(recommendations, "Focus on building a broader skill foundation. Aim for at least 5-7 core skills.")
	}

	totalHours := 0.0
	for _, skill := range skills {
		totalHours += skill.TotalHours
	}

	if totalHours < 50 {
		recommendations = append(recommendations, "Increase learning commitment. Aim for at least 50 hours across your skills.")
	}

	lowestTrack := ""
	lowestScore := 100
	for track, score := range scores {
		if score < lowestScore {
			lowestScore = score
			lowestTrack = track
		}
	}

	if lowestTrack != "" && lowestScore < 50 {
		recommendations = append(recommendations, "Prioritize building "+lowestTrack+" skills to diversify your profile.")
	}

	recommendations = append(recommendations, "Build real projects to demonstrate your skills to potential employers.")

	return recommendations
}

func generateSummary(skills []analyticsSkill, strengthAreas []string, growthAreas []string) string {
	var summary strings.Builder
	summary.WriteString("Professional profile based on ")
	summary.WriteString(strconv.Itoa(len(skills)))
	summary.WriteString(" tracked skills")

	if len(strengthAreas) > 0 {
		summary.WriteString(". Strength areas: ")
		summary.WriteString(strings.Join(strengthAreas, ", "))
	}

	if len(growthAreas) > 0 {
		summary.WriteString(". Focus on growth in: ")
		summary.WriteString(strings.Join(growthAreas, ", "))
	}

	return summary.String()
}
