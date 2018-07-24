-- ===========================================================
-- @File	: app.lua
-- @Brief	: 简单的MVC框架
-- @Author	: Leo Zhao
-- @Date	: 2017-05-03
-- ===========================================================

--------------------- 全局工具函数 ----------------------------

-- 取Model
function M(class)
	return require(config.dirs.model .. class);
end

-- 取Controller
function C(class)
	local c = config.dirs.controller .. class;
	if not os.exists(c .. '.lua') then return end;
	return require(c);
end

-- 取Vender
function T(class)
	return require(config.dirs.vendor .. class);
end

------------------------- 路由定义 ----------------------------

-- 入口
local entry = function(controller, action, req, rsp)
	template.template_root = config.dirs.view;

	req.controller	= controller;
	req.action		= action;

	-- 调用权限管理
	if not M('authorization'):check(req, rsp) then
		return;
	end

	-- 检测URL是否合法
	local tb = C(controller);
	if not tb or not tb[action] then return rsp:error(404) end;

	-- 调用
	xpcall(function()
		tb[action](tb, req, rsp);
	end, function(err)
		log.error(err);
		rsp:error(500);
	end);
end

-- 主页
router:get('^/$', function(req, rsp)
	entry('dashboard/tasks', 'index', req, rsp);
end);

-- MVC
router:any('^/([A-Za-z][A-Za-z0-9_/]*)/([A-Za-z][A-Za-z0-9_]*)$', function(req, rsp, controller, action)
	entry(controller, action, req, rsp);
end);

-- MVC（默认ACTION）
router:any('^/([A-Za-z][A-Za-z0-9_]*)$', function(req, rsp, controller)
	entry(controller, 'index', req, rsp);
end);

-- MVC（默认ACTION，多层controller）
router:any('^/([A-Za-z][A-Za-z0-9_/]*)/$', function(req, rsp, controller)
	entry(controller, 'index', req, rsp);
end);
