<script>
	import { goto } from '$app/navigation';

	async function handleLogout() {
		await fetch('/api/logout', { method: 'POST' });
		goto('/dashboard/login');
	}
</script>

<h1>Welcome to your dashboard!</h1>
<button on:click={handleLogout}>Logout</button>
