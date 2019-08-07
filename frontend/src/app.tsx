import * as React from 'react';
import * as ReactDOM from 'react-dom';

import {Install} from './pages/install';
import {Login} from './pages/login';
import {Home} from './pages/home';

const App = () => {
    switch (location.pathname) {
    case '/install':
        return <Install/>;
    case '/login':
        return <Login/>;
    default:
        return <Home/>;
    }
};

ReactDOM.render(<App/>, document.getElementById('app'));
