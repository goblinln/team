import {Loading, Notification} from '../components';

interface Response {
    err?: string;
    data?: any;
};

interface Parameter {
    url: string;
    method?: 'GET'|'POST'|'PUT'|'DELETE';
    data?: any;
    dontShowLoading?: boolean;
    success?: (data: any) => void;
}

export const request = (param: Parameter) => {
    let init : RequestInit = { method: param.method||'GET', body: param.data, credentials: "include" };
    let finish = () => { !param.dontShowLoading&&Loading.hide(); }

    if (!param.dontShowLoading) Loading.show();

    if (param.success == null) {
        fetch(param.url, init).catch(() => null).then(finish);
    } else {
        fetch(param.url, init)
            .then(res => res.json())
            .then((rsp: Response) => rsp.err?Notification.alert(rsp.err, 'error'):param.success(rsp.data))
            .catch(() => null)
            .then(finish);
    }
};