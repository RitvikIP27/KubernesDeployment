package handlers

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/RitvikIP27/KubernesDeployment/database"
	"github.com/RitvikIP27/KubernesDeployment/models"
	"github.com/gin-gonic/gin"
)

func GetDashboard(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var dash models.Dashboard
	database.DB.QueryRow("SELECT COUNT(*) FROM skills WHERE user_id = $1", userID).Scan(&dash.TotalSkills)
	database.DB.QueryRow("SELECT COALESCE(SUM(hours), 0) FROM learning_logs WHERE user_id = $1", userID).Scan(&dash.TotalHours)
	database.DB.QueryRow("SELECT COUNT(*) FROM learning_logs WHERE user_id = $1", userID).Scan(&dash.TotalLogs)

	err := database.DB.QueryRow(`
		SELECT s.name FROM skills s
		LEFT JOIN learning_logs l ON s.id = l.skill_id AND l.user_id = s.user_id
		WHERE s.user_id = $1
		GROUP BY s.id, s.name
		ORDER BY COALESCE(SUM(l.hours), 0) DESC
		LIMIT 1
	`, userID).Scan(&dash.TopSkill)
	if err != nil {
		dash.TopSkill = "N/A"
	}

	c.JSON(http.StatusOK, dash)
}

func HealthCheck(c *gin.Context) {
	err := database.DB.Ping()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

type analyticsPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

type activityDay struct {
	Date  string  `json:"date"`
	Hours float64 `json:"hours"`
	Level int     `json:"level"`
}

type topSkill struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Hours    float64 `json:"hours"`
	Progress float64 `json:"progress"`
}

type streakSummary struct {
	Current int `json:"current"`
	Longest int `json:"longest"`
}

type readinessScore struct {
	Track   string   `json:"track"`
	Score   int      `json:"score"`
	Matched []string `json:"matched"`
	Missing []string `json:"missing"`
}

type jobMatch struct {
	Role            string   `json:"role"`
	Readiness       int      `json:"readiness"`
	MatchedSkills   []string `json:"matched_skills"`
	MissingSkills   []string `json:"missing_skills"`
	LearningPath    []string `json:"learning_path"`
	RecommendedNext string   `json:"recommended_next"`
}

type analyticsResponse struct {
	LearningHours    []analyticsPoint `json:"learning_hours"`
	SkillGrowth      []analyticsPoint `json:"skill_growth"`
	WeeklyProgress   []analyticsPoint `json:"weekly_progress"`
	MonthlyProgress  []analyticsPoint `json:"monthly_progress"`
	Streaks          streakSummary    `json:"streaks"`
	TopSkills        []topSkill       `json:"top_skills"`
	ActivityCalendar []activityDay    `json:"activity_calendar"`
	ReadinessScores  []readinessScore `json:"readiness_scores"`
	JobMatches       []jobMatch       `json:"job_matches"`
}

type analyticsSkill struct {
	Name        string
	Category    string
	TargetHours int
	TotalHours  float64
}

type analyticsLog struct {
	SkillName string
	Hours     float64
	LogDate   time.Time
}

var readinessRequirements = map[string][]string{
	"Frontend Readiness": {"HTML & CSS", "JavaScript", "React", "Tailwind CSS", "Web Performance", "Git"},
	"Backend Readiness":  {"Go", "REST APIs", "SQL & PostgreSQL", "Redis & Caching", "Docker", "System Design"},
	"DevOps Readiness":   {"Linux Basics", "Docker", "Kubernetes", "AWS Cloud", "Terraform", "CI/CD"},
}

var jobRequirements = map[string][]string{
	"DevOps Engineer":  {"Linux Basics", "Docker", "Kubernetes", "AWS Cloud", "Terraform", "CI/CD", "Python Scripting"},
	"Backend Engineer": {"Go", "REST APIs", "SQL & PostgreSQL", "Redis & Caching", "Docker", "System Design"},
	"Cloud Engineer":   {"Linux Basics", "Docker", "Kubernetes", "AWS Cloud", "Terraform", "CI/CD", "Security"},
}

func GetAnalytics(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	skills, err := loadAnalyticsSkills(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logs, err := loadAnalyticsLogs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analyticsResponse{
		LearningHours:    buildLearningHours(logs),
		SkillGrowth:      buildSkillGrowth(skills),
		WeeklyProgress:   buildWeeklyProgress(logs),
		MonthlyProgress:  buildMonthlyProgress(logs),
		Streaks:          buildStreaks(logs),
		TopSkills:        buildTopSkills(skills),
		ActivityCalendar: buildActivityCalendar(logs),
		ReadinessScores:  buildReadinessScores(skills),
		JobMatches:       buildJobMatches(skills),
	})
}

func loadAnalyticsSkills(userID int) ([]analyticsSkill, error) {
	rows, err := database.DB.Query(`
		SELECT s.name, s.category, s.target_hours, COALESCE(SUM(l.hours), 0) AS total_hours
		FROM skills s
		LEFT JOIN learning_logs l ON s.id = l.skill_id AND l.user_id = s.user_id
		WHERE s.user_id = $1
		GROUP BY s.id, s.name, s.category, s.target_hours
		ORDER BY s.created_at
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	skills := []analyticsSkill{}
	for rows.Next() {
		var skill analyticsSkill
		if err := rows.Scan(&skill.Name, &skill.Category, &skill.TargetHours, &skill.TotalHours); err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}
	return skills, rows.Err()
}

func loadAnalyticsLogs(userID int) ([]analyticsLog, error) {
	rows, err := database.DB.Query(`
		SELECT s.name, l.hours, l.log_date
		FROM learning_logs l
		JOIN skills s ON s.id = l.skill_id AND s.user_id = l.user_id
		WHERE l.user_id = $1
		ORDER BY l.log_date
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []analyticsLog{}
	for rows.Next() {
		var log analyticsLog
		if err := rows.Scan(&log.SkillName, &log.Hours, &log.LogDate); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}

func buildLearningHours(logs []analyticsLog) []analyticsPoint {
	bySkill := map[string]float64{}
	for _, log := range logs {
		bySkill[log.SkillName] += log.Hours
	}

	points := make([]analyticsPoint, 0, len(bySkill))
	for name, hours := range bySkill {
		points = append(points, analyticsPoint{Label: name, Value: round1(hours)})
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Value > points[j].Value
	})
	if len(points) > 8 {
		points = points[:8]
	}
	return points
}

func buildSkillGrowth(skills []analyticsSkill) []analyticsPoint {
	points := make([]analyticsPoint, 0, len(skills))
	for _, skill := range skills {
		progress := 0.0
		if skill.TargetHours > 0 {
			progress = minFloat((skill.TotalHours/float64(skill.TargetHours))*100, 100)
		}
		points = append(points, analyticsPoint{Label: skill.Name, Value: round1(progress)})
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Value > points[j].Value
	})
	if len(points) > 8 {
		points = points[:8]
	}
	return points
}

func buildWeeklyProgress(logs []analyticsLog) []analyticsPoint {
	now := time.Now()
	start := startOfWeek(now).AddDate(0, 0, -7*5)
	weeks := make([]analyticsPoint, 6)
	indexByWeek := map[string]int{}
	for i := range weeks {
		weekStart := start.AddDate(0, 0, i*7)
		label := weekStart.Format("Jan 2")
		weeks[i] = analyticsPoint{Label: label}
		indexByWeek[weekStart.Format("2006-01-02")] = i
	}

	for _, log := range logs {
		weekStart := startOfWeek(log.LogDate)
		if idx, ok := indexByWeek[weekStart.Format("2006-01-02")]; ok {
			weeks[idx].Value += log.Hours
		}
	}
	roundPoints(weeks)
	return weeks
}

func buildMonthlyProgress(logs []analyticsLog) []analyticsPoint {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -5, 0)
	months := make([]analyticsPoint, 6)
	indexByMonth := map[string]int{}
	for i := range months {
		month := start.AddDate(0, i, 0)
		label := month.Format("Jan")
		months[i] = analyticsPoint{Label: label}
		indexByMonth[month.Format("2006-01")] = i
	}

	for _, log := range logs {
		if idx, ok := indexByMonth[log.LogDate.Format("2006-01")]; ok {
			months[idx].Value += log.Hours
		}
	}
	roundPoints(months)
	return months
}

func buildStreaks(logs []analyticsLog) streakSummary {
	if len(logs) == 0 {
		return streakSummary{}
	}

	activeDays := map[string]bool{}
	for _, log := range logs {
		activeDays[dateKey(log.LogDate)] = true
	}

	keys := make([]string, 0, len(activeDays))
	for key := range activeDays {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	longest := 0
	currentRun := 0
	var previous time.Time
	for i, key := range keys {
		day, _ := time.Parse("2006-01-02", key)
		if i == 0 || day.Sub(previous).Hours() == 24 {
			currentRun++
		} else {
			currentRun = 1
		}
		if currentRun > longest {
			longest = currentRun
		}
		previous = day
	}

	current := 0
	cursor := normalizeDate(time.Now())
	if !activeDays[dateKey(cursor)] {
		cursor = cursor.AddDate(0, 0, -1)
	}
	for activeDays[dateKey(cursor)] {
		current++
		cursor = cursor.AddDate(0, 0, -1)
	}

	return streakSummary{Current: current, Longest: longest}
}

func buildTopSkills(skills []analyticsSkill) []topSkill {
	top := make([]topSkill, 0, len(skills))
	for _, skill := range skills {
		progress := 0.0
		if skill.TargetHours > 0 {
			progress = minFloat((skill.TotalHours/float64(skill.TargetHours))*100, 100)
		}
		top = append(top, topSkill{
			Name:     skill.Name,
			Category: skill.Category,
			Hours:    round1(skill.TotalHours),
			Progress: round1(progress),
		})
	}
	sort.Slice(top, func(i, j int) bool {
		return top[i].Hours > top[j].Hours
	})
	if len(top) > 5 {
		top = top[:5]
	}
	return top
}

func buildActivityCalendar(logs []analyticsLog) []activityDay {
	now := normalizeDate(time.Now())
	start := now.AddDate(0, 0, -83)
	hoursByDay := map[string]float64{}
	maxHours := 0.0
	for _, log := range logs {
		key := dateKey(log.LogDate)
		hoursByDay[key] += log.Hours
		if hoursByDay[key] > maxHours {
			maxHours = hoursByDay[key]
		}
	}

	days := make([]activityDay, 0, 84)
	for i := 0; i < 84; i++ {
		day := start.AddDate(0, 0, i)
		hours := round1(hoursByDay[dateKey(day)])
		level := 0
		if hours > 0 && maxHours > 0 {
			level = int((hours / maxHours) * 4)
			if level < 1 {
				level = 1
			}
			if level > 4 {
				level = 4
			}
		}
		days = append(days, activityDay{Date: dateKey(day), Hours: hours, Level: level})
	}
	return days
}

func buildReadinessScores(skills []analyticsSkill) []readinessScore {
	scores := make([]readinessScore, 0, len(readinessRequirements))
	tracks := []string{"Frontend Readiness", "Backend Readiness", "DevOps Readiness"}
	for _, track := range tracks {
		matched, missing, score := evaluateRequirements(skills, readinessRequirements[track])
		scores = append(scores, readinessScore{Track: track, Score: score, Matched: matched, Missing: missing})
	}
	return scores
}

func buildJobMatches(skills []analyticsSkill) []jobMatch {
	roles := []string{"DevOps Engineer", "Backend Engineer", "Cloud Engineer"}
	matches := make([]jobMatch, 0, len(roles))
	for _, role := range roles {
		matched, missing, score := evaluateRequirements(skills, jobRequirements[role])
		next := "Keep strengthening your project portfolio."
		if len(missing) > 0 {
			next = "Start with " + missing[0] + "."
		}
		matches = append(matches, jobMatch{
			Role:            role,
			Readiness:       score,
			MatchedSkills:   matched,
			MissingSkills:   missing,
			LearningPath:    missing,
			RecommendedNext: next,
		})
	}
	return matches
}

func evaluateRequirements(skills []analyticsSkill, requirements []string) ([]string, []string, int) {
	matched := []string{}
	missing := []string{}
	total := 0.0
	for _, req := range requirements {
		best := 0.0
		for _, skill := range skills {
			if skillMatchesRequirement(skill, req) {
				score := 0.25
				if skill.TargetHours > 0 {
					score = minFloat(skill.TotalHours/float64(skill.TargetHours), 1)
				} else if skill.TotalHours > 0 {
					score = minFloat(skill.TotalHours/20, 1)
				}
				if score > best {
					best = score
				}
			}
		}
		total += best
		if best >= 0.35 {
			matched = append(matched, req)
		} else {
			missing = append(missing, req)
		}
	}
	if len(requirements) == 0 {
		return matched, missing, 0
	}
	return matched, missing, int(minFloat((total/float64(len(requirements)))*100, 100) + 0.5)
}

func skillMatchesRequirement(skill analyticsSkill, requirement string) bool {
	name := strings.ToLower(skill.Name)
	category := strings.ToLower(skill.Category)
	req := strings.ToLower(requirement)
	for _, part := range strings.FieldsFunc(req, func(r rune) bool {
		return r == ' ' || r == '&' || r == '/' || r == '-'
	}) {
		part = strings.TrimSpace(part)
		if len(part) >= 3 && strings.Contains(name, part) {
			return true
		}
	}
	return strings.Contains(name, req) || strings.Contains(req, name) || strings.Contains(category, req)
}

func startOfWeek(t time.Time) time.Time {
	day := normalizeDate(t)
	offset := int(day.Weekday())
	if offset == 0 {
		offset = 7
	}
	return day.AddDate(0, 0, -(offset - 1))
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func dateKey(t time.Time) string {
	return normalizeDate(t).Format("2006-01-02")
}

func roundPoints(points []analyticsPoint) {
	for i := range points {
		points[i].Value = round1(points[i].Value)
	}
}

func round1(value float64) float64 {
	return float64(int(value*10+0.5)) / 10
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
