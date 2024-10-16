"use client";

import React, { useState } from 'react';
import { Input, Button, Modal, Text } from '@nextui-org/react';
import Image from 'next/image';
import styles from './SignUp.module.css';

const SignUpPage = () => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    const response = await fetch('http://localhost:1325/user_v1/users/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name, email, password }),
    });

    if (response.ok) {
      // Show the success modal
      setIsModalOpen(true);
    } else {
      // Handle registration error
      console.error('Registration failed');
    }
  };

  const closeModal = () => {
    setIsModalOpen(false);
  };

  return (
    <div className={styles.container}>
      <div className={styles.left}>
        <Image className={styles.logo} src="/assets/logo-mfu-v2.png" alt="Logo" width={70} height={70} />
        <h1 className={styles.header}>WELCOME NEW USER</h1>
        <p className={styles.underheader}>Welcome to MFU Sport complex.</p>
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
      </div>
      <div className={styles.right}>
        <Image
          src="/assets/loginpicture.png"
          alt="Sports Image"
          layout="fill"
          className={styles.rightImage}
        />
      </div>
      <Modal
        closeButton
        aria-labelledby="modal-title"
        open={isModalOpen}
        onClose={closeModal}
      >
        <Modal.Header>
          <Text id="modal-title" size={18}>
            Signup Successful
          </Text>
        </Modal.Header>
        <Modal.Body>
          <Text>Welcome to MFU Sport Complex! Your signup was successful.</Text>
        </Modal.Body>
        <Modal.Footer>
          <Button auto flat color="primary" onClick={closeModal}>
            Close
          </Button>
        </Modal.Footer>
      </Modal>
    </div>
  );
};

export default SignUpPage;
