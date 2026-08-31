import { Outlet, Link, useLocation } from 'react-router-dom';
import { useState } from 'react';
import { useAuth } from '../context/auth';
import { useTheme } from '../context/theme';
import {
    LayoutDashboard, Users, Building2, Briefcase, CalendarClock, CalendarDays,
    Wallet, UserSearch, Target, GraduationCap, Bell, Megaphone, Settings,
    Menu, X, Sun, Moon, LogOut, ChevronDown,
} from 'lucide-react';

const navItems = [
    { label: 'Dashboard', icon: LayoutDashboard, path: '/dashboard' },
    { label: 'Employees', icon: Users, path: '/employees' },
    { label: 'Departments', icon: Building2, path: '/departments' },
    { label: 'Designations', icon: Briefcase, path: '/designations' },
    { label: 'Attendance', icon: CalendarClock, path: '/attendance' },
    { label: 'Leaves', icon: CalendarDays, path: '/leaves' },
    { label: 'Payroll', icon: Wallet, path: '/payroll' },
    { label: 'Recruitment', icon: UserSearch, path: '/recruitment' },
    { label: 'Performance', icon: Target, path: '/performance' },
    { label: 'Training', icon: GraduationCap, path: '/training' },
    { label: 'Notifications', icon: Bell, path: '/notifications' },
    { label: 'Announcements', icon: Megaphone, path: '/announcements' },
    { label: 'Settings', icon: Settings, path: '/settings/roles' },
];

export default function DashboardLayout() {
    const [sidebarOpen, setSidebarOpen] = useState(false);
    const { user, logout } = useAuth();
    const { theme, toggleTheme } = useTheme();
    const location = useLocation();
    const [userMenuOpen, setUserMenuOpen] = useState(false);

    return (
        <div className="min-h-screen flex">
            {/* Mobile overlay */}
            {sidebarOpen && (
                <div
                    className="fixed inset-0 bg-black/40 z-40 lg:hidden"
                    onClick={() => setSidebarOpen(false)}
                />
            )}

            {/* Sidebar */}
            <aside
                className={`
          fixed inset-y-0 left-0 z-50 w-64 bg-white dark:bg-surface-900 border-r border-surface-200 dark:border-surface-800
          transform transition-transform duration-200 ease-in-out
          lg:translate-x-0 lg:static lg:z-auto
          ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
        `}
            >
                {/* Logo */}
                <div className="h-16 flex items-center justify-between px-5 border-b border-surface-200 dark:border-surface-800">
                    <Link to="/dashboard" className="text-lg font-bold text-primary-600 dark:text-primary-400">
                        HRMS
                    </Link>
                    <button onClick={() => setSidebarOpen(false)} className="lg:hidden p-1 text-surface-400 hover:text-surface-600">
                        <X className="w-5 h-5" />
                    </button>
                </div>

                {/* Nav */}
                <nav className="p-3 space-y-0.5 overflow-y-auto h-[calc(100vh-4rem)]">
                    {navItems.map((item) => {
                        const isActive = location.pathname.startsWith(item.path);
                        return (
                            <Link
                                key={item.path}
                                to={item.path}
                                onClick={() => setSidebarOpen(false)}
                                className={`
                  flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors
                  ${isActive
                                        ? 'bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-400'
                                        : 'text-surface-600 hover:bg-surface-100 dark:text-surface-400 dark:hover:bg-surface-800'
                                    }
                `}
                            >
                                <item.icon className="w-4 h-4 shrink-0" />
                                {item.label}
                            </Link>
                        );
                    })}
                </nav>
            </aside>

            {/* Main area */}
            <div className="flex-1 flex flex-col min-w-0">
                {/* Top bar */}
                <header className="h-16 flex items-center justify-between px-4 sm:px-6 bg-white dark:bg-surface-900 border-b border-surface-200 dark:border-surface-800 sticky top-0 z-30">
                    <button
                        onClick={() => setSidebarOpen(true)}
                        className="lg:hidden p-2 -ml-2 rounded-lg text-surface-500 hover:bg-surface-100 dark:hover:bg-surface-800"
                    >
                        <Menu className="w-5 h-5" />
                    </button>

                    <div className="hidden lg:block" />

                    <div className="flex items-center gap-2">
                        {/* Theme toggle */}
                        <button
                            onClick={toggleTheme}
                            className="p-2 rounded-lg text-surface-500 hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors"
                            title={theme === 'dark' ? 'Switch to light' : 'Switch to dark'}
                        >
                            {theme === 'dark' ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
                        </button>

                        {/* Notifications */}
                        <Link
                            to="/notifications"
                            className="p-2 rounded-lg text-surface-500 hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors"
                        >
                            <Bell className="w-4 h-4" />
                        </Link>

                        {/* User menu */}
                        <div className="relative">
                            <button
                                onClick={() => setUserMenuOpen(!userMenuOpen)}
                                className="flex items-center gap-2 px-3 py-1.5 rounded-lg hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors"
                            >
                                <div className="w-7 h-7 rounded-full bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center text-xs font-semibold text-primary-700 dark:text-primary-400">
                                    {user?.email?.charAt(0).toUpperCase() || 'U'}
                                </div>
                                <span className="hidden sm:block text-sm text-surface-700 dark:text-surface-300 max-w-[120px] truncate">
                                    {user?.email || 'User'}
                                </span>
                                <ChevronDown className="w-3 h-3 text-surface-400" />
                            </button>

                            {userMenuOpen && (
                                <>
                                    <div className="fixed inset-0 z-40" onClick={() => setUserMenuOpen(false)} />
                                    <div className="absolute right-0 mt-1 w-48 bg-white dark:bg-surface-800 rounded-lg shadow-lg border border-surface-200 dark:border-surface-700 py-1 z-50">
                                        <button
                                            onClick={() => { logout(); setUserMenuOpen(false); }}
                                            className="w-full flex items-center gap-2 px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-surface-50 dark:hover:bg-surface-700"
                                        >
                                            <LogOut className="w-4 h-4" />
                                            Logout
                                        </button>
                                    </div>
                                </>
                            )}
                        </div>
                    </div>
                </header>

                {/* Page content */}
                <main className="flex-1 p-4 sm:p-6">
                    <Outlet />
                </main>
            </div>
        </div>
    );
}
