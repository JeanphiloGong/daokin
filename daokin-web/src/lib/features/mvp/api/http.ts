import { ApiError } from '../types';

export type HttpFetch = typeof fetch;

interface ApiErrorEnvelope {
	code?: string;
	message?: string;
	request_id?: string;
	details?: unknown;
}

interface JsonResponse {
	message?: string;
	data?: unknown;
	error?: string | ApiErrorEnvelope;
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

	const requestId = response.headers.get('x-request-id') ?? undefined;
	let payload: JsonResponse | null = null;
	const rawText = await response.text();
	if (rawText.length > 0) {
		try {
			payload = JSON.parse(rawText) as JsonResponse;
		} catch {
			throw new ApiError(
				`Server returned invalid JSON: ${response.status}`,
				response.status,
				rawText,
				'INVALID_JSON',
				requestId
			);
		}
	}

	if (!response.ok) {
		const errorField = payload?.error;
		const errorObj = typeof errorField === 'object' && errorField !== null ? errorField : null;
		const message =
			errorObj?.message ??
			(typeof errorField === 'string' ? errorField : undefined) ??
			payload?.message ??
			`Request failed: ${response.status}`;
		const code = errorObj?.code;
		const details = errorObj?.details ?? payload;
		const errorRequestID = errorObj?.request_id ?? requestId;
		throw new ApiError(message, response.status, details, code, errorRequestID);
	}

	if (payload && 'data' in payload && payload.data !== undefined) {
		return payload.data as T;
	}

	return (payload ?? {}) as T;
}
