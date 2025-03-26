# ChatVui System Design

## Functional Requirements
- one-one
- Group chat 
- Push notifications for new messages
- Online/offline status indicators
- Read messages
- Message history and search

## Non-Functional Requirements
- Low latency (<100ms for message delivery)
- High availability (99.9% uptime)
- High reliability (no message loss)
- Data persistence for chat history

## System Architecture

### Backend Components
currentlt it don't store user info
1. **API Gateway** (feature)
   - Routes requests to appropriate services
   - Handles authentication/authorization

2. **User Service** (feature)
   - User registration and management
   - User profile and preferences
   - Online status tracking

3. **Chat Service**
   - Message handling and routing
   - Supports one-to-one and group chats
   - Real-time communication via WebSockets

4. **Notification Service**
   - Push notifications for new messages
   - Email notifications for important events

5. **Storage Services**
   - MySQL for user data and relationships
   - Redis for caching and real-time features
   - Object storage for media files

### Frontend Components
1. **Web Application (React)**
   - WebSocket connection for real-time updates
## Database Schema

### Users Table
```sql
CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  password_hash VARCHAR(256) NOT NULL,
  email VARCHAR(100) UNIQUE,
  profile_pic_url VARCHAR(255),
  last_online TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Messages Table
```sql
CREATE TABLE messages (
  id INT AUTO_INCREMENT PRIMARY KEY,
  sender_id INT NOT NULL,
  conversation_id INT NOT NULL,
  content TEXT NOT NULL,
  media_url VARCHAR(255),
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (sender_id) REFERENCES users(id),
  FOREIGN KEY (conversation_id) REFERENCES conversations(id)
);
```

### Conversations Table
```sql
CREATE TABLE conversations (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100),
  is_group BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Conversation_Users Table
```sql
CREATE TABLE conversation_users (
  conversation_id INT NOT NULL,
  user_id INT NOT NULL,
  is_admin BOOLEAN DEFAULT FALSE,
  joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (conversation_id, user_id),
  FOREIGN KEY (conversation_id) REFERENCES conversations(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);
```

## API Design

### User Management
- `POST /api/users` - Register new user
- `GET /api/users` - Get all users
- `GET /api/users/active` - Get active users
- `GET /api/users/:id` - Get user profile
- `PUT /api/users/:id` - Update user profile

### Chat Operations
- `POST /api/conversations` - Create new conversation
- `GET /api/conversations` - Get user's conversations
- `GET /api/conversations/:id` - Get conversation details
- `POST /api/conversations/:id/messages` - Send message
- `GET /api/conversations/:id/messages` - Get conversation messages
- `PUT /api/messages/:id/read` - Mark message as read

### WebSocket Events
- `connect` - Establish WebSocket connection
- `message` - New message event
- `typing` - User typing indicator
- `read` - Message read receipt
- `user_status` - User online/offline status update

## Security Measures
- TLS for all HTTP/WebSocket connections
- JWT for authentication
- Input validation and sanitization