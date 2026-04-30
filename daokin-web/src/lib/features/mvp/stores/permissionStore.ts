import { get, writable } from 'svelte/store';

import type { ApiClient } from '../api/client';
import type { PermissionRecord, StepState } from '../types';

import { toErrorMessage } from './errors';

interface PermissionState {
	granteeWallet: string;
	scope: string;
	permission: PermissionRecord | null;
	grantStep: StepState;
	revokeStep: StepState;
}

const initialState: PermissionState = {
	granteeWallet: '0x2222222222222222222222222222222222222222',
	scope: 'view',
	permission: null,
	grantStep: { status: 'idle', error: '' },
	revokeStep: { status: 'idle', error: '' }
};

export function createPermissionStore() {
	const store = writable<PermissionState>(initialState);
	const { subscribe, update } = store;

	return {
		subscribe,
		setGranteeWallet(granteeWallet: string): void {
			update((state) => ({ ...state, granteeWallet }));
		},
		setScope(scope: string): void {
			update((state) => ({ ...state, scope }));
		},
		setGrantError(message: string): void {
			update((state) => ({
				...state,
				grantStep: { status: 'error', error: message }
			}));
		},
		resetForWalletChange(): void {
			update((state) => ({
				...state,
				permission: null,
				grantStep: { status: 'idle', error: '' },
				revokeStep: { status: 'idle', error: '' }
			}));
		},
		resetForArtifactChange(): void {
			update((state) => ({
				...state,
				permission: null,
				grantStep: { status: 'idle', error: '' },
				revokeStep: { status: 'idle', error: '' }
			}));
		},
		async grantPermission(
			api: ApiClient,
			artifactId: string,
			granterWallet: string
		): Promise<boolean> {
			const state = get(store);
			if (state.scope.trim() === '') {
				update((current) => ({
					...current,
					grantStep: { status: 'error', error: 'Permission scope is required.' }
				}));
				return false;
			}

			update((current) => ({ ...current, grantStep: { status: 'loading', error: '' } }));
			try {
				const permission = await api.permission.grant({
					artifactId,
					granterWallet,
					granteeWallet: state.granteeWallet,
					scope: state.scope
				});
				update((current) => ({
					...current,
					permission,
					grantStep: { status: 'success', error: '' },
					revokeStep: { status: 'idle', error: '' }
				}));
				return true;
			} catch (error) {
				update((current) => ({
					...current,
					grantStep: { status: 'error', error: toErrorMessage(error) }
				}));
				return false;
			}
		},
		async revokePermission(api: ApiClient, granterWallet: string, reason = ''): Promise<boolean> {
			const state = get(store);
			if (!state.permission || state.permission.status !== 'active') {
				update((current) => ({
					...current,
					revokeStep: { status: 'error', error: 'Active permission is required before revoke.' }
				}));
				return false;
			}

			update((current) => ({ ...current, revokeStep: { status: 'loading', error: '' } }));
			try {
				const permission = await api.permission.revoke({
					permissionId: state.permission.id,
					granterWallet,
					reason
				});
				update((current) => ({
					...current,
					permission,
					revokeStep: { status: 'success', error: '' }
				}));
				return true;
			} catch (error) {
				update((current) => ({
					...current,
					revokeStep: { status: 'error', error: toErrorMessage(error) }
				}));
				return false;
			}
		}
	};
}
