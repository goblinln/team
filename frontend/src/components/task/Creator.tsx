import * as React from 'react';
import * as moment from 'moment';

import {
    Button,
    Checkbox,
    Col,
    DatePicker,
    Empty,
    Form,
    Icon,
    Input,
    Row,
    Select,
    Tag,
    Upload,
    message,
} from 'antd';

import { TaskWeight, TaskTag, ProjectRole } from '../../common/Consts';
import { IProject, IProjectMember } from '../../common/Protocol';
import { Fetch } from '../../common/Request';
import * as Markdown from '../markdown/Markdown';

/**
 * 任务创建表单的可配置属性
 */
export interface IProps {
    /**
     * 任务发布完成后的动作
     */
    onFinish: () => void;

    /**
     * 使用Form.create接口自动填充的属性
     */
    form: any,
}

/**
 * 历史记录
 */
class SelectHistory {
    proj: IProject;
    c: number[];
    d: number[];
    t: number[];

    constructor() {
        this.proj = null;
        this.c = [];
        this.d = [];
        this.t = [];
    }

    static load(proj: IProject) {
        let record = JSON.parse(window.localStorage.getItem(`creator_history_${proj.id}`));
        let ret = new SelectHistory();
        ret.proj = proj;

        if (record) {
            ret.c = record.c;
            ret.d = record.d;
            ret.t = record.t;
        }

        return ret;
    }

    save() {
        window.localStorage.setItem(`creator_history_${this.proj.id}`, JSON.stringify({
            c: this.c,
            d: this.d,
            t: this.t,
        }));
    }

    get(group: string) {
        switch (group) {
        case 'creator':
            return this.safeGet(this.c)
        case 'developer':
            return this.safeGet(this.d)
        case 'tester':
            return this.safeGet(this.t)
        }
    }

    getFirst(group: string) {
        switch (group) {
        case 'creator':
            return this.safeGetFirst(this.c)
        case 'developer':
            return this.safeGetFirst(this.d)
        case 'tester':
            return this.safeGetFirst(this.t)
        }
    }

    safeGet(desire: number[]) {
        let ret: IProjectMember[] = [];

        if (this.proj) {
            this.proj.members.forEach(m => {
                let idx = desire.indexOf(m.user.id);
                if (idx >= 0) ret.push(m);
            });
        }

        return ret;
    }

    safeGetFirst(desire: number[]) {
        let ret = this.safeGet(desire);
        if (ret.length > 0) return ret[0].user.id;
        return this.proj.members[0].user.id;
    }

    push(group: string, id: number) {
        switch (group) {
        case 'creator':
            this.c = this.safePush(this.c, id);
            break;
        case 'developer':
            this.d = this.safePush(this.d, id);
            break;
        case 'tester':
            this.t = this.safePush(this.t, id);
            break;
        }
    }

    safePush(old: number[], id: number) {
        let idx = old.indexOf(id);

        if (idx >= 0) {
            old.splice(idx, 1)
            old = [id, ...old]
        } else {
            old = [id, ...old]
        }

        return old.length > 3 ? old.slice(0, 3) : old;
    }
}

/**
 * 发布任务子页
 */
export const Creator = Form.create<IProps>()((props: IProps) => {
    const {getFieldDecorator, getFieldValue, setFieldsValue, resetFields, validateFields} = props.form;

    /**
     * 状态列表
     */
    const [isModifyCreator, setIsModifyCreator] = React.useState<boolean>(false);
    const [attachments, setAttachments] = React.useState<File[]>([]);
    const [projs, setProjs] = React.useState<IProject[]>([]);
    const [selectedProj, setSelectedProj] = React.useState<IProject>(null);
    const [history, setHistory] = React.useState<SelectHistory>(new SelectHistory());

    /**
     * 取项目列表
     */
    React.useEffect(() => {
        Fetch.get(`/api/project/mine`, rsp => {
            if (rsp.err) {
                message.error(rsp.err, 1)
            } else {
                rsp.data.forEach((proj: IProject) => {
                    proj.members = proj.members.sort((a, b) => {
                        if (a.role != b.role) {
                            return a.role - b.role;
                        } else {
                            return a.user.account.localeCompare(b.user.account);
                        }
                    })
                })

                setProjs(rsp.data);
            }
        });
    }, []);

    /**
     * 选择项目一系列操作
     */
    const selProject = (pid: any) => {
        setFieldsValue({proj: pid});

        for (let i = 0; i < projs.length; ++i) {
            if (projs[i].id == pid) {
                let last = SelectHistory.load(projs[i]);
                let lastCreator = last.getFirst('creator');
                let lastDeveloper = last.getFirst('developer');
                let lastTester = last.getFirst('tester');

                setHistory(last);
                setSelectedProj(projs[i]);

                setFieldsValue({
                    branch: 0,
                    branch_mask: 0,
                    creator: isModifyCreator ? lastCreator : null,
                    creator_mask: isModifyCreator ? lastCreator : null,
                    developer: lastDeveloper,
                    developer_mask: lastDeveloper,
                    tester: lastTester,
                    tester_mask: lastTester,
                });
                return;
            }
        }

        resetFields(['branch', 'branch_mask', 'creator', 'creator_mask', 'developer', 'developer_mask', 'tester', 'tester_mask']);
        setHistory(new SelectHistory());
        setSelectedProj(null);
    }

    /**
     * 添加附件
     */
    const addAttachment = (file: File) => {
        setAttachments(old => [...old, file]);
        return false;
    }

    /**
     * 删除附件
     */
    const delAttachment = (file: File) => {
        setAttachments(old => old.slice().splice(old.indexOf(file), 1));
        return true;
    }

    /**
     * 编写Markdown上传图片
     */
    const uploadImage = (img: File, done: (url: string) => void) => {
        let param = new FormData();
        param.append('img', img, img.name);

        Fetch.post('/api/file/upload', param, rsp => {
            rsp.err ? message.error(rsp.err, 1) : done(rsp.data.url);
        });
    };

    /**
     * 创建
     */
    const submit = (ev: React.FormEvent<HTMLFormElement>) => {
        ev.preventDefault();
        validateFields((err: any) => {
            if (!err) {
                let param = new FormData(ev.currentTarget);
                attachments.forEach(file => {param.append('files[]', file, file.name)});

                let tags: number[] = getFieldValue('tags[]') || []
                tags.forEach(tag => {param.append('tags[]', tag.toString())});

                let record = SelectHistory.load(selectedProj)
                let creator = getFieldValue('creator')
                if (creator) record.push('creator', creator)
                record.push('developer', getFieldValue('developer'))
                record.push('tester', getFieldValue('tester'))
                record.save()

                Fetch.post('/api/task', param, rsp => {
                    rsp.err ? message.error(rsp.err, 1) : props.onFinish();
                });
            }
        });
    };

    return (projs.length == 0 ? <Empty style={{marginTop: "10%"}} description="您还未加入任何项目，无法创建任务"/> :
        <Form style={{padding: 16}} onSubmit={submit}>
                <Row>
                    <Col span={6} style={{padding: '0px 2px'}}>
                        <Form.Item label='标题' style={{marginBottom: 8}} required={true}>
                            {getFieldDecorator('name', {
                                rules: [{required: true, message: '必须指定任务标题'}, {max: 64, message: '最大64个字符'}]
                            })(<Input id='name' name='name'/>)}
                        </Form.Item>
                    </Col>

                    <Col span={4} style={{padding: '0px 2px'}}>
                        <Form.Item label='项目' style={{marginBottom: 8}} required={true}>
                            {getFieldDecorator('proj', {
                                rules: [{required: true, message: '请指定所属项目'}],
                            })(<Input id='proj' name='proj' hidden={true}/>)}
                            
                            <Select id='proj_mask' defaultValue={null} onChange={(ev: number) => { selProject(ev) }}>
                                {projs.map(proj => {
                                    return <Select.Option key={proj.id} value={proj.id}>{proj.name}</Select.Option>
                                })}
                            </Select>
                        </Form.Item>
                    </Col>

                    <Col span={4} style={{padding: '0px 2px'}}>
                        <Form.Item label='分支' style={{marginBottom: 8}} required={true}>
                            {getFieldDecorator('branch', {
                                rules: [{required: true, message: '请指定所属分支'}],
                            })(<Input id='branch' name='branch' hidden={true}/>)}

                            {getFieldDecorator('branch_mask', {})(
                                <Select id='branch_mask' onChange={(ev: number) => { setFieldsValue({branch: ev}) }}>
                                    {selectedProj && selectedProj.branches.map((branch, idx) => {
                                        return <Select.Option key={idx} value={idx}>{branch}</Select.Option>
                                    })}
                                </Select>
                            )}                            
                        </Form.Item>
                    </Col>

                    <Col span={4} style={{padding: '0px 2px'}}>
                        <Form.Item label='优先级' style={{marginBottom: 8}} required={true}>
                            {getFieldDecorator('weight', {
                                rules: [{required: true, message: '请指定优先级'}],
                                initialValue: 0,
                            })(<Input id='weight' name='weight' hidden={true}/>)}

                            <Select defaultValue={0} onChange={(ev: number) => { setFieldsValue({weight: ev.valueOf()}) }}>
                                {TaskWeight.map((weight, idx) => {
                                    return <Select.Option key={idx} value={idx}><span style={{color: weight.color}}>{weight.name}</span></Select.Option>
                                })}
                            </Select>
                        </Form.Item>
                    </Col>
                </Row>

                <Row>
                    <Col span={6} style={{padding: '0 2px'}}>
                        <Form.Item style={{marginBottom: 8}} label={<Checkbox onChange={() => setIsModifyCreator(old => !old)}>指定负责人</Checkbox>}>
                            {getFieldDecorator('creator', {})(<Input id='creator' name='creator' hidden={true}/>)}
                            {getFieldDecorator('creator_mask', {})(
                                <Select id='creator_mask' onChange={(ev: number) => { setFieldsValue({creator: ev}) }} disabled={!isModifyCreator}>
                                    <Select.OptGroup label='最近选择'>
                                        {history.get('creator').map(m => {
                                            return <Select.Option key={m.user.id} value={m.user.id}>【{ProjectRole[m.role]}】{m.user.name}</Select.Option>
                                        })}
                                    </Select.OptGroup>
                                    <Select.OptGroup label='全部人员'>
                                        {selectedProj && selectedProj.members.map(member => {
                                            return <Select.Option key={member.user.id} value={member.user.id}>【{ProjectRole[member.role]}】{member.user.name}</Select.Option>
                                        })}
                                    </Select.OptGroup>                                    
                                </Select>
                            )}                            
                        </Form.Item>
                    </Col>

                    <Col span={6} style={{padding: '0 2px'}}>
                        <Form.Item label='开发人员' style={{marginBottom: 8}} required={true}>
                            {getFieldDecorator('developer', {
                                rules: [{required: true, message: '未指定开发人员'}]
                            })(<Input id='developer' name='developer' hidden={true}/>)}

                            {getFieldDecorator('developer_mask', {})(
                                <Select id='developer_mask' onChange={(ev: number) => { setFieldsValue({developer: ev}) }}>
                                    <Select.OptGroup label='最近选择'>
                                        {history.get('developer').map(m => {
                                            return <Select.Option key={m.user.id} value={m.user.id}>【{ProjectRole[m.role]}】{m.user.name}</Select.Option>
                                        })}
                                    </Select.OptGroup>
                                    <Select.OptGroup label='全部人员'>
                                        {selectedProj && selectedProj.members.map(member => {
                                            return <Select.Option key={member.user.id} value={member.user.id}>【{ProjectRole[member.role]}】{member.user.name}</Select.Option>
                                        })}
                                    </Select.OptGroup>
                                </Select>
                            )}
                        </Form.Item>
                    </Col>

                    <Col span={6} style={{padding: '0 2px'}}>
                        <Form.Item label='测试人员/验收人员' style={{marginBottom: 8}} required={true}>
                            {getFieldDecorator('tester', {
                                rules: [{required: true, message: '请指定测试人员'}]
                            })(<Input id='tester' name='tester' hidden={true}/>)}

                            {getFieldDecorator('tester_mask', {})(
                                <Select id='tester_mask' onChange={(ev: number) => { setFieldsValue({tester: ev}) }}>
                                    <Select.OptGroup label='最近选择'>
                                        {history.get('tester').map(m => {
                                            return <Select.Option key={m.user.id} value={m.user.id}>【{ProjectRole[m.role]}】{m.user.name}</Select.Option>
                                        })}
                                    </Select.OptGroup>
                                    <Select.OptGroup label='全部人员'>
                                        {selectedProj && selectedProj.members.map(member => {
                                            return <Select.Option key={member.user.id} value={member.user.id}>【{ProjectRole[member.role]}】{member.user.name}</Select.Option>
                                        })}
                                    </Select.OptGroup>
                                </Select>
                            )}
                        </Form.Item>
                    </Col>
                </Row>

                <Row>
                    <Col span={4} style={{padding: '0 2px'}}>
                        <Form.Item label='计划开始时间' style={{marginBottom: 8}} required={true}>
                            {getFieldDecorator('startTime', {
                                rules: [{required: true, message: '请选择开始时间'}],
                                initialValue: moment().add(1, 'day')
                            })(<DatePicker id='startTime' name='startTime'/>)}
                        </Form.Item>
                    </Col>

                    <Col span={4} style={{padding: '0 2px'}}>
                        <Form.Item label='计划截止时间' style={{marginBottom: 8}} required={true}>
                            {getFieldDecorator('endTime', {
                                rules: [{required: true, message: '请选择截止时间'}],
                                initialValue: moment().add(1, 'day')
                            })(<DatePicker id='endTime' name='endTime'/>)}                           
                        </Form.Item>
                    </Col>

                    <Col span={16} style={{padding: '0 2px'}}>
                        <Form.Item label='任务标签' style={{marginBottom: 8}}>
                            {getFieldDecorator("tags[]", {})(
                                <Checkbox.Group>
                                    {TaskTag.map((tag, i) => {
                                        return (
                                            <Checkbox key={i} id='tags[]' name='tags[]' value={i}>
                                                <Tag color={tag.color}><span>{tag.name}</span></Tag>
                                            </Checkbox>
                                        );
                                    })}
                                </Checkbox.Group>
                            )}                            
                        </Form.Item>
                    </Col>
                </Row>

                <Row>
                    <Form.Item label='任务描述' style={{marginBottom: 8}} required={true}>
                        {getFieldDecorator('content', {
                            rules: [{required: true, message: '描述不可为空'}]
                        })(<Input.TextArea id='content' name='content' hidden={true}/>)}

                        <Markdown.Editor
                            rows={8}
                            onChange={data => setFieldsValue({content: data})}
                            onUpload={uploadImage}/>
                    </Form.Item>              
                </Row>

                <Row>
                    <Upload className='upload-list-inline' beforeUpload={addAttachment} onRemove={file => delAttachment(file.originFileObj)}>
                        <Button type='link'><Icon type='upload'/> 添加附件</Button>
                    </Upload>
                </Row>

                <Button type='primary' htmlType='submit' style={{marginTop: 16}}>发布任务</Button>
            </Form>
    );
});