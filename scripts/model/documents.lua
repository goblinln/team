-- ========================================================
-- @File    : /model/documents.lua
-- @Brief   : 文档系统数据
-- @Author  : Leo Zhao
-- @Date    : 2018-07-16
-- ========================================================
local documents = inherit(M('base'));

-- 获取图书列表
function documents:table_of_contents()
    local articles  = self:query('SELECT `id`, `name`, `parent_id` FROM `documents`');
    local map       = {};
    local contents  = {};

    -- 生成目录树算法
    local function add_node(node)
        if node.added then return end;

        local parent    = map[node.parent_id];
        local this      = { text = node.name, id = node.id, data = node.parent_id };
        if parent then
            if not parent.added then add_node(parent) end;
            parent.added.children = parent.added.children or {};
            table.insert(parent.added.children, this);
        else
            table.insert(contents, this);
        end

        node.added = this;
    end

    for _, info in ipairs(articles) do
        map[info.id] = info;
    end

    for _, info in ipairs(articles) do
        add_node(info);
    end

    return contents;
end

-- 取得一个文档
function documents:get(id)
    local find = self:query('SELECT `author`, `modify_user`, `modify_time`, `content` FROM `documents` WHERE id=?1', id);
    local info = find[1];
    
    if not info then return end;

    local names = M('user'):get_names(info.author, info.modify_user);
    info.author = names[info.author] or '神秘人';
    info.modify_user = names[info.modify_user] or '神秘人';
    return info;
end

-- 新建文档
function documents:add(name, parent_id)
    self:exec(
        "INSERT INTO `documents`(`name`, `parent_id`, `author`, `modify_user`, `modify_time`, `content`) VALUES(?1, ?2, ?3, ?3, NOW(), '')",
        name, parent_id, session.uid);
    return true;
end

-- 编辑文档
function documents:edit(id, content)
    self:exec(
        "UPDATE `documents` SET modify_user=?1, modify_time=NOW(), content=?2 WHERE id=?3",
        session.uid, content, id);
    return true;
end

-- 移动文档
function documents:move(id, to)
    self:exec('UPDATE `documents` SET parent_id=?1 WHERE id=?2', to, id);
    return true;
end

-- 删除文档
function documents:delete(id, name)
    local find = self:query('SELECT `name`, `author`, `parent_id` FROM `documents` WHERE id=?1', id);
    local info = find[1];

    if not info then return false, '文档不存在！' end;
    if name ~= info.name then return false, '文档名不匹配！' end;
    if info.author ~= session.uid and (not session.is_su) then return false, '您不能删除其他人的文档' end;
    
    self:exec('UPDATE `documents` SET `parent_id`=?1 WHERE `parent_id`=?2', info.parent_id, id);
    self:exec('DELETE FROM `documents` WHERE id=?1', id);
    return true;
end

return documents;