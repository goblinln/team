-- ========================================================
-- @File    : /model/authorization.lua
-- @Brief   : 权限管理模块
-- @Author  : Leo Zhao
-- @Date    : 2018-07-16
-- ========================================================
local authorization = inherit(M('base'));

-- 权限管理入口
function authorization:check(req, rsp)    
    -- 检测是否需要安装
    if req.controller == 'install/setup' then
        if os.exists(config.app.install_lock) then
            rsp:redirect('/user/login/');
            return false;
        else
            return true;
        end
    elseif not os.exists(config.app.install_lock) then
        rsp:redirect('/install/setup/?override=1');
        return false;
    end

    -- 检测是否需要登录
    if not session.uid then
		local ok = M('user'):auto_login(req, rsp);
		if ok then
			if req.controller == 'user/login' then
                rsp:redirect('/dashboard/tasks/');
                return false;
			end
		elseif req.controller ~= 'user/login' then
            rsp:redirect('/user/login/');
            return false;
		end
	elseif req.controller == 'user/login' then
        rsp:redirect('/dashboard/tasks/');
        return false;
    end

    return true;
end

return authorization;