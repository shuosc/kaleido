<template>
    <div>
        <h1>kaleido</h1>
        <h2>A high-speed http redirector for open-source mirror site.</h2>
        本工具将会根据您的IP所处在的地理位置，将您的请求自动定向到较近的镜像站。
        <h3>How to use</h3>
        <p>
            将对应的官方源的URL中的域名改成<code>mirrors.rocks</code>即可。
        </p>
        <p>
            以Ubuntu为例，更改<code>/etc/apt/sources.list</code>文件中Ubuntu 默认的源地址<code>http://archive.ubuntu.com/</code> 为
            <code>http://mirrors.rocks</code>
            即可。
        <p>
            ⚠️：暂时没有HTTPS的支持，对于有些镜像（如pip）请强制允许使用http协议。
        </p>
    </div>
</template>

<script lang="ts">
    import {Component, Vue} from 'vue-property-decorator';

    @Component({})
    export default class Home extends Vue {
        drawer: boolean = false;
    }
</script>
