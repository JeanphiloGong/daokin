import { get, writable } from 'svelte/store';

import type { ApiClient } from '../api/client';
import type { DaoMembership, DaoSummary, StepState } from '../types';

import { toErrorMessage } from './errors';

interface DaoState {
	globalError: string;
	daos: DaoSummary[];
	selectedDaoId: string;
	membership: DaoMembership | null;
	joinStep: StepState;
}

const initialState: DaoState = {
	globalError: '',
	daos: [],
	selectedDaoId: '',
	membership: null,
	joinStep: { status: 'idle', error: '' }
};

export function createDaoStore() {
	const store = writable<DaoState>(initialState);
	const { subscribe, update } = store;

	return {
		subscribe,
		setSelectedDaoId(selectedDaoId: string): void {
			update((state) => ({ ...state, selectedDaoId }));
		},
		setJoinError(message: string): void {
			update((state) => ({
				...state,
				joinStep: { status: 'error', error: message }
			}));
		},
		resetForWalletChange(): void {
			update((state) => ({
				...state,
				membership: null,
				joinStep: { status: 'idle', error: '' }
			}));
		},
		resetForArtifactChange(): void {
			update((state) => ({
				...state,
				membership: null,
				joinStep: { status: 'idle', error: '' }
			}));
		},
		async loadDaos(api: ApiClient): Promise<void> {
			try {
				const daos = await api.dao.list();
				update((state) => ({
					...state,
					daos,
					selectedDaoId: state.selectedDaoId || daos[0]?.id || '',
					globalError: ''
				}));
			} catch (error) {
				update((state) => ({
					...state,
					globalError: `Failed to load DAO options: ${toErrorMessage(error)}`
				}));
			}
		},
		async joinDao(api: ApiClient, userId: string, artifactId: string): Promise<boolean> {
			const state = get(store);
			if (state.selectedDaoId.trim() === '') {
				update((current) => ({
					...current,
					joinStep: { status: 'error', error: 'Pick a DAO first.' }
				}));
				return false;
			}

			update((current) => ({ ...current, joinStep: { status: 'loading', error: '' } }));
			try {
				const membership = await api.dao.join({
					daoId: state.selectedDaoId,
					userId,
					artifactId
				});
				update((current) => ({
					...current,
					membership,
					joinStep: { status: 'success', error: '' }
				}));
				return true;
			} catch (error) {
				update((current) => ({
					...current,
					joinStep: { status: 'error', error: toErrorMessage(error) }
				}));
				return false;
			}
		}
	};
}
