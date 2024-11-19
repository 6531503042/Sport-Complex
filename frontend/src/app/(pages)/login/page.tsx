"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import Image from 'next/image';
import Link from 'next/link';
import styles from './Login.module.css';
import { useAuth } from '../../context/AuthContext';
import { motion } from 'framer-motion';
import { Email, Lock, ArrowForward } from '@mui/icons-material';

const shakeAnimation = {
  initial: { x: 0 },
  animate: {
    x: [0, -10, 10, -10, 10, 0],
    transition: { duration: 0.4 },
  },
};

const LoginPage = () => {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showError, setShowError] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const { setUser } = useAuth();
  const [isLoading, setIsLoading] = useState(false);
  const [isEmailValid, setIsEmailValid] = useState(true);

  const validateEmail = (email: string) => {
    const emailPattern = /^[a-zA-Z0-9._%+-]+@(gmail\.com|hotmail\.com|lamduan\.mfu\.ac\.th)$/;
    return emailPattern.test(email);
  };

  const handleLogin = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setIsLoading(true);

    // Validate email field
    if (!email) {
      setErrorMessage('Email is required');
      setShowError(true);
      setIsEmailValid(false);
      setIsLoading(false);
      resetEmailValidation();
      return;
    }

    // Validate email format
    if (!validateEmail(email)) {
      setErrorMessage('Email must end with @gmail.com, @hotmail.com, or @lamduan.mfu.ac.th');
      setShowError(true);
      setIsEmailValid(false);
      setIsLoading(false);
      resetEmailValidation();
      return;
    }

    setIsEmailValid(true);

    try {
      const response = await fetch('http://localhost:1323/auth_v1/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        setErrorMessage('Wrong email or password');
        setShowError(true);
      } else {
        const data = await response.json();
        localStorage.setItem('user', JSON.stringify(data));
        localStorage.setItem('access_token', data.credential.access_token);
        localStorage.setItem('refresh_token', data.credential.refresh_token);
        setUser(data);
        setShowError(false);

        if (data.credential.role_code === 0) {
          router.push('/homepage'); 
        } else if (data.credential.role_code === 1) {
          router.push('/admin_dashboard'); 
        }
      }
    } catch (error) {
      console.error('Login error:', error);
      setErrorMessage('An error occurred. Please try again.');
      setShowError(true);
    } finally {
      setIsLoading(false);
    }
  };

  const resetEmailValidation = () => {
    setTimeout(() => {
      setIsEmailValid(true);
    }, 400); 
  };

  return (
    <motion.div 
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.5 }}
      className={styles.container}
    >
      <motion.div 
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ delay: 0.2 }}
        className={styles.left}
      >
        <motion.div
          whileHover={{ scale: 1.1 }}
          whileTap={{ scale: 0.9 }}
          className="flex justify-center"
        >
          <Image 
            className={styles.logo} 
            src="/assets/logo-mfu-v2.png" 
            alt="Logo" 
            width={75} 
            height={75} 
            priority
          />
        </motion.div>
        
        <motion.h1 
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.3 }}
          className={styles.header}
        >
          Welcome Back
        </motion.h1>
        
        <motion.p 
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.4 }}
          className={styles.underheader}
        >
          Welcome back to MFU Sport complex.
        </motion.p>

        {showError && (
          <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className={styles.error}
          >
            {errorMessage}
          </motion.div>
        )}

        <motion.form 
          className={styles.form} 
          onSubmit={handleLogin}
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.5 }}
        >
          <motion.div
            className={`${styles.inputGroup} ${!isEmailValid ? styles.errorInputGroup : ''}`}
            initial={isEmailValid ? {} : shakeAnimation.initial}
            animate={isEmailValid ? {} : shakeAnimation.animate}
          >
            <div className={styles.inputWrapper}>
              <Email className={styles.inputIcon} />
              <input
                type="email"
                name="email"
                placeholder="Enter your lamduan email"
                className={`${styles.input} ${!isEmailValid ? styles.errorInput : ''}`}
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
          </motion.div>

          <div className={styles.inputGroup}>
            <div className={styles.inputWrapper}>
              <Lock className={styles.inputIcon} />
              <input
                type="password"
                name="password"
                placeholder="Enter your password"
                className={styles.input}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
          </div>

          <div className={styles.forgotPassword}>
            <Link 
              href="https://www.facebook.com/mfusportcomplex" 
              className={styles.link}
            >
              Forgot password?
            </Link>
          </div>

          <motion.button
            type="submit"
            className={styles.button}
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
            disabled={isLoading}
          >
            {isLoading ? (
              <div className={styles.spinner} />
            ) : (
              <div className={styles.buttonContent}>
                <span>Sign in</span>
                <ArrowForward className={styles.buttonIcon} />
              </div>
            )}
          </motion.button>

          <motion.p 
            className={styles.textCenter}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.6 }}
          >
            Are you an outsider?{' '}
            <Link href="signup" className={styles.link}>
              Sign up for free!
            </Link>
          </motion.p>
        </motion.form>
      </motion.div>
    </motion.div>
  );
};

export default LoginPage;
