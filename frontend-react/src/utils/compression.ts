import { compress, decompress } from 'lz-string';

export const compressionUtils = {
	compressData(data: any): string {
		return compress(JSON.stringify(data));
	},

	decompressData<T>(compressed: string): T {
		return JSON.parse(decompress(compressed));
	},

	// Helper to estimate storage savings
	getCompressionStats(data: any) {
		const raw = JSON.stringify(data);
		const compressed = this.compressData(data);
		return {
			originalSize: raw.length,
			compressedSize: compressed.length,
			savings:
				(((raw.length - compressed.length) / raw.length) * 100).toFixed(
					1
				) + '%',
		};
	},
};
