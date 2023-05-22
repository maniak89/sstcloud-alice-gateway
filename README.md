# sstcloud-alice-gateway
Шлюз для Алисы до SST Cloud

# Переменные окружения

Так же поддерживается установка значение через .env файл созданный в рабочем каталоге

| Переменная               | Описание                                                                               | Значение по умолчанию                            | Обязательно             |
|--------------------------|----------------------------------------------------------------------------------------|--------------------------------------------------|-------------------------|
| HTTP_ADDRESS             | string                                                                                 | Адрес на котором поднять REST интерфейс          | :80                     | Да |
| LOGGER_DISABLE_TIMESTAMP | bool                                                                                   | Не показывать в логах timestamp                  | false                   | Нет |
| LOGGER_ENABLE_CALLER     | bool                                                                                   | Показывать откуда был вызов логера               | false                   | Нет |
| LOGGER_ENABLE_CONSOLE    | bool                                                                                   | Форматирование в "человеческий" вид с подсветкой | false                   | Нет |
| LOGGER_LEVEL             | string                                                                                 | Уровень логирования                              | info                    | Нет |
| SST_TIMEOUT              | Таймаут до SST                                                                         | 5s                                               | Нет                     |
| SST_URL                  | Адрес REST SST                                                                         | https://api.sst-cloud.com                        | Нет                     |

# OAuth2
Для корректной работы с Yandex.Cloud и Алисой в частности требуется иметь некий OAuth2 аутификатор. 
Сервис ожидает X-User-Id по которому найдет в бд учетные записи пользователя и будет использовать их для обращения к sst


