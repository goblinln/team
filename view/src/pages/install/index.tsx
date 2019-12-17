import * as React from 'react';

import {Button, Card, Form, Input, Row, Steps} from '../../components';
import {request} from '../../common/request';

export const Install = () => {
    const [step, setStep] = React.useState<number>(0);
    const [dbInited, setDBInited] = React.useState<boolean>(false);
    const [dbStatus, setDBStatus] = React.useState<string[]>([]);
    const dbTimer = React.useRef<number>(-1);

    const configureForm = Form.useForm({
        port: {required: '端口号不可为空', range: {min: 80, max: 65535, message: '端口号需要在80至65535之间'}},
        mysqlHost: {required: 'MySQL服务器的地址不可为空'},
        mysqlUser: {required: 'MySQL登录用户未填写'},
        mysqlDB: {required: '数据未填写'},
    });

    const suForm = Form.useForm({
        account: {required: '登录帐号不可为空'},
        name: {required: '登录用户显示名称未填写'},
        pswd: {required: '密码不可为空', length: {min: 6, max: 512, message: '密码最少6位'}, pattern: {test: /[\w_\d]+/, message: '密码只能包含字母、数字和下划线'}},
    });

    const goBack = (ev: React.MouseEvent<HTMLButtonElement>) => {
        ev.preventDefault();
        if (step == 0) return;
        if (step <= 1) setDBInited(false);
        setStep(prev => prev - 1); 
    };

    const goNext = () => {
        if (step == 0) {
            startDBTimer();
            setDBStatus([]);
            setDBInited(false);
        } else if (step == 1) {
            clearDBTimer();
        }

        setStep(prev => prev + 1);
    };

    const submitConfigure = (ev: React.FormEvent<HTMLFormElement>) => {
        ev.preventDefault();

        request({
            url: '/install/configure',
            method: 'POST',
            data: new FormData(ev.currentTarget),
            success: goNext,
        });
    };

    const fetchDBStatus = () => {
        request({
            url: '/install/status',
            dontShowLoading: true,
            success: (data: any) => {
                setDBStatus(data.status);
                if (data.done) {
                    clearDBTimer();
                    setDBInited(true);
                } else if (data.isError) {
                    clearDBTimer();
                }
            }
        });
    };

    const startDBTimer = () => {
        fetchDBStatus();
        dbTimer.current = window.setInterval(() => {
            fetchDBStatus();
        }, 500);
    };

    const clearDBTimer = () => {
        if (dbTimer.current >= 0) {
            clearInterval(dbTimer.current);
            dbTimer.current = -1;
        }
    };

    const submitSuperUser = (ev: React.FormEvent<HTMLFormElement>) => {
        ev.preventDefault();
        
        request({
            url: '/install/admin',
            method: 'POST',
            data: new FormData(ev.currentTarget),
            success: goNext,
        });
    };

    return (
        <div className='fullscreen pt-4 bg-light' style={{display: 'flex', justifyContent: 'center'}}>
            <div>
                <p className='text-logo fg-muted'>系统初始化配置</p>

                <Steps current={step}>
                    <Steps.Step label='基本配置'>
                        <Form form={configureForm} style={{width: 400}} onSubmit={submitConfigure}>
                            <Card header='监听端口（需要重新启动）' bordered>
                                <Form.Field htmlFor='port'>
                                    <Input name='port' value='8080'/>
                                </Form.Field>
                            </Card>

                            <Card className='mt-2' header='MySQL配置' bordered>
                                <Form.Field htmlFor='mysqlHost' label='服务器地址'>
                                    <Input name='mysqlHost' value='127.0.0.1:3306'/>
                                </Form.Field>

                                <Form.Field htmlFor='mysqlUser' label='登录用户'>
                                    <Input name='mysqlUser' value='root'/>
                                </Form.Field>

                                <Form.Field htmlFor='mysqlPswd' label='登录密码'>
                                    <Input.Password name='mysqlPswd'/>
                                </Form.Field>

                                <Form.Field htmlFor='mysqlDB' label='使用的数据库'>
                                    <Input name='mysqlDB' value='team'/>
                                </Form.Field>
                            </Card>

                            <Row className='mt-2' flex={{align: 'middle', justify: 'center'}}>
                                <Button theme='primary' size='sm' onClick={ev => {ev.preventDefault(); configureForm.submit()}}>下一步</Button>
                            </Row>
                        </Form>
                    </Steps.Step>

                    <Steps.Step label='系统初始化'>
                        <Card header='进度信息' style={{width: 400}} bordered>
                            <ul style={{marginLeft: 20}}>
                                {dbStatus.map((msg, idx) => <li key={idx}>{msg}</li>)}
                            </ul>
                        </Card>

                        <Row className='mt-2' flex={{align: 'middle', justify: 'center'}}>
                            <Button size='sm' onClick={goBack}>上一步</Button>
                            <Button theme='primary' size='sm' onClick={() => goNext()} disabled={!dbInited}>下一步</Button>
                        </Row>
                    </Steps.Step>

                    <Steps.Step label='超级用户'>
                        <Form form={suForm} style={{width: 400}} onSubmit={submitSuperUser}>
                            <Card header='默认管理员' bordered>
                                <Form.Field htmlFor='account' label='帐号'>
                                    <Input name='account' value='admin'/>
                                </Form.Field>

                                <Form.Field htmlFor='name' label='昵称'>
                                    <Input name='name' value='超级管理员'/>
                                </Form.Field>

                                <Form.Field htmlFor='pswd' label='登录密码'>
                                    <Input.Password name='pswd'/>
                                </Form.Field>
                            </Card>

                            <Row className='mt-2' flex={{align: 'middle', justify: 'center'}}>
                                <Button size='sm' onClick={goBack}>上一步</Button>
                                <Button theme='primary' size='sm' onClick={ev => {ev.preventDefault(); suForm.submit()}}>下一步</Button>
                            </Row>
                        </Form>
                    </Steps.Step>

                    <Steps.Step label='部署完成'>
                        <Button theme='primary' size='sm' onClick={() => location.href = '/'}>访问网站</Button>
                    </Steps.Step>
                </Steps>
            </div>
        </div>
    );
};