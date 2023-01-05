### Перед началом работы
1. Задайте переменные окружения *__TARANTOOL_LOGIN__* и *__TARANTOOL_PASSWORD__*
2. Запустите команду docker-compose up -d из каталога docker
3. Соберите прокси-сервер, который будет принимать входящие запросы и распределять
   их между инстансами приложения, запустив команду 
>__go build proxy.go__
4. Соберите приложение, исполнив команду 
 >__go build main.go__

Запуск приложения:
1. Запустите по одному инстансу приложения, указав в качестве аргументов endpoint 
   для Tarantool и Host, по которому приложение будет слушать запросы от прокси
   Пример запуска 
>__./app_bin_name -storage-type tarantool -tarantool-endpoint localhost:3301 -http-endpoint localhost:8080__
2. Запустите прокси ./proxy
   Пример запуска
>__./app_bin_name -app-hosts-list localhost:8080,localhost:8081 -proxy-host-and-port localhost:9000__

__ВАЖНО!!!!__ В качестве хранилища данных приложение может использовать оперативную память. 
Для использования этой возможности, достаточно указать в качестве значения для параметра
-storage-type inmemory или не указывать данный параметр вовсе. В этом случае, команда 
запуска приложения будет выглядеть так:

>__./app_bin_name -storage-type inmemory -http-endpoint localhost:8080__

или так

>__./app_bin_name -http-endpoint localhost:8080__