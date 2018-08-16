-- ========================================================
-- @File    : /controller/user/notification.lua
-- @Brief   : 消息通知
-- @Author  : Leo Zhao
-- @Date    : 2018-08-16
-- ========================================================
local notification = {};

-- 主页
function notification:index(req, rsp)
    rsp:html('user/notification/index.html', {
        messages = M('notification'):all()
    });
end

-- 获取通知数量
function notification:count(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;
    rsp:json{ count = M('notification'):count() };
end

-- 标记一条消息已读
function notification:read(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;
    
    M('notification'):del(req.post.id);
    rsp:json{ ok = true };
end

-- 标记所有消息已读
function notification:read_all(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;
    
    M('notification'):clear();
    rsp:json{ ok = true };
end

return notification;