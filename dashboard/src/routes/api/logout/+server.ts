import { json } from '@sveltejs/kit';

export function POST({ cookies }) {
	cookies.delete('session', { path: '/' });
	return json({ success: true });
}
