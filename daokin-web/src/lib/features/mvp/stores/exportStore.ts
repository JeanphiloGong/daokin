import { writable } from 'svelte/store';

import type { ApiClient } from '../api/client';
import type { StepState, UserExportBundle } from '../types';

import { toErrorMessage } from './errors';

interface ExportState {
	bundle: UserExportBundle | null;
	step: StepState;
}

const initialState: ExportState = {
	bundle: null,
	step: { status: 'idle', error: '' }
};

export function createExportStore() {
	const store = writable<ExportState>(initialState);
	const { subscribe, update } = store;

	return {
		subscribe,
		resetForWalletChange(): void {
			update(() => initialState);
		},
		async loadUserExport(api: ApiClient, wallet: string): Promise<boolean> {
			update((state) => ({ ...state, step: { status: 'loading', error: '' } }));
			try {
				const bundle = await api.userExport.get(wallet);
				update(() => ({
					bundle,
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
