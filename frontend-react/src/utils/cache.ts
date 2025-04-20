export class Cache {
    private storage = new Map<string, { value: unknown; expiry: number }>();

    set<T>(key: string, value: T, ttlMinutes = 5): void {
        this.storage.set(key, {
            value,
            expiry: Date.now() + ttlMinutes * 60 * 1000,
        });
    }

    get<T>(key: string): T | null {
        const item = this.storage.get(key);
        if (!item) return null;
        if (Date.now() > item.expiry) {
            this.storage.delete(key);
            return null;
        }
        return item.value as T;
    }
} 