export type StepStatus = 'idle' | 'loading' | 'success' | 'error';
export type StepKey = 'walletAuth' | 'createArtifact' | 'joinDao' | 'tipPay' | 'viewAttribution';

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
