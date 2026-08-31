import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import { Search } from 'lucide-react';
import type { Candidate } from '../../types';

const data: Candidate[] = [];

export default function CandidatesPage() {
    return (
        <div>
            <PageHeader title="Candidates" subtitle="All candidate profiles" />
            <div className="relative max-w-sm mb-4">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-surface-400" />
                <input
                    type="text"
                    placeholder="Search candidates..."
                    className="w-full pl-9 pr-3 py-2 rounded-lg border border-surface-300 dark:border-surface-600 bg-white dark:bg-surface-800 text-sm text-surface-900 dark:text-surface-100 focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
            </div>
            <DataTable<Candidate & Record<string, unknown>>
                columns={[
                    { key: 'first_name', label: 'Name', render: (row) => `${row.first_name} ${row.last_name || ''}` },
                    { key: 'email', label: 'Email' },
                    { key: 'phone', label: 'Phone' },
                    { key: 'source', label: 'Source' },
                    { key: 'created_at', label: 'Added' },
                ]}
                data={data as unknown as (Candidate & Record<string, unknown>)[]}
                emptyMessage="No candidates found."
            />
        </div>
    );
}
