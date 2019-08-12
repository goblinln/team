import * as React from 'react';

import {Avatar, Badge, Drawer, Layout, Menu, Icon} from '../../components';
import {request} from '../../common/request';
import {User, Notice} from '../../common/protocol';

import {UserPage} from '../user';
import {TaskPage} from '../task';
import {ProjectPage} from '../project';
import {DocumentPage} from '../document';
import {SharePage} from '../share';
import {AdminPage} from '../admin';

interface MainMenu {
    name: string;
    id: string;
    icon: string;
    click: () => void;
    needAdmin?: boolean;
};

export const Home = () => {
    const [user, setUser] = React.useState<User>({account: 'Unknown', id: 0});
    const [notices, setNotices] = React.useState<Notice[]>([]);
    const [page, setPage] = React.useState<JSX.Element>();

    const menus: MainMenu[] = [
        {name: '任务', id: 'task', icon: 'calendar-check', click: () => setPage(<TaskPage uid={user.id}/>)},
        {name: '项目', id: 'project', icon: 'pie-chart', click: () => setPage(<ProjectPage uid={user.id}/>)},
        {name: '文档', id: 'document', icon: 'read', click: () => setPage(<DocumentPage/>)},
        {name: '分享', id: 'share', icon: 'cloud-upload', click: () => setPage(<SharePage/>)},
        {name: '管理', id: 'admin', icon: 'setting', click: () => setPage(<AdminPage/>), needAdmin: true},
    ];

    React.useEffect(() => {
        fetchUserInfo();
        fetchNotices();
        setInterval(fetchNotices, 60000);
    }, []);

    const fetchUserInfo = () => {
        request({url: '/api/user', success: setUser, dontShowLoading: true});
    };

    const fetchNotices = () => {
        request({url: '/api/notice/list', success: setNotices, dontShowLoading: true});
    };

    const openProfiler = () => {
        Drawer.open({
            width: 350,
            header: '用户信息',
            body: <UserPage
                user={user} 
                notices={notices} 
                onAvatarChanged={a => setUser(prev => {let ret = {...prev}; ret.avatar = a; return ret;})} 
                onNoticeChanged={fetchNotices}/>,
        });
    };

    return (
        <Layout style={{width: '100vw', height: '100vh'}}>
            <Layout.Sider width={64}>
                <div className='text-center my-3'>
                    <div onClick={openProfiler}>
                        <Avatar size={48} src={user.avatar}/>                        
                        {notices.length > 0&&<div style={{marginTop: -20}}><Badge theme='danger' className='r-1'>{notices.length}</Badge></div>}
                    </div>
                </div>

                <Menu defaultActive='task' theme='dark' style={{fontSize: 24}}>
                    {user&&menus.map(m => {
                        if (m.needAdmin && !user.isSu) return null;

                        return (
                            <Menu.Item key={m.id} id={m.id} className='py-2' title={m.name} onClick={m.click}>
                                <Icon type={m.icon}/>
                            </Menu.Item>
                        );
                    })}
                </Menu>

                <div style={{position: 'absolute', left: 0, bottom: 16, width: '100%', fontSize: 24, textAlign: 'center'}}>
                    <Icon type='export' title='退出' onClick={() => location.href = '/logout'}/>
                </div>
            </Layout.Sider>

            <Layout.Content>
                {page||<TaskPage uid={user.id}/>}
            </Layout.Content>
        </Layout>
    );
};