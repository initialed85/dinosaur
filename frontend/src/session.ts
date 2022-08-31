import { PORT } from './config';

interface CreateSessionResponse {
    uuid: string;
    port: number;
    internal_url: string;
}

export interface Session {
    uuid: string;
    port: number;
    internalUrl: string;
}

export const getSessionAPIURL = (path: string): string => {
    const url = `${window.location.protocol}//${window.location.hostname}:${PORT}/${path}`;
    console.log(url);
    return url;
};

let createSessionInFlight = false;

export const createSession = async (language: string): Promise<Session> => {
    if (createSessionInFlight) {
        throw new Error('createSession already in-flight');
    }

    createSessionInFlight = true;

    const r = await fetch(getSessionAPIURL(`create_session/${language}`));

    const response = (await r.json()) as CreateSessionResponse;

    createSessionInFlight = false;

    return {
        uuid: response.uuid,
        port: response.port,
        internalUrl: response.internal_url
    } as Session;
};

export const pushToSession = async (session: Session, data: string): Promise<void> => {
    const r = await fetch(getSessionAPIURL(`push_to_session/${session.uuid}/`), {
        method: 'POST',
        body: JSON.stringify({
            data: data
        })
    });

    // TODO read and validate
    await r.json();
};
