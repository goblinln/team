-- ========================================================
-- @File    : /controller/dashboard/files.lua
-- @Brief   : 上传文件管理
-- @Author  : Leo Zhao
-- @Date    : 2018-07-14
-- ========================================================
local file = {};

-- 上传API
function file:upload(req, rsp)
    if req.method ~= 'POST' then return rsp:error(405) end;

    local dir = 'www/upload';
    if not os.exists(dir) then os.mkdir(dir) end;

    dir = dir .. '/' .. session.uid;    
    if not os.exists(dir) then os.mkdir(dir) end;

    local uploaded = { };
    for name, path in pairs(req.file) do
        local to = dir .. '/' .. os.time() .. '_' .. name;
        if os.cp(path, to) then table.insert(uploaded, { ok = true, name = name, url = '/' .. to}) end;
    end

    rsp:json(uploaded[1] or {});
end

return file;