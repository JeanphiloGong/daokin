import { get, writable } from 'svelte/store';

import type { ApiClient } from '../api/client';
import type { Artifact, StepState } from '../types';

import { toErrorMessage } from './errors';

interface ArtifactState {
	artifactTitle: string;
	artifactDescription: string;
	artifact: Artifact | null;
	step: StepState;
}

const initialState: ArtifactState = {
	artifactTitle: 'My First Dao Artifact',
	artifactDescription: 'A simple manifesto draft for contributor-owned communities.',
	artifact: null,
	step: { status: 'idle', error: '' }
};

export function createArtifactStore() {
	const store = writable<ArtifactState>(initialState);
	const { subscribe, update } = store;

	return {
		subscribe,
		setArtifactTitle(artifactTitle: string): void {
			update((state) => ({ ...state, artifactTitle }));
		},
		setArtifactDescription(artifactDescription: string): void {
			update((state) => ({ ...state, artifactDescription }));
		},
		setError(message: string): void {
			update((state) => ({
				...state,
				step: { status: 'error', error: message }
			}));
		},
		resetForWalletChange(): void {
			update((state) => ({
				...state,
				artifact: null,
				step: { status: 'idle', error: '' }
			}));
		},
		async createArtifact(api: ApiClient, creatorId: string): Promise<boolean> {
			if (creatorId.trim() === '') {
				update((state) => ({
					...state,
					step: {
						status: 'error',
						error: 'Connect wallet before creating artifacts.'
					}
				}));
				return false;
			}

			update((state) => ({ ...state, step: { status: 'loading', error: '' } }));
			const state = get(store);
			try {
				const artifact = await api.artifact.create({
					title: state.artifactTitle,
					description: state.artifactDescription,
					creatorId
				});
				update((current) => ({
					...current,
					artifact,
					step: { status: 'success', error: '' }
				}));
				return true;
			} catch (error) {
				update((current) => ({
					...current,
					step: { status: 'error', error: toErrorMessage(error) }
				}));
				return false;
			}
		}
	};
}
