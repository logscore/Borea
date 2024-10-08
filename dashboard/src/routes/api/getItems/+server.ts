// src/routes/api/getItems/+server.ts
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ request, fetch }) => {
	const backendUrl = 'http://localhost:8080/getItems';

	try {
		// Forward the original request body to the backend
		const { query, params } = await request.json();

		const backendResponse = await fetch(backendUrl, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ query, params })
		});

		if (!backendResponse.ok) {
			throw new Error(`Backend responded with status: ${backendResponse.status}`);
		}

		const data = await backendResponse.json();
		return json(data);
	} catch (error) {
		console.error('Error proxying request to backend:', error);
		return new Response(JSON.stringify({ error: 'Internal Server Error' }), {
			status: 500,
			headers: { 'Content-Type': 'application/json' }
		});
	}
};
