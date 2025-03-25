// WebSocketService.ts
interface Message {
  id: string;
  text: string;
  user_id: string;
  username: string;
  timestamp: string;
  type: 'user' | 'system';
}

class WebSocketService {
  private socket: WebSocket | null = null;
  private userId: string | null = null;
  private messageHandlers: ((message: Message) => void)[] = [];

  connect(userId: string): Promise<void> {
    return new Promise((resolve, reject) => {
      this.userId = userId;
      
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = window.location.hostname === 'localhost' ? 
        `${window.location.hostname}:8080` : window.location.host;
        
      this.socket = new WebSocket(`${protocol}//${host}/api/ws/${userId}`);
      
      this.socket.onopen = () => {
        console.log('WebSocket connection established');
        resolve();
      };
      
      this.socket.onmessage = (event) => {
        try {
          const messages = event.data.split('\n');
          messages.forEach((msgStr: string) => {
            if (msgStr.trim()) {
              const message = JSON.parse(msgStr) as Message;
              this.notifyMessageHandlers(message);
            }
          });
        } catch (error) {
          console.error('Error parsing message:', error);
        }
      };
      
      this.socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        reject(error);
      };
      
      this.socket.onclose = () => {
        console.log('WebSocket connection closed');
        // Attempt to reconnect after a delay
        setTimeout(() => {
          if (this.userId) {
            this.connect(this.userId).catch(console.error);
          }
        }, 3000);
      };
    });
  }
  
  sendMessage(text: string): void {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      console.error('WebSocket is not connected');
      return;
    }
    
    const message = {
      text: text
    };
    
    this.socket.send(JSON.stringify(message));
  }
  
  addMessageHandler(handler: (message: Message) => void): void {
    this.messageHandlers.push(handler);
  }
  
  removeMessageHandler(handler: (message: Message) => void): void {
    const index = this.messageHandlers.indexOf(handler);
    if (index !== -1) {
      this.messageHandlers.splice(index, 1);
    }
  }
  
  private notifyMessageHandlers(message: Message): void {
    this.messageHandlers.forEach(handler => {
      handler(message);
    });
  }
  
  disconnect(): void {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }
}

// Create a singleton instance
const webSocketService = new WebSocketService();

export default webSocketService;
export type { Message }; 