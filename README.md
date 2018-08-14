# 项目任务管理系统

一个模仿[Team@OSC](https://team.oschina.net)的`radmine`?

## 重要提醒

1. 代码库中没有`omni`的可执行文件，请到**发行版**中下载带预编译的版本，或自行到[OmniWeb](https://gitee.com/love_linger/OmniWeb.git)中下载编译

2. 现在有些属性只有项目管理员可以更改

    * 项目成员管理
    * 更改已存在任务的时间计划（其他成员只能通知项目管理员协调更改时间，自己不能随意更改）
    * 项目视图中的任务编辑
    * 任务的验收管理

3. 关于验收

    * 测试通过后状态为【已完成】的任务才会出现在【验收管理】中。
    * 验收操作也是归档操作。验收后，任务不再出现在任务视图中，可在项目验收管理及周报中查看。

## 更新日志

* 2018/08/14  
    1. 调整页面布局
    2. Markdown编辑器定制
        * 禁用快捷键，以支持tab输入
        * 禁止Shift, Ctrl, Atl等键按起时乱补全
    3. 上传新头像时，删除旧头像
    4. 修复内存泄露的BUG：主要原因是mysql_close()未能完全释放占用资源，需要调用mysql_server_end()
    5. 修复Gantt图未对齐BUG

* 2018/08/10
    1. 多个页面布局的修改
    2. 增加错误静态页面
    3. 修复SQL在旧版本的MySQL(<5.6)中执行错误

* 2018/08/06

    1. Markdown编辑器增加语法说明
    2. 一些视图Layout修改
    3. 任务现在支持上传附件
    4. 使用最新的omni 4.4.2

* 2018/08/01

    1. 修复SQL语句中不能使用超过9个参数的BUG
    2. 重新设计任务的状态，新增加`测试中`状态
    3. 任务常规视图中，不再显示【已滞后】列表，改为截止时间显示红字
    4. 在任务表中增加`cooperator`用于指定测试/验收人员
    5. 重新定制`jQuery.GanttView`用于支持显示测试验收人员 

## 预览

* 任务视图 - 【看板】模式
![Tasks](/preview/preview.png)

* 任务视图 - 【甘特图】模式
![Gantt](/preview/gantt.png)

* 任务视图 - 任务详情
![TaskInfo](/preview/task.png)

* 任务视图 - 发布任务预览
![Publish](/preview/publish.png)

* 项目视图 - 项目周报
![Reports](/preview/reports.png)

* 项目视图 - 验收管理
![Archive](/preview/archive.png)

* 项目视图 - 成员管理
![Members](/preview/members.png)

* 文档视图
![Documents](/preview/documents.png)

* 文件分享
![Share](/preview/share.png)

## 使用技术

* [OmniWeb](https://gitee.com/love_linger/OmniWeb.git)
* [JQuery](https://jquery.com)
* [JQuery Datetimepicker](https://github.com/xdan/datetimepicker)
* [jsTree](https://www.jstree.com)
* [Bootstrap 4](http://getbootstrap.com/)
* [Bootstrap Markdown](https://github.com/toopay/bootstrap-markdown)
* [FontAwesome 4](http://www.fontawesome.com.cn/)
* [Gantt Chart](https://github.com/982964399/jQuery-ganttView)

> 【注】`Bootstrap markdown`是个人定制版，集成了`markedjs`, `flowchart.js`, `katexjs`。

## 安装

发行版中提供了基于`Window 10`及`CentOS 7`的预编译版，直接下载，配置完`omni.ini`中数据库，运行`omni`或`omni.exe`即可。

## 默认帐号

**帐号:** admin  
**密码:** team

>【注】可修改`scripts/controller/install/setup.lua`中`setup.build_in`配置，自定义默认帐号