import { json } from '@sveltejs/kit';
import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken';
const HOST_ADDRESS = process.env.HOST_ADDRESS;
const GO_PORT = process.env.GO_PORT;
const SERVER_KEY = process.env.SERVER_KEY;

type Auth_User = {
	id: number;
	username: string;
	password_hash: string;
};

async function postQueryToServer(handle: string, query: string, params: string[]) {
	const url = `http://${HOST_ADDRESS}:${GO_PORT}/${handle}`;
	try {
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ query, params })
		});

		if (!response.ok) {
			throw new Error(`HTTP Error: ${response.status}`);
		}

		const result = await response.json();
		return result;
	} catch (error) {
		console.error('Error:', error);
		throw error;
	}
}

export async function POST({ request, cookies }) {
	try {
		const { username, password } = await request.json();

		const query = `SELECT ID, username, password_hash FROM admin_users WHERE username = $1`;

		const data = (await postQueryToServer('getItem', query, [username])) as Auth_User;

		// Check if user exists
		const user = data.username === username;

		if (user) {
			// Compare the provided password with the stored password hash
			const isMatch = await bcrypt.compare(password, data.password_hash);

			if (isMatch) {
				// Generate authentication token
				const authToken = jwt.sign({ username: data.username }, SERVER_KEY, { expiresIn: '7d' });

				// Set cookie with the authentication token
				cookies.set('session', authToken, {
					path: '/',
					secure: false,
					httpOnly: true,
					sameSite: 'strict',
					maxAge: 60 * 60 * 24 * 7
				});

				return json({ success: true });
			}
		}

		return json({ success: false }, { status: 401 });
	} catch (error) {
		console.error('Error during authentication:', error);
		return json({ success: false, error: 'An unexpected error occurred' }, { status: 500 });
	}
}
