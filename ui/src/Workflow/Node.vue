<template>
<b-row>
	<b-col v-if="isPod()">
		<b-link :to="`/watch/${namespace}/pods/${content.id}`">{{content.displayName}}</b-link>
	</b-col>
	<b-col v-else>
		{{content.displayName}}
	</b-col>
	<b-col>{{content.type}}</b-col>
	<b-col>{{content.phase}}</b-col>
	<b-col>{{content.startedAt}}</b-col>
</b-row>
</template>

<script>
export default {
	props: ['content', 'namespace'],
	methods: {
		isPod () {
			return this.content.type == 'Pod'
		}
	}
}
</script>