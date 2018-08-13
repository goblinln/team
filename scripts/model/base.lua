-- ========================================================
-- @File    : /model/base.lua
-- @Brief   : 数据类基类定义
-- @Author  : Leo Zhao
-- @Date    : 2018-07-16
-- ========================================================
local base = {};

-- 配置，可在自己的Model中修改配置
base.settings	= {
    host	= config.database.host,
    port	= tonumber(config.database.port),
    user	= config.database.user,
    pswd	= config.database.pswd,
    db		= config.database.db,
    charset	= config.database.charset,
}

-----------------------------------------------------------

-- 取得最后一次插入记录的ID
function base:last_id()
    self:__init();
    return tonumber(self.__conn:insert_id());
end

-- 取得最近一次操作影响的行数
function base:affected_rows()
    self:__init();
    return self.__conn:affected_rows();
end

-- 更新数据库。INSERT/DELETE/UPDATE等非查询操作，无返回值，如果有参数请使用'?N'的方式
-- 如 self:exec('SELECT * FROM `users` WHERE name=?1 AND level=?2', name, level);
function base:exec(sql, ...)
    self:__init();
    self.__conn:query(self:__mksql(sql, ...));
end

-- 查询数据库。SELECT操作，带参数的同exec
function base:query(sql, ...)
    self:__init();
    self.__conn:query(self:__mksql(sql, ...));

    local res = self.__conn:use_result();
    local ret = {};

    while true do
        local row = res:fetch('a');
        if not row then break end;
        table.insert(ret, row);
    end

    return ret;
end

-----------------------------------------------------------

-- 初始化
function base:__init()
    if not self.__conn then
        self.__conn = T('mysql/mysql').config().connect(self.settings);
        if not _G.cleanup_db then
            _G.cleanup_db = T('mysql/mysql').C.mysql_server_end;
        end
    end
end

-- 生成MYSQL的值
function base:__mkval(v)
    local t = type(v);

    if t == 'nil' then
        return 'NULL';
    elseif t == 'boolean' then
        return v and 1 or 0;
    elseif t == 'string' then
        return "'" .. self.__conn:escape(v, string.len(v)) .. "'";
    elseif t == 'table' then
        if v.year and v.month and v.day then
            if v.hour then
                return string.format("'%04d-%02d-%02d %02d:%02d:%02d'", v.year, v.month, v.day, v.hour, v.min or 0, v.sec or 0);
            else
                return string.format("'%04d-%02d-%02d'", v.year, v.month, v.day);
            end
        else
            local serialized	= json.encode(v);
            local encrypt		= "'" .. self.__conn:escape(serialized, string.len(serialized)) .. "'";
            return encrypt;
        end
    else
        return tostring(v);
    end
end

-- 生成SQL
function base:__mksql(sql, ...)
    local count	= select('#', ...);
    if count == 0 then return sql end;

    local args	= {...};
    local param = {};

    for i = 1, count do param[tostring(i)] = self:__mkval(args[i]) end;
    local sql = string.gsub(sql, '%?(%d+)', param);
    return sql;
end

return base;