import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import StatusBadge from '../../components/StatusBadge';
import { Plus } from 'lucide-react';
import type { JobPosting } from '../../types';

const data: JobPosting[] = [];

export default function JobsPage() {
    return (
        <div>
            <PageHeader
                title="Recruitment"
                subtitle="Job postings and candidate pipeline"
                actions={
                    <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors">
                        <Plus className="w-4 h-4" />
                        Post Job
                    </button>
                }
            />
            <DataTable<JobPosting & Record<string, unknown>>
                columns={[
                    { key: 'title', label: 'Title' },
                    { key: 'type', label: 'Type' },
                    { key: 'openings', label: 'Openings' },
                    { key: 'location', label: 'Location' },
                    { key: 'status', label: 'Status', render: (row) => <StatusBadge status={row.status} /> },
                    { key: 'created_at', label: 'Posted' },
                ]}
                data={data as unknown as (JobPosting & Record<string, unknown>)[]}
                emptyMessage="No job postings found."
            />
        </div>
    );
}
