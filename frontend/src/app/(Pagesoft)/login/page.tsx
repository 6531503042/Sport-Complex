"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Input, Button, Checkbox, Link } from '@nextui-org/react';
import Image from 'next/image';
import styles from './Login.module.css';
import { useAuth } from '../../context/AuthContext'; 

const LoginPage = () => {
  const router = useRouter();
  const [showError, setShowError] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const { setUser } = useAuth(); // Destructure setUser from useAuth

  const handleLogin = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const form = event.currentTarget;
    const email = (form.elements.namedItem('email') as HTMLInputElement).value;
    const password = (form.elements.namedItem('password') as HTMLInputElement).value;

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
        router.push('/homepage'); 
      }
    } catch (error) {
      console.error('Login error:', error);
      setErrorMessage('An error occurred. Please try again.');
      setShowError(true);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.left}>
        <Image className={styles.logo} src="/assets/logo-mfu-v2.png" alt="Logo" width={75} height={75} />
        <h1 className={styles.header}>Welcome Back</h1>
        <p className={styles.underheader}>Welcome back to MFU Sport complex.</p>
        {showError && <div className={styles.error}>{errorMessage}</div>}
        <form className={styles.form} onSubmit={handleLogin}>
          <div className={styles.input}>
            <Input
              fullWidth
              isClearable
              label="Email"
              placeholder="Enter your lamduan email"
              type="email"
              name="email"
            />
          </div>
          <div className={styles.input}>
            <Input
              fullWidth
              isClearable
              label="Password"
              placeholder="********"
              type="password"
              name="password"
            />
          </div>
          <div className={styles.checkboxContainer}>
            <div className={styles.checkboxWrapper}>
              <Checkbox className={styles.checkbox} />
              <span className={styles.checkboxText}>Remember me</span>
            </div>
            <Link href="#" className={styles.link}>Forgot password</Link>
          </div>
          <Button type="submit" className={styles.button} color="primary">
            Sign in
          </Button>
          <p className={`${styles.textCenter} ${styles.signupText}`}>
            Are you an outsider? <Link href="signup" className={styles.link}>Sign up for free!</Link>
          </p>
        </form>
      </div>
      <div className={styles.right}>
        <Image src="/assets/loginpicture.png" alt="Sports Image" layout="fill" className={styles.rightImage} />
      </div>
    </div>
  );
};

export default LoginPage;
