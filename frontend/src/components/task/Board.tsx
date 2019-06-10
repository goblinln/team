import * as React from 'react';
import * as moment from 'moment';

import {
    Button,
    Card,
    Col,
    Empty,
    Icon,
    Row,
    Tag,
} from 'antd';

import { ITask } from '../../common/Protocol';
import { TaskWeight, TaskStatus } from '../../common/Consts';
import * as TaskViewer from './Viewer';

/**
 * 看板视图的可配置属性
 */
export interface IBoardProps {
    /**
     * 可见任务列表
     */
    tasks: ITask[],
    /**
     * 是否是只读模式
     */
    isReadonly?: boolean;
}

/**
 * 看板视图
 */
export const Board = (props: IBoardProps) => {

    const [groups, setGroups] = React.useState<ITask[][]>([[], [], [], []]);

    React.useEffect(() => {
        let result: ITask[][] = [[], [], [], []];
        props.tasks.forEach(task => result[task.state].push(task));
        setGroups(result);
    }, [props.tasks]);

    return (
        <Row type='flex' justify='space-between' style={{margin: 8}}>
            {groups.map((group, idx) => {
                const state = TaskStatus[idx];
                
                return (
                    <Col key={state.type} span={6}>
                        <Card
                            title={<div style={{ color: 'white' }}><Icon type={state.icon} style={{marginRight: 4}} />{state.name}</div>}
                            extra={<Tag color='#fff'><label style={{ color: 'black' }}>{group.length}</label></Tag>}
                            headStyle={{backgroundColor: state.color, height: 40}}
                            bodyStyle={{padding: 0, margin: 4}}
                            style={{margin: 4}}>
                            
                            {group.length == 0 ? <Empty description='暂无数据' /> : group.map(task => {
                                const weight = TaskWeight[task.weight];
                                const now = moment();
                                const endTime = moment(task.endTime);
                                const border = task.bringTop ? 'purple' : (endTime.diff(now) <= 0 ? 'red' : 'gray');

                                return (
                                    <Card key={task.id} bodyStyle={{padding: 8, borderLeft: '4px solid ' + border}} style={{marginBottom: 4}}>
                                        <Row type='flex' justify='space-between' style={{fontSize: '.6em'}}>
                                            <Col><Icon type='pie-chart' /> {task.proj.name}</Col>
                                            <Col><Icon type='branches' /> {task.proj.branches[task.branch] || '默认'}</Col>
                                        </Row>

                                        <Row>
                                            <Button
                                                type='link'
                                                style={{textAlign: 'left', margin: 0, padding: 0, fontSize: '1em', fontWeight: 'bold', textOverflow: 'ellipsis', whiteSpace: 'nowrap', overflow: 'hidden' }}
                                                onClick={() => TaskViewer.default.open(task.id, props.isReadonly)}
                                                block>
                                                <label style={{ color: weight.color }}>{weight.name}</label>{task.name}
                                            </Button>
                                        </Row>

                                        <Row type='flex' justify='space-between' style={{fontSize: '.4em'}}>
                                            <Col><Icon type="calendar" /> {task.endTime}</Col>
                                            <Col>{task.creator.name}<Icon type='right' />{task.developer.name}<Icon type='right' />{task.tester.name}</Col>
                                        </Row>
                                    </Card>
                                );
                            })}
                        </Card>
                    </Col>
                );
            })}
        </Row>
    );
};