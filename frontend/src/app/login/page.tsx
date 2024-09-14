import React from 'react';
import { Input, Button, Checkbox, Link } from '@nextui-org/react';
import Image from 'next/image';
import styles from './Login.module.css';

const LoginPage = () => {
  return (
    <div className={styles.container}>
      <div className={styles.left}>
        <Image className={styles.logo} src="/assets/logo-mfu-v2.png" alt="Logo" width={75} height={75} />
        <h1 className={styles.header}>Welcome Back</h1>
        <p className={styles.uderheader}>Welcome back to MFU Sport complex.</p>
        <form className={styles.form}>
          <Input
            fullWidth
            isClearable
            label="Email"
            placeholder="Enter your email"
            type="email"
            className={styles.input} 
          />
          <Input
            fullWidth
            isClearable
            label="Password"
            placeholder="********"
            type="password"
            className={styles.input} 
          />
          <div className={styles.checkboxContainer}>
            <Checkbox className={styles.checkbox} />
            <span>Remember me</span>
            <Link href="#">Forgot password</Link>
          </div>
          <Button type="submit" className={styles.button} color="primary">
            Sign in
          </Button>
          <p className={styles.textCenter}>
            Are you an outsider? <Link href="#">Sign up for free!</Link>
          </p>
        </form>
      </div>
      <div className={styles.right}>
        <Image src="/assets/loginpicture.png" alt="Sports Image" layout="fill" objectFit="cover" />
      </div>
    </div>
  );
};

export default LoginPage;
