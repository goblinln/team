-- ========================================================
-- @File    : /controller/user/profile.lua
-- @Brief   : 用户信息管理
-- @Author  : Leo Zhao
-- @Date    : 2018-07-24
-- ========================================================
local profile = {};

-- 主页
function profile:index(req, rsp)
    rsp:html('user/profile/index.html');
end

-- 修改名字
function profile:set_name(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;
    if req.post.name == session.name then return rsp:json{ ok = false } end;
    local ok, err = M('user'):set_name(req.post.name);
    rsp:json{ok = ok, err_msg = err};
end

-- 修改头像
function profile:set_avatar(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local dir = 'www/upload';
    if not os.exists(dir) then os.mkdir(dir) end;

    dir = dir .. '/' .. session.uid;  
    if not os.exists(dir) then os.mkdir(dir) end;

    local to = dir .. '/' .. os.time() .. '_' .. req.post.img;
    local from = req.file[req.post.img];
    if not os.cp(from, to) then return rsp:error(500) end;

    local url = '/' .. to;
    
    M('user'):set_avatar(url);
    rsp:json({ ok = true, url = url });
end

-- 修改密码
function profile:set_pswd(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    if req.post.new_pswd ~= req.post.cfm_pswd then
        return rsp:json{ ok = false, err_msg = '两次输入的密码不一致！' };
    end

    if not req.post.new_pswd then
        return rsp:json{ ok = false, err_msg = '新密码不能为空！' };
    end

    if req.post.new_pswd == req.post.org_pswd then
        return rsp:json{ ok = true };
    end

    local ok, err = M('user'):set_pswd(md5(req.post.org_pswd), md5(req.post.new_pswd));
    rsp:json{ok = ok, err_msg = err};
end

return profile;