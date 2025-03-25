// ApiService.ts
interface User {
  id: string;
  username: string;
  created_at: string;
}

class ApiService {
  private apiUrl: string;
  
  constructor() {
    // Determine API URL based on environment
    const isLocalhost = window.location.hostname === 'localhost';
    this.apiUrl = isLocalhost 
      ? 'http://localhost:8080/api' 
      : `${window.location.origin}/api`;
  }
  
  /**
   * Register a new user with the given username
   */
  async registerUser(username: string): Promise<User> {
    try {
      const response = await fetch(`${this.apiUrl}/users`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username }),
      });
      
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      
      const data = await response.json();
      return data.user;
    } catch (error) {
      console.error('Error registering user:', error);
      throw error;
    }
  }
  
  /**
   * Get all registered users
   */
  async getUsers(): Promise<User[]> {
    try {
      const response = await fetch(`${this.apiUrl}/users`);
      
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      
      const data = await response.json();
      return data.users;
    } catch (error) {
      console.error('Error fetching users:', error);
      throw error;
    }
  }
  
  /**
   * Get active users in the chat
   */
  async getActiveUsers(): Promise<User[]> {
    try {
      const response = await fetch(`${this.apiUrl}/users/active`);
      
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      
      const data = await response.json();
      return data.users;
    } catch (error) {
      console.error('Error fetching active users:', error);
      throw error;
    }
  }
}

// Create a singleton instance
const apiService = new ApiService();

export default apiService;
export type { User }; 