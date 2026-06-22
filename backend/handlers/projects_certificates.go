package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/RitvikIP27/KubernesDeployment/database"
	"github.com/RitvikIP27/KubernesDeployment/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// Phase 5: Projects Handlers

// GetProjects returns all projects for the authenticated user
func GetProjects(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, title, description, technologies, link, duration_months, completion_date, impact, created_at, updated_at
		FROM projects
		WHERE user_id = $1
		ORDER BY completion_date DESC NULLS LAST, created_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch projects"})
		return
	}
	defer rows.Close()

	projects := []models.Project{}
	for rows.Next() {
		var p models.Project
		var techs pq.StringArray
		var completionDate sql.NullString
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &techs, &p.Link, &p.DurationMonths, &completionDate, &p.Impact, &p.CreatedAt, &p.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan project"})
			return
		}
		p.Technologies = techs
		if completionDate.Valid {
			p.CompletionDate = completionDate.String
		}
		projects = append(projects, p)
	}

	if projects == nil {
		projects = []models.Project{}
	}
	c.JSON(http.StatusOK, projects)
}

// GetProject returns a single project by ID
func GetProject(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	var p models.Project
	var techs pq.StringArray
	var completionDate sql.NullString
	err = database.DB.QueryRow(`
		SELECT id, title, description, technologies, link, duration_months, completion_date, impact, created_at, updated_at
		FROM projects
		WHERE id = $1 AND user_id = $2
	`, projectID, userID).Scan(&p.ID, &p.Title, &p.Description, &techs, &p.Link, &p.DurationMonths, &completionDate, &p.Impact, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch project"})
		return
	}

	p.Technologies = techs
	if completionDate.Valid {
		p.CompletionDate = completionDate.String
	}
	c.JSON(http.StatusOK, p)
}

// CreateProject creates a new project
func CreateProject(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var completionDate sql.NullString
	if req.CompletionDate != "" {
		completionDate = sql.NullString{String: req.CompletionDate, Valid: true}
	}

	err := database.DB.QueryRow(`
		INSERT INTO projects (user_id, title, description, technologies, link, duration_months, completion_date, impact, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`, userID, req.Title, req.Description, pq.Array(req.Technologies), req.Link, req.DurationMonths, completionDate, req.Impact, time.Now(), time.Now()).
		Scan(&req, &req, &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "project created successfully"})
}

// UpdateProject updates an existing project
func UpdateProject(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var completionDate sql.NullString
	if req.CompletionDate != "" {
		completionDate = sql.NullString{String: req.CompletionDate, Valid: true}
	}

	result, err := database.DB.Exec(`
		UPDATE projects
		SET title = $1, description = $2, technologies = $3, link = $4, duration_months = $5, completion_date = $6, impact = $7, updated_at = $8
		WHERE id = $9 AND user_id = $10
	`, req.Title, req.Description, pq.Array(req.Technologies), req.Link, req.DurationMonths, completionDate, req.Impact, time.Now(), projectID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update project"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project updated successfully"})
}

// DeleteProject deletes a project
func DeleteProject(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	result, err := database.DB.Exec(`DELETE FROM projects WHERE id = $1 AND user_id = $2`, projectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete project"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project deleted successfully"})
}

// Phase 5: Certificates Handlers

// GetCertificates returns all certificates for the authenticated user
func GetCertificates(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, name, issuer, credential_id, credential_url, issue_date, expiry_date, skills_covered, created_at
		FROM certificates
		WHERE user_id = $1
		ORDER BY issue_date DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch certificates"})
		return
	}
	defer rows.Close()

	certs := []models.Certificate{}
	for rows.Next() {
		var cert models.Certificate
		var skills pq.StringArray
		var expiryDate sql.NullString
		if err := rows.Scan(&cert.ID, &cert.Name, &cert.Issuer, &cert.CredentialID, &cert.CredentialURL, &cert.IssueDate, &expiryDate, &skills, &cert.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan certificate"})
			return
		}
		cert.SkillsCovered = skills
		if expiryDate.Valid {
			cert.ExpiryDate = expiryDate.String
		}
		certs = append(certs, cert)
	}

	if certs == nil {
		certs = []models.Certificate{}
	}
	c.JSON(http.StatusOK, certs)
}

// CreateCertificate creates a new certificate
func CreateCertificate(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.CreateCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expiryDate sql.NullString
	if req.ExpiryDate != "" {
		expiryDate = sql.NullString{String: req.ExpiryDate, Valid: true}
	}

	var certID int
	err := database.DB.QueryRow(`
		INSERT INTO certificates (user_id, name, issuer, credential_id, credential_url, issue_date, expiry_date, skills_covered, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, userID, req.Name, req.Issuer, req.CredentialID, req.CredentialURL, req.IssueDate, expiryDate, pq.Array(req.SkillsCovered), time.Now()).
		Scan(&certID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create certificate"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "certificate created successfully", "id": certID})
}

// DeleteCertificate deletes a certificate
func DeleteCertificate(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	certID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid certificate ID"})
		return
	}

	result, err := database.DB.Exec(`DELETE FROM certificates WHERE id = $1 AND user_id = $2`, certID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete certificate"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "certificate not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certificate deleted successfully"})
}
