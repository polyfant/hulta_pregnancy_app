// Types for IndexedDB data structures

export interface FoalGrowthData {
  foalId: number;
  weightData: number[];
  heightData: number[];
  timestamp: number;
}

export type ChartDataKey = `foal-${number}`;

export function getChartDataKey(foalId: number): ChartDataKey {
  return `foal-${foalId}`;
}
