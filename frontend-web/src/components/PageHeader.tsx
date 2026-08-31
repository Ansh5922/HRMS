import type { ReactNode } from 'react';

interface PageHeaderProps {
    title: string;
    subtitle?: string;
    actions?: ReactNode;
}

export default function PageHeader({ title, subtitle, actions }: PageHeaderProps) {
    return (
        <div className="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between mb-6">
            <div>
                <h1 className="text-xl font-semibold text-surface-900 dark:text-surface-50">{title}</h1>
                {subtitle && <p className="text-sm text-surface-500 dark:text-surface-400 mt-0.5">{subtitle}</p>}
            </div>
            {actions && <div className="flex items-center gap-2 mt-3 sm:mt-0">{actions}</div>}
        </div>
    );
}
