# sre-solution-cup-2023

## Run

```shell
make run
```

Перейти по адресу http://localhost:3000

## TODO

- [ ] Следить за изменениями конфигурации и при их выявлении также перестраивать расписание
- [ ] Продление работ
- [ ] Перенос работ
- [x] Если работы в желаемое время не могут быть проведены сообщать об этом, предлагая несколько ближайших по времени вариантов с учетом окон времени из конфигурации сервиса.
- [x] Гарантировать, что в одной зоне доступности в любой момент времени проводятся не более, чем одни работы.
- [x] При полной невозможности проведения работ сообщать об этом сообщением о невозможности при приеме заявки.
- [x] Критичные работы помещать в расписание вне очереди, принудительно отменяя обычные ручные работы.

