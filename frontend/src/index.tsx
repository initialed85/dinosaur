import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import { App } from './App';
import { getParams } from './params';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);

const params = getParams();

root.render(
    <React.StrictMode>
        <App language={params.language} sessionUUID={params.sessionUUID} />
    </React.StrictMode>
);
