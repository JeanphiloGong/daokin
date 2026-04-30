export type StepStatus = 'idle' | 'loading' | 'success' | 'error';
export type StepKey =
	| 'walletAuth'
	| 'createArtifact'
	| 'joinDao'
	| 'grantPermission'
	| 'tipPay'
	| 'viewAttribution'
	| 'exportUserData'
	| 'leaveDao';

export interface StepState {
	status: StepStatus;
	error: string;
}

export interface DaoSummary {
	id: string;
	name: string;
	description: string;
}

export interface AuthSession {
	token: string;
	userId: string;
	walletAddress: string;
	displayName: string;
}

export interface Artifact {
	id: string;
	title: string;
	description: string;
	creatorId: string;
	createdAt: string;
}

export interface DaoMembership {
	daoId: string;
	userId: string;
	role: 'member' | 'steward';
	joinedAt: string;
	leftAt?: string | null;
	leaveReason?: string;
}

export interface PermissionRecord {
	id: string;
	artifactId: string;
	granterWallet: string;
	granteeWallet: string;
	scope: string;
	status: 'active' | 'revoked';
	grantedAt: string;
	revokedAt?: string | null;
	revokeReason?: string;
}

export interface TransferRecord {
	id: string;
	artifactId: string;
	payerWallet: string;
	recipientWallet: string;
	token: string;
	amountAtomic: string;
	txHash: string;
	createdAt: string;
}

export interface TipPayment {
	paymentId: string;
	artifactId: string;
	fromUserId: string;
	toUserId: string;
	amount: number;
	currency: 'USDC';
	status: 'confirmed';
	createdAt: string;
}

export interface AttributionRecord {
	contributorId: string;
	contributorName: string;
	walletAddress: string;
	role: string;
	sharePercent: number;
}

export interface AttributionTrail {
	artifactId: string;
	permissions: PermissionRecord[];
	transfers: TransferRecord[];
	totalRecords: number;
}

export interface UserExportBundle {
	wallet: string;
	artifacts: Artifact[];
	memberships: DaoMembership[];
	permissions: PermissionRecord[];
	transfers: TransferRecord[];
}

export interface ConnectWalletInput {
	walletAddress: string;
}

export interface CreateArtifactInput {
	title: string;
	description: string;
	creatorId: string;
}

export interface JoinDaoInput {
	daoId: string;
	userId: string;
	artifactId?: string;
}

export interface LeaveDaoInput {
	daoId: string;
	userId: string;
	reason?: string;
}

export interface GrantPermissionInput {
	artifactId: string;
	granterWallet: string;
	granteeWallet: string;
	scope: string;
}

export interface RevokePermissionInput {
	permissionId: string;
	granterWallet: string;
	reason?: string;
}

export interface TipInput {
	artifactId: string;
	fromUserId: string;
	toUserId: string;
	amount: number;
}

export type ApiMode = 'mock' | 'http';

export class ApiError extends Error {
	status?: number;
	details?: unknown;
	code?: string;
	requestId?: string;

	constructor(
		message: string,
		status?: number,
		details?: unknown,
		code?: string,
		requestId?: string
	) {
		super(message);
		this.name = 'ApiError';
		this.status = status;
		this.details = details;
		this.code = code;
		this.requestId = requestId;
	}
}
