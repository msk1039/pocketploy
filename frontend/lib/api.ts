import {
  SignupRequest,
  LoginRequest,
  RefreshRequest,
  LogoutRequest,
  UpdateUserRequest,
  AuthResponse,
  RefreshResponse,
  UserResponse,
  ErrorResponse,
} from "@/types/auth";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

// Custom error class for API errors
export class ApiError extends Error {
  constructor(
    message: string,
    public statusCode: number,
    public details?: Record<string, string>
  ) {
    super(message);
    this.name = "ApiError";
  }
}

// Generic fetch wrapper with error handling
async function fetchAPI<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;

  const config: RequestInit = {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  };

  try {
    const response = await fetch(url, config);
    const data = await response.json();

    if (!response.ok) {
      const errorData = data as ErrorResponse;
      throw new ApiError(
        errorData.error || "An error occurred",
        response.status,
        errorData.details
      );
    }

    return data as T;
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }
    throw new ApiError("Network error or server unavailable", 0);
  }
}

// Token management functions
let accessToken: string | null = null;

export function setAccessToken(token: string | null) {
  accessToken = token;
}

export function getAccessToken(): string | null {
  return accessToken;
}

export function clearTokens() {
  accessToken = null;
  if (typeof window !== "undefined") {
    localStorage.removeItem("refresh_token");
  }
}

function getRefreshToken(): string | null {
  if (typeof window !== "undefined") {
    return localStorage.getItem("refresh_token");
  }
  return null;
}

function setRefreshToken(token: string) {
  if (typeof window !== "undefined") {
    localStorage.setItem("refresh_token", token);
  }
}

// Auth API functions

export async function signup(data: SignupRequest): Promise<AuthResponse> {
  const response = await fetchAPI<AuthResponse>("/auth/signup", {
    method: "POST",
    body: JSON.stringify(data),
  });

  // Store tokens
  setAccessToken(response.data.access_token);
  setRefreshToken(response.data.refresh_token);

  return response;
}

export async function login(data: LoginRequest): Promise<AuthResponse> {
  const response = await fetchAPI<AuthResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify(data),
  });

  // Store tokens
  setAccessToken(response.data.access_token);
  setRefreshToken(response.data.refresh_token);

  return response;
}

export async function logout(): Promise<void> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) {
    clearTokens();
    return;
  }

  try {
    await fetchAPI<{ success: boolean; message: string }>("/auth/logout", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${getAccessToken()}`,
      },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
  } finally {
    clearTokens();
  }
}

export async function refreshAccessToken(): Promise<string> {
  const refreshToken = getRefreshToken();
  console.log("refreshAccessToken called, refresh token:", refreshToken ? "present" : "missing");
  
  if (!refreshToken) {
    throw new ApiError("No refresh token available", 401);
  }

  const response = await fetchAPI<RefreshResponse>("/auth/refresh", {
    method: "POST",
    body: JSON.stringify({ refresh_token: refreshToken }),
  });

  console.log("New access token received");
  setAccessToken(response.data.access_token);
  return response.data.access_token;
}

export async function getCurrentUser(): Promise<UserResponse> {
  const token = getAccessToken();
  console.log("getCurrentUser called, access token:", token ? "present" : "missing");
  
  return fetchAPI<UserResponse>("/auth/me", {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function updateUserProfile(
  data: UpdateUserRequest
): Promise<UserResponse> {
  return fetchAPI<UserResponse>("/users/me", {
    method: "PATCH",
    headers: {
      Authorization: `Bearer ${getAccessToken()}`,
    },
    body: JSON.stringify(data),
  });
}

// Health check functions
export async function checkHealth(): Promise<{ status: string; timestamp: string }> {
  return fetchAPI<{ status: string; timestamp: string }>("/health", {
    method: "GET",
  });
}

export async function checkDatabaseHealth(): Promise<{
  status: string;
  message: string;
  timestamp: string;
}> {
  return fetchAPI<{ status: string; message: string; timestamp: string }>(
    "/health/db",
    {
      method: "GET",
    }
  );
}
