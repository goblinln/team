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
    end
end

-- 生成SQL
function base:__mksql(sql, ...)
    local count	= select('#', ...);
    if count == 0 then return sql end;

    local args	= {...};
    for i = 1, count do
        local data	= args[i];
        local dt	= type(data);
        local key	= '?' .. i;

        if dt == 'nil' then
            sql = string.gsub(sql, key, 'NULL');
        elseif dt == 'boolean' then
            sql = string.gsub(sql, key, data and 1 or 0);
        elseif dt == 'string' then
            str = self.__conn:escape(data, string.len(data));
            str = string.gsub(str, '%%', '%%%%');
            sql = string.gsub(sql, key, "'" .. str .. "'");
        elseif dt == 'table' then
            if data.is_datetime then
                local str = string.format('%04d-%02d-%02d %02d:%02d:%02d', data.year, data.month, data.day, data.hour, data.min, data.sec);
                sql = string.gsub(sql, key, "'" .. str .. "'");
            elseif data.is_date then
                local str = string.format('%04d-%02d-%02d', data.year, data.month, data.day);
                sql = string.gsub(sql, key, "'" .. str .. "'");
            elseif data.is_time then
                local str = string.format('%02d:%02d:%02d', data.hour, data.min, data.sec);
                sql = string.gsub(sql, key, "'" .. str .. "'");
            elseif data.is_timestamp then
                sql = string.gsub(sql, key, tostring(os.time(data)));
            else
                local serialized	= json.encode(data);  
                local encrypt		= self.__conn:escape(serialized, string.len(serialized));
                
                encrypt = string.gsub(encrypt, '%%', '%%%%');
                sql = string.gsub(sql, key, "'" .. encrypt .. "'");
            end
        else
            sql = string.gsub(sql, key, tostring(data));
        end
    end
    
    return sql;
end

return base;