export let SESSION_API_URL = '';

if (process.env.NODE_ENV === 'development') {
    SESSION_API_URL = 'http://localhost:8080';
}
