import React from 'react';
import { Input, Button } from '@nextui-org/react';
import Image from 'next/image';
import styles from './SignUp.module.css';

const SignUpPage = () => {
  return (
    <div className={styles.container}>
      <div className={styles.left}>
        <Image className={styles.logo} src="/assets/logo-mfu-v2.png" alt="Logo" width={70} height={70} />
        <h1 className={styles.header}>WELCOME NEW USER</h1>
        <p className={styles.underheader}>Welcome to MFU Sport complex.</p>
        <form className={styles.form}>
          <div className={styles.input}>
            <Input
              fullWidth
              isClearable
              label="Name"
              placeholder="Enter your name"
              type="text"
            />
          </div>
          <div className={styles.input}>
            <Input
              fullWidth
              isClearable
              label="Email"
              placeholder="Enter your email"
              type="email"
            />
          </div>
          <div className={styles.input}>
            <Input
              fullWidth
              isClearable
              label="Password"
              placeholder="********"
              type="password"
            />
          </div>
          <div className={styles.input}>
            <Input
              fullWidth
              isClearable
              label="Confirm Password"
              placeholder="********"
              type="password"
            />
          </div>
          <Button type="submit" className={styles.button} color="primary">
            Sign up
          </Button>
        </form>
      </div>
      <div className={styles.right}>
        <Image src="/assets/loginpicture.png" alt="Sports Image" layout="fill" objectFit="cover" className={styles.rightImage} />
      </div>
    </div>
  );
};

export default SignUpPage;
