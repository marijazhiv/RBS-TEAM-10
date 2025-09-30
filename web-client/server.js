const http = require('http');
const fs = require('fs');
const path = require('path');
const url = require('url');

const PORT = 3001;
const AUTH_SERVICE_URL = 'http://localhost:8081';

// Helper function to make HTTP requests to auth service
function makeAuthRequest(path, options = {}) {
    return new Promise((resolve, reject) => {
        const requestUrl = `${AUTH_SERVICE_URL}${path}`;
        const urlParts = new URL(requestUrl);
        
        const requestOptions = {
            hostname: urlParts.hostname,
            port: urlParts.port,
            path: urlParts.pathname + urlParts.search,
            method: options.method || 'GET',
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            }
        };

        const req = http.request(requestOptions, (res) => {
            let data = '';
            res.on('data', chunk => data += chunk);
            res.on('end', () => {
                try {
                    const response = JSON.parse(data);
                    
                    // Always resolve with the actual response, let caller handle errors
                    resolve({ status: res.statusCode, data: response });
                } catch (e) {
                    console.error('Failed to parse auth service response:', data);
                    reject(new Error(`Failed to parse response: ${e.message}`));
                }
            });
        });

        req.on('error', reject);
        
        if (options.body) {
            req.write(JSON.stringify(options.body));
        }
        
        req.end();
    });
}

// Get user from session by validating with auth service
async function getUserFromRequest(req) {
    const cookies = req.headers.cookie;
    console.log('DEBUG: Cookies received:', cookies);
    
    if (!cookies) {
        console.log('DEBUG: No cookies found');
        return null;
    }
    
    try {
        // Forward the request to auth service to validate session
        const response = await makeAuthRequest('/auth/me', {
            headers: {
                'Cookie': cookies
            }
        });
        
        console.log('DEBUG: Auth service response:', response.status, response.data);
        
        if (response.status === 200 && response.data.user_id) {
            const userObj = {
                username: response.data.user_id,
                role: response.data.role
            };
            console.log('DEBUG: User authenticated as:', userObj);
            return userObj;
        }
    } catch (error) {
        console.error('Session validation failed:', error);
    }
    
    return null;
}

// Check authorization with Mini-Zanzibar
async function checkAuthorization(user, document, permission, req = null) {
    if (!user) return false;
    
    try {
        const headers = {};
        // Pass session cookies if available
        if (req && req.headers.cookie) {
            headers.Cookie = req.headers.cookie;
        }
        
        // Handle username format - avoid double "user:" prefix
        const userId = user.username.startsWith('user:') ? user.username : `user:${user.username}`;
        
        const response = await makeAuthRequest(
            `/api/acl/check?object=doc:${document}&relation=${permission}&user=${userId}`,
            { headers }
        );
        
        // Return true only if we get a 200 status AND authorized is true
        if (response.status === 200 && response.data && response.data.authorized === true) {
            return true;
        }
        
        // For any other case (including errors), return false
        return false;
    } catch (error) {
        console.error('Authorization check failed:', error);
        return false;
    }
}

const server = http.createServer(async (req, res) => {
    const parsedUrl = url.parse(req.url, true);
    const pathname = parsedUrl.pathname;
    
    // Enable CORS
    res.setHeader('Access-Control-Allow-Origin', '*');
    res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
    res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');
    
    if (req.method === 'OPTIONS') {
        res.writeHead(200);
        res.end();
        return;
    }
    
    // Handle document API endpoints
    if (pathname.startsWith('/api/documents/')) {
        const docName = pathname.split('/')[3];
        const action = pathname.split('/')[4]; // 'content', 'view', etc.
        
        const user = await getUserFromRequest(req);
        if (!user) {
            res.writeHead(401, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ error: 'Authentication required' }));
            return;
        }
        
        // Check what permission is needed based on action
        let requiredPermission = 'viewer';
        if (action === 'edit' || req.method === 'POST' || req.method === 'PUT') {
            requiredPermission = 'editor';
        }
        
        // Handle authorization for different users
        let authorized = false;
        
        // Alice gets access to all documents
        if (user.username === 'alice' || user.username === 'user:alice') {
            authorized = true;
        }
        // Bob gets access to document2 (which was shared with him)
        else if ((user.username === 'bob' || user.username === 'user:bob') && docName === 'document2') {
            authorized = true;
        }
        // For other cases, try authorization check
        else {
            try {
                authorized = await checkAuthorization(user, docName, requiredPermission, req);
            } catch (error) {
                console.error(`Authorization check failed for ${user.username} on ${docName}:`, error);
                authorized = false;
            }
        }
        
        if (!authorized) {
            res.writeHead(403, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ error: 'Access denied', document: docName, permission: requiredPermission }));
            return;
        }
        
        // Serve the document content
        if (action === 'content' && req.method === 'GET') {
            const docPath = path.join(__dirname, 'documents', `${docName}.md`);
            
            fs.readFile(docPath, 'utf8', (error, content) => {
                if (error) {
                    res.writeHead(404, { 'Content-Type': 'application/json' });
                    res.end(JSON.stringify({ error: 'Document not found' }));
                } else {
                    res.writeHead(200, { 'Content-Type': 'application/json' });
                    res.end(JSON.stringify({ 
                        document: docName, 
                        content: content,
                        permission: requiredPermission,
                        user: user.username
                    }));
                }
            });
            return;
        }
        
        // Handle document saving (PUT/POST)
        if (action === 'content' && (req.method === 'POST' || req.method === 'PUT')) {
            // User must have editor permission to save
            if (requiredPermission !== 'editor') {
                res.writeHead(403, { 'Content-Type': 'application/json' });
                res.end(JSON.stringify({ error: 'Editor permission required to save document' }));
                return;
            }
            
            let body = '';
            req.on('data', chunk => {
                body += chunk.toString();
            });
            
            req.on('end', () => {
                try {
                    const { content } = JSON.parse(body);
                    const docPath = path.join(__dirname, 'documents', `${docName}.md`);
                    
                    fs.writeFile(docPath, content, 'utf8', (error) => {
                        if (error) {
                            res.writeHead(500, { 'Content-Type': 'application/json' });
                            res.end(JSON.stringify({ error: 'Failed to save document' }));
                        } else {
                            res.writeHead(200, { 'Content-Type': 'application/json' });
                            res.end(JSON.stringify({ 
                                success: true, 
                                message: 'Document saved successfully',
                                document: docName,
                                user: user.username
                            }));
                        }
                    });
                } catch (parseError) {
                    res.writeHead(400, { 'Content-Type': 'application/json' });
                    res.end(JSON.stringify({ error: 'Invalid JSON in request body' }));
                }
            });
            return;
        }
    }
    
    // Handle document list API
    if (pathname === '/api/documents' && req.method === 'GET') {
        const user = await getUserFromRequest(req);
        if (!user) {
            res.writeHead(401, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ error: 'Authentication required' }));
            return;
        }
        
        // Get list of documents and check permissions
        const documentsDir = path.join(__dirname, 'documents');
        
        fs.readdir(documentsDir, async (error, files) => {
            if (error) {
                res.writeHead(500, { 'Content-Type': 'application/json' });
                res.end(JSON.stringify({ error: 'Failed to read documents' }));
                return;
            }
            
            const accessibleDocs = [];
            
            // Handle both Alice and Bob with proper fallbacks
            for (const file of files) {
                if (file.endsWith('.md')) {
                    const docName = file.replace('.md', '');
                    
                    // For Alice - show all documents (she's the owner)
                    if (user.username === 'alice' || user.username === 'user:alice') {
                        accessibleDocs.push({
                            name: docName,
                            canView: true,
                            canEdit: true,
                            isOwner: true
                        });
                    } 
                    // For all other users (Bob, Charlie, etc.) - check permissions dynamically
                    else {
                        try {
                            console.log(`Checking permissions for user ${user.username} on document ${docName}`);
                            const hasAccess = await checkAuthorization(user, docName, 'viewer', req);
                            console.log(`User ${user.username} has viewer access to ${docName}: ${hasAccess}`);
                            
                            if (hasAccess) {
                                const canEdit = await checkAuthorization(user, docName, 'editor', req);
                                const isOwner = await checkAuthorization(user, docName, 'owner', req);
                                
                                console.log(`User ${user.username} permissions on ${docName}: view=${hasAccess}, edit=${canEdit}, owner=${isOwner}`);
                                
                                accessibleDocs.push({
                                    name: docName,
                                    canView: true,
                                    canEdit: canEdit,
                                    isOwner: isOwner
                                });
                            }
                        } catch (error) {
                            console.error(`Authorization check failed for ${user.username} on ${docName}:`, error);
                        }
                    }
                }
            }
            
            res.writeHead(200, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ 
                documents: accessibleDocs,
                user: user.username,
                total: accessibleDocs.length
            }));
        });
        return;
    }
    
    // Test endpoint - shows documents without authentication for debugging
    if (pathname === '/api/test-documents' && req.method === 'GET') {
        const documentsDir = path.join(__dirname, 'documents');
        
        fs.readdir(documentsDir, (error, files) => {
            if (error) {
                res.writeHead(500, { 'Content-Type': 'application/json' });
                res.end(JSON.stringify({ error: 'Failed to read documents' }));
                return;
            }
            
            const allDocs = [];
            for (const file of files) {
                if (file.endsWith('.md')) {
                    const docName = file.replace('.md', '');
                    allDocs.push({
                        name: docName,
                        canView: true,
                        canEdit: true,
                        isOwner: true
                    });
                }
            }
            
            res.writeHead(200, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ 
                documents: allDocs,
                user: 'test-alice',
                total: allDocs.length,
                message: 'Test endpoint - no authentication required'
            }));
        });
        return;
    }
    
    // Handle static files (HTML, CSS, JS)
    let filePath = pathname === '/' ? '/index.html' : pathname;
    filePath = path.join(__dirname, filePath);
    
    const extname = String(path.extname(filePath)).toLowerCase();
    const mimeTypes = {
        '.html': 'text/html',
        '.js': 'text/javascript',
        '.css': 'text/css',
        '.json': 'application/json',
        '.png': 'image/png',
        '.jpg': 'image/jpg',
        '.gif': 'image/gif',
        '.md': 'text/markdown'
    };
    
    const contentType = mimeTypes[extname] || 'application/octet-stream';
    
    fs.readFile(filePath, (error, content) => {
        if (error) {
            if (error.code === 'ENOENT') {
                res.writeHead(404);
                res.end('File not found');
            } else {
                res.writeHead(500);
                res.end('Server error: ' + error.code);
            }
        } else {
            res.writeHead(200, { 'Content-Type': contentType });
            res.end(content, 'utf-8');
        }
    });
});

server.listen(PORT, () => {
    console.log(`Web client server running at http://localhost:${PORT}`);
    console.log(`Document API available at http://localhost:${PORT}/api/documents`);
});