<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';
	import { processDateData } from '../../dataProcessors';

	export let data;
	export let title = 'Chart Title';
	export let xAxisTitle = 'Time';
	export let yAxisTitle = 'Value';
	export let seriesName = 'Series';
	export let lineColor = '#4fee25';

	let chart: ApexCharts;
	let chartElement: HTMLDivElement;

	$: x_data = processDateData(data);

	onMount(async () => {
		if (browser) {
			const ApexCharts = (await import('apexcharts')).default;

			const options = {
				chart: {
					type: 'line',
					height: 450,
					background: '#1e1e1e',
					toolbar: {
						show: true,
						tools: {
							download: true,
							selection: false,
							zoom: false,
							zoomin: false,
							zoomout: false,
							pan: false,
							reset: false
						}
					}
				},
				noData: {
					text: 'No Data To Display'
				},
				grid: {
					show: true
				},
				series: [
					{
						name: seriesName,
						data: x_data
					}
				],
				xaxis: {
					categories: Array.from({ length: 14 }, (_, i) => {
						let date = new Date();
						date.setDate(date.getDate() - (13 - i)); // Adjust for 14-day range
						return date.toLocaleDateString('en-US', { month: '2-digit', day: '2-digit' });
					}),
					title: {
						text: xAxisTitle
					},
					labels: {
						rotate: -45,
						rotateAlways: true,
						hideOverlappingLabels: false
					}
				},
				yaxis: {
					show: true,
					title: {
						text: yAxisTitle
					},
					min: 0
				},
				fill: {
					opacity: 1
				},
				title: {
					text: title,
					align: 'left'
				},
				dataLabels: {
					enabled: false
				},
				stroke: {
					show: true,
					width: 1.5,
					colors: [lineColor]
				},
				tooltip: {
					enabled: true,
					y: {
						formatter: function (val: number) {
							return val + ' ' + yAxisTitle.toLowerCase();
						}
					}
				},
				theme: {
					mode: 'dark'
				}
			};

			chart = new ApexCharts(chartElement, options);
			chart.render();
		}
	});

	onDestroy(() => {
		if (browser && chart) {
			chart.destroy();
		}
	});
</script>

<div bind:this={chartElement}></div>

<style>
	div {
		width: 100%;
		height: 100%;
	}
</style>
