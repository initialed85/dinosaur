import React, { useEffect, useState } from 'react';
import './Shell.css';
import { getSessionAPIURL, heartbeatForSession, Session } from './session';

export interface ShellProps {
    session: Session;
}

export function Shell(props: ShellProps) {
    const [heartbeat, setHeartbeat] = useState(null);

    useEffect(() => {
        if (!heartbeat) {
            setHeartbeat(
                setInterval(async () => {
                    await heartbeatForSession(props.session);
                }, 1_000) as any
            );
        }
    });

    return (
        // TODO properly integrate I guess xterm.js rather than iframe to gotty's usage thereof
        <iframe
            title="shell"
            className="shell-iframe"
            src={getSessionAPIURL(`proxy_session/${props.session.uuid}`)}
        ></iframe>
    );
}
