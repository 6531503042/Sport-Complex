"use client";

import React, { useState, useEffect } from 'react';
import { Modal, Button, Text } from '@nextui-org/react';
import { useAuth } from '../context/AuthContext';
import { getToken, isTokenExpired } from '../utils/auth';

const SessionTimeoutPopup = () => {
  const { logout } = useAuth();
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    const interval = setInterval(() => {
      const token = getToken();
      if (!token || isTokenExpired(token)) {
        setIsOpen(true);
      }
    }, 60000);

    return () => clearInterval(interval);
  }, []);

  return (
    <Modal open={isOpen} onClose={() => setIsOpen(false)}>
      <Modal.Header>
        <Text h3>Session Expired</Text>
      </Modal.Header>
      <Modal.Body>
        <Text>Your session has expired. Please login again.</Text>
      </Modal.Body>
      <Modal.Footer>
        <Button auto flat color="primary" onClick={logout}>
          Go to Login
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

export default SessionTimeoutPopup;
