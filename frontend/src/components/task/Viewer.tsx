import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as moment from 'moment';

import {
    Affix,
    Avatar,
    Button,
    Col,
    Comment,
    DatePicker,
    Divider,
    Drawer,
    Dropdown,
    Icon,
    Mentions,
    Menu,
    Popover,
    Row,
    Select,
    Tabs,
    Tag,
    Timeline,
    Tooltip,
    message,
} from 'antd';

import { ProjectRole, TaskStatus, TaskWeight, TaskTag } from '../../common/Consts';
import { ITask, IUser } from '../../common/Protocol';
import { Fetch } from '../../common/Request';
import * as Markdown from '../markdown/Markdown';

/**
 * 任务详情查看接口
 */
export default class Viewer {

    /**
     * 初始化
     * 
     * @param anchor 锚点
     * @param onModify 技能有更改时回调
     */
    public static init(anchor: React.RefObject<any>, onModify?: (task: ITask) => void) {
        Viewer._anchor = anchor;
        Viewer._onModify = onModify;
    }

    /**
     * 显示任务详情
     * 
     * @param taskId 任务ID
     * @param isReadonly 是否只读模式打开
     */
    public static open(taskId: number, isReadonly?: boolean) {
        Fetch.get(`/api/task/${taskId}`, rsp => {
            if (rsp.err) {
                message.error(rsp.err, 1);
            } else {
                let task: ITask = rsp.data;
                ReactDOM.render(
                    isReadonly ? <ReadOnlyViewer task={task}/> : <EditableViewer task={task}/>, 
                    Viewer._anchor.current);
            }
        })      
    }

    /**
     * 关闭任务展示
     * 
     * @param task 显示的任务
     * @param isModified 是否有更改
     */
    public static close(task: ITask, isModified?: boolean) {
        ReactDOM.render(null, Viewer._anchor.current);
        if (isModified && Viewer._onModify) Viewer._onModify(task);
    }

    private static _anchor: React.RefObject<any> = null;
    private static _onModify: (task: ITask) => void = null;
}

/**
 * 公用标题展示
 */
const CommonHeader = (props: {task: ITask, titleWidth: number}) => {
    const {task} = props;

    return (
        <Row type='flex' justify='space-between' align='middle'>
            <Col style={{fontWeight: 'bold', fontSize: '1.2em', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', maxWidth: props.titleWidth}}>
                {task.name}
                <div style={{fontWeight: 'normal', fontSize: '.4em', display: 'inline', marginLeft: 4}}>
                    {task.tags.map(tag => {
                        return <Tag key={tag} color={TaskTag[tag].color}><span>{TaskTag[tag].name}</span></Tag>
                    })}
                </div>
            </Col>

            <Col style={{ fontWeight: 'normal', fontSize: '.5em' }}>
                <span style={{ marginRight: 16 }}><Icon type='pie-chart' /> {task.proj.name}</span>
                <span style={{ marginRight: 16 }}><Icon type='branches' /> {task.proj.branches[task.branch] || '默认'}</span>
                <span><Icon type={TaskStatus[task.state].icon}/> {TaskStatus[task.state].name}</span>
            </Col>
        </Row>
    );
}

/**
 * 只读模式查看面板
 */
const ReadOnlyViewer = (props: {task: ITask}) => {
    const {task} = props;

    return (
        <Drawer
            title={<CommonHeader task={task} titleWidth={500}/>}
            width={600}
            bodyStyle={{ padding: '4px 0px' }}
            visible={true}
            closable={false}
            onClose={() => { Viewer.close(task); }}>
            <Row type='flex' justify='start' style={{ padding: '4px 16px' }}>
                <Col style={{ marginRight: 16 }}><Icon type='notification' /> {task.creator.name}</Col>
                <Col style={{ marginRight: 16 }}><Icon type='code' /> {task.developer.name}</Col>
                <Col style={{ marginRight: 16 }}><Icon type='experiment' /> {task.tester.name}</Col>
                <Col style={{ marginRight: 16 }}><Icon type='calendar' /> {moment(task.startTime).format('MM月DD日')} - {moment(task.endTime).format('MM月DD日')}</Col>
                <Col style={{ marginRight: 16 }}><Icon type='tag' /> {TaskWeight[task.weight].name}</Col>
            </Row>

            <Divider style={{ marginTop: 4, marginBottom: 4 }} />

            <div style={{ padding: '16px 16px' }}>
                <Markdown.Renderer source={task.content} />
            </div>

            {task.attachments.length > 0 && (
                <Row type='flex' justify='start' align='middle' style={{ padding: '4px 16px', fontSize: '.4em' }}>
                    附件：
                    {task.attachments.map((attachment, idx) => {
                        return <Button key={idx} href={attachment.url} size='small' type='link' style={{ marginRight: 4, fontSize: '.4em' }}><Icon type='link' />{attachment.name}</Button>
                    })}
                </Row>
            )}

            <div style={{ padding: '0px 4px', marginTop: 16 }}>
                <Tabs defaultActiveKey='comments'>
                    <Tabs.TabPane key='comments' tab='评论'>
                        <div style={{ marginLeft: 8 }}>
                            {task.comments.map((comment, idx) => {
                                return (
                                    <Comment
                                        key={idx}
                                        author={comment.user}
                                        avatar={<Avatar icon='user' src={comment.avatar} />}
                                        content={<Markdown.Renderer source={comment.content} />}
                                        datetime={<Tooltip title={comment.time}><span>{moment(comment.time).fromNow()}</span></Tooltip>} />
                                );
                            })}
                        </div>
                    </Tabs.TabPane>

                    <Tabs.TabPane key='events' tab='事件'>
                        <div style={{ marginLeft: 8 }}>
                            <Timeline>
                                {task.events.map((ev, idx) => {
                                    return (
                                        <Timeline.Item style={{ padding: '0 0 4px' }} key={idx}>
                                            <small>{ev.time}</small>  <strong>{ev.operator}</strong> {ev.desc}
                                        </Timeline.Item>
                                    );
                                })}
                            </Timeline>
                        </div>
                    </Tabs.TabPane>
                </Tabs>
            </div>
        </Drawer>
    );
};

/**
 * 编辑模式预览
 */
const EditableViewer = (props: {task: ITask}) => {
    
    /**
     * 状态列表
     */
    const [task, setTask] = React.useState<ITask>(props.task);
    const [isDirty, setDirty] = React.useState<boolean>(false);
    const [needFetch, setNeedFetch] = React.useState<boolean>(false);
    const [isEditorShow, setEditorShow] = React.useState<boolean>(false);
    const [isCommentEditorShow, setCommentEditorShow] = React.useState<boolean>(false);

    /**
     * 重新加载
     */
    React.useEffect(() => {
        if (!isDirty) return;
        setNeedFetch(true);

        Fetch.get(`/api/task/${task.id}`, rsp => {
            if (rsp.err) {
                message.error(rsp.err, 1);
            } else {
                setTask(rsp.data);
                setDirty(false);
            }
        });
    }, [isDirty]);

    /**
     * 各种修改器
     */
    const MemberEditor = (props: {icon: string, current: IUser, kind: string}) => {
        const [selected, setSelected] = React.useState<IUser>(props.current);
        const titles: any = {'creator': '修改负责人', 'developer': '修改开发人员', 'tester': '修改测试人员'}
        const sorted = task.proj.members.sort((a, b) => {
            if (a.role != b.role) {
                return a.role - b.role;
            } else {
                return a.user.account.localeCompare(b.user.account);
            }
        })

        const modify = (val: number) => {
            let find: IUser = null;
            for (let i = 0; i < task.proj.members.length; ++i) {
                if (task.proj.members[i].user.id == val) {
                    find = task.proj.members[i].user;
                    break;
                }
            }

            if (!find) return;

            let param = new URLSearchParams({
                member: val.toString(),
            });

            Fetch.patch(`/api/task/${task.id}/${props.kind}`, param, rsp => {
                rsp.err ? message.error(rsp.err, 1) : setSelected(find)
            });
        }

        const onVisibleChange = (visible: boolean) => {
            setDirty(!visible && selected.id != props.current.id);
        }

        return (
            <Popover
                title={titles[props.kind]}
                trigger='hover'
                onVisibleChange={onVisibleChange}
                content={
                    <Select value={selected.id} onChange={modify} style={{minWidth: 160}}>
                        {sorted.map(member => {
                            return <Select.Option key={member.user.id} value={member.user.id}>【{ProjectRole[member.role]}】{member.user.name}</Select.Option>
                        })}
                    </Select>
                }>
                <Button type='link' size='small' style={{padding: '0 4px', color: 'rgba(0,0,0,.65)'}}><Icon type={props.icon} />{selected.name}</Button>
            </Popover>
        );
    }
    const WeightEditor = (props: {current: number}) => {
        const [w, setW] = React.useState<number>(props.current);

        const modify = (val: number) => {
            let param = new URLSearchParams({
                weight: val.toString(),
                old: TaskWeight[props.current].name,
            })
            Fetch.patch(`/api/task/${task.id}/weight`, param, rsp => {
                rsp.err ? message.error(rsp.err, 1) : setW(val);
            });
        }

        const onVisibleChange = (visible: boolean) => {
            setDirty(!visible && w != props.current);
        }

        return (
            <Dropdown
                onVisibleChange={onVisibleChange}
                overlay={
                    <Menu>
                        {TaskWeight.map((weight, idx) => {
                            return <Menu.Item key={idx} onClick={() => modify(idx)}><span style={{color: weight.color}}>{weight.name}</span></Menu.Item>
                        })}
                    </Menu>
                }>
                <Button type='link' size='small' style={{padding: '0 4px', color: 'rgba(0,0,0,.65)'}}><Icon type='tag' />{TaskWeight[w].name}</Button>
            </Dropdown>
        );
    }
    const TimeEditor = (props: {start: string, end: string}) => {
        const [startTime, setStartTime] = React.useState<string>(props.start);
        const [endTime, setEndTime] = React.useState<string>(props.end);
        const [changed, setChanged] = React.useState<boolean>(false);

        const modify = () => {
            if (startTime == task.startTime && endTime == task.endTime) {
                return;
            }

            let param = new URLSearchParams({
                startTime: startTime,
                endTime: endTime,
            })

            Fetch.patch(`/api/task/${task.id}/time`, param, rsp => {
                if (rsp.err) {
                    message.error(rsp.err, 1)
                } else {
                    setChanged(true);
                }
            });
        }

        const onVisibleChange = (visible: boolean) => {
            setDirty(!visible && changed);
        }

        return (
            <Popover
                title='修改'
                trigger='hover'
                onVisibleChange={onVisibleChange}
                content={
                    <div>
                        <Row style={{marginBottom: 8}}>开始时间</Row>
                        <Row style={{marginBottom: 8}}>
                            <DatePicker defaultValue={moment(props.start)} onChange={(date: moment.Moment, dateString: string) => setStartTime(dateString)}/>
                        </Row>

                        <Row style={{marginBottom: 8}}>截止时间</Row>
                        <Row style={{marginBottom: 8}}>
                            <DatePicker defaultValue={moment(props.end)} onChange={(date: moment.Moment, dateString: string) => setEndTime(dateString)}/>
                        </Row>

                        <Row>
                            <Button type='primary' onClick={() => modify()} block>修改时间</Button>
                        </Row>
                    </div>
                }>
                <Button type='link' size='small' style={{padding: '0 4px', color: 'rgba(0,0,0,.65)'}}>
                    <Icon type='calendar' />{moment(props.start).format('MM月DD日')} - {moment(props.end).format('MM月DD日')}
                </Button>
            </Popover>
        );
    }
    const ContentEditor = (props: {content: string}) => {
        const [content, setContent] = React.useState<string>(props.content);

        const modify = () => {
            Fetch.patch(`/api/task/${task.id}/content`, new URLSearchParams({content: content}), rsp => {
                if (rsp.err) {
                    message.error(rsp.err, 1);
                } else {
                    setEditorShow(false);
                    setDirty(content != task.content);
                }
            });
        }

        return (
            <div>
                <Row>
                    <Markdown.Editor rows={16} value={content} onChange={data => setContent(data)}/>
                </Row>
                <Row type='flex' justify='center' style={{marginTop: 8}}>
                    <Button type='primary' style={{marginRight: 8}} onClick={() => modify()}>修改</Button>
                    <Button onClick={() => setEditorShow(false)}>取消</Button>
                </Row>
            </div>
        );
    }
    const CommentEditor = () => {
        const [content, setContent] = React.useState<string>('');

        const sendComment = () => {
            Fetch.post(`/api/task/${task.id}/comment`, new URLSearchParams({content: content}), rsp => {
                if (rsp.err) {
                    message.error(rsp.err, 1);
                } else {
                    setCommentEditorShow(false);
                    setDirty(true);
                }
            });
        }

        return (
            <div>
                <Row>
                    <Mentions onChange={txt => setContent(txt)} rows={4} style={{width: 600}}>
                        {task.proj.members.map(member => {
                            return <Mentions.Option key={member.user.id} value={member.user.name}>【{ProjectRole[member.role]}】{member.user.name}</Mentions.Option>
                        })}
                    </Mentions>
                </Row>
                <Row type='flex' justify='center' style={{marginTop: 8}}>
                    <Button type='primary' style={{marginRight: 8}} onClick={() => sendComment()}>发表</Button>
                    <Button onClick={() => setCommentEditorShow(false)}>取消</Button>
                </Row>
            </div>
        );
    }

    /**
     * 任务进入下一个流程
     */
    const goNext = () => {
        Fetch.post(`/api/task/${task.id}/next`, {}, rsp => {
            if (rsp.err) {
                message.error(rsp.err, 1);
            } else {
                setDirty(true);
            }
        })
    };

    /**
     * 删除任务
     */
    const deleteTask = () => {
        Fetch.delete(`/api/task/${task.id}`, rsp => {
            rsp.err ? message.error(rsp.err, 1) : Viewer.close(task, true);
        });
    };

    return (
        <Drawer
            title={<CommonHeader task={task} titleWidth={600}/>}
            width={750}
            bodyStyle={{ padding: '4px 0px' }}
            visible={true}
            closable={false}
            onClose={() => { Viewer.close(task, needFetch); }}>

            <Row type='flex' justify='start' align='middle' style={{ padding: '4px 16px' }}>
                <Col><MemberEditor current={task.creator} kind='creator' icon='notification'/></Col>
                <Col><MemberEditor current={task.developer} kind='developer' icon='code'/></Col>
                <Col><MemberEditor current={task.tester} kind='tester' icon='experiment'/></Col>
                <Col><WeightEditor current={task.weight}/></Col>
                <Col><TimeEditor start={task.startTime} end={task.endTime}/></Col>
                <Col>
                    <Button type='link' size='small' style={{padding: '0 4px', color: 'rgba(0,0,0,.65)'}} onClick={() => setEditorShow(true)}>
                        <Icon type='edit' />编辑
                    </Button>
                </Col>
                <Col>
                    <Button type='link' size='small' style={{padding: '0 4px', color: 'rgba(0,0,0,.65)'}} onClick={() => goNext() }>
                        <Icon type="right-circle" />下一步
                    </Button>
                </Col>
                <Col>
                    <Button type='link' size='small' style={{padding: '0 4px', color: 'rgba(0,0,0,.65)'}} onClick={() => deleteTask()}>
                        <Icon type='delete' />删除
                    </Button>
                </Col>
            </Row>

            <Divider style={{ marginTop: 4, marginBottom: 4 }} />

            <div style={{ padding: '16px 16px' }}>
                {isEditorShow ? <ContentEditor content={task.content} /> : <Markdown.Renderer source={task.content} />}
            </div>

            {task.attachments.length > 0 && (
                <Row type='flex' justify='start' align='middle' style={{ padding: '4px 16px', fontSize: '.4em' }}>
                    附件：
                    {task.attachments.map((attachment, idx) => {
                        return <Button key={idx} href={attachment.url} size='small' type='link' style={{ marginRight: 4, fontSize: '.4em' }}><Icon type='link' />{attachment.name}</Button>
                    })}
                </Row>
            )}

            <div style={{ padding: '0px 4px', marginTop: 16 }}>
                <Tabs defaultActiveKey='comments'>
                    <Tabs.TabPane key='comments' tab='评论'>
                        <div style={{ marginLeft: 8 }}>
                            {task.comments.map((comment, idx) => {
                                return (
                                    <Comment
                                        key={idx}
                                        author={comment.user}
                                        avatar={<Avatar icon='user' src={comment.avatar} />}
                                        content={<Markdown.Renderer source={comment.content} />}
                                        datetime={<Tooltip title={comment.time}><span>{moment(comment.time).fromNow()}</span></Tooltip>} />
                                );
                            })}
                        </div>
                    </Tabs.TabPane>

                    <Tabs.TabPane key='events' tab='事件'>
                        <div style={{ marginLeft: 8 }}>
                            <Timeline>
                                {task.events.map((ev, idx) => {
                                    return (
                                        <Timeline.Item style={{ padding: '0 0 4px' }} key={idx}>
                                            <small>{ev.time}</small>  <strong>{ev.operator}</strong> {ev.desc}
                                        </Timeline.Item>
                                    );
                                })}
                            </Timeline>
                        </div>
                    </Tabs.TabPane>
                </Tabs>
            </div>

            <Affix style={{ position: 'absolute', bottom: 16, right: 16 }}>
                <Popover
                    content={<CommentEditor/>}
                    title='发表评论'
                    trigger='click'
                    visible={isCommentEditorShow}
                    onVisibleChange={() => setCommentEditorShow(prev => !prev)}>
                    <Button type='link'><Icon type='message' style={{ fontSize: 24 }} title='评论' /></Button>
                </Popover>
            </Affix>
        </Drawer>
    );
};