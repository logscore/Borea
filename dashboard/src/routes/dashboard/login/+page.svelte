<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import InputUsername from '$lib/components/InputUsername.svelte';
	import InputPassword from '$lib/components/InputPassword.svelte';
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

<div class="flex h-screen w-screen items-center justify-center">
	<Card.Root class="w-full max-w-sm">
		<Card.Header>
			<Card.Title class="text-2xl">Login</Card.Title>
			<Card.Description>Enter your admin username and password.</Card.Description>
		</Card.Header>
		<form on:submit|preventDefault={handleSubmit}>
			<Card.Content class="grid gap-4">
				<div class="grid gap-2">
					<Label for="username">Email</Label>
					<InputUsername bind:value={username} />
				</div>
				<div class="grid gap-2">
					<Label for="password">Password</Label>
					<InputPassword bind:value={password} />
				</div>
				{#if error}
					<div class="text-red-500">
						<p class="pb-1 text-center">{error}</p>
					</div>
				{/if}
			</Card.Content>
			<Card.Footer>
				<Button type="submit" class="w-full">Sign in</Button>
			</Card.Footer>
		</form>
	</Card.Root>
</div>
