import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import PageHeader from '../../components/PageHeader';
import { ArrowLeft, Send } from 'lucide-react';
import { Link } from 'react-router-dom';

export default function LeaveApplyPage() {
    const navigate = useNavigate();
    const [form, setForm] = useState({
        leave_type_id: '',
        from_date: '',
        to_date: '',
        session: 'full',
        reason: '',
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
        setForm({ ...form, [e.target.name]: e.target.value });
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        // TODO: API call POST /api/v1/leaves/apply/:emp_id
        navigate('/leaves');
    };

    return (
        <div>
            <PageHeader
                title="Apply for Leave"
                actions={
                    <Link to="/leaves" className="flex items-center gap-2 text-sm text-surface-600 dark:text-surface-400 hover:text-surface-900 dark:hover:text-surface-100">
                        <ArrowLeft className="w-4 h-4" /> Back
                    </Link>
                }
            />

            <div className="max-w-lg">
                <form onSubmit={handleSubmit} className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 p-5 space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-surface-700 dark:text-surface-300 mb-1.5">Leave Type</label>
                        <select
                            name="leave_type_id"
                            value={form.leave_type_id}
                            onChange={handleChange}
                            required
                            className="w-full px-3 py-2 rounded-lg border border-surface-300 dark:border-surface-600 bg-white dark:bg-surface-800 text-sm text-surface-900 dark:text-surface-100"
                        >
                            <option value="">Select leave type</option>
                        </select>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-surface-700 dark:text-surface-300 mb-1.5">From</label>
                            <input type="date" name="from_date" value={form.from_date} onChange={handleChange} required
                                className="w-full px-3 py-2 rounded-lg border border-surface-300 dark:border-surface-600 bg-white dark:bg-surface-800 text-sm text-surface-900 dark:text-surface-100" />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-surface-700 dark:text-surface-300 mb-1.5">To</label>
                            <input type="date" name="to_date" value={form.to_date} onChange={handleChange} required
                                className="w-full px-3 py-2 rounded-lg border border-surface-300 dark:border-surface-600 bg-white dark:bg-surface-800 text-sm text-surface-900 dark:text-surface-100" />
                        </div>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-surface-700 dark:text-surface-300 mb-1.5">Session</label>
                        <select name="session" value={form.session} onChange={handleChange}
                            className="w-full px-3 py-2 rounded-lg border border-surface-300 dark:border-surface-600 bg-white dark:bg-surface-800 text-sm text-surface-900 dark:text-surface-100">
                            <option value="full">Full Day</option>
                            <option value="first_half">First Half</option>
                            <option value="second_half">Second Half</option>
                        </select>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-surface-700 dark:text-surface-300 mb-1.5">Reason</label>
                        <textarea name="reason" value={form.reason} onChange={handleChange} rows={3}
                            className="w-full px-3 py-2 rounded-lg border border-surface-300 dark:border-surface-600 bg-white dark:bg-surface-800 text-sm text-surface-900 dark:text-surface-100 resize-none"
                            placeholder="Brief reason for leave" />
                    </div>

                    <button type="submit"
                        className="w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors">
                        <Send className="w-4 h-4" /> Submit Application
                    </button>
                </form>
            </div>
        </div>
    );
}
