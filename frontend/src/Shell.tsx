import React from 'react';
import './Shell.css';
import { getSessionAPIURL } from './session';

export interface ShellProps {
    sessionUUID: string;
}

export function Shell(props: ShellProps) {
    return (
        // TODO properly integrate I guess xterm.js rather than iframe to gotty's usage thereof
        <iframe
            title="shell"
            className="shell-iframe"
            src={getSessionAPIURL(`proxy_session/${props.sessionUUID}`)}
        ></iframe>
    );
}
