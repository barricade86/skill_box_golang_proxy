box.cfg
{
    pid_file = nil,
    background = false,
    log_level = 5
}

local function init()
    local users = box.schema.space.create('users', {
        if_not_exists = true
    })
    users:format({
                                {name = 'id', type = 'unsigned', is_nullable = false},
                                {name = 'user_name', type = 'string', is_nullable = false},
                                {name = 'age', type = 'unsigned', is_nullable = false},
                                {name = 'friends', type = 'array', is_nullable = true}
                             })
    users:create_index('primary', {unique = true, sequence=true})
end

box.once('init', init)