import { requestJson, type HttpFetch } from './http';
import {
	connectWalletMock,
	createArtifactMock,
	getUserExportMock,
	grantPermissionMock,
	joinDaoMock,
	leaveDaoMock,
	listAttributionMock,
	listDaosMock,
	revokePermissionMock,
	tipMock
} from './mock';
import {
	ApiError,
	type ApiMode,
	type Artifact,
	type AttributionRecord,
	type AttributionTrail,
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
		leave: (input: LeaveDaoInput) => Promise<DaoMembership>;
	};
	permission: {
		grant: (input: GrantPermissionInput) => Promise<PermissionRecord>;
		revoke: (input: RevokePermissionInput) => Promise<PermissionRecord>;
	};
	payment: {
		tip: (input: TipInput) => Promise<TipPayment>;
	};
	attribution: {
		listByArtifact: (artifactId: string) => Promise<AttributionRecord[]>;
	};
	userExport: {
		get: (wallet: string) => Promise<UserExportBundle>;
	};
}

export interface ApiClientConfig {
	mode?: ApiMode;
	baseUrl?: string;
	httpFetch?: HttpFetch;
}

interface EthereumProvider {
	request: (args: { method: string; params?: unknown[] }) => Promise<unknown>;
}

interface BackendArtifact {
	id: string;
	creator_wallet: string;
	content_hash: string;
	content_uri: string;
	created_at: string;
}

interface BackendMembership {
	dao_id: string;
	wallet: string;
	joined_at: string;
	left_at?: string | null;
	leave_reason?: string;
}

interface BackendPermission {
	id: string;
	artifact_id: string;
	granter_wallet: string;
	grantee_wallet: string;
	scope: string;
	status: 'active' | 'revoked';
	granted_at: string;
	revoked_at?: string | null;
	revoke_reason?: string;
}

interface BackendTransfer {
	id: string;
	artifact_id: string;
	payer_wallet: string;
	recipient_wallet: string;
	token: string;
	amount_atomic: string;
	tx_hash: string;
	created_at: string;
}

interface BackendAttributionTrail {
	artifact_id: string;
	permissions: BackendPermission[];
	transfers: BackendTransfer[];
	total_records: number;
}

interface BackendUserExportBundle {
	wallet: string;
	artifacts: BackendArtifact[];
	memberships: BackendMembership[];
	permissions: BackendPermission[];
	transfers: BackendTransfer[];
}

const DEFAULT_BASE_URL = '/api';
const USDC_DECIMALS = 6;

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

function joinPath(baseUrl: string, path: string): string {
	const cleanedBase = baseUrl.endsWith('/') ? baseUrl.slice(0, -1) : baseUrl;
	const cleanedPath = path.startsWith('/') ? path : `/${path}`;
	return `${cleanedBase}${cleanedPath}`;
}

function normalizeWallet(walletAddress: string): string {
	return walletAddress.trim().toLowerCase();
}

function ensureWallet(walletAddress: string): string {
	const normalized = normalizeWallet(walletAddress);
	if (!/^0x[a-f0-9]{40}$/.test(normalized)) {
		throw new ApiError('Wallet format invalid for current API contract.', 400, { walletAddress });
	}
	return normalized;
}

function makeDisplayName(walletAddress: string): string {
	const normalized = normalizeWallet(walletAddress);
	return `Member ${normalized.slice(0, 6)}...${normalized.slice(-4)}`;
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

async function signChallengeMessage(wallet: string, message: string): Promise<string> {
	const provider = (globalThis as typeof globalThis & { ethereum?: EthereumProvider }).ethereum;
	if (!provider) {
		throw new ApiError('Browser wallet provider is required for HTTP API mode.', 400);
	}

	const signature = await provider.request({
		method: 'personal_sign',
		params: [message, wallet]
	});
	if (typeof signature !== 'string' || signature.trim() === '') {
		throw new ApiError('Wallet signature was not returned.', 400);
	}
	return signature;
}

function toUSDCAtomic(amount: number): string {
	if (!Number.isFinite(amount) || amount <= 0) {
		throw new ApiError('Tip amount must be greater than 0.', 400, { amount });
	}
	return Math.round(amount * 10 ** USDC_DECIMALS).toString();
}

function fromUSDCAtomic(amountAtomic: string): number {
	const parsed = Number.parseInt(amountAtomic, 10);
	if (!Number.isFinite(parsed)) {
		return 0;
	}
	return parsed / 10 ** USDC_DECIMALS;
}

function makeLocalTxHash(): string {
	const random = Math.random().toString(16).slice(2).padEnd(64, '0');
	return `0x${random.slice(0, 64)}`;
}

function mapMembership(membership: BackendMembership): DaoMembership {
	return {
		daoId: membership.dao_id,
		userId: normalizeWallet(membership.wallet),
		role: 'member',
		joinedAt: membership.joined_at,
		leftAt: membership.left_at ?? null,
		leaveReason: membership.leave_reason
	};
}

function mapPermission(permission: BackendPermission): PermissionRecord {
	return {
		id: permission.id,
		artifactId: permission.artifact_id,
		granterWallet: normalizeWallet(permission.granter_wallet),
		granteeWallet: normalizeWallet(permission.grantee_wallet),
		scope: permission.scope,
		status: permission.status,
		grantedAt: permission.granted_at,
		revokedAt: permission.revoked_at ?? null,
		revokeReason: permission.revoke_reason
	};
}

function mapTransfer(transfer: BackendTransfer): TransferRecord {
	return {
		id: transfer.id,
		artifactId: transfer.artifact_id,
		payerWallet: normalizeWallet(transfer.payer_wallet),
		recipientWallet: normalizeWallet(transfer.recipient_wallet),
		token: transfer.token,
		amountAtomic: transfer.amount_atomic,
		txHash: transfer.tx_hash,
		createdAt: transfer.created_at
	};
}

function mapArtifact(artifact: BackendArtifact, fallback?: CreateArtifactInput): Artifact {
	return {
		id: artifact.id,
		title: fallback?.title ?? artifact.content_uri,
		description: fallback?.description ?? artifact.content_hash,
		creatorId: normalizeWallet(artifact.creator_wallet),
		createdAt: artifact.created_at
	};
}

function mapUserExport(bundle: BackendUserExportBundle): UserExportBundle {
	return {
		wallet: normalizeWallet(bundle.wallet),
		artifacts: bundle.artifacts.map((artifact) => mapArtifact(artifact)),
		memberships: bundle.memberships.map(mapMembership),
		permissions: bundle.permissions.map(mapPermission),
		transfers: bundle.transfers.map(mapTransfer)
	};
}

function mapAttributionTrail(trail: BackendAttributionTrail): AttributionTrail {
	return {
		artifactId: trail.artifact_id,
		permissions: trail.permissions.map(mapPermission),
		transfers: trail.transfers.map(mapTransfer),
		totalRecords: trail.total_records
	};
}

function buildAttributionRecords(
	artifact: BackendArtifact,
	trail: AttributionTrail
): AttributionRecord[] {
	const creatorWallet = normalizeWallet(artifact.creator_wallet);
	const rows = new Map<string, AttributionRecord>();

	rows.set(`creator:${creatorWallet}`, {
		contributorId: creatorWallet,
		contributorName: makeDisplayName(creatorWallet),
		walletAddress: creatorWallet,
		role: 'creator',
		sharePercent: 100
	});

	for (const permission of trail.permissions) {
		if (permission.status !== 'active') {
			continue;
		}
		const wallet = normalizeWallet(permission.granteeWallet);
		if (wallet === creatorWallet) {
			continue;
		}
		rows.set(`permission:${permission.id}`, {
			contributorId: wallet,
			contributorName: makeDisplayName(wallet),
			walletAddress: wallet,
			role: `authorized ${permission.scope}`,
			sharePercent: 0
		});
	}

	for (const transfer of trail.transfers) {
		const wallet = normalizeWallet(transfer.payerWallet);
		if (wallet === creatorWallet || rows.has(`payer:${wallet}`)) {
			continue;
		}
		rows.set(`payer:${wallet}`, {
			contributorId: wallet,
			contributorName: makeDisplayName(wallet),
			walletAddress: wallet,
			role: 'payer',
			sharePercent: 0
		});
	}

	return Array.from(rows.values());
}

function createHttpClient(baseUrl: string, httpFetch: HttpFetch): ApiClient {
	const sessionsByWallet = new Map<string, AuthSession>();

	function rememberSession(session: AuthSession): AuthSession {
		sessionsByWallet.set(normalizeWallet(session.walletAddress), session);
		return session;
	}

	function getSession(walletAddress: string): AuthSession {
		const wallet = ensureWallet(walletAddress);
		const session = sessionsByWallet.get(wallet);
		if (!session) {
			throw new ApiError(`Wallet ${wallet} must be authenticated before this action.`, 401);
		}
		return session;
	}

	function authorizedJson<T>(path: string, walletAddress: string, init?: RequestInit): Promise<T> {
		const session = getSession(walletAddress);
		return requestJson<T>(httpFetch, joinPath(baseUrl, path), {
			...init,
			headers: {
				...(init?.headers ?? {}),
				authorization: `Bearer ${session.token}`
			}
		});
	}

	return {
		auth: {
			walletConnect: async (walletAddress) => {
				const wallet = ensureWallet(walletAddress);
				const challenge = await requestJson<{ nonce: string; message: string }>(
					httpFetch,
					joinPath(baseUrl, '/v1/auth/challenge'),
					{
						method: 'POST',
						body: JSON.stringify({ wallet })
					}
				);
				const signature = await signChallengeMessage(wallet, challenge.message);
				const session = await requestJson<{
					wallet?: string;
					access_token?: string;
					verified_at?: string;
				}>(httpFetch, joinPath(baseUrl, '/v1/auth/verify'), {
					method: 'POST',
					body: JSON.stringify({ wallet, nonce: challenge.nonce, signature })
				});
				const verifiedWallet = ensureWallet(session.wallet ?? wallet);
				return rememberSession({
					token: session.access_token ?? makeLocalID('session'),
					userId: verifiedWallet,
					walletAddress: verifiedWallet,
					displayName: makeDisplayName(verifiedWallet)
				});
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
				const created = await authorizedJson<BackendArtifact>('/v1/artifacts', creatorWallet, {
					method: 'POST',
					body: JSON.stringify(payload)
				});
				return mapArtifact(created, input);
			}
		},
		dao: {
			list: async () => [...defaultDaos],
			join: async (input) => {
				const wallet = ensureWallet(input.userId);
				const joined = await authorizedJson<BackendMembership>(
					`/v1/daos/${input.daoId}/join`,
					wallet,
					{
						method: 'POST',
						body: JSON.stringify({ wallet })
					}
				);
				return mapMembership(joined);
			},
			leave: async (input) => {
				const wallet = ensureWallet(input.userId);
				const left = await authorizedJson<BackendMembership>(
					`/v1/daos/${input.daoId}/leave`,
					wallet,
					{
						method: 'POST',
						body: JSON.stringify({ wallet, reason: input.reason ?? '' })
					}
				);
				return mapMembership(left);
			}
		},
		permission: {
			grant: async (input) => {
				const granterWallet = ensureWallet(input.granterWallet);
				const granteeWallet = ensureWallet(input.granteeWallet);
				const granted = await authorizedJson<BackendPermission>('/v1/permissions', granterWallet, {
					method: 'POST',
					body: JSON.stringify({
						artifact_id: input.artifactId,
						granter_wallet: granterWallet,
						grantee_wallet: granteeWallet,
						scope: input.scope
					})
				});
				return mapPermission(granted);
			},
			revoke: async (input) => {
				const granterWallet = ensureWallet(input.granterWallet);
				const revoked = await authorizedJson<BackendPermission>(
					`/v1/permissions/${input.permissionId}/revoke`,
					granterWallet,
					{
						method: 'POST',
						body: JSON.stringify({
							granter_wallet: granterWallet,
							reason: input.reason ?? ''
						})
					}
				);
				return mapPermission(revoked);
			}
		},
		payment: {
			tip: async (input) => {
				const payerWallet = ensureWallet(input.fromUserId);
				const transfer = await authorizedJson<BackendTransfer>('/v1/transfers', payerWallet, {
					method: 'POST',
					body: JSON.stringify({
						artifact_id: input.artifactId,
						payer_wallet: payerWallet,
						token: 'USDC',
						amount_atomic: toUSDCAtomic(input.amount),
						tx_hash: makeLocalTxHash()
					})
				});
				return {
					paymentId: transfer.id,
					artifactId: transfer.artifact_id,
					fromUserId: normalizeWallet(transfer.payer_wallet),
					toUserId: normalizeWallet(transfer.recipient_wallet),
					amount: fromUSDCAtomic(transfer.amount_atomic),
					currency: 'USDC',
					status: 'confirmed',
					createdAt: transfer.created_at
				};
			}
		},
		attribution: {
			listByArtifact: async (artifactId) => {
				const [artifact, trail] = await Promise.all([
					requestJson<BackendArtifact>(httpFetch, joinPath(baseUrl, `/v1/artifacts/${artifactId}`)),
					requestJson<BackendAttributionTrail>(
						httpFetch,
						joinPath(baseUrl, `/v1/artifacts/${artifactId}/attribution`)
					)
				]);
				return buildAttributionRecords(artifact, mapAttributionTrail(trail));
			}
		},
		userExport: {
			get: async (walletAddress) => {
				const wallet = ensureWallet(walletAddress);
				const bundle = await authorizedJson<BackendUserExportBundle>(
					`/v1/users/${wallet}/export`,
					wallet
				);
				return mapUserExport(bundle);
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
			join: joinDaoMock,
			leave: leaveDaoMock
		},
		permission: {
			grant: grantPermissionMock,
			revoke: revokePermissionMock
		},
		payment: {
			tip: tipMock
		},
		attribution: {
			listByArtifact: listAttributionMock
		},
		userExport: {
			get: getUserExportMock
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
