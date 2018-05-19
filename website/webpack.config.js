var path = require('path')
var webpack = require('webpack')

var app_dir = path.resolve(__dirname, 'source')
var build_dir = path.resolve(__dirname, 'public')

module.exports = {
    mode: 'development',
    devServer: {
        inline: true,
        contentBase: build_dir,
        port: 3000
    },
    entry: app_dir + '/index.js',
    output: {
        path: build_dir,
        filename: 'bundle.js'
    },
    module: {
        rules: [
            {
                test: /\.js?$/,
                include: app_dir,
                exclude: /node_modules/,
                use: [
                    {
                        loader: 'babel-loader',
                        options: {
                            presets: ['react']
                        }
                    }
                ]
            }
        ]
    }
}