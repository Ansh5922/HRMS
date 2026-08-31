import PageHeader from '../../components/PageHeader';
import { Users, Building2, CalendarDays, Wallet, UserSearch, CalendarClock } from 'lucide-react';
import { mockDashboardStats } from '../../lib/mockData';

const stats = [
    { label: 'Total Employees', value: mockDashboardStats.totalEmployees, icon: Users, color: 'text-blue-500 bg-blue-50 dark:bg-blue-900/20' },
    { label: 'Departments', value: mockDashboardStats.departments, icon: Building2, color: 'text-purple-500 bg-purple-50 dark:bg-purple-900/20' },
    { label: 'Present Today', value: mockDashboardStats.presentToday, icon: CalendarClock, color: 'text-emerald-500 bg-emerald-50 dark:bg-emerald-900/20' },
    { label: 'Pending Leaves', value: mockDashboardStats.pendingLeaves, icon: CalendarDays, color: 'text-amber-500 bg-amber-50 dark:bg-amber-900/20' },
    { label: 'Payroll This Month', value: `₹${(mockDashboardStats.payrollThisMonth / 100000).toFixed(1)}L`, icon: Wallet, color: 'text-indigo-500 bg-indigo-50 dark:bg-indigo-900/20' },
    { label: 'Open Positions', value: mockDashboardStats.openPositions, icon: UserSearch, color: 'text-rose-500 bg-rose-50 dark:bg-rose-900/20' },
];

export default function DashboardPage() {
    return (
        <div>
            <PageHeader title="Dashboard" subtitle="Overview of your organization" />

            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                {stats.map((stat) => (
                    <div
                        key={stat.label}
                        className="flex items-center gap-4 p-4 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800"
                    >
                        <div className={`p-3 rounded-lg ${stat.color}`}>
                            <stat.icon className="w-5 h-5" />
                        </div>
                        <div>
                            <p className="text-2xl font-semibold text-surface-900 dark:text-surface-50">{stat.value}</p>
                            <p className="text-sm text-surface-500 dark:text-surface-400">{stat.label}</p>
                        </div>
                    </div>
                ))}
            </div>

            {/* Placeholder sections */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-6">
                <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 p-5">
                    <h3 className="text-sm font-semibold text-surface-900 dark:text-surface-50 mb-4">Recent Activity</h3>
                    <p className="text-sm text-surface-400">No recent activity</p>
                </div>
                <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 p-5">
                    <h3 className="text-sm font-semibold text-surface-900 dark:text-surface-50 mb-4">Upcoming Events</h3>
                    <p className="text-sm text-surface-400">No upcoming events</p>
                </div>
            </div>
        </div>
    );
}
