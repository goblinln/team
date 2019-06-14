import * as React from 'react';
import { WrappedFormInternalProps } from 'antd/lib/form/Form';

import {
    Button,
    Card,
    Checkbox,
    Form,
    Icon,
    Input,
    Layout,
    Row,
    message,
} from 'antd';

import { Fetch } from '../../common/Request';

/**
 * 登录页
 */
export const Login = () => {
    /**
     * 登录的表单组件
     */
    const LoginForm = Form.create()((props: WrappedFormInternalProps) => {
        const { getFieldDecorator, validateFields } = props.form;

        /**
         * 登录处理函数
         */
        const login = (ev: React.FormEvent<HTMLFormElement>) => {
            ev.preventDefault();

            validateFields(err => {
                if (!err) {
                    Fetch.post('/login', new FormData(ev.currentTarget), rsp => {
                        if (rsp.err) {
                            message.error(rsp.err, 1);
                        } else {
                            location.href = '/';
                        }
                    });
                }
            });
        }

        return (
            <Form onSubmit={login}>
                <Form.Item>
                    {getFieldDecorator('account', {
                        rules: [{required: true, message: '登录帐号不可为空'}]
                    })(<Input id='account' name='account' prefix={<Icon type='user' style={{color: 'rgba(0, 0, 0, .25)'}}/>} placeholder='登录帐号'/>)}
                </Form.Item>

                <Form.Item>
                    {getFieldDecorator('password', {
                        rules: [{required: true, message: '请输入登录密码'}, {whitespace: true, message: '密码不可使用空白字符'}]
                    })(<Input.Password id='password' name='password' type='password' prefix={<Icon type='lock' style={{color: 'rgba(0, 0, 0, .25)'}}/>} placeholder='登录密码'/>)}
                </Form.Item>

                <Form.Item style={{marginBottom: 4}}>
                    <Checkbox id='remember' name='remember' value="1">一个月内自动登录</Checkbox>
                </Form.Item>

                <Button type='primary' htmlType='submit' block>登录</Button>
            </Form>
        );
    });

    return (
        <Layout style={{width: '100vw', height: '100vh', backgroundColor: 'rgba(0, 0, 0, .5)'}}>
            <Layout.Content>
                <Row type='flex' justify='center' style={{marginTop: '10%', marginBottom: 16}}>
                    <span style={{fontSize: '3.2em', fontWeight: 'bolder', color: 'rgb(82,82,82)'}}>
                        团队协作平台
                    </span>
                </Row>

                <Row type='flex' justify='center'>
                    <Card style={{width: 360, textAlign: 'left', boxShadow: '0 .5rem 1rem rgba(0, 0, 0, .15)'}}>
                        <LoginForm />
                    </Card>
                </Row>
            </Layout.Content>
        </Layout>
    );
};
