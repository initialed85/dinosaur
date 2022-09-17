import React, { useEffect, useState } from 'react';
import './Shell.css';
import { getSessionAPIURL, heartbeatForSession, Session } from './session';

export interface IFrameProps {
    src: string;
}

export interface ShellProps {
    session: Session;
}

function IFrame(props: IFrameProps) {
    return (
        // TODO properly integrate I guess xterm.js rather than iframe to gotty's usage thereof
        <iframe title="shell" className="shell-iframe" src={props.src}></iframe>
    );
}

export function Shell(props: ShellProps) {
    const [heartbeat, setHeartbeat] = useState(null);

    useEffect(() => {
        if (!heartbeat) {
            setHeartbeat(
                setInterval(async () => {
                    try {
                        await heartbeatForSession(props.session);
                    } catch (e) {
                        // noop
                    }
                }, 1_000) as any
            );
        }
    }, [heartbeat, props.session]);

    return <IFrame src={getSessionAPIURL(`proxy_session/${props.session.uuid}`)}></IFrame>;
}
