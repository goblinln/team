/**
 * 服务器返回协议
 */
export interface IResponse {
    /**
     * 错误信息
     */
    err?: string;
    /**
     * 正确响应时的数据
     */
    data?: any;
}

/**
 * 网络请求状态通知
 */
export class FetchStateNotifier {
    _cb: (state: boolean) => void;

    constructor(cb: (state: boolean) => void) {
        this._cb = cb;
    }

    invoke(state: boolean) {
        this._cb && this._cb(state);
    }
}

/**
 * 针对项目对Fetch简单封装一下
 */
export const Fetch = {
    notifier: new FetchStateNotifier(null),

    _p: (method: string, url: string, data?: any, callback?: (rsp: IResponse) => void) => {
        let param : RequestInit = { method: method, body: data };
        Fetch.notifier.invoke(true);

        if (callback == null) {
            fetch(url, param).catch(e => Fetch.notifier.invoke(false)).then(() => Fetch.notifier.invoke(false))
        } else {
            fetch(url, param)
                .then(res => res.json())
                .then(rsp => {
                    callback(rsp);
                    Fetch.notifier.invoke(false);
                })
                .catch(() => Fetch.notifier.invoke(false));
        }
    },

    get: (url: string, callback: (rsp: IResponse) => void) => {
        Fetch.notifier.invoke(true);

        fetch(url)
            .then(res => res.json())
            .then(rsp => {
                callback(rsp);
                Fetch.notifier.invoke(false);
            })
            .catch(() => Fetch.notifier.invoke(false));
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
        Fetch.notifier.invoke(true)

        fetch(url, {method: 'DELETE'})
            .then(res => res.json())
            .then(rsp => {
                callback(rsp);
                Fetch.notifier.invoke(false);
            })
            .catch(() => Fetch.notifier.invoke(false));
    }
}