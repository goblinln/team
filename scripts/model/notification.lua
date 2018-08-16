-- ========================================================
-- @File    : /model/notification.lua
-- @Brief   : 消息通知
-- @Author  : Leo Zhao
-- @Date    : 2018-08-16
-- ========================================================
local m = inherit(M('base'));

-- 取得一个人的通知数量
function m:count()
    local data = self:query([[SELECT COUNT(*) AS num FROM `notifications` WHERE `uid`=?1]], session.uid)[1];
    if not data then return 0 end;
    return data.num;
end

-- 取得一个人的所有的消息通知
function m:all()
    return self:query([[SELECT * FROM `notifications` WHERE `uid`=?1]], session.uid);
end

-- 删除一条消息
function m:del(id)
    self:exec([[DELETE FROM `notifications` WHERE `id`=?1 AND `uid`=?2]], id, session.uid);
end

-- 删除全部消息
function m:clear()
    self:exec([[DELETE FROM `notifications` WHERE `uid`=?1]], session.uid);
end

-- 发送通知
function m:add(uid, message)
    self:exec([[INSERT INTO `notifications`(`uid`, `message`) VALUES(?1, ?2)]], uid, message);
end

return m;