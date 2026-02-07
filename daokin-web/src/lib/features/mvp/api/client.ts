import { requestJson, type HttpFetch } from './http';
import {
	connectWalletMock,
	createArtifactMock,
	joinDaoMock,
	listAttributionMock,
	listDaosMock,
	tipMock
} from './mock';
import type {
	ApiMode,
	Artifact,
	AttributionRecord,
	AuthSession,
	CreateArtifactInput,
	DaoMembership,
	DaoSummary,
	JoinDaoInput,
	TipInput,
	TipPayment
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

function joinPath(baseUrl: string, path: string): string {
	const cleanedBase = baseUrl.endsWith('/') ? baseUrl.slice(0, -1) : baseUrl;
	const cleanedPath = path.startsWith('/') ? path : `/${path}`;
	return `${cleanedBase}${cleanedPath}`;
}

function createHttpClient(baseUrl: string, httpFetch: HttpFetch): ApiClient {
	return {
		auth: {
			walletConnect: (walletAddress) =>
				requestJson<AuthSession>(httpFetch, joinPath(baseUrl, '/v1/auth/wallet'), {
					method: 'POST',
					body: JSON.stringify({ walletAddress })
				})
		},
		artifact: {
			create: (input) =>
				requestJson<Artifact>(httpFetch, joinPath(baseUrl, '/v1/artifacts'), {
					method: 'POST',
					body: JSON.stringify(input)
				})
		},
		dao: {
			list: () => requestJson<DaoSummary[]>(httpFetch, joinPath(baseUrl, '/v1/daos')),
			join: (input) =>
				requestJson<DaoMembership>(httpFetch, joinPath(baseUrl, `/v1/daos/${input.daoId}/join`), {
					method: 'POST',
					body: JSON.stringify(input)
				})
		},
		payment: {
			tip: (input) =>
				requestJson<TipPayment>(httpFetch, joinPath(baseUrl, '/v1/payments/tips'), {
					method: 'POST',
					body: JSON.stringify(input)
				})
		},
		attribution: {
			listByArtifact: (artifactId) =>
				requestJson<AttributionRecord[]>(
					httpFetch,
					joinPath(baseUrl, `/v1/artifacts/${artifactId}/attributions`)
				)
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
