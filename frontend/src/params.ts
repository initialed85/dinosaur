export interface Params {
    language: string;
}

export const getParams = (): Params => {
    const urlSearchParams = new URLSearchParams(window.location.search);

    return {
        language: urlSearchParams.get('language')
    } as Params;
};
