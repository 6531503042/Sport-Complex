"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Input, Button, Card, Spacer } from '@nextui-org/react';
import Image from 'next/image';
import styles from './SignUp.module.css';

const SignUpPage = () => {
  const router = useRouter();
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [isSuccessful, setIsSuccessful] = useState(false);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (password !== confirmPassword) {
      setError('Passwords do not match');
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
        console.error('Registration failed');
      }
    } catch (error) {
      console.error('An error occurred during registration:', error);
    }
  };

  const handleLoginRedirect = () => {
    router.push('/login');
  };

  return (
    <div className={styles.container}>
      <div className={styles.left}>
        <Image className={styles.logo} src="/assets/logo-mfu-v2.png" alt="Logo" width={70} height={70} />
        <h1 className={styles.header}>WELCOME NEW USER</h1>
        <p className={styles.underheader}>Welcome to MFU Sport complex.</p>
        {isSuccessful ? (
          <Card className={styles.successCard}>
            <h2 className={styles.successHeader}>Signup Successful</h2>
            <p className={styles.successText}>Welcome to MFU Sport Complex! Your signup was successful.</p>
            <Spacer y={1} />
            <Button className={styles.button} color="primary" onClick={handleLoginRedirect}>
              OK
            </Button>
          </Card>
        ) : (
          <form className={styles.form} onSubmit={handleSubmit}>
            <div className={styles.input}>
              <Input
                fullWidth
                isClearable
                label="Name"
                placeholder="Enter your name"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
            <div className={styles.input}>
              <Input
                fullWidth
                isClearable
                label="Email"
                placeholder="Enter your email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
            <div className={styles.input}>
              <Input
                fullWidth
                isClearable
                label="Password"
                placeholder="********"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            <div className={styles.input}>
              <Input
                fullWidth
                isClearable
                label="Confirm Password"
                placeholder="********"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
              />
            </div>
            {error && <p className={styles.error}>{error}</p>}
            <Button type="submit" className={styles.button} color="primary">
              Sign up
            </Button>
          </form>
        )}
      </div>
    </div>
  );
};

export default SignUpPage;
