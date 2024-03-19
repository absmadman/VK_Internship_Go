# Rest API в рамках VK Internship
REST API service implementation Go

# Быстрый старт
Для быстрого старта вам нужно иметь установленный Docker и docker-compose, перейти в корневую \
папку проекта(VK_Internship_Go) и выполнить команду `make docker-compose-api`, проект развернется \
и будет принимать запросы по адресу `localhost:8080`

# Установка и запуск
Проект запускается в Docker контейнере с помощью docker-compose, для удобства Makefile содержит следующие зависимости:

`make help` - выводит подсказки по зависимостям \
`make docker-compose-api` - запускает приложение в контейнере с помощью compose \
`make clean-pgdata` - очищает данные из базы данных\
`make docker-stop-api` - останавливает работу контейнеров\
`make docker-clean-api` - удаляет контейнеры\
`make server-logs` - выводит логи с контейнера сервера \
`make database-logs` - выводит логи с контейнера базы данных\
`make all-logs` - выводит все логи вместе\

# О проекте
Проект был создан с помощью пакета Gin для роутинга запросов, также использовалась библиотека для работы с LRU \
кешированием, в качестве хранилища данных был использован Postgresql. \
Проект запускается в двух Docker контейнерах один для работы API, а второй для базы данных. \
Соединение между конейтенерами осуществляется с помощью links в docker-compose файле и переменных окружения

Регулировка различных функций по типу имени базы данных, размеры кешей, пользователя базы данных происходит с помощью \
параметров в файле ***.env*** переменные из этого файла передаются в docker-compose

Проект реализует следующие "ручки" и методы (см. примеры):
* /users 
    * GET с параметрами name(строка) или id(целое число): 
      * Например: `curl -X GET 'http://localhost:8080/users?id=1 или /users?name=Dmitriy`
    * POST без параметров, принимает json с полями name(строка) и balance(число с плавающей точкой):
      * Например: `curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d {'"name": "Dmitriy", "balance": 2.5'}`
    * PUT с параметрами name(строка) или id(целое число) а также json с одним или несколькими полями которые следует обновить
      * Например: `curl -X PUT 'http://localhost:8080/users?name=Dmitriy' -H "Content-Type: application/json" -d '{"name": "Misha", "balance": 15.3}'`
    * DELETE с параметрами name(строка) или id(целое число)
      * Например: `curl -X DELETE 'http://localhost:8080/users?id=1' или /users?name=Misha`
* /quests
    * GET с параметрами name(строка) или id(целое число)
      * Например: `curl -X GET 'http://localhost:8080/quests?id=1 или /quests?name=Quest1`
    * POST без параметров, принимает json с полями name(строка) и cost(число с плавающей точкой)
      * Например: `curl -X POST http://localhost:8080/quests -H "Content-Type: application/json" -d {'"name": "Quest1", "cost": 2.5'}`
    * PUT с параметрами name(строка) или id(целое число) а также json с полями которые следует обновить
      * Например: `curl -X PUT 'http://localhost:8080/quests?name=Quest1' -H "Content-Type: application/json" -d '{"name": "Quest2", "cost": 15.3}'`
    * DELETE с параметрами name(строка) или id(целое число)
      * Например: `curl -X DELETE 'http://localhost:8080/quests?id=1' или /quests?name=Quest1`
* /events
    * GET с параметрами user_id(целое число)
      * Например: `curl -X GET 'http://localhost:8080/events?user_id=2'`
    * POST без параметров, принимает json с полями user_id(целое число) и quest_id(целое число)
      * Например: `curl -X POST 'http://localhost:8080/events' -H "Content-Type: application/json" -d '{"user_id": 2, "quest_id": 1}'`

# Функционал
Помимо базовых CRUD операций также было реализован дополнительный функционал:
* Кеширование с настраиваемым размером кеша
* Логирование в os.Stdout контейнера, который в последствии можно посмотреть с помощью зависимости make all-logs
* GET, PUT, DELETE запросы поддерживают операции не только по id но и по имени
* Проект запускается в Docker контейнерах
* Проект разбит на директории и файлы
* Корректная обработка ошибок                                            
