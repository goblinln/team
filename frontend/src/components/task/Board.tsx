import * as React from 'react';
import * as moment from 'moment';

import {
    Button,
    Card,
    Col,
    Dropdown,
    Empty,
    Icon,
    Menu,
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
    interface ISorter { name: string; exec: (a:ITask, b:ITask) => number; }
    interface IGroup { sorter: number, tasks: ITask[] }

    /**
     * 状态列表
     */
    const [sortMethod, setSortMethod] = React.useState<number[]>([0, 0, 0, 0]);
    const [groups, setGroups] = React.useState<IGroup[]>([]);

    /**
     * 排序规则
     */
    const sorters: ISorter[] = [
        {
            name: '默认排序',
            exec: (a, b) => {
                if (a.bringTop != b.bringTop) {
                    return a.bringTop ? -1 : 1;
                } else if (a.weight != b.weight) {
                    return b.weight - a.weight;
                } else {
                    moment(a.endTime).diff(moment(b.endTime), 'd');
                }
            }
        },
        {
            name: '按发布时间',
            exec: (a, b) => {
                let offset = moment(a.startTime).diff(moment(b.startTime))
                if (offset != 0) return offset;
                
                if (a.weight != b.weight) {
                    return b.weight - a.weight;
                } else if (a.bringTop != b.bringTop) {
                    return a.bringTop ? -1 : 1;
                } else {
                    return 0;
                }
            }
        },
        {
            name: '按截止时间',
            exec: (a, b) => {
                let offset = moment(a.endTime).diff(moment(b.endTime))
                if (offset != 0) return offset;
                
                if (a.weight != b.weight) {
                    return b.weight - a.weight;
                } else if (a.bringTop != b.bringTop) {
                    return a.bringTop ? -1 : 1;
                } else {
                    return 0;
                }
            }
        },
    ];

    /**
     * 任务列表变化时重新分组
     */
    React.useEffect(() => {
        let result: IGroup[] = [
            { sorter: sortMethod[0], tasks: [] },
            { sorter: sortMethod[1], tasks: [] },
            { sorter: sortMethod[2], tasks: [] },
            { sorter: sortMethod[3], tasks: [] },
        ];

        props.tasks.forEach(task => result[task.state].tasks.push(task));
        result.forEach(group => group.tasks = group.tasks.sort(sorters[group.sorter].exec));

        setGroups(result);
    }, [props.tasks]);

    /**
     * 排序规则发生变化时重新排序
     */
    React.useEffect(() => {
        let targets = groups.slice();

        for (let i = 0; i < 4; ++i) {
            let group = targets[i];
            if (!group) continue;

            if (group.sorter != sortMethod[i]) {
                group.sorter = sortMethod[i];
                group.tasks = group.tasks.sort(sorters[group.sorter].exec);
            }
        }

        setGroups(targets);
    }, [sortMethod]);

    /**
     * 修改排序规则
     */
    const changeSortMethod = (idx: number, method: number) => {
        setSortMethod(prev => {
            let old = prev.slice();
            old[idx] = method;
            return old;
        })
    };

    /**
     * 组表头
     */
    const makeGroupExtra = (group: IGroup, idx: number) => {
        return (
            <div>
                <Dropdown
                    overlay={
                        <Menu>
                            {sorters.map((sorter, method) => {
                                return <Menu.Item key={idx} onClick={() => changeSortMethod(idx, method)}>{sorter.name}</Menu.Item>
                            })}
                        </Menu>
                    }>
                    <Button type='link' size='small' style={{color: 'rgba(255,255,255,.65)', fontWeight: 'bold'}}>{sorters[sortMethod[idx]].name}</Button>
                </Dropdown>
                <Tag color='#fff'><label style={{ color: 'black' }}>{group.tasks.length}</label></Tag>
            </div>
        );
    };

    return (
        <Row type='flex' justify='space-between' style={{margin: 8}}>
            {groups.map((group, idx) => {
                const state = TaskStatus[idx];
                
                return (
                    <Col key={state.type} span={6}>
                        <Card
                            title={<div style={{ color: 'white' }}><Icon type={state.icon} style={{marginRight: 4}} />{state.name}</div>}
                            extra={makeGroupExtra(group, idx)}
                            headStyle={{backgroundColor: state.color, height: 40}}
                            bodyStyle={{padding: 0, margin: 4}}
                            style={{margin: 4}}>
                            
                            {group.tasks.length == 0 ? <Empty description='暂无数据' /> : group.tasks.map(task => {
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