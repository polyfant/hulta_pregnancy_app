import { wordlist } from './wordlist'; // BIP39 wordlist

export const securityService = {
	generateRecoveryPhrase(): string {
		// Generate 12 random words (easier to remember than a password)
		const words = [];
		for (let i = 0; i < 12; i++) {
			const index = Math.floor(Math.random() * wordlist.length);
			words.push(wordlist[index]);
		}
		return words.join(' ');
	},

	async phraseToKey(phrase: string): Promise<string> {
		// Convert recovery phrase to encryption key
		const encoder = new TextEncoder();
		const data = encoder.encode(phrase);
		const hash = await crypto.subtle.digest('SHA-256', data);
		return Array.from(new Uint8Array(hash))
			.map((b) => b.toString(16).padStart(2, '0'))
			.join('');
	},

	// Simple interface for users
	async protect(data: any, recoveryPhrase: string) {
		const key = await this.phraseToKey(recoveryPhrase);
		return encryptionService.encrypt(data, key);
	},

	async recover(encrypted: any, recoveryPhrase: string) {
		const key = await this.phraseToKey(recoveryPhrase);
		return encryptionService.decrypt(encrypted, key);
	},
};
