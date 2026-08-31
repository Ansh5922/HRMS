import PageHeader from '../../components/PageHeader';
import { Bell, Check, CheckCheck } from 'lucide-react';
import type { Notification } from '../../types';

const data: Notification[] = [];

export default function NotificationsPage() {
    return (
        <div>
            <PageHeader
                title="Notifications"
                subtitle="Stay updated with recent activity"
                actions={
                    <button className="flex items-center gap-2 px-3 py-2 rounded-lg text-sm text-surface-600 dark:text-surface-400 hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors">
                        <CheckCheck className="w-4 h-4" />
                        Mark all read
                    </button>
                }
            />

            {data.length === 0 ? (
                <div className="text-center py-16">
                    <Bell className="w-10 h-10 mx-auto text-surface-300 dark:text-surface-600 mb-3" />
                    <p className="text-sm text-surface-400">No notifications yet</p>
                </div>
            ) : (
                <div className="space-y-1">
                    {data.map((n) => (
                        <div
                            key={n.id}
                            className={`flex items-start gap-3 p-4 rounded-lg border transition-colors ${n.is_read
                                    ? 'bg-white dark:bg-surface-900 border-surface-200 dark:border-surface-800'
                                    : 'bg-primary-50/50 dark:bg-primary-900/10 border-primary-100 dark:border-primary-900/30'
                                }`}
                        >
                            <div className="shrink-0 mt-0.5">
                                {n.is_read ? (
                                    <Check className="w-4 h-4 text-surface-400" />
                                ) : (
                                    <div className="w-2 h-2 rounded-full bg-primary-500 mt-1" />
                                )}
                            </div>
                            <div className="flex-1 min-w-0">
                                <p className="text-sm font-medium text-surface-900 dark:text-surface-50">{n.title}</p>
                                {n.message && <p className="text-sm text-surface-500 dark:text-surface-400 mt-0.5">{n.message}</p>}
                                <p className="text-xs text-surface-400 mt-1">{n.created_at}</p>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
