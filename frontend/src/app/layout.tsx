import { AuthProvider } from './context/AuthContext';
import { ThemeProvider } from 'next-themes';
import '../app/globals.css';
import { ReactNode } from 'react';

interface LayoutProps {
  children: ReactNode;
}

const Layout = ({ children }: LayoutProps) => {
  return (
    <html lang="en">
      <head>
        <title>Sport Complex</title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      </head>
      <body>
        <ThemeProvider attribute="class" defaultTheme="light">
          <AuthProvider>
            <div className="site-wrapper">
              <main className="main-content">
                {children}
              </main>
            </div>
          </AuthProvider>
        </ThemeProvider>
      </body>
    </html>
  );
};

export default Layout;
