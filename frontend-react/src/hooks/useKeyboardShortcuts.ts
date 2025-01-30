export function useKeyboardShortcuts() {
    useEffect(() => {
        const handleKeyPress = (e: KeyboardEvent) => {
            if (e.ctrlKey && e.key === 'k') {
                e.preventDefault();
                // Focus search
            }
            if (e.ctrlKey && e.key === 'n') {
                e.preventDefault();
                // New horse
            }
        };

        window.addEventListener('keydown', handleKeyPress);
        return () => window.removeEventListener('keydown', handleKeyPress);
    }, []);
} 