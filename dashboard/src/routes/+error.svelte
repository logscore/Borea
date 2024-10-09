<script lang="ts">
	import { page } from '$app/stores';

	$: console.error($page.error);
</script>

<div class="p-4">
	<div class="grid gap-4">
		<span class="text-2xl font-bold">Error: {$page.status}</span>

		{#if $page.error}
			<div>
				<h2 class="mb-1 text-lg font-bold">Message:</h2>
				<pre
					class="border-danger-300 bg-danger-50 text-danger-content ml-2 rounded-md border p-4 text-xs">{$page
						.error.message}</pre>
			</div>
		{/if}
	</div>
</div>
