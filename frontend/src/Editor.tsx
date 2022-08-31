import React from 'react';
import './Editor.css';

import MonacoEditor, {Monaco} from '@monaco-editor/react';
import {editor} from 'monaco-editor';

export interface EditorProps {
    language: string;
    code: string;
    setEditorValue: (x: string) => void;
}

export function Editor(props: EditorProps) {
    let timeout: NodeJS.Timeout | null;

    const handleEditorDidMount = (e: editor.ICodeEditor, m: Monaco) => {
        // TODO this kills Ctrl + C and Ctrl + V; kinda just wanna kill it externally?
        // e.onKeyDown(event => {
        //     const { keyCode, ctrlKey, metaKey } = event;
        //     if (
        //         (keyCode === KeyCode.KeyC ||
        //             keyCode === KeyCode.KeyV) &&
        //         (metaKey || ctrlKey)
        //     ) {
        //         event.preventDefault();
        //     }
        // });

        props.setEditorValue(e.getValue());

        e.onDidChangeModelContent((event: editor.IModelContentChangedEvent) => {
            if (timeout) {
                clearTimeout(timeout);
                timeout = null;
            }

            timeout = setTimeout(() => {
                timeout = null;
                props.setEditorValue(e.getValue());
            }, 1_000);
        });
    };

    return (
        <MonacoEditor
            height="100%"
            theme="vs-dark"
            defaultLanguage={props.language}
            defaultValue={props.code}
            options={{
                minimap: {enabled: false},
                wordBasedSuggestions: false,
                contextmenu: false,
                fontSize: 12,
                fontFamily: 'monospace',
                formatOnPaste: true,
                formatOnType: true,
                scrollBeyondLastLine: false
            }}
            onMount={handleEditorDidMount}
        />
    );
}
