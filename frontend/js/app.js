const API = '/api';

// Roadmaps config definition
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
    const stored = localStorage.getItem('skillpulse-theme');
    if (stored) return stored;
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function applyTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('skillpulse-theme', theme);
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
    if (!localStorage.getItem('skillpulse-theme')) {
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
    await loadSettings();
    loadDashboard();
    loadSkills();
});

// API Calls
async function loadSettings() {
    try {
        const res = await fetch(`${API}/settings`);
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
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ key: 'target_role', value: role }),
        });
        if (!res.ok) throw new Error('Failed to save target role setting');
    } catch (err) {
        console.error('Failed to save settings:', err);
    }
}

async function loadDashboard() {
    try {
        const res = await fetch(`${API}/dashboard`);
        dashboard = await res.json();
        renderStats();
    } catch (err) {
        console.error('Failed to load dashboard:', err);
    }
}

async function loadSkills() {
    try {
        const res = await fetch(`${API}/skills`);
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
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to create skill');
    return res.json();
}

async function deleteSkill(id) {
    const res = await fetch(`${API}/skills/${id}`, { method: 'DELETE' });
    if (!res.ok) throw new Error('Failed to delete skill');
    return res.json();
}

async function logSession(skillId, data) {
    const res = await fetch(`${API}/skills/${skillId}/log`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
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

// Close modals on backdrop click
document.querySelectorAll('.modal-backdrop').forEach(el => {
    el.addEventListener('click', (e) => {
        if (e.target === el) {
            el.classList.remove('active');
        }
    });
});
