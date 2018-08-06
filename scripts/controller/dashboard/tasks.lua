-- ========================================================
-- @File    : /controller/dashboard/tasks.lua
-- @Brief   : 任务管理
-- @Author  : Leo Zhao
-- @Date    : 2018-07-19
-- ========================================================
local tasks = {};

-- 任务状态。注意：‘滞后’不是一个状态，在规定时间内未开发或未完成！
tasks.status = {
    CREATED     = 1,
    UNDERWAY    = 2,
    TESTING     = 3,
    FINISHED    = 4,
    ARCHIVED    = 5,
}

-- 任务标记
tasks.tags = {
    { title = 'BUG', desc = "缺陷", color = 'bg-danger' },
    { title = 'QUICKFIX', desc = "快速修正", color = 'bg-secondary' },
    { title = 'ENHANCEMENT', desc = "优化", color = 'bg-info' },
    { title = 'FEATURE', desc = "特性", color = 'bg-success' },
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
    CREATED             = 1,
    UNDERWAY            = 2,
    TESTING             = 3,
    FINISHED            = 4,
    ARCHIVED            = 5,
    MODIFY_STARTTIME    = 6,
    MODIFY_ENDTIME      = 7,
    MODIFY_ASSIGNED     = 8,
    MODIFY_COOPERATOR   = 9,
    MODIFY_WEIGHT       = 10,
    MODIFY_TAGS         = 11,
    MODIFY_CONTANT      = 12,
    ADD_COMMENT         = 13,
    DEL_COMMENT         = 14,
    ADD_ATTACHMENT      = 15,
    DEL_ATTACHMENT      = 16,
};

-- 任务主页
function tasks:index(req, rsp)
    local m         = M('tasks');
    local tasks_    = m:get_mine(m.filters.ABOUT);
    local parsed    = self:__process(session.name .. '【所有】', tasks_);

    parsed.dashboard_menu = 'tasks';
    parsed.mine = self:__process_mine(tasks_);
    rsp:html('dashboard/tasks/index.html', parsed);
end

-- 取得一个任务的信息
function tasks:info(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local task_ = M('tasks'):get(req.post.id);
    self:__process_event(task_.events);
    rsp:html('dashboard/tasks/task_info.html', {
        task        = task_,
        tags        = self.tags,
        weights     = self.weights,
        is_admin    = M('projects'):is_admin(task_.info.pid),
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
    param.cooperator = tonumber(param.cooperator);
    param.weight = tonumber(param.weight);
    param.tags = param.tags or {};
    for n, t in ipairs(param.tags) do param.tags[n] = tonumber(t) end;

    local ok, err = M('tasks'):add(param, req.file);
    if not ok then return rsp:json({ok = false, err_msg = err}) end;
    rsp:json{ ok = true };
end

-- 切换预览模式
function tasks:switch_viewmode(req, rsp)
    session.gantt_mode = (not session.gantt_mode);
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

-- 修改协作者（测试或验收人员）
function tasks:mod_cooperator(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end

    local mod = { tonumber(req.post.org_uid), tonumber(req.post.new_uid) };
    local ok = true;
    
    if mod[1] ~= mod[2] then
        ok = M('tasks'):mod_cooperator(req.post.tid, mod);
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

-- 评论
function tasks:add_comment(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;
    local ok = M('tasks'):add_comment(req.post.id, req.post.comment);
    rsp:json{ok = true};
end

-- 撤销评论
function tasks:del_comment(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;
    local ok = M('tasks'):del_comment(req.post.id);
    rsp:json{ok = true};
end

-- 添加附件
function tasks:add_attachment(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local data = C('dashboard/files'):do_upload(req.file);
    for _, info in ipairs(data) do
        if info.ok then
            if not M('tasks'):add_attachment(req.post.tid, info.name, info.url) then
                return rsp:json{ ok = false, err_msg = '上传文件失败' };
            end
        else
            return rsp:json{ ok = false, err_msg = '上传文件失败' };
        end
    end
    
    return rsp:json{ok = true};
end

-- 删除附件
function tasks:del_attachment(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local ok, err = M('tasks'):del_attachment(req.post.tid, req.post.aid);
    rsp:json{ ok = ok, err_msg = err };
end

----------------------- 任务筛选 --------------------------

-- 我的所有任务
function tasks:all_of_mine(req, rsp)
    local m         = M('tasks');
    local tasks_    = m:get_mine(m.filters.ABOUT);
    self:__layout_tasks(rsp, session.name .. '【所有】', tasks_);
end

-- 选出我发布的任务
function tasks:create_by_me(req, rsp)
    local m         = M('tasks');
    local tasks_    = m:get_mine(m.filters.CREATE_BY);
    self:__layout_tasks(rsp, session.name .. '【发布】', tasks_);
end

-- 选出指派给我的任务
function tasks:assigned_to_me(req, rsp)
    local m         = M('tasks');
    local tasks_    = m:get_mine(m.filters.ASSIGNED_TO);
    self:__layout_tasks(rsp, '【指派给】' .. session.name, tasks_);
end

-- 选出指定任务等级的任务
function tasks:filter_weight(req, rsp)
    local m         = M('tasks');
    local tasks_    = m:get_mine(m.filters.BY_WEIGHT, req.post.p1);
    self:__layout_tasks(rsp, session.name .. string.format('【%s】', self.weights[tonumber(req.post.p1)].title), tasks_);
end

-- 选出指定项目的任务
function tasks:filter_proj(req, rsp)
    local m         = M('tasks');
    local tasks_    = m:get_mine(m.filters.BY_PROJ, req.post.p1);
    self:__layout_tasks(rsp, session.name .. string.format('@【%s】', req.post.p2), tasks_);
end

-----------------------------------------------------------

-- 布局任务列表
function tasks:__layout_tasks(rsp, title, tasks_, readonly) 
    rsp:html('dashboard/tasks/view_tasks.html', self:__process(title, tasks_, readonly));
end

-- 分析任务，分组等
function tasks:__process(title, tasks_, readonly)
    local created, underway, testing, finished, gantt_data, gantt_map = {}, {}, {}, {}, {}, {};
    local summary = { title = title, created = 0, underway = 0, testing = 0, finished = 0 };
    local now = os.time();

    for _, info in ipairs(tasks_) do
        local gantt_color = 'grey';

        if info.status == self.status.CREATED then
            table.insert(created, info);
            summary.created = summary.created + 1;

            if os.time(info.start_time) <= now then
                info.delayed = true;
                gantt_color = 'red';
            end
        elseif info.status == self.status.UNDERWAY then
            table.insert(underway, info);
            summary.underway = summary.underway + 1;
            
            if os.time(info.end_time) <= now then
                info.delayed = true;
                gantt_color = 'red';
            else
                gantt_color = '#17a2b8';
            end
        elseif info.status == self.status.TESTING then
            table.insert(testing, info);
            summary.testing = summary.testing + 1;

            if os.time(info.end_time) <= now then
                info.delayed = true;
                gantt_color = 'red';
            else
                gantt_color = '#007bff';
            end
        elseif info.status == self.status.FINISHED then
            table.insert(finished, info);
            summary.finished = summary.finished + 1;
            gantt_color = 'lightgreen';
        end

        if not gantt_map[info.assigned] then
            table.insert(gantt_data, {
                id = info.assigned,
                name = info.assigned_name,
                series = {},
            });
            gantt_map[info.assigned] = #gantt_data;
        end
        
        table.insert(gantt_data[gantt_map[info.assigned]].series, {
            id = info.id,
            name = info.name,
            cooperator = info.cooperator_name,
            start = tostring(info.start_time),
            ['end'] = tostring(info.end_time),
            options = { color = gantt_color },
        });
    end

    return {
        summary     = summary,
        tags        = self.tags,
        weights     = self.weights,
        readonly    = readonly,
        gantt_data  = #gantt_data == 0 and '[]' or json.encode(gantt_data),
        tasks       = {
            created     = created,
            underway    = underway,
            testing     = testing,
            finished    = finished,
        }
    }
end

-- 分析当前用户的任务
function tasks:__process_mine(tasks_)
    local mine = {
        create_by_me    = 0,
        assigned_to_me  = 0,
        weights         = {},
        projects        = {},
    };

    for _, info in ipairs(tasks_) do
        if info.creator == session.uid then
            mine.create_by_me = mine.create_by_me + 1;
        end
        
        if info.assigned == session.uid or info.cooperator == session.uid then
            mine.assigned_to_me = mine.assigned_to_me + 1;
        end

        mine.weights[info.weight] = mine.weights[info.weight] or 0;
        mine.weights[info.weight] = mine.weights[info.weight] + 1;

        mine.projects[info.pid] = mine.projects[info.pid] or { name = info.pname, count = 0 };
        mine.projects[info.pid].count = mine.projects[info.pid].count + 1;
    end

    return mine;
end

-- 分析事件，解析等
function tasks:__process_event(evs)
    for _, info in ipairs(evs) do
        info.timepoint = tostring(info.timepoint);
        if info.event == self.events.CREATED then
            info.event_desc = (_ ~= #evs and '重新' or '') .. '发布了任务';
        elseif info.event == self.events.UNDERWAY then
            info.event_desc = '开始了任务';
        elseif info.event == self.events.TESTING then
            info.event_desc = '开启了测试流程';
        elseif info.event == self.events.FINISHED then
            info.event_desc = '完成了任务';
        elseif info.event == self.events.ARCHIVED then
            info.event_desc = '验收了任务';
        elseif info.event == self.events.MODIFY_STARTTIME then
            info.event_desc = '修改任务开始时间 : ' .. info.addition[1] .. ' > ' .. info.addition[2];
        elseif info.event == self.events.MODIFY_ENDTIME then      
            info.event_desc = '修改任务结束时间 : ' .. info.addition[1] .. ' > ' .. info.addition[2];
        elseif info.event == self.events.MODIFY_ASSIGNED then
            info.event_desc = '修改任务指派：' .. info.addition[1] .. ' > ' .. info.addition[2];
        elseif info.event == self.events.MODIFY_COOPERATOR then
            info.event_desc = '修改协作人员：' .. info.addition[1] .. ' > ' .. info.addition[2];
        elseif info.event == self.events.MODIFY_WEIGHT then            
            info.event_desc = '修改任务优先级 : ' .. self.weights[info.addition[1]].title .. ' > ' .. self.weights[info.addition[2]].title;
        elseif info.event == self.events.MODIFY_TAGS then
            info.event_desc = '修改任务标签';
        elseif info.event == self.events.MODIFY_CONTANT then
            info.event_desc = '修改任务内容';
        elseif info.event == self.events.ADD_COMMENT then
            info.event_desc = '评论了该任务';
        elseif info.event == self.events.DEL_COMMENT then
            info.event_desc = '撤销了一条评论';
        elseif info.event == self.events.ADD_ATTACHMENT then
            info.event_desc = '上传了任务附件: ' .. info.addition[1];
        elseif info.event == self.events.DEL_ATTACHMENT then
            info.event_desc = '删除了任务附件：' .. info.addition[1];
        else
            info.event_desc = '对任务其他内容进行了修改';
        end
    end
end

return tasks;