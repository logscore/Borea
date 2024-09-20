import { type RequestEvent, redirect, type Handle } from '@sveltejs/kit';
import { SERVER_KEY } from '$env/static/private';
import jwt from 'jsonwebtoken';

const nonPublicRoutes = ['/', '/dashboard'];

const authenticatedUser = async (event: RequestEvent) => {
	const token = event.cookies.get('session');

	try {
		await jwt.verify(token ?? '', SERVER_KEY);
		return true;
	} catch {
		return false;
	}
};

export const handle: Handle = async ({ event, resolve }) => {
	const verified = await authenticatedUser(event);
	const { pathname } = event.url;

	if (nonPublicRoutes.includes(pathname) && !verified) {
		throw redirect(303, '/dashboard/login');
	}

	return await resolve(event);
};
