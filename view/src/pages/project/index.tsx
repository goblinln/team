import * as React from 'react';

import {Layout, Icon, Menu, Empty, Badge, Row} from '../../components';
import {Project} from '../../common/protocol';
import {request} from '../../common/request';

import {Tasks} from './tasks';
import {Members} from './members';
import {Milestones} from './milestones';

export const ProjectPage = (props: {uid: number}) => {
    const [projs, setProjs] = React.useState<Project[]>([]);
    const [page, setPage] = React.useState<JSX.Element>();

    React.useEffect(() => {
        request({url: '/api/project/mine', success: setProjs});
    }, []);

    return (
        <Layout style={{width: '100%', height: '100%'}}>
            <Layout.Sider width={200} theme='light'>            
                <div style={{padding: 8, borderBottom: '1px solid #E2E2E2'}}>
                    <label className='text-bold fg-muted' style={{fontSize: '1.2em'}}><Icon type='pie-chart' className='mr-2'/>项目列表</label>
                </div>

                {projs.length == 0?<Empty label='您还未加入任何项目'/>:(
                    <Menu theme='light'>
                        {projs.map(p => {
                            let isAdmin = false;

                            for (let i = 0; i < p.members.length; ++i) {
                                if (p.members[i].user.id == props.uid) {
                                    isAdmin = p.members[i].isAdmin || false;
                                    break;
                                }
                            }

                            return (
                                <Menu.SubMenu key={p.id} collapse='disabled' label={<Row flex={{align: 'middle', justify: 'space-between'}}>{p.name}<Badge className='ml-2' theme='info'>{isAdmin?'管理员':'成员'}</Badge></Row>}>
                                    <Menu.Item onClick={() => setPage(<Tasks proj={p} isAdmin={isAdmin}/>)}>任务列表</Menu.Item>
                                    <Menu.Item onClick={() => setPage(<Milestones pid={p.id} isAdmin={isAdmin}/>)}>里程计划</Menu.Item>
                                    {isAdmin&&<Menu.Item onClick={() => setPage(<Members pid={p.id}/>)}>成员管理</Menu.Item>}
                                </Menu.SubMenu>
                            );
                        })}
                    </Menu>
                )}
            </Layout.Sider>

            <Layout.Content>
                {page}
            </Layout.Content>
        </Layout>
    );
};