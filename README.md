# Проект с двумя клиентами и одним сервером

```
Первый клиент считывает огромное количество информации при нажатии на кнопку. Он должен работать без сбоев даже при большом количестве одновременных нажатий.
Второй клиент создает много соединений с разных горутин и шлет данные в формате Base64. Сервер должен потоково обрабатывать данные, декодировать их из Base64 и сохранять в БД работая без сбоев.
```

client_receiver - клиент с получением данных с сервера и дальнейшей записью в файл  
client_sender - клиент с отправкой данных через сокет с горутин  
server - сервер, слушающий http(для отправки данных клиенту) и tcp-соединения(для подключения через socket)  

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
