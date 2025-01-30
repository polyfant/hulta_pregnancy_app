import { useState } from 'react';

export function useRetry(callback: () => Promise<void>, maxRetries = 3) {
    const [retries, setRetries] = useState(0);
    const [error, setError] = useState<Error | null>(null);

    const retry = async () => {
        try {
            await callback();
            setError(null);
            setRetries(0);
        } catch (err) {
            if (retries < maxRetries) {
                setRetries(r => r + 1);
                setTimeout(retry, Math.pow(2, retries) * 1000);
            } else {
                setError(err as Error);
            }
        }
    };

    return { retry, error, retries };
} 