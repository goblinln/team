import * as React from 'react';
import { UploadChangeParam } from 'antd/lib/upload/interface';

import {
    Button,
    Divider,
    Icon,
    Layout,
    Popconfirm,
    Progress,
    Row,
    Table,
    Upload,
    message,
} from 'antd';

import { IShare } from '../../common/Protocol';
import { Fetch } from '../../common/Request';

/**
 * 分享页
 */
export const Page = () => {

    /**
     * 状态列表
     */
    const [files, setFiles] = React.useState<IShare[]>([]);
    const [progress, setProgress] = React.useState<{value: number}>(null);

    /**
     * 分享列表数据结构
     */
    const shareTableColumn = [
        {
            title: '文件',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: '上传者',
            dataIndex: 'uploader',
            key: 'uploader',
        },
        {
            title: '上传时间',
            dataIndex: 'time',
            key: 'time',
        },
        {
            title: '大小',
            dataIndex: 'size',
            key: 'size',
            render: (text: any, record: IShare, index: number) => {
                let size = record.size;
                if (size > 1024 * 1024) return `${(size / (1024 * 1024)).toFixed(3)} MB`;
                if (size > 1024) return `${(size / 1024).toFixed(3)} KB`;
                return `${size} B`;
            }
        },
        {
            title: '操作',
            key: 'options',
            render: (text: any, record: IShare, index: number) => (
                <span>
                    <a href={`/api/file/share/${record.id}`}>下载</a>
                    <Divider type="vertical" />
                    <Popconfirm okText='是的' cancelText='手滑了' title={`确定要删除文件【${record.name}】吗？`} onConfirm={() => deleteShare(record.id)}>
                        <a>删除</a>
                    </Popconfirm>
                </span>
            )
        }
    ];

    /**
     * 初始化时拉取一次
     */
    React.useEffect(() => fetchAll(), []);

    /**
     * 拉取列表
     */
    const fetchAll = () => {
        Fetch.get('/api/file/share/list', rsp => {
            rsp.err ? message.error(rsp.err, 1) : setFiles(rsp.data);
        });
    };

    /**
     * 上传回调
     */
    const handleUpload = (ev: UploadChangeParam) => {
        if (ev.file.status == 'done') {
            setProgress(null);
            fetchAll();
        } else if (ev.file.status == 'uploading') {
            setProgress({value: ev.event.percent || 0});
        } else {
            setProgress(null);
        }
    };

    /**
     * 删除文件
     */
    const deleteShare = (id: number) => {
        Fetch.delete(`/api/file/share/${id}`, rsp => {
            rsp.err ? message.error(rsp.err, 1) : fetchAll();
        });
    };

    return (
        <Layout style={{width: '100%', height: '100%'}}>
            <Layout.Content style={{padding: 32}}>
                <Row>
                    <Upload name='uploader' action='/api/file/share' showUploadList={false} onChange={handleUpload}>
                        <Button><Icon type="upload" />分享文件</Button>
                    </Upload>
                    {progress && <Progress type='circle' percent={progress.value} />}
                </Row>
                <Row style={{marginTop: 16}}>
                    <Table columns={shareTableColumn} dataSource={files} bordered/>
                </Row>
            </Layout.Content>
        </Layout>
    );
}