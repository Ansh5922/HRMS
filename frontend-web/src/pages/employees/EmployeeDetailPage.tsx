import { useParams, Link } from 'react-router-dom';
import PageHeader from '../../components/PageHeader';
import StatusBadge from '../../components/StatusBadge';
import { ArrowLeft, Mail, Phone, MapPin, Building2, Briefcase, Calendar, User } from 'lucide-react';
import { mockEmployees } from '../../lib/mockData';

export default function EmployeeDetailPage() {
    const { id } = useParams();
    const emp = mockEmployees.find(e => e.id === id);

    if (!emp) {
        return (
            <div>
                <PageHeader title="Employee Not Found" actions={
                    <Link to="/employees" className="flex items-center gap-2 text-sm text-surface-600 dark:text-surface-400 hover:text-surface-900 dark:hover:text-surface-100">
                        <ArrowLeft className="w-4 h-4" /> Back to list
                    </Link>
                } />
                <p className="text-surface-400">Employee with ID {id} was not found.</p>
            </div>
        );
    }

    const details = [
        { icon: Mail, label: 'Email', value: emp.personal_email },
        { icon: Phone, label: 'Phone', value: emp.phone },
        { icon: Building2, label: 'Department', value: emp.department },
        { icon: Briefcase, label: 'Designation', value: emp.designation },
        { icon: User, label: 'Manager', value: emp.manager_name || 'None' },
        { icon: MapPin, label: 'Location', value: `${emp.city}, ${emp.state}` },
        { icon: Calendar, label: 'Joined', value: emp.joining_date },
    ];

    return (
        <div>
            <PageHeader
                title="Employee Details"
                actions={
                    <Link to="/employees" className="flex items-center gap-2 text-sm text-surface-600 dark:text-surface-400 hover:text-surface-900 dark:hover:text-surface-100">
                        <ArrowLeft className="w-4 h-4" /> Back to list
                    </Link>
                }
            />

            {/* Header card */}
            <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 p-5 mb-6">
                <div className="flex flex-col sm:flex-row gap-4">
                    <div className="w-16 h-16 rounded-full bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center text-xl font-bold text-primary-700 dark:text-primary-400 shrink-0">
                        {emp.first_name.charAt(0)}{emp.last_name.charAt(0)}
                    </div>
                    <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-3 flex-wrap">
                            <h2 className="text-lg font-semibold text-surface-900 dark:text-surface-50">{emp.first_name} {emp.last_name}</h2>
                            <StatusBadge status={emp.status} />
                            <span className="text-xs px-2 py-0.5 rounded-full bg-surface-100 dark:bg-surface-800 text-surface-500 font-mono">{emp.emp_code}</span>
                        </div>
                        <p className="text-sm text-surface-500 dark:text-surface-400 mt-1">
                            {emp.designation} &middot; {emp.department} &middot; <span className="capitalize">{emp.employment_type.replace('_', ' ')}</span>
                        </p>
                    </div>
                </div>
            </div>

            {/* Details grid */}
            <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 p-5 mb-6">
                <h3 className="text-sm font-semibold text-surface-900 dark:text-surface-50 mb-4">Personal Information</h3>
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                    {details.map((d) => (
                        <div key={d.label} className="flex items-start gap-3">
                            <d.icon className="w-4 h-4 text-surface-400 mt-0.5 shrink-0" />
                            <div>
                                <p className="text-xs text-surface-400">{d.label}</p>
                                <p className="text-sm text-surface-700 dark:text-surface-300">{d.value || '--'}</p>
                            </div>
                        </div>
                    ))}
                    <div className="flex items-start gap-3">
                        <Calendar className="w-4 h-4 text-surface-400 mt-0.5 shrink-0" />
                        <div>
                            <p className="text-xs text-surface-400">Date of Birth</p>
                            <p className="text-sm text-surface-700 dark:text-surface-300">{emp.dob || '--'}</p>
                        </div>
                    </div>
                    <div className="flex items-start gap-3">
                        <User className="w-4 h-4 text-surface-400 mt-0.5 shrink-0" />
                        <div>
                            <p className="text-xs text-surface-400">Gender</p>
                            <p className="text-sm text-surface-700 dark:text-surface-300">{emp.gender || '--'}</p>
                        </div>
                    </div>
                    <div className="flex items-start gap-3">
                        <MapPin className="w-4 h-4 text-surface-400 mt-0.5 shrink-0" />
                        <div>
                            <p className="text-xs text-surface-400">Work Location</p>
                            <p className="text-sm text-surface-700 dark:text-surface-300">{emp.work_location || '--'}</p>
                        </div>
                    </div>
                </div>
            </div>

            {/* Tabs placeholder */}
            <div className="border-b border-surface-200 dark:border-surface-800 mb-6 overflow-x-auto">
                <div className="flex gap-0 min-w-max">
                    {['Documents', 'Bank Details', 'Attendance', 'Leaves', 'Payslips'].map((tab, i) => (
                        <button
                            key={tab}
                            className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors whitespace-nowrap ${i === 0
                                    ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                                    : 'border-transparent text-surface-500 hover:text-surface-700 dark:hover:text-surface-300'
                                }`}
                        >
                            {tab}
                        </button>
                    ))}
                </div>
            </div>

            <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 p-5">
                <p className="text-sm text-surface-400">Tab content will be loaded from API</p>
            </div>
        </div>
    );
}
