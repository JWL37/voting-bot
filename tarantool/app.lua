-- app.lua

-- Подключаем Tarantool
box.cfg {
    listen = 3301, -- Порт для подключения
}

-- Создаём последовательность, если она не существует
if not box.sequence.poll_id_seq then
    box.schema.sequence.create('poll_id_seq', {min = 1, start = 1})
end

-- Создаём спейсы, если они не существуют
if not box.space.polls then
    box.schema.space.create('polls')
    box.space.polls:format({
        {name = 'id', type = 'unsigned'}, -- Используем unsigned для автоинкрементного ID
        {name = 'options', type = 'map'},
        {name = 'creator', type = 'string'},
        {name = 'active', type = 'boolean'},
    })
    box.space.polls:create_index('primary', {
        parts = {'id'},
        type = 'tree',
        sequence = 'poll_id_seq', -- Привязываем последовательность к индексу
    })
end

if not box.space.voters then
    box.schema.space.create('voters')
    box.space.voters:format({
        {name = 'poll_id', type = 'unsigned'}, -- Используем unsigned для ID
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

-- Выводим сообщение о успешном создании спейсов
print("Спейсы 'polls' и 'voters' успешно созданы.")