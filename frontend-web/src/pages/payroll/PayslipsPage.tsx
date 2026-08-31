import PageHeader from '../../components/PageHeader';
import DataTable from '../../components/DataTable';
import type { Payslip } from '../../types';

const data: Payslip[] = [];

export default function PayslipsPage() {
    return (
        <div>
            <PageHeader title="Payslips" subtitle="Employee payslip history" />
            <DataTable<Payslip & Record<string, unknown>>
                columns={[
                    { key: 'emp_id', label: 'Employee' },
                    { key: 'working_days', label: 'Working Days' },
                    { key: 'lop_days', label: 'LOP Days' },
                    { key: 'gross', label: 'Gross' },
                    { key: 'total_deductions', label: 'Deductions' },
                    { key: 'net_pay', label: 'Net Pay' },
                ]}
                data={data as unknown as (Payslip & Record<string, unknown>)[]}
                emptyMessage="No payslips found."
            />
        </div>
    );
}
