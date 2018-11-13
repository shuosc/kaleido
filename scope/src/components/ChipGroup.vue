<template>
    <div class="chip-group">
        <div>
            <v-btn fab dark small color="teal" @click="expanded = !expanded">
                <v-icon :style="{ transform:'rotate('+180*expanded+'deg)' }">keyboard_arrow_down</v-icon>
            </v-btn>
        </div>
        <transition name="fade">
            <div v-if="expanded">
                <v-chip v-for="d in data" :key="d.name">{{d.name}}</v-chip>
            </div>
        </transition>
    </div>
</template>

<script lang="ts">
    import {Component, Prop, Vue} from 'vue-property-decorator';

    @Component({})
    export default class Home extends Vue {
        @Prop()
        data!: Array<{ name: string }>;

        expanded: boolean = false;
    }
</script>

<style lang="stylus" scoped>
    .chip-group
        max-width 400px
        min-width 400px
        transition all 0.5s

    .fade-enter-active,
    .fade-leave-active
        transition: opacity .5s;

    .fade-enter,
    .fade-leave-to
        opacity: 0;
</style>