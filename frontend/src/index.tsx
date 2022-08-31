import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import { App } from './App';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);

const language = window.location.pathname.split('/')[1] || 'go';

root.render(
    <React.StrictMode>
        <App language={language} />
    </React.StrictMode>
);
