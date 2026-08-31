import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import StatusBadge from '../../components/StatusBadge';
import type { ReviewCycle } from '../../types';

const data: ReviewCycle[] = [];

export default function ReviewsPage() {
    return (
        <div>
            <PageHeader title="Performance Reviews" subtitle="Review cycles and evaluations" />
            <DataTable<ReviewCycle & Record<string, unknown>>
                columns={[
                    { key: 'name', label: 'Cycle' },
                    { key: 'type', label: 'Type' },
                    { key: 'year', label: 'Year' },
                    { key: 'start_date', label: 'Start' },
                    { key: 'end_date', label: 'End' },
                    { key: 'status', label: 'Status', render: (row) => <StatusBadge status={row.status} /> },
                ]}
                data={data as unknown as (ReviewCycle & Record<string, unknown>)[]}
                emptyMessage="No review cycles found."
            />
        </div>
    );
}
