import { Navigate } from 'react-router-dom';
import { useAuth } from '../context/auth';
import { Loader } from 'lucide-react';

const DEV_BYPASS = import.meta.env.DEV; // true when running `npm run dev`

export default function ProtectedRoute({ children }: { children: React.ReactNode }) {
    const { isAuthenticated, isLoading } = useAuth();

    // Skip auth in dev mode so all pages are accessible without backend
    if (DEV_BYPASS) {
        return <>{children}</>;
    }

    if (isLoading) {
        return (
            <div className="flex items-center justify-center h-screen">
                <Loader className="w-6 h-6 animate-spin text-primary-500" />
            </div>
        );
    }

    if (!isAuthenticated) {
        return <Navigate to="/login" replace />;
    }

    return <>{children}</>;
}
