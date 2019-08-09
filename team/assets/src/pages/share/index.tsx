import * as React from 'react';
import * as moment from 'moment';

import {Input, Button, Icon, Progress, Table, TableColumn, Modal, Notification} from '../../components';
import {Share} from '../../common/protocol';
import { request } from '../../common/request';

export const SharePage = () => {
    const [isUploading, setIsUploading] = React.useState<boolean>(false);
    const [progress, setProgress] = React.useState<number>(0);
    const [files, setFiles] = React.useState<Share[]>([]);

    React.useEffect(() => fetchAll(), []);

    const schema: TableColumn[] = [
        {label: '文件', dataIndex: 'name', align: 'left'},
        {label: '上传者', dataIndex: 'uploader'},
        {label: '上传时间', dataIndex: 'time', sorter: (a: Share, b: Share) => moment(a.time).diff(moment(b.time))},
        {label: '大小', align: 'right', sorter: (a: Share, b: Share) => a.size - b.size, renderer: (data: Share) => formatSize(data)},
        {label: '操作', renderer: (data: Share) => (
            <span>
                <a className='link' href={`/api/file/share/${data.id}`}>下载</a>
                <div className='divider-v'/>
                <a className='link' href='javascript:void(0)' onClick={() => delShare(data)}>删除</a>
            </span>
        )}
    ];

    const fetchAll = () => {
        request({url: '/api/file/share/list', success: setFiles});
    };

    const uploader = (file: File) => {
        setIsUploading(true);
        setProgress(0);

        new Promise((resolve, reject) => {
            let param = new FormData();
            param.append('file', file, file.name);

            let request = new XMLHttpRequest();
            request.open('POST', '/api/file/share');
            request.upload.onprogress = (e: ProgressEvent) => setProgress(e.loaded * 100 / e.total);
            request.onerror = () => reject();
            request.onload = () => request.status == 200?resolve():reject();
            request.send(param);
        })
        .then(() => Notification.alert('上传成功', 'info'), () => Notification.alert('上传失败', 'error'))
        .then(() => {setIsUploading(false); setProgress(0); fetchAll()});
    };

    const formatSize = (data: Share) => {
        let size = data.size;
        if (size > 1024 * 1024 * 1024) return `${(size / (1024 * 1024 * 1024)).toFixed(3)} GB`;
        if (size > 1024 * 1024) return `${(size / (1024 * 1024)).toFixed(3)} MB`;
        if (size > 1024) return `${(size / 1024).toFixed(3)} KB`;
        return `${size} B`;
    };

    const delShare = (data: Share) => {
        Modal.open({
            title: '删除确认',
            body: <div className='my-2'>确定要删除文件【{data.name}】吗？</div>,
            onOk: () => {
                request({url: `/api/file/share/${data.id}`, method: 'DELETE', success: fetchAll});
            }
        });
    };

    return (
        <div className='m-4'>
            <div className='mb-3'>
                <Input.Uploader customUpload={uploader}>
                    <Button size='sm'><Icon type='upload' className='mr-1'/>上传文件</Button>
                </Input.Uploader>
            </div>
            {isUploading&&(
                <div>
                    <Progress percent={progress}/>
                </div>
            )}
            <Table size='lg' dataSource={files} columns={schema} pagination={15}/>
        </div>
    );
};