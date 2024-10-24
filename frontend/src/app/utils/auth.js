export const setToken = (token) => {
    localStorage.setItem('token', token);
  };
  
  export const getToken = () => {
    return localStorage.getItem('token');
  };
  
  export const removeToken = () => {
    localStorage.removeItem('token');
  };
  
  export const setRefreshToken = (refreshToken) => {
    localStorage.setItem('refreshToken', refreshToken);
  };
  
  export const getRefreshToken = () => {
    return localStorage.getItem('refreshToken');
  };
  
  export const removeRefreshToken = () => {
    localStorage.removeItem('refreshToken');
  };
  
  export const isTokenExpired = (token) => {
    if (!token) return true;
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.exp * 1000 < Date.now();
  };
  
  export const checkTokenExpiry = () => {
    const token = getToken();
    return isTokenExpired(token);
  };
  
  export const refreshToken = async () => {
    const refreshToken = getRefreshToken();
    if (!refreshToken || isTokenExpired(refreshToken)) {
      removeToken();
      removeRefreshToken();
      return false;
    }
  
    const response = await fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ refreshToken }),
    });
  
    if (response.ok) {
      const { token, refreshToken: newRefreshToken } = await response.json();
      setToken(token);
      setRefreshToken(newRefreshToken);
      return true;
    } else {
      removeToken();
      removeRefreshToken();
      return false;
    }
  };
  