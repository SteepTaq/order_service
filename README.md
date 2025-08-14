Видео:
[Ссылка на видео](https://drive.google.com/file/d/182X4hIK8aWvbDQnVxgtatSk0prEqi8k5/view?usp=sharing)

Быстрый старт:
make all

другие команды:
1) поднять postgres и kafka:
make infra-up
2) Миграции БД:
make migrate-up
make migrate-down
3) Запуск сервиса:
make run

этой командой можно отправить несколько тестовых заказов в кафку:
make producer


