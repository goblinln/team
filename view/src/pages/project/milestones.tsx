import * as React from 'react';
import * as moment from 'moment';

import { TableColumn, FormProxy, FormFieldValidator, Modal, Form, Input, Card, Row, Icon, Button, Table } from '../../components';
import { ProjectMilestone } from '../../common/protocol';
import { request } from '../../common/request';

export const Milestones = (props: {pid: number, isAdmin: boolean}) => {
    const [milestones, setMilestones] = React.useState<ProjectMilestone[]>([]);

    const schema: TableColumn[] = [
        {label: 'ID', dataIndex: 'id'},
        {label: '名称', dataIndex: 'name'},
        {label: '开始时间', dataIndex: 'startTime'},
        {label: '结束时间', dataIndex: 'endTime'},
        {label: '操作', renderer: (data: ProjectMilestone) => (
            <span>
                <a className='link' onClick={() => null}>详情</a>
                {props.isAdmin && [
                    <div key='d-0' className='divider-v'/>,
                    <a key='edit' className='link' onClick={() => editMilestone(data)}>编辑</a>,
                    <div key='d-1' className='divider-v'/>,
                    <a key='delete' className='link' onClick={() => delMilestone(data)}>删除</a>,
                ]}                
            </span>
        )}
    ];

    const validator: {[k:string]:FormFieldValidator} = {
        name: {required: '里程碑名称不可为空'},
        startTime: {required: '请设置计划开始时间'},
        endTime: {required: '请设置计划结束时间'},
    };

    React.useEffect(() => fetchMilestones(), [props]);

    const fetchMilestones = () => {
        request({url: `/api/project/${props.pid}/milestone/list`, success: (data: ProjectMilestone[]) => {
            data.sort((a, b) => {
                let offset = moment(b.startTime).diff(moment(a.startTime))
                return offset == 0 ? moment(b.endTime).diff(moment(a.endTime)) : offset;
            })

            setMilestones(data);
        }});
    };

    const addMilestone = () => {
        let form: FormProxy = null;
        let closer: () => void = null;

        const submit = (ev: React.FormEvent<HTMLFormElement>) => {
            ev.preventDefault();
            request({
                url: `/api/project/${props.pid}/milestone`, 
                method: 'POST', 
                data: new FormData(ev.currentTarget), 
                success: () => {closer(); fetchMilestones()}});
        };

        closer = Modal.open({
            title: '新建里程碑',
            body: (
                <Form style={{width: 300}} form={() => {form = Form.useForm(validator); return form}} onSubmit={submit}>
                    <Form.Field htmlFor='name' label='名称'>
                        <Input name='name' autoComplete='off'/>
                    </Form.Field>
                    <Form.Field htmlFor='startTime' label='开始时间'>
                        <Input.DatePicker name='startTime' mode='date' value={moment().format('YYYY-MM-DD')}/>
                    </Form.Field>
                    <Form.Field htmlFor='endTime' label='结束时间'>
                        <Input.DatePicker name='endTime' mode='date' value={moment().add(1, 'd').format('YYYY-MM-DD')}/>
                    </Form.Field>
                </Form> 
            ),
            onOk: () => {form.submit(); return false}
        });
    };

    const editMilestone = (m: ProjectMilestone) => {
        let form: FormProxy = null;
        let closer: () => void = null;

        const submit = (ev: React.FormEvent<HTMLFormElement>) => {
            ev.preventDefault();
            request({
                url: `/api/project/${props.pid}/milestone/${m.id}`, 
                method: 'PUT', 
                data: new FormData(ev.currentTarget), 
                success: () => {closer(); fetchMilestones()}});
        };

        closer = Modal.open({
            title: '编辑里程碑',
            body: (
                <Form style={{width: 300}} form={() => {form = Form.useForm(validator); return form}} onSubmit={submit}>
                    <Form.Field htmlFor='name' label='名称'>
                        <Input name='name' autoComplete='off' value={m.name}/>
                    </Form.Field>
                    <Form.Field htmlFor='startTime' label='开始时间'>
                        <Input.DatePicker name='startTime' mode='date' value={m.startTime}/>
                    </Form.Field>
                    <Form.Field htmlFor='endTime' label='结束时间'>
                        <Input.DatePicker name='endTime' mode='date' value={m.endTime}/>
                    </Form.Field>
                </Form>
            ),
            onOk: () => {form.submit(); return false}
        });
    };

    const delMilestone = (m: ProjectMilestone) => {
        Modal.open({
            title: '删除确认',
            body: <div className='my-2'>确定要删除里程碑【{m.name}】吗（相关任务的里程碑会被置空）？</div>,
            onOk: () => {
                request({
                    url: `/api/project/${props.pid}/milestone/${m.id}`, 
                    method: 'DELETE', 
                    success: fetchMilestones
                });
            }
        });
    };

    return (
        <div className='m-4'>
            <Card
                className='mt-3'
                bodyProps={{className: 'p-2'}}
                header={props.isAdmin ? (
                    <Row flex={{align: 'middle', justify: 'space-between'}}>
                        <span>
                            <Icon type='idcard' className='mr-2'/>
                            里程碑列表
                        </span>
                        <Button theme='link' onClick={addMilestone}><Icon type='plus' className='mr-1'/>新建里程碑</Button>
                    </Row>) : <span><Icon type='idcard' className='mr-2'/>里程碑列表</span>
                }
                bordered
                shadowed>
                <Table dataSource={milestones} columns={schema}/>
            </Card>
        </div>
    );
}