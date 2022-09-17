export interface Params {
    language: string;
    sessionUUID: string;
}

export const getParams = (): Params => {
    const urlSearchParams = new URLSearchParams(window.location.search);

    return {
        language: urlSearchParams.get('language'),
        sessionUUID: urlSearchParams.get('session_uuid')
    } as Params;
};

export const setParams = (params: Params) => {
    const urlSearchParams = new URLSearchParams(window.location.search);

    if (params.language) {
        urlSearchParams.set('language', params.language);
    }

    if (params.sessionUUID) {
        urlSearchParams.set('session_uuid', params.sessionUUID);
    }

    window.history.pushState('dinosaur', '__unused__', `?${urlSearchParams.toString()}`);
};
