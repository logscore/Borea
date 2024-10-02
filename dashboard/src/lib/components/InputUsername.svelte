<script lang="ts">
	import { onMount } from 'svelte';
	import { Input } from '$lib/components/ui/input/index.js';

	let ref: HTMLInputElement | null;
	let props = {
		type: 'text',
		name: 'username',
		placeholder: 'Username',
		autocomplete: 'username' as const,
		required: true
	};
	export let value = '';

	onMount(() => {
		ref?.focus();
	});
</script>

<Input bind:ref bind:value {...props} {...$$restProps} />
