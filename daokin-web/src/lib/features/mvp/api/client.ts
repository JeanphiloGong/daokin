import { requestJson, type HttpFetch } from './http';
import {
	connectWalletMock,
	createArtifactMock,
	joinDaoMock,
	listAttributionMock,
	listDaosMock,
	tipMock
} from './mock';
import {
	ApiError,
	type ApiMode,
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

export interface ApiClient {
	auth: {
		walletConnect: (walletAddress: string) => Promise<AuthSession>;
	};
	artifact: {
		create: (input: CreateArtifactInput) => Promise<Artifact>;
	};
	dao: {
		list: () => Promise<DaoSummary[]>;
		join: (input: JoinDaoInput) => Promise<DaoMembership>;
	};
	payment: {
		tip: (input: TipInput) => Promise<TipPayment>;
	};
	attribution: {
		listByArtifact: (artifactId: string) => Promise<AttributionRecord[]>;
	};
}

export interface ApiClientConfig {
	mode?: ApiMode;
	baseUrl?: string;
	httpFetch?: HttpFetch;
}

const DEFAULT_BASE_URL = '/api';

const defaultDaos: DaoSummary[] = [
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
];

const httpTipLedger = new Map<string, TipPayment[]>();

function joinPath(baseUrl: string, path: string): string {
	const cleanedBase = baseUrl.endsWith('/') ? baseUrl.slice(0, -1) : baseUrl;
	const cleanedPath = path.startsWith('/') ? path : `/${path}`;
	return `${cleanedBase}${cleanedPath}`;
}

function normalizeWallet(walletAddress: string): string {
	return walletAddress.trim().toLowerCase();
}

function makeDisplayName(walletAddress: string): string {
	const normalized = normalizeWallet(walletAddress);
	if (normalized.length < 10) {
		return 'Member';
	}
	return `Member ${normalized.slice(0, 6)}...${normalized.slice(-4)}`;
}

function nowISO(): string {
	return new Date().toISOString();
}

function makeLocalID(prefix: string): string {
	return `${prefix}_${Math.random().toString(36).slice(2, 10)}`;
}

function makePseudoHash(title: string, description: string): string {
	const source = `${title.trim()}|${description.trim()}`;
	let hash = 0;
	for (let i = 0; i < source.length; i += 1) {
		hash = (hash << 5) - hash + source.charCodeAt(i);
		hash |= 0;
	}
	const hex = Math.abs(hash).toString(16).padStart(64, '0');
	return `0x${hex}`;
}

function ensureWallet(userID: string): string {
	const normalized = normalizeWallet(userID);
	if (!normalized.startsWith('0x') || normalized.length < 10) {
		throw new ApiError('Wallet format invalid for current API contract.', 400, { userID });
	}
	return normalized;
}

function buildCreatorAttribution(artifactId: string, creatorWallet: string): AttributionRecord[] {
	const tips = httpTipLedger.get(artifactId) ?? [];
	const latestTip = tips[tips.length - 1] ?? null;

	const creator: AttributionRecord = {
		contributorId: creatorWallet,
		contributorName: makeDisplayName(creatorWallet),
		walletAddress: creatorWallet,
		role: 'creator',
		sharePercent: latestTip && latestTip.toUserId !== creatorWallet ? 85 : 100
	};

	if (!latestTip || latestTip.toUserId === creatorWallet) {
		return [creator];
	}

	const steward: AttributionRecord = {
		contributorId: latestTip.toUserId,
		contributorName: 'Dao steward',
		walletAddress: latestTip.toUserId,
		role: 'dao steward',
		sharePercent: 15
	};

	return [creator, steward];
}

function createHttpClient(baseUrl: string, httpFetch: HttpFetch): ApiClient {
	return {
		auth: {
			walletConnect: async (walletAddress) => {
				const wallet = ensureWallet(walletAddress);
				const challenge = await requestJson<{ nonce: string }>(
					httpFetch,
					joinPath(baseUrl, '/v1/auth/challenge'),
					{
						method: 'POST',
						body: JSON.stringify({ wallet })
					}
				);
				const session = await requestJson<{
					wallet?: string;
					access_token?: string;
					verified_at?: string;
				}>(httpFetch, joinPath(baseUrl, '/v1/auth/verify'), {
					method: 'POST',
					body: JSON.stringify({ wallet, nonce: challenge.nonce, signature: 'local-dev-signature' })
				});
				const verifiedWallet = normalizeWallet(session.wallet ?? wallet);
				return {
					token: session.access_token ?? makeLocalID('session'),
					userId: verifiedWallet,
					walletAddress: verifiedWallet,
					displayName: makeDisplayName(verifiedWallet)
				};
			}
		},
		artifact: {
			create: async (input) => {
				const creatorWallet = ensureWallet(input.creatorId);
				const payload = {
					creator_wallet: creatorWallet,
					content_hash: makePseudoHash(input.title, input.description),
					content_uri: `daokin://artifact/${encodeURIComponent(input.title.trim())}`
				};
				const created = await requestJson<{
					id: string;
					creator_wallet: string;
					created_at: string;
				}>(httpFetch, joinPath(baseUrl, '/v1/artifacts'), {
					method: 'POST',
					body: JSON.stringify(payload)
				});
				return {
					id: created.id,
					title: input.title,
					description: input.description,
					creatorId: normalizeWallet(created.creator_wallet),
					createdAt: created.created_at
				};
			}
		},
		dao: {
			list: async () => [...defaultDaos],
			join: async (input) => {
				const wallet = ensureWallet(input.userId);
				const joined = await requestJson<{ dao_id?: string; wallet?: string; joined_at?: string }>(
					httpFetch,
					joinPath(baseUrl, `/v1/daos/${input.daoId}/join`),
					{
						method: 'POST',
						body: JSON.stringify({ wallet })
					}
				);
				return {
					daoId: joined.dao_id ?? input.daoId,
					userId: normalizeWallet(joined.wallet ?? wallet),
					role: 'member',
					joinedAt: joined.joined_at ?? nowISO()
				};
			}
		},
		payment: {
			tip: async (input) => {
				const payment: TipPayment = {
					paymentId: makeLocalID('tip'),
					artifactId: input.artifactId,
					fromUserId: ensureWallet(input.fromUserId),
					toUserId: ensureWallet(input.toUserId),
					amount: input.amount,
					currency: 'USDC',
					status: 'confirmed',
					createdAt: nowISO()
				};
				const current = httpTipLedger.get(input.artifactId) ?? [];
				httpTipLedger.set(input.artifactId, [...current, payment]);
				return payment;
			}
		},
		attribution: {
			listByArtifact: async (artifactId) => {
				const artifact = await requestJson<{ creator_wallet?: string }>(
					httpFetch,
					joinPath(baseUrl, `/v1/artifacts/${artifactId}`)
				);
				const creatorWallet = ensureWallet(artifact.creator_wallet ?? '');
				return buildCreatorAttribution(artifactId, creatorWallet);
			}
		}
	};
}

function createMockClient(): ApiClient {
	return {
		auth: {
			walletConnect: connectWalletMock
		},
		artifact: {
			create: createArtifactMock
		},
		dao: {
			list: listDaosMock,
			join: joinDaoMock
		},
		payment: {
			tip: tipMock
		},
		attribution: {
			listByArtifact: listAttributionMock
		}
	};
}

export function createApiClient(config: ApiClientConfig = {}): ApiClient {
	const mode = config.mode ?? 'mock';

	if (mode === 'mock') {
		return createMockClient();
	}

	const httpFetch = config.httpFetch ?? fetch;
	const baseUrl = config.baseUrl ?? DEFAULT_BASE_URL;
	return createHttpClient(baseUrl, httpFetch);
}
