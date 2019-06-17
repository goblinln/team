import * as React from 'react';

import {
    Avatar,
    Card,
    Checkbox,
    Col,
    Button,
    Divider,
    Empty,
    Icon,
    Input,
    Modal,
    Row,
    Select,
    Table,
    Tag,
    Popconfirm,
    message,
} from 'antd';

import { IProject, IProjectMember, IUser } from '../../common/Protocol';
import { ProjectRole } from '../../common/Consts';
import { Fetch } from '../../common/Request';

/**
 * 项目管理页
 */
export const Manage = (props: {pid: number}) => {

    /**
     * 自己管理项目数据
     */
    const [proj, setProj] = React.useState<IProject>(null);

    /**
     * 启动拉取项目信息
     */
    React.useEffect(() => fetchProject(), []);

    /**
     * 成员列表的表格格式定义
     */
    const memeberColumns = [
        {
            title: '头像',
            dataIndex: 'user.avatar',
            key: 'avatar',
            render: (text: any, record: IProjectMember, index: number) => <Avatar icon='user' src={record.user.avatar}/>
        },
        {
            title: '昵称',
            dataIndex: 'user.name',
            key: 'name',
        },
        {
            title: '帐号',
            dataIndex: 'user.account',
            key: 'account',
        },
        {
            title: '角色',
            dataIndex: 'role',
            key: 'role',
            render: (text: any, record: IProjectMember, index: number) => <Tag color='#108ee9'>{ProjectRole[record.role]}</Tag>
        },
        {
            title: '管理权限',
            dataIndex: 'isAdmin',
            key: 'isAdmin',
            render: (text: any, record: IProjectMember, index: number) => <Checkbox checked={record.isAdmin}/>
        },
        {
            title: '操作',
            key: 'options',
            render: (text: any, record: IProjectMember, index: number) => (
                <span>
                    <a onClick={() => editMember(record)}>编辑</a>
                    <Divider type="vertical" />
                    <Popconfirm okText='是的' cancelText='手滑了' title={`确定要删除成员【${record.user.name}】吗？`} onConfirm={() => delMember(record)}>
                        <a>删除</a>
                    </Popconfirm>
                </span>
            )
        }
    ];

    /**
     * 拉取项目信息
     */
    const fetchProject = () => {
        Fetch.get(`/api/project/${props.pid}`, rsp => {
            if (rsp.err) {
                message.error(rsp.err, 1);
            } else {
                rsp.data.members = rsp.data.members.sort((a: IProjectMember, b: IProjectMember) => {
                    return a.user.account.localeCompare(b.user.account);
                })
                setProj(rsp.data);
            }
        });
    }

    /**
     * 添加分支
     */
    const addBranch = () => {
        let branch: string = '';

        Modal.confirm({
            title: '添加分支',
            width: 200,
            maskClosable: true,
            icon: null,
            content: <Input style={{marginTop: 24}} onChange={ev => branch = ev.target.value}/>,
            okText: '提交',
            onOk: () => {
                if (branch.length == 0) {
                    message.error('分支不可为空', 1);
                    return;
                }

                if (proj.branches.indexOf(branch) >= 0) {
                    message.error('分支已存在', 1);
                    return;
                }

                let param = new FormData();
                param.append('branch', branch);

                Fetch.post(`/api/project/${props.pid}/branch`, param, rsp => {
                    rsp.err ? message.error(rsp.err, 1) : fetchProject();
                });
            },
            cancelText: '取消',
        });
    };

    /**
     * 添加成员
     */
    const addMember = () => {
        Fetch.get(`/api/project/${props.pid}/invites`, rsp => {
            if (rsp.err) {
                message.error(rsp.err, 1);
                return;
            }

            let invites: IUser[] = rsp.data;
            let selUser: number = -1;
            let role: number = 0;
            let isAdmin: boolean = false;

            Modal.confirm({
                title: '邀请成员',
                width: 500,
                maskClosable: true,
                icon: null,
                content: (
                    <Row style={{marginTop: 24}} gutter={8} type='flex' justify='center'>
                        <Col>
                            用户：
                            <Select size='small' style={{width: 180}} onChange={(ev: any) => selUser = ev.valueOf()}>
                                {invites.map((user, idx) => {
                                    return <Select.Option key={idx} value={user.id}>{user.name}</Select.Option>
                                })}
                            </Select>
                        </Col>                           

                        <Col>
                            职能：
                            <Select size='small' defaultValue={role} style={{width: 80}} onChange={(ev: number) => role = ev}>
                                {ProjectRole.map((name, idx) => {
                                    return <Select.Option key={idx} value={idx}>{name}</Select.Option>
                                })}
                            </Select>
                        </Col>

                        <Col>            
                            管理权限：<Checkbox defaultChecked={isAdmin} onChange={ev => isAdmin = ev.target.checked} />
                        </Col>
                    </Row>
                ),
                okText: '添加',
                onOk: () => {
                    let param = new FormData();
                    param.append('uid', selUser.toString());
                    param.append('isAdmin', isAdmin ? '1' : '0');
                    param.append('role', role.toString());

                    Fetch.post(`/api/project/${props.pid}/member`, param, rsp => {
                        rsp.err ? message.error(rsp.err, 1) : fetchProject();
                    })
                },
                cancelText: '取消',
            });
        })                
    };

    /**
     * 编辑成员
     */
    const editMember = (member: IProjectMember) => {
        let isAdmin = member.isAdmin || false;
        let role = member.role;

        Modal.confirm({
            title: `编辑成员 - ${member.user.name}`,
            width: 280,
            maskClosable: true,
            icon: null,
            content: (
                <Row style={{marginTop: 24}} type='flex' justify='start' gutter={32}>
                    <Col>
                        职能：
                        <Select size='small' defaultValue={role} style={{width: 80}} onChange={(ev: number) => role = ev}>
                            {ProjectRole.map((name, idx) => {
                                return <Select.Option key={idx} value={idx}>{name}</Select.Option>
                            })}
                        </Select>
                    </Col>

                    <Col>            
                        管理权限：<Checkbox defaultChecked={isAdmin} onChange={ev => isAdmin = ev.target.checked} />
                    </Col>
                </Row>
            ),
            okText: '提交',
            onOk: () => {
                let param = new FormData();
                param.append('role', role.toString());
                param.append('isAdmin', isAdmin ? '1' : '0');

                Fetch.put(`/api/project/${props.pid}/member/${member.user.id}`, param, rsp => {
                    rsp.err ? message.error(rsp.err, 1) : fetchProject();
                })
            },
            cancelText: '取消',
        });
    }

    /**
     * 删除成员
     */
    const delMember = (member: IProjectMember) => {
        Fetch.delete(`/api/project/${props.pid}/member/${member.user.id}`, rsp => {
            rsp.err ? message.error(rsp.err, 1) : fetchProject();
        });
    };

    return (
        <div style={{padding: 16}}>
            <Row>
                <Card title={<div><Icon type='branches' style={{margin: '0 8px'}}/>分支列表</div>} extra={<Button type='link' onClick={() => addBranch()}><Icon type='plus'/>添加分支</Button>}>
                    <div>
                        {proj ? proj.branches.map((branch, idx) => {
                            return <Tag key={idx} color='#108ee9'>{branch}</Tag>
                        }) : <Empty description='暂无数据'/>}
                    </div>
                </Card>
            </Row>

            <Row style={{marginTop: 16}}>
                <Card title={<div><Icon type='idcard' style={{margin: '0 8px'}}/>成员列表</div>} extra={<Button type='link' onClick={() => addMember()}><Icon type='plus'/>添加成员</Button>}>
                    <Table size='middle' columns={memeberColumns} dataSource={proj ? proj.members : null} bordered/>
                </Card>
            </Row>
        </div>
    );
};