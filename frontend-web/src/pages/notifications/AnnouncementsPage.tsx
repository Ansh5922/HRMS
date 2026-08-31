import PageHeader from '../../components/PageHeader';
import { Plus, Megaphone } from 'lucide-react';
import type { Announcement } from '../../types';

const data: Announcement[] = [];

export default function AnnouncementsPage() {
    return (
        <div>
            <PageHeader
                title="Announcements"
                subtitle="Company-wide announcements"
                actions={
                    <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors">
                        <Plus className="w-4 h-4" />
                        New Announcement
                    </button>
                }
            />

            {data.length === 0 ? (
                <div className="text-center py-16">
                    <Megaphone className="w-10 h-10 mx-auto text-surface-300 dark:text-surface-600 mb-3" />
                    <p className="text-sm text-surface-400">No announcements yet</p>
                </div>
            ) : (
                <div className="space-y-4">
                    {data.map((a) => (
                        <div key={a.id} className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 p-5">
                            <h3 className="text-sm font-semibold text-surface-900 dark:text-surface-50">{a.title}</h3>
                            <p className="text-sm text-surface-600 dark:text-surface-400 mt-2">{a.content}</p>
                            <p className="text-xs text-surface-400 mt-3">{a.published_at}</p>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
