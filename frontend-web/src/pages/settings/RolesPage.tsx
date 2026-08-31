import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import { Plus, Shield } from 'lucide-react';
import type { Role } from '../../types';

const data: Role[] = [];

export default function RolesPage() {
    return (
        <div>
            <PageHeader
                title="Roles & Permissions"
                subtitle="Manage access control"
                actions={
                    <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors">
                        <Plus className="w-4 h-4" />
                        Create Role
                    </button>
                }
            />

            {data.length === 0 ? (
                <div className="text-center py-16">
                    <Shield className="w-10 h-10 mx-auto text-surface-300 dark:text-surface-600 mb-3" />
                    <p className="text-sm text-surface-400">No custom roles created yet</p>
                </div>
            ) : (
                <DataTable<Role & Record<string, unknown>>
                    columns={[
                        { key: 'name', label: 'Role Name' },
                        { key: 'description', label: 'Description' },
                        {
                            key: 'is_system', label: 'Type', render: (row) => (
                                <span className={`text-xs font-medium ${row.is_system ? 'text-blue-500' : 'text-surface-500'}`}>
                                    {row.is_system ? 'System' : 'Custom'}
                                </span>
                            )
                        },
                    ]}
                    data={data as unknown as (Role & Record<string, unknown>)[]}
                />
            )}
        </div>
    );
}
