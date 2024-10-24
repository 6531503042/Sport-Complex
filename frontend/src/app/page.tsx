"use client";
import React from 'react';
import { useAuth } from '../app/context/AuthContext';

const HomePage = () => {
  const { user, logout } = useAuth();

  return (
    <div>
      {user ? (
        <div>
          <h1>Welcome, {user.name}</h1>
          <button onClick={logout}>Logout</button>
        </div>
      ) : (
        <p>Loading...</p>
      )}
    </div>
  );
};

export default HomePage;
