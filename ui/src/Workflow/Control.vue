<template>
<b-container fluid>
	<b-row>
		<b-col><b-button @click="retry()" variant="danger">Retry</b-button>
		</b-col>
	</b-row>
</b-container>
</template>

<script>
export default {
	props: ['name', 'namespace'],
	data () {
		return {
			nodes: []
		}
	},
	methods: {
		retry () {
			let re = this.$axios.post(`/workflow/${this.namespace}/${this.name}/retry`)
			window.console.log(re)
		}
	},
}
</script>


