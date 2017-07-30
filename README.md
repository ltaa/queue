# Состав
Проект содержит следующие образы:    

- rabbit
- queue_db
- queue

# Сборка
Для сборки текущего проекта необходимо находясь в папке проекта выполнить:
```
docker build -t queue .

cd docker/queue_db && docker build -t queue_db .

cd ../rabbit && docker build -t rabbit .
```

# Запуск
Запуск с помощью docker-compose
```
docker-compose -f compose.yml up
```

# db credentials:
```
postgres/postgres
```

