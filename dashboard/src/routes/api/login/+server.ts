import { json } from '@sveltejs/kit';
import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken';
import { SERVER_KEY, GO_PORT } from '$env/static/private';

type Auth_User = {
	id: number;
	username: string;
	password_hash: string;
};

// Note, the back end is built to take an array of params that will be inserted into the query you pass in.
// Note the '?' in the query below.
// That will be where your params are inserted in the order you put in the array.
async function postQueryToServer(handle: string, query: string, params: string[]) {
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

		// Using parameterized query to prevent SQL injection, even though the backend functions are safe.
		// TODO: implement a ORM to make this function work for any GET request and prevent sql injection
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
