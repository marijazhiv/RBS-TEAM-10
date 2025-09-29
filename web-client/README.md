# 📁 DocuManager - Web Client

A proof-of-concept web client that demonstrates the Mini-Zanzibar authorization system in action.

## 🎯 Purpose

This web application showcases how Mini-Zanzibar can be integrated into a real-world document management system, providing:
- **Document access control** based on user permissions
- **Real-time authorization checks** for view/edit/share operations
- **User-friendly interface** for testing Mini-Zanzibar functionality
- **Live demonstration** of hierarchical permissions (owner → editor → viewer)

## 🚀 Features

### 📋 **Document Management**
- Upload new documents (automatically grants owner permission)
- View documents (requires viewer permission or higher)
- Edit documents (requires editor permission or higher)
- Share documents with other users (requires owner permission)

### 🔐 **Authorization Testing**
- Test any user's permission for any document
- Real-time feedback on authorization decisions
- Clear success/denial messages with explanations

### 👥 **Share Management**
- Grant permissions to other users
- Three permission levels: Owner, Editor, Viewer
- Hierarchical inheritance (owners can edit and view, editors can view)

### 📊 **Activity Monitoring**
- Live activity log showing all operations
- Connection status to Mini-Zanzibar
- Success/error tracking for all actions

## 🛠️ Technical Integration

### **Mini-Zanzibar API Integration**
The web client communicates with Mini-Zanzibar through REST API calls:

```javascript
// Check authorization
GET /acl/check?object=doc:report1&relation=viewer&user=user:alice

// Grant permission
POST /acl
{
  "object": "doc:report1",
  "relation": "editor",
  "user": "user:bob"
}
```

### **Permission Hierarchy**
Demonstrates the Mini-Zanzibar namespace configuration:
- **Owner**: Full control (can view, edit, share)
- **Editor**: Can view and edit (inherits from owner)
- **Viewer**: Can only view (inherits from editor)

## 🎮 Demo Users

The application includes pre-configured demo users for testing:

- **👤 Alice**: Owner of most documents
- **👤 Bob**: Editor access to some documents
- **👤 Charlie**: Viewer access to some documents  
- **👤 David**: No access (for testing denials)

## 📱 How to Use

### **1. Start Mini-Zanzibar**
Make sure Mini-Zanzibar is running on port 8080:
```bash
cd mini-zanzibar
go run cmd/server/main.go
```

### **2. Open Web Client**
Simply open `index.html` in your web browser.

### **3. Login**
Click on any demo user button to login and start testing.

### **4. Test Features**
- **Documents**: Upload, view, edit, and share documents
- **Auth Test**: Test authorization for any user/document combination
- **Share**: Grant permissions to other users

## 🔍 What This Demonstrates

### **Real-World Authorization**
Shows how Mini-Zanzibar would work in a production document management system:
- Users can only access documents they have permission for
- Permission inheritance works automatically
- Sharing is restricted to document owners

### **API Integration**
Demonstrates proper integration patterns:
- Error handling for API failures
- Real-time permission checking
- User-friendly feedback for authorization decisions

### **Security Model**
Illustrates Zanzibar's security principles:
- Explicit permission grants
- Hierarchical access control
- Granular permission checking

## 🎨 User Interface

### **Modern Design**
- Clean, responsive interface
- Intuitive navigation between sections
- Real-time status indicators
- Activity logging for transparency

### **Accessibility**
- Keyboard shortcuts (Ctrl+1/2/3 for navigation)
- Clear visual feedback for all actions
- Error messages with actionable guidance

## 🔧 Technical Details

### **Frontend Technologies**
- **HTML5**: Semantic structure
- **CSS3**: Modern styling with gradients and animations
- **Vanilla JavaScript**: No dependencies, lightweight
- **Fetch API**: For Mini-Zanzibar communication

### **Authorization Flow**
1. User attempts an action (view/edit/share document)
2. Client sends permission check to Mini-Zanzibar
3. Mini-Zanzibar evaluates ACLs and namespace rules
4. Client receives authorization decision
5. UI updates to show success or denial

### **Error Handling**
- Network failures are handled gracefully
- Clear error messages for users
- Fallback behavior when Mini-Zanzibar is offline

## 🎯 Educational Value

This client serves as a **proof-of-concept** that demonstrates:
- How authorization systems integrate with applications
- The practical benefits of hierarchical permissions
- Real-world usage patterns for access control
- The simplicity of Zanzibar-style authorization

Perfect for understanding how Mini-Zanzibar fits into larger systems while keeping the implementation focused and educational.