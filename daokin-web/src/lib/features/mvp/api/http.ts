import { ApiError } from '../types';

export type HttpFetch = typeof fetch;

interface JsonResponse {
	message?: string;
	data?: unknown;
	[key: string]: unknown;
}

export async function requestJson<T>(
	httpFetch: HttpFetch,
	url: string,
	init?: RequestInit
): Promise<T> {
	const response = await httpFetch(url, {
		headers: {
			'content-type': 'application/json',
			...(init?.headers ?? {})
		},
		...init
	});

	let payload: JsonResponse | null = null;
	const rawText = await response.text();
	if (rawText.length > 0) {
		try {
			payload = JSON.parse(rawText) as JsonResponse;
		} catch {
			throw new ApiError('Server returned invalid JSON', response.status, rawText);
		}
	}

	if (!response.ok) {
		throw new ApiError(
			payload?.message ?? `Request failed: ${response.status}`,
			response.status,
			payload
		);
	}

	if (payload && 'data' in payload && payload.data !== undefined) {
		return payload.data as T;
	}

	return (payload ?? {}) as T;
}
