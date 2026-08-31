import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import StatusBadge from '../../components/StatusBadge';
import { Plus } from 'lucide-react';
import type { WorkflowTemplate } from '../../types';

const data: WorkflowTemplate[] = [];

export default function WorkflowsPage() {
    return (
        <div>
            <PageHeader
                title="Workflows"
                subtitle="Custom approval workflows"
                actions={
                    <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors">
                        <Plus className="w-4 h-4" />
                        Create Workflow
                    </button>
                }
            />
            <DataTable<WorkflowTemplate & Record<string, unknown>>
                columns={[
                    { key: 'name', label: 'Name' },
                    { key: 'module', label: 'Module' },
                    { key: 'steps', label: 'Steps', render: (row) => `${(row.steps as unknown[]).length} steps` },
                    {
                        key: 'is_active', label: 'Status', render: (row) => (
                            <StatusBadge status={row.is_active ? 'active' : 'draft'} />
                        )
                    },
                ]}
                data={data as unknown as (WorkflowTemplate & Record<string, unknown>)[]}
                emptyMessage="No workflow templates found."
            />
        </div>
    );
}
