import { useState } from 'react';
import { Link } from 'react-router-dom';
import { Mail } from 'lucide-react';
import api from '../../lib/api';

export default function ForgotPasswordPage() {
    const [email, setEmail] = useState('');
    const [sent, setSent] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);
        try {
            await api.post('/auth/forgot-password', { email });
            setSent(true);
        } catch {
            setError('Failed to send reset email.');
        } finally {
            setLoading(false);
        }
    };

    if (sent) {
        return (
            <div className="text-center">
                <Mail className="w-10 h-10 mx-auto text-primary-500 mb-4" />
                <h2 className="text-lg font-semibold text-surface-900 dark:text-surface-50 mb-2">Check your email</h2>
                <p className="text-sm text-surface-500 dark:text-surface-400 mb-6">
                    We've sent a password reset link to <strong>{email}</strong>
                </p>
                <Link to="/login" className="text-sm text-primary-600 dark:text-primary-400 hover:underline">
                    Back to sign in
                </Link>
            </div>
        );
    }

    return (
        <div>
            <h2 className="text-lg font-semibold text-surface-900 dark:text-surface-50 mb-2">Reset your password</h2>
            <p className="text-sm text-surface-500 dark:text-surface-400 mb-6">Enter your email and we'll send you a reset link.</p>

            {error && (
                <div className="mb-4 p-3 rounded-lg bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 text-sm">{error}</div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                    <label className="block text-sm font-medium text-surface-700 dark:text-surface-300 mb-1.5">Email</label>
                    <input
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                        className="w-full px-3 py-2 rounded-lg border border-surface-300 dark:border-surface-600 bg-white dark:bg-surface-800 text-surface-900 dark:text-surface-100 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                        placeholder="you@company.com"
                    />
                </div>

                <button
                    type="submit"
                    disabled={loading}
                    className="w-full px-4 py-2.5 rounded-lg bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium transition-colors disabled:opacity-50"
                >
                    {loading ? 'Sending...' : 'Send reset link'}
                </button>
            </form>

            <p className="mt-6 text-center text-sm text-surface-500 dark:text-surface-400">
                <Link to="/login" className="text-primary-600 dark:text-primary-400 hover:underline">Back to sign in</Link>
            </p>
        </div>
    );
}
