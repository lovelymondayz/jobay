import { useState, useEffect } from 'react';
import AdminDashboard from './pages/AdminDashboard';
import UploadPage from './pages/UploadPage';
import UserDashboard from './pages/UserDashboard';

const SLUG_REGEX = /^[a-z0-9-]+$/;

export default function App() {
  const [path, setPath] = useState(window.location.pathname);

  useEffect(() => {
    const onPopState = () => setPath(window.location.pathname);
    window.addEventListener('popstate', onPopState);
    return () => window.removeEventListener('popstate', onPopState);
  }, []);

  // Navigation helper (used by Link buttons)
  const navigate = (to: string) => {
    window.history.pushState({}, '', to);
    setPath(to);
  };

  // Route matching
  if (path === '/upload') {
    return <UploadPage />;
  }

  // Check if path is a user slug (not /api, /ws, /assets, /uploads)
  const slug = path.slice(1); // remove leading /
  if (slug && SLUG_REGEX.test(slug) && !slug.startsWith('api') && !slug.startsWith('ws') && !slug.startsWith('assets') && !slug.startsWith('uploads')) {
    return <UserDashboard slug={slug} />;
  }

  // Default: admin dashboard
  return <AdminDashboard />;
}
