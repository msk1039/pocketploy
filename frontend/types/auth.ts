// User type
export interface User {
  id: string;
  username: string;
  email: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  last_login_at?: string;
}

// Auth API Request types
export interface SignupRequest {
  username: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RefreshRequest {
  refresh_token: string;
}

export interface LogoutRequest {
  refresh_token: string;
}

export interface UpdateUserRequest {
  username?: string;
  email?: string;
}

// Auth API Response types
export interface AuthResponse {
  success: boolean;
  message?: string;
  data: {
    user: User;
    access_token: string;
    refresh_token: string;
    expires_at: string;
  };
}

export interface RefreshResponse {
  success: boolean;
  data: {
    access_token: string;
    expires_at: string;
  };
}

export interface UserResponse {
  success: boolean;
  data: {
    user: User;
  };
}

export interface ErrorResponse {
  success: false;
  error: string;
  details?: Record<string, string>;
}

// Auth Context types
export interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  signup: (username: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}
