import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import { Plus } from 'lucide-react';
import type { Designation } from '../../types';

const data: Designation[] = [];

export default function DesignationsPage() {
    return (
        <div>
            <PageHeader
                title="Designations"
                subtitle="Manage job titles and grades"
                actions={
                    <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors">
                        <Plus className="w-4 h-4" />
                        Add Designation
                    </button>
                }
            />
            <DataTable<Designation & Record<string, unknown>>
                columns={[
                    { key: 'title', label: 'Title' },
                    { key: 'grade', label: 'Grade' },
                    { key: 'level', label: 'Level' },
                ]}
                data={data as unknown as (Designation & Record<string, unknown>)[]}
                emptyMessage="No designations found."
            />
        </div>
    );
}
