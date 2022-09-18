export let SESSION_API_URL = '';

if (process.env.NODE_ENV === 'development') {
    SESSION_API_URL = 'http://localhost:8080';
}

export const PUSH_TO_SESSION_DEBOUNCE = 1_000;
export const HEARTBEAT_FOR_SESSION_INTERVAL = 5_000;
