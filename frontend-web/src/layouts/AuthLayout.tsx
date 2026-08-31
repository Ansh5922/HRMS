import { Outlet } from 'react-router-dom';
import { useTheme } from '../context/theme';
import { Sun, Moon } from 'lucide-react';

export default function AuthLayout() {
    const { theme, toggleTheme } = useTheme();

    return (
        <div className="min-h-screen flex items-center justify-center bg-surface-50 dark:bg-surface-950 px-4 py-8">
            {/* Theme toggle */}
            <button
                onClick={toggleTheme}
                className="fixed top-4 right-4 p-2 rounded-lg text-surface-500 hover:bg-surface-200 dark:hover:bg-surface-800 transition-colors"
            >
                {theme === 'dark' ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
            </button>

            <div className="w-full max-w-md">
                {/* Logo */}
                <div className="text-center mb-8">
                    <h1 className="text-2xl font-bold text-primary-600 dark:text-primary-400">HRMS</h1>
                    <p className="text-sm text-surface-500 dark:text-surface-400 mt-1">Human Resource Management System</p>
                </div>

                {/* Card */}
                <div className="bg-white dark:bg-surface-900 rounded-xl shadow-sm border border-surface-200 dark:border-surface-800 p-6 sm:p-8">
                    <Outlet />
                </div>
            </div>
        </div>
    );
}
