import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import { Plus } from 'lucide-react';
import type { Department } from '../../types';

const data: Department[] = [];

export default function DepartmentsPage() {
    return (
        <div>
            <PageHeader
                title="Departments"
                subtitle="Manage organization departments"
                actions={
                    <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors">
                        <Plus className="w-4 h-4" />
                        Add Department
                    </button>
                }
            />
            <DataTable<Department & Record<string, unknown>>
                columns={[
                    { key: 'name', label: 'Name' },
                    { key: 'code', label: 'Code' },
                    { key: 'description', label: 'Description' },
                    {
                        key: 'is_active', label: 'Status', render: (row) => (
                            <span className={`text-xs font-medium ${row.is_active ? 'text-emerald-600' : 'text-surface-400'}`}>
                                {row.is_active ? 'Active' : 'Inactive'}
                            </span>
                        )
                    },
                ]}
                data={data as unknown as (Department & Record<string, unknown>)[]}
                emptyMessage="No departments found."
            />
        </div>
    );
}
