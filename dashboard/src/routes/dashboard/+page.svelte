<script lang="ts">
	import Nav from '$lib/components/Nav.svelte';
	export let data;
</script>

<div class="flex">
	<h1>Welcome to your dashboard!</h1>
	<Nav />
</div>

<div>{JSON.stringify(data.result)}</div>

<div>{JSON.stringify(data.error)}</div>
