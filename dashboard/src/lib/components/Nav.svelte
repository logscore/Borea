<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { toggleMode } from 'mode-watcher';
	import Moon from 'lucide-svelte/icons/moon';
	import Sun from 'lucide-svelte/icons/sun';
	import BookIcon from '$lib/components/icons/BookIcon.svelte';
	import PersonIcon from '$lib/components/icons/PersonIcon.svelte';

	async function handleLogout() {
		await fetch('/api/logout', { method: 'POST' });
		goto('/dashboard/login');
	}
</script>

<nav
	class="fixed right-0 top-0 flex h-full flex-col justify-between bg-[#fef5ea] p-2 pb-3 dark:bg-gray-900"
>
	<Tooltip.Root>
		<Tooltip.Trigger>
			<Button
				href="https://receiptsniffer.com/docs"
				variant="outline"
				class="flex w-auto items-center justify-center rounded-full bg-[#fef5ea] p-2 hover:bg-[#e6d7c4] dark:bg-gray-900 dark:hover:bg-gray-800"
			>
				<PersonIcon />
			</Button>
		</Tooltip.Trigger>
		<Tooltip.Content>
			<p class="text-gray-800 dark:text-gray-200">Visit docs</p>
		</Tooltip.Content>
	</Tooltip.Root>
	<Button
		on:click={toggleMode}
		variant="outline"
		class="mt-3 flex w-auto items-center justify-center rounded-full bg-[#fef5ea] p-2 hover:bg-[#e6d7c4] dark:bg-gray-900 dark:hover:bg-gray-800"
	>
		<div class="relative flex h-7 w-7 items-center justify-center">
			<Sun
				class="absolute h-6 w-6 rotate-0 scale-100 text-gray-800 transition-all dark:-rotate-90 dark:scale-0 dark:text-gray-300"
			/>
			<Moon
				class="absolute h-6 w-6 rotate-90 scale-0 text-gray-800 transition-all dark:rotate-0 dark:scale-100 dark:text-gray-300"
			/>
		</div>
		<span class="sr-only">Toggle theme</span>
	</Button>
	<div class="flex h-full flex-col justify-end">
		<Tooltip.Root>
			<Tooltip.Trigger>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger asChild let:builder>
						<Button
							builders={[builder]}
							variant="outline"
							class="flex items-center justify-center rounded-full bg-[#fef5ea] p-2 hover:bg-[#e6d7c4] dark:bg-gray-900 dark:hover:bg-gray-800"
						>
							<BookIcon />
						</Button>
					</DropdownMenu.Trigger>
					<DropdownMenu.Content
						class="w-56 bg-[#f5e6d3] text-gray-800 dark:bg-gray-800 dark:text-gray-200"
					>
						<DropdownMenu.Label>My Account</DropdownMenu.Label>
						<DropdownMenu.Separator />
						<DropdownMenu.Item
							href="https://github.com/logscore/Borea"
							class="hover:bg-[#e6d7c4] dark:hover:bg-gray-700">GitHub</DropdownMenu.Item
						>
						<DropdownMenu.Item
							on:click={handleLogout}
							class="text-red-600 hover:bg-[#e6d7c4] dark:text-red-400 dark:hover:bg-gray-700"
							>Log out</DropdownMenu.Item
						>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</Tooltip.Trigger>
			<Tooltip.Content>
				<p class="text-gray-800 dark:text-gray-200">User account</p>
			</Tooltip.Content>
		</Tooltip.Root>
	</div>
</nav>
