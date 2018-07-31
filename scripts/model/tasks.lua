-- ========================================================
-- @File    : /model/tasks.lua
-- @Brief   : 任务管理
-- @Author  : Leo Zhao
-- @Date    : 2018-07-20
-- ========================================================
local tasks = inherit(M('base'));

-- 取得一个项目的任务
function tasks:get_by_proj(pid)
    local find  = self:query([[
            SELECT `id`, `creator`, `assigned`, `name`, `weight`, `tags`, `start_time`, `end_time`, `status`
            FROM `tasks`
            WHERE `pid`=?1 AND `status`<>?2]],
            pid, C('dashboard/tasks').status.ARCHIVED);
            
    self:__on_loaded(find);
    return find;
end

-- 取得一个项目可归档的任务列表
function tasks:get_archivable_by_proj(pid)
    local find  = self:query([[
            SELECT `id`, `creator`, `assigned`, `name`, `weight`, `tags`, `start_time`, `end_time`, `status`
            FROM `tasks`
            WHERE `pid`=?1 AND `status`>?2]],
            pid, C('dashboard/tasks').status.UNDERWAY);

    self:__on_loaded(find);
    return find;
end

-- 生成任务周报
function tasks:report_for_proj(pid, start_time, end_time, to)
    -- 本周验收的任务
    to.archived = self:query([[
        SELECT `id`, `creator`, `assigned`, `name`, `weight`, `tags`, `start_time`, `end_time`, `status`
        FROM `tasks`
        WHERE `pid`=?1 AND (`archive_time`>=?2 AND `archive_time`<=?3)]],
        pid, start_time, end_time);

    -- 本周该验收但未验收的任务
    to.not_archived = self:query([[
        SELECT `id`, `creator`, `assigned`, `name`, `weight`, `tags`, `start_time`, `end_time`, `status`
        FROM `tasks`
        WHERE `pid`=?1 AND UNIX_TIMESTAMP(`end_time`)<=?2 AND (`archive_time`=-1 OR `archive_time`>?2)]],
        pid, end_time, start_time);

    local users = M('user'):get_names();

    for _, info in ipairs(to.archived) do
        info.creator_name = users[info.creator] or '神秘人';
        info.assigned_name = users[info.assigned] or '神秘人';
    end

    for _, info in ipairs(to.not_archived) do
        info.creator_name = users[info.creator] or '神秘人';
        info.assigned_name = users[info.assigned] or '神秘人';
    end
end

-- 取得当前用户的所有任务
function tasks:get_mine(cond)
    local sql   = [[
        SELECT `tasks`.`id` as id, `pid`, `tasks`.`name` as name, `projects`.`name` as pname, `creator`, `assigned`, `weight`, `tags`, `start_time`, `end_time`, `status`
        FROM `tasks` LEFT JOIN `projects` ON `tasks`.`pid`=`projects`.`id`
        WHERE `status`<>]] .. C('dashboard/tasks').status.ARCHIVED .. ' AND ' .. cond;
        
    local find = self:query(sql);
    self:__on_loaded(find);
    return find;
end

-- 新建
function tasks:add(param)
    local ok, err = false, '';

    xpcall(function()
        self:exec([[
            INSERT INTO
            `tasks`(`pid`, `creator`, `assigned`, `name`, `weight`, `tags`, `start_time`, `end_time`, `content`)
            VALUES(?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9)]],
            param.pid, session.uid, param.assigned, param.name, param.weight, param.tags, param.start_time, param.end_time, param.content or '');

        self:__add_event(self:last_id(), C('dashboard/tasks').events.CREATE);
        ok = true;
    end, function(stack)
        log.error(stack);
        err = stack;
    end);

    return ok, err;
end

-- 取提指定id的任务
function tasks:get(id)
    local ret = {};

    xpcall(function()
        local tasks_ = self:query([[SELECT * FROM `tasks` WHERE `id`=?1]], id);
        if #tasks_ == 0 then return end;

        local events_ = self:query([[SELECT * FROM `task_events` WHERE `tid`=?1 ORDER BY `timepoint` DESC]], id);
        local comments_ = self:query([[SELECT * FROM `task_comments` WHERE `tid`=?1 ORDER BY `timepoint` DESC]], id);
        self:__on_loaded(tasks_, events_, comments_);

        ret.info = tasks_[1];
        ret.events = events_;
        ret.comments = comments_;
    end, function(stack)
        log.error(stack);
    end);

    return ret;
end

-- 修改指派人
function tasks:mod_assign(id, mods)
    local names = M('user'):get_names();
    local changes = { names[mods[1]] or '神秘人', names[mods[2]] or '神秘人' };

    self:exec([[
        UPDATE `tasks`
        SET `assigned`=?1
        WHERE `id`=?2]], mods[2], id);

    self:__add_event(id, C('dashboard/tasks').events.MODIFY_ASSIGNED, changes);
end

-- 修改时间
function tasks:mod_time(id, is_start, times)
    local col = 'end_time';
    local ev = C('dashboard/tasks').events.MODIFY_ENDTIME;

    if is_start then
        col = 'start_time';
        ev = C('dashboard/tasks').events.MODIFY_STARTTIME;
    end

    self:exec([[
        UPDATE `tasks`
        SET `]] .. col .. [[`=?1
        WHERE `id`=?2]], times[2], id);

    self:__add_event(id, ev, times);
    return true;
end

-- 修改状态
function tasks:mod_status(id, status)
    if status == C('dashboard/tasks').status.ARCHIVED then
        self:exec([[
            UPDATE `tasks`
            SET `status`=?1, `archive_time`=?2
            WHERE `id`=?3]], status, os.time(), id);
    else
        self:exec([[
            UPDATE `tasks`
            SET `status`=?1
            WHERE `id`=?2]], status, id);
    end
    
    self:__add_event(id, status);
    return true;
end

-- 修改优先级
function tasks:mod_weight(id, weights)
    self:exec([[
        UPDATE `tasks`
        SET `weight`=?1
        WHERE `id`=?2]], weights[2], id);
    self:__add_event(id, C('dashboard/tasks').events.MODIFY_WEIGHT, weights);
    return true;
end

-- 修改内容
function tasks:mod_content(id, content)
    self:exec([[
        UPDATE `tasks`
        SET `content`=?1
        WHERE `id`=?2]], content, id);
    self:__add_event(id, C('dashboard/tasks').events.MODIFY_CONTANT);
    return true;
end

-- 添加回复
function tasks:add_comment(id, content)
    self:exec([[
        INSERT INTO `task_comments`(`tid`, `uid`, `comment`)
        VALUES(?1, ?2, ?3)]], id, session.uid, content);
    self:__add_event(id, C('dashboard/tasks').events.ADD_COMMENT);
    return true;
end

-- 撤销回复
function tasks:del_comment(id)
    self:exec([[DELETE FROM `task_comments` WHERE `id`=?1 AND `uid`=?2]], id, session.uid);
    self:__add_event(id, C('dashboard/tasks').events.DEL_COMMENT);
    return true;
end

-- 添加事件
function tasks:__add_event(tid, event, addition)
    self:exec([[
        INSERT INTO
        `task_events`(`tid`, `uid`, `event`, `addition`)
        VALUES(?1, ?2, ?3, ?4)]],
        tid, session.uid, event, addition);
end

-- 后处理
function tasks:__on_loaded(tasks_, events_, comments_)    
    local users, avatars = M('user'):get_names();

    for _, info in ipairs(tasks_) do
        info.creator_name = users[info.creator] or '神秘人';
        info.assigned_name = users[info.assigned] or '神秘人';
        info.tags = json.decode(info.tags or '[]') or {};
    end

    for _, info in ipairs(events_ or {}) do
        info.user = users[info.uid] or '神秘人';
        info.addition = json.decode(info.addition or 'null');
    end

    for _, info in ipairs(comments_ or {}) do
        info.user = users[info.uid] or '神秘人';
        info.user_avatar = avatars[info.uid] or '/www/images/default_avatar.png';
    end
end

return tasks;