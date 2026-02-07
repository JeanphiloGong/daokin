<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { StepStatus } from '../types';

	let {
		title,
		description,
		status = 'idle',
		error = '',
		children
	}: {
		title: string;
		description: string;
		status?: StepStatus;
		error?: string;
		children?: Snippet;
	} = $props();

	const statusLabel: Record<StepStatus, string> = {
		idle: 'Not started',
		loading: 'Working',
		success: 'Done',
		error: 'Failed'
	};
</script>

<section class="step-card" data-status={status}>
	<header class="step-header">
		<div>
			<h2>{title}</h2>
			<p>{description}</p>
		</div>
		<span class="status-pill">{statusLabel[status]}</span>
	</header>

	<div class="step-body">
		{@render children?.()}
	</div>

	{#if error}
		<p class="step-error" role="alert">{error}</p>
	{/if}
</section>

<style>
	.step-card {
		border: 1px solid #d1d5db;
		border-radius: 0.75rem;
		padding: 1rem;
		background: #ffffff;
		box-shadow: 0 1px 2px rgba(15, 23, 42, 0.06);
	}

	.step-header {
		display: flex;
		justify-content: space-between;
		gap: 1rem;
		align-items: start;
		margin-bottom: 0.85rem;
	}

	h2 {
		font-size: 1rem;
		font-weight: 700;
		margin: 0;
	}

	p {
		margin: 0.2rem 0 0;
		font-size: 0.9rem;
		color: #475569;
	}

	.status-pill {
		font-size: 0.75rem;
		line-height: 1;
		padding: 0.35rem 0.55rem;
		border-radius: 999px;
		background: #e2e8f0;
		color: #334155;
		white-space: nowrap;
	}

	.step-body {
		display: grid;
		gap: 0.75rem;
	}

	.step-error {
		margin-top: 0.75rem;
		font-size: 0.85rem;
		padding: 0.65rem 0.75rem;
		border-radius: 0.5rem;
		background: #fee2e2;
		border: 1px solid #fecaca;
		color: #991b1b;
	}

	.step-card[data-status='loading'] .status-pill {
		background: #dbeafe;
		color: #1d4ed8;
	}

	.step-card[data-status='success'] .status-pill {
		background: #dcfce7;
		color: #15803d;
	}

	.step-card[data-status='error'] .status-pill {
		background: #fee2e2;
		color: #b91c1c;
	}
</style>
