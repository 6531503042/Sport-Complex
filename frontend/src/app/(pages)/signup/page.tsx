"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import Image from 'next/image';
import Link from 'next/link';
import styles from './SignUp.module.css';
import { motion } from 'framer-motion';
import { Email, Lock, Person, ArrowForward } from '@mui/icons-material';

const shakeAnimation = {
  initial: { x: 0 },
  animate: {
    x: [0, -10, 10, -10, 10, 0],
    transition: { duration: 0.4 },
  },
};

const SignUpPage = () => {
  const router = useRouter();
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [isSuccessful, setIsSuccessful] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [shake, setShake] = useState(false);

  const validateEmail = (email: string) => {
    const emailPattern = /^[a-zA-Z0-9._%+-]+@(gmail\.com|hotmail\.com|lamduan\.mfu\.ac\.th)$/;
    return emailPattern.test(email);
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setIsLoading(true);

    if (!validateEmail(email)) {
      setError('Invalid email address');
      setShake(true);
      setIsLoading(false);
      setTimeout(() => {
        setShake(false);
        setError('');
      }, 1000);
      return;
    }

    if (password !== confirmPassword) {
      setError('Passwords do not match');
      setShake(true);
      setIsLoading(false);
      setTimeout(() => {
        setShake(false);
        setError('');
      }, 1000);
      return;
    }

    try {
      const response = await fetch('http://localhost:1325/user_v1/users/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, email, password }),
      });

      if (response.ok) {
        setIsSuccessful(true);
      } else {
        setError('Registration failed. Please try again.');
        setShake(true);
        setTimeout(() => {
          setShake(false);
          setError('');
        }, 1000);
      }
    } catch (error) {
      console.error('An error occurred during registration:', error);
      setError('An error occurred. Please try again.');
      setShake(true);
      setTimeout(() => {
        setShake(false);
        setError('');
      }, 1000);
    } finally {
      setIsLoading(false);
    }
  };

  const handleLoginRedirect = () => {
    router.push('/login');
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
            width={70}
            height={70}
            priority
          />
        </motion.div>

        <motion.h1
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.3 }}
          className={styles.header}
        >
          Create Account
        </motion.h1>

        <motion.p
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.4 }}
          className={styles.underheader}
        >
          Join MFU Sport Complex today
        </motion.p>

        {error && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className={styles.error}
          >
            {error}
          </motion.div>
        )}

        {isSuccessful ? (
          <motion.div
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            className={styles.successCard}
          >
            <h2 className={styles.successHeader}>Registration Successful!</h2>
            <p className={styles.successText}>
              Welcome to MFU Sport Complex! You can now log in to your account.
            </p>
            <motion.button
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              onClick={handleLoginRedirect}
              className={styles.button}
            >
              Go to Login
            </motion.button>
          </motion.div>
        ) : (
          <motion.form
            className={styles.form}
            onSubmit={handleSubmit}
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.5 }}
          >
            <motion.div
              className={styles.inputGroup}
              initial={{ x: 0 }}
              animate={shake ? shakeAnimation.animate : {}}
            >
              <div className={styles.inputWrapper}>
                <Person className={styles.inputIcon} />
                <input
                  type="text"
                  placeholder="Enter your name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className={styles.input}
                  required
                />
              </div>
            </motion.div>

            <motion.div
              className={styles.inputGroup}
              initial={{ x: 0 }}
              animate={shake ? shakeAnimation.animate : {}}
            >
              <div className={styles.inputWrapper}>
                <Email className={styles.inputIcon} />
                <input
                  type="email"
                  placeholder="Enter your email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className={styles.input}
                  required
                />
              </div>
            </motion.div>

            <motion.div
              className={styles.inputGroup}
              initial={{ x: 0 }}
              animate={shake ? shakeAnimation.animate : {}}
            >
              <div className={styles.inputWrapper}>
                <Lock className={styles.inputIcon} />
                <input
                  type="password"
                  placeholder="Enter your password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className={styles.input}
                  required
                />
              </div>
            </motion.div>

            <motion.div
              className={styles.inputGroup}
              initial={{ x: 0 }}
              animate={shake ? shakeAnimation.animate : {}}
            >
              <div className={styles.inputWrapper}>
                <Lock className={styles.inputIcon} />
                <input
                  type="password"
                  placeholder="Confirm your password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  className={styles.input}
                  required
                />
              </div>
            </motion.div>

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
                  <span>Sign up</span>
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
              Already have an account?{' '}
              <Link href="/login" className={styles.link}>
                Sign in here
              </Link>
            </motion.p>
          </motion.form>
        )}
      </motion.div>
    </motion.div>
  );
};

export default SignUpPage;
