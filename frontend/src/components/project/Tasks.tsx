import * as React from 'react';
import { WrappedFormUtils } from 'antd/lib/form/Form';

import {
    Button,
    Form,
    Icon,
    Input,
    PageHeader,
    Row,
    Select,
    message,
} from 'antd';

import { IProject, ITask } from '../../common/Protocol';
import { ProjectRole } from '../../common/Consts';
import { Fetch } from '../../common/Request';
import { Gantt } from '../task/Gantt';
import { Board } from '../task/Board';
import * as TaskViewer from '../task/Viewer';

/**
 * 项目任务视图
 */
export const Tasks = (props: {proj: IProject, isReadonly: boolean}) => {
    /**
     * 任务过滤Form的可配置属性
     */
    interface IFilterFormProps {
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
    const FilterForm = Form.create<IFilterFormProps>()((formProps: IFilterFormProps) => {
        const {getFieldDecorator, resetFields} = formProps.form;

        /**
         * 重置功能
         */
        const reset = () => {
            resetFields();
            formProps.onChange({user: null, branch: null, name: ''});
        };

        return (
            <Form layout='inline'>
                <Form.Item label='相关人员'>
                    {getFieldDecorator('user', {})(
                        <Select id='user' style={{minWidth: 128}} onChange={ev => formProps.onChange({user: ev.valueOf()})}>
                            {props.proj.members.map(member => {
                                return <Select.Option key={member.user.id} value={member.user.id}>【{ProjectRole[member.role]}】{member.user.name}</Select.Option>
                            })}
                        </Select>
                    )}
                </Form.Item>

                <Form.Item label='分支'>
                    {getFieldDecorator('branch', {})(
                        <Select id='branch' style={{minWidth: 128}} onChange={ev => formProps.onChange({branch: ev.valueOf()})}>
                            {props.proj.branches.map((branch, idx) => {
                                return <Select.Option key={idx} value={idx}>{branch}</Select.Option>
                            })}
                        </Select>
                    )}
                </Form.Item>

                <Form.Item label='任务名'>
                    {getFieldDecorator('name', {})(
                        <Input id='name' onChange={ev => formProps.onChange({name: ev.target.value})}/>
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
    const [isGantt, setIsGantt] = React.useState<boolean>(false);
    const [isFilterShow, setIsFilterShow] = React.useState<boolean>(false);
    const [filter, setFilter] = React.useState<{[k:string]: any}>({name: ''});
    const [tasks, setTasks] = React.useState<ITask[]>([]);
    const [visibleTasks, setVisibleTasks] = React.useState<ITask[]>([]);
    const taskDetailAnchor = React.useRef<any>(null);

    /**
     * 进入页面时拉取任务列表
     */
    React.useEffect(() => {
        TaskViewer.default.init(taskDetailAnchor, () => fetchAll());
        fetchAll();
    }, []);

    /**
     * 任务列表或过滤选项发生变化后，重新过滤任务
     */
    React.useEffect(() => {
        let result: ITask[] = [];

        tasks.forEach(task => {
            if (filter.user && filter.user != task.developer.id && filter.user != task.creator.id && filter.user != task.tester.id) return;
            if (filter.branch && filter.branch != task.branch) return;
            if (filter.name.length > 0 && task.name.indexOf(filter.name) == -1) return;
            result.push(task);
        });

        setVisibleTasks(result);
    }, [tasks, filter]);

    /**
     * 取任务列表
     */
    const fetchAll = () => {
        Fetch.get(`/api/task/project/${props.proj.id}`, rsp => {
            rsp.err ? message.error(rsp.err, 1) : setTasks(rsp.data)
        })
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
     * 搜索栏，防止刷新
     */
    const filterbar = React.useMemo(() => {
        return <Row type='flex' justify='center'><FilterForm onChange={modifyFilter} /></Row>;
    }, [props.proj]);

    return (
        <div>
            <PageHeader
                    title={`【${props.proj.name}】任务列表`}
                    extra={(
                        <div style={{marginTop: 4}}>
                            <Button style={{marginRight: 4}} onClick={() => fetchAll()}><Icon type='reload' />刷新</Button>,
                            <Button style={{marginRight: 4}} onClick={() => setIsGantt(prev => !prev)}><Icon type='switcher' />{isGantt ? '看板模式' : '甘特图'}</Button>,
                            <Button type={isFilterShow ? 'primary' : 'default'} style={{marginRight: 4}} onClick={() => setIsFilterShow(prev => !prev)}><Icon type='filter' />任务过滤</Button>,
                        </div>
                    )}
                    style={{ padding: '12px 16px', borderBottom: '1px solid #dee2e6' }}>
                    {isFilterShow && filterbar}
            </PageHeader>

            {isGantt ? <Gantt tasks={visibleTasks} isReadonly={props.isReadonly} /> : <Board tasks={visibleTasks} isReadonly={props.isReadonly}/>}

            <div ref={taskDetailAnchor}/>
        </div>
    );
};