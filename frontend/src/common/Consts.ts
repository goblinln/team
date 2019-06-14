/**
 * 用户在项目中的角色
 */
export const ProjectRole = [
    '策划',
    '研发',
    '测试',
    '运营',
    '美术',
];

/**
 * 任务状态定义
 */
export const TaskStatus = [
    {
        type: 'created',
        name: '待办中',
        icon: 'schedule',
        color: '#6c757d',
    },
    {
        type: 'underway',
        name: '进行中',
        icon: 'build',
        color: '#17a2b8',
    },
    {
        type: 'testing',
        name: '测试中',
        icon: 'experiment',
        color: '#007bff',
    },
    {
        type: 'finished',
        name: '已完成',
        icon: 'check-circle',
        color: '#28a745',
    },
    {
        type: 'archived',
        name: '已验收',
        icon: 'file-done',
        color: '#28a745',
    }
];

/**
 * 任务的权重定义
 */
export const TaskWeight = [
    { name: '[一般] ', color: 'green' },
    { name: '[次要] ', color: 'blue' },
    { name: '[主要] ', color: 'orange' },
    { name: '[严重] ', color: '#f50' },
];

/**
 * 任务标签定义
 */
export const TaskTag = [
    { name: '缺陷', color: '#f50' },
    { name: '快速修正', color: '#2db7f5' },
    { name: '优化', color: '#87d068' },
    { name: '新功能', color: '#108ee9' },
];
