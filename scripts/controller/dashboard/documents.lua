-- ========================================================
-- @File    : /controller/dashboard/documents.lua
-- @Brief   : 文档管理
-- @Author  : Leo Zhao
-- @Date    : 2018-07-17
-- ========================================================
local documents = {};

-- 文档主页
function documents:index(req, rsp)
    rsp:html('dashboard/documents/index.html', { dashboard_menu = 'documents' });
end

-- 获取列表
function documents:table_of_contents(req, rsp)
    local tree = M('documents'):table_of_contents();
    if #tree == 0 then
        table.insert(tree, { icon = 'fa fa-trash-o', text = '空目录，右键新建', id = -1, data = -1 });
    end

    rsp:json(tree);
end

-- 获取文档信息
function documents:info(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local info = M('documents'):get(req.post.id);
    if not info then return rsp:json{ ok = false, err_msg = '没找到ID为:' .. req.post.id .. '的文档！' } end;
    return rsp:json{
        ok = true,
        author = info.author,
        last_modify = tostring(info.modify_time),
        modify_user = info.modify_user,
        content = info.content };
end

-- 新建文档
function documents:add(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local ok, err = M('documents'):add(req.post.name, req.post.parent_id);
    if not ok then return rsp:json{ ok = false, err_msg = err } end;
    return rsp:json{ok = true};
end

-- 重命名文档
function documents:rename(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local ok, err = M('documents'):rename(req.post.name, req.post.rename_id);
    if not ok then return rsp:json{ ok = false, err_msg = err } end;
    return rsp:json{ok = true};
end

-- 编辑文档
function documents:edit(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local ok, err = M('documents'):edit(req.post.id, req.post.content);
    if not ok then return rsp:json{ ok = false, err_msg = err } end;
    return rsp:json{ok = true};
end

-- 移动文档
function documents:move(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local ok, err = M('documents'):move(req.post.id, req.post.to);
    if not ok then return rsp:json{ ok = false, err_msg = err } end;
    return rsp:json{ok = true};
end

-- 删除文档
function documents:delete(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local ok, err = M('documents'):delete(req.post.id, req.post.name);
    if not ok then return rsp:json{ ok = false, err_msg = err } end;
    return rsp:json{ok = true};
end

return documents;