"use client";

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useRouter } from 'next/navigation';
import { getToken, removeToken, removeRefreshToken, refreshToken as refreshUserToken } from '../utils/auth';

interface AuthContextType {
  user: any; 
  setUser: (user: any) => void; 
  logout: () => void;
}

// Create a default value for the AuthContext
const defaultAuthContext: AuthContextType = {
  user: null,
  setUser: () => {},
  logout: () => {},
};

const AuthContext = createContext<AuthContextType>(defaultAuthContext);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [user, setUser] = useState(null);
  const router = useRouter();

  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    const token = getToken();

    if (storedUser && token) {
      setUser(JSON.parse(storedUser));
    }

    const interval = setInterval(async () => {
      if (token) {
        const refreshed = await refreshUserToken();
        if (!refreshed) {
          logout();
        }
      }
    }, 60000); 

    return () => clearInterval(interval);
  }, []);

  const logout = () => {
    removeToken();
    removeRefreshToken();
    localStorage.removeItem('user');
    setUser(null);
    router.push('/login');
  };

  return (
    <AuthContext.Provider value={{ user, setUser, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
