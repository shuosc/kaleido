module.exports = {
    chainWebpack: config => {
        config.module
            .rule('graphql')
            .test(/\.graphql$/)
            .use('gq-loader')
            .loader('gq-loader')
            .options({
                url: 'http://localhost:8086/graphql'
            })
            .end()
    }
};