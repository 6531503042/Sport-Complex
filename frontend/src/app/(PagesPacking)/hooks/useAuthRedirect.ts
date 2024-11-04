"use client";

import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

const useAuthRedirect = () => {
  const router = useRouter();

  useEffect(() => {
    const accessToken = localStorage.getItem('access_token');
    
    if (!accessToken) {
      router.replace('/login');
    }
  }, [router]);
};

export default useAuthRedirect;
