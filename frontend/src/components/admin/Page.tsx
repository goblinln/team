import * as React from 'react';
import { WrappedFormUtils } from 'antd/lib/form/Form';

import {
    Button,
    Card,
    Checkbox,
    Col,
    Divider,
    Form,
    Icon,
    Input,
    Layout,
    Modal,
    Popconfirm,
    Row,
    Select,
    Table,
    message,
} from 'antd';

import { IUser, IProject } from '../../common/Protocol';
import { ProjectRole } from '../../common/Consts';
import { Fetch } from '../../common/Request';

/**
 * 系统管理页
 */
export const Page = () => {

    /**
     * 状态列表
     */
    const [users, setUsers] = React.useState<IUser[]>([]);
    const [projs, setProjs] = React.useState<IProject[]>([]);
    const [modal, setModal] = React.useState<React.ReactNode>(null);
    const modalForm = React.useRef<any>(null);

    /**
     * 用户数据表结构
     */
    const userSchema = [
        {
            title: '昵称',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: '帐号',
            dataIndex: 'account',
            key: 'account',
        },
        {
            title: '管理员',
            key: 'isAdmin',
            render: (text: any, record: IUser, index: number) => <Checkbox checked={record.isSu}/>
        },
        {
            title: '操作',
            key: 'options',
            render: (text: any, record: IUser, index: number) => (
                <span>
                    <a onClick={() => showModal("修改帐号信息", <FormEditUser user={record} ref={modalForm}/>)}>编辑</a>
                    <Divider type="vertical" />
                    <a onClick={() => lockUser(record)}>{record.isLocked ? '解锁' : '禁用'}</a>
                    <Divider type="vertical" />
                    <Popconfirm okText='是的' cancelText='手滑了' title={`确定要删除用户【${record.name}】吗？`} onConfirm={() => delUser(record)}>
                        <a>删除</a>
                    </Popconfirm>
                </span>
            )
        }
    ];

    /**
     * 项目数据表结构
     */
    const projectSchema = [
        {
            title: '项目名',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: '操作',
            key: 'options',
            render: (text: any, record: IProject, index: number) => (
                <span>
                    <a onClick={() => showModal("修改项目信息", <FormEditProject proj={record} ref={modalForm}/>)}>编辑</a>
                    <Divider type="vertical" />
                    <Popconfirm okText='是的' cancelText='手滑了' title={`确定要删除项目【${record.name}】吗？`} onConfirm={() => delProject(record)}>
                        <a>删除</a>
                    </Popconfirm>
                </span>
            )
        }
    ];

    /**
     * 初始拉取一次
     */
    React.useEffect(() => {
        fetchUsers();
        fetchProjs();
    }, []);

    /**
     * 拉取所有用户列表
     */
    const fetchUsers = () => {
        Fetch.get('/admin/user/list', rsp => {
            rsp.err ? message.error(rsp.err, 1) : setUsers(rsp.data);
        });
    };

    /**
     * 拉取所有项目列表
     */
    const fetchProjs = () => {
        Fetch.get('/admin/project/list', rsp => {
            rsp.err ? message.error(rsp.err, 1) : setProjs(rsp.data);
        });
    };

    /**
     * 显示相关弹出框
     */
    const showModal = (title: string, form: JSX.Element) => {
        setModal(
            <Modal
                title={title}
                visible={true}
                onCancel={() => setModal(null)}
                footer={null}>
                {form}
            </Modal>
        );
    };

    /**
     * 添加用户
     */
    const FormAddUser = Form.create()((props: {form: WrappedFormUtils}) => {
        const { getFieldDecorator, getFieldValue, setFieldsValue, validateFields } = props.form;

        const validateNewPswd = (rule: any, value: string, callback: (err?: string) => void) => {
            if (value && value !== getFieldValue('pswd')) {
                callback('两次输入的密码不一致！');
            } else {
                callback();
            }
        };

        const submit = (ev: React.FormEvent<HTMLFormElement>) => {
            ev.preventDefault();

            validateFields(err => {
                if (!err) {
                    setModal(null);
                    Fetch.post('/admin/user', new FormData(ev.currentTarget), rsp => {
                        rsp.err ? message.error(rsp.err, 1) : fetchUsers();
                    });
                }
            })
        }

        return (
            <Form onSubmit={submit}>
                <Form.Item label='登录帐号' required={true} style={{marginBottom: 8}}>
                    {getFieldDecorator('account', {
                        rules: [
                            {required: true, message: '登录帐号未填写'},
                            {whitespace: true, message: '帐号不能包含空白字符'},
                            {max: 32, message: '最大32个字符'},
                        ]
                    })(
                        <Input id='account' name='account' />
                    )}
                </Form.Item>

                <Form.Item label='显示名称' required={true} style={{marginBottom: 8}}>
                    {getFieldDecorator('name', {
                        rules: [
                            {required: true, message: '显示名称未填写'},
                            {whitespace: true, message: '名称不能包含空白字符'},
                            {max: 32, message: '最大32个字符'},
                        ]
                    })(
                        <Input id='name' name='name' />
                    )}
                </Form.Item>

                <Form.Item label='初始登录密码' required={true} style={{marginBottom: 8}}>
                    {getFieldDecorator('pswd', {
                        rules: [
                            {required: true, message: '密码不可为空'},
                            {whitespace: true, message: '密码不能包含空白字符'},
                        ]
                    })(
                        <Input.Password type='password' id='pswd' name='pswd' />
                    )}
                </Form.Item>

                <Form.Item label='确认密码' required={true} style={{marginBottom: 8}}>
                    {getFieldDecorator('cfmPswd', {
                        rules: [
                            {required: true, message: '请确认密码'},
                            {whitespace: true, message: '密码不能包含空白字符'},
                            {validator: validateNewPswd}
                        ]
                    })(
                        <Input.Password type='password' id='cfmPswd' name='cfmPswd' />
                    )}
                </Form.Item>

                <Form.Item>
                    {getFieldDecorator("isSu", {
                        initialValue: '0'
                    })(
                        <Input hidden={true} id='isSu' name='isSu'/>
                    )}

                    <Checkbox onChange={ev => setFieldsValue({isSu: ev.target.checked ? "1" : "0"})} /> 拥有超级管理员权限
                </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit">添加</Button>
                </Form.Item>
            </Form>
        );
    });

    /**
     * 编辑用户
     */
    interface IFormEditUserProps { user: IUser; form: WrappedFormUtils }
    const FormEditUser = Form.create<IFormEditUserProps>()((props: IFormEditUserProps) => {
        const {getFieldDecorator, setFieldsValue, validateFields} = props.form;

        const submit = (ev: React.FormEvent<HTMLFormElement>) => {
            ev.preventDefault();

            validateFields(err => {
                if (!err) {
                    setModal(null);
                    Fetch.put(`/admin/user/${props.user.id}`, new FormData(ev.currentTarget), rsp => {
                        rsp.err ? message.error(rsp.err, 1) : fetchUsers();
                    });
                }
            })
        }

        return (
            <Form onSubmit={submit}>
                <Form.Item label='登录帐号' required={true} style={{marginBottom: 8}}>
                    {getFieldDecorator('account', {
                        rules: [
                            {required: true, message: '登录帐号未填写'},
                            {whitespace: true, message: '帐号不能包含空白字符'},
                            {max: 32, message: '最大32个字符'},
                        ],
                        initialValue: props.user.account,
                    })(
                        <Input id='account' name='account'/>
                    )}
                </Form.Item>

                <Form.Item label='显示名称' required={true} style={{marginBottom: 8}}>
                    {getFieldDecorator('name', {
                        rules: [
                            {required: true, message: '显示名称未填写'},
                            {whitespace: true, message: '名称不能包含空白字符'},
                            {max: 32, message: '最大32个字符'},
                        ],
                        initialValue: props.user.name,
                    })(
                        <Input id='name' name='name'/>
                    )}
                </Form.Item>

                <Form.Item>
                    {getFieldDecorator("isSu", {
                        initialValue: props.user.isSu ? '1' : '0'
                    })(
                        <Input hidden={true} id='isSu' name='isSu'/>
                    )}
                    
                    <Checkbox onChange={ev => setFieldsValue({isSu: ev.target.checked ? "1" : "0"})} defaultChecked={props.user.isSu} /> 拥有超级管理员权限
                </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit">修改</Button>
                </Form.Item>
            </Form>
        );
    });

    /**
     * 锁定或解禁用户
     */
    const lockUser = (user: IUser) => {
        Fetch.patch(`/admin/user/${user.id}/lock`, null, rsp => {
            rsp.err ? message.error(rsp.err) : fetchUsers()
        });
    };

    /**
     * 删除用户
     */
    const delUser = (user: IUser) => {
        Fetch.delete(`/admin/user/${user.id}`, rsp => {
            rsp.err ? message.error(rsp.err) : fetchUsers();
        });
    };

    /**
     * 添加项目
     */
    const FormAddProject = Form.create()((props: {form: WrappedFormUtils}) => {
        const { getFieldDecorator, setFieldsValue, validateFields } = props.form;

        const submit = (ev: React.FormEvent<HTMLFormElement>) => {
            ev.preventDefault();
            
            validateFields(err => {
                if (!err) {
                    setModal(null);
                    Fetch.post('/admin/project', new FormData(ev.currentTarget), rsp => {
                        rsp.err ? message.error(rsp.err, 1) : fetchProjs();
                    })
                }
            })
        }

        return (
            <Form onSubmit={submit}>
                <Form.Item label='项目名'>
                    {getFieldDecorator('name', {
                        rules: [
                            {required: true, message: '项目名称不可为空'},
                        ]
                    })(
                        <Input id='name' name='name'/>
                    )}
                </Form.Item>

                <Form.Item label='初始管理员'>
                    {getFieldDecorator('admin', {
                        rules: [
                            {required: true, message: '需要指定一个默认的管理员'},
                        ]
                    })(
                        <Input id='admin' name='admin' hidden/>
                    )}

                    <Select onChange={ev => setFieldsValue({admin: ev.valueOf()})}>
                        {users.map(user => {
                            return <Select.Option key={user.id} value={user.id}>{user.name}</Select.Option>
                        })}
                    </Select>
                </Form.Item>

                <Form.Item label='初始管理员角色'>
                    {getFieldDecorator('role', {
                        rules: [
                            {required: true, message: '这个还没填写呢'},
                        ]
                    })(
                        <Input id='role' name='role' hidden/>
                    )}

                    <Select onChange={ev => setFieldsValue({role: ev.valueOf()})}>
                        {ProjectRole.map((role, idx) => {
                            return <Select.Option key={idx} value={idx}>{role}</Select.Option>
                        })}
                    </Select>
                </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit">添加</Button>
                </Form.Item>
            </Form>
        );
    });

    /**
     * 编辑项目
     */
    interface IFormEditProjectProps { proj: IProject, form: WrappedFormUtils }
    const FormEditProject = Form.create<IFormEditProjectProps>()((props: IFormEditProjectProps) => {
        const {getFieldDecorator, validateFields} = props.form;

        const submit = (ev: React.FormEvent<HTMLFormElement>) => {
            ev.preventDefault();

            validateFields(err => {
                if (!err) {
                    setModal(null);
                    Fetch.put(`/admin/project/${props.proj.id}`, new FormData(ev.currentTarget), rsp => {
                        rsp.err ? message.error(rsp.err, 1) : fetchProjs();
                    })
                }
            })
        }

        return (
            <Form onSubmit={submit}>
                <Form.Item label='项目名称' required={true}>
                    {getFieldDecorator('name', {
                        rules: [
                            {required: true, message: '项目名称不可为空'}
                        ],
                        initialValue: props.proj.name,
                    })(
                        <Input id='name' name='name'/>
                    )}
                </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit">修改</Button>
                </Form.Item>
            </Form>
        )
    });

    /**
     * 删除项目
     */
    const delProject = (proj: IProject) => {
        Fetch.delete(`/admin/project/${proj.id}`, rsp => {
            rsp.err ? message.error(rsp.err, 1) : fetchProjs();
        });
    };

    return (
        <Layout style={{width: '100%', height: '100%'}}>
            <Layout.Content style={{padding: 32}}>
                <Row gutter={16}>
                    <Col span={16}>
                        <Card title={<div><Icon type='idcard' style={{margin: '0 8px'}}/>用户管理</div>} extra={<Button type='link' onClick={() => showModal('添加用户', <FormAddUser ref={modalForm}/>)}><Icon type='plus'/>添加用户</Button>}>
                            <Table size='middle' columns={userSchema} dataSource={users} bordered/>
                        </Card>
                    </Col>
                    <Col span={8}>
                        <Card title={<div><Icon type='pie-chart' style={{margin: '0 8px'}}/>项目管理</div>} extra={<Button type='link' onClick={() => showModal('添加项目', <FormAddProject ref={modalForm}/>)}><Icon type='plus'/>添加项目</Button>}>
                            <Table size='middle' columns={projectSchema} dataSource={projs} bordered/>
                        </Card>
                    </Col>    
                </Row>

                <div>
                    {modal}
                </div>
            </Layout.Content>
        </Layout>
    )
};