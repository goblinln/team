import * as React from 'react';
import * as moment from 'moment';

import {Badge, Button, Col, Form, Icon, Input, Markdown, Row, Notification, Empty} from '../../components';
import {TaskTag, TaskWeight, ProjectRole} from '../../common/consts';
import {Project} from '../../common/protocol';
import {CreateTaskHistory} from '../../common/storage';
import {request} from '../../common/request';

export const Creator = (props: {onDone: () => void}) => {
    const [canModifyCreator, setCanModifyCreator] = React.useState<boolean>(false);
    const [projs, setProjs] = React.useState<Project[]>([]);
    const [selectedProj, setSelectedProj] = React.useState<Project>(null);
    const [history, setHistory] = React.useState<CreateTaskHistory>(new CreateTaskHistory());
    const [initCreator, setInitCreator] = React.useState<number>(undefined);
    const [initDeveloper, setInitDeveloper] = React.useState<number>(undefined);
    const [initTester, setInitTester] = React.useState<number>(undefined);

    React.useEffect(() => {
        request({
            url: '/api/project/mine',
            success: (data: Project[]) => {
                data.forEach(p => {
                    p.members.sort((a, b) => {
                        if (a.role != b.role) {
                            return a.role - b.role; 
                        } else {
                            return a.user.account.localeCompare(b.user.account);
                        }
                    })
                });

                setProjs(data);

                let proj = data[0];
                if (proj) {
                    setHistory(CreateTaskHistory.load(proj));
                    setSelectedProj(proj);
                }
            }
        });
    }, []);

    const form = Form.useForm({
        name: {required: '任务标题不可为空', length: {min: 2, max: 64, message: '任务名在2到64个字符之间'}},
        proj: {required: '请指定所属项目'},
        branch: {required: '请指定所属分支'},
        weight: {required: '请指定优先级'},
        developer: {required: '未指定开发人员'},
        tester: {required: '请指定测试人员'},
        startTime: {required: '请选择开始时间'},
        endTime: {required: '请选择截止时间'},
        content: {required: '描述不可为空'},
    });

    const handleProjectChange = (ev: React.ChangeEvent<HTMLSelectElement>) => {
        let idx = ev.currentTarget.selectedIndex;
        let proj = projs[idx];
        let his = CreateTaskHistory.load(proj)
        setHistory(his);
        setInitCreator(his.creators.length > 0?his.creators[0].user.id:undefined);
        setInitDeveloper(his.developers.length > 0?his.developers[0].user.id:undefined);
        setInitTester(his.testers.length > 0?his.testers[0].user.id:undefined);
        setSelectedProj(proj);
    };

    const handleSubmit = (ev: React.FormEvent<HTMLFormElement>) => {
        ev.preventDefault();

        let param = new FormData(ev.currentTarget);
        let d: number = parseInt(param.get('developer') as string);
        let t: number = parseInt(param.get('tester') as string);
        let c: number = undefined;
        if (canModifyCreator) {
            c = parseInt(param.get('creator') as string);
        }
        history.save(d, t, c);

        request({
            url: '/api/task',
            method: 'POST',
            data: param,
            success: () => {
                Notification.alert('发布任务成功', 'info');
                props.onDone();
            }
        })
    };

    const uploadForMarkdown = (file: File, done: (url: string) => void) => {
        let param = new FormData();
        param.append('img', file, file.name);        
        request({url: '/api/file/upload', method: 'POST', data: param, success: (data: any) => done(data.url)});
    };

    return projs.length == 0?<div className='mt-4'><Empty label='您还未加入项目，无法发布任务'/></div>:(
        <Form form={form} onSubmit={handleSubmit}>
            <Row space={8}>
                <Col span={{xs: 4}}>
                    <Form.Field htmlFor='name' label='标题'>
                        <Input name='name' autoComplete='off'/>
                    </Form.Field>
                </Col>
                <Col span={{xs: 2}}>
                    <Form.Field htmlFor='proj' label='项目'>
                        <Input.Select name='proj' onChange={handleProjectChange}>
                            {projs.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
                        </Input.Select>
                    </Form.Field>
                </Col>
                <Col span={{xs: 2}}>
                    <Form.Field htmlFor='branch' label='分支'>
                        <Input.Select name='branch'>
                            {selectedProj&&selectedProj.branches.map((b, i) => <option key={i} value={i}>{b}</option>)}
                        </Input.Select>
                    </Form.Field>
                </Col>
                <Col span={{xs: 2}}>
                    <Form.Field htmlFor='weight' label='优先级'>
                        <Input.Select name='weight'>
                            {TaskWeight.map((w, i) => <option key={i} value={i} style={{color: w.color}}>{w.name}</option>)}
                        </Input.Select>
                    </Form.Field>
                </Col>
            </Row>

            <Row space={8}>
                <Col span={{xs: 3}}>
                    <Form.Field 
                        htmlFor='creator' 
                        label={
                            <Input.Checkbox 
                                style={{fontSize: 12}} 
                                label='指定负责人' 
                                value='1' 
                                checked={canModifyCreator} 
                                onChange={b => setCanModifyCreator(b)}/>
                        }>
                        <Input.Select name='creator' value={initCreator} disabled={!canModifyCreator}>
                            <optgroup label='最近选择'>
                                {history.creators.map(m => <option key={m.user.id} value={m.user.id}>【{ProjectRole[m.role]}】{m.user.name}</option>)}
                            </optgroup>
                            <optgroup label='其他人员'>
                                {selectedProj&&selectedProj.members.map(m => {
                                    if (history.creators.indexOf(m) == -1) {
                                        return <option key={m.user.id} value={m.user.id}>【{ProjectRole[m.role]}】{m.user.name}</option>
                                    } else {
                                        return null;
                                    }
                                })}
                            </optgroup>
                        </Input.Select>
                    </Form.Field>
                </Col>

                <Col span={{xs: 3}}>
                    <Form.Field htmlFor='developer' label='开发人员'>
                        <Input.Select name='developer' value={initDeveloper}>
                            <optgroup label='最近选择'>
                                {history.developers.map(m => <option key={m.user.id} value={m.user.id}>【{ProjectRole[m.role]}】{m.user.name}</option>)}
                            </optgroup>
                            <optgroup label='其他人员'>
                                {selectedProj&&selectedProj.members.map(m => {
                                    if (history.developers.indexOf(m) == -1) {
                                        return <option key={m.user.id} value={m.user.id}>【{ProjectRole[m.role]}】{m.user.name}</option>
                                    } else {
                                        return null;
                                    }
                                })}
                            </optgroup>
                        </Input.Select>
                    </Form.Field>
                </Col>

                <Col span={{xs: 3}}>
                    <Form.Field htmlFor='tester' label='测试人员/验收人员'>
                        <Input.Select name='tester' value={initTester}>
                            <optgroup label='最近选择'>
                                {history.testers.map(m => <option key={m.user.id} value={m.user.id}>【{ProjectRole[m.role]}】{m.user.name}</option>)}
                            </optgroup>
                            <optgroup label='其他人员'>
                                {selectedProj&&selectedProj.members.map(m => {
                                    if (history.testers.indexOf(m) == -1) {
                                        return <option key={m.user.id} value={m.user.id}>【{ProjectRole[m.role]}】{m.user.name}</option>
                                    } else {
                                        return null;
                                    }
                                })}
                            </optgroup>
                        </Input.Select>
                    </Form.Field>
                </Col>
            </Row>

            <Row space={8}>
                <Col span={{xs: 2}}>
                    <Form.Field htmlFor='startTime' label='计划开始时间'>
                        <Input.DatePicker name='startTime' mode='date' value={moment().add(1, 'd').format('YYYY-MM-DD')}/>
                    </Form.Field>
                </Col>

                <Col span={{xs: 2}}>
                    <Form.Field htmlFor='endTime' label='计划截止时间'>
                        <Input.DatePicker name='endTime' mode='date' value={moment().add(1, 'd').format('YYYY-MM-DD')}/>
                    </Form.Field>
                </Col>

                <Col span={{xs: 8}}>
                    <Form.Field htmlFor='tags[]' label='任务标签'>
                        {TaskTag.map((t, i) => (
                            <Input.Checkbox name='tags[]' value={`${i}`} label={<Badge theme={t.theme as any}>{t.name}</Badge>}/>
                        ))}
                    </Form.Field>
                </Col>
            </Row>

            <Form.Field htmlFor='content' label='任务描述'>
                <Markdown.Editor name='content' rows={10} onUpload={uploadForMarkdown}/>
            </Form.Field>

            <Form.Field>
                <Input.Uploader name='files[]'>
                    <Button theme='link' size='sm' className='p-0'><Icon type='upload' className='mr-1'/>添加附件</Button>
                </Input.Uploader>
            </Form.Field>

            <Button theme='primary' onClick={ev => {ev.preventDefault(); form.submit()}}><Icon type='plus' className='mr-1'/>发布任务</Button>
        </Form>
    );
};