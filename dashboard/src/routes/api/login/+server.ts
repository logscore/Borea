import { json } from '@sveltejs/kit';
import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken'; // Add this import
import { SERVER_KEY } from '$env/static/private';

//replace with SQLite storage implementation)

const passwordhsh = await bcrypt.hash('gabagool', 10);

const users = [
	{
		username: 'admin',
		passwordHash: `${passwordhsh}`
	} // You'll need to set this password hash
];
export async function POST({ request, cookies }) {
	try {
		const { username, password } = await request.json();
		const user = users.find((u) => u.username === username);

		if (user) {
			const isMatch = await bcrypt.compare(password, user.passwordHash);

			if (isMatch) {
				const authToken = jwt.sign({ username: user.username }, SERVER_KEY, { expiresIn: '7d' });
				cookies.set('session', authToken, {
					path: '/',
					httpOnly: true,
					sameSite: 'strict',
					maxAge: 60 * 60 * 24 * 7 // 7 day expiration
				});
				return json({ success: true });
			}
		}

		return json({ success: false }, { status: 401 });
	} catch (error) {
		return json({ success: false, error: 'An unexpected error occurred' }, { status: 500 });
	}
}
