var path = require('path')
const webpack = require('webpack')
const public = path.resolve(__dirname, 'public')

module.exports = {
    mode: 'production',
    entry: public + '/index.js',
    output: {
        path: public,
        filename: 'bundle.js'
    },
    module: {
        rules: [
            {
                test: /\.js?$/,
                include: public,
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