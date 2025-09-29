// Configuration
const AUTH_BASE_URL = 'http://127.0.0.1:8081';
const API_BASE_URL = 'http://127.0.0.1:8081/api';
let currentUser = null;
let isAuthenticated = false;

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
});

function initializeApp() {
    // Check authentication status first
    checkAuthStatus();
    
    // Check if auth service is running
    checkAuthServiceStatus();
    
    // Initialize UI
    initializeUI();
    
    logActivity('Application initialized. Please login to continue.');
    
    // Test connection to auth service
    testAuthServiceConnection();
}

// Test connection to auth service
async function testAuthServiceConnection() {
    try {
        console.log('Testing connection to auth service...');
        const response = await fetch(`${AUTH_BASE_URL}/health`);
        if (response.ok) {
            const data = await response.json();
            console.log('Auth service connection test successful:', data);
            logActivity('✅ Auth service connection test successful', 'success');
        } else {
            throw new Error(`Auth service returned status: ${response.status}`);
        }
    } catch (error) {
        console.error('Auth service connection test failed:', error);
        logActivity(`❌ Auth service connection test failed: ${error.message}`, 'error');
    }
}

// Authentication functions
async function checkAuthStatus() {
    try {
        const response = await makeAuthCall('/auth/me');
        if (response.user_id) {
            currentUser = {
                id: response.user_id,
                username: response.username,
                role: response.role
            };
            isAuthenticated = true;
            showSection('documents');
            updateUserDisplay();
            logActivity(`✅ Welcome back, ${currentUser.username}!`, 'success');
        } else {
            showLoginForm();
        }
    } catch (error) {
        showLoginForm();
        logActivity('Please login to continue', 'info');
    }
}

async function login(username, password) {
    try {
        const response = await makeAuthCall('/auth/login', 'POST', {
            username: username,
            password: password
        });
        
        if (response.success) {
            currentUser = response.user;
            isAuthenticated = true;
            hideLoginForm();
            showSection('documents');
            updateUserDisplay();
            logActivity(`✅ Login successful! Welcome, ${currentUser.username}`, 'success');
            loadDocuments();
        } else {
            showLoginError(response.message);
            logActivity(`❌ Login failed: ${response.message}`, 'error');
        }
    } catch (error) {
        showLoginError('Login failed. Please check your credentials.');
        logActivity(`❌ Login error: ${error.message}`, 'error');
    }
}

async function logout() {
    try {
        await makeAuthCall('/auth/logout', 'POST');
        currentUser = null;
        isAuthenticated = false;
        showLoginForm();
        logActivity('Logged out successfully', 'info');
    } catch (error) {
        logActivity('Logout error', 'error');
    }
}

// API call functions
async function makeAuthCall(endpoint, method = 'GET', data = null) {
    try {
        console.log(`Making ${method} request to ${AUTH_BASE_URL}${endpoint}`);
        const config = {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include' // Include cookies for session management
        };
        
        if (data) {
            config.body = JSON.stringify(data);
            console.log('Request data:', data);
        }
        
        const response = await fetch(`${AUTH_BASE_URL}${endpoint}`, config);
        console.log('Response status:', response.status);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const result = await response.json();
        console.log('Response data:', result);
        return result;
    } catch (error) {
        console.error('Auth API call failed:', error);
        console.error('Error details:', {
            message: error.message,
            name: error.name,
            stack: error.stack
        });
        throw error;
    }
}

async function makeApiCall(endpoint, method = 'GET', data = null) {
    if (!isAuthenticated) {
        throw new Error('Authentication required');
    }
    
    try {
        const config = {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include' // Include cookies for session management
        };
        
        if (data) {
            config.body = JSON.stringify(data);
        }
        
        const response = await fetch(`${API_BASE_URL}${endpoint}`, config);
        
        if (response.status === 401) {
            // Session expired
            isAuthenticated = false;
            currentUser = null;
            showLoginForm();
            throw new Error('Session expired. Please login again.');
        }
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error('API call failed:', error);
        logActivity(`API Error: ${error.message}`, 'error');
        throw error;
    }
}

// Service status checks
async function checkAuthServiceStatus() {
    try {
        const response = await fetch(`${AUTH_BASE_URL}/health`);
        if (response.ok) {
            logActivity('✅ Connected to Auth Service successfully', 'success');
            updateConnectionStatus(true);
        } else {
            throw new Error('Auth service unhealthy');
        }
    } catch (error) {
        logActivity('❌ Failed to connect to Auth Service. Make sure it\'s running on port 8081.', 'error');
        updateConnectionStatus(false);
    }
}

// UI Management
function initializeUI() {
    // Setup login form
    setupLoginForm();
    
    // Setup navigation
    setupNavigation();
    
    // Setup logout button
    setupLogoutButton();
    
    // Initially hide main content if not authenticated
    if (!isAuthenticated) {
        showLoginForm();
    }
}

function setupLoginForm() {
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            await login(username, password);
        });
    }
}

function setupNavigation() {
    // Navigation button event listeners
    document.getElementById('nav-documents').addEventListener('click', () => {
        if (isAuthenticated) {
            showSection('documents');
            loadDocuments();
        }
    });
    
    document.getElementById('nav-access').addEventListener('click', () => {
        if (isAuthenticated) {
            showSection('access-control');
        }
    });
    
    document.getElementById('nav-test').addEventListener('click', () => {
        if (isAuthenticated) {
            showSection('test-authorization');
        }
    });
}

function setupLogoutButton() {
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', logout);
    }
}

function showLoginForm() {
    document.getElementById('login-section').style.display = 'block';
    document.getElementById('main-content').style.display = 'none';
    document.getElementById('user-info').style.display = 'none';
}

function hideLoginForm() {
    document.getElementById('login-section').style.display = 'none';
    document.getElementById('main-content').style.display = 'block';
    document.getElementById('user-info').style.display = 'block';
}

function showLoginError(message) {
    const errorDiv = document.getElementById('login-error');
    if (errorDiv) {
        errorDiv.textContent = message;
        errorDiv.style.display = 'block';
    }
}

function updateUserDisplay() {
    if (currentUser) {
        const userInfo = document.getElementById('user-info');
        const userDisplay = document.getElementById('current-user');
        const roleDisplay = document.getElementById('current-role');
        
        if (userDisplay) userDisplay.textContent = currentUser.username;
        if (roleDisplay) roleDisplay.textContent = currentUser.role;
        if (userInfo) userInfo.style.display = 'block';
    }
}

// Document Management
async function loadDocuments() {
    try {
        const response = await makeAuthCall('/documents');
        const documentsContainer = document.getElementById('documents-list');
        
        if (response.documents && response.documents.length > 0) {
            documentsContainer.innerHTML = response.documents.map(doc => `
                <div class="document-item">
                    <h3>${doc}</h3>
                    <div class="document-actions">
                        <button onclick="openDocument('${doc}')" class="btn-primary">Open</button>
                        <button onclick="editDocument('${doc}')" class="btn-secondary">Edit</button>
                        <button onclick="shareDocument('${doc}')" class="btn-tertiary">Share</button>
                    </div>
                </div>
            `).join('');
        } else {
            documentsContainer.innerHTML = '<p>No documents available for your role.</p>';
        }
        
        logActivity(`📄 Loaded ${response.documents.length} documents`, 'info');
    } catch (error) {
        logActivity(`Failed to load documents: ${error.message}`, 'error');
    }
}

async function openDocument(docName) {
    try {
        const response = await makeAuthCall(`/documents/${docName}/access?permission=viewer`);
        
        if (response.authorized) {
            logActivity(`📖 Opened document: ${docName}`, 'success');
            showDocumentContent(docName, 'Viewing document content...');
        } else {
            logActivity(`❌ Access denied: Cannot view ${docName}`, 'error');
            showAccessDenied();
        }
    } catch (error) {
        logActivity(`Error opening document: ${error.message}`, 'error');
    }
}

async function editDocument(docName) {
    try {
        const response = await makeAuthCall(`/documents/${docName}/access?permission=editor`);
        
        if (response.authorized) {
            logActivity(`✏️ Editing document: ${docName}`, 'success');
            showDocumentContent(docName, 'Editing document content...', true);
        } else {
            logActivity(`❌ Access denied: Cannot edit ${docName}`, 'error');
            showAccessDenied();
        }
    } catch (error) {
        logActivity(`Error editing document: ${error.message}`, 'error');
    }
}

function shareDocument(docName) {
    if (currentUser.role !== 'admin') {
        logActivity(`❌ Access denied: Only admins can share documents`, 'error');
        return;
    }
    
    // Show share modal (simplified for demo)
    logActivity(`📤 Share functionality for ${docName} (Admin only)`, 'info');
}

// Authorization Testing
async function testAuthorization() {
    const user = document.getElementById('test-user').value.trim();
    const document = document.getElementById('test-document').value.trim();
    const permission = document.getElementById('test-permission').value;
    
    if (!user || !document) {
        showTestError('Please fill in all fields');
        return;
    }
    
    try {
        const response = await makeApiCall(`/acl/check?object=doc:${document}&relation=${permission}&user=user:${user}`);
        
        const resultText = response.authorized ? 
            `✅ AUTHORIZED: ${user} has ${permission} access to ${document}` :
            `❌ DENIED: ${user} does NOT have ${permission} access to ${document}`;
        
        const resultClass = response.authorized ? 'result-success' : 'result-error';
        
        showTestResult(resultText, resultClass);
        logActivity(`🔐 Authorization test: ${user} -> ${document} (${permission}): ${response.authorized ? 'ALLOWED' : 'DENIED'}`, response.authorized ? 'success' : 'error');
        
    } catch (error) {
        showTestResult(`Error: ${error.message}`, 'result-error');
        logActivity(`🔐 Authorization test error: ${error.message}`, 'error');
    }
}

// ACL Management (Admin only)
async function createACL() {
    if (currentUser.role !== 'admin') {
        logActivity(`❌ Access denied: Only admins can manage ACLs`, 'error');
        return;
    }
    
    const object = document.getElementById('acl-object').value.trim();
    const relation = document.getElementById('acl-relation').value.trim();
    const user = document.getElementById('acl-user').value.trim();
    
    if (!object || !relation || !user) {
        showTestError('Please fill in all ACL fields');
        return;
    }
    
    try {
        const response = await makeApiCall('/acl', 'POST', {
            object: object,
            relation: relation,
            user: user
        });
        
        logActivity(`✅ ACL created: ${object}#${relation}@${user}`, 'success');
        clearACLForm();
    } catch (error) {
        logActivity(`❌ Failed to create ACL: ${error.message}`, 'error');
    }
}

// Utility functions
function showSection(sectionName) {
    // Hide all sections
    const sections = ['documents', 'access-control', 'test-authorization'];
    sections.forEach(section => {
        const element = document.getElementById(section);
        if (element) {
            element.style.display = 'none';
        }
    });
    
    // Show selected section
    const targetSection = document.getElementById(sectionName);
    if (targetSection) {
        targetSection.style.display = 'block';
    }
    
    // Update navigation
    updateNavigationState(sectionName);
}

function updateNavigationState(activeSection) {
    const navButtons = document.querySelectorAll('.nav-btn');
    navButtons.forEach(btn => {
        btn.classList.remove('active');
    });
    
    const activeBtn = document.getElementById(`nav-${activeSection}`);
    if (activeBtn) {
        activeBtn.classList.add('active');
    }
}

function logActivity(message, type = 'info') {
    const timestamp = new Date().toLocaleTimeString();
    const logEntry = document.createElement('div');
    logEntry.className = `log-entry log-${type}`;
    logEntry.innerHTML = `<span class="timestamp">[${timestamp}]</span> ${message}`;
    
    const logContainer = document.getElementById('activity-log');
    if (logContainer) {
        logContainer.appendChild(logEntry);
        logContainer.scrollTop = logContainer.scrollHeight;
    }
    
    console.log(`[${timestamp}] ${message}`);
}

function updateConnectionStatus(isConnected) {
    const statusIndicator = document.createElement('span');
    statusIndicator.className = `status-indicator ${isConnected ? 'status-online' : 'status-offline'}`;
    
    const statusText = document.createElement('span');
    statusText.textContent = isConnected ? 'Connected to Auth Service' : 'Auth Service Offline';
    
    const header = document.querySelector('header p');
    if (header) {
        header.innerHTML = '';
        header.appendChild(statusIndicator);
        header.appendChild(statusText);
    }
}

function showDocumentContent(docName, content, isEditing = false) {
    const modal = document.createElement('div');
    modal.className = 'modal';
    modal.innerHTML = `
        <div class="modal-content">
            <div class="modal-header">
                <h3>${isEditing ? 'Editing' : 'Viewing'}: ${docName}</h3>
                <button onclick="this.closest('.modal').remove()" class="close-btn">&times;</button>
            </div>
            <div class="modal-body">
                ${isEditing ? 
                    `<textarea rows="10" cols="50">${content}</textarea>` : 
                    `<p>${content}</p>`
                }
            </div>
            <div class="modal-footer">
                ${isEditing ? 
                    '<button onclick="saveDocument()" class="btn-primary">Save</button>' : 
                    ''
                }
                <button onclick="this.closest(\'.modal\').remove()" class="btn-secondary">Close</button>
            </div>
        </div>
    `;
    document.body.appendChild(modal);
}

function showAccessDenied() {
    const modal = document.createElement('div');
    modal.className = 'modal';
    modal.innerHTML = `
        <div class="modal-content">
            <div class="modal-header">
                <h3>Access Denied</h3>
                <button onclick="this.closest('.modal').remove()" class="close-btn">&times;</button>
            </div>
            <div class="modal-body">
                <p>❌ You don't have permission to access this resource.</p>
                <p>Current user: <strong>${currentUser.username}</strong> (${currentUser.role})</p>
            </div>
            <div class="modal-footer">
                <button onclick="this.closest('.modal').remove()" class="btn-secondary">Close</button>
            </div>
        </div>
    `;
    document.body.appendChild(modal);
}

function showTestResult(text, className) {
    const resultDiv = document.getElementById('test-result');
    if (resultDiv) {
        resultDiv.textContent = text;
        resultDiv.className = className;
        resultDiv.style.display = 'block';
    }
}

function showTestError(message) {
    const errorDiv = document.getElementById('test-error');
    if (errorDiv) {
        errorDiv.textContent = message;
        errorDiv.style.display = 'block';
    }
}

function clearACLForm() {
    const fields = ['acl-object', 'acl-relation', 'acl-user'];
    fields.forEach(field => {
        const element = document.getElementById(field);
        if (element) element.value = '';
    });
}