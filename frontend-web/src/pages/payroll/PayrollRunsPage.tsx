import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import StatusBadge from '../../components/StatusBadge';
import { Plus } from 'lucide-react';
import type { PayrollRun } from '../../types';

const data: PayrollRun[] = [];

export default function PayrollRunsPage() {
    return (
        <div>
            <PageHeader
                title="Payroll"
                subtitle="Monthly payroll processing"
                actions={
                    <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors">
                        <Plus className="w-4 h-4" />
                        New Payroll Run
                    </button>
                }
            />
            <DataTable<PayrollRun & Record<string, unknown>>
                columns={[
                    { key: 'month', label: 'Month' },
                    { key: 'year', label: 'Year' },
                    { key: 'total_gross', label: 'Gross' },
                    { key: 'total_deductions', label: 'Deductions' },
                    { key: 'total_net', label: 'Net Pay' },
                    { key: 'status', label: 'Status', render: (row) => <StatusBadge status={row.status} /> },
                ]}
                data={data as unknown as (PayrollRun & Record<string, unknown>)[]}
                emptyMessage="No payroll runs found."
            />
        </div>
    );
}
