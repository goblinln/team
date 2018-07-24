# 项目任务管理系统

一个模仿[Team@OSC](https://team.oschina.net)的`radmine`?

## 预览

* 【看板】模式
![Tasks](/preview/preview.png)

* 【甘特图】模式
![Gantt](/preview/gantt.png)

* 发布任务预览
![Publish](/preview/publish.png)

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
2. 修改omni.ini中数据库配置
3. 运行omni.exe。需要Linux版的，请自行到`OmniWeb`中下载预编译好的可执行文件

## 默认帐号

**帐号:** admin  
**密码:** team

>【注】可修改`scripts/controller/install/setup.lua`中`setup.build_in`配置，自定义默认帐号