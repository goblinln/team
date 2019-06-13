import * as React from 'react';
import { AntTreeNodeMouseEvent, AntTreeNodeSelectedEvent } from 'antd/lib/tree/Tree';

import {
    Button,
    Divider,
    Dropdown,
    Icon,
    Input,
    Layout,
    Menu,
    Modal,
    Row,
    Tag,
    Tree,
    message,
} from 'antd';

import { IDocument } from '../../common/Protocol';
import { Fetch } from '../../common/Request'; 
import * as Markdown from '../markdown/Markdown';

/**
 * WIKI模块
 */
export const Page = () => {
    /**
     * 内容编辑器
     */
    const ContentEditor = (props: {id: number, value: string}) => {
        const [content, setContent] = React.useState<string>(props.value);

        const uploadImage = (img: File, done: (url: string) => void) => {
            let param = new FormData();
            param.append('img', img, img.name);

            Fetch.post('/api/file/upload', param, rsp => {
                rsp.err ? message.error(rsp.err, 1) : done(rsp.data.url);
            });
        };

        const modify = () => {
            Fetch.patch(`/api/document/${props.id}/content`, new URLSearchParams({content: content}), rsp => {
                rsp.err ? message.error(rsp.err, 1) : fetchOne(props.id, false);
            })
        }

        return (
            <div>
                <Row>
                    <Markdown.Editor rows={24} value={content} onChange={data => setContent(data)} onUpload={uploadImage}/>
                </Row>
                <Row type='flex' justify='center' style={{marginTop: 8}}>
                    <Button type='primary' style={{marginRight: 8}} onClick={() => modify()}>修改</Button>
                    <Button onClick={() => setView(prev => ({doc: prev.doc, isEditing: false}))}>取消</Button>
                </Row>
            </div>
        );
    };

    /**
     * 状态列表
     */
    const [docs, setDocs] = React.useState<{[key:number]: IDocument}>(null);
    const [nodes, setNodes] = React.useState<{[key:number]: number[]}>(null);
    const [view, setView] = React.useState<{doc: IDocument, isEditing: boolean}>(null);
    const [contextMenu, setContextMenu] = React.useState<{rect: DOMRect, docId: number}>(null);

    /**
     * 初始拉取
     */
    React.useEffect(() => fetchAll(), []);

    /**
     * 拉取文档列表
     */
    const fetchAll = () => {
        Fetch.get('/api/document/list', rsp => {
            if (rsp.err) {
                message.error(rsp.err, 1);
                return;
            }

            if (rsp.data.length == 0) {
                setNodes(null);
                setDocs(null);
                setView(null);
                return;
            }

            let map: {[key:number]: IDocument} = {}
            let parsed: {[key:number]: number[]} = {}

            rsp.data.forEach((doc: IDocument) => {
                map[doc.id] = doc;

                if (parsed[doc.parent] != null) {
                    parsed[doc.parent].push(doc.id);
                } else {
                    parsed[doc.parent] = [doc.id];
                }
            });

            setDocs(map);
            setNodes(parsed);
        });
    };

    /**
     * 取文档，用于编辑或查看
     */
    const fetchOne = (id: number, isEditing: boolean) => {
        Fetch.get(`/api/document/${id}`, rsp => {
            rsp.data.err ? message.error(rsp.data.err, 1) : setView({doc: rsp.data, isEditing: isEditing});
        });
    }

    /**
     * 递归生成树
     */
    const makeTreeNode = (id: number) => {
        if (!nodes) {
            return <Tree.TreeNode title="右键点击，新建文档" key="-1"/>
        }

        let children = nodes[id];
        if (!children) return null;

        return children.map(child => {
            if (!docs) return null;

            let doc = docs[child];
            if (!doc) return null;

            return (                   
                <Tree.TreeNode title={doc.title} key={doc.id.toString()} doc={doc.id}>
                    {makeTreeNode(doc.id)}
                </Tree.TreeNode>
            );
        })
    };

    /**
     * 选中查看文档
     */
    const selectTree = (selectedKeys: string[], ev: AntTreeNodeSelectedEvent) => {
        if (!nodes) return;
        fetchOne(ev.node.props.doc, false);
    };

    /**
     * 右键菜单
     */
    const rightClickTree = (ev: AntTreeNodeMouseEvent) => {            
        setContextMenu({
            rect: ev.event.currentTarget.getBoundingClientRect(),
            docId: ev.node.props.doc,
        });
    };

    /**
     * 新建
     */
    const addDoc = (parent: number) => {
        setContextMenu(null);

        let name: string = '';
        let parentNode = parent == -1 ? '无' : docs[parent].title;

        Modal.confirm({
            title: '新建文档',
            width: 200,
            maskClosable: true,
            icon: null,
            content: (
                <div style={{marginTop: 24, textAlign: 'center'}}>
                    <p style={{marginBottom: 0}}>父节点 <Tag color='#108ee9'>{parentNode}</Tag></p>                       
                    <Input style={{marginTop: 24}} onChange={ev => name = ev.target.value}/>
                </div>
            ),
            okText: '提交',
            onOk: () => {
                if (name.length == 0) {
                    message.error('文档名不可为空', 1);
                    return;
                }

                if (nodes) {
                    let childs = nodes[parent] || [];
                    for (let i = 0; i < childs.length; ++i) {
                        let doc = docs[childs[i]];
                        if (doc && doc.title == name) {
                            message.error(`同级目录下，已存在文档名为【${name}】`, 1);
                            return;
                        }
                    }
                }                

                let param = new FormData();
                param.append('title', name);
                param.append('parent', parent.toString());

                Fetch.post('/api/document', param, rsp => {
                    rsp.err ? message.error(rsp.err, 1) : fetchAll();
                });
            },
            cancelText: '取消',
        });
    };

    /**
     * 重命名文档
     */
    const renameDoc = (id: number) => {
        setContextMenu(null);

        let name: string = docs[id].title;

        Modal.confirm({
            title: '重命名文档',
            width: 200,
            maskClosable: true,
            icon: null,
            content: (
                <Input style={{marginTop: 24}} defaultValue={name} onChange={ev => name = ev.target.value}/>
            ),
            okText: '提交',
            onOk: () => {
                if (name == docs[id].title) {
                    return;
                }

                if (name.length == 0) {
                    message.error('文档名不可为空', 1);
                    return;
                }

                let parent = docs[id].parent;
                let childs = nodes[parent] || [];
                for (let i = 0; i < childs.length; ++i) {
                    let doc = docs[childs[i]];
                    if (doc && doc.title == name) {
                        message.error(`同级目录下，已存在文档名为【${name}】`, 1);
                        return;
                    }
                }
                
                let param = new FormData();
                param.append('title', name);

                Fetch.patch(`/api/document/${id}/title`, param, rsp => {
                    rsp.err ? message.error(rsp.err, 1) : fetchAll();
                });
            },
            cancelText: '取消',
        });
    };

    /**
     * 编辑文档
     */
    const editDoc = (id: number) => {
        setContextMenu(null);
        fetchOne(id, true);
    };

    /**
     * 删除文档
     */
    const delDoc = (id: number) => {
        setContextMenu(null);

        Fetch.delete(`/api/document/${id}`, rsp => {
            rsp.err ? message.error(rsp.err, 1) : fetchAll();
        })
    };

    return (
        <Layout style={{width: '100%', height: '100%'}}>
            <Layout.Sider width={250} theme='light' style={{borderRight: '1px solid rgba(179,179,179,1)'}}>
                <Row style={{padding: 8}}>
                    <span style={{fontSize: '1.2em', fontWeight: 'bolder', color: 'rgba(0,0,0,.7)'}}><Icon type='book' style={{marginRight: 8}}/>文档列表</span>
                </Row>

                <Divider style={{margin: 0, background: '#dad5d5'}}/>

                <Tree style={{marginLeft: 8}} showLine defaultExpandAll={true} onSelect={selectTree} onRightClick={rightClickTree}>                    
                    {makeTreeNode(-1)}
                </Tree>

                {contextMenu && (
                    <Dropdown          
                        visible={true}
                        onVisibleChange={v => !v && setContextMenu(null)} 
                        overlay={(
                            <Menu mode='vertical'>
                                {!nodes ? (
                                    <Menu.Item key="create" onClick={() => addDoc(-1)}>新建文档</Menu.Item>
                                ) : [
                                    <Menu.SubMenu key='add' title='新建文档'>
                                        <Menu.Item key='sibling' onClick={() => addDoc(docs[contextMenu.docId].parent)}>同级文档</Menu.Item>
                                        <Menu.Item key='child' onClick={() => addDoc(contextMenu.docId)}>子级文档</Menu.Item>
                                    </Menu.SubMenu>,                              
                                    <Menu.Item key='rename' onClick={() => renameDoc(contextMenu.docId)}>重命名</Menu.Item>,
                                    <Menu.Item key='edit' onClick={() => editDoc(contextMenu.docId)}>编辑</Menu.Item>,
                                    <Menu.Item key='delete' onClick={() => delDoc(contextMenu.docId)}>删除</Menu.Item>,
                                ]}                                
                            </Menu>
                        )}>
                        <div
                            style={{
                                position: 'absolute',
                                left: contextMenu.rect.x - 64, 
                                top: contextMenu.rect.y, 
                                width: contextMenu.rect.width, 
                                height: contextMenu.rect.height}}/>    
                    </Dropdown>
                )}
            </Layout.Sider>

            <Layout.Content>
                {!view ? null : (
                    view.isEditing ? (
                        <div style={{margin: 16}}><ContentEditor id={view.doc.id} value={view.doc.content}/></div>
                    ) : (
                        <div style={{margin: 16}}><Markdown.Renderer source={view.doc.content}/></div>
                    )
                )}
            </Layout.Content>
        </Layout>
    );
};