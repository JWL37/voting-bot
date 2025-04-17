
box.cfg {
    listen = 3301, 
}

if not box.sequence.poll_id_seq then
    box.schema.sequence.create('poll_id_seq', {min = 1, start = 1})
end

if not box.space.polls then
    box.schema.space.create('polls')
    box.space.polls:format({
        {name = 'id', type = 'unsigned'},
        {name = 'options', type = 'map'},
        {name = 'creator', type = 'string'},
        {name = 'active', type = 'boolean'},
    })
    box.space.polls:create_index('primary', {
        parts = {'id'},
        type = 'tree',
        sequence = 'poll_id_seq', 
    })
end

if not box.space.voters then
    box.schema.space.create('voters')
    box.space.voters:format({
        {name = 'poll_id', type = 'unsigned'}, 
        {name = 'user_id', type = 'string'},
    })
    box.space.voters:create_index('primary', {
        parts = {'poll_id', 'user_id'},
        type = 'hash',
    })
    box.space.voters:create_index('poll_id_index', {
        parts = {'poll_id'},
        type = 'tree',
    })
end

print("Спейсы 'polls' и 'voters' успешно созданы.")