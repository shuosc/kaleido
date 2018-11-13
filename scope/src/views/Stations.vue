<template>
    <v-data-table :headers="headers" :items="mirrorStations" hide-actions>
        <template slot="items" slot-scope="props">
            <td>{{ props.item.name }}</td>
            <td class="text-xs-right">
                <chip-group :data="props.item.mirrors"></chip-group>
            </td>
        </template>
    </v-data-table>
</template>

<script lang="ts">
    import {Component, Vue} from "vue-property-decorator";
    import ChipGroup from "../components/ChipGroup.vue";

    const getMirrorStations = require('../graphql/mirrorStations.graphql');

    @Component({
        components: {ChipGroup}
    })
    export default class Mirrors extends Vue {
        mirrorStations: Array<{
            names: string,
            mirrors: Array<{ name: string }>
        }> = [];
        headers = [
            {
                text: 'name',
                align: 'left',
                value: 'name'
            },
            {
                text: 'mirrors',
                value: 'mirrors',
                align: 'right',
                sortable: false,
            }
        ];

        mounted() {
            getMirrorStations().then((data: any) => {
                this.mirrorStations = data.data.data.mirrorStations;
            });
        }
    }
</script>

<style scoped>

</style>