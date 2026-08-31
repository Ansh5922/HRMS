import { Loader } from 'lucide-react';

export default function LoadingSpinner() {
    return (
        <div className="flex items-center justify-center py-20">
            <Loader className="w-6 h-6 animate-spin text-primary-500" />
        </div>
    );
}
