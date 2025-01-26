import { AES, enc } from 'crypto-js';

export const encryptionService = {
	// Use device-specific info + user password to create encryption key
	async generateKey(password: string): Promise<string> {
		const deviceId = await this.getDeviceFingerprint();
		return `${deviceId}-${password}`;
	},

	async getDeviceFingerprint(): Promise<string> {
		const { userAgent, language } = navigator;
		const { width, height } = screen;
		const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;

		const fingerprint = `${userAgent}-${language}-${width}x${height}-${timeZone}`;
		const encoder = new TextEncoder();
		const data = encoder.encode(fingerprint);
		const hash = await crypto.subtle.digest('SHA-256', data);
		return Array.from(new Uint8Array(hash))
			.map((b) => b.toString(16).padStart(2, '0'))
			.join('');
	},

	encrypt(data: any, key: string): string {
		const jsonStr = JSON.stringify(data);
		return AES.encrypt(jsonStr, key).toString();
	},

	decrypt(encrypted: string, key: string): any {
		const decrypted = AES.decrypt(encrypted, key);
		return JSON.parse(decrypted.toString(enc.Utf8));
	},
};
