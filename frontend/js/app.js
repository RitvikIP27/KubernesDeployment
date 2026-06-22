const API = '/api';

// Adding authentication service
let authToken = localStorage.getItem('helixacore-token');

function readAuthRedirect() {
    const params = new URLSearchParams(window.location.search);
    const token = params.get('token');
    const authError = params.get('auth_error');

    if (token) {
        authToken = token;
        localStorage.setItem('helixacore-token', authToken);
        window.history.replaceState({}, document.title, window.location.pathname);
    } else if (authError) {
        window.history.replaceState({}, document.title, window.location.pathname);
    }

    return { token, authError };
}

const authRedirect = readAuthRedirect();

function getAuthHeaders() {
    const headers = { 'Content-Type': 'application/json' };
    if (authToken) {
        headers['Authorization'] = `Bearer ${authToken}`;
    }
    return headers;
}

function logout() {
    localStorage.removeItem('helixacore-token');
    authToken = null;
    showLoginPage();
}

function showLoginPage() {
    const authContainer = document.getElementById('auth-container');
    const appContainer = document.querySelector('main.container');
    const header = document.querySelector('.header');

    if (authContainer) authContainer.hidden = false;
    if (appContainer) appContainer.hidden = true;
    if (header) header.hidden = true;
}

function hideLoginPage() {
    const authContainer = document.getElementById('auth-container');
    const appContainer = document.querySelector('main.container');
    const header = document.querySelector('.header');

    if (authContainer) authContainer.hidden = true;
    if (appContainer) appContainer.hidden = false;
    if (header) header.hidden = false;
}

function switchAuthTab(event, tab) {
    document.querySelectorAll('.auth-tab').forEach(el => el.classList.remove('active'));
    document.querySelectorAll('.auth-tab-btn').forEach(el => el.classList.remove('active'));
    const targetTab = document.getElementById(tab + '-tab');
    if (targetTab) targetTab.classList.add('active');
    if (event && event.currentTarget) event.currentTarget.classList.add('active');
}

function startGoogleLogin() {
    window.location.href = `${API}/auth/google`;
}

async function readResponse(res) {
    const contentType = res.headers.get('content-type') || '';
    if (contentType.includes('application/json')) {
        return res.json();
    }

    const text = await res.text();
    return { error: text || `Request failed with status ${res.status}` };
}

async function handleLogin(e) {
    e.preventDefault();
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    try {
        const res = await fetch(`${API}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        if (!res.ok) {
            const err = await readResponse(res);
            document.getElementById('login-error').textContent = err.error || 'Login failed';
            return;
        }

        const data = await readResponse(res);
        authToken = data.token;
        localStorage.setItem('helixacore-token', authToken);
        triggerConfetti();
        setTimeout(() => location.reload(), 900);
    } catch (err) {
        document.getElementById('login-error').textContent = 'An error occurred: ' + err.message;
    }
}

async function handleRegister(e) {
    e.preventDefault();
    const name = document.getElementById('register-name').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;

    try {
        const res = await fetch(`${API}/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, email, password })
        });

        if (!res.ok) {
            const err = await readResponse(res);
            document.getElementById('register-error').textContent = err.error || 'Registration failed';
            return;
        }

        const data = await readResponse(res);
        authToken = data.token;
        localStorage.setItem('helixacore-token', authToken);
        triggerConfetti();
        setTimeout(() => location.reload(), 900);
    } catch (err) {
        document.getElementById('register-error').textContent = 'An error occurred: ' + err.message;
    }
}

// Roadmaps config  KI definition
const ROADMAPS = {
    "Site Reliability Engineer (SRE)": [
        { name: "Linux Basics", category: "DevOps", target_hours: 20 },
        { name: "Docker", category: "DevOps", target_hours: 40 },
        { name: "Kubernetes", category: "DevOps", target_hours: 60 },
        { name: "Prometheus & Grafana", category: "Other", target_hours: 30 },
        { name: "Terraform", category: "DevOps", target_hours: 35 },
        { name: "Go", category: "Programming", target_hours: 50 },
        { name: "CI/CD", category: "DevOps", target_hours: 25 }
    ],
    "DevOps Engineer": [
        { name: "Linux Basics", category: "DevOps", target_hours: 25 },
        { name: "Docker", category: "DevOps", target_hours: 40 },
        { name: "Kubernetes", category: "DevOps", target_hours: 60 },
        { name: "AWS Cloud", category: "Cloud", target_hours: 50 },
        { name: "Terraform", category: "DevOps", target_hours: 40 },
        { name: "CI/CD", category: "DevOps", target_hours: 30 },
        { name: "Python Scripting", category: "Programming", target_hours: 30 }
    ],
    "Backend Developer": [
        { name: "Go", category: "Programming", target_hours: 50 },
        { name: "REST APIs", category: "Programming", target_hours: 30 },
        { name: "SQL & PostgreSQL", category: "Databases", target_hours: 40 },
        { name: "Redis & Caching", category: "Databases", target_hours: 20 },
        { name: "Docker", category: "DevOps", target_hours: 25 },
        { name: "System Design", category: "Other", target_hours: 45 }
    ],
    "Frontend Developer": [
        { name: "HTML & CSS", category: "Other", target_hours: 20 },
        { name: "JavaScript", category: "Programming", target_hours: 50 },
        { name: "React", category: "Programming", target_hours: 60 },
        { name: "Tailwind CSS", category: "Other", target_hours: 15 },
        { name: "Web Performance", category: "Other", target_hours: 25 },
        { name: "Git", category: "Other", target_hours: 15 }
    ]
};

// Theme Management
function getPreferredTheme() {
    const stored = localStorage.getItem('helixacore-theme');
    if (stored) return stored;
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function applyTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('helixacore-theme', theme);
    const btn = document.getElementById('theme-toggle');
    if (btn) btn.textContent = theme === 'dark' ? '\u2600\uFE0F' : '\uD83C\uDF19';
}

document.addEventListener('DOMContentLoaded', () => {
    applyTheme(getPreferredTheme());

    const toggleBtn = document.getElementById('theme-toggle');
    if (toggleBtn) {
        toggleBtn.addEventListener('click', () => {
            const current = document.documentElement.getAttribute('data-theme');
            applyTheme(current === 'dark' ? 'light' : 'dark');
        });
    }
});

window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    if (!localStorage.getItem('helixacore-theme')) {
        applyTheme(e.matches ? 'dark' : 'light');
    }
});

// State
let skills = [];
let dashboard = {};
let currentTargetRole = "Site Reliability Engineer (SRE)";
let activeTab = "roadmap";

// DOM Elements
const statsContainer = document.getElementById('stats');
const skillsGrid = document.getElementById('skills-grid');
const addSkillModal = document.getElementById('add-skill-modal');
const logSessionModal = document.getElementById('log-session-modal');
const addSkillForm = document.getElementById('add-skill-form');
const logSessionForm = document.getElementById('log-session-form');

// Initialize
document.addEventListener('DOMContentLoaded', async () => {
    if (!authToken) {
        showLoginPage();
        if (authRedirect.authError) {
            const loginError = document.getElementById('login-error');
            if (loginError) loginError.textContent = authRedirect.authError;
        }
        return;
    }

    // ShOwnig  logout button to header
    const header = document.querySelector('.header-content');
    if (header) {
        const logoutBtn = document.createElement('button');
        logoutBtn.className = 'logout-btn';
        logoutBtn.textContent = 'Logout';
        logoutBtn.style.cssText = 'position: absolute; top: 1rem; right: 3rem; padding: 0.5rem 1rem; background: #f44; color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: 600;';
        logoutBtn.onclick = logout;
        header.parentElement.appendChild(logoutBtn);
    }

    await loadSettings();
    loadDashboard();
    loadSkills();
});

// API Calls
async function loadSettings() {
    try {
        const res = await fetch(`${API}/settings`, {
            headers: getAuthHeaders()
        });
        if (res.ok) {
            const settings = await res.json();
            if (settings && settings.target_role) {
                currentTargetRole = settings.target_role;
                const select = document.getElementById("role-select");
                if (select) select.value = currentTargetRole;
            }
        }
    } catch (err) {
        console.error('Failed to load settings:', err);
    }
}

async function saveTargetRole(role) {
    try {
        const res = await fetch(`${API}/settings`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({ key: 'target_role', value: role }),
        });
        if (!res.ok) throw new Error('Failed to save target role setting');
    } catch (err) {
        console.error('Failed to save settings:', err);
    }
}

async function loadDashboard() {
    try {
        const res = await fetch(`${API}/dashboard`, {
            headers: getAuthHeaders()
        });
        dashboard = await res.json();
        renderStats();
    } catch (err) {
        console.error('Failed to load dashboard:', err);
    }
}

async function loadSkills() {
    try {
        const res = await fetch(`${API}/skills`, {
            headers: getAuthHeaders()
        });
        skills = await res.json();
        renderSkills();
        renderRoadmap();
        updatePreparednessAndBadges();
    } catch (err) {
        console.error('Failed to load skills:', err);
    }
}

async function createSkill(data) {
    const res = await fetch(`${API}/skills`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to create skill');
    return res.json();
}

async function deleteSkill(id) {
    const res = await fetch(`${API}/skills/${id}`, {
        method: 'DELETE',
        headers: getAuthHeaders()
    });
    if (!res.ok) throw new Error('Failed to delete skill');
    return res.json();
}

async function logSession(skillId, data) {
    const res = await fetch(`${API}/skills/${skillId}/log`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to log session');
    return res.json();
}

// Render Functions
function renderStats() {
    if (!statsContainer) return;
    statsContainer.innerHTML = `
        <div class="stat-card">
            <div class="label">Total Skills</div>
            <div class="value">${dashboard.total_skills || 0}</div>
        </div>
        <div class="stat-card">
            <div class="label">Hours Logged</div>
            <div class="value">${(dashboard.total_hours || 0).toFixed(1)}</div>
        </div>
        <div class="stat-card">
            <div class="label">Sessions</div>
            <div class="value">${dashboard.total_logs || 0}</div>
        </div>
        <div class="stat-card">
            <div class="label">Top Skill</div>
            <div class="value" style="font-size:1.2rem">${dashboard.top_skill || 'N/A'}</div>
        </div>
    `;
}

function renderSkills() {
    if (!skillsGrid) return;
    if (!skills || skills.length === 0) {
        skillsGrid.innerHTML = `
            <div class="empty-state" style="grid-column: 1 / -1">
                <h3>No skills yet</h3>
                <p>Click "Add Custom Skill" to start tracking your learning journey.</p>
            </div>
        `;
        return;
    }

    skillsGrid.innerHTML = skills.map(skill => {
        const progress = skill.target_hours > 0
            ? Math.min((skill.total_hours / skill.target_hours) * 100, 100)
            : 0;

        return `
            <div class="skill-card">
                <div class="skill-header">
                    <span class="skill-name">${escapeHtml(skill.name)}</span>
                    ${skill.category ? `<span class="skill-category">${escapeHtml(skill.category)}</span>` : ''}
                </div>
                <div class="progress-bar">
                    <div class="fill" style="width: ${progress}%"></div>
                </div>
                <div class="progress-text">
                    <span>${skill.total_hours.toFixed(1)} hrs logged</span>
                    <span>${skill.target_hours > 0 ? skill.target_hours + ' hrs goal' : 'No goal set'}</span>
                </div>
                <div class="skill-actions">
                    <button class="btn btn-primary btn-sm" onclick="openLogModal(${skill.id}, '${escapeHtml(skill.name)}')">
                        + Log Session
                    </button>
                    <button class="btn btn-danger btn-sm" onclick="handleDelete(${skill.id})">
                        Delete
                    </button>
                </div>
            </div>
        `;
    }).join('');
}

function renderRoadmap() {
    const roadmapGrid = document.getElementById("roadmap-grid");
    if (!roadmapGrid) return;

    const requirements = ROADMAPS[currentTargetRole] || [];
    roadmapGrid.innerHTML = requirements.map(req => {
        // Find matching active skill
        const activeSkill = skills.find(s => s.name.toLowerCase() === req.name.toLowerCase());
        const loggedHours = activeSkill ? activeSkill.total_hours : 0;
        const progress = Math.min((loggedHours / req.target_hours) * 100, 100);

        let cardClass = "roadmap-card";
        let statusText = "Not Started";
        let statusClass = "unstarted";

        if (activeSkill) {
            if (progress >= 100) {
                cardClass += " completed";
                statusText = "Mastered";
                statusClass = "completed";
            } else {
                statusText = "Learning";
                statusClass = "learning";
            }
        } else {
            cardClass += " unstarted";
        }

        return `
            <div class="${cardClass}">
                <div>
                    <div class="skill-header">
                        <span class="skill-name">${escapeHtml(req.name)}</span>
                        <span class="skill-category">${escapeHtml(req.category)}</span>
                    </div>
                    <div class="progress-bar">
                        <div class="fill" style="width: ${progress}%"></div>
                    </div>
                    <div class="progress-text">
                        <span>${loggedHours.toFixed(1)} hrs logged</span>
                        <span>${req.target_hours} hrs target</span>
                    </div>
                </div>
                <div class="skill-meta">
                    <span class="status-badge ${statusClass}">${statusText}</span>
                    ${!activeSkill ? `
                        <button class="btn btn-primary btn-sm" onclick="startLearningRoadmapSkill('${escapeHtml(req.name)}', '${escapeHtml(req.category)}', ${req.target_hours})">
                            + Start Learning
                        </button>
                    ` : `
                        <button class="btn btn-primary btn-sm" onclick="openLogModal(${activeSkill.id}, '${escapeHtml(activeSkill.name)}')">
                            + Log Session
                        </button>
                    `}
                </div>
            </div>
        `;
    }).join('');
}

function updatePreparednessAndBadges() {
    const requirements = ROADMAPS[currentTargetRole] || [];
    if (requirements.length === 0) return;

    let totalProgress = 0;
    requirements.forEach(req => {
        const activeSkill = skills.find(s => s.name.toLowerCase() === req.name.toLowerCase());
        const loggedHours = activeSkill ? activeSkill.total_hours : 0;
        const progress = Math.min((loggedHours / req.target_hours) * 100, 100);
        totalProgress += progress;
    });

    const preparedness = totalProgress / requirements.length;
    
    // Update preparedness UI
    const prepPercentage = document.getElementById("preparedness-percentage");
    const prepFill = document.getElementById("preparedness-fill");
    
    if (prepPercentage) prepPercentage.textContent = `${preparedness.toFixed(0)}%`;
    if (prepFill) prepFill.style.width = `${preparedness}%`;

    // Calculate total hours
    const totalHours = skills.reduce((acc, curr) => acc + curr.total_hours, 0);

    // Badges definitions
    const badges = [
        { id: "level-1", name: "Beginner", icon: "🌱", description: "Prep > 10% or 5h logged", level: "Level 1", check: (p, h) => p >= 10 || h >= 5 },
        { id: "level-2", name: "Intermediate", icon: "🛡️", description: "Prep > 35% or 20h logged", level: "Level 2", check: (p, h) => p >= 35 || h >= 20 },
        { id: "level-3", name: "Advanced", icon: "🚀", description: "Prep > 60% or 50h logged", level: "Level 3", check: (p, h) => p >= 60 || h >= 50 },
        { id: "level-4", name: "Expert", icon: "🔮", description: "Prep > 80% or 100h logged", level: "Level 4", check: (p, h) => p >= 80 || h >= 100 },
        { id: "level-5", name: "Master", icon: "👑", description: "Prep > 95% or 200h logged", level: "Level 5", check: (p, h) => p >= 95 || h >= 200 }
    ];

    const badgesGrid = document.getElementById("badges-grid");
    if (badgesGrid) {
        badgesGrid.innerHTML = badges.map(badge => {
            const unlocked = badge.check(preparedness, totalHours);
            const statusClass = unlocked ? `unlocked ${badge.id}` : 'locked';
            const statusTooltip = unlocked ? 'Unlocked!' : `Locked: Requires ${badge.description}`;
            return `
                <div class="badge-card ${statusClass}" title="${statusTooltip}">
                    <span class="badge-icon">${badge.icon}</span>
                    <span class="badge-name">${badge.name}</span>
                    <span class="badge-level">${badge.level}</span>
                </div>
            `;
        }).join('');
    }
}

// Quick start roadmap skill
async function startLearningRoadmapSkill(name, category, targetHours) {
    try {
        await createSkill({
            name: name,
            category: category,
            target_hours: targetHours
        });
        showToast(`Started learning "${name}"!`, 'success');
        await loadDashboard();
        await loadSkills();
    } catch (err) {
        showToast('Failed to start learning skill', 'error');
    }
}

// Change Target Role
async function changeTargetRole(role) {
    currentTargetRole = role;
    await saveTargetRole(role);
    renderRoadmap();
    updatePreparednessAndBadges();
    showToast(`Target role updated to: ${role}`, 'success');
}

// Switch Tab
function switchTab(tab) {
    activeTab = tab;
    document.querySelectorAll(".tab-btn").forEach(el => el.classList.remove("active"));
    document.querySelectorAll(".tab-content").forEach(el => el.classList.remove("active"));

    const activeTabBtn = document.getElementById(`tab-${tab}`);
    const activeTabContent = document.getElementById(`content-${tab}`);
    
    if (activeTabBtn) activeTabBtn.classList.add("active");
    if (activeTabContent) activeTabContent.classList.add("active");
}

// Modal Handlers
function openAddModal() {
    addSkillForm.reset();
    addSkillModal.classList.add('active');
}

function closeAddModal() {
    addSkillModal.classList.remove('active');
}

let currentLogSkillId = null;

function openLogModal(skillId, skillName) {
    currentLogSkillId = skillId;
    document.getElementById('log-skill-name').textContent = skillName;
    document.getElementById('log-date').value = new Date().toISOString().split('T')[0];
    logSessionForm.reset();
    document.getElementById('log-date').value = new Date().toISOString().split('T')[0];
    logSessionModal.classList.add('active');
}

function closeLogModal() {
    logSessionModal.classList.remove('active');
    currentLogSkillId = null;
}

// Form Handlers
addSkillForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
        await createSkill({
            name: document.getElementById('skill-name').value,
            category: document.getElementById('skill-category').value,
            target_hours: parseInt(document.getElementById('skill-target').value) || 0,
        });
        closeAddModal();
        showToast('Skill added!', 'success');
        loadDashboard();
        loadSkills();
    } catch (err) {
        showToast('Failed to add skill', 'error');
    }
});

logSessionForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
        await logSession(currentLogSkillId, {
            hours: parseFloat(document.getElementById('log-hours').value),
            notes: document.getElementById('log-notes').value,
            log_date: document.getElementById('log-date').value,
        });
        closeLogModal();
        showToast('Session logged!', 'success');
        loadDashboard();
        loadSkills();
    } catch (err) {
        showToast('Failed to log session', 'error');
    }
});

async function handleDelete(id) {
    if (!confirm('Delete this skill and all its logs?')) return;
    try {
        await deleteSkill(id);
        showToast('Skill deleted', 'success');
        loadDashboard();
        loadSkills();
    } catch (err) {
        showToast('Failed to delete skill', 'error');
    }
}

// Utilities
function escapeHtml(str) {
    if (!str) return '';
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

function showToast(message, type = 'success') {
    const toast = document.getElementById('toast');
    if (!toast) return;
    toast.textContent = message;
    toast.className = `toast ${type} show`;
    setTimeout(() => toast.classList.remove('show'), 3000);
}

function triggerConfetti() {
    const container = document.getElementById('confetti-container');
    if (!container) return;

    const colors = ['#8b5cf6', '#22c55e', '#38bdf8', '#f97316', '#ec4899', '#fbbf24'];
    for (let i = 0; i < 24; i++) {
        const piece = document.createElement('div');
        const size = Math.floor(Math.random() * 10) + 6;
        piece.className = 'confetti-piece';
        piece.style.background = colors[Math.floor(Math.random() * colors.length)];
        piece.style.left = `${Math.random() * 100}%`;
        piece.style.width = `${size}px`;
        piece.style.height = `${size * 0.4}px`;
        piece.style.opacity = Math.random() * 0.9 + 0.4;
        piece.style.transform = `rotate(${Math.random() * 360}deg)`;
        container.appendChild(piece);

        requestAnimationFrame(() => {
            piece.style.transform += ` translateY(120vh) rotate(${Math.random() * 720}deg)`;
            piece.style.opacity = '0';
        });

        setTimeout(() => {
            piece.remove();
        }, 1400 + Math.random() * 400);
    }
}

// Close modals on backdrop click
function attachBackdropListeners() {
    document.querySelectorAll('.modal-backdrop').forEach(el => {
        el.addEventListener('click', (e) => {
            if (e.target === el) {
                el.classList.remove('active');
            }
        });
    });
}

attachBackdropListeners();

// ============================================
// PHASE 2: ANALYTICS DASHBOARD
// ============================================

async function loadAnalytics() {
    if (!authToken) return;
    
    try {
        const res = await fetch(`${API}/analytics`, {
            headers: getAuthHeaders()
        });
        
        if (!res.ok) throw new Error('Failed to load analytics');
        
        const data = await res.json();
        displayAnalyticsDashboard(data);
    } catch (err) {
        showToast('Failed to load analytics', 'error');
        console.error(err);
    }
}

function displayAnalyticsDashboard(data) {
    // Learning Hours Chart
    if (data.learning_hours && document.getElementById('chart-learning-hours')) {
        new Chart(document.getElementById('chart-learning-hours'), {
            type: 'bar',
            data: {
                labels: data.learning_hours.map(p => p.label),
                datasets: [{
                    label: 'Hours',
                    data: data.learning_hours.map(p => p.value),
                    backgroundColor: 'rgba(79, 70, 229, 0.5)',
                    borderColor: 'rgba(79, 70, 229, 1)',
                    borderRadius: 6
                }]
            },
            options: { responsive: true, maintainAspectRatio: true }
        });
    }

    // Weekly Progress Chart
    if (data.weekly_progress && document.getElementById('chart-weekly')) {
        new Chart(document.getElementById('chart-weekly'), {
            type: 'line',
            data: {
                labels: data.weekly_progress.map(p => p.label),
                datasets: [{
                    label: 'Weekly Hours',
                    data: data.weekly_progress.map(p => p.value),
                    borderColor: 'rgba(79, 70, 229, 1)',
                    backgroundColor: 'rgba(79, 70, 229, 0.1)',
                    tension: 0.4,
                    fill: true
                }]
            },
            options: { responsive: true, maintainAspectRatio: true }
        });
    }

    // Monthly Progress Chart
    if (data.monthly_progress && document.getElementById('chart-monthly')) {
        new Chart(document.getElementById('chart-monthly'), {
            type: 'line',
            data: {
                labels: data.monthly_progress.map(p => p.label),
                datasets: [{
                    label: 'Monthly Hours',
                    data: data.monthly_progress.map(p => p.value),
                    borderColor: 'rgba(16, 185, 129, 1)',
                    backgroundColor: 'rgba(16, 185, 129, 0.1)',
                    tension: 0.4,
                    fill: true
                }]
            },
            options: { responsive: true, maintainAspectRatio: true }
        });
    }

    // Activity Heatmap
    if (data.activity_calendar) {
        displayActivityHeatmap(data.activity_calendar);
    }

    // Streaks
    if (data.streaks) {
        document.getElementById('current-streak').textContent = data.streaks.current;
        document.getElementById('longest-streak').textContent = data.streaks.longest;
    }

    // Top Skills
    if (data.top_skills) {
        const topSkillsList = document.getElementById('top-skills-list');
        topSkillsList.innerHTML = data.top_skills.map(skill => `
            <div class="skill-item">
                <div class="skill-item-name">${skill.name}</div>
                <div class="skill-item-value">${skill.hours}h (${Math.round(skill.progress)}%)</div>
            </div>
        `).join('');
    }
}

function displayActivityHeatmap(days) {
    const heatmapEl = document.getElementById('activity-heatmap');
    if (!heatmapEl) return;
    
    heatmapEl.innerHTML = days.map(day => {
        const title = `${day.date}: ${day.hours}h`;
        return `<div class="heatmap-day level-${day.level}" title="${title}"></div>`;
    }).join('');
}

// ============================================
// PHASE 3: CAREER READINESS
// ============================================

async function loadCareerReadiness() {
    if (!authToken) return;
    
    try {
        const res = await fetch(`${API}/career-readiness`, {
            headers: getAuthHeaders()
        });
        
        if (!res.ok) throw new Error('Failed to load career readiness');
        
        const data = await res.json();
        displayCareerReadiness(data);
    } catch (err) {
        showToast('Failed to load career readiness', 'error');
        console.error(err);
    }
}

function displayCareerReadiness(readinessScores) {
    const grid = document.getElementById('readiness-grid');
    grid.innerHTML = readinessScores.map(score => `
        <div class="readiness-card">
            <h3>${score.track}</h3>
            <div class="readiness-score">
                <div class="readiness-circle" style="--score-angle: ${(score.score / 100) * 360}deg">
                    ${score.score}%
                </div>
                <div class="readiness-details">
                    <h4>Readiness Score</h4>
                    <p>${score.progress_percent}% Complete</p>
                </div>
            </div>
            <div class="readiness-skills">
                <h4>✓ Matched Skills (${score.matched_skills.length})</h4>
                <div class="skill-tags">
                    ${score.matched_skills.map(s => `<span class="skill-tag">${s}</span>`).join('')}
                </div>
            </div>
            <div class="readiness-skills">
                <h4>→ Skills to Learn (${score.missing_skills.length})</h4>
                <div class="skill-tags">
                    ${score.missing_skills.map(s => `<span class="skill-tag missing">${s}</span>`).join('')}
                </div>
            </div>
        </div>
    `).join('');
}

// ============================================
// PHASE 4: JOB MATCHING
// ============================================

async function loadJobMatches() {
    if (!authToken) return;
    
    try {
        const res = await fetch(`${API}/job-matches`, {
            headers: getAuthHeaders()
        });
        
        if (!res.ok) throw new Error('Failed to load job matches');
        
        const data = await res.json();
        displayJobMatches(data);
    } catch (err) {
        showToast('Failed to load job matches', 'error');
        console.error(err);
    }
}

function displayJobMatches(matches) {
    const grid = document.getElementById('job-matches-grid');
    grid.innerHTML = matches.map(match => `
        <div class="job-match-card">
            <h3>${match.role}</h3>
            <div class="job-readiness">
                <div class="job-readiness-value">${match.readiness_score}%</div>
                <div class="job-readiness-label">
                    <small>Role Readiness</small>
                    <span>
                        ${match.readiness_score >= 70 ? '🚀 Ready' : match.readiness_score >= 40 ? '📈 Building' : '🌱 Starting'}
                    </span>
                </div>
            </div>
            
            <div class="job-skills-section">
                <h4><span>${match.matched_skills.length}</span> Matched</h4>
                <div class="job-skills-list">
                    ${match.matched_skills.map(s => `<span class="job-skill">${s}</span>`).join('')}
                </div>
            </div>
            
            <div class="job-skills-section">
                <h4><span>${match.missing_skills.length}</span> Missing</h4>
                <div class="job-skills-list">
                    ${match.missing_skills.map(s => `<span class="job-skill missing">${s}</span>`).join('')}
                </div>
            </div>
            
            <div class="job-recommendation">
                <strong>Next Step:</strong> ${match.recommended_next}
            </div>
        </div>
    `).join('');
}

// ============================================
// PHASE 5: PROJECTS & CERTIFICATES
// ============================================

async function loadProjects() {
    if (!authToken) return;
    
    try {
        const res = await fetch(`${API}/projects`, {
            headers: getAuthHeaders()
        });
        
        if (!res.ok) throw new Error('Failed to load projects');
        
        const data = await res.json();
        displayProjects(data);
    } catch (err) {
        showToast('Failed to load projects', 'error');
        console.error(err);
    }
}

function displayProjects(projects) {
    const grid = document.getElementById('projects-grid');
    if (!projects || projects.length === 0) {
        grid.innerHTML = '<div class="empty-state" style="grid-column: 1/-1;"><h3>No projects yet</h3><p>Add your projects to showcase your work.</p></div>';
        return;
    }
    
    grid.innerHTML = projects.map(p => `
        <div class="project-card">
            <h3>${p.title}</h3>
            <p class="project-desc">${p.description || 'No description'}</p>
            ${p.impact ? `<div class="project-impact"><label>Impact</label><p>${p.impact}</p></div>` : ''}
            ${p.technologies && p.technologies.length ? `
                <div class="project-tech">
                    <label>Technologies</label>
                    <div class="tech-tags">
                        ${p.technologies.map(t => `<span class="tech-tag">${t}</span>`).join('')}
                    </div>
                </div>
            ` : ''}
            <div class="project-meta">
                ${p.link ? `<a href="${p.link}" target="_blank" class="project-link">View Project →</a>` : ''}
                ${p.duration_months ? `<span>${p.duration_months} months</span>` : ''}
            </div>
            <div class="card-actions">
                <button class="btn-delete" onclick="deleteProject(${p.id})">Delete</button>
            </div>
        </div>
    `).join('');
}

async function loadCertificates() {
    if (!authToken) return;
    
    try {
        const res = await fetch(`${API}/certificates`, {
            headers: getAuthHeaders()
        });
        
        if (!res.ok) throw new Error('Failed to load certificates');
        
        const data = await res.json();
        displayCertificates(data);
    } catch (err) {
        showToast('Failed to load certificates', 'error');
        console.error(err);
    }
}

function displayCertificates(certificates) {
    const grid = document.getElementById('certificates-grid');
    if (!certificates || certificates.length === 0) {
        grid.innerHTML = '<div class="empty-state" style="grid-column: 1/-1;"><h3>No certificates yet</h3><p>Add your certifications and credentials.</p></div>';
        return;
    }
    
    grid.innerHTML = certificates.map(c => `
        <div class="certificate-card">
            <h3>${c.name}</h3>
            <p class="certificate-issuer">${c.issuer || 'Certification'}</p>
            ${c.skills_covered && c.skills_covered.length ? `
                <div class="project-tech">
                    <label>Skills Covered</label>
                    <div class="tech-tags">
                        ${c.skills_covered.map(t => `<span class="tech-tag">${t}</span>`).join('')}
                    </div>
                </div>
            ` : ''}
            <div class="certificate-dates">
                <label>Dates</label>
                <div class="certificate-meta">
                    <span>Issued: ${c.issue_date}</span>
                    ${c.expiry_date ? `<span>Expires: ${c.expiry_date}</span>` : '<span>No expiry</span>'}
                </div>
            </div>
            ${c.credential_url ? `
                <p style="margin-top: 1rem;">
                    <a href="${c.credential_url}" target="_blank" class="project-link">Verify →</a>
                </p>
            ` : ''}
            <div class="card-actions">
                <button class="btn-delete" onclick="deleteCertificate(${c.id})">Delete</button>
            </div>
        </div>
    `).join('');
}

// Modal functions
function openProjectModal() {
    document.getElementById('project-modal').classList.add('active');
}

function closeProjectModal() {
    document.getElementById('project-modal').classList.remove('active');
    document.getElementById('project-form').reset();
}

function openCertificateModal() {
    document.getElementById('certificate-modal').classList.add('active');
}

function closeCertificateModal() {
    document.getElementById('certificate-modal').classList.remove('active');
    document.getElementById('certificate-form').reset();
}

async function submitProject(e) {
    e.preventDefault();
    
    const project = {
        title: document.getElementById('project-title').value,
        description: document.getElementById('project-description').value,
        technologies: document.getElementById('project-technologies').value.split(',').map(t => t.trim()),
        link: document.getElementById('project-link').value,
        duration_months: parseInt(document.getElementById('project-duration').value) || 0,
        completion_date: document.getElementById('project-completion-date').value,
        impact: document.getElementById('project-impact').value
    };
    
    try {
        const res = await fetch(`${API}/projects`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify(project)
        });
        
        if (!res.ok) throw new Error('Failed to create project');
        
        closeProjectModal();
        loadProjects();
        showToast('Project added successfully!');
    } catch (err) {
        showToast('Failed to add project', 'error');
        console.error(err);
    }
}

async function submitCertificate(e) {
    e.preventDefault();
    
    const cert = {
        name: document.getElementById('cert-name').value,
        issuer: document.getElementById('cert-issuer').value,
        credential_id: document.getElementById('cert-credential-id').value,
        credential_url: document.getElementById('cert-credential-url').value,
        issue_date: document.getElementById('cert-issue-date').value,
        expiry_date: document.getElementById('cert-expiry-date').value,
        skills_covered: document.getElementById('cert-skills').value.split(',').map(s => s.trim())
    };
    
    try {
        const res = await fetch(`${API}/certificates`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify(cert)
        });
        
        if (!res.ok) throw new Error('Failed to create certificate');
        
        closeCertificateModal();
        loadCertificates();
        showToast('Certificate added successfully!');
    } catch (err) {
        showToast('Failed to add certificate', 'error');
        console.error(err);
    }
}

async function deleteProject(id) {
    if (!confirm('Delete this project?')) return;
    
    try {
        const res = await fetch(`${API}/projects/${id}`, {
            method: 'DELETE',
            headers: getAuthHeaders()
        });
        
        if (!res.ok) throw new Error('Failed to delete project');
        loadProjects();
        showToast('Project deleted');
    } catch (err) {
        showToast('Failed to delete project', 'error');
    }
}

async function deleteCertificate(id) {
    if (!confirm('Delete this certificate?')) return;
    
    try {
        const res = await fetch(`${API}/certificates/${id}`, {
            method: 'DELETE',
            headers: getAuthHeaders()
        });
        
        if (!res.ok) throw new Error('Failed to delete certificate');
        loadCertificates();
        showToast('Certificate deleted');
    } catch (err) {
        showToast('Failed to delete certificate', 'error');
    }
}

// ============================================
// PHASE 5: USER PROFILE
// ============================================

async function loadProfile() {
    if (!authToken) return;
    
    try {
        const res = await fetch(`${API}/profile`, {
            headers: getAuthHeaders()
        });
        
        if (!res.ok) throw new Error('Failed to load profile');
        
        const data = await res.json();
        populateProfileForm(data);
    } catch (err) {
        console.error(err);
    }
}

function populateProfileForm(profile) {
    if (profile.headline) document.getElementById('profile-headline').value = profile.headline;
    if (profile.bio) document.getElementById('profile-bio').value = profile.bio;
    if (profile.location) document.getElementById('profile-location').value = profile.location;
    if (profile.website) document.getElementById('profile-website').value = profile.website;
    if (profile.github_url) document.getElementById('profile-github').value = profile.github_url;
    if (profile.linkedin_url) document.getElementById('profile-linkedin').value = profile.linkedin_url;
    if (profile.visibility) document.getElementById('profile-visibility').value = profile.visibility;
}

async function submitProfile(e) {
    e.preventDefault();
    
    const profile = {
        headline: document.getElementById('profile-headline').value,
        bio: document.getElementById('profile-bio').value,
        location: document.getElementById('profile-location').value,
        website: document.getElementById('profile-website').value,
        github_url: document.getElementById('profile-github').value,
        linkedin_url: document.getElementById('profile-linkedin').value,
        visibility: document.getElementById('profile-visibility').value
    };
    
    try {
        const res = await fetch(`${API}/profile`, {
            method: 'PUT',
            headers: getAuthHeaders(),
            body: JSON.stringify(profile)
        });
        
        if (!res.ok) throw new Error('Failed to save profile');
        
        showToast('Profile updated successfully!');
    } catch (err) {
        showToast('Failed to save profile', 'error');
        console.error(err);
    }
}

async function generateProfessionalProfile() {
    try {
        const res = await fetch(`${API}/profile/professional`, {
            headers: getAuthHeaders()
        });
        
        if (!res.ok) throw new Error('Failed to generate profile');
        
        const data = await res.json();
        displayProfessionalProfile(data);
    } catch (err) {
        showToast('Failed to generate profile', 'error');
        console.error(err);
    }
}

function displayProfessionalProfile(profile) {
    const summaryEl = document.getElementById('professional-summary');
    
    let html = `
        <div>
            <h4>Summary</h4>
            <p>${profile.summary}</p>
            
            <h4>Strength Areas</h4>
            <ul>
                ${(profile.strength_areas || []).map(a => `<li>${a}</li>`).join('')}
            </ul>
            
            <h4>Growth Areas</h4>
            <ul>
                ${(profile.growth_areas || []).map(a => `<li>${a}</li>`).join('')}
            </ul>
            
            <h4>Recommendations</h4>
            <ul>
                ${(profile.recommendations || []).map(r => `<li>${r}</li>`).join('')}
            </ul>
    `;
    
    if (profile.readiness_scores && Object.keys(profile.readiness_scores).length > 0) {
        html += `
            <h4>Readiness Scores</h4>
            <ul>
                ${Object.entries(profile.readiness_scores).map(([track, score]) => `<li>${track}: ${score}%</li>`).join('')}
            </ul>
        `;
    }
    
    html += '</div>';
    summaryEl.innerHTML = html;
}

// ============================================
// TAB SWITCHING
// ============================================

function switchTab(tab) {
    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(el => el.classList.remove('active'));
    document.querySelectorAll('.tab-btn').forEach(el => el.classList.remove('active'));
    
    // Show selected tab
    const tabEl = document.getElementById(`content-${tab}`);
    const btnEl = document.getElementById(`tab-${tab}`);
    if (tabEl) tabEl.classList.add('active');
    if (btnEl) btnEl.classList.add('active');
    
    // Load data for tab
    if (tab === 'dashboard') loadAnalytics();
    else if (tab === 'readiness') loadCareerReadiness();
    else if (tab === 'jobs') loadJobMatches();
    else if (tab === 'skills') loadSkills();
    else if (tab === 'portfolio') {
        loadProjects();
        loadCertificates();
    } else if (tab === 'profile') loadProfile();
}

function switchPortfolioTab(tab) {
    document.querySelectorAll('.portfolio-section').forEach(el => el.classList.remove('active'));
    document.querySelectorAll('.portfolio-tab-btn').forEach(el => el.classList.remove('active'));
    
    const sectionEl = document.getElementById(`portfolio-${tab}`);
    const btnEl = event.currentTarget;
    
    if (sectionEl) sectionEl.classList.add('active');
    if (btnEl) btnEl.classList.add('active');
}

// ============================================
// INITIALIZE ON LOAD
// ============================================

document.addEventListener('DOMContentLoaded', () => {
    if (authToken) {
        hideLoginPage();
        loadDashboard();
    } else {
        showLoginPage();
    }
    
    // Attach form handlers
    if (document.getElementById('profile-form')) {
        document.getElementById('profile-form').addEventListener('submit', submitProfile);
    }
    if (document.getElementById('project-form')) {
        document.getElementById('project-form').addEventListener('submit', submitProject);
    }
    if (document.getElementById('certificate-form')) {
        document.getElementById('certificate-form').addEventListener('submit', submitCertificate);
    }
});
