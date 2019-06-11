import * as React from 'react';
import * as ReactDOM from 'react-dom';

import {
    Avatar,
    Badge,
    Icon,
    Layout,
    Menu,
    Spin,
} from 'antd';

import { INotice, IUser } from '../../common/Protocol';
import { Fetch, FetchStateNotifier } from '../../common/Request';
import * as Login from '../user/Login';
import * as Profile from '../user/Profile';
import * as Task from '../task/Page';
import * as Project from '../project/Page';
import * as Document from '../document/Page';
import * as Share from '../share/Page';
import * as Admin from '../admin/Page';

/**
 * 项目主页
 */
export const Index = () => {

    /**
     * 状态列表
     */
    const [user, setUser] = React.useState<IUser>(null);
    const [notices, setNoticies] = React.useState<INotice[]>([]);
    const [subpage, setSubpage] = React.useState<JSX.Element>(null);
    const [loading, setLoading] = React.useState<boolean>(false);
    const popupAnchor = React.useRef<any>();
    const loadingDely = React.useRef<any>();

    /**
     * 主菜单
     */
    const menu = [
        {
            name: '任务',
            key: 'task',
            icon: 'schedule',
            click: () => setSubpage(<Task.Page/>),
        },
        {
            name: '项目',
            key: 'project',
            icon: 'pie-chart',
            click: () => setSubpage(<Project.Page uid={user.id}/>),
        },
        {
            name: '文档',
            key: 'wiki',
            icon: 'read',
            click: () => setSubpage(<Document.Page/>),
        },
        {
            name: '分享',
            key: 'share',
            icon: 'cloud-upload',
            click: () => setSubpage(<Share.Page/>),
        },
        {
            name: '管理',
            key: 'admin',
            icon: 'setting',
            needAdmin: true,
            click: () => setSubpage(<Admin.Page/>),
        }
    ];

    /**
     * 取帐号信息并且取一次通知
     */
    React.useEffect(() => {
        Fetch.notifier = new FetchStateNotifier(state => {
            if (loadingDely.current) {
                clearTimeout(loadingDely.current)
                loadingDely.current = null;
            }

            loadingDely.current = setTimeout(isFetching => {
                setLoading(isFetching);
                loadingDely.current = null;
            }, 20, state);
        });
        fetchUser();
    }, []);

    /**
     * 定时更新通知内容
     */
    React.useEffect(() => {
        const noticeFetcher = setInterval(fetchNotice, 60000);
        return () => clearInterval(noticeFetcher);
    }, []);

    /**
     * 取一次用户信息
     */
    const fetchUser = () => {
        Fetch.get('/api/user', rsp => {
            if (!rsp.err) {
                setUser(rsp.data);
                fetchNotice(true);
            }
        })
    };

    /**
     * 取一次通知消息
     */
    const fetchNotice = (force?: boolean) => {
        if (!user && !force) return;
        Fetch.get('/api/notice/list', rsp => { !rsp.err && setNoticies(rsp.data) });
    };

    /**
     * 右侧弹出子页面（抽屉模式）
     */
    const popup = (page: JSX.Element) => {
        ReactDOM.render(page, popupAnchor.current);
    };

    /**
     * 关闭弹出子界面（抽屉模式）
     */
    const closePopup = () => {
        ReactDOM.render(null, popupAnchor.current);
    }

    /**
     * 退出系统
     */
    const logout = () => {
        Fetch.post('/logout', null, () => { setUser(null) });
    }

    return !user ? <Login.View onLogined={() => fetchUser()} /> : (
        <Layout style={{ width: '100vw', height: '100vh' }}>
            <Layout.Sider width={64}>
                <div style={{ textAlign: 'center', marginTop: 16, marginBottom: 16 }}>
                    <div
                        onClick={() => {
                            popup(<Profile.View
                                name = {user.name}
                                account = {user.account}
                                avatar = {user.avatar}
                                notices = {notices}
                                onClose = {(newAvatar, refreshNotice) => {
                                    if (newAvatar != null) setUser(prev => ({id: prev.id, name: prev.name, avatar: newAvatar}));
                                    if (refreshNotice) fetchNotice();
                                    closePopup();
                                }}/>);
                        }}>
                        <Badge count={notices.length} offset={[-24, 44]}>                            
                            <Avatar icon='user' src={user.avatar} size={48} />
                        </Badge>
                    </div>
                </div>

                <Menu theme='dark' defaultSelectedKeys={['task']}>
                    {menu.map(item => {
                        if (item.needAdmin && !user.isSu) return null;

                        return (
                            <Menu.Item title={item.name} key={item.key} onClick={() => item.click()} style={{textAlign: 'center'}}>
                                <Icon type={item.icon} style={{fontSize: '1.5em', margin: 0}}/>
                            </Menu.Item>
                        );
                    })}
                </Menu>

                <div style={{position: 'absolute', left: 0, bottom: 16, width: '100%', textAlign: 'center'}}>
                    <Icon type='logout' style={{color: 'white', fontSize: '1.5em'}} onClick={() => logout()}/>
                </div>
            </Layout.Sider>

            <Layout.Content>                   
                {subpage || <Task.Page />}
                <div ref={popupAnchor} />
                {loading && (
                    <div style={{position: "absolute", left: 64, right: 0, top: 0, bottom: 0, display: 'flex', flexDirection: 'row', justifyContent: 'center', alignItems: 'center', backgroundColor: 'rgba(0, 0, 0, .25)'}} onClick={() => null}>
                        <Spin tip="加载中..." size='large' />
                    </div>
                )}
            </Layout.Content>
        </Layout>
    );
};
