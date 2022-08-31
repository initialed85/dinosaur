import React from 'react';
import './Editor.css';

import MonacoEditor, { Monaco } from '@monaco-editor/react';
import { editor } from 'monaco-editor';

export interface EditorProps {
    language: string;
    setEditorValue: (x: string) => void;
}

const goDefaultValue = `
package main

import "log"

func main() {
    log.Printf("Hello, world.")
}
`;

const defaultValueByLanguage = new Map();
defaultValueByLanguage.set('go', goDefaultValue);

export function Editor(props: EditorProps) {
    let timeout: NodeJS.Timeout | null;

    const handleEditorDidMount = (e: editor.ICodeEditor, m: Monaco) => {
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
            defaultValue={defaultValueByLanguage.get(props.language).trim() + '\n' || ''}
            options={{
                minimap: { enabled: false },
                wordBasedSuggestions: false,
                contextmenu: false,
                fontSize: 12,
                fontFamily: 'monospace',
                formatOnPaste: true,
                formatOnType: true
            }}
            onMount={handleEditorDidMount}
        />
    );
}
