import { useState } from 'react'
import './App.css'
import InputForm from './components/InputForm'
import ChatPage from './components/ChatPage'

function App() {
  const [currentPage, setCurrentPage] = useState<'input' | 'chat'>('input');
  const [username, setUsername] = useState<string | null>(null);

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
            <button className="back-button" onClick={() => setCurrentPage('input')}>
              Quay lại
            </button>
          </div>
        )}
      </header>
      <main>
        {currentPage === 'input' && <InputForm onStart={handleStartChat} />}
        {currentPage === 'chat' && <ChatPage username={username || 'Khách'} />}
      </main>
    </div>
  )
}

export default App
