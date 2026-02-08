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

<section class="step-card" data-status={status} aria-busy={status === 'loading'}>
	<header class="step-header">
		<div>
			<h2>{title}</h2>
			<p>{description}</p>
		</div>
		<span class="status-pill" role="status" aria-live="polite">{statusLabel[status]}</span>
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
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: var(--space-4);
		background: var(--color-surface);
		box-shadow: 0 1px 2px rgba(15, 23, 42, 0.06);
	}

	.step-header {
		display: flex;
		justify-content: space-between;
		gap: var(--space-4);
		align-items: start;
		margin-bottom: var(--space-3);
	}

	h2 {
		font-size: 1rem;
		font-weight: 700;
		margin: 0;
	}

	p {
		margin: 0.2rem 0 0;
		font-size: 0.9rem;
		color: var(--color-text-muted);
	}

	.status-pill {
		font-size: 0.75rem;
		line-height: 1;
		padding: 0.35rem 0.55rem;
		border-radius: 999px;
		background: var(--color-pill-neutral-bg);
		color: var(--color-pill-neutral-text);
		white-space: nowrap;
	}

	.step-body {
		display: grid;
		gap: var(--space-3);
	}

	.step-error {
		margin-top: var(--space-3);
		font-size: 0.85rem;
		padding: 0.65rem 0.75rem;
		border-radius: var(--radius-sm);
		background: var(--color-error-bg);
		border: 1px solid var(--color-error-border);
		color: var(--color-error-text);
	}

	.step-card[data-status='loading'] .status-pill {
		background: var(--color-pill-loading-bg);
		color: var(--color-pill-loading-text);
	}

	.step-card[data-status='success'] .status-pill {
		background: var(--color-pill-success-bg);
		color: var(--color-pill-success-text);
	}

	.step-card[data-status='error'] .status-pill {
		background: var(--color-pill-error-bg);
		color: var(--color-pill-error-text);
	}
</style>
