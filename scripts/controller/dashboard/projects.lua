-- ========================================================
-- @File    : /controller/dashboard/projects.lua
-- @Brief   : 项目模块
-- @Author  : Leo Zhao
-- @Date    : 2018-07-19
-- ========================================================
local projects = {};

-- 人员分工
projects.roles = {
    '未指定',
    '策划',
    '程序',
    '测试',
    '运营',
    '美术',
}

-- 项目主页
function projects:index(req, rsp)
    local valids = M('projects'):all_of_mine();
    rsp:html('dashboard/projects/index.html', {
        dashboard_menu  = 'projects',
        projs           = valids, 
        roles           = self.roles});
end

-- 取得指定项目任务信息
function projects:get_tasks(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local tasks = M('tasks'):get_by_proj(req.post.pid);
    local is_admin = M('projects'):is_admin(req.post.pid);
    C('dashboard/tasks'):__layout_tasks(rsp, req.post.pname, tasks, not is_admin);
end

-- 取得指定项目成员信息
function projects:get_members(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;
    if not M('projects'):is_admin(req.post.pid) then return rsp:error(403) end;

    local members = M('projects'):get_members(req.post.pid);
    rsp:html('dashboard/projects/members.html', {
        proj_id = req.post.pid,
        proj_name = req.post.pname,
        roles = self.roles,
        members = members,
    });
end

-- 取得项目可邀请的成员
function projects:get_invite_users(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;    
    if not M('projects'):is_admin(req.post.pid) then return rsp:error(403) end;

    local users, err = M('projects'):get_invite_users(req.post.pid);
    if not users then return rsp:json{ ok = false, err_msg = err } end;
    rsp:json{ ok = true, users = users };
end

-- 添加成员
function projects:add_member(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;    
    if not M('projects'):is_admin(req.post.pid) then return rsp:error(403) end;

    local ok, err = M('projects'):add_member(req.post.pid, req.post.uid, req.post.role, req.post.is_admin);
    if not ok then return rsp:json{ ok = false, err_msg = err } end;
    rsp:json{ ok = true };
end

-- 修改成员属性
function projects:edit_member(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;    
    if not M('projects'):is_admin(req.post.pid) then return rsp:error(403) end;

    local ok, err = M('projects'):edit_member(req.post.pid, req.post.uid, req.post.role, req.post.is_admin);
    if not ok then return rsp:json{ ok = false, err_msg = err } end;
    rsp:json{ ok = true, remove_admin = (not req.post.is_admin) };
end

-- 移除团队成员
function projects:del_member(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;    
    if not M('projects'):is_admin(req.post.pid) then return rsp:error(403) end;

    local ok, err = M('projects'):del_member(tonumber(req.post.pid), tonumber(req.post.uid));
    if not ok then return rsp:json{ ok = false, err_msg = err } end;
    rsp:json{ ok = true };
end

-- 取得周报
function projects:get_reports(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local report = M('projects'):get_reports(req.post.pid, req.post.week_offset or 0);
    rsp:html('dashboard/projects/reports.html', {
        proj_id = req.post.pid,
        proj_name = req.post.pname,
        weights = C('dashboard/tasks').weights,
        report = report,
    });
end

-- 取得可归档的任务
function projects:get_can_archive(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local tasks = M('tasks'):get_archivable_by_proj(req.post.pid);
    rsp:html('dashboard/projects/archive.html', {
        tasks = tasks,
        tags = C('dashboard/tasks').tags,
        weights = C('dashboard/tasks').weights,
    });
end

return projects;