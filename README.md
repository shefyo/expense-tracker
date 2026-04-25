# Expense-tracker

Идея для pet-проекта взята отсюда:
https://roadmap.sh/projects/expense-tracker

Простая утилита для учёта расходов, написанная на Go. Работает из командной строки, хранит данные в JSON-файле рядом с собой.

---

## Установка

```bash
git clone https://github.com/shefyo/expense-tracker.git
cd expense-tracker
go build -o expense-tracker .
```

Если хочешь вызывать команду из любого места, скопируй бинарник в PATH:

```bash
sudo mv expense-tracker /usr/local/bin/
```

---

## Использование

### Добавить расход

```bash
expense-tracker add --description "Lunch" --amount 20
# Expense added successfully (ID: 1)

expense-tracker add --description "Taxi home" --amount 8.50 --category Transport
# Expense added successfully (ID: 2)
```

`--category` — необязательный флаг, но удобный.

---

### Посмотреть все расходы

```bash
expense-tracker list
# ID    Date          Description           Amount      Category
# 1     2024-08-06    Lunch                 $20.00      -
# 2     2024-08-06    Taxi home             $8.50       Transport
```

Можно фильтровать по категории:

```bash
expense-tracker list --category Transport
```

---

### Обновить расход

Если ошибся в сумме или описании:

```bash
expense-tracker update --id 1 --amount 22 --description "Lunch with colleague"
# Expense 1 updated successfully
```

Необязательно передавать все поля — обновится только то, что указал.

---

### Удалить расход

```bash
expense-tracker delete --id 2
# Expense deleted successfully
```

---

### Сводка

Общая по всем расходам:

```bash
expense-tracker summary
# Total expenses: $28.50
```

За конкретный месяц текущего года:

```bash
expense-tracker summary --month 8
# Total expenses for August: $28.50
```

Если для месяца установлен бюджет, сводка покажет остаток или предупреждение о превышении.

---

### Бюджет

```bash
expense-tracker budget --month 8 --amount 200
# Budget for August set to $200.00

expense-tracker summary --month 8
# Total expenses for August: $28.50
# Budget for August: $200.00
# Remaining budget: $171.50
```

Если расходы превысят бюджет:

```
⚠  Warning: you have exceeded your budget by $15.00
```

---

### Экспорт в CSV

```bash
expense-tracker export
# Expenses exported to expenses.csv

expense-tracker export --output ~/Documents/august.csv
```

Формат файла:

```
ID,Date,Description,Amount,Category
1,2024-08-06,"Lunch",20.00,-
2,2024-08-06,"Taxi home",8.50,Transport
```

---

## Хранение данных

Все данные хранятся в файле `expenses.json` в той директории, из которой запускается команда. Файл создаётся автоматически при первом добавлении расхода. Можно копировать, бэкапить или редактировать руками — формат простой.

---

## Команды

| Команда   | Описание                          |
|-----------|-----------------------------------|
| `add`     | Добавить расход                   |
| `update`  | Обновить существующий расход      |
| `delete`  | Удалить расход по ID              |
| `list`    | Показать все расходы              |
| `summary` | Итог по всем или по месяцу        |
| `budget`  | Установить бюджет на месяц        |
| `export`  | Экспортировать расходы в CSV      |

Подробная справка:

```bash
expense-tracker help
```

---

## Требования

- Go 1.21+
- Никаких внешних зависимостей

---

## Лицензия

MIT
