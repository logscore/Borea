<script lang="ts">
	import Icon from 'svelte-fa';
	import { faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons';
	import { Input } from '$lib/components/ui/input/index.js';

	let ref: HTMLInputElement | null;
	let showPassword = false;
	let props = {
		type: 'password',
		name: 'password',
		placeholder: 'Password',
		autocomplete: 'current-password' as const,
		required: true
	};
	export let value = '';

	function toggleShowPassword() {
		showPassword = !showPassword;
		const type = showPassword ? 'text' : 'password';
		props.type = type;
	}
</script>

<div class="password-input">
	<Input bind:ref bind:value {...props} {...$$restProps} />
	<button type="button" tabindex="-1" on:click={toggleShowPassword}>
		<div>
			<Icon class="eye-icon" icon={showPassword ? faEyeSlash : faEye} />
		</div>
	</button>
</div>

<style>
	.password-input {
		position: relative;
	}

	.password-input button {
		color: gray;
		position: absolute;
		right: 0;
		top: 0;
		bottom: 0;
		background: none;
		border: none;
		cursor: pointer;

		& div {
			display: flex;
			padding: 0.5rem;
		}
	}
</style>
