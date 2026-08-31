interface Column<T> {
    key: string;
    label: string;
    render?: (row: T) => React.ReactNode;
}

interface DataTableProps<T> {
    columns: Column<T>[];
    data: T[];
    emptyMessage?: string;
}

export default function DataTable<T extends Record<string, unknown>>({
    columns,
    data,
    emptyMessage = 'No records found.',
}: DataTableProps<T>) {
    if (data.length === 0) {
        return (
            <div className="text-center py-12 text-surface-400 dark:text-surface-500">
                {emptyMessage}
            </div>
        );
    }

    return (
        <div className="overflow-x-auto rounded-lg border border-surface-200 dark:border-surface-700">
            <table className="w-full text-sm">
                <thead>
                    <tr className="bg-surface-50 dark:bg-surface-800 border-b border-surface-200 dark:border-surface-700">
                        {columns.map((col) => (
                            <th
                                key={col.key}
                                className="text-left px-4 py-3 font-medium text-surface-600 dark:text-surface-300 whitespace-nowrap"
                            >
                                {col.label}
                            </th>
                        ))}
                    </tr>
                </thead>
                <tbody className="divide-y divide-surface-100 dark:divide-surface-800">
                    {data.map((row, i) => (
                        <tr key={i} className="hover:bg-surface-50 dark:hover:bg-surface-800/50 transition-colors">
                            {columns.map((col) => (
                                <td key={col.key} className="px-4 py-3 text-surface-700 dark:text-surface-300 whitespace-nowrap">
                                    {col.render ? col.render(row) : String(row[col.key] ?? '-')}
                                </td>
                            ))}
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}
