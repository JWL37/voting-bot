### Команды
```
poll create "Вопрос?" "Вариант 1" "Вариант 2" ...
poll vote <ID голосования> <номер варианта>
poll results <ID голосования>
poll close <ID голосования>
poll delete <ID голосования>
```

### Запуск
сначала ставим mattermost(я делал через докер).Можно сделать командой:
```
    make mattermost
```
потом в `.env` заполняем поля с названием команды,канала и токена бота(получить токен можно внутри приложения mattermost)
Остаётся запустить 
```
    make
```

P.S. бота также хотел засунуть в докер,но он ругался на github.com/tarantool/go-openssl.Пробовал и многоступенчатую,и обычную сборку,также пробовал скачивать в ось внутри докера openssl,но в итоге пофиксить эту проблему не смог.

### Линтер

Для корректной работы линтера нужна версия 2.0 и выше:

```
➜  voting-bot git:(main) ✗ golangci-lint --version
golangci-lint has version 2.0.2 built with go1.24.1 from 2b224c2c on 2025-03-25T21:36:18Z
```

Запуск линтера:
```
    make lint 
```
