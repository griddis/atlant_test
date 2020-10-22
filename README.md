### Генерация прото файлов

```bash
make proto
```

### запуск тестов

```bash
make test
```

### запуск сервиса

```bash
docker-compose -f quickstart.yml up --build
```

### проверка grpc

#### импорт данных

```bash
grpcurl -d '{ "files":[ "https://gist.githubusercontent.com/griddis/cf62cdfaa46d779dd6f7f7b436ceb77d/raw/e143b62d01bef89d77f3a9e1e1001d73027da224/gistfile1.txt" ] }' -import-path ./api/proto -proto service.proto -plaintext localhost:1443 service.Service/Fetch
```

можно указывать несколько файлов на загрузку, файлы загружать будет краулер по ограниченному количеству паралельных загрузок


#### реализация бесконечного срола 
```bash
grpcurl -d '{"limiter":{ "offsetbyid": "5f90a61f40abd998fd4051dd" } }' -import-path ./api/proto -proto service.proto -plaintext localhost:1443 service.Service/List
```

где offsetbyid указывается id записи, после которой надо делать выборку и даже если что-то появится более свежее или новое, что бы на странице не было повторов

### ограничение вывода

```bash
grpcurl -d '{"limiter":{ "limit": 10 } }' -import-path ./api/proto -proto service.proto -plaintext localhost:1443 service.Service/List
```

вывод 10 записей

### сортировка

```bash
grpcurl -d '{"sorter":{ "price": 1 } }' -import-path ./api/proto -proto service.proto -plaintext localhost:1443 service.Service/List
```

что бы произвести сортировку по 1 из полей, надо в обьекте sorter сделать ключ с именем поля сортировки, а значени поставить "1" - по возростанию, "-1" - по убыванию, можно указывать несколько сортировок