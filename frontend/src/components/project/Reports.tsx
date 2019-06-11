import * as React from 'react';
import * as moment from 'moment';

import {
    Button,
    Col,
    Divider,
    Icon,
    Row,
    Tag,
    message,
} from 'antd';

import { IProject, ITask } from '../../common/Protocol';
import { Fetch } from '../../common/Request';
import * as TaskViewer from '../task/Viewer';

/**
 * 项目周报页
 */
export const Reports = (props: {proj: IProject, isReadonly: boolean}) => {

    /**
     * 状态列表
     */
    const [week, setWeek] = React.useState<moment.Moment>(moment().startOf('week'));
    const [data, setData] = React.useState<{archived: ITask[], unarchived: ITask[]}>({archived: [], unarchived: []});
    const [showOption, setShowOption] = React.useState<boolean>(false);
    const taskDetailAnchor = React.useRef<any>(null);

    /**
     * 初始化
     */
    React.useEffect(() => {
        TaskViewer.default.init(taskDetailAnchor, () => null);
    }, []);

    /**
     * 拉取任务列表
     */
    React.useEffect(() => {
        setShowOption(!props.isReadonly && week.diff(moment().startOf('week')) == 0);
        fetchReport(week);
    }, [week]);

    /**
     * 拉取周报
     */
    const fetchReport = (whichWeek?: moment.Moment) => {
        let w = whichWeek || moment().startOf('week');

        Fetch.get(`/api/project/${props.proj.id}/report/${w.unix()}`, rsp => {
            if (rsp.err) {
                message.error(rsp.err);
            } else {
                rsp.data.unarchived = rsp.data.unarchived.sort((a: ITask, b: ITask) => {
                    if (a.state != b.state) {
                        return b.state - a.state;
                    } else {
                        return moment(a.endTime).diff(moment(b.endTime), 'd');
                    }
                })
                setData(rsp.data);
            }
        })
    };

    /**
     * 验收一个任务
     */
    const archiveOne = (tid: number) => {
        Fetch.put(`/api/project/${props.proj.id}/archive/${tid}`, null, rsp => {
            rsp.err ? message.error(rsp.err, 1) : fetchReport()
        });
    }

    /**
     * 一键验收全部已完成的任务
     */
    const archiveAll = () => {
        Fetch.post(`/api/project/${props.proj.id}/archive/all`, null, rsp => {
            rsp.err ? message.error(rsp.err, 1) : fetchReport()
        });
    };

    return (
        <div style={{padding: 16}}>
            <Row type='flex' justify='center' align='middle'>
                <Icon type="left-circle" theme="filled" style={{fontSize: '2.5em'}} onClick={() => setWeek(moment(week).subtract(1, 'week'))}/>
                <div style={{margin: '0 16px', textAlign: 'center'}}>
                    <p style={{fontSize: '2em', fontWeight: 'bolder', marginBottom: 8}}>第{week.weeks()}周项目周报</p>
                    <small><Tag color='#6c757d'>{moment(week).format('YYYY/MM/DD')} - {moment(week).endOf('week').format('YYYY/MM/DD')}</Tag></small>
                </div>
                <Icon type="right-circle" theme="filled" style={{fontSize: '2.5em'}} onClick={() => setWeek(moment(week).add(1, 'week'))}/>
            </Row>

            <Row gutter={8} style={{marginTop: 16}}>
                <Col span={12}>
                    <Row type='flex' justify='space-between' align='bottom'>
                        <div style={{fontSize: '1.5em', fontWeight: 'bold'}}>
                            <Icon type='frown'/> 未验收任务
                        </div>

                        {showOption && <Button type='link' size='small' onClick={() => archiveAll()}>验收全部已完成</Button>}
                    </Row>

                    <Divider style={{margin: '8px 0', background: '#cccccc'}}/>

                    {data.unarchived.map(task => {
                        return (
                            <Row type='flex' justify='space-between' align='middle'>
                                <div style={{padding: '0 8px'}}>
                                    {task.state == 3 ? <Icon type='question-circle' style={{color: 'yellow'}}/> : <Icon type='close-circle' style={{color: 'red'}}/>}

                                    <span style={{margin: '0 8px'}}>{task.endTime}</span>

                                    <Button
                                        type='link'
                                        style={{padding: 0}}
                                        onClick={() => TaskViewer.default.open(task.id, props.isReadonly)}>{task.name}</Button>
                                </div>

                                <div style={{padding: '0 8px'}}>
                                    <span>{task.creator.name}<Icon type='right' />{task.developer.name}<Icon type='right' />{task.tester.name}</span>
                                    {showOption && task.state == 3 && <Button type='link' size='small' style={{padding: 0, paddingLeft: 8}} onClick={() => archiveOne(task.id)}>验收</Button>}
                                </div>
                            </Row>
                        );
                    })}
                </Col>

                <Col span={12}>
                    <Row type='flex' justify='space-between' align='bottom'>
                        <div style={{fontSize: '1.5em', fontWeight: 'bold'}}>
                            <Icon type='smile'/> 已验收任务
                        </div>
                    </Row>

                    <Divider style={{margin: '8px 0', background: '#cccccc'}}/>

                    {data.archived.map(task => {
                        return (
                            <Row type='flex' justify='space-between' align='middle'>
                                <div style={{padding: '0 8px'}}>
                                    <Icon type='check' style={{color: 'green', paddingRight: 8}}/>
                                    <Button
                                        type='link'
                                        style={{overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', padding: 0}}
                                        onClick={() => TaskViewer.default.open(task.id, props.isReadonly)}>{task.name}</Button>
                                </div>

                                <div style={{padding: '0 8px'}}>
                                    {task.creator.name}<Icon type='right' />{task.developer.name}<Icon type='right' />{task.tester.name}
                                </div>
                            </Row>
                        );
                    })}
                </Col>
            </Row>

            <div ref={taskDetailAnchor}/>
        </div>
    );
}