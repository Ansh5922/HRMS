import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import StatusBadge from '../../components/StatusBadge';
import { Plus } from 'lucide-react';
import type { LeaveApplication } from '../../types';
import { Link } from 'react-router-dom';

const data: LeaveApplication[] = [];

export default function LeaveListPage() {
    return (
        <div>
            <PageHeader
                title="Leave Management"
                subtitle="View and manage leave applications"
                actions={
                    <Link
                        to="/leaves/apply"
                        className="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors"
                    >
                        <Plus className="w-4 h-4" />
                        Apply Leave
                    </Link>
                }
            />

            {/* Filters */}
            <div className="flex flex-col sm:flex-row gap-3 mb-4">
                <select className="px-3 py-2 rounded-lg border border-surface-300 dark:border-surface-600 bg-white dark:bg-surface-800 text-sm text-surface-700 dark:text-surface-300">
                    <option value="">All Status</option>
                    <option value="pending">Pending</option>
                    <option value="approved">Approved</option>
                    <option value="rejected">Rejected</option>
                </select>
                <select className="px-3 py-2 rounded-lg border border-surface-300 dark:border-surface-600 bg-white dark:bg-surface-800 text-sm text-surface-700 dark:text-surface-300">
                    <option value="">All Types</option>
                </select>
            </div>

            <DataTable<LeaveApplication & Record<string, unknown>>
                columns={[
                    { key: 'employee_name', label: 'Employee' },
                    { key: 'leave_type', label: 'Type' },
                    { key: 'from_date', label: 'From' },
                    { key: 'to_date', label: 'To' },
                    { key: 'days', label: 'Days' },
                    { key: 'status', label: 'Status', render: (row) => <StatusBadge status={row.status} /> },
                ]}
                data={data as unknown as (LeaveApplication & Record<string, unknown>)[]}
                emptyMessage="No leave applications found."
            />
        </div>
    );
}
