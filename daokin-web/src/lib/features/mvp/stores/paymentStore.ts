import { get, writable } from 'svelte/store';

import type { ApiClient } from '../api/client';
import type { AttributionRecord, StepState, TipPayment } from '../types';

import { toErrorMessage } from './errors';

interface PaymentState {
	tipAmount: string;
	tip: TipPayment | null;
	attributions: AttributionRecord[];
	tipStep: StepState;
	attributionStep: StepState;
}

const initialState: PaymentState = {
	tipAmount: '1.00',
	tip: null,
	attributions: [],
	tipStep: { status: 'idle', error: '' },
	attributionStep: { status: 'idle', error: '' }
};

interface TipPayload {
	artifactId: string;
	fromUserId: string;
	toUserId: string;
}

export function createPaymentStore() {
	const store = writable<PaymentState>(initialState);
	const { subscribe, update } = store;

	return {
		subscribe,
		setTipAmount(tipAmount: string): void {
			update((state) => ({ ...state, tipAmount }));
		},
		setTipError(message: string): void {
			update((state) => ({
				...state,
				tipStep: { status: 'error', error: message }
			}));
		},
		setAttributionError(message: string): void {
			update((state) => ({
				...state,
				attributionStep: { status: 'error', error: message }
			}));
		},
		resetForWalletChange(): void {
			update((state) => ({
				...state,
				tip: null,
				attributions: [],
				tipStep: { status: 'idle', error: '' },
				attributionStep: { status: 'idle', error: '' }
			}));
		},
		resetForArtifactChange(): void {
			update((state) => ({
				...state,
				tip: null,
				attributions: [],
				tipStep: { status: 'idle', error: '' },
				attributionStep: { status: 'idle', error: '' }
			}));
		},
		resetForJoinChange(): void {
			update((state) => ({
				...state,
				tip: null,
				tipStep: { status: 'idle', error: '' },
				attributionStep: { status: 'idle', error: '' }
			}));
		},
		async tipPay(api: ApiClient, payload: TipPayload): Promise<boolean> {
			const state = get(store);
			const amount = Number.parseFloat(state.tipAmount);
			if (!Number.isFinite(amount) || amount <= 0) {
				update((current) => ({
					...current,
					tipStep: { status: 'error', error: 'Tip amount must be a positive number.' }
				}));
				return false;
			}

			update((current) => ({ ...current, tipStep: { status: 'loading', error: '' } }));
			try {
				const tip = await api.payment.tip({
					...payload,
					amount
				});
				update((current) => ({
					...current,
					tip,
					tipStep: { status: 'success', error: '' }
				}));
				return true;
			} catch (error) {
				update((current) => ({
					...current,
					tipStep: { status: 'error', error: toErrorMessage(error) }
				}));
				return false;
			}
		},
		async loadAttributions(api: ApiClient, artifactId: string): Promise<boolean> {
			update((state) => ({ ...state, attributionStep: { status: 'loading', error: '' } }));
			try {
				const attributions = await api.attribution.listByArtifact(artifactId);
				update((state) => ({
					...state,
					attributions,
					attributionStep: { status: 'success', error: '' }
				}));
				return true;
			} catch (error) {
				update((state) => ({
					...state,
					attributionStep: { status: 'error', error: toErrorMessage(error) }
				}));
				return false;
			}
		}
	};
}
