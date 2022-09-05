import { SESSION_API_URL } from './config';

interface ErrorResponse {
    error: string;
}

interface CreateSessionResponse {
    uuid: string;
    port: number;
    internal_url: string;
    code: string;
}

export interface Session {
    uuid: string;
    port: number;
    internalUrl: string;
    code: string;
}

export const getSessionAPIURL = (path: string): string => {
    return `${SESSION_API_URL}/${path}`;
};

let createSessionInFlight = false;

export const createSession = async (language: string): Promise<Session> => {
    if (createSessionInFlight) {
        throw new Error('createSession already in-flight');
    }

    createSessionInFlight = true;

    const r = await fetch(getSessionAPIURL(`create_session/${language}`));

    if (r.status !== 201) {
        const errorResponse = (await r.json()) as ErrorResponse;
        throw new Error(errorResponse.error);
    }

    const response = (await r.json()) as CreateSessionResponse;

    createSessionInFlight = false;

    return {
        uuid: response.uuid,
        port: response.port,
        internalUrl: response.internal_url,
        code: response.code
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

export const heartbeatForSession = async (session: Session): Promise<void> => {
    const r = await fetch(getSessionAPIURL(`heartbeat_for_session/${session.uuid}/`));

    // TODO read and validate
    await r.json();
};
