import React, { useEffect, useState } from 'react';
import './App.css';
import { Editor } from './Editor';
import { Shell } from './Shell';
import { createSession, pushToSession, Session } from './session';

export interface AppProps {
    language: string;
}

export function App(props: AppProps) {
    const [editorValue, setEditorValue] = useState('');
    const [session, setSession] = useState(null);

    useEffect(() => {
        if (!session) {
            createSession(props.language)
                .then((x: Session) => {
                    setSession(x as any);
                })
                .catch(e => {});
        }

        if (session && editorValue) {
            pushToSession(session, editorValue)
                .then(x => {
                    console.log(x);
                })
                .catch(e => {
                    console.log(e);
                });
        }
    });

    return (
        <div className="outer-container">
            <div></div>
            <div className="inner-container">
                <div className="editor-item">
                    <Editor
                        language={'go'}
                        setEditorValue={(x: string): void => {
                            setEditorValue(x);
                        }}
                    />
                </div>
                <div className="shell-item">
                    {session ? (
                        <Shell sessionUUID={(session as Session).uuid} />
                    ) : (
                        <div className="shell-iframe">ERROR: Failed to contact session service</div>
                    )}
                </div>
            </div>
            <div></div>
        </div>
    );
}
