# 项目任务管理系统

一个模仿[Team@OSC](https://team.oschina.net)的`radmine`?

## 重要提醒

1. 因项目在开发期，所以数据表经常变动，升级时请比对`scripts/controller/install/setup.lua`文件，如果表结构发生变化，自行决定升级数据库方式

    * 方式一：手动删除所有表，删除`omni.lock`文件，重启即可。
    * 方式二：比对文件，自行更改表结构

2. 现在有些属性只有项目管理员可以更改

    * 已存在任务的时间计划
    * 项目视图中的任务列表
    * 任务的归档管理

3. 验收操作也是归档操作！

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

## 使用技术

* [OmniWeb](https://gitee.com/love_linger/OmniWeb.git)
* [JQuery](https://jquery.com)
* [JQuery Form](http://plugins.jquery.com/form/)
* [JQuery Datetimepicker](https://github.com/xdan/datetimepicker)
* [Bootstrap 4](http://getbootstrap.com/)
* [Bootstrap Markdown](https://github.com/toopay/bootstrap-markdown)
* [FontAwesome 4](http://www.fontawesome.com.cn/)
* [Gantt Chart](https://github.com/982964399/jQuery-ganttView)
* [jsTree](https://www.jstree.com)

> 【注】`Bootstrap markdown`是个人定制版，集成了`markedjs`, `mermaidjs`, `katexjs`。

## 安装

1. Clone
2. 修改omni.ini中数据库配置。
3. 运行omni.exe。需要Linux版的，请自行到`OmniWeb`中下载并编译

## 默认帐号

**帐号:** admin  
**密码:** team

>【注】可修改`scripts/controller/install/setup.lua`中`setup.build_in`配置，自定义默认帐号