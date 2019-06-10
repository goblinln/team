import * as React from 'react';
import {WrappedFormUtils} from 'antd/lib/form/Form';

import {
    Button,
    Form,
    Icon,
    Input,
    Layout,
    PageHeader,
    Select,
    Row,
    message,
} from 'antd';

import { ITask, IProject } from '../../common/Protocol';
import { Fetch } from '../../common/Request';
import { Creator } from './Creator';
import { Board } from './Board';
import { Gantt } from './Gantt';
import * as TaskViewer from './Viewer';

/**
 * 任务模块主页
 */
export const Page = () => {
    /**
     * 任务过滤Form的可配置属性
     */
    interface IFilterFormProps {
        /**
         * 可选的项目列表
         */
        projs: IProject[];
        /**
         * 输入发生变化时的回调
         */
        onChange: (fields: {[k:string]: any}) => void;
        /**
         * 自动生成的Form操作工具
         */
        form: WrappedFormUtils;
    }

    /**
     * 任务过滤Form
     */
    const FilterForm = Form.create<IFilterFormProps>()((props: IFilterFormProps) => {
        const {getFieldDecorator, resetFields} = props.form;

        /**
         * 状态列表
         */
        const [validBranches, setValidBranches] = React.useState<string[]>([]);

        /**
         * 选择项目回调
         */
        const selectProj = (projId: any) => {
            resetFields(['branch']);
            let find = false;
            for (let i = 0; i < props.projs.length; ++i) {
                if (props.projs[i].id == projId) {
                    setValidBranches(props.projs[i].branches);
                    find = true;
                    break;
                }
            }
            if (!find) setValidBranches([]);            
            props.onChange({proj: projId});
        }

        /**
         * 重置功能
         */
        const reset = () => {
            resetFields();
            setValidBranches([]);
            props.onChange({prop: null, branch: null, name: ''});
        };

        return (
            <Form layout='inline'>
                <Form.Item label='项目'>
                    {getFieldDecorator('proj', {})(
                        <Select id='proj' style={{minWidth: 128}} onChange={ev => selectProj(ev.valueOf())}>
                            {props.projs.map(proj => {
                                return <Select.Option key={proj.id} value={proj.id}>{proj.name}</Select.Option>
                            })}
                        </Select>
                    )}
                </Form.Item>

                <Form.Item label='分支'>
                    {getFieldDecorator('branch', {})(
                        <Select id='branch' style={{minWidth: 128}} onChange={ev => props.onChange({branch: ev.valueOf()})}>
                            {validBranches.map((branch, idx) => {
                                return <Select.Option key={idx} value={idx}>{branch}</Select.Option>
                            })}
                        </Select>
                    )}
                </Form.Item>

                <Form.Item label='任务名'>
                    {getFieldDecorator('name', {})(
                        <Input id='name' onChange={ev => props.onChange({name: ev.target.value})}/>
                    )}
                </Form.Item>

                <Form.Item>
                    <Button onClick={reset}>重置</Button>
                </Form.Item>
            </Form>
        );
    });

    /**
     * 状态列表
     */
    const [tasks, setTasks] = React.useState<ITask[]>([]);
    const [projs, setProjs] = React.useState<IProject[]>([]);
    const [visibleTasks, setVisibleTasks] = React.useState<ITask[]>([]);
    const [filter, setFilter] = React.useState<{[k:string]: any}>({name: ''});
    const [isCreator, setIsCreator] = React.useState<boolean>(false);
    const [isGantt, setIsGantt] = React.useState<boolean>(false);
    const [isFilterShow, setIsFilterShow] = React.useState<boolean>(false);
    const taskDetailAnchor = React.useRef<any>();

    /**
     * 进入页面时拉取任务列表
     */
    React.useEffect(() => {
        TaskViewer.default.init(taskDetailAnchor, task => {
            fetchAll();
        });
        fetchAll();
    }, []);

    /**
     * 任务列表发生变化后，重新过滤任务
     */
    React.useEffect(() => {
        let result: ITask[] = [];
        let projects: IProject[] = [];

        tasks.forEach(task => {
            let find = false;
            for (let idx: number = 0; idx < projects.length; ++idx) {
                if (projects[idx].id == task.proj.id) {
                    find = true;
                    break;
                }
            }

            if (!find) {
                projects.push({
                    id: task.proj.id,
                    name: task.proj.name,
                    branches: [...task.proj.branches],
                });
            }

            if (filter.proj && filter.proj != task.proj.id) return;
            if (filter.branch && filter.branch != task.branch) return;
            if (filter.name.length > 0 && task.name.indexOf(filter.name) == -1) return;
            result.push(task);
        });

        setProjs(projects);
        setVisibleTasks(result);
    }, [tasks]);

    /**
     * 过滤选项发生变化后，重新过滤任务
     */
    React.useEffect(() => {
        let result: ITask[] = [];

        tasks.forEach(task => {
            if (filter.proj && filter.proj != task.proj.id) return;
            if (filter.branch && filter.branch != task.branch) return;
            if (filter.name.length > 0 && task.name.indexOf(filter.name) == -1) return;
            result.push(task);
        });

        setVisibleTasks(result);
    }, [filter]);

    /**
     * 取任务列表
     */
    const fetchAll = () => {
        Fetch.get(`/api/task/mine`, rsp => {
            rsp.err ? message.error(rsp.err, 1) : setTasks(rsp.data);
        });
    };

    /**
     * 修改过滤选项
     */
    const modifyFilter = (fields: {[k: string]: any}) => {
        let result: {[k:string]: any} = {};
        for (let k in filter) result[k] = filter[k];
        for (let k in fields) result[k] = fields[k];
        setFilter(result);
    };

    /**
     * 搜索栏，防止页面刷新时搜索内容重置
     */
    const filterbar = React.useMemo(() => {
        return <Row type='flex' justify='center'><FilterForm projs={projs} onChange={modifyFilter} /></Row>;
    }, [projs]);

    return (
        <Layout style={{ width: '100%', height: '100%' }}>
            <Layout.Content>
                <PageHeader
                        title={isCreator ? '发布任务' : '任务列表'}
                        onBack={isCreator ? () => setIsCreator(false) : null}
                        extra={isCreator ? null : (
                            <div style={{marginTop: 4}}>
                                <Button style={{marginRight: 4}} onClick={() => fetchAll()}><Icon type='reload' />刷新</Button>,
                                <Button style={{marginRight: 4}} onClick={() => setIsGantt(prev => !prev)}><Icon type='switcher' />{isGantt ? '看板模式' : '甘特图'}</Button>,
                                <Button type={isFilterShow ? 'primary' : 'default'} style={{marginRight: 4}} onClick={() => setIsFilterShow(prev => !prev)}><Icon type='filter' />任务过滤</Button>,
                                <Button type='primary' onClick={() => setIsCreator(true)}><Icon type='plus' />发布任务</Button>,
                            </div>
                        )}
                        style={{ padding: '12px 16px', borderBottom: '1px solid #dee2e6' }}>
                        {!isCreator && isFilterShow && filterbar}
                </PageHeader>
                
                {isCreator ? <Creator onFinish={() => { setIsCreator(false); fetchAll(); }} /> : (isGantt ? <Gantt tasks={visibleTasks}/> : <Board tasks={visibleTasks}/>)}

                <div ref={taskDetailAnchor}/>
            </Layout.Content>
        </Layout>
    );
};