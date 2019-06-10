import * as React from 'react';
import * as ReactDOM from 'react-dom';

import { Index } from './components/index/Index';
import { Install } from './components/install/Install';

import './App.css';

let appMount = document.getElementById('app');
if (appMount) {
    ReactDOM.render(<Index />, appMount);
} else {
    ReactDOM.render(<Install />, document.getElementById('install'));
}