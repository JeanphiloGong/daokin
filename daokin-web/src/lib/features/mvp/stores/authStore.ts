import { get, writable } from 'svelte/store';

import type { ApiClient } from '../api/client';
import type { AuthSession, StepState } from '../types';

import { toErrorMessage } from './errors';

interface AuthState {
	walletAddress: string;
	session: AuthSession | null;
	step: StepState;
}

const initialState: AuthState = {
	walletAddress: '0x1111111111111111111111111111111111111111',
	session: null,
	step: { status: 'idle', error: '' }
};

export function createAuthStore() {
	const store = writable<AuthState>(initialState);
	const { subscribe, update } = store;

	return {
		subscribe,
		setWalletAddress(walletAddress: string): void {
			update((state) => ({ ...state, walletAddress }));
		},
		resetForWalletRefresh(): void {
			update((state) => ({
				...state,
				session: null,
				step: { status: 'idle', error: '' }
			}));
		},
		async connectWallet(api: ApiClient): Promise<boolean> {
			update((state) => ({ ...state, step: { status: 'loading', error: '' } }));
			const walletAddress = get(store).walletAddress;

			try {
				const session = await api.auth.walletConnect(walletAddress);
				update((state) => ({
					...state,
					session,
					step: { status: 'success', error: '' }
				}));
				return true;
			} catch (error) {
				update((state) => ({
					...state,
					step: { status: 'error', error: toErrorMessage(error) }
				}));
				return false;
			}
		}
	};
}
