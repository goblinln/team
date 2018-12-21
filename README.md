# 项目任务管理系统

一个模仿[Team@OSC](https://team.oschina.net)的`radmine`?

## 重要提醒

1. 代码库中没有`omni`的可执行文件，请到[OmniWeb](https://gitee.com/love_linger/OmniWeb.git)中下载编译好的二进制文件

2. 依赖mysql库  

    * Windows：下载[Mysql Connector C](http://iso.mirrors.ustc.edu.cn/mysql-ftp/Downloads/Connector-C/mysql-connector-c-6.1.11-winx64.zip)，将lib目录下的libmysql.dll，放在项目根目录，并重命名为mysql.dll
    * CentOS: `yum install mariadb-devel`, `ln -s /usr/lib64/mysql/libmysqlclient.so /usr/lib64/libmysqlclient.so`
    * MacOS: `brew install libmysqlclient`

2. 现在有些属性只有项目管理员可以更改

    * 项目成员管理
    * 更改已存在任务的时间计划（其他成员只能通知项目管理员协调更改时间，自己不能随意更改）
    * 任务的验收管理

3. 关于验收

    * 测试通过后状态为【已完成】的任务才会出现在【验收管理】中。
    * 验收操作也是归档操作。验收后，任务不再出现在任务视图中，可在项目验收管理及周报中查看。 

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
* [JQuery TOC](http://github.com/ndabas/toc)
* [Bootstrap 4](http://getbootstrap.com/)
* [Bootstrap Markdown](https://github.com/toopay/bootstrap-markdown)
* [Bootstrap Datepicker](https://github.com/uxsolutions/bootstrap-datepicker)
* [Bootstrap Year Calendar](https://github.com/Paul-DS/bootstrap-year-calendar)
* [Bootstrap Notify](https://github.com/mouse0270/bootstrap-notify)
* [FontAwesome 4](http://www.fontawesome.com.cn/)
* [Gantt Chart](https://github.com/982964399/jQuery-ganttView)
* [jsTree](https://www.jstree.com)

> 【注】`Bootstrap markdown`是个人定制版，集成了`markedjs`, `flowchart.js`, `katexjs`。

## 安装

1. Clone代码库
2. 将下载omni可执行文件到项目根目录
3. 根据平台，提供mysql动态库，供脚本调用，方法见 #重要提示
4. 更改omni.ini，配置数据库
5. 运行omni

## 默认帐号

**帐号:** admin  
**密码:** team

>【注】可修改`scripts/controller/install/setup.lua`中`setup.build_in`配置，自定义默认帐号