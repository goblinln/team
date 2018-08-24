-- ========================================================
-- @File    : /controller/dashboard/files.lua
-- @Brief   : 上传文件管理
-- @Author  : Leo Zhao
-- @Date    : 2018-07-14
-- ========================================================
local files = {};

-- 分享首页
function files:index(req, rsp)
    local shared = M('base'):query([[SELECT * FROM `files`]]);
    local names = M('user'):get_names();
    local min_kb = 1024;
    local min_mb = min_kb * 1024;
    local min_gb = min_mb * 1024;

    for _, info in ipairs(shared) do
        info.creator_name = names[info.creator] or '神秘人';
        if info.size < min_mb then
            info.size = string.format("%.3f KB", info.size * 1.0 / min_kb);
        elseif info.size < min_gb then
            info.size = string.format("%.3f MB", info.size * 1.0 / min_mb);
        else
            info.size = string.format("%.3f GB", info.size * 1.0 / min_gb);
        end
    end

    rsp:html('dashboard/files/index.html', {
        dashboard_menu = 'files',
        shared = shared,
    });
end

-- 上传API
function files:upload(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local data = self:do_upload(req.file)[1];
    if not data then return rsp:json{ ok = false } end;
    if (not data.ok) or (not req.post.is_share) then return rsp:json(data) end;

    if req.post.desc then
        req.post.desc = string.gsub(req.post.desc, '\n', '<br>');
    end

    xpcall(function()
        M('base'):exec([[
            INSERT INTO `files`(`name`, `path`, `creator`, `size`, `desc`)
            VALUES(?1, ?2, ?3, ?4, ?5)]], data.name, data.url, session.uid, data.size, req.post.desc);
    end, function(e)
        data.ok = false;
    end);

    rsp:json(data);
end

-- 删除上传文件
function files:delete(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local shared = M('base'):query([[SELECT * FROM `files` WHERE `id`=?1]], req.post.id)[1];
    if not shared then
        return rsp:json{ ok = false, err_msg = '文件不存在' };
    elseif shared.creator ~= session.uid and not session.is_su then
        return rsp:json{ ok = false, err_msg = '您不能删除这个文件' };
    end

    os.rm('.' .. shared.path);
    M('base'):exec([[DELETE FROM `files` WHERE `id`=?1]], shared.id);
    rsp:json{ok = true};
end

-------------------------- 内部调用 --------------------------

-- 上传处理
function files:do_upload(file_map)
    local dir = 'www/upload/' .. session.uid;
    if not os.exists(dir) then os.mkdir(dir) end;

    local uploaded = { };
    for name, path in pairs(file_map) do
        local to = dir .. '/' .. os.time() .. '_' .. name;
        if os.cp(path, to) then table.insert(uploaded, { ok = true, name = name, url = '/' .. to, size = os.filesize(to) }) end;
    end

    return uploaded;
end

return files;