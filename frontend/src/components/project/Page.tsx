import * as React from 'react';

import {
    Collapse,
    Divider,
    Empty,
    Icon,
    Layout,
    Menu,
    Row,
    Tag,
    message,
} from 'antd';

import { IProject } from '../../common/Protocol';
import { Fetch } from '../../common/Request';
import { Tasks } from './Tasks';
import { Reports } from './Reports';
import { Manage } from './Manage';

/**
 * 项目模块主页
 */
export const Page = (props: {uid: number}) => {

    /**
     * 状态列表
     */
    const [projs, setProjs] = React.useState<IProject[]>([]);
    const [subpage, setSubpage] = React.useState<JSX.Element>(<Empty description='请选择相应操作' style={{marginTop: '10%'}}/>);

    /**
     * 拉取项目列表
     */
    React.useEffect(() => {
        Fetch.get('/api/project/mine', rsp => {
            rsp.err ? message.error(rsp.err, 1) : setProjs(rsp.data);
        })
    }, []);

    return (
        <Layout style={{width: '100%', height: '100%'}}>
            <Layout.Sider width={200} theme='light' style={{borderRight: '1px solid rgba(179,179,179,1)', background: '#f0f2f5'}}>
                <Row style={{padding: 8}}>
                    <span style={{fontSize: '1.2em', fontWeight: 'bolder', color: 'rgba(0,0,0,.7)'}}><Icon type='pie-chart' style={{marginRight: 8}}/>项目列表</span>
                </Row>

                <Divider style={{margin: 0, background: '#dad5d5'}}/>

                {projs.length == 0 ? <Empty style={{marginTop: 8}} description='未加入任何项目'/> : (
                    <Collapse accordion={true} defaultActiveKey={['0']} style={{margin: 8}} onChange={() => setSubpage(<Empty description='请选择相应操作' style={{marginTop: '10%'}}/>)}>
                        {projs.map((proj, idx) => {
                            let isAdmin = false;

                            for (let i = 0; i < proj.members.length; ++i) {
                                if (proj.members[i].user.id == props.uid) {
                                    isAdmin = proj.members[i].isAdmin || false;
                                    break;
                                }
                            }

                            return (
                                <Collapse.Panel key={idx.toString()} header={proj.name} extra={isAdmin ? <Tag color='green' style={{margin: 0}}>管理员</Tag> : null}>
                                    <Menu mode='vertical'>
                                        <Menu.Item key={`task_${idx}`} style={{margin: 0}} onClick={() => setSubpage(<Tasks proj={proj} isReadonly={!isAdmin}/>)}>任务列表</Menu.Item>
                                        <Menu.Item key={`reports_${idx}`} style={{margin: 0}} onClick={() => setSubpage(<Reports proj={proj} isReadonly={!isAdmin}/>)}>验收相关</Menu.Item>
                                        {isAdmin && <Menu.Item key={`manager_${idx}`} style={{margin: 0}} onClick={() => setSubpage(<Manage pid={proj.id} />)}>项目管理</Menu.Item>}
                                    </Menu>
                                </Collapse.Panel>
                            );
                        })}
                    </Collapse>
                )}
            </Layout.Sider>

            <Layout.Content>
                {subpage}
            </Layout.Content>
        </Layout>
    );
};