import React, { useEffect, useState } from 'react';
import './App.css';
import { Editor } from './Editor';
import { Shell } from './Shell';
import { createSession, pushToSession, Session } from './session';

const supportedLanguages = [
    {
        name: 'go',
        friendlyName: 'Go'
    },
    {
        name: 'python',
        friendlyName: 'Python'
    },
    {
        name: 'typescript',
        friendlyName: 'TypeScript'
    },
    {
        name: 'c',
        friendlyName: 'C'
    },
    {
        name: 'rust',
        friendlyName: 'Rust'
    },
    {
        name: 'java',
        friendlyName: 'Java'
    }
];

export interface AppProps {
    language: string | null;
}

export function App(props: AppProps) {
    const [language, setLanguage] = useState(props.language);
    const [session, setSession] = useState(null);
    const [editorValue, setEditorValue] = useState('');
    const [error, setError] = useState(null);

    useEffect(() => {
        if (!language) {
            return;
        }

        if (!session) {
            createSession(language)
                .then((x: Session) => {
                    setSession(x as any);
                })
                .catch(e => {
                    if (!error) {
                        setError(e.toString());
                    }
                });
            return;
        }

        if (!editorValue) {
            return;
        }

        pushToSession(session, editorValue)
            .then(x => {
                // noop
            })
            .catch(e => {
                // noop
            });
    }, [language, editorValue, session, error, props.language]);

    const buttons: any[] = [];

    supportedLanguages.forEach((x, i) => {
        buttons.push(
            <button
                key={`button-language-selection-${i}`}
                className="button-language-selection"
                onClick={() => {
                    setLanguage(x.name);
                }}
            >
                {x.friendlyName}
            </button>
        );
    });

    return (
        <div className="outer-container">
            <div></div>
            {!language ? (
                <div className="inner-container-language-selection">{buttons}</div>
            ) : (
                <div className="inner-container">
                    <div className="editor-item">
                        {session ? (
                            <Editor
                                language={language}
                                code={(session as unknown as Session).code}
                                setEditorValue={(x: string): void => {
                                    setEditorValue(x);
                                }}
                            />
                        ) : (
                            <div className="shell-iframe">
                                Attempting to interact with backend... <br />
                                <br />
                                {error}
                            </div>
                        )}
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
            )}
            <div></div>
        </div>
    );
}
