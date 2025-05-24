import { useState, useEffect } from 'react'
import './App.css'
import InputForm from './components/InputForm'
import ChatPage from './components/ChatPage'

function App() {
  const [currentPage, setCurrentPage] = useState<'input' | 'chat'>('input');
  const [username, setUsername] = useState<string | null>(null);
  const [isOnline, setIsOnline] = useState(false);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    if (username) {
      const websocket = new WebSocket('ws://localhost:8080/ws');
      
      websocket.onopen = () => {
        setIsOnline(true);
        websocket.send(JSON.stringify({
          type: 'user_connected',
          username: username
        }));
      };

      websocket.onclose = () => {
        setIsOnline(false);
      };

      setWs(websocket);

      return () => {
        websocket.close();
      };
    }
  }, [username]);

  const handleStartChat = (name: string) => {
    setUsername(name);
    setCurrentPage('chat');
  };

  return (
    <div className="app-container">
      <header className="app-header">
        <h1>Chat vui hjhj</h1>
        {currentPage === 'chat' && username && (
          <div className="user-info">
            Xin chào, <span className="username">{username}</span>!
            <span className={`status-indicator ${isOnline ? 'online' : 'offline'}`}>
              {isOnline ? 'Online' : 'Offline'}
            </span>
            <button className="back-button" onClick={() => setCurrentPage('input')}>
              Quay lại
            </button>
          </div>
        )}
      </header>
      <main>
        {currentPage === 'input' && <InputForm onStart={handleStartChat} />}
        {currentPage === 'chat' && <ChatPage username={username || 'Khách'} ws={ws} />}
      </main>
    </div>
  )
}

export default App
