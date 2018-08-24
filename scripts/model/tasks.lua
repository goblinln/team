-- ========================================================
-- @File    : /model/tasks.lua
-- @Brief   : 任务管理
-- @Author  : Leo Zhao
-- @Date    : 2018-07-20
-- ========================================================
local tasks = inherit(M('base'));

-- 任务所属过滤
tasks.filters = {
    ABOUT       = 1,
    CREATE_BY   = 2,
    ASSIGNED_TO = 3,
    BY_WEIGHT   = 4,
    BY_PROJ     = 5,
}

-- 取得一个项目的任务
function tasks:get_by_proj(pid)
    local find  = self:query([[
            SELECT `id`, `pid`, `creator`, `assigned`, `cooperator`, `name`, `weight`, `tags`, `start_time`, `end_time`, `status`
            FROM `tasks`
            WHERE `pid`=?1 AND `status`<>?2]],
            pid, C('dashboard/tasks').status.ARCHIVED);
            
    self:__on_loaded(find);
    return find;
end

-- 取得一个项目可归档的任务列表
function tasks:get_archivable_by_proj(pid)
    local timepoint = os.time();
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
    local find = self:query([[
            SELECT `id`, `creator`, `assigned`, `cooperator`, `name`, `weight`, `tags`, `start_time`, `end_time`, `status`
            FROM `tasks`
            WHERE `pid`=?1 AND `status`>?2 AND (`archive_time`=-1 OR `archive_time`>=?3)]],
            pid, C('dashboard/tasks').status.UNDERWAY, week_start);

    self:__on_loaded(find);
    return find;
end

-- 生成任务周报
function tasks:report_for_proj(pid, start_time, end_time, to)
    -- 本周验收的任务
    to.archived = self:query([[
        SELECT `id`, `creator`, `assigned`, `cooperator`, `name`, `weight`, `tags`, `start_time`, `end_time`, `status`
        FROM `tasks`
        WHERE `pid`=?1 AND (`archive_time`>=?2 AND `archive_time`<=?3)]],
        pid, start_time, end_time);

    -- 本周该验收但未验收的任务
    to.not_archived = self:query([[
        SELECT `id`, `creator`, `assigned`, `cooperator`, `name`, `weight`, `tags`, `start_time`, `end_time`, `status`
        FROM `tasks`
        WHERE `pid`=?1 AND UNIX_TIMESTAMP(`end_time`)<=?2 AND (`archive_time`=-1 OR `archive_time`>?2)]],
        pid, end_time, start_time);

    local users = M('user'):get_names();

    for _, info in ipairs(to.archived) do
        info.creator_name = users[info.creator] or '神秘人';
        info.assigned_name = users[info.assigned] or '神秘人';
        info.cooperator_name = users[info.cooperator] or '神秘人';
    end

    for _, info in ipairs(to.not_archived) do
        info.creator_name = users[info.creator] or '神秘人';
        info.assigned_name = users[info.assigned] or '神秘人';
        info.cooperator_name = users[info.cooperator] or '神秘人';
    end
end

-- 取得当前用户的所有任务
function tasks:get_mine(filter, v)
    local cond = 'false';

    if filter == self.filters.ABOUT then
        cond = string.gsub('(`creator`=__ME OR `assigned`=__ME OR `cooperator`=__ME)', '__ME', session.uid);
    elseif filter == self.filters.CREATE_BY then
        cond = '`creator`=' .. session.uid;
    elseif filter == self.filters.ASSIGNED_TO then
        cond = string.format('(`assigned`=%d OR `cooperator`=%d)', session.uid, session.uid);
    elseif filter == self.filters.BY_WEIGHT then
        cond = string.format('(`creator`=%d OR `assigned`=%d OR `cooperator`=%d) AND (`weight`=%d)', session.uid, session.uid, session.uid, tonumber(v));
    elseif filter == self.filters.BY_PROJ then
        cond = string.format('(`creator`=%d OR `assigned`=%d OR `cooperator`=%d) AND (`pid`=%d)', session.uid, session.uid, session.uid, tonumber(v));
    end

    local sql   = [[
        SELECT `tasks`.`id` as id, `pid`, `tasks`.`name` as name, `projects`.`name` as pname, `creator`, `assigned`, `cooperator`, `weight`, `tags`, `start_time`, `end_time`, `status`
        FROM `tasks` LEFT JOIN `projects` ON `tasks`.`pid`=`projects`.`id`
        WHERE `status`<>]] .. C('dashboard/tasks').status.ARCHIVED .. ' AND ' .. cond;
        
    local find = self:query(sql);
    self:__on_loaded(find);
    return find;
end

-- 新建
function tasks:add(param, files)
    local ok, err = false, '';

    if not param.content or string.len(param.content) == 0 then
        return false, '任务详情必须写明！';
    end

    local uploaded = {};

    for name, path in pairs(files) do
        local dir = 'www/upload/' .. session.uid;
        if not os.exists(dir) then os.mkdir(dir) end;

        local to = dir .. '/' .. os.time() .. '_' .. name;
        if os.cp(path, to) then
            table.insert(uploaded, { name = name, url = '/' .. to });
        end
    end

    xpcall(function()
        self:exec([[
            INSERT INTO
            `tasks`(`pid`, `creator`, `assigned`, `cooperator`, `name`, `weight`, `tags`, `start_time`, `end_time`, `content`)
            VALUES(?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10)]],
            param.pid, session.uid, param.assigned, param.cooperator, param.name, param.weight, param.tags, param.start_time, param.end_time, param.content);

        local tid = self:last_id();
        self:__add_event(tid, C('dashboard/tasks').events.CREATED, {1});

        for _, attachment in ipairs(uploaded) do
            self:exec([[
                INSERT INTO `task_attachments`(`tid`, `name`, `url`)
                VALUES(?1, ?2, ?3)]], tid, attachment.name, attachment.url);
        end
        
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
        local attachments_ = self:query([[SELECT * FROM `task_attachments` WHERE `tid`=?1]], id);
        self:__on_loaded(tasks_, events_, comments_);

        ret.info = tasks_[1];
        ret.events = events_;
        ret.comments = comments_;
        ret.attachments = attachments_;
    end, function(stack)
        log.error(stack);
    end);

    return ret;
end

-- 删除
function tasks:delete(id)
    local info = self:query([[SELECT * FROM `tasks` WHERE `id`=?1]], id)[1];
    if session.uid ~= info.creator and not M('projects'):is_admin(info.pid) then
        return false, '您没有权限删除任务！';
    end

    self:exec([[DELETE FROM `tasks` WHERE `id`=?1]], id);
    return true;
end

-- 修改名字
function tasks:mod_name(id, name)
    self:exec([[
        UPDATE `tasks`
        SET `name`=?1
        WHERE `id`=?2]], name, id);

    self:__add_event(id, C('dashboard/tasks').events.RENAME);
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

-- 修改协作者（测试或验收人员）
function tasks:mod_cooperator(id, mods)
    local names = M('user'):get_names();
    local changes = { names[mods[1]] or '神秘人', names[mods[2]] or '神秘人' };

    self:exec([[
        UPDATE `tasks`
        SET `cooperator`=?1
        WHERE `id`=?2]], mods[2], id);

    self:__add_event(id, C('dashboard/tasks').events.MODIFY_COOPERATOR, changes);
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

-- 添加附件
function tasks:add_attachment(id, name, url)
    self:exec([[
        INSERT INTO `task_attachments`(`tid`, `name`, `url`)
        VALUES(?1, ?2, ?3)]], id, name, url);
    self:__add_event(id, C('dashboard/tasks').events.ADD_ATTACHMENT, {name});
    return true;
end

-- 删除附件
function tasks:del_attachment(id, aid)
    local info = self:query([[
        SELECT *
        FROM `task_attachments`
        WHERE `id`=?1 AND `tid`=?2]], aid, id)[1];
    if not info then return false, '附件不存在或已被删除' end;

    os.rm('.' .. info.url);
    self:exec([[DELETE FROM `task_attachments` WHERE `id`=?1]], aid);
    return true;
end

-------------------------- 内部调用 --------------------------

-- 添加事件
function tasks:__add_event(tid, event, addition)
    self:exec([[
        INSERT INTO
        `task_events`(`tid`, `uid`, `event`, `addition`)
        VALUES(?1, ?2, ?3, ?4)]],
        tid, session.uid, event, addition);

    local info = self:query([[SELECT * FROM `tasks` WHERE `id`=?1]], tid)[1];
    local evs = C('dashboard/tasks').events;
    local who = '<label class="text-muted mr-1 mb-0">' .. session.name .. '</label>';
    local action = '';
    local task_link = '<a href="#" class="task-link mx-1" onclick="return open_task_via_notice(' .. tid .. ');">' .. info.name .. '</a>';

    if event == evs.CREATED then
        action = who .. (addition and '' or '重新') .. '发布了任务' .. task_link;
    elseif event == evs.UNDERWAY then
        action = who .. '开始了任务' .. task_link;
    elseif event == evs.TESTING then
        action = who .. '对任务' .. task_link .. '开启了测试流程';
    elseif event == evs.FINISHED then
        action = who .. '完成了任务' .. task_link;
    elseif event == evs.ARCHIVED then
        action = who .. '验收了任务' .. task_link;
    elseif event == evs.MODIFY_STARTTIME then
        action = who .. '修改任务' .. task_link .. '的开始时间 : ' .. addition[1] .. ' > ' .. addition[2];
    elseif event == evs.MODIFY_ENDTIME then      
        action = who .. '修改任务' .. task_link .. '的结束时间 : ' .. addition[1] .. ' > ' .. addition[2];
    elseif event == evs.MODIFY_ASSIGNED then
        action = who .. '修改任务' .. task_link .. '的指派：' .. addition[1] .. ' > ' .. addition[2];
    elseif event == evs.MODIFY_COOPERATOR then
        action = who .. '修改任务' .. task_link .. '的协作人员：' .. addition[1] .. ' > ' .. addition[2];
    elseif event == evs.MODIFY_WEIGHT then            
        action = who .. '修改任务' .. task_link .. '的优先级 : ' .. C('dashboard/tasks').weights[addition[1]].title .. ' > ' .. C('dashboard/tasks').weights[addition[2]].title;
    elseif event == evs.MODIFY_TAGS then
        action = who .. '修改任务' .. task_link .. '的标签';
    elseif event == evs.MODIFY_CONTANT then
        action = who .. '修改任务' .. task_link .. '的内容';
    elseif event == evs.ADD_COMMENT then
        action = who .. '评论了' .. task_link;
    elseif event == evs.DEL_COMMENT then
        action = who .. '撤销了' .. task_link .. '中一条评论';
    elseif event == evs.ADD_ATTACHMENT then
        action = who .. '上传了附件: ' .. addition[1] .. '到' .. task_link;
    elseif event == evs.DEL_ATTACHMENT then
        action = who .. '删除了任务' .. task_link .. '的附件：' .. addition[1];
    elseif event == evs.RENAME then
        action = who .. '修改了任务' .. task_link .. '的名称';
    else
        action = who .. '对任务' .. task_link .. '其他内容进行了修改';
    end

    local notified = {};

    if info.creator ~= session.uid and not notified[info.creator] then
        M('notification'):add(info.creator, action);
        notified[info.creator] = true;
    end

    if info.assigned ~= session.uid and not notified[info.assigned] then
        M('notification'):add(info.assigned, action);
        notified[info.assigned] = true;
    end

    if info.cooperator ~= session.uid and not notified[info.cooperator] then
        M('notification'):add(info.cooperator, action);
        notified[info.cooperator] = true;
    end
end

-- 后处理
function tasks:__on_loaded(tasks_, events_, comments_)    
    local users, avatars = M('user'):get_names();

    for _, info in ipairs(tasks_) do
        info.creator_name = users[info.creator] or '神秘人';
        info.assigned_name = users[info.assigned] or '神秘人';
        info.cooperator_name = users[info.cooperator] or '神秘人';
        info.tags = json.decode(info.tags or '[]') or {};
    end

    for _, info in ipairs(events_ or {}) do
        info.user = users[info.uid] or '神秘人';
        info.addition = json.decode(info.addition or '[]');
    end

    for _, info in ipairs(comments_ or {}) do
        info.user = users[info.uid] or '神秘人';
        info.user_avatar = avatars[info.uid] or '/www/images/default_avatar.png';
    end
end

return tasks;