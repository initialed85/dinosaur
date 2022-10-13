import React, {useEffect, useState} from 'react';
import './App.css';
import {Editor} from './Editor';
import {Shell} from './Shell';
import {createSession, getSession, getSupportedLanguages, pushToSession, Session, SupportedLanguage} from './session';
import {Params, setParams} from './params';

export interface AppProps {
    language: string | null;
    sessionUUID: string | null;
}

export function App(props: AppProps) {
    const [supportedLanguages, setSupportedLanguages] = useState(null);
    const [language, setLanguage] = useState(props.language);
    const [session, setSession] = useState(null);
    const [editorValue, setEditorValue] = useState('');
    const [error, setError] = useState(null);

    const languageSelections: any[] = [];

    useEffect(() => {
        if (!supportedLanguages) {
            getSupportedLanguages()
                .then(receivedSupportedLanguages => {
                    setSupportedLanguages(receivedSupportedLanguages as any);
                })
                .catch(e => {
                    console.error(e);
                });
        }

        if (!language) {
            return;
        }

        setParams({
            language: language
        } as Params);

        if (!session) {
            if (props?.sessionUUID) {
                getSession(props.sessionUUID)
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

        setParams({
            language: language,
            sessionUUID: (session as Session).uuid
        } as Params);

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
    }, [supportedLanguages, language, editorValue, session, error, props]);

    if (supportedLanguages) {
        (supportedLanguages as SupportedLanguage[]).forEach((x, i) => {
            languageSelections.push(
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
    }

    return (
        <div className="outer-container">
            <div></div>
            {!language ? (
                <div className="inner-container-language-selection">{languageSelections}</div>
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
                                Attempting to interact with backend... <br/>
                                <br/>
                                {error}
                            </div>
                        )}
                    </div>
                    <div className="shell-item">
                        {session ? (
                            <Shell session={session}/>
                        ) : (
                            <div className="shell-iframe">
                                Attempting to interact with backend... <br/>
                                <br/>
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
