СУБД: PostgreSQL
## Список вопросов
1) Какой счет считать резервным? Принято решение, резервный счет имеет ID= -1

## Инструкция по запуску
1) `git clone https://github.com/HennOgyrchik/turbo-carnival.git`
2) `cd turbo-carnival`
3) `docker compose up -d`
4) `docker compose exec -T  postgres psql -U test -d postgres < <(cat internal/postgresql/struct_db.sql)`

## Описание методов
Все методы `PUT` производят запись в таблицу `transactions` о тех или иных изменениях.

Операции резервирования средств и признания выручки могут выполняться не более одного раза для каждого заказа.
Т.е. номер заказа уникален для операции резервирования и признания выручки.

Если данные введены корректно и операция прошла успешно, то возвращается http-код `200 OK`, иначе `400 Bad Request`.

### 1) (GET) /balance
Возвращает баланс пользователя по его ID.

Принимает:
```json
{
  "user_id":5
}
```

Возвращает:
```json
{
  "Cash": 27
}
```

### 2) (PUT) /replenish
Производит пополнение счета пользователя на заданную сумму. Если пользователь с указанным id отствутствует, создает нового.

Принимает:
```json
{
  "user_id":2,
  "count":7
}
```

### 3) (PUT) /reserve
Списывает указанную сумму со счета пользователя и зачисляет ее пользователю с ID= -1.

Принимает:
```json
{
  "user_id":4,
  "service_id":1,
  "order_id":77,
  "count":10
}
```

### 4) (PUT) /revenue
Списывает с резервного счета указанную сумму, только если существует соответствующая запись
о резервировании средств.

Для корректного выполнения все поля данного метода должны быть равны полям метода `/reserve` соответственно.
Если значения полей разнятся, списание не произойдет.

Принимает:
```json
{
  "user_id":4,
  "service_id":1,
  "order_id":77,
  "count":10
}
```
### 5) (GET) /monthly_report
Формирует отчет с указанием сумм выручки по каждой из предоставленной услуге за соответствующий месяц.

Принимает:
```json
{
  "date":"2022-11"
}
```

Возвращает ссылку на скачивание отчета в формате `.csv`

### 6) (GET) /report
Скачивает сформированный отчет
