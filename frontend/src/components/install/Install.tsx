import * as React from 'react';
import { WrappedFormUtils } from 'antd/lib/form/Form';

import {
    Button,
    Card,
    Form,
    Input,
    InputNumber,
    Layout,
    Steps,
    Row,
    message,
} from 'antd';

import { Fetch } from '../../common/Request'

/**
 * 安装部署页
 */
export const Install = () => {
    /**
     * 状态列表
     */
    const [current, setCurrent] = React.useState<number>(0);

    /**
     * 配置表单
     */
    const Configure = Form.create()((props: {form: WrappedFormUtils}) => {
        const { getFieldDecorator, validateFields } = props.form;

        const submit = (ev: React.FormEvent<HTMLFormElement>) => {
            ev.preventDefault();

            validateFields(err => {
                if (!err) {
                    Fetch.post('/install/configure', new FormData(ev.currentTarget), rsp => {
                        rsp.err ? message.error(rsp.err, 1) : setCurrent(prev => prev + 1);
                    });                    
                }
            });
        }

        return (
            <Form style={{padding: '0 100px'}} onSubmit={submit}>
                <Card title='监听端口（完成后需要重新启动）' style={{marginBottom: 8}} bodyStyle={{padding: '8px 16px'}}>
                    {getFieldDecorator('port', {
                        initialValue: '8080',
                    })(
                        <InputNumber id='port' name='port' min={80} max={65535}/>
                    )}
                </Card>

                <Card title='MySQL配置' style={{marginBottom: 8}} bodyStyle={{padding: 8}}>
                    <Form.Item label='服务器地址' required style={{marginBottom: 4}}>
                        {getFieldDecorator('mysqlHost', {
                            rules: [
                                { required: true, message: 'MySQL服务器的地址不可为空'},
                            ],
                            initialValue: '127.0.0.1:3306',
                        })(
                            <Input id='mysqlHost' name='mysqlHost'/>
                        )}
                    </Form.Item>

                    <Form.Item label='登录用户' required style={{marginBottom: 4}}>
                        {getFieldDecorator('mysqlUser', {
                            rules: [
                                { required: true, message: 'MySQL登录用户未填写'},
                            ],
                            initialValue: 'root',
                        })(
                            <Input id='mysqlUser' name='mysqlUser'/>
                        )}
                    </Form.Item>

                    <Form.Item label='登录密码' required style={{marginBottom: 4}}>
                        {getFieldDecorator('mysqlPswd', {})(
                            <Input type='password' id='mysqlPswd' name='mysqlPswd'/>
                        )}
                    </Form.Item>

                    <Form.Item label='使用的数据库' required>
                        {getFieldDecorator('mysqlDB', {
                            rules: [
                                { required: true, message: '数据库'},
                            ],
                            initialValue: 'team',
                        })(
                            <Input id='mysqlDB' name='mysqlDB'/>
                        )}
                    </Form.Item>
                </Card>

                <Row type='flex' justify='center'>
                    <Button type='primary' htmlType='submit'>下一步</Button>
                </Row>
            </Form>       
        ); 
    });

    /**
     * 部署进度子页
     */
    const Progress = () => {
        const [progress, setProgress] = React.useState<string[]>([]);
        const [canNext, setCanNext] = React.useState<boolean>(false);

        const fetchStatus = (timer: any) => {
            Fetch.get('/install/status', rsp => {
                setProgress(rsp.data.status);
                if (rsp.data.done) {
                    clearInterval(timer);
                    setCanNext(true);
                } else if (rsp.data.isError) {                    
                    clearInterval(timer);
                }
            })
        }

        React.useEffect(() => {
            fetchStatus(-1);

            let id = setInterval(() => {
                fetchStatus(id)
            }, 1000);
        }, [])

        return (
            <div style={{padding: '0 100px'}}>
                <Card title='进度信息'>
                    <ul style={{marginInlineStart: -10}}>
                        {progress.map((msg, idx) => {
                            return <li key={idx}>{msg}</li>
                        })}
                    </ul>
                </Card>

                <Row type='flex' justify='center' style={{marginTop: 8}}>
                    <Button type='primary' disabled={!canNext} onClick={() => setCurrent(prev => prev + 1)}>下一步</Button>
                </Row>
            </div>            
        );
    };

    /**
     * 配置默认管理员帐号子页
     */
    const CreateAdmin = Form.create()((props: {form: WrappedFormUtils}) => {
        const { getFieldDecorator, validateFields } = props.form;

        const submit = (ev: React.FormEvent<HTMLFormElement>) => {
            ev.preventDefault();
            
            validateFields(err => {
                if (!err) {
                    Fetch.post('/install/admin', new FormData(ev.currentTarget), rsp => {
                        rsp.err ? message.error(rsp.err, 1) : setCurrent(prev => prev + 1);
                    });
                }
            });
        }

        return (
            <Form style={{padding: '0 100px'}} onSubmit={submit}>
                <Card title='默认管理员' style={{marginBottom: 8}} bodyStyle={{padding: 8}}>
                    <Form.Item label='帐号' required style={{marginBottom: 4}}>
                        {getFieldDecorator('account', {
                            rules: [
                                { required: true, message: '登录帐号不可为空'},
                            ],
                            initialValue: 'admin',
                        })(
                            <Input id='account' name='account'/>
                        )}
                    </Form.Item>

                    <Form.Item label='登录用户' required style={{marginBottom: 4}}>
                        {getFieldDecorator('name', {
                            rules: [
                                { required: true, message: '登录用户显示名称未填写'},
                            ],
                            initialValue: '超级管理员',
                        })(
                            <Input id='name' name='name'/>
                        )}
                    </Form.Item>

                    <Form.Item label='登录密码' required style={{marginBottom: 4}}>
                        {getFieldDecorator('pswd', {
                            rules: [
                                {required: true, message: '密码不可为空'},
                                {min: 6, message: '密码最少6位'},
                                {whitespace: true, message: '密码不可使用空白字符'}
                            ]
                        })(
                            <Input type='password' id='pswd' name='pswd'/>
                        )}
                    </Form.Item>
                </Card>

                <Row type='flex' justify='center'>
                    <Button type='primary' htmlType='submit'>下一步</Button>
                </Row>
            </Form>
        );
    });

    /**
     * 步骤列表
     */
    const steps = [
        {
            title: '基本配置',
            description: '配置所需的参数',
            content: <Configure />,
        },
        {
            title: '系统初始化',
            description: '应用配置',
            content: <Progress />,
        },
        {
            title: '默认用户',
            description: '创建默认超级管理员',
            content: <CreateAdmin />,
        },
        {
            title: '部署完成',
            description: null,
            content: <Row type='flex' justify='center' style={{marginTop: 60}}><Button onClick={() => location.href = '/'}>访问主页</Button></Row>,
        }
    ];

    return (
        <Layout style={{width: '100vw', height: '100vh', backgroundColor: 'rgba(0, 0, 0, .5)'}}>
            <Layout.Content style={{padding: 32}}>
                <p style={{marginBottom: 32, fontSize: '3.2em', fontWeight: 'bolder', color: 'rgb(82,82,82)', textAlign: 'center'}}>
                    系统安装部署
                </p>

                <Row type='flex' justify='center' style={{marginBottom: 32}}>
                    <Steps current={current} style={{width: 640}}>
                        {steps.map((step, index) => {
                            return <Steps.Step key={index} title={step.title} description={step.description}/>
                        })}
                    </Steps>
                </Row>                

                <Row type='flex' justify='center'>
                    <div style={{width: 600}}>
                        {steps[current].content}
                    </div>
                </Row>
            </Layout.Content>
        </Layout>
    )
}