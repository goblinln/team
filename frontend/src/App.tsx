import * as React from 'react';
import * as ReactDOM from 'react-dom';

import { Install } from './components/install/Install';
import { Login } from './components/user/Login';
import { Index } from './components/index/Index';

import './App.css';

let mount = document.getElementById('app');
switch (location.pathname) {
case '/': ReactDOM.render(<Index />, mount); break;
case '/install': ReactDOM.render(<Install />, mount); break;
case '/login': ReactDOM.render(<Login />, mount); break;
}