import { expect, type Page, test } from "@playwright/test";

const ADMIN_PATH = process.env.STATIC_SMOKE_ADMIN_PATH ?? "/manage";
const LEGACY_ADMIN_PATH = "/admin";

function trackAssetFailures(page: Page) {
	const failures: string[] = [];

	page.on("response", (response) => {
		const resourceType = response.request().resourceType();
		if (resourceType !== "script" && resourceType !== "stylesheet") {
			return;
		}
		if (response.status() >= 400) {
			failures.push(`${response.status()} ${new URL(response.url()).pathname}`);
		}
	});

	return failures;
}

async function expectAppBooted(page: Page) {
	await page.waitForFunction(() => {
		const root = document.querySelector("#root");
		return Boolean(
			root &&
				root.childElementCount > 0 &&
				(window as Window & { __ENV?: unknown }).__ENV,
		);
	});
}

async function collectScriptPaths(page: Page) {
	return page
		.locator("script[src]")
		.evaluateAll((nodes) =>
			nodes.map((node) => new URL((node as HTMLScriptElement).src).pathname),
		);
}

test("admin custom path serves runtime-rewritten assets and survives reload", async ({
	page,
}) => {
	const failures = trackAssetFailures(page);

	const response = await page.goto(`${ADMIN_PATH}/dashboard/workplace`);
	expect(response?.status()).toBe(200);

	await expectAppBooted(page);
	await expect(page).toHaveURL(
		new RegExp(`${ADMIN_PATH}/dashboard/workplace$`),
	);
	expect(
		await page.evaluate(
			() =>
				(window as Window & { __ENV?: Record<string, string> }).__ENV
					?.VITE_ADMIN_PATH,
		),
	).toBe(ADMIN_PATH);

	const scriptPaths = await collectScriptPaths(page);
	expect(
		scriptPaths.some((path) => path.startsWith(`${ADMIN_PATH}/assets/`)),
	).toBe(true);

	const reloadResponse = await page.reload();
	expect(reloadResponse?.status()).toBe(200);
	await expectAppBooted(page);
	expect(failures).toEqual([]);
});

test("legacy /admin routes redirect to the runtime admin path", async ({
	request,
}) => {
	const response = await request.get(
		`${LEGACY_ADMIN_PATH}/dashboard/workplace?from=legacy`,
		{
			maxRedirects: 0,
		},
	);

	expect(response.status()).toBe(308);
	expect(response.headers().location).toBe(
		`${ADMIN_PATH}/dashboard/workplace?from=legacy`,
	);
});

test("user auth route loads from embedded assets and survives reload", async ({
	page,
}) => {
	const failures = trackAssetFailures(page);

	const response = await page.goto("/auth");
	expect(response?.status()).toBe(200);

	await expectAppBooted(page);
	await expect(page).toHaveURL(/\/auth$/);
	expect(
		await page.evaluate(
			() =>
				(window as Window & { __ENV?: Record<string, string> }).__ENV
					?.VITE_SITE_URL,
		),
	).toContain("127.0.0.1");

	const scriptPaths = await collectScriptPaths(page);
	expect(scriptPaths.some((path) => path.startsWith("/assets/"))).toBe(true);

	const reloadResponse = await page.reload();
	expect(reloadResponse?.status()).toBe(200);
	await expectAppBooted(page);
	expect(failures).toEqual([]);
});

test("user dashboard route returns SPA HTML on direct entry and reload", async ({
	page,
}) => {
	const failures = trackAssetFailures(page);

	const response = await page.goto("/dashboard");
	expect(response?.status()).toBe(200);

	await expectAppBooted(page);
	await expect(page).toHaveURL(/\/dashboard$/);

	const reloadResponse = await page.reload();
	expect(reloadResponse?.status()).toBe(200);
	await expectAppBooted(page);
	expect(failures).toEqual([]);
});
