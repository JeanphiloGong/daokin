import {
	ApiError,
	type Artifact,
	type AttributionRecord,
	type AuthSession,
	type CreateArtifactInput,
	type DaoMembership,
	type DaoSummary,
	type JoinDaoInput,
	type TipInput,
	type TipPayment
} from '../types';

interface MockDatabase {
	users: Record<string, AuthSession>;
	daos: DaoSummary[];
	artifacts: Artifact[];
	memberships: DaoMembership[];
	payments: TipPayment[];
	attributions: Record<string, AttributionRecord[]>;
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
	payments: [],
	attributions: {}
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

function ensureWallet(walletAddress: string): void {
	if (!/^0x[a-f0-9]{8,}$/i.test(walletAddress.trim())) {
		throw new ApiError('Wallet address format is invalid. Example: 0xabc12345');
	}
}

export async function connectWalletMock(walletAddress: string): Promise<AuthSession> {
	ensureWallet(walletAddress);
	const normalized = normalizeWallet(walletAddress);
	const existing = Object.values(db.users).find(
		(user) => normalizeWallet(user.walletAddress) === normalized
	);

	if (existing) {
		return delay(existing);
	}

	const session: AuthSession = {
		token: uid('jwt'),
		userId: uid('user'),
		walletAddress,
		displayName: `Member ${Object.keys(db.users).length + 1}`
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
	if (!db.users[input.creatorId]) {
		throw new ApiError('Wallet session not found. Please reconnect wallet.');
	}

	const artifact: Artifact = {
		id: uid('artifact'),
		title: input.title.trim(),
		description: input.description.trim(),
		creatorId: input.creatorId,
		createdAt: now()
	};

	db.artifacts.push(artifact);

	const creator = db.users[input.creatorId];
	db.attributions[artifact.id] = [
		{
			contributorId: creator.userId,
			contributorName: creator.displayName,
			walletAddress: creator.walletAddress,
			role: 'creator',
			sharePercent: 100
		}
	];

	return delay(artifact);
}

export async function joinDaoMock(input: JoinDaoInput): Promise<DaoMembership> {
	const dao = db.daos.find((item) => item.id === input.daoId);
	if (!dao) {
		throw new ApiError('DAO not found.');
	}
	if (!db.users[input.userId]) {
		throw new ApiError('User session expired. Reconnect wallet.');
	}

	const existing = db.memberships.find((m) => m.daoId === input.daoId && m.userId === input.userId);
	if (existing) {
		if (input.artifactId) {
			rewriteAttributionForDao(input.artifactId, dao.name);
		}
		return delay(existing);
	}

	const membership: DaoMembership = {
		daoId: input.daoId,
		userId: input.userId,
		role: 'member',
		joinedAt: now()
	};
	db.memberships.push(membership);

	if (input.artifactId) {
		rewriteAttributionForDao(input.artifactId, dao.name);
	}

	return delay(membership);
}

function rewriteAttributionForDao(artifactId: string, daoName: string): void {
	const rows = db.attributions[artifactId];
	if (!rows || rows.length === 0) {
		return;
	}
	const creator = rows[0];
	db.attributions[artifactId] = [
		{ ...creator, sharePercent: 85 },
		{
			contributorId: uid('steward'),
			contributorName: `${daoName} Steward`,
			walletAddress: '0xdao000001',
			role: 'dao steward',
			sharePercent: 15
		}
	];
}

export async function tipMock(input: TipInput): Promise<TipPayment> {
	if (!db.artifacts.some((artifact) => artifact.id === input.artifactId)) {
		throw new ApiError('Artifact not found for tip.');
	}
	if (!db.users[input.fromUserId]) {
		throw new ApiError('Tip sender is not authenticated.');
	}
	if (input.amount <= 0) {
		throw new ApiError('Tip amount must be greater than 0.');
	}

	const payment: TipPayment = {
		paymentId: uid('tip'),
		artifactId: input.artifactId,
		fromUserId: input.fromUserId,
		toUserId: input.toUserId,
		amount: input.amount,
		currency: 'USDC',
		status: 'confirmed',
		createdAt: now()
	};

	db.payments.push(payment);
	return delay(payment);
}

export async function listAttributionMock(artifactId: string): Promise<AttributionRecord[]> {
	if (!db.artifacts.some((artifact) => artifact.id === artifactId)) {
		throw new ApiError('Artifact does not exist.');
	}

	return delay([...(db.attributions[artifactId] ?? [])]);
}
