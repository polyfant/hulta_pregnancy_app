import { useState } from 'react';

export function useSearchHistory(key: string, maxItems = 5) {
    const [history, setHistory] = useState<string[]>(() => {
        const saved = localStorage.getItem(`search_history_${key}`);
        return saved ? JSON.parse(saved) : [];
    });

    const addToHistory = (term: string) => {
        setHistory(prev => {
            const newHistory = [term, ...prev.filter(t => t !== term)].slice(0, maxItems);
            localStorage.setItem(`search_history_${key}`, JSON.stringify(newHistory));
            return newHistory;
        });
    };

    return { history, addToHistory };
} 