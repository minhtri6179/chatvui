import { useState, useEffect, useRef } from 'react';
import './ChatPage.css';
import webSocketService from '../services/WebSocketService';
import apiService from '../services/ApiService';

interface ChatPageProps {
  username: string;
  ws: WebSocket | null;
}

interface Message {
  id: number | string;
  text: string;
  sender: 'user' | 'system';
  timestamp: Date;
}

const ChatPage: React.FC<ChatPageProps> = ({ username, ws }) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [activeUsers, setActiveUsers] = useState<string[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  
  // Initialize user and websocket connection
  useEffect(() => {
    const initChat = async () => {
      try {
        // Register user
        await apiService.registerUser(username);
        
        // Connect to WebSocket
        await webSocketService.connect(username);
        setIsConnected(true);
        
        // Add welcome message
        setMessages([{
          id: 'welcome',
          text: `Chào mừng ${username} đến với cuộc trò chuyện!`,
          sender: 'system',
          timestamp: new Date()
        }]);
      } catch (error) {
        console.error('Error initializing chat:', error);
        setMessages([{
          id: 'error',
          text: 'Không thể kết nối đến máy chủ. Vui lòng thử lại sau.',
          sender: 'system',
          timestamp: new Date()
        }]);
      }
    };
    
    initChat();
    
    // Cleanup on unmount
    return () => {
      webSocketService.disconnect();
      setIsConnected(false);
    };
  }, [username]);

  // Handle WebSocket messages
  useEffect(() => {
    if (!ws) return;

    const handleMessage = (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      
      if (data.type === 'online_users') {
        setActiveUsers(data.users);
      } else if (data.type === 'user_status') {
        // Update messages with user status change
        setMessages(prev => [...prev, {
          id: Date.now(),
          text: `${data.username} is now ${data.status}`,
          sender: 'system',
          timestamp: new Date()
        }]);
      }
    };

    ws.addEventListener('message', handleMessage);
    return () => ws.removeEventListener('message', handleMessage);
  }, [ws]);
  
  // Auto-scroll to bottom when messages change
  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (newMessage.trim() === '' || !isConnected) return;
    
    // Send message via WebSocket
    webSocketService.sendMessage(newMessage);
    setNewMessage('');
  };

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString('vi-VN', { 
      hour: '2-digit', 
      minute: '2-digit'
    });
  };

  return (
    <div className="chat-container">
      <div className="chat-header">
        <h2>Trò Chuyện</h2>
        <div className="active-users">
          <h3>Người dùng đang online ({activeUsers.length})</h3>
          <ul className="online-users-list">
            {activeUsers.map((user, index) => (
              <li key={index} className="online-user">
                <span className="status-indicator online"></span>
                {user}
              </li>
            ))}
          </ul>
        </div>
      </div>
      
      <div className="messages-container">
        {messages.map(message => (
          <div 
            key={message.id} 
            className={`message ${message.sender === 'user' ? 'user-message' : 'system-message'}`}
          >
            <div className="message-content">
              {message.text}
            </div>
            <div className="message-time">
              {formatTime(message.timestamp)}
            </div>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>
      
      <form className="input-container" onSubmit={handleSendMessage}>
        <input
          type="text"
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Nhập tin nhắn..."
          className="message-input"
          disabled={!isConnected}
        />
        <button 
          type="submit" 
          className="send-button"
          disabled={!isConnected}
        >
          Gửi
        </button>
      </form>
    </div>
  );
}

export default ChatPage; 