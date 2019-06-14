/**
 * 用户信息
 */
export interface IUser {
    /**
     * 唯一ID
     */
    id: number,
    /**
     * 帐号
     */
    account?: string,
    /**
     * 显示名称
     */
    name?: string,
    /**
     * 头像
     */
    avatar?: string,
    /**
     * 是否是超级管理员
     */
    isSu?: boolean,
    /**
     * 是否被锁定登录
     */
    isLocked?: boolean,
}

/**
 * 项目成员
 */
export interface IProjectMember {
    /**
     * 成员信息
     */
    user: IUser;
    /**
     * 成员角色，参见`ProjectRole`
     */
    role: number;
    /**
     * 是否是项目管理员
     */
    isAdmin?: boolean;
}

/**
 * 项目信息
 */
export interface IProject {
    /**
     * 唯一ID
     */
    id: number;
    /**
     * 显示名称
     */
    name: string;
    /**
     * 成员列表
     */
    members?: IProjectMember[];
    /**
     * 分支列表
     */
    branches?: string[];
}

/**
 * 评论
 */
export interface ITaskComment {
    /**
     * 时间
     */
    time: string;
    /**
     * 人
     */
    user: string,
    /**
     * 头像
     */
    avatar: string;
    /**
     * 内容
     */
    content: string;
}

/**
 * 任务相关操作事件
 */
export interface ITaskEvent {
    /**
     * 时间
     */
    time: string;
    /**
     * 操作人员
     */
    operator: string;
    /**
     * 事件
     */
    event: number;
    /**
     * 事件附加参数
     */
    extra: string;
}

/**
 * 服务器返回的任务数据类型
 */
export interface ITask {
    /**
     * 任务唯一ID
     */
    id: number;

    /**
     * 任务标题
     */
    name: string;

    /**
     * 所属项目
     */
    proj: IProject,

    /**
     * 所属分支
     */
    branch: number,

    /**
     * 是否置顶
     */
    bringTop?: boolean;

    /**
     * 权重
     */
    weight: number;

    /**
     * 当前的状态
     */
    state: number;

    /**
     * 创建者/需求发起方
     */
    creator: IUser,

    /**
     * 开发者/乙方
     */
    developer: IUser,

    /**
     * 测试者/验收方
     */
    tester: IUser,

    /**
     * 计划开始时间
     */
    startTime: string;

    /**
     * 计划截止时间
     */
    endTime: string;

    /**
     * 任务标签列表
     */
    tags?: number[];

    /**
     * 任务内容
     */
    content?: string;

    /**
     * 任务附件列表
     */
    attachments?: {name: string, url: string}[];

    /**
     * 评论
     */
    comments?: ITaskComment[];

    /**
     * 事件
     */
    events?: ITaskEvent[];
}

/**
 * 通知事件
 */
export interface INotice {
    /**
     * 唯一ID
     */
    id: number;
    /**
     * 相关任务ID
     */
    tid: number;
    /**
     * 相关任务名
     */
    tname: string;
    /**
     * 相关操作人员
     */
    operator: string;
    /**
     * 时间
     */
    time: string;
    /**
     * 事件
     */
    ev: number;
}

/**
 * WIKI文档
 */
export interface IDocument {
    /**
     * 唯一ID
     */
    id: number;
    /**
     * 父节点ID
     */
    parent: number;
    /**
     * 标题
     */
    title: string;
    /**
     * 创建者
     */
    creator?: string;
    /**
     * 最近更新人
     */
    modifier?: string;
    /**
     * 最近更新时间
     */
    time?: string;
    /**
     * 内容
     */
    content?: string;
}

/**
 * 分享文件
 */
export interface IShare {
    /**
     * 唯一ID
     */
    id: number;
    /**
     * 文件名
     */
    name: string;
    /**
     * 上传者
     */
    uploader: string;
    /**
     * 时间
     */
    time: string;
    /**
     * 大小
     */
    size: number;
}