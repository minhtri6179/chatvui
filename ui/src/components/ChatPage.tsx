import { useState, useEffect, useRef } from 'react';
import './ChatPage.css';
import webSocketService, { Message as WSMessage } from '../services/WebSocketService';
import apiService, { User } from '../services/ApiService';

interface ChatPageProps {
  username: string;
}

interface Message {
  id: number | string;
  text: string;
  sender: 'user' | 'system';
  timestamp: Date;
}

function ChatPage({ username }: ChatPageProps) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [user, setUser] = useState<User | null>(null);
  const [activeUsers, setActiveUsers] = useState<User[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  
  // Initialize user and websocket connection
  useEffect(() => {
    const initChat = async () => {
      try {
        // Register user
        const registeredUser = await apiService.registerUser(username);
        setUser(registeredUser);
        
        // Connect to WebSocket
        await webSocketService.connect(registeredUser.id);
        setIsConnected(true);
        
        // Add welcome message
        setMessages([{
          id: 'welcome',
          text: `Chào mừng ${username} đến với cuộc trò chuyện!`,
          sender: 'system',
          timestamp: new Date()
        }]);
        
        // Fetch active users
        const users = await apiService.getActiveUsers();
        setActiveUsers(users);
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
  
  // Handler for incoming websocket messages
  useEffect(() => {
    const handleMessage = (wsMessage: WSMessage) => {
      const message: Message = {
        id: wsMessage.id,
        text: wsMessage.text,
        sender: wsMessage.type === 'user' && wsMessage.user_id === user?.id ? 'user' : 'system',
        timestamp: new Date(wsMessage.timestamp)
      };
      
      setMessages(prevMessages => [...prevMessages, message]);
    };
    
    webSocketService.addMessageHandler(handleMessage);
    
    return () => {
      webSocketService.removeMessageHandler(handleMessage);
    };
  }, [user]);
  
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
          {activeUsers.length > 0 && (
            <span>{activeUsers.length} người dùng đang hoạt động</span>
          )}
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