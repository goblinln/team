-- ========================================================
-- @File    : /controller/user/logout.lua
-- @Brief   : 登出流程管理
-- @Author  : Leo Zhao
-- @Date    : 2018-07-16
-- ========================================================
local logout = {};

function logout:index(req, rsp)
    M('user'):logout(req, rsp);
    rsp:redirect('/user/login/');
end

return logout;