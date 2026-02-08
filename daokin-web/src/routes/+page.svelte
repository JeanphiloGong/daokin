<script lang="ts">
	import { onMount } from 'svelte';
	import { derived } from 'svelte/store';

	import { createApiClient } from '$lib/features/mvp/api/client';
	import StepCard from '$lib/features/mvp/components/StepCard.svelte';
	import { createArtifactStore } from '$lib/features/mvp/stores/artifactStore';
	import { createAuthStore } from '$lib/features/mvp/stores/authStore';
	import { createDaoStore } from '$lib/features/mvp/stores/daoStore';
	import { createPaymentStore } from '$lib/features/mvp/stores/paymentStore';

	const api = createApiClient({ mode: 'mock' });

	const authStore = createAuthStore();
	const artifactStore = createArtifactStore();
	const daoStore = createDaoStore();
	const paymentStore = createPaymentStore();

	const selectedDao = derived(
		daoStore,
		($daoStore) => $daoStore.daos.find((dao) => dao.id === $daoStore.selectedDaoId) ?? null
	);
	const canCreateArtifact = derived(authStore, ($authStore) => Boolean($authStore.session));
	const canJoinDao = derived(
		[authStore, artifactStore, daoStore],
		([$authStore, $artifactStore, $daoStore]) =>
			Boolean($authStore.session && $artifactStore.artifact && $daoStore.selectedDaoId)
	);
	const canTip = derived(
		[authStore, artifactStore, daoStore],
		([$authStore, $artifactStore, $daoStore]) =>
			Boolean($authStore.session && $artifactStore.artifact && $daoStore.membership)
	);
	const canViewAttribution = derived(artifactStore, ($artifactStore) =>
		Boolean($artifactStore.artifact)
	);

	onMount(async () => {
		await daoStore.loadDaos(api);
	});

	async function connectWallet(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		const connected = await authStore.connectWallet(api);
		if (!connected) {
			return;
		}

		artifactStore.resetForWalletChange();
		daoStore.resetForWalletChange();
		paymentStore.resetForWalletChange();
	}

	async function createArtifact(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		const session = $authStore.session;
		if (!session) {
			artifactStore.setError('Connect wallet before creating artifacts.');
			return;
		}

		const created = await artifactStore.createArtifact(api, session.userId);
		if (!created) {
			return;
		}

		daoStore.resetForArtifactChange();
		paymentStore.resetForArtifactChange();
	}

	async function joinDao(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		const session = $authStore.session;
		const artifact = $artifactStore.artifact;
		if (!session || !artifact) {
			daoStore.setJoinError('Create an artifact before joining a DAO.');
			return;
		}

		const joined = await daoStore.joinDao(api, session.userId, artifact.id);
		if (!joined) {
			return;
		}

		paymentStore.resetForJoinChange();
	}

	async function tipPay(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		const session = $authStore.session;
		const artifact = $artifactStore.artifact;
		const membership = $daoStore.membership;
		if (!session || !artifact || !membership) {
			paymentStore.setTipError('Join DAO before tipping.');
			return;
		}

		const stewardRecord = $paymentStore.attributions.find((item) => item.role === 'dao steward');
		const toUserId = stewardRecord?.contributorId ?? `${membership.daoId}-treasury`;
		await paymentStore.tipPay(api, {
			artifactId: artifact.id,
			fromUserId: session.userId,
			toUserId
		});
	}

	async function viewAttribution(): Promise<void> {
		const artifact = $artifactStore.artifact;
		if (!artifact) {
			paymentStore.setAttributionError('Create an artifact first.');
			return;
		}

		await paymentStore.loadAttributions(api, artifact.id);
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

	{#if $daoStore.globalError}
		<p class="global-error" role="alert">{$daoStore.globalError}</p>
	{/if}

	<section class="step-stack">
		<StepCard
			title="1) Wallet auth"
			description="Connect wallet and establish a simple user session."
			status={$authStore.step.status}
			error={$authStore.step.error}
		>
			<form
				class="stack-form"
				onsubmit={connectWallet}
				aria-busy={$authStore.step.status === 'loading'}
			>
				<label>
					<span>Wallet address</span>
					<input
						type="text"
						value={$authStore.walletAddress}
						oninput={(event) =>
							authStore.setWalletAddress((event.currentTarget as HTMLInputElement).value)}
						placeholder="0xabc12345"
						required
					/>
				</label>
				<button type="submit" disabled={$authStore.step.status === 'loading'}>
					Connect wallet
				</button>
			</form>
			{#if $authStore.session}
				<p class="detail">
					Session: {$authStore.session.displayName} ({$authStore.session.userId})
				</p>
			{/if}
		</StepCard>

		<StepCard
			title="2) Create artifact"
			description="Submit title + description and receive a traceable artifact ID."
			status={$artifactStore.step.status}
			error={$artifactStore.step.error}
		>
			<form
				class="stack-form"
				onsubmit={createArtifact}
				aria-busy={$artifactStore.step.status === 'loading'}
			>
				<label>
					<span>Title</span>
					<input
						type="text"
						value={$artifactStore.artifactTitle}
						oninput={(event) =>
							artifactStore.setArtifactTitle((event.currentTarget as HTMLInputElement).value)}
						minlength="3"
						required
					/>
				</label>
				<label>
					<span>Description</span>
					<textarea
						value={$artifactStore.artifactDescription}
						oninput={(event) =>
							artifactStore.setArtifactDescription(
								(event.currentTarget as HTMLTextAreaElement).value
							)}
						rows="3"
					></textarea>
				</label>
				<button
					type="submit"
					disabled={!$canCreateArtifact || $artifactStore.step.status === 'loading'}
				>
					Create artifact
				</button>
			</form>
			{#if $artifactStore.artifact}
				<p class="detail">Artifact ID: {$artifactStore.artifact.id}</p>
			{/if}
		</StepCard>

		<StepCard
			title="3) Join DAO"
			description="Join a DAO circle so artifact contribution can be contextualized."
			status={$daoStore.joinStep.status}
			error={$daoStore.joinStep.error}
		>
			<form
				class="stack-form"
				onsubmit={joinDao}
				aria-busy={$daoStore.joinStep.status === 'loading'}
			>
				<label>
					<span>DAO</span>
					<select
						value={$daoStore.selectedDaoId}
						onchange={(event) =>
							daoStore.setSelectedDaoId((event.currentTarget as HTMLSelectElement).value)}
						required
					>
						{#each $daoStore.daos as dao (dao.id)}
							<option value={dao.id}>{dao.name}</option>
						{/each}
					</select>
				</label>
				<button type="submit" disabled={!$canJoinDao || $daoStore.joinStep.status === 'loading'}>
					Join DAO
				</button>
			</form>
			{#if $daoStore.membership}
				<p class="detail">Joined: {$selectedDao?.name ?? $daoStore.membership.daoId}</p>
			{/if}
		</StepCard>

		<StepCard
			title="4) Tip / pay"
			description="Send a small payment to test value flow after joining a DAO."
			status={$paymentStore.tipStep.status}
			error={$paymentStore.tipStep.error}
		>
			<form
				class="stack-form"
				onsubmit={tipPay}
				aria-busy={$paymentStore.tipStep.status === 'loading'}
			>
				<label>
					<span>Amount (USDC)</span>
					<input
						type="number"
						min="0.01"
						step="0.01"
						value={$paymentStore.tipAmount}
						oninput={(event) =>
							paymentStore.setTipAmount((event.currentTarget as HTMLInputElement).value)}
						required
					/>
				</label>
				<button type="submit" disabled={!$canTip || $paymentStore.tipStep.status === 'loading'}>
					Send tip
				</button>
			</form>
			{#if $paymentStore.tip}
				<p class="detail">Payment confirmed: {$paymentStore.tip.paymentId}</p>
			{/if}
		</StepCard>

		<StepCard
			title="5) View attribution"
			description="Fetch attribution records to verify ownership and contributor shares."
			status={$paymentStore.attributionStep.status}
			error={$paymentStore.attributionStep.error}
		>
			<button type="button" onclick={viewAttribution} disabled={!$canViewAttribution}
				>Refresh attribution</button
			>

			{#if $paymentStore.attributions.length === 0}
				<p class="detail">No attribution loaded yet.</p>
			{:else}
				<table class="attribution-table">
					<caption class="sr-only">Attribution contributors and share percentages</caption>
					<thead>
						<tr>
							<th scope="col">Contributor</th>
							<th scope="col">Role</th>
							<th scope="col">Wallet</th>
							<th scope="col">Share</th>
						</tr>
					</thead>
					<tbody>
						{#each $paymentStore.attributions as item (item.contributorId + item.role)}
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
	.page-shell {
		max-width: 960px;
		margin: 0 auto;
		padding: var(--space-6) var(--space-4) calc(var(--space-6) + 0.5rem);
	}

	.page-header h1 {
		font-size: 1.5rem;
		margin: 0;
	}

	.page-header p {
		margin: 0.45rem 0 0;
		color: var(--color-pill-neutral-text);
	}

	.mode-note {
		font-size: 0.9rem;
		color: var(--color-text-muted);
	}

	.global-error {
		margin: var(--space-4) 0;
		padding: 0.75rem 0.9rem;
		border: 1px solid var(--color-error-border);
		background: var(--color-error-bg);
		border-radius: 0.55rem;
		color: var(--color-error-text);
	}

	.step-stack {
		display: grid;
		gap: 0.9rem;
		margin-top: var(--space-5);
	}

	.stack-form {
		display: grid;
		gap: var(--space-2);
	}

	label {
		display: grid;
		gap: var(--space-1);
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
		border: 1px solid var(--color-border-strong);
		border-radius: var(--radius-sm);
		padding: 0.55rem 0.65rem;
		background: var(--color-surface);
		color: var(--color-text);
	}

	button {
		width: fit-content;
		padding: 0.5rem 0.8rem;
		border-radius: var(--radius-sm);
		border: 1px solid var(--color-primary);
		background: var(--color-primary);
		color: var(--color-primary-contrast);
		font-weight: 600;
		cursor: pointer;
	}

	button:disabled {
		opacity: 0.55;
		cursor: not-allowed;
	}

	.detail {
		font-size: 0.88rem;
		color: var(--color-pill-neutral-text);
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
		border-bottom: 1px solid var(--color-pill-neutral-bg);
	}

	.sr-only {
		border: 0;
		clip: rect(0 0 0 0);
		height: 1px;
		margin: -1px;
		overflow: hidden;
		padding: 0;
		position: absolute;
		white-space: nowrap;
		width: 1px;
	}

	@media (max-width: 640px) {
		.page-shell {
			padding-top: var(--space-5);
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
