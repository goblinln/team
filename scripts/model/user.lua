-- ========================================================
-- @File    : /model/user.lua
-- @Brief   : 用户数据管理
-- @Author  : Leo Zhao
-- @Date    : 2018-07-16
-- ========================================================
local user = inherit(M('base'));

user.SECRET_AUTOLOGIN   = '@user.auto_login_secret_TEAM';

-- 登录
function user:login(req, rsp)
    local account   = req.post.account;
    local pswd      = md5(req.post.pswd);
    local remember  = req.post.remember_me;
    local find      = self:query("SELECT * FROM `users` WHERE `account`=?1", account);

    if #find == 0 then
        return false, '帐号不存在！';
    elseif find[1].pswd ~= pswd then
        return false, '帐号名或密码不正确！';
    else
        if remember then
            local expire    = 30 * 3600 * 24;
            local token     = b64.encode(json.encode{
                account = find[1].account,
                ip      = req.remote,
                sign    = md5(find[1].account .. '|' .. req.remote .. '|' .. self.SECRET_AUTOLOGIN),
            });
            token = string.gsub(token, '\n', '');
            rsp:cookie('login_token', token, expire, '/');            
            self:exec("UPDATE `users` SET auto_login_expire=?1 WHERE id=?2", os.time() + expire, find[1].id)
        end

        self:__on_login(find[1]);
        return true;
    end
end

-- 自动登录
function user:auto_login(req, rsp)
    if not req.cookie.login_token then return false end;

    local data  = json.decode(b64.decode(req.cookie.login_token));
    local ok    = false;

    if data.account and data.ip == req.remote then
        local sign = md5(data.account .. '|' .. req.remote .. '|' .. self.SECRET_AUTOLOGIN);
        if sign == data.sign then
            local find = self:query("SELECT * FROM `users` WHERE account=?1", data.account);
            if #find == 1 and find[1].auto_login_expire > os.time() then
                self:__on_login(find[1]);
                ok = true;
            end
        end
    end

    if not ok then rsp:cookie('login_token', '', -3600, '/') end;
    return ok;
end

-- 登出
function user:logout(req, rsp)
    if not session.uid then return end;

    self:exec("UPDATE `users` SET auto_login_expire=?1 WHERE id=?2", 0, session.uid);
    session = {};

    for k, _ in pairs(req.cookie) do
        rsp:cookie(k, '', -1);
    end
end

-- 取得指定ID用户的名字，如果没有参数，则取所有
function user:get_names(...)
    local np    = select('#', ...);
    local find  = {};
    local ret   = {};

    if np > 0 then
        find = self:query('SELECT `id`, `name` FROM `users` WHERE `id` IN(' .. table.concat({...}, ',') .. ')')    
    else    
        find = self:query([[SELECT `id`, `name` FROM `users`]]);
    end

    for _, info in ipairs(find) do
        ret[info.id] = info.name;
    end

    return ret;
end

-- 修改头像
function user:set_avatar(url)
    if not session.uid then return end;
    self:exec("UPDATE `users` SET avatar=?1 WHERE id=?2", url, session.uid);
end

-- 修改密码
function user:set_pswd(old, new)
    if not session.uid then return end;
    self:exec([[UPDATE `users` SET pswd=?1 WHERE id=?2 AND pswd=?3]], new, session.uid, old);
    if self:affected_rows() <= 0 then
        return false, '原始密码不正确！';
    else
        return true;
    end
end

-- 取出所有的用户
function user:all()
    local find = self:query([[SELECT `id`,`account`,`name`,`is_su` FROM `users`]]);
    return find;
end

-- 添加用户
function user:add(account, name, pswd, is_su)
    local ok, err = false, '';

    xpcall(function()
        self:exec([[
            INSERT INTO `users`(`account`, `name`, `pswd`, `is_su`)
            VALUES(?1, ?2, ?3, ?4)]], account, name, pswd, is_su);
        ok = true;
    end, function(stack)
        err = stack;
    end)

    return ok, err;
end

-- 编辑用户
function user:edit(id, account, name, is_su)
    local ok, err = false, '';

    xpcall(function()
        self:exec([[
            UPDATE `users`
            SET `account`=?1, `name`=?2, `is_su`=?3
            WHERE `id`=?4]], account, name, is_su, id);
        ok = true;
    end, function(stack)
        err = stack;
    end)

    return ok, err;
end

-- 删除用户
function user:del(id)
    self:exec([[DELETE FROM `users` WHERE `id`=?1]], id);
end

-----------------------------------------------------------

-- 登录完成后操作
function user:__on_login(him)
    session.uid     = him.id;
    session.name    = him.name;
    session.account = him.account;
    session.avatar  = him.avatar;
    session.is_su   = him.is_su == 1;
end

return user;
