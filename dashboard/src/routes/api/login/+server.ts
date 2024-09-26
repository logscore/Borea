import { json } from '@sveltejs/kit';
import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken';
import { SERVER_KEY, GO_PORT } from '$env/static/private';

type AuthUser = {
	ID: number;
	username: string;
	passwordHash: string;
};

async function postQueryToServer(handle: string, query: string, params: (string | number)[]) {
	const url = `http://localhost:${GO_PORT}/${handle}`;
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

		// Using parameterized query to prevent SQL injection
		// TODO: implement a ORM to make this function wokr for any GET request and prevent sql injection
		const query = `SELECT * FROM auth_users WHERE username = ?`;

		const data = (await postQueryToServer('getItems', query, [username])) as AuthUser[];

		// Check if user exists
		const user = data.find((u) => u.username === username);

		if (user) {
			// Compare the provided password with the stored password hash
			const isMatch = await bcrypt.compare(password, user.passwordHash);

			if (isMatch) {
				// Generate authentication token
				const authToken = jwt.sign({ username: user.username }, SERVER_KEY, { expiresIn: '7d' });

				// Set cookie with the authentication token
				cookies.set('session', authToken, {
					path: '/',
					secure: true,
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
