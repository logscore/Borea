<script>
	import { goto } from '$app/navigation';

	let username = '';
	let password = '';
	let error = '';

	async function handleSubmit() {
		const response = await fetch('/api/login', {
			method: 'POST',
			body: JSON.stringify({ username, password }),
			headers: {
				'Content-Type': 'application/json'
			}
		});
		const result = await response.json();
		if (result.success) {
			goto('/dashboard');
		} else {
			error = 'Invalid credentials';
		}
	}
</script>

<form on:submit|preventDefault={handleSubmit}>
	<input bind:value={username} type="text" placeholder="Username" required />
	<input bind:value={password} type="password" placeholder="Password" required />
	<button type="submit">Login</button>
</form>
{#if error}
	<p>{error}</p>
{/if}
