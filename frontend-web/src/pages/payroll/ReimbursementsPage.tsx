import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import StatusBadge from '../../components/StatusBadge';
import { Plus } from 'lucide-react';
import type { Reimbursement } from '../../types';

const data: Reimbursement[] = [];

export default function ReimbursementsPage() {
    return (
        <div>
            <PageHeader
                title="Reimbursements"
                subtitle="Expense claims and approvals"
                actions={
                    <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors">
                        <Plus className="w-4 h-4" />
                        New Claim
                    </button>
                }
            />
            <DataTable<Reimbursement & Record<string, unknown>>
                columns={[
                    { key: 'category', label: 'Category' },
                    { key: 'amount', label: 'Amount' },
                    { key: 'description', label: 'Description' },
                    { key: 'created_at', label: 'Date' },
                    { key: 'status', label: 'Status', render: (row) => <StatusBadge status={row.status} /> },
                ]}
                data={data as unknown as (Reimbursement & Record<string, unknown>)[]}
                emptyMessage="No reimbursement claims found."
            />
        </div>
    );
}
