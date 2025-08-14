import './globals.css';

import { Inter } from 'next/font/google';
import { ReactNode } from 'react';
import ClientProvider from './client-provider';


const inter = Inter({ subsets: ['latin'] });

export const metadata = {
  title: 'Grompt - AI Prompt Engineering Tool',
  description: 'A powerful tool for building prompts with AI assistance using real engineering practices.',
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en" className="dark">
      <body className={`${inter.className} bg-gray-900 text-white`}>
        <ClientProvider>
          {children}
        </ClientProvider>
      </body>
    </html>
  );
}
