<script lang="ts">
	import { onMount } from 'svelte';
	import StepCard from '$lib/features/mvp/components/StepCard.svelte';
	import { createApiClient } from '$lib/features/mvp/api/client';
	import { ApiError } from '$lib/features/mvp/types';
	import type {
		Artifact,
		AttributionRecord,
		AuthSession,
		DaoMembership,
		DaoSummary,
		StepStatus,
		TipPayment
	} from '$lib/features/mvp/types';

	type StepKey = 'walletAuth' | 'createArtifact' | 'joinDao' | 'tipPay' | 'viewAttribution';
	interface StepState {
		status: StepStatus;
		error: string;
	}

	const api = createApiClient({ mode: 'mock' });

	const state = $state({
		globalError: '',
		walletAddress: '0xabc12345',
		session: null as AuthSession | null,
		artifactTitle: 'My First Dao Artifact',
		artifactDescription: 'A simple manifesto draft for contributor-owned communities.',
		artifact: null as Artifact | null,
		daos: [] as DaoSummary[],
		selectedDaoId: '',
		membership: null as DaoMembership | null,
		tipAmount: '1.00',
		tip: null as TipPayment | null,
		attributions: [] as AttributionRecord[],
		steps: {
			walletAuth: { status: 'idle', error: '' },
			createArtifact: { status: 'idle', error: '' },
			joinDao: { status: 'idle', error: '' },
			tipPay: { status: 'idle', error: '' },
			viewAttribution: { status: 'idle', error: '' }
		} as Record<StepKey, StepState>
	});

	const selectedDao = $derived(state.daos.find((dao) => dao.id === state.selectedDaoId) ?? null);
	const canCreateArtifact = $derived(Boolean(state.session));
	const canJoinDao = $derived(Boolean(state.session && state.artifact && state.selectedDaoId));
	const canTip = $derived(Boolean(state.session && state.membership && state.artifact));
	const canViewAttribution = $derived(Boolean(state.artifact));

	onMount(async () => {
		await loadDaos();
	});

	function resetDownstream(from: StepKey): void {
		const order: StepKey[] = [
			'walletAuth',
			'createArtifact',
			'joinDao',
			'tipPay',
			'viewAttribution'
		];
		const index = order.indexOf(from);
		for (const step of order.slice(index + 1)) {
			state.steps[step] = { status: 'idle', error: '' };
		}
		if (from === 'walletAuth') {
			state.artifact = null;
			state.membership = null;
			state.tip = null;
			state.attributions = [];
		}
		if (from === 'createArtifact') {
			state.membership = null;
			state.tip = null;
			state.attributions = [];
		}
		if (from === 'joinDao') {
			state.tip = null;
		}
	}

	function setStep(step: StepKey, status: StepStatus, error = ''): void {
		state.steps[step] = { status, error };
	}

	function toErrorMessage(error: unknown): string {
		if (error instanceof ApiError) {
			return error.message;
		}
		if (error instanceof Error) {
			return error.message;
		}
		return 'Unexpected error. Please retry.';
	}

	async function loadDaos(): Promise<void> {
		try {
			const daos = await api.dao.list();
			state.daos = daos;
			if (!state.selectedDaoId && daos.length > 0) {
				state.selectedDaoId = daos[0].id;
			}
			state.globalError = '';
		} catch (error) {
			state.globalError = `Failed to load DAO options: ${toErrorMessage(error)}`;
		}
	}

	async function connectWallet(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		setStep('walletAuth', 'loading');
		try {
			state.session = await api.auth.walletConnect(state.walletAddress);
			setStep('walletAuth', 'success');
			resetDownstream('walletAuth');
		} catch (error) {
			setStep('walletAuth', 'error', toErrorMessage(error));
		}
	}

	async function createArtifact(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		if (!state.session) {
			setStep('createArtifact', 'error', 'Connect wallet before creating artifacts.');
			return;
		}

		setStep('createArtifact', 'loading');
		try {
			state.artifact = await api.artifact.create({
				title: state.artifactTitle,
				description: state.artifactDescription,
				creatorId: state.session.userId
			});
			setStep('createArtifact', 'success');
			resetDownstream('createArtifact');
		} catch (error) {
			setStep('createArtifact', 'error', toErrorMessage(error));
		}
	}

	async function joinDao(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		if (!state.session || !state.artifact) {
			setStep('joinDao', 'error', 'Create an artifact before joining a DAO.');
			return;
		}
		if (!state.selectedDaoId) {
			setStep('joinDao', 'error', 'Pick a DAO first.');
			return;
		}

		setStep('joinDao', 'loading');
		try {
			state.membership = await api.dao.join({
				daoId: state.selectedDaoId,
				userId: state.session.userId,
				artifactId: state.artifact.id
			});
			setStep('joinDao', 'success');
			resetDownstream('joinDao');
		} catch (error) {
			setStep('joinDao', 'error', toErrorMessage(error));
		}
	}

	async function tipPay(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		if (!state.session || !state.artifact || !state.membership) {
			setStep('tipPay', 'error', 'Join DAO before tipping.');
			return;
		}

		const amount = Number.parseFloat(state.tipAmount);
		if (!Number.isFinite(amount) || amount <= 0) {
			setStep('tipPay', 'error', 'Tip amount must be a positive number.');
			return;
		}

		setStep('tipPay', 'loading');
		try {
			const stewardRecord = state.attributions.find((item) => item.role === 'dao steward');
			state.tip = await api.payment.tip({
				artifactId: state.artifact.id,
				fromUserId: state.session.userId,
				toUserId: stewardRecord?.contributorId ?? `${state.membership.daoId}-treasury`,
				amount
			});
			setStep('tipPay', 'success');
		} catch (error) {
			setStep('tipPay', 'error', toErrorMessage(error));
		}
	}

	async function viewAttribution(): Promise<void> {
		if (!state.artifact) {
			setStep('viewAttribution', 'error', 'Create an artifact first.');
			return;
		}

		setStep('viewAttribution', 'loading');
		try {
			state.attributions = await api.attribution.listByArtifact(state.artifact.id);
			setStep('viewAttribution', 'success');
		} catch (error) {
			setStep('viewAttribution', 'error', toErrorMessage(error));
		}
	}
</script>

<svelte:head>
	<title>DaoKin MVP Flow</title>
</svelte:head>

<main class="page-shell">
	<header class="page-header">
		<h1>DaoKin MVP Flow</h1>
		<p>wallet auth -&gt; create artifact -&gt; join dao -&gt; tip/pay -&gt; view attribution</p>
		<p class="mode-note">API mode: mock (v1 endpoints already defined in client)</p>
	</header>

	{#if state.globalError}
		<p class="global-error" role="alert">{state.globalError}</p>
	{/if}

	<section class="step-stack">
		<StepCard
			title="1) Wallet auth"
			description="Connect wallet and establish a simple user session."
			status={state.steps.walletAuth.status}
			error={state.steps.walletAuth.error}
		>
			<form class="stack-form" onsubmit={connectWallet}>
				<label>
					<span>Wallet address</span>
					<input type="text" bind:value={state.walletAddress} placeholder="0xabc12345" required />
				</label>
				<button type="submit" disabled={state.steps.walletAuth.status === 'loading'}>
					Connect wallet
				</button>
			</form>
			{#if state.session}
				<p class="detail">Session: {state.session.displayName} ({state.session.userId})</p>
			{/if}
		</StepCard>

		<StepCard
			title="2) Create artifact"
			description="Submit title + description and receive a traceable artifact ID."
			status={state.steps.createArtifact.status}
			error={state.steps.createArtifact.error}
		>
			<form class="stack-form" onsubmit={createArtifact}>
				<label>
					<span>Title</span>
					<input type="text" bind:value={state.artifactTitle} minlength="3" required />
				</label>
				<label>
					<span>Description</span>
					<textarea bind:value={state.artifactDescription} rows="3"></textarea>
				</label>
				<button
					type="submit"
					disabled={!canCreateArtifact || state.steps.createArtifact.status === 'loading'}
				>
					Create artifact
				</button>
			</form>
			{#if state.artifact}
				<p class="detail">Artifact ID: {state.artifact.id}</p>
			{/if}
		</StepCard>

		<StepCard
			title="3) Join DAO"
			description="Join a DAO circle so artifact contribution can be contextualized."
			status={state.steps.joinDao.status}
			error={state.steps.joinDao.error}
		>
			<form class="stack-form" onsubmit={joinDao}>
				<label>
					<span>DAO</span>
					<select bind:value={state.selectedDaoId} required>
						{#each state.daos as dao}
							<option value={dao.id}>{dao.name}</option>
						{/each}
					</select>
				</label>
				<button type="submit" disabled={!canJoinDao || state.steps.joinDao.status === 'loading'}>
					Join DAO
				</button>
			</form>
			{#if state.membership}
				<p class="detail">Joined: {selectedDao?.name ?? state.membership.daoId}</p>
			{/if}
		</StepCard>

		<StepCard
			title="4) Tip / pay"
			description="Send a small payment to test value flow after joining a DAO."
			status={state.steps.tipPay.status}
			error={state.steps.tipPay.error}
		>
			<form class="stack-form" onsubmit={tipPay}>
				<label>
					<span>Amount (USDC)</span>
					<input type="number" min="0.01" step="0.01" bind:value={state.tipAmount} required />
				</label>
				<button type="submit" disabled={!canTip || state.steps.tipPay.status === 'loading'}
					>Send tip</button
				>
			</form>
			{#if state.tip}
				<p class="detail">Payment confirmed: {state.tip.paymentId}</p>
			{/if}
		</StepCard>

		<StepCard
			title="5) View attribution"
			description="Fetch attribution records to verify ownership and contributor shares."
			status={state.steps.viewAttribution.status}
			error={state.steps.viewAttribution.error}
		>
			<button type="button" onclick={viewAttribution} disabled={!canViewAttribution}
				>Refresh attribution</button
			>

			{#if state.attributions.length === 0}
				<p class="detail">No attribution loaded yet.</p>
			{:else}
				<table class="attribution-table">
					<thead>
						<tr>
							<th>Contributor</th>
							<th>Role</th>
							<th>Wallet</th>
							<th>Share</th>
						</tr>
					</thead>
					<tbody>
						{#each state.attributions as item}
							<tr>
								<td>{item.contributorName}</td>
								<td>{item.role}</td>
								<td>{item.walletAddress}</td>
								<td>{item.sharePercent}%</td>
							</tr>
						{/each}
					</tbody>
				</table>
			{/if}
		</StepCard>
	</section>
</main>

<style>
	:global(body) {
		background: #f8fafc;
		color: #0f172a;
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	}

	.page-shell {
		max-width: 960px;
		margin: 0 auto;
		padding: 2rem 1rem 2.5rem;
	}

	.page-header h1 {
		font-size: 1.5rem;
		margin: 0;
	}

	.page-header p {
		margin: 0.45rem 0 0;
		color: #334155;
	}

	.mode-note {
		font-size: 0.9rem;
		color: #475569;
	}

	.global-error {
		margin: 1rem 0;
		padding: 0.75rem 0.9rem;
		border: 1px solid #fecaca;
		background: #fee2e2;
		border-radius: 0.55rem;
		color: #991b1b;
	}

	.step-stack {
		display: grid;
		gap: 0.9rem;
		margin-top: 1.25rem;
	}

	.stack-form {
		display: grid;
		gap: 0.65rem;
	}

	label {
		display: grid;
		gap: 0.35rem;
		font-size: 0.9rem;
		font-weight: 600;
	}

	input,
	textarea,
	select,
	button {
		font: inherit;
	}

	input,
	textarea,
	select {
		border: 1px solid #cbd5e1;
		border-radius: 0.5rem;
		padding: 0.55rem 0.65rem;
		background: #ffffff;
	}

	button {
		width: fit-content;
		padding: 0.5rem 0.8rem;
		border-radius: 0.5rem;
		border: 1px solid #0f172a;
		background: #0f172a;
		color: #ffffff;
		font-weight: 600;
		cursor: pointer;
	}

	button:disabled {
		opacity: 0.55;
		cursor: not-allowed;
	}

	.detail {
		font-size: 0.88rem;
		color: #334155;
	}

	.attribution-table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.88rem;
	}

	.attribution-table th,
	.attribution-table td {
		text-align: left;
		padding: 0.45rem 0.35rem;
		border-bottom: 1px solid #e2e8f0;
	}

	@media (max-width: 640px) {
		.page-shell {
			padding-top: 1.25rem;
		}

		button {
			width: 100%;
		}

		.attribution-table {
			display: block;
			overflow-x: auto;
		}
	}
</style>
