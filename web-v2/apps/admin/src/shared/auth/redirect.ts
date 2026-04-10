export function normalizeRedirectTarget(redirect?: string | null) {
	if (typeof redirect !== 'string') {
		return undefined;
	}

	if (!redirect.startsWith('/') || redirect.startsWith('//')) {
		return undefined;
	}

	return redirect;
}

export function buildRedirectTarget(input: {
	pathname: string;
	searchStr?: string;
	hash?: string;
}) {
	const redirect = `${input.pathname}${input.searchStr ?? ''}${input.hash ?? ''}`;

	return normalizeRedirectTarget(redirect) ?? '/';
}
