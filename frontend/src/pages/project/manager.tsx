import * as React from 'react';

import {Table, Modal, Input, Notification, Form, FormProxy, FormFieldValidator, TableColumn, Avatar, Badge, Icon, Row, Card, Button} from '../../components';
import {Project, ProjectMember, User} from '../../common/protocol';
import {request} from '../../common/request';
import { ProjectRole } from '../../common/consts';

export const Manager = (props: {pid: number}) => {
    const [proj, setProj] = React.useState<Project>();

    const memberSchema: TableColumn[] = [
        {label: '头像', renderer: (data: ProjectMember) => <Avatar size={32} src={data.user.avatar}/>},
        {label: '昵称', renderer: (data: ProjectMember) => data.user.name},
        {label: '帐号', renderer: (data: ProjectMember) => data.user.account},
        {label: '角色', renderer: (data: ProjectMember) => <Badge theme='primary'>{ProjectRole[data.role]}</Badge>},
        {label: '管理权限', renderer: (data: ProjectMember) => <Input.Switch on={data.isAdmin} disabled/>},
        {label: '操作', renderer: (data: ProjectMember) => (
            <span>
                <a className='link' onClick={() => editMember(data)}>编辑</a>
                <div className='divider-v'/>
                <a className='link' onClick={() => delMember(data)}>删除</a>
            </span>
        )}
    ];

    React.useEffect(() => fetchProject(), []);

    const fetchProject = () => {
        request({url: `/api/project/${props.pid}`, success: (data: Project) => {
            data.members.sort((a, b) => a.user.account.localeCompare(b.user.account));
            setProj(data);
        }});
    };

    const addBranch = () => {
        let branch: string = '';

        Modal.open({
            title: '添加分支',
            body: <Input className='my-2' onChange={v => branch = v}/>,
            onOk: () => {
                if (branch.length == 0) {
                    Notification.alert('分支不可为空', 'error');
                    return;
                }

                if (proj.branches.indexOf(branch) >= 0) {
                    Notification.alert('同名分支已存在', 'error');
                    return;
                }

                let param = new FormData();
                param.append('branch', branch);
                request({url: `/api/project/${props.pid}/branch`, data: param, success: fetchProject});
            }
        });
    };

    const addMember = () => {
        request({
            url: `/api/project/${props.pid}/invites`,
            success: (data: User[]) => {
                let form: FormProxy = null;
                let closer: () => void = null;

                const validate: {[k:string]:FormFieldValidator} = {
                    uid: {required: '请选择邀请的成员'},
                    role: {required: '请设置成员的职能'},
                };

                const submit = (ev: React.FormEvent<HTMLFormElement>) => {
                    ev.preventDefault();
                    request({
                        url: `/api/project/${props.pid}/member`, 
                        method: 'POST', 
                        data: new FormData(ev.currentTarget), 
                        success: () => {closer(); fetchProject()}});
                };

                closer = Modal.open({
                    title: '邀请成员',
                    body: (
                        <Form style={{width: 300}} form={() => {form = Form.useForm(validate); return form}} onSubmit={submit}>
                            <Form.Field htmlFor='uid' label='可邀请用户'>
                                <Input.Select name='uid'>
                                    {data.map(u => <option key={u.id} value={u.id}>{u.name}</option>)}
                                </Input.Select>
                            </Form.Field>
                            <Form.Field htmlFor='role' label='职能'>
                                <Input.Select name='role'>
                                    {ProjectRole.map((r, i) => <option key={i} value={i}>{r}</option>)}
                                </Input.Select>
                            </Form.Field>
                            <Form.Field>
                                <Input.Checkbox name='isAdmin' value='1' label='管理权限'/>
                            </Form.Field>
                        </Form> 
                    ),
                    onOk: () => {form.submit(); return false}
                })
            }
        });
    };

    const editMember = (m: ProjectMember) => {
        let form: FormProxy = null;
        let closer: () => void = null;

        const validate: {[k:string]:FormFieldValidator} = {
            role: {required: '请设置成员的职能'},
        };

        const submit = (ev: React.FormEvent<HTMLFormElement>) => {
            ev.preventDefault();
            request({
                url: `/api/project/${props.pid}/member/${m.user.id}`, 
                method: 'PUT', 
                data: new FormData(ev.currentTarget), 
                success: () => {closer(); fetchProject()}});
        };

        closer = Modal.open({
            title: '编辑成员',
            body: (
                <Form style={{width: 300}} form={() => {form = Form.useForm(validate); return form}} onSubmit={submit}>
                    <Form.Field htmlFor='role' label='职能'>
                        <Input.Select name='role' value={m.role}>
                            {ProjectRole.map((r, i) => <option key={i} value={i}>{r}</option>)}
                        </Input.Select>
                    </Form.Field>
                    <Form.Field>
                        <Input.Checkbox name='isAdmin' value='1' label='管理权限' checked={m.isAdmin}/>
                    </Form.Field>
                </Form> 
            ),
            onOk: () => {form.submit(); return false}
        });
    };

    const delMember = (m: ProjectMember) => {
        Modal.open({
            title: '删除确认',
            body: <div className='my-2'>确定要删除成员【{m.user.name}】吗？</div>,
            onOk: () => {
                request({
                    url: `/api/project/${props.pid}/member/${m.user.id}`, 
                    method: 'DELETE', 
                    success: fetchProject
                });
            }
        });
    };

    return (
        <div className='m-4'>
            <Card
                header={
                    <Row flex={{align: 'middle', justify: 'space-between'}}>
                        <span>
                            <Icon type='branches' className='mr-2'/>
                            分支列表
                        </span>
                        <Button theme='link' onClick={addBranch}><Icon type='plus' className='mr-1'/>添加分支</Button>
                    </Row>
                }
                bordered
                shadowed>
                <div className='p-2'>
                    {proj&&proj.branches.map((b, i) => <Badge key={i} theme='primary'>{b}</Badge>)}
                </div>
            </Card>

            <Card
                className='mt-3'
                bodyProps={{className: 'p-2'}}
                header={
                    <Row flex={{align: 'middle', justify: 'space-between'}}>
                        <span>
                            <Icon type='idcard' className='mr-2'/>
                            成员列表
                        </span>
                        <Button theme='link' onClick={addMember}><Icon type='plus' className='mr-1'/>添加成员</Button>
                    </Row>
                }
                bordered
                shadowed>
                <Table dataSource={proj?proj.members:[]} columns={memberSchema}/>
            </Card>
        </div>
    );
};