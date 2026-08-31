import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import StatusBadge from '../../components/StatusBadge';
import { Plus } from 'lucide-react';
import type { Goal } from '../../types';

const data: Goal[] = [];

export default function GoalsPage() {
    return (
        <div>
            <PageHeader
                title="Goals"
                subtitle="Track and manage employee goals"
                actions={
                    <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors">
                        <Plus className="w-4 h-4" />
                        Set Goal
                    </button>
                }
            />
            <DataTable<Goal & Record<string, unknown>>
                columns={[
                    { key: 'title', label: 'Goal' },
                    { key: 'type', label: 'Type' },
                    { key: 'target_value', label: 'Target' },
                    { key: 'current_value', label: 'Current' },
                    { key: 'weight', label: 'Weight %' },
                    { key: 'due_date', label: 'Due Date' },
                    { key: 'status', label: 'Status', render: (row) => <StatusBadge status={row.status} /> },
                ]}
                data={data as unknown as (Goal & Record<string, unknown>)[]}
                emptyMessage="No goals set."
            />
        </div>
    );
}
