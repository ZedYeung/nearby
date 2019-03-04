import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './components/App';
import { CookiesProvider } from 'react-cookie';
import { BrowserRouter } from 'react-router-dom';
import * as serviceWorker from './serviceWorker';

ReactDOM.render(
    <BrowserRouter>
      <CookiesProvider>
        <App />
      </CookiesProvider>
    </BrowserRouter>, document.getElementById('root'));
serviceWorker.register();
