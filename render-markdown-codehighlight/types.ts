import { RefObject } from 'react';

export type pluginHookProps = HTMLElement | RefObject<HTMLElement> | null;

export interface Request {
    get: (url: string) => Promise<Response>;
}
