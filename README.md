# medods-service

Запуск сервиса

-docker-compose build  

-docker-compose up  (возможно, надо будет 2 раза прописать ее)

После запуска у сервиса будет два маршрута
localhost:8080/api/auth/create - Принимает на вход json с guid и отправляет два токена, заполняет cookie.


localhost:8080/api/auth/refresh - обрабатывает cookie с названием refreshToken и заголовок Email, обновляет access и refresh токен