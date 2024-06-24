# sync-manager

## Остается сделать: 
### Взаимодействие
1. Сделать CRUD для добавления систем (добавляем систему,
чтобы отправлять ей сообщение из консумера по веб хуку)
   Нужны таблицы:
   | source: id, name, name_lat, code(index), receive_method[amqp/http]
   | route: id, name(index), url, description, source_fk
### Консумер
1. Логика сохранения сообщений батчем из 10 штук [**DONE**] 
2. Отправка по вебхуку на адрес системы
### REST
1. Отправка в RabbitMq при сохранении сообщения [**DONE**]
2. Потвторная отправка всех неудачных сообщений с retried = false
### Cron
1. Раз в 5 минут отправлять неудачные сообщения с dead = true && 
retried = false

Кажется, я запутался. 
Можно биндить как в одну, так и в другую сторону.