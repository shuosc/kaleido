<template>
    <v-data-table :headers="headers" :items="mirrors" hide-actions>
        <template slot="items" slot-scope="props">
            <td>{{ props.item.name }}</td>
            <td class="text-xs-right">
                <chip-group :data="props.item.mirrorStations"></chip-group>
            </td>
        </template>
    </v-data-table>
</template>

<script lang="ts">
    import {Component, Vue} from 'vue-property-decorator';
    import ChipGroup from "../components/ChipGroup.vue";

    const getMirrors = require('../graphql/mirrors.graphql');

    @Component({
        components: {ChipGroup}
    })
    export default class Mirrors extends Vue {
        mirrors: Array<{
            names: string,
            mirrorStations: Array<{ name: string }>
        }> = [];
        headers = [
            {
                text: 'name',
                align: 'left',
                value: 'name'
            },
            {
                text: 'stations',
                value: 'mirrorStations',
                align: 'right',
                sortable: false,
            }
        ];

        mounted() {
            getMirrors().then((data: any) => {
                this.mirrors = data.data.data.mirrors;
            });
        }
    }
</script>

<style scoped>

</style>