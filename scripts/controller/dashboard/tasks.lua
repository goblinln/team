-- ========================================================
-- @File    : /controller/dashboard/tasks.lua
-- @Brief   : 任务管理
-- @Author  : Leo Zhao
-- @Date    : 2018-07-19
-- ========================================================
local tasks = {};

-- 任务状态。注意：‘滞后’不是一个状态，而是处于‘进行中’的任务，在规定时间内未完成！
tasks.status = {
    OPENED      = 1,
    UNDERWAY    = 2,
    CLOSED      = 3,
    ARCHIVED    = 4,
}

-- 任务状态说明
tasks.status_desc = {
    '待办中',
    '进行中',
    '已完成',
    '已归档',
}

-- 任务标记
tasks.tags = {
    { title = 'BUG', desc = "缺陷", color = 'bg-danger' },
    { title = 'DUPLICATED', desc = "重复", color = 'bg-warning' },
    { title = 'ENHANCEMENT', desc = "优化", color = 'bg-info' },
    { title = 'FEATURE', desc = "特性", color = 'bg-success' },
    { title = 'INVALID', desc = "无效", color = 'bg-dark' },
    { title = 'QUESTION', desc = "问题", color = 'bg-primary' },
    { title = 'WONTFIX', desc = "不修复", color = 'bg-secondary' },
}

-- 任务优先级
tasks.weights = {
    { title = '一般', color = 'text-secondary' },
    { title = '次要', color = 'text-info' },
    { title = '主要', color = 'text-warning' },
    { title = '严重', color = 'text-danger' },
}

-- 任务事件定义，前面的需要与tasks.status对应
tasks.events = {
    CREATE              = 1,
    START               = 2,
    CLOSED              = 3,
    ARCHIVED            = 4,
    MODIFY_STARTTIME    = 5,
    MODIFY_ENDTIME      = 6,
    MODIFY_ASSIGNED     = 7,
    MODIFY_WEIGHT       = 8,
    MODIFY_TAGS         = 9,
    MODIFY_CONTANT      = 10,
};

-- 任务主页
function tasks:index(req, rsp)
    local tasks_    = M('tasks'):get_mine('(creator=' .. session.uid .. ' OR assigned=' .. session.uid .. ')');
    local parsed    = self:__process(session.name .. '【所有】', tasks_, true);

    parsed.dashboard_menu = 'tasks';
    rsp:html('dashboard/tasks/index.html', parsed);
end

-- 取得一个任务的信息
function tasks:info(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;
    local task_ = M('tasks'):get(req.post.id);
    self:__process_event(task_.events);
    rsp:html('dashboard/tasks/task_info.html', {
        task        = task_,
        events      = self.events,
        tags        = self.tags,
        weights     = self.weights,
        members     = M('projects'):get_members(task_.info.pid),
    });
end

-- 查看任务信息，不可编辑
function tasks:readonly_info(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;
    local task_ = M('tasks'):get(req.post.id);
    self:__process_event(task_.events);
    rsp:html('dashboard/tasks/task_info_readonly.html', {
        enable_archive  = req.post.enable_archive,
        task            = task_,
        events          = self.events,
        tags            = self.tags,
        weights         = self.weights
    });
end

-- 取得选中项目的所有成员
function tasks:members_of_proj(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;
    rsp:json(M('projects'):get_members(req.post.pid));
end

-- 发布任务页面
function tasks:try_create(req, rsp)
    local projs = M('projects'):all_of_mine();
    if #projs == 0 then
        rsp:html('dashboard/tasks/ask_create_proj.html', { roles = C('dashboard/projects').roles });
    else
        rsp:html('dashboard/tasks/create_task.html', {
            projs = projs,
            tags = self.tags,
            weights = self.weights,
        });
    end
end

-- 确认发布任务
function tasks:do_create(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    -- 将数值参数处理一下
    local param = req.post;
    param.pid = tonumber(param.pid);
    param.assigned = tonumber(param.assigned);
    param.weight = tonumber(param.weight);
    param.tags = param.tags or {};
    for n, t in ipairs(param.tags) do param.tags[n] = tonumber(t) end;

    local ok, err = M('tasks'):add(param);
    if not ok then return rsp:json({ok = false, err_msg = err}) end;
    rsp:json{ok = true};
end

----------------------- 修改任务 --------------------------

-- 修改指派人
function tasks:mod_assign(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end

    local mod = { tonumber(req.post.org_uid), tonumber(req.post.new_uid) };
    local ok = true;
    
    if mod[1] ~= mod[2] then
        ok = M('tasks'):mod_assign(req.post.tid, mod);
    end
    
    rsp:json{ok = true};
end

-- 修改时间
function tasks:mod_time(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local tid = req.post.tid;
    local org_start = req.post.org_start;
    local org_end = req.post.org_end;
    local start_time = req.post.start_time;
    local end_time = req.post.end_time;

    if org_start ~= start_time then
        M('tasks'):mod_time(tid, true, {org_start, start_time});
    end

    if org_end ~= end_time then
        M('tasks'):mod_time(tid, false, {org_end, end_time});
    end
    
    rsp:json{ok = true};
end

-- 修改状态
function tasks:mod_status(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local mod = { tonumber(req.post.org_status), tonumber(req.post.new_status) };
    local ok = true;
    
    if mod[1] ~= mod[2] then
        ok = M('tasks'):mod_status(req.post.tid, mod[2]);
    end
    
    rsp:json{ok = true};
end

-- 修改优先级
function tasks:mod_weight(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local mod = { tonumber(req.post.org_weight), tonumber(req.post.new_weight) };
    local ok = true;
    
    if mod[1] ~= mod[2] then
        ok = M('tasks'):mod_weight(req.post.tid, mod);
    end
    
    rsp:json{ok = true};
end

-- 修改内容
function tasks:mod_content(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;
    local ok = M('tasks'):mod_content(req.post.id, req.post.content);
    rsp:json{ok = true};
end

----------------------- 任务筛选 --------------------------

-- 我的所有任务
function tasks:all_of_mine(req, rsp)
    local tasks_ = M('tasks'):get_mine('(creator=' .. session.uid .. ' OR assigned=' .. session.uid .. ')');
    self:__layout_tasks(rsp, session.name .. '【所有】', tasks_, true);
end

-- 选出我发布的任务
function tasks:create_by_me(req, rsp)
    self:__layout_tasks(rsp, session.name .. '【发布】', M('tasks'):get_mine('creator=' .. session.uid), true);
end

-- 选出指派给我的任务
function tasks:assigned_to_me(req, rsp)
    self:__layout_tasks(rsp, '【指派给】' .. session.name, M('tasks'):get_mine('assigned=' .. session.uid), true);
end

-- 选出指定任务等级的任务
function tasks:filter_weight(req, rsp)
    local tasks_ = M('tasks'):get_mine('weight=' .. req.post.p1 .. ' AND (creator=' .. session.uid .. ' OR assigned=' .. session.uid .. ')');
    self:__layout_tasks(rsp, session.name .. string.format('【%s】', self.weights[tonumber(req.post.p1)].title), tasks_, true);
end

-- 选出指定项目的任务
function tasks:filter_proj(req, rsp)
    local tasks_ = M('tasks'):get_mine('pid=' .. req.post.p1 .. ' AND (creator=' .. session.uid .. ' OR assigned=' .. session.uid .. ')');
    self:__layout_tasks(rsp, session.name .. string.format('@【%s】', req.post.p2), tasks_, true);
end

-----------------------------------------------------------

-- 布局任务列表
function tasks:__layout_tasks(rsp, title, tasks_, process_mine) 
    rsp:html('dashboard/tasks/view_tasks.html', self:__process(title, tasks_, process_mine));
end

-- 分析事件，解析等
function tasks:__process_event(evs)
    for _, info in ipairs(evs) do
        info.timepoint = tostring(info.timepoint);
        if info.event == self.events.CREATE then
            info.event_desc = (_ ~= #evs and '重新' or '') .. '发布了任务';
        elseif info.event == self.events.START then
            info.event_desc = '开始了任务';
        elseif info.event == self.events.CLOSED then
            info.event_desc = '关闭了任务';
        elseif info.event == self.events.ARCHIVED then
            info.event_desc = '归档了任务';
        elseif info.event == self.events.MODIFY_STARTTIME then
            info.event_desc = '修改任务开始时间 : ' .. info.addition[1] .. ' > ' .. info.addition[2];
        elseif info.event == self.events.MODIFY_ENDTIME then      
            info.event_desc = '修改任务结束时间 : ' .. info.addition[1] .. ' > ' .. info.addition[2];
        elseif info.event == self.events.MODIFY_ASSIGNED then
            info.event_desc = '修改任务指派：' .. info.addition[1] .. ' > ' .. info.addition[2];
        elseif info.event == self.events.MODIFY_WEIGHT then            
            info.event_desc = '修改任务优先级 : ' .. self.weights[info.addition[1]].title .. ' > ' .. self.weights[info.addition[2]].title;
        elseif info.event == self.events.MODIFY_TAGS then
            info.event_desc = '修改任务标签';
        elseif info.event == self.events.MODIFY_CONTANT then
            info.event_desc = '修改任务内容';
        else
            info.event_desc = '对任务其他内容进行了修改';
        end
    end
end

-- 分析任务，分组等
function tasks:__process(title, tasks_, process_mine)
    local opened, closed, underway, delayed, gantt_groups, gantt_map = {}, {}, {}, {}, {}, {}, {};
    local summary = { title = title, opened = 0, closed = 0, underway = 0, delayed = 0 };
    local now = os.time();

    local mine = {
        opened          = 0,
        underway        = 0,
        delayed         = 0,
        closed          = 0,
        create_by_me    = 0,
        assigned_to_me  = 0,
        weights         = {},
        projects        = {},
    };

    for _, info in ipairs(tasks_) do
        local percent = 0;

        if info.status == self.status.OPENED then
            table.insert(opened, info);
            summary.opened = summary.opened + 1;
            mine.opened = mine.opened + 1;

            local deadline = os.time(info.end_time);
            if deadline <= now then
                table.insert(delayed, info);
                summary.delayed = summary.delayed + 1;
                mine.delayed = mine.delayed + 1;
            end
        elseif info.status == self.status.UNDERWAY then
            table.insert(underway, info);
            summary.underway = summary.underway + 1;
            mine.underway = mine.underway + 1;

            local deadline = os.time(info.end_time);
            if deadline <= now then
                table.insert(delayed, info);
                summary.delayed = summary.delayed + 1;                
                mine.delayed = mine.delayed + 1;
            end
            percent = 0.5;
        elseif info.status == self.status.CLOSED then
            table.insert(closed, info);
            summary.closed = summary.closed + 1;
            mine.closed = mine.closed + 1;
            percent = 1;
        end

        if process_mine then
            if info.creator == session.uid then
                mine.create_by_me = mine.create_by_me + 1;
            end
            
            if info.assigned == session.uid then
                mine.assigned_to_me = mine.assigned_to_me + 1;
            end

            mine.weights[info.weight] = mine.weights[info.weight] or 0;
            mine.weights[info.weight] = mine.weights[info.weight] + 1;

            mine.projects[info.pid] = mine.projects[info.pid] or { name = info.pname, count = 0 };
            mine.projects[info.pid].count = mine.projects[info.pid].count + 1;
        end

        if not gantt_map[info.assigned] then
            table.insert(gantt_groups, {
                id = info.assigned,
                name = info.assigned_name,
                children = {},
            });
            gantt_map[info.assigned] = #gantt_groups;
        end
        
        table.insert(gantt_groups[gantt_map[info.assigned]].children, {
            id = info.id,
            name = info.name,
            start_time = tostring(info.start_time),
            end_time = tostring(info.end_time),
            percent = percent,
        });
    end

    return {
        summary     = summary,
        mine        = process_mine and mine or nil;
        tags        = self.tags,
        weights     = self.weights,
        gantt_data  = #gantt_groups == 0 and '[]' or json.encode(gantt_groups),
        tasks       = {
            opened      = opened,
            closed      = closed,
            underway    = underway,
            delayed     = delayed,
        }
    }
end

return tasks;