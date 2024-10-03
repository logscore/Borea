import { type RequestEvent, redirect, type Handle } from '@sveltejs/kit';
import { SERVER_KEY } from '$env/static/private';
import jwt from 'jsonwebtoken';

const nonPublicRoutes = ['/', '/dashboard'];

const authenticatedUser = (event: RequestEvent) => {
	const token = event.cookies.get('session');

	try {
		jwt.verify(token ?? '', SERVER_KEY);
		return true;
	} catch {
		return false;
	}
};

export const handle: Handle = async ({ event, resolve }) => {
	const verified = authenticatedUser(event);
	const { pathname } = event.url;

	if (!verified && nonPublicRoutes.includes(pathname)) {
		throw redirect(303, '/dashboard/login');
	}

	const response = await resolve(event);
	return response;
};
