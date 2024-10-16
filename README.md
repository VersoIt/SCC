# Проект с двумя клиентами и одним сервером

```
Первый клиент считывает огромное количество информации при нажатии на кнопку и при этом он должен работать без сбоев даже при большом количестве одновременных нажатий.
Второй клиент создает много соединений с разных горутин и шлет данные в формате Base64. Сервер должен потоково обрабатывать данные и сохранять их в БД, при этом работать без сбоев.
```

Папка client_receiver - проект клиента с получением данных с сервера и дальнейшей записью в файл
Папка client_sender - проект клиента с отправкой данных через сокет с горутин
Папка server - сам сервер, слушающий http(для отправки данных) и tcp-соединения(для подключения через socket)

```
Клиенты реализованы с использованием фреймворка QT
Сервер разработан по принципам чистой архитектуры
```

### Для запуска приложения:

```
docker-compose -f server/docker-compose.yml up
Миграции для БД применяются автоматически при старте сервиса
make build && make run
```
