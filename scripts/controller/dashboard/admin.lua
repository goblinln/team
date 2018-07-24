-- ========================================================
-- @File    : /controller/dashboard/admin.lua
-- @Brief   : 超级用户管理
-- @Author  : Leo Zhao
-- @Date    : 2018-07-24
-- ========================================================
local admin = {};

-- 主页
function admin:index(req, rsp)
    if not session.is_su then return rsp:error(403) end;
    rsp:html('dashboard/admin/index.html', {dashboard_menu = 'admin'});
end

-- 管理用户列表
function admin:users(req, rsp)
    if not session.is_su then return rsp:error(403) end;
    rsp:html('dashboard/admin/users.html', {
        users = M('user'):all()
    });
end

-- 添加用户
function admin:add_user(req, rsp)
    if not session.is_su then return rsp:error(403) end;
    if req.method ~= 'POST' then return rsp:error(405) end;

    local account = req.post.account;
    local name = req.post.name;
    local pswd = req.post.pswd;
    local cfm_pswd = req.post.cfm_pswd;
    local is_su = req.post.is_su and 1 or 0;
    
    if not account then
        return rsp:json{ok = false, err_msg = '帐号不能为空'};
    elseif not name then
        return rsp:json{ok = false, err_msg = '用户名不能为空'};
    elseif not pswd then
        return rsp:json{ok = false, err_msg = '密码不能为空'};
    elseif pswd ~= cfm_pswd then
        return rsp:json{ok = false, err_msg = '两次输入的密码不一致'};
    else
        local ok, err = M('user'):add(account, name, md5(pswd), is_su);
        rsp:json{ok = ok, err_msg = err};
    end
end

-- 编辑用户
function admin:edit_user(req, rsp)
    if not session.is_su then return rsp:error(403) end;
    if req.method ~= 'POST' then return rsp:error(405) end;

    local id = req.post.id;
    local account = req.post.account;
    local name = req.post.name;
    local is_su = req.post.is_su and 1 or 0;

    local ok, err = M('user'):edit(id, account, name, is_su);
    rsp:json{ok = ok, err_msg = err};
end

-- 删除用户
function admin:del_user(req, rsp)
    if not session.is_su then return rsp:error(403) end;
    if req.method ~= 'POST' then return rsp:error(405) end;

    local id = req.post.id;
    if not id then
        return rsp:json{ok = false, err_msg = '需要指定删除的角色'};
    else
        M('user'):del(id);
        return rsp:json{ok = true};
    end
end

return admin;