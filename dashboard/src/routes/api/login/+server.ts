import { json } from '@sveltejs/kit';
import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken'; // Add this import
import { SERVER_KEY } from '$env/static/private';
import sqlite3 from 'sqlite3';
import path from 'path';

const dbPath = path.resolve('../auth.db');

interface User {
	username: string;
	passwordHash: string;
}

function queryUser(username: string): Promise<User[]> {
	const sql = `SELECT * FROM users WHERE username = ?`;
	return new Promise((resolve, reject) => {
		const db = new sqlite3.Database(dbPath, sqlite3.OPEN_READONLY, (err) => {
			if (err) {
				console.error('Error opening database:', err.message);
				return reject(err);
			}
		});

		db.all(sql, [username], (err: Error, rows: User[]) => {
			if (err) {
				console.error(err.message);
				return reject(err);
			}
			db.close();
			resolve(rows);
		});
	});
}

export async function POST({ request, cookies }) {
	try {
		const { username, password } = await request.json();
		const users = await queryUser(username);
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
