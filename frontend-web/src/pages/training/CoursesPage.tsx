import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import StatusBadge from '../../components/StatusBadge';
import { Plus } from 'lucide-react';
import type { Course } from '../../types';

const data: Course[] = [];

export default function CoursesPage() {
    return (
        <div>
            <PageHeader
                title="Training"
                subtitle="Courses and learning management"
                actions={
                    <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors">
                        <Plus className="w-4 h-4" />
                        Add Course
                    </button>
                }
            />
            <DataTable<Course & Record<string, unknown>>
                columns={[
                    { key: 'title', label: 'Title' },
                    { key: 'category', label: 'Category' },
                    { key: 'type', label: 'Type' },
                    { key: 'duration_hrs', label: 'Duration (hrs)' },
                    {
                        key: 'is_mandatory', label: 'Mandatory', render: (row) => (
                            <StatusBadge status={row.is_mandatory ? 'active' : 'draft'} />
                        )
                    },
                ]}
                data={data as unknown as (Course & Record<string, unknown>)[]}
                emptyMessage="No courses found."
            />
        </div>
    );
}
