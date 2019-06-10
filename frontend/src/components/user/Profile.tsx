import * as React from 'react';
import {WrappedFormUtils} from 'antd/lib/form/Form';

import {
    Avatar,
    Button,
    Card,
    Drawer,
    Empty,
    Form,
    Icon,
    Input,
    Row,
    Upload,
    message,
} from 'antd';

import { INotice } from '../../common/Protocol';
import { Fetch } from '../../common/Request';

/**
 * 个人信息页的可配置属性
 */
export interface IProps {
    /**
     * 用户名
     */
    name: string;
    /**
     * 帐号
     */
    account: string;
    /**
     * 当前的头像信息
     */
    avatar?: string;
    /**
     * 通知列表
     */
    notices: INotice[];
    /**
     * 关闭个人信息页的回调
     */
    onClose: (avatar: string, needUpdateNotice: boolean) => void;
}

/**
 * 个人信息弹出页
 */
export const View = (props: IProps) => {
    /**
     * 修改密码表单属性
     */
    interface IResetPswdFormProps {
        onFinish: () => void;
        form: WrappedFormUtils;
    }

    /**
     * 修改密码表单
     */
    const ResetPswdForm = Form.create<IResetPswdFormProps>()((props: IResetPswdFormProps) => {
        const {getFieldDecorator, getFieldValue, validateFields} = props.form;

        /**
         * 验证两次输入的新密码是否一致
         */
        const validateNewPswd = (rule: any, value: string, callback: (err?: string) => void) => {
            if (value && value !== getFieldValue('cfmPswd')) {
                callback('两次输入的密码不一致！');
            } else {
                callback();
            }
        };

        /**
         * 修改密码
         */
        const modifyPswd = (ev: React.FormEvent<HTMLFormElement>) => {
            ev.preventDefault();

            validateFields(err => {
                if (!err) {
                    Fetch.patch(`/api/user/pswd`, new FormData(ev.currentTarget), rsp => {
                        rsp.err ? message.error(rsp.err, 1) : props.onFinish();
                    });
                }
            });
        };

        return (
            <Form onSubmit={modifyPswd}>
                <Form.Item>
                    {getFieldDecorator('oldPswd', {
                        rules: [{required: true, message: '原始密码不可为空'}]
                    })(<Input.Password id='oldPswd' name='oldPswd' prefix={<Icon type='lock' style={{color: 'rgba(0, 0, 0, .25)'}} />} type='password' placeholder='原始密码'/>)}
                </Form.Item>

                <Form.Item>
                    {getFieldDecorator('newPswd', {
                        rules: [{required: true, message: '新密码不可为空'}, {whitespace: true, message: '密码不可有空格'}]
                    })(<Input.Password id='newPswd' name='newPswd' prefix={<Icon type='lock' style={{color: 'rgba(0, 0, 0, .25)'}} />} type='password' placeholder='新的密码'/>)}
                </Form.Item>

                <Form.Item>
                    {getFieldDecorator('cfmPswd', {
                        rules: [{required: true, message: '请再次输入新密码'}, {validator: validateNewPswd}]
                    })(<Input.Password id='cfmPswd' name='cfmPswd' prefix={<Icon type='lock' style={{color: 'rgba(0, 0, 0, .25)'}} />} type='password' placeholder='确认新密码'/>)}
                </Form.Item>

                <Button type='primary' htmlType='submit' block>提交</Button>
            </Form>
        );
    });

    /**
     * 状态机
     */
    const [avatar, setAvatar] = React.useState<string>(props.avatar);
    const [isResetPassword, setIsResetPassword] = React.useState<boolean>(false);
    const [needUpdateNotice, setNeedUpdateNotice] = React.useState<boolean>(false);

    /**
     * 修改头像
     */
    const modifyAvatar = (image: File) => {
        let param = new FormData();
        param.append('img', image, image.name);

        Fetch.patch(`/api/user/avatar`, param, rsp => {
            rsp.err ? message.error(rsp.err, 1) : setAvatar(rsp.data);
        });

        return false;
    };

    /**
     * 删除一条信息
     */
    const delNoticeOne = (id: number) => {
        Fetch.delete(`/api/notice/${id}`, rsp => {
            rsp.err ? message.error(rsp.err) : setNeedUpdateNotice(true);
        });
    };

    /**
     * 删除全部消息
     */
    const delNoticeAll = () => {
        Fetch.delete(`/api/notice/all`, rsp => {
            rsp.err ? message.error(rsp.err) : setNeedUpdateNotice(true);
        });
    };

    return (
        <Drawer title='用户信息' closable={true} width={350} visible={true} onClose={() => props.onClose(avatar, needUpdateNotice)}>
            <div style={{textAlign: 'center'}}>
                <Row type='flex' justify='center'>
                    <Upload name='avatar' listType='picture-card' showUploadList={false} accept='image/*' action='/file/upload' beforeUpload={modifyAvatar}>
                        <Avatar icon='user' size={80} src={avatar} style={{marginBottom: '.2em'}} />
                        点击修改
                    </Upload>
                </Row>

                <p style={{fontSize: '2em', fontWeight: 'bolder', marginBottom: 0}}>{props.name}</p>
                <p>{props.account}</p>

                <Button size='small' onClick={() => setIsResetPassword(true)}>修改密码</Button>
            </div>

            <Card title='消息列表' extra={<Button size='small' type='link' onClick={() => delNoticeAll()}>清空</Button>} bodyStyle={{padding: 0}} style={{marginTop: 16}}>
                {props.notices.length == 0 ? <Empty description='暂无数据'/> : (
                    <ul style={{margin: '0 8px', paddingInlineStart: 12}}>
                        {props.notices.map(notice => {
                            return (
                                <li key={notice.id} style={{margin: 4}}>
                                    <p style={{marginBottom: 2}}>{notice.content}</p>
                                    <Row type='flex' justify='space-between' align='middle'>                                    
                                        <small><Icon type='calendar' /> {notice.time}</small>
                                        <small><a onClick={() => delNoticeOne(notice.id)}>删除</a></small>
                                    </Row>
                                </li>
                            );
                        })}
                    </ul>                    
                )}
            </Card>

            <Drawer title='修改密码' closable={true} width={200} visible={isResetPassword} onClose={() => setIsResetPassword(prev => !prev)} >
                <ResetPswdForm onFinish={() => setIsResetPassword(false)}/>
            </Drawer>
        </Drawer>
    );
};
