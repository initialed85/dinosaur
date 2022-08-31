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
    const [error, setError] = useState(null);

    useEffect(() => {
        if (!session) {
            createSession(props.language)
                .then((x: Session) => {
                    setSession(x as any);
                })
                .catch(e => {
                    if (!error) {
                        setError(e.toString());
                    }
                });
        }

        if (session && editorValue) {
            pushToSession(session, editorValue)
                .then(x => {
                    // noop
                })
                .catch(e => {
                    // noop
                });
        }
    }, [editorValue, session, error, props.language]);

    return (
        <div className="outer-container">
            <div></div>
            <div className="inner-container">
                <div className="editor-item">
                    <Editor
                        language={props.language}
                        setEditorValue={(x: string): void => {
                            setEditorValue(x);
                        }}
                    />
                </div>
                <div className="shell-item">
                    {session ? (
                        <Shell session={session} />
                    ) : (
                        <div className="shell-iframe">
                            Attempting to interact with backend... <br />
                            <br />
                            {error}
                        </div>
                    )}
                </div>
            </div>
            <div></div>
        </div>
    );
}
