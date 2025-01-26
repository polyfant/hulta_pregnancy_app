// Custom type definitions for idb

declare module 'idb' {
  export interface OpenDBOptions {
    upgrade?: (db: IDBPDatabase, oldVersion: number, newVersion: number | null) => void;
    blocked?: () => void;
    blocking?: () => void;
    terminated?: () => void;
  }

  export interface IDBPDatabase<T = any> {
    name: string;
    version: number;
    objectStoreNames: DOMStringList;
    
    // Add missing methods
    createObjectStore(name: string, options?: IDBObjectStoreParameters): IDBObjectStore;
    transaction(storeNames: string | string[], mode?: IDBTransactionMode): IDBPTransaction<T>;
  }

  export interface IDBPTransaction<T = any> {
    objectStore<K extends keyof T>(name: K): IDBPObjectStore<T[K]>;
    done: Promise<void>;
  }

  export interface IDBPObjectStore<T = any> {
    put(value: T, key?: IDBValidKey): Promise<IDBValidKey>;
    get(key: IDBValidKey): Promise<T | undefined>;
  }

  export function openDB<T = any>(
    name: string, 
    version?: number, 
    options?: OpenDBOptions
  ): Promise<IDBPDatabase<T>>;
}
