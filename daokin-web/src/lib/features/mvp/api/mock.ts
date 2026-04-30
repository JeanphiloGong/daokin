import {
	ApiError,
	type Artifact,
	type AttributionRecord,
	type AuthSession,
	type CreateArtifactInput,
	type DaoMembership,
	type DaoSummary,
	type GrantPermissionInput,
	type JoinDaoInput,
	type LeaveDaoInput,
	type PermissionRecord,
	type RevokePermissionInput,
	type TipInput,
	type TipPayment,
	type TransferRecord,
	type UserExportBundle
} from '../types';

interface MockDatabase {
	users: Record<string, AuthSession>;
	daos: DaoSummary[];
	artifacts: Artifact[];
	memberships: DaoMembership[];
	permissions: PermissionRecord[];
	payments: TipPayment[];
	transfers: TransferRecord[];
}

const db: MockDatabase = {
	users: {},
	daos: [
		{
			id: 'dao-builders',
			name: 'The Dao of Builders',
			description: 'Ship practical tools that preserve contributor ownership.'
		},
		{
			id: 'dao-stewards',
			name: 'The Dao of Stewards',
			description: 'Hold governance lines and keep community decisions transparent.'
		}
	],
	artifacts: [],
	memberships: [],
	permissions: [],
	payments: [],
	transfers: []
};

function uid(prefix: string): string {
	const random = Math.random().toString(36).slice(2, 8);
	return `${prefix}_${random}`;
}

function now(): string {
	return new Date().toISOString();
}

function delay<T>(data: T): Promise<T> {
	return new Promise((resolve) => setTimeout(() => resolve(data), 180));
}

function normalizeWallet(walletAddress: string): string {
	return walletAddress.trim().toLowerCase();
}

function ensureWallet(walletAddress: string): string {
	const normalized = normalizeWallet(walletAddress);
	if (!/^0x[a-f0-9]{40}$/i.test(normalized)) {
		throw new ApiError('Wallet address format is invalid. Use a full EVM address.');
	}
	return normalized;
}

function makeDisplayName(walletAddress: string): string {
	const normalized = normalizeWallet(walletAddress);
	return `Member ${normalized.slice(0, 6)}...${normalized.slice(-4)}`;
}

function amountToAtomic(amount: number): string {
	return Math.round(amount * 1_000_000).toString();
}

export async function connectWalletMock(walletAddress: string): Promise<AuthSession> {
	const normalized = ensureWallet(walletAddress);
	const existing = db.users[normalized];

	if (existing) {
		return delay(existing);
	}

	const session: AuthSession = {
		token: uid('jwt'),
		userId: normalized,
		walletAddress: normalized,
		displayName: makeDisplayName(normalized)
	};

	db.users[session.userId] = session;
	return delay(session);
}

export async function listDaosMock(): Promise<DaoSummary[]> {
	return delay([...db.daos]);
}

export async function createArtifactMock(input: CreateArtifactInput): Promise<Artifact> {
	if (input.title.trim().length < 3) {
		throw new ApiError('Artifact title needs at least 3 characters.');
	}
	const creatorId = ensureWallet(input.creatorId);
	if (!db.users[creatorId]) {
		throw new ApiError('Wallet session not found. Please reconnect wallet.');
	}

	const artifact: Artifact = {
		id: uid('artifact'),
		title: input.title.trim(),
		description: input.description.trim(),
		creatorId,
		createdAt: now()
	};

	db.artifacts.push(artifact);
	return delay(artifact);
}

export async function joinDaoMock(input: JoinDaoInput): Promise<DaoMembership> {
	const dao = db.daos.find((item) => item.id === input.daoId);
	if (!dao) {
		throw new ApiError('DAO not found.');
	}
	const userId = ensureWallet(input.userId);
	if (!db.users[userId]) {
		throw new ApiError('User session expired. Reconnect wallet.');
	}

	const existing = db.memberships.find((m) => m.daoId === input.daoId && m.userId === userId);
	if (existing) {
		if (existing.leftAt) {
			existing.leftAt = null;
			existing.leaveReason = '';
			existing.joinedAt = now();
		}
		return delay(existing);
	}

	const membership: DaoMembership = {
		daoId: input.daoId,
		userId,
		role: 'member',
		joinedAt: now(),
		leftAt: null
	};
	db.memberships.push(membership);

	return delay(membership);
}

export async function leaveDaoMock(input: LeaveDaoInput): Promise<DaoMembership> {
	const userId = ensureWallet(input.userId);
	const membership = db.memberships.find(
		(item) => item.daoId === input.daoId && item.userId === userId && !item.leftAt
	);
	if (!membership) {
		throw new ApiError('Active DAO membership not found.');
	}

	membership.leftAt = now();
	membership.leaveReason = input.reason?.trim() ?? '';
	return delay(membership);
}

export async function grantPermissionMock(input: GrantPermissionInput): Promise<PermissionRecord> {
	const granterWallet = ensureWallet(input.granterWallet);
	const granteeWallet = ensureWallet(input.granteeWallet);
	const artifact = db.artifacts.find((item) => item.id === input.artifactId);
	if (!artifact) {
		throw new ApiError('Artifact not found for permission grant.');
	}
	if (normalizeWallet(artifact.creatorId) !== granterWallet) {
		throw new ApiError('Only the artifact creator can grant permission.');
	}
	if (granterWallet === granteeWallet) {
		throw new ApiError('Permission grantee must be different from the owner.');
	}

	const existing = db.permissions.find(
		(item) =>
			item.artifactId === input.artifactId &&
			item.granteeWallet === granteeWallet &&
			item.scope === input.scope &&
			item.status === 'active'
	);
	if (existing) {
		return delay(existing);
	}

	const permission: PermissionRecord = {
		id: uid('perm'),
		artifactId: input.artifactId,
		granterWallet,
		granteeWallet,
		scope: input.scope.trim(),
		status: 'active',
		grantedAt: now(),
		revokedAt: null
	};
	db.permissions.push(permission);
	return delay(permission);
}

export async function revokePermissionMock(
	input: RevokePermissionInput
): Promise<PermissionRecord> {
	const granterWallet = ensureWallet(input.granterWallet);
	const permission = db.permissions.find((item) => item.id === input.permissionId);
	if (!permission) {
		throw new ApiError('Permission not found.');
	}
	if (permission.granterWallet !== granterWallet) {
		throw new ApiError('Only the granter can revoke permission.');
	}
	if (permission.status !== 'active') {
		throw new ApiError('Permission is already revoked.');
	}

	permission.status = 'revoked';
	permission.revokedAt = now();
	permission.revokeReason = input.reason?.trim() ?? '';
	return delay(permission);
}

export async function tipMock(input: TipInput): Promise<TipPayment> {
	const artifact = db.artifacts.find((item) => item.id === input.artifactId);
	if (!artifact) {
		throw new ApiError('Artifact not found for tip.');
	}
	const fromUserId = ensureWallet(input.fromUserId);
	const toUserId = ensureWallet(input.toUserId);
	if (!db.users[fromUserId]) {
		throw new ApiError('Tip sender is not authenticated.');
	}
	if (input.amount <= 0) {
		throw new ApiError('Tip amount must be greater than 0.');
	}
	if (normalizeWallet(artifact.creatorId) !== fromUserId) {
		const hasPermission = db.permissions.some(
			(item) =>
				item.artifactId === input.artifactId &&
				item.granteeWallet === fromUserId &&
				item.scope === 'view' &&
				item.status === 'active'
		);
		if (!hasPermission) {
			throw new ApiError('Permission is required before this transfer.');
		}
	}

	const payment: TipPayment = {
		paymentId: uid('tip'),
		artifactId: input.artifactId,
		fromUserId,
		toUserId,
		amount: input.amount,
		currency: 'USDC',
		status: 'confirmed',
		createdAt: now()
	};
	const transfer: TransferRecord = {
		id: payment.paymentId,
		artifactId: input.artifactId,
		payerWallet: fromUserId,
		recipientWallet: toUserId,
		token: 'USDC',
		amountAtomic: amountToAtomic(input.amount),
		txHash: `0x${uid('mocktx').padEnd(64, '0').slice(0, 64)}`,
		createdAt: payment.createdAt
	};

	db.payments.push(payment);
	db.transfers.push(transfer);
	return delay(payment);
}

export async function listAttributionMock(artifactId: string): Promise<AttributionRecord[]> {
	const artifact = db.artifacts.find((item) => item.id === artifactId);
	if (!artifact) {
		throw new ApiError('Artifact does not exist.');
	}

	const rows: AttributionRecord[] = [
		{
			contributorId: artifact.creatorId,
			contributorName: makeDisplayName(artifact.creatorId),
			walletAddress: artifact.creatorId,
			role: 'creator',
			sharePercent: 100
		}
	];

	for (const permission of db.permissions) {
		if (permission.artifactId !== artifactId || permission.status !== 'active') {
			continue;
		}
		rows.push({
			contributorId: permission.granteeWallet,
			contributorName: makeDisplayName(permission.granteeWallet),
			walletAddress: permission.granteeWallet,
			role: `authorized ${permission.scope}`,
			sharePercent: 0
		});
	}

	return delay(rows);
}

export async function getUserExportMock(walletAddress: string): Promise<UserExportBundle> {
	const wallet = ensureWallet(walletAddress);
	if (!db.users[wallet]) {
		throw new ApiError('Wallet session not found. Please reconnect wallet.');
	}

	return delay({
		wallet,
		artifacts: db.artifacts.filter((artifact) => normalizeWallet(artifact.creatorId) === wallet),
		memberships: db.memberships.filter((membership) => membership.userId === wallet),
		permissions: db.permissions.filter(
			(permission) => permission.granterWallet === wallet || permission.granteeWallet === wallet
		),
		transfers: db.transfers.filter(
			(transfer) => transfer.payerWallet === wallet || transfer.recipientWallet === wallet
		)
	});
}
