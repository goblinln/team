-- ========================================================
-- @File    : /model/projects.lua
-- @Brief   : 项目数据管理
-- @Author  : Leo Zhao
-- @Date    : 2018-07-19
-- ========================================================
local projects = inherit(M('base'));

-- 取得当前用户参与的项目列表
function projects:all_of_mine()
    local find = self:query([[
        SELECT * 
        FROM `projects` INNER JOIN `project_members`
        WHERE projects.id = project_members.pid AND project_members.uid = ?1]],
        session.uid);

    return find;
end

-- 创建项目
function projects:add(name, role, repo)
    local ok, err = false, '';

    xpcall(function()
        self:exec('INSERT INTO `projects`(`name`, `repo`) VALUES(?1, ?2)', name, repo or '');
        if self:affected_rows() == 0 then
            err = '未知原因';
            return;
        end

        self:exec('INSERT INTO `project_members`(`uid`, `pid`, `role`, `is_admin`) VALUES(?1, ?2, ?3, 1)', session.uid, self:last_id(), role or 1);
        if self:affected_rows() == 0 then
            err = '未知原因';
        else
            ok  = true;
        end
    end, function(stack)
        err = stack;
    end);

    return ok, err;
end

-- 检测是否是项目的管理员
function projects:is_admin(pid)
    local find = self:query([[
        SELECT `is_admin` 
        FROM `project_members`
        WHERE pid = ?1 AND uid = ?2]],
        pid, session.uid);

    return #find > 0 and find[1].is_admin == 1;
end

-- 删除项目
function projects:delete(pid)
    local ok, err = false, '';

    xpcall(function()
        self:exec('DELETE FROM `projects` WHERE `id`=?1', pid);
        if self:affected_rows() == 0 then
            err = '项目不存在！';
            return;
        end

        self:exec('DELETE FROM `project_members` WHERE `pid`=?1', pid);
        self:exec('DELETE FROM `tasks` WHERE `pid`=?1', pid);
        ok = true;
    end, function(stack)
        err = stack;
    end);

    return ok, err;
end

-- 添加成员
function projects:add_member(pid, uid, role, is_admin)
    local ok, err = false, '';

    xpcall(function()
        self:exec('INSERT INTO `project_members`(`uid`, `pid`, `role`, `is_admin`) VALUES(?1, ?2, ?3, ?4)', uid, pid, role, is_admin and 1 or 0);
        ok = true;
    end, function(stack)
        err = stack;
    end);

    return ok, err;
end

-- 取得项目的成员
function projects:get_members(pid)
    local find  = self:query([[SELECT * FROM `project_members` WHERE `pid`=?1]], pid);
    local names = M('user'):get_names();

    for _, info in ipairs(find) do
        info.user = names[info.uid] or '神秘人';
        info.user_role = C('dashboard/projects').roles[info.role];
    end

    return find;
end

-- 取得可邀请的成员
function projects:get_invite_users(pid)
    local find  = self:query([[SELECT `uid` FROM `project_members` WHERE `pid`=?1]], pid);
    local names = M('user'):get_names();
    local ret   = {};
    local map   = {};

    for _, info in ipairs(find) do
        map[info.uid] = true;
    end

    for id, name in pairs(names) do
        if not map[id] then table.insert(ret, {uid = id, uname=name}) end;
    end
    
    if #ret == 0 then return false, '没有可邀请的成员！' end;
    return ret;
end

-- 修改成员信息
function projects:edit_member(pid, uid, role, is_admin)
    local ok, err = false, '';

    xpcall(function()
        self:exec('UPDATE `project_members` SET `role`=?1, `is_admin`=?2 WHERE `uid`=?3 AND `pid`=?4', role, is_admin and 1 or 0, uid, pid);
        ok = true;
    end, function(stack)
        err = stack;
    end);

    return ok, err;
end

-- 移除项目成员
function projects:del_member(pid, uid)
    if uid == session.uid then return false, '您不能移除自己，请将管理权限移交，再由其他人移除您！' end;

    local ok, err = false, '';

    xpcall(function()
        self:exec('DELETE FROM `project_members` WHERE `uid`=?1 AND `pid`=?2', uid, pid);
        ok = true;
    end, function(stack)
        err = stack;
    end);

    return ok, err;
end

-- 生成周报
function projects:get_reports(pid, week_offset)
    local report = {};
    local timepoint = week_offset * 3600 * 24 * 7 + os.time();
    local calc_start = os.date('*t', timepoint);

    -- 偏移到星期日
    if calc_start.wday ~= 1 then
        local sunday = timepoint - (calc_start.wday - 1) * 3600 * 24;
        calc_start = os.date('*t', sunday);
    end

    calc_start.hour = 0;
    calc_start.min = 0;
    calc_start.sec = 0;

    local week_start = os.time(calc_start);
    local week_end = week_start + 3600 * 24 * 7;

    report.week = os.date('%U', week_start);
    report.week_offset = week_offset;
    report.week_start = os.date('%Y/%m/%d', week_start);
    report.week_end = os.date('%Y/%m/%d', week_end);
    
    M('tasks'):report_for_proj(pid, week_start, week_end, report);
    return report;
end

return projects;