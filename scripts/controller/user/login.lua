-- ========================================================
-- @File    : /controller/user/login.lua
-- @Brief   : 登录流程管理
-- @Author  : Leo Zhao
-- @Date    : 2018-07-16
-- ========================================================
local login = {};

-- 登录页
function login:index(req, rsp)
    rsp:html('user/login/index.html');
end

-- 登录处理
function login:do_login(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local ok, err = M('user'):login(req, rsp);
    if not ok then return rsp:json{ ok = false, err_msg = err } end;
    rsp:json{ ok = true };
end

-- 忘记密码
function login:forget_pswd(req, rsp)
    rsp:html('user/login/forget_pswd.html');
end

return login;