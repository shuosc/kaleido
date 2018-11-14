<template>
    <div class="home">
        <v-app id="inspire">
            <v-navigation-drawer v-model="drawer" fixed app>
                <v-list dense>
                    <v-list-tile @click="$router.push('mirrors')">
                        <v-list-tile-action>
                            <v-icon>save</v-icon>
                        </v-list-tile-action>
                        <v-list-tile-content>
                            <v-list-tile-title>镜像</v-list-tile-title>
                        </v-list-tile-content>
                    </v-list-tile>
                    <v-list-tile @click="$router.push('stations')">
                        <v-list-tile-action>
                            <v-icon>web</v-icon>
                        </v-list-tile-action>
                        <v-list-tile-content>
                            <v-list-tile-title>镜像站</v-list-tile-title>
                        </v-list-tile-content>
                    </v-list-tile>
                </v-list>
            </v-navigation-drawer>
            <v-toolbar color="#219c59" dark fixed app>
                <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
                <v-toolbar-title>KaleidoScope</v-toolbar-title>
            </v-toolbar>
            <v-content>
                <v-container fluid fill-height>
                    <v-layout justify-center align-center>
                        <router-view></router-view>
                    </v-layout>
                </v-container>
            </v-content>
            <v-footer color="#219c59" app>
                <span class="white--text">&copy; SHUOSC</span>
            </v-footer>
        </v-app>
    </div>
</template>

<style lang="stylus">
    #app
        font-family 'Avenir', Helvetica, Arial, sans-serif
        -webkit-font-smoothing antialiased
        -moz-osx-font-smoothing grayscale
        text-align center
        color #2c3e50

    #nav
        padding 30px
        a
            font-weight bold
            color #2c3e50
            &.router-link-exact-active
                color #42b983
</style>

<script>
    export default {
        name: 'app',
        components: {},
        data() {
            return {
                drawer: false
            }
        }
    }
</script>
