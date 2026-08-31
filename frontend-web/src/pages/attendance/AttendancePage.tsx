import { useState } from 'react';
import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import StatusBadge from '../../components/StatusBadge';
import { CalendarDays } from 'lucide-react';
import type { AttendanceRecord } from '../../types';

const data: AttendanceRecord[] = [];

export default function AttendancePage() {
    const [date, setDate] = useState(new Date().toISOString().split('T')[0]);

    return (
        <div>
            <PageHeader
                title="Attendance"
                subtitle="Daily attendance records"
                actions={
                    <div className="flex items-center gap-2">
                        <CalendarDays className="w-4 h-4 text-surface-400" />
                        <input
                            type="date"
                            value={date}
                            onChange={(e) => setDate(e.target.value)}
                            className="px-3 py-2 rounded-lg border border-surface-300 dark:border-surface-600 bg-white dark:bg-surface-800 text-sm text-surface-700 dark:text-surface-300"
                        />
                    </div>
                }
            />

            {/* Summary cards */}
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-6">
                {[
                    { label: 'Present', value: '--', color: 'text-emerald-600' },
                    { label: 'Absent', value: '--', color: 'text-red-500' },
                    { label: 'Late', value: '--', color: 'text-orange-500' },
                    { label: 'On Leave', value: '--', color: 'text-blue-500' },
                ].map((s) => (
                    <div key={s.label} className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 p-4 text-center">
                        <p className={`text-2xl font-semibold ${s.color}`}>{s.value}</p>
                        <p className="text-xs text-surface-500 mt-1">{s.label}</p>
                    </div>
                ))}
            </div>

            <DataTable<AttendanceRecord & Record<string, unknown>>
                columns={[
                    { key: 'emp_id', label: 'Employee' },
                    { key: 'check_in', label: 'Check In' },
                    { key: 'check_out', label: 'Check Out' },
                    { key: 'method', label: 'Method' },
                    { key: 'duration_mins', label: 'Duration' },
                    { key: 'status', label: 'Status', render: (row) => <StatusBadge status={row.status} /> },
                ]}
                data={data as unknown as (AttendanceRecord & Record<string, unknown>)[]}
                emptyMessage="No attendance records for this date."
            />
        </div>
    );
}
