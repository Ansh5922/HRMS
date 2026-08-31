import { useState } from 'react';
import { Link } from 'react-router-dom';
import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import StatusBadge from '../../components/StatusBadge';
import { Plus, Search } from 'lucide-react';
import type { Employee } from '../../types';
import { mockEmployees } from '../../lib/mockData';

export default function EmployeeListPage() {
    const [search, setSearch] = useState('');
    const [deptFilter, setDeptFilter] = useState('');
    const [statusFilter, setStatusFilter] = useState('');

    const filtered = mockEmployees.filter((emp) => {
        const matchesSearch = !search || `${emp.first_name} ${emp.last_name} ${emp.emp_code}`.toLowerCase().includes(search.toLowerCase());
        const matchesDept = !deptFilter || emp.department === deptFilter;
        const matchesStatus = !statusFilter || emp.status === statusFilter;
        return matchesSearch && matchesDept && matchesStatus;
    });

    const departments = [...new Set(mockEmployees.map(e => e.department).filter(Boolean))];

    return (
        <div>
            <PageHeader
                title="Employees"
                subtitle={`${mockEmployees.length} total employees`}
                actions={
                    <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors">
                        <Plus className="w-4 h-4" />
                        Add Employee
                    </button>
                }
            />

            <div className="flex flex-col sm:flex-row gap-3 mb-4">
                <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-surface-400" />
                    <input
                        type="text"
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        placeholder="Search by name or code..."
                        className="w-full pl-9 pr-3 py-2 rounded-lg border border-surface-300 dark:border-surface-600 bg-white dark:bg-surface-800 text-sm text-surface-900 dark:text-surface-100 focus:outline-none focus:ring-2 focus:ring-primary-500"
                    />
                </div>
                <select value={deptFilter} onChange={(e) => setDeptFilter(e.target.value)} className="px-3 py-2 rounded-lg border border-surface-300 dark:border-surface-600 bg-white dark:bg-surface-800 text-sm text-surface-700 dark:text-surface-300">
                    <option value="">All Departments</option>
                    {departments.map(d => <option key={d} value={d}>{d}</option>)}
                </select>
                <select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)} className="px-3 py-2 rounded-lg border border-surface-300 dark:border-surface-600 bg-white dark:bg-surface-800 text-sm text-surface-700 dark:text-surface-300">
                    <option value="">All Status</option>
                    <option value="active">Active</option>
                    <option value="probation">Probation</option>
                    <option value="notice_period">Notice Period</option>
                </select>
            </div>

            <DataTable<Employee & Record<string, unknown>>
                columns={[
                    { key: 'emp_code', label: 'Code' },
                    {
                        key: 'first_name', label: 'Name', render: (row) => (
                            <Link to={`/employees/${row.id}`} className="text-primary-600 dark:text-primary-400 hover:underline font-medium">
                                {row.first_name} {row.last_name}
                            </Link>
                        )
                    },
                    { key: 'department', label: 'Department' },
                    { key: 'designation', label: 'Designation' },
                    {
                        key: 'employment_type', label: 'Type', render: (row) => (
                            <span className="capitalize">{row.employment_type.replace('_', ' ')}</span>
                        )
                    },
                    { key: 'work_location', label: 'Location' },
                    { key: 'status', label: 'Status', render: (row) => <StatusBadge status={row.status} /> },
                ]}
                data={filtered as unknown as (Employee & Record<string, unknown>)[]}
                emptyMessage="No employees match your filters."
            />
        </div>
    );
}
