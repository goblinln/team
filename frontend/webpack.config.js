const path = require('path');
const webpack = require('webpack');
const tsImport = require('ts-import-plugin');

module.exports = {
    mode: 'production',
    entry: path.resolve(__dirname, './src/App.tsx'),
    output: {
        path: path.resolve(__dirname, './dist'),
        filename: 'app.bundle.js',
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.js']
    },
//  devtool: 'inline-source-map',
    module: {
        rules: [
            {
                test: /\.ts[x]?$/,
                exclude: /node_modules/,
                use: {
                    loader: 'ts-loader',
                    options: {
                        getCustomTransformers: () => ({
                            before: [
                                tsImport({
                                    libraryName: 'antd',
                                    libraryDirectory: 'lib',
                                    style: 'css',
                                }),
                            ]
                        })
                    }
                }
            },
            {
                test: /\.css$/,
                loader: 'style-loader!css-loader',
            }
        ]
    },
    plugins: [
        new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/)
    ]
};