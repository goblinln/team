-- ========================================================
-- @File    : /controller/install/setup.lua
-- @Brief   : 安装部署
-- @Author  : Leo Zhao
-- @Date    : 2018-07-19
-- ========================================================
local setup = {};

-- 数据库配置
setup.struct = {
    users = [[
        `id` INTEGER NOT NULL AUTO_INCREMENT,
        `account` VARCHAR(64) UNIQUE NOT NULL,
        `name` VARCHAR(32) UNIQUE NOT NULL,
        `avatar` VARCHAR(128) DEFAULT '/www/images/default_avatar.png',
        `pswd` CHAR(32) NOT NULL,
        `is_su` BOOLEAN DEFAULT 0,
        `is_locked` BOOLEAN DEFAULT 0,
        `auto_login_expire` INTEGER UNSIGNED DEFAULT 0,
        PRIMARY KEY(`id`)]],
    documents = [[
        `id` INTEGER NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(64) NOT NULL,
        `parent_id` INTEGER DEFAULT -1,
        `author` INTEGER NOT NULL,
        `modify_user` INTEGER NOT NULL,
        `modify_time` DATETIME,
        `content` TEXT,
        PRIMARY KEY(`id`)]],
    projects = [[
        `id` INTEGER NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(32) UNIQUE NOT NULL,
        `repo` VARCHAR(128) DEFAULT '',
        PRIMARY KEY(`id`)]],
    project_members = [[
        `uid` INTEGER NOT NULL,
        `pid` INTEGER NOT NULL,
        `role` INTEGER NOT NULL,
        `is_admin` BOOLEAN DEFAULT FALSE]],
    tasks = [[
        `id` INTEGER NOT NULL AUTO_INCREMENT,
        `pid` INTEGER DEFAULT -1,
        `creator` INTEGER NOT NULL,
        `assigned` INTEGER NOT NULL,
        `cooperator` INTEGER NOT NULL,
        `name` VARCHAR(32) NOT NULL,
        `weight` INTEGER DEFAULT 1,
        `tags` VARCHAR(64) DEFAULT '[]',
        `start_time` DATE NOT NULL,
        `end_time` DATE NOT NULL,
        `archive_time` INTEGER DEFAULT -1,
        `status` INTEGER DEFAULT 1,
        `content` TEXT,
        PRIMARY KEY(`id`)]],
    task_events = [[
        `timepoint` DATETIME DEFAULT NOW(),
        `tid` INTEGER NOT NULL,
        `uid` INTEGER NOT NULL,
        `event` INTEGER NOT NULL,
        `addition` VARCHAR(64)]],
    task_comments = [[
        `id` INTEGER NOT NULL AUTO_INCREMENT,
        `timepoint` DATETIME DEFAULT NOW(),
        `tid` INTEGER NOT NULL,
        `uid` INTEGER NOT NULL,
        `comment` TEXT,
        PRIMARY KEY(`id`)]],
    files = [[
        `id` INTEGER NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(128) NOT NULL,
        `path` VARCHAR(512) NOT NULL,
        `creator` INTEGER NOT NULL,
        `size` INTEGER NOT NULL,
        `upload_time` DATETIME DEFAULT NOW(),
        PRIMARY KEY(`id`)]]
};

-- 部署时写入数据
setup.build_in = {
    users = {
        {
            account = 'admin',
            name    = '管理员',
            pswd    = md5('team'),
            is_su   = 1,
        }
    }
};

-- 安装
function setup:index(req, rsp)
    local processing, code;
    local use_drop = req.get.override;

    xpcall(function()
        local db = M('base');

        for name, sql in pairs(self.struct) do
            log.info('CREATE TABLE : ' .. name);

            if use_drop then
                code = string.format('CREATE TABLE `%s`(%s);', name, sql);            
                db:exec(string.format('DROP TABLE IF EXISTS `%s`', name));
            else
                code = string.format('CREATE TABLE IF NOT EXISTS `%s`(%s);', name, sql);  
            end
            
            db:exec(code);
        end

        for name, rows in pairs(self.build_in) do         
            for _, row in ipairs(rows) do
                code = 'INSERT INTO `' .. name .. '`(';
                
                local vals  = {};
                local param = ' VALUES('

                for key, val in pairs(row) do
                    table.insert(vals, val);

                    code    = code .. '`' .. key .. '`,';
                    param   = param .. '?' .. #vals .. ',';
                end

                code = string.sub(code, 1, -2) .. ') ' ..  string.sub(param, 1, -2) .. ')';
                db:exec(code, unpack(vals));
            end
        end

        local lock = io.open(config.app.install_lock, 'w+');
        lock:write('Remove this file & visite /install/setup/?override=1 will resetup databases');
        lock:close();

        rsp:html('install/success.html');
    end, function(err)
        log.error(err);

        err = string.gsub(err, '\n', '<br>');
        rsp:html('install/failed.html', { err_msg = '[RUN SQL FAILED] : ' .. code .. '<br><br>' .. err });
    end);
end

return setup;