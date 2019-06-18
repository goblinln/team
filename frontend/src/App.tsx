import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as moment from 'moment';

import { Install } from './components/install/Install';
import { Login } from './components/user/Login';
import { Index } from './components/index/Index';

import 'moment/locale/zh-cn';
import './App.css';

moment.locale('zh-cn');

let mount = document.getElementById('app');
switch (location.pathname) {
case '/': ReactDOM.render(<Index />, mount); break;
case '/install': ReactDOM.render(<Install />, mount); break;
case '/login': ReactDOM.render(<Login />, mount); break;
}