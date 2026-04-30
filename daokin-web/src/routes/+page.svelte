<script lang="ts">
	import { env } from '$env/dynamic/public';
	import { onMount } from 'svelte';
	import { derived } from 'svelte/store';

	import { createApiClient } from '$lib/features/mvp/api/client';
	import StepCard from '$lib/features/mvp/components/StepCard.svelte';
	import { createArtifactStore } from '$lib/features/mvp/stores/artifactStore';
	import { createAuthStore } from '$lib/features/mvp/stores/authStore';
	import { createDaoStore } from '$lib/features/mvp/stores/daoStore';
	import { createExportStore } from '$lib/features/mvp/stores/exportStore';
	import { createPaymentStore } from '$lib/features/mvp/stores/paymentStore';
	import { createPermissionStore } from '$lib/features/mvp/stores/permissionStore';

	const apiMode = env.PUBLIC_DAOKIN_API_MODE === 'http' ? 'http' : 'mock';
	const apiBaseUrl = env.PUBLIC_DAOKIN_API_BASE_URL || undefined;
	const api = createApiClient({ mode: apiMode, baseUrl: apiBaseUrl });

	const authStore = createAuthStore();
	const artifactStore = createArtifactStore();
	const daoStore = createDaoStore();
	const permissionStore = createPermissionStore();
	const paymentStore = createPaymentStore();
	const exportStore = createExportStore();

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
	const canGrantPermission = derived([authStore, artifactStore], ([$authStore, $artifactStore]) =>
		Boolean($authStore.session && $artifactStore.artifact)
	);
	const canTip = derived(
		[authStore, artifactStore, permissionStore],
		([$authStore, $artifactStore, $permissionStore]) =>
			Boolean(
				$authStore.session &&
				$artifactStore.artifact &&
				$permissionStore.permission?.status === 'active'
			)
	);
	const canViewAttribution = derived(artifactStore, ($artifactStore) =>
		Boolean($artifactStore.artifact)
	);
	const canExport = derived(authStore, ($authStore) => Boolean($authStore.session));
	const canLeaveDao = derived([authStore, daoStore], ([$authStore, $daoStore]) =>
		Boolean($authStore.session && $daoStore.membership?.leftAt == null)
	);
	const permissionCardState = derived(permissionStore, ($permissionStore) => {
		if ($permissionStore.revokeStep.status !== 'idle') {
			return $permissionStore.revokeStep;
		}
		return $permissionStore.grantStep;
	});

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
		permissionStore.resetForWalletChange();
		paymentStore.resetForWalletChange();
		exportStore.resetForWalletChange();
	}

	async function publishArtifact(event: SubmitEvent): Promise<void> {
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
		permissionStore.resetForArtifactChange();
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

	async function grantPermission(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		const session = $authStore.session;
		const artifact = $artifactStore.artifact;
		if (!session || !artifact) {
			permissionStore.setGrantError('Create an artifact before granting permission.');
			return;
		}

		await permissionStore.grantPermission(api, artifact.id, session.userId);
	}

	async function revokePermission(): Promise<void> {
		const session = $authStore.session;
		if (!session) {
			permissionStore.setGrantError('Connect wallet before revoking permission.');
			return;
		}

		await permissionStore.revokePermission(api, session.userId, 'MVP manual revoke check');
	}

	async function tipPay(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		const session = $authStore.session;
		const artifact = $artifactStore.artifact;
		const permission = $permissionStore.permission;
		if (!session || !artifact || !permission || permission.status !== 'active') {
			paymentStore.setTipError('Grant active permission before recording a transfer.');
			return;
		}

		try {
			await api.auth.walletConnect(permission.granteeWallet);
		} catch (error) {
			paymentStore.setTipError(
				error instanceof Error ? error.message : 'Permission grantee authentication failed.'
			);
			return;
		}

		await paymentStore.tipPay(api, {
			artifactId: artifact.id,
			fromUserId: permission.granteeWallet,
			toUserId: session.userId
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

	async function exportUserData(): Promise<void> {
		const session = $authStore.session;
		if (!session) {
			return;
		}

		await exportStore.loadUserExport(api, session.userId);
	}

	async function leaveDao(): Promise<void> {
		const session = $authStore.session;
		if (!session) {
			daoStore.setLeaveError('Connect wallet before leaving a DAO.');
			return;
		}

		await daoStore.leaveDao(api, session.userId, 'MVP reversibility check');
	}
</script>

<svelte:head>
	<title>DaoKin MVP Flow</title>
</svelte:head>

<main class="page-shell">
	<header class="page-header">
		<h1>DaoKin MVP Flow</h1>
		<p>
			wallet auth -&gt; create artifact -&gt; join dao -&gt; grant permission -&gt; record transfer
			-&gt; attribution -&gt; export -&gt; leave dao
		</p>
		<p class="mode-note">
			API mode: {apiMode}{apiMode === 'http' && apiBaseUrl ? ` (${apiBaseUrl})` : ''}
		</p>
	</header>

	{#if $daoStore.globalError}
		<p class="global-error" role="alert">{$daoStore.globalError}</p>
	{/if}

	<section class="step-stack">
		<StepCard
			title="1) Wallet auth"
			description="Connect the artifact owner wallet and establish a user session."
			status={$authStore.step.status}
			error={$authStore.step.error}
		>
			<form
				class="stack-form"
				onsubmit={connectWallet}
				aria-busy={$authStore.step.status === 'loading'}
			>
				<label>
					<span>Owner wallet address</span>
					<input
						type="text"
						value={$authStore.walletAddress}
						oninput={(event) =>
							authStore.setWalletAddress((event.currentTarget as HTMLInputElement).value)}
						placeholder="0x1111111111111111111111111111111111111111"
						required
					/>
				</label>
				<button type="submit" disabled={$authStore.step.status === 'loading'}>
					Connect owner wallet
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
				onsubmit={publishArtifact}
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
			description="Join a DAO circle explicitly; no onboarding step assigns one automatically."
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
				<p class="detail">
					Joined: {$selectedDao?.name ?? $daoStore.membership.daoId}
					{#if $daoStore.membership.leftAt}
						(left at {$daoStore.membership.leftAt})
					{/if}
				</p>
			{/if}
		</StepCard>

		<StepCard
			title="4) Grant permission"
			description="Create an explicit consent record before a non-owner transfer is recorded."
			status={$permissionCardState.status}
			error={$permissionCardState.error}
		>
			<form
				class="stack-form"
				onsubmit={grantPermission}
				aria-busy={$permissionStore.grantStep.status === 'loading'}
			>
				<label>
					<span>Grantee wallet</span>
					<input
						type="text"
						value={$permissionStore.granteeWallet}
						oninput={(event) =>
							permissionStore.setGranteeWallet((event.currentTarget as HTMLInputElement).value)}
						placeholder="0x2222222222222222222222222222222222222222"
						required
					/>
				</label>
				<label>
					<span>Scope</span>
					<input
						type="text"
						value={$permissionStore.scope}
						oninput={(event) =>
							permissionStore.setScope((event.currentTarget as HTMLInputElement).value)}
						required
					/>
				</label>
				<div class="action-row">
					<button
						type="submit"
						disabled={!$canGrantPermission || $permissionStore.grantStep.status === 'loading'}
					>
						Grant permission
					</button>
					<button
						type="button"
						class="secondary-button"
						onclick={revokePermission}
						disabled={$permissionStore.permission?.status !== 'active' ||
							$permissionStore.revokeStep.status === 'loading'}
					>
						Revoke permission
					</button>
				</div>
			</form>
			{#if $permissionStore.permission}
				<p class="detail">
					Permission: {$permissionStore.permission.id} ({$permissionStore.permission.status})
				</p>
			{/if}
		</StepCard>

		<StepCard
			title="5) Record transfer"
			description="Authenticate the grantee wallet and record a permissioned transfer to the artifact owner."
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
					Record transfer
				</button>
			</form>
			{#if $paymentStore.tip}
				<p class="detail">Transfer confirmed: {$paymentStore.tip.paymentId}</p>
			{/if}
		</StepCard>

		<StepCard
			title="6) View attribution"
			description="Fetch attribution records to verify permission and value flow references."
			status={$paymentStore.attributionStep.status}
			error={$paymentStore.attributionStep.error}
		>
			<button type="button" onclick={viewAttribution} disabled={!$canViewAttribution}>
				Refresh attribution
			</button>

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

		<StepCard
			title="7) Export user data"
			description="Fetch owner-owned and owner-associated records in one machine-readable bundle."
			status={$exportStore.step.status}
			error={$exportStore.step.error}
		>
			<button
				type="button"
				onclick={exportUserData}
				disabled={!$canExport || $exportStore.step.status === 'loading'}
			>
				Export owner data
			</button>

			{#if $exportStore.bundle}
				<ul class="summary-list">
					<li>Artifacts: {$exportStore.bundle.artifacts.length}</li>
					<li>Memberships: {$exportStore.bundle.memberships.length}</li>
					<li>Permissions: {$exportStore.bundle.permissions.length}</li>
					<li>Transfers: {$exportStore.bundle.transfers.length}</li>
				</ul>
			{/if}
		</StepCard>

		<StepCard
			title="8) Leave DAO"
			description="Close the reversibility check by leaving the joined DAO."
			status={$daoStore.leaveStep.status}
			error={$daoStore.leaveStep.error}
		>
			<button
				type="button"
				onclick={leaveDao}
				disabled={!$canLeaveDao || $daoStore.leaveStep.status === 'loading'}
			>
				Leave DAO
			</button>
			{#if $daoStore.membership?.leftAt}
				<p class="detail">Left DAO at {$daoStore.membership.leftAt}</p>
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

	.action-row {
		display: flex;
		flex-wrap: wrap;
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

	.secondary-button {
		border-color: var(--color-border-strong);
		background: var(--color-surface);
		color: var(--color-text);
	}

	button:disabled {
		opacity: 0.55;
		cursor: not-allowed;
	}

	.detail {
		font-size: 0.88rem;
		color: var(--color-pill-neutral-text);
	}

	.summary-list {
		margin: 0;
		padding-left: 1.2rem;
		color: var(--color-pill-neutral-text);
		font-size: 0.88rem;
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

		.action-row {
			display: grid;
		}

		.attribution-table {
			display: block;
			overflow-x: auto;
		}
	}
</style>
