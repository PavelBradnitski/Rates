# Валютный API для NBRB
## Данный проект собирает данные о курсах валют НБРБ по отношению к белорусскому рублю ежедневно, используя API НБРБ:
## https://api.nbrb.by/exrates/rates?periodicity=0.

## Проект включает в себя HTTP API сервер, который обрабатывает несколько GET-запросов для получения данных о курсах валют.

## Технологии
### Проект разработан с использованием следующих технологий:

### Go (Golang) - для разработки серверной логики.
### MySQL - для хранения данных о курсах валют.
### Docker - для контейнеризации проекта и упрощения развертывания.
## Установка
### Для работы с проектом на локальной машине выполните следующие шаги.

### 1. Клонировать репозиторий
### git clone git@github.com:PavelBradnitski/Rates.git
### cd Rates
## 2. Настроить Docker
### Проект использует Docker для развертывания всех зависимостей (например, MySQL).

### Запустите контейнеры с помощью Docker Compose:
### docker-compose up --build

## API
### Доступные запросы
### 1. Получение всех записей курсов валют, собранных с API НБРБ

### Метод: GET
### URL: http://localhost:8080/rate
### Описание: Возвращает все записи, собранные с API НБРБ.
### Пример запроса:
### GET http://localhost:8080/rate

### Получение записей за выбранный день

### Метод: GET
### URL: http://localhost:8080/rate/{date}
### Описание: Возвращает записи за указанный день. Дата должна быть в формате YYYY-MM-DD.
### Пример запроса:
### GET http://localhost:8080/rate/2025-01-12
### Пример ответа:
### json
[
  {
    "Cur_ID": 1,
    "Date": "2025-01-12",
    "Cur_Abbreviation": "USD",
    "Cur_Scale": 1,
    "Cur_Name": "Доллар США",
    "Cur_OfficialRate": 2.5
  },
  ...
]
### Формат ответа
### Ответы на запросы будут представлены в формате JSON.

### Пример ответа на запрос всех курсов валют:

json
[
  {
    "Cur_ID": 1,
    "Date": "2025-01-12",
    "Cur_Abbreviation": "USD",
    "Cur_Scale": 1,
    "Cur_Name": "Доллар США",
    "Cur_OfficialRate": 2.5
  },
  {
    "Cur_ID": 2,
    "Date": "2025-01-12",
    "Cur_Abbreviation": "EUR",
    "Cur_Scale": 1,
    "Cur_Name": "Евро",
    "Cur_OfficialRate": 3.0
  }
]
### Примечания
## Проект автоматически обновляет данные о курсах валют один раз в сутки в 03:00.
## Сервер работает по умолчанию на порту 8080.