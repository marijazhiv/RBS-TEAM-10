// Configuration
const API_BASE_URL = 'http://localhost:8080';
let currentUser = null;
let shareModalDocument = null;

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
});

function initializeApp() {
    // Check if Mini-Zanzibar is running
    checkMiniZanzibarStatus();
    
    // Initialize default data
    initializeDefaultData();
    
    // Show login section by default
    showSection('documents');
    
    logActivity('Application initialized. Please login to continue.');
}

// Mini-Zanzibar API functions
async function makeApiCall(endpoint, method = 'GET', data = null) {
    try {
        const config = {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            }
        };
        
        if (data) {
            config.body = JSON.stringify(data);
        }
        
        const response = await fetch(`${API_BASE_URL}${endpoint}`, config);
        
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

async function checkMiniZanzibarStatus() {
    try {
        await makeApiCall('/health');
        logActivity('Connected to Mini-Zanzibar successfully', 'success');
        updateConnectionStatus(true);
    } catch (error) {
        logActivity('Failed to connect to Mini-Zanzibar. Make sure it\'s running on port 8080.', 'error');
        updateConnectionStatus(false);
    }
}

function updateConnectionStatus(isConnected) {
    const statusIndicator = document.createElement('span');
    statusIndicator.className = `status-indicator ${isConnected ? 'status-online' : 'status-offline'}`;
    
    const statusText = document.createElement('span');
    statusText.textContent = isConnected ? 'Connected to Mini-Zanzibar' : 'Mini-Zanzibar Offline';
    
    const header = document.querySelector('header p');
    header.innerHTML = '';
    header.appendChild(statusIndicator);
    header.appendChild(statusText);
}

// Initialize default data (Demo ACLs)
async function initializeDefaultData() {
    const defaultACLs = [
        { object: 'doc:report1', relation: 'owner', user: 'user:alice' },
        { object: 'doc:report1', relation: 'editor', user: 'user:bob' },
        { object: 'doc:manual2', relation: 'owner', user: 'user:alice' },
        { object: 'doc:manual2', relation: 'viewer', user: 'user:charlie' },
        { object: 'doc:secret3', relation: 'owner', user: 'user:alice' }
    ];

    // Try to create default ACLs
    for (const acl of defaultACLs) {
        try {
            await makeApiCall('/acl', 'POST', acl);
        } catch (error) {
            // Ignore errors for existing ACLs
        }
    }
    
    logActivity('Default demo data initialized');
}

// Authentication functions
function quickLogin(username) {
    currentUser = username;
    document.getElementById('username').textContent = username;
    document.getElementById('login-btn').style.display = 'none';
    document.getElementById('logout-btn').style.display = 'inline-block';
    
    // Hide login section and show documents
    document.getElementById('login-section').style.display = 'none';
    showSection('documents');
    
    logActivity(`👤 Logged in as ${username}`, 'success');
    
    // Update UI based on user permissions
    updateDocumentPermissions();
}

function logout() {
    currentUser = null;
    document.getElementById('username').textContent = 'guest';
    document.getElementById('login-btn').style.display = 'inline-block';
    document.getElementById('logout-btn').style.display = 'none';
    
    // Show login section
    document.getElementById('login-section').style.display = 'block';
    document.getElementById('documents-section').style.display = 'none';
    document.getElementById('auth-test-section').style.display = 'none';
    document.getElementById('share-section').style.display = 'none';
    
    logActivity('👋 Logged out', 'info');
}

// Document operations
async function openDocument(documentId) {
    if (!currentUser) {
        showError('Please login first');
        return;
    }
    
    try {
        const result = await makeApiCall(`/acl/check?object=doc:${documentId}&relation=viewer&user=user:${currentUser}`);
        
        if (result.authorized) {
            logActivity(`${currentUser} opened document: ${documentId}`, 'success');
            showSuccess(`Successfully opened ${documentId}`);
        } else {
            logActivity(`${currentUser} denied access to view: ${documentId}`, 'error');
            showError(`Access denied: You don't have permission to view ${documentId}`);
        }
    } catch (error) {
        showError(`Failed to check permissions: ${error.message}`);
    }
}

async function editDocument(documentId) {
    if (!currentUser) {
        showError('Please login first');
        return;
    }
    
    try {
        const result = await makeApiCall(`/acl/check?object=doc:${documentId}&relation=editor&user=user:${currentUser}`);
        
        if (result.authorized) {
            logActivity(`${currentUser} edited document: ${documentId}`, 'success');
            showSuccess(`Successfully edited ${documentId}`);
        } else {
            logActivity(`${currentUser} denied access to edit: ${documentId}`, 'error');
            showError(`Access denied: You don't have permission to edit ${documentId}`);
        }
    } catch (error) {
        showError(`Failed to check permissions: ${error.message}`);
    }
}

async function uploadDocument() {
    if (!currentUser) {
        showError('Please login first');
        return;
    }
    
    const docName = document.getElementById('doc-name').value.trim();
    if (!docName) {
        showError('Please enter a document name');
        return;
    }
    
    try {
        // Grant owner permission to the current user for the new document
        const aclData = {
            object: `doc:${docName}`,
            relation: 'owner',
            user: `user:${currentUser}`
        };
        
        await makeApiCall('/acl', 'POST', aclData);
        
        logActivity(`${currentUser} uploaded document: ${docName}`, 'success');
        showSuccess(`Successfully uploaded ${docName}. You are now the owner.`);
        
        // Clear the input
        document.getElementById('doc-name').value = '';
        
        // Add document to the list
        addDocumentToList(docName);
        
    } catch (error) {
        showError(`Failed to upload document: ${error.message}`);
    }
}

function addDocumentToList(docName) {
    const container = document.getElementById('documents-container');
    
    const docItem = document.createElement('div');
    docItem.className = 'document-item';
    docItem.setAttribute('data-doc', docName);
    
    docItem.innerHTML = `
        <span class="doc-icon">📄</span>
        <span class="doc-name">${docName}</span>
        <div class="doc-actions">
            <button onclick="openDocument('${docName}')" class="action-btn view">👁️ View</button>
            <button onclick="editDocument('${docName}')" class="action-btn edit">✏️ Edit</button>
            <button onclick="showShareModal('${docName}')" class="action-btn share">🔗 Share</button>
        </div>
    `;
    
    container.appendChild(docItem);
}

// Share functionality
function showShareModal(documentId) {
    if (!currentUser) {
        showError('Please login first');
        return;
    }
    
    shareModalDocument = documentId;
    document.getElementById('modal-doc-name').textContent = documentId;
    document.getElementById('modal-share-user').value = '';
    document.getElementById('modal-share-permission').value = 'viewer';
    document.getElementById('modal-result').innerHTML = '';
    document.getElementById('share-modal').style.display = 'block';
}

function closeShareModal() {
    document.getElementById('share-modal').style.display = 'none';
    shareModalDocument = null;
}

async function executeShare() {
    const shareUser = document.getElementById('modal-share-user').value.trim();
    const sharePermission = document.getElementById('modal-share-permission').value;
    
    if (!shareUser) {
        document.getElementById('modal-result').innerHTML = '<div class="result-container result-error">Please enter a username</div>';
        return;
    }
    
    try {
        // First check if current user can share (needs owner permission)
        const canShare = await makeApiCall(`/acl/check?object=doc:${shareModalDocument}&relation=owner&user=user:${currentUser}`);
        
        if (!canShare.authorized) {
            document.getElementById('modal-result').innerHTML = '<div class="result-container result-error">Access denied: Only owners can share documents</div>';
            return;
        }
        
        // Grant the permission
        const aclData = {
            object: `doc:${shareModalDocument}`,
            relation: sharePermission,
            user: `user:${shareUser}`
        };
        
        await makeApiCall('/acl', 'POST', aclData);
        
        logActivity(`🔗 ${currentUser} shared ${shareModalDocument} with ${shareUser} as ${sharePermission}`, 'success');
        document.getElementById('modal-result').innerHTML = `<div class="result-container result-success">Successfully shared ${shareModalDocument} with ${shareUser} as ${sharePermission}</div>`;
        
        setTimeout(() => {
            closeShareModal();
        }, 2000);
        
    } catch (error) {
        document.getElementById('modal-result').innerHTML = `<div class="result-container result-error">Failed to share: ${error.message}</div>`;
    }
}

// Share management
async function shareDocument() {
    const docName = document.getElementById('share-document').value.trim();
    const userName = document.getElementById('share-user').value.trim();
    const permission = document.getElementById('share-permission').value;
    
    if (!docName || !userName) {
        showError('Please fill in all fields');
        return;
    }
    
    if (!currentUser) {
        showError('Please login first');
        return;
    }
    
    try {
        // First check if current user can share (needs owner permission)
        const canShare = await makeApiCall(`/acl/check?object=doc:${docName}&relation=owner&user=user:${currentUser}`);
        
        if (!canShare.authorized) {
            showError('Access denied: Only owners can share documents');
            return;
        }
        
        // Grant the permission
        const aclData = {
            object: `doc:${docName}`,
            relation: permission,
            user: `user:${userName}`
        };
        
        await makeApiCall('/acl', 'POST', aclData);
        
        logActivity(`${currentUser} shared ${docName} with ${userName} as ${permission}`, 'success');
        showShareSuccess(`Successfully shared ${docName} with ${userName} as ${permission}`);
        
        // Clear form
        document.getElementById('share-document').value = '';
        document.getElementById('share-user').value = '';
        
    } catch (error) {
        showShareError(`Failed to share: ${error.message}`);
    }
}

// Authorization testing
async function testAuthorization() {
    const user = document.getElementById('test-user').value.trim();
    const document = document.getElementById('test-document').value.trim();
    const permission = document.getElementById('test-permission').value;
    
    if (!user || !document) {
        showTestError('Please fill in all fields');
        return;
    }
    
    try {
        const result = await makeApiCall(`/acl/check?object=doc:${document}&relation=${permission}&user=user:${user}`);
        
        const resultText = result.authorized ? 
            `AUTHORIZED: ${user} has ${permission} access to ${document}` :
            `DENIED: ${user} does NOT have ${permission} access to ${document}`;
        
        const resultClass = result.authorized ? 'result-success' : 'result-error';
        
        showTestResult(resultText, resultClass);
        logActivity(`Authorization test: ${user} -> ${document} (${permission}): ${result.authorized ? 'ALLOWED' : 'DENIED'}`, result.authorized ? 'success' : 'error');
        
    } catch (error) {
        showTestError(`Failed to test authorization: ${error.message}`);
    }
}

// Update document permissions UI
async function updateDocumentPermissions() {
    if (!currentUser) return;
    
    const documents = document.querySelectorAll('.document-item');
    
    for (const docElement of documents) {
        const docId = docElement.getAttribute('data-doc');
        const viewBtn = docElement.querySelector('.action-btn.view');
        const editBtn = docElement.querySelector('.action-btn.edit');
        const shareBtn = docElement.querySelector('.action-btn.share');
        
        try {
            // Check view permission
            const canView = await makeApiCall(`/acl/check?object=doc:${docId}&relation=viewer&user=user:${currentUser}`);
            viewBtn.disabled = !canView.authorized;
            
            // Check edit permission
            const canEdit = await makeApiCall(`/acl/check?object=doc:${docId}&relation=editor&user=user:${currentUser}`);
            editBtn.disabled = !canEdit.authorized;
            
            // Check owner permission (for sharing)
            const canShare = await makeApiCall(`/acl/check?object=doc:${docId}&relation=owner&user=user:${currentUser}`);
            shareBtn.disabled = !canShare.authorized;
            
        } catch (error) {
            console.error(`Failed to check permissions for ${docId}:`, error);
        }
    }
}

// Navigation
function showSection(sectionName) {
    // Hide all sections
    document.getElementById('documents-section').style.display = 'none';
    document.getElementById('auth-test-section').style.display = 'none';
    document.getElementById('share-section').style.display = 'none';
    
    // Show selected section
    document.getElementById(`${sectionName}-section`).style.display = 'block';
    
    // Update navigation buttons
    document.querySelectorAll('.nav-btn').forEach(btn => btn.classList.remove('active'));
    event?.target?.classList.add('active');
    
    // Update document permissions if showing documents
    if (sectionName === 'documents' && currentUser) {
        updateDocumentPermissions();
    }
}

// Utility functions
function showError(message) {
    const result = document.getElementById('test-result') || document.createElement('div');
    result.innerHTML = `<div class="result-container result-error">${message}</div>`;
}

function showSuccess(message) {
    const result = document.getElementById('test-result') || document.createElement('div');
    result.innerHTML = `<div class="result-container result-success">${message}</div>`;
}

function showTestResult(message, className) {
    const result = document.getElementById('test-result');
    result.innerHTML = `<div class="result-container ${className}">${message}</div>`;
}

function showTestError(message) {
    showTestResult(message, 'result-error');
}

function showShareSuccess(message) {
    const result = document.getElementById('share-result');
    result.innerHTML = `<div class="result-container result-success">${message}</div>`;
}

function showShareError(message) {
    const result = document.getElementById('share-result');
    result.innerHTML = `<div class="result-container result-error">${message}</div>`;
}

function logActivity(message, type = 'info') {
    const log = document.getElementById('activity-log');
    const timestamp = new Date().toLocaleTimeString();
    
    const logItem = document.createElement('div');
    logItem.className = `log-item ${type}`;
    logItem.textContent = `[${timestamp}] ${message}`;
    
    log.appendChild(logItem);
    log.scrollTop = log.scrollHeight;
    
    // Keep only last 20 log items
    while (log.children.length > 20) {
        log.removeChild(log.firstChild);
    }
}

// Login modal functions (for future enhancement)
function showLoginModal() {
    // For now, just scroll to login section
    document.getElementById('login-section').scrollIntoView({ behavior: 'smooth' });
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('share-modal');
    if (event.target === modal) {
        closeShareModal();
    }
}

// Keyboard shortcuts
document.addEventListener('keydown', function(event) {
    // ESC to close modal
    if (event.key === 'Escape') {
        closeShareModal();
    }
    
    // Ctrl+1, Ctrl+2, Ctrl+3 for navigation
    if (event.ctrlKey) {
        switch(event.key) {
            case '1':
                showSection('documents');
                break;
            case '2':
                showSection('auth-test');
                break;
            case '3':
                showSection('share');
                break;
        }
    }
});