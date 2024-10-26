import { AuthProvider } from './context/AuthContext';
import '../app/globals.css';

import { ReactNode } from 'react';

interface LayoutProps {
  children: ReactNode;
}

const Layout = ({ children }: LayoutProps) => {
  return (
    <html lang="en">
      <head>
        <title>Your Website Title</title>
      </head>
      <body>
       
          <main>{children}</main>
       
      </body>
    </html>
  );
};

export default Layout;
