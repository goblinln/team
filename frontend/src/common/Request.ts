export interface IResponse {
    err?: string;
    data?: any;
}

/**
 * 针对项目对Fetch简单封装一下
 */
export const Fetch = {

    _p: (method: string, url: string, data?: any, callback?: (rsp: IResponse) => void) => {
        let param : RequestInit = { method: method, body: data };
        if (callback == null) {
            fetch(url, param).catch(e => console.error(e))
        } else {
            fetch(url, param)
                .then(res => res.json())
                .then(rsp => callback(rsp))
                .catch(e => console.error(e))
        }
    },

    get: (url: string, callback: (rsp: IResponse) => void) => {
        fetch(url)
            .then(res => res.json())
            .then(rsp => callback(rsp))
            .catch(e => console.error(e));
    },

    post: (url: string, data?: any, callback?: (rsp: IResponse) => void) => {
        Fetch._p('POST', url, data, callback)
    },

    put: (url: string, data?: any, callback?: (rsp: IResponse) => void) => {
        Fetch._p('PUT', url, data, callback)
    },

    patch: (url: string, data?: any, callback?: (rsp: IResponse) => void) => {
        Fetch._p('PATCH', url, data, callback)
    },

    delete: (url: string, callback?: (rsp: IResponse) => void) => {
        fetch(url, {method: 'DELETE'})
            .then(res => res.json())
            .then(rsp => callback(rsp))
            .catch(e => console.log(e));
    }
}