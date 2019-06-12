const path = require('path');
const tsImport = require('ts-import-plugin');

module.exports = {
    mode: 'production',
    entry: path.resolve(__dirname, './src/App.tsx'),
    output: {
        path: path.resolve(__dirname, '../publish/www'),
        filename: 'app.js',
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
                        transpileOnly: true,
                        compilerOptions: {
                            target: 'es5'
                        },
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
    }
};