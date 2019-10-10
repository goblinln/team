import * as React from 'react';

import {Button, Icon, Row, Input, Drawer} from '../../components';
import {Task, Project} from '../../common/protocol';
import {ProjectRole} from '../../common/consts';
import {request} from '../../common/request';

import {Board} from '../task/board';
import {Gantt} from '../task/gantt';

export const Tasks = (props: {proj: Project, isAdmin: boolean}) => {
    const {proj, isAdmin} = props;

    const [isGantt, setIsGantt] = React.useState<boolean>(false);
    const [isFilterVisible, setFilterVisible] = React.useState<boolean>(false);
    const [tasks, setTasks] = React.useState<Task[]>([]);
    const [visibleTasks, setVisibleTask] = React.useState<Task[]>([]);
    const [filter, setFilter] = React.useState<{m: number, b: number, n: string}>({m: -1, b: -1, n: ''});

    React.useEffect(() => {
        fetchTasks();
    }, [proj]);

    React.useEffect(() => {
        let ret: Task[] = [];

        tasks.forEach(t => {
            if (filter.b != -1 && t.branch != filter.b) return;
            if (filter.n.length > 0 && t.name.indexOf(filter.n) == -1) return;
            if (filter.m != -1) {
                if (t.creator.id != filter.m && t.developer.id != filter.m && t.tester.id != filter.m) return;
            }

            ret.push(t);
        });

        setVisibleTask(ret);
    }, [tasks, filter]);

    const fetchTasks = () => {
        request({url: `/api/task/project/${proj.id}`, success: setTasks});
    };

    const handleMemberChange = (ev: React.ChangeEvent<HTMLSelectElement>) => {
        let selected = parseInt(ev.target.value);
        setFilter(prev => {
            return {
                m: selected,
                b: prev.b,
                n: prev.n
            }
        });
    };

    const handleBranchChange = (ev: React.ChangeEvent<HTMLSelectElement>) => {
        let selected = parseInt(ev.target.value);
        setFilter(prev => {
            return {
                m: prev.m,
                b: selected,
                n: prev.n
            }
        });
    };

    const handleNameChange = (v: string) => {
        setFilter(prev => {
            return {
                m: prev.m,
                b: prev.b,
                n: v
            }
        });
    };

    const board = React.useMemo(() => <Board tasks={visibleTasks} onModified={isAdmin?fetchTasks:null}/>, [visibleTasks]);
    const gantt = React.useMemo(() => <Gantt tasks={visibleTasks} onModified={isAdmin?fetchTasks:null}/>, [visibleTasks]);

    return (
        <div>
            <div style={{padding: '8px 16px', borderBottom: '1px solid #E2E2E2'}}>
                <Row flex={{align: 'middle', justify: 'space-between'}}>
                    <label className='text-bold fg-muted' style={{fontSize: '1.2em'}}>{`【${proj.name}】任务列表`}</label>
                    <div>
                        <Button size='sm' onClick={() => fetchTasks()}><Icon className='mr-1' type='reload'/>刷新</Button>
                        <Button size='sm' onClick={() => setIsGantt(prev => !prev)}><Icon className='mr-1' type='view'/>{isGantt?'看板模式':'甘特图'}</Button>
                        <Button size='sm' theme={isFilterVisible?'primary':'default'} onClick={() => setFilterVisible(prev => !prev)}><Icon className='mr-1' type='filter'/>任务过滤</Button>
                    </div>
                </Row>

                <div className={`mt-2 center-child ${isFilterVisible?'':' hide'}`}>
                    <div>
                        <label className='mr-1'>选择成员</label>
                        <Input.Select style={{width: 150}} value={filter.m} onChange={handleMemberChange}>
                            <option key={'none'} value={-1}>无要求</option>
                            {proj.members.map(m => <option key={m.user.id} value={m.user.id}>【{ProjectRole[m.role]}】{m.user.name}</option>)}
                        </Input.Select>
                    </div>

                    <div className='ml-3'>
                        <label className='mr-1'>选择分支</label>
                        <Input.Select style={{width: 150}} value={filter.b} onChange={handleBranchChange}>
                            <option key={'none'} value={-1}>无要求</option>
                            {proj.branches.map((b, i) => <option key={i} value={i}>{b}</option>)}
                        </Input.Select>
                    </div>

                    <div className='ml-3'>
                        <label className='mr-1'>任务名</label>
                        <Input style={{width: 150}} value={filter.n} onChange={handleNameChange}/>
                    </div>

                    <Button className='ml-3' size='sm' onClick={() => setFilter({m: -1, b: -1, n: ''})}>重置</Button>
                </div>
            </div>
            
            <div className='px-2 mt-3'>
                {isGantt?gantt:board}
            </div>
        </div>
    );
};
