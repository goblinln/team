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
//    devtool: 'inline-source-map',
    devServer: {
        historyApiFallback: true,
        before: function (app, server) {
            app.get('/', (req, rsp) => {
                rsp.writeHead(200, {'Content-Type': 'text/html'});
                rsp.end(`<!DOCTYPE html>
                    <html>
                        <head>
                            <meta charset="utf-8">
                            <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
                            <title>模拟 - 协作系统</title>
                            <link rel="shortcut icon" href="./dist/favicon.ico" />
                        </head>
                        <body>
                            <div id="app"></div>                
                            <script src="./app.bundle.js"></script>
                        </body>
                    </html>`);
            });
        }
    },
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