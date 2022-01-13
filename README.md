# sstcloud-alice-gateway
Шлюз для Алисы до SST Cloud

# Переменные окружения

|Переменная|Описание|Значение по умолчанию|Обязательно|
|---|---|---|---|
| HTTP_ADDRESS | string | Адрес на котором поднять REST интерфейс | :80 | Да |
| LOGGER_DISABLE_TIMESTAMP | bool | Не показывать в логах timestamp | false | Нет |
| LOGGER_ENABLE_CALLER | bool | Показывать откуда был вызов логера | false | Нет |
| LOGGER_ENABLE_CONSOLE | bool | Форматирование в "человеческий" вид с подсветкой | false | Нет |
| LOGGER_LEVEL | string | Уровень логирования | info | Нет |
| OAUTH2_ENABLED | Встроенная реализация oauth2 и соответствующие endpoint's активированы | false | Нет |
| OAUTH2_JWT_ALGO | Алгоритм подписи JWT | HS256 | Нет |
| OAUTH2_KEY | Ключ для проверки токена jwt | - | Да |
| OAUTH2_KEY_IN_BASE64 | Ключ в формате base64 | false | Нет |
| OAUTH2_TOKEN_STORE_FILE | Файл хранения oauth2 токенов при включенной встроенной реализации oauth2 | tokens | Нет |
| OAUTH2_USER_DOMAIN | Домен куда пользователя должно редиректить при включенной встроенной реализации oauth2 | http://localhost | Нет |
| OAUTH2_USER_ID | Идентификатор пользователя используемый при включенной встроенной реализации oauth2 | - | Да, если OAUTH2_ENABLED |
| OAUTH2_USER_SECRET | Секрет пользователя используемый при включенной встроенной реализации oauth2  | - | Да, если OAUTH2_ENABLED |
| SST_TIMEOUT | Таймаут до SST | 5s | Нет |
| SST_URL | Адрес REST SST | https://api.sst-cloud.com | Нет |

# OAuth2
Для корректной работы с Yandex.Cloud и Алисой в частности требуется иметь некий OAuth2 аутификатор.

Можно использовать сторонний, однако это не отменяет необходимости проверки токена на стороне шлюза, по этому необходимо его указывать в OAUTH2_KEY.

В случае если поднимать сторонний лень, можно использовать встроенный - OAUTH2_ENABLED. Он будет аутифицировать только одного
пользователя OAUTH2_USER_ID с токеном OAUTH2_USER_SECRET. При проверке токена также используется домен, для Yandex.Cloud 
будет https://social.yandex.net

