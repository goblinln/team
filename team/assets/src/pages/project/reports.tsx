import * as React from 'react';
import * as moment from 'moment';

import {Row, Col, Icon, Badge} from '../../components';
import {Project, Task} from '../../common/protocol';
import {request} from '../../common/request';
import {Viewer} from '../task/viewer';

interface ReportData {
    archived: Task[];
    unarchived: Task[];
};

export const Reports = (props: {proj: Project, isReadOnly: boolean}) => {
    const [week, setWeek] = React.useState<moment.Moment>(moment().startOf('week'));
    const [data, setData] = React.useState<ReportData>({archived: [], unarchived: []});

    React.useEffect(() => {
        fetchReports();
    }, [props, week]);

    const fetchReports = () => {
        request({
            url: `/api/project/${props.proj.id}/report/${week.unix()}`,
            success: (rsp: ReportData) => {
                rsp.unarchived.sort((a, b) => {
                    if (a.state != b.state) {
                        return a.state - b.state;
                    } else {
                        return moment(a.endTime).diff(b.endTime);
                    }
                });
                setData(rsp);
            }
        })
    };

    return (
        <div className='m-2 fg-muted'>
            <Row flex={{align: 'middle', justify: 'center'}}>
                <Icon type='left-circle-fill' style={{fontSize: '2.5em'}} onClick={() => setWeek(moment(week).subtract(1, 'week'))}/>
                <div className='mx-3 text-center'>
                    <p style={{fontSize: '2em', fontWeight: 'bolder'}}>第{week.weeks()}周项目周报</p>
                    <small><Badge theme='highlight' className='bg-darkgray'>{moment(week).format('YYYY/MM/DD')} - {moment(week).endOf('week').format('YYYY/MM/DD')}</Badge></small>
                </div>
                <Icon type='right-circle-fill' style={{fontSize: '2.5em'}} onClick={() => setWeek(moment(week).add(1, 'week'))}/>
            </Row>

            <Row className='mt-3' space={8}>
                <Col span={{xs: 6}}>
                    <div style={{fontSize: '1.5em', fontWeight: 'bold'}}>
                        <Icon type='frown' className='mr-2'/>未验收任务
                    </div>

                    <div className='divider-h my-2'/>

                    {data.unarchived.map(t => (
                        <Row flex={{align: 'middle', justify: 'space-between'}} className='mb-2'>
                            <span className='px-1 text-ellipsis'>
                                <Icon type={t.state==3?'question-circle':'close-circle'} className={t.state==3?'mr-1 fg-info':'mr-1 fg-danger'}/>
                                {t.endTime}
                                <a className='ml-1 link' onClick={() => Viewer.open(t.id, props.isReadOnly?null:fetchReports)}>{t.name}</a>
                            </span>

                            <span className='px-1'>
                                {t.creator.name}<Icon type='right'/>{t.developer.name}<Icon type='right'/>{t.tester.name}
                            </span>
                        </Row>
                    ))}
                </Col>

                <Col span={{xs: 6}}>
                    <div style={{fontSize: '1.5em', fontWeight: 'bold'}}>
                        <Icon type='smile-circle' className='mr-2'/>已验收任务
                    </div>

                    <div className='divider-h my-2'/>

                    {data.archived.map(t => (
                        <Row flex={{align: 'middle', justify: 'space-between'}} className='mb-2'>
                            <span className='px-1 text-ellipsis'>
                                <Icon type='check' className='mr-1 fg-success'/>
                                <a className='ml-1 link' onClick={() => Viewer.open(t.id)}>{t.name}</a>
                            </span>

                            <span className='px-1'>
                                {t.creator.name}<Icon type='right'/>{t.developer.name}<Icon type='right'/>{t.tester.name}
                            </span>
                        </Row>
                    ))}
                </Col>
            </Row>
        </div>
    );
};