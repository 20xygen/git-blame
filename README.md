# Blame

---

## Общее описание
Blame - это CLI linux-утилита  для анализа статистики авторов Git репозитория.
Работа программы основана на парсинге раличиных команд `git`.

---

## Сборка проекта

1. Убедитесь, что у вас установлен Golang

2. Клонируйте репозиторий:
```bash
git clone https://github.com/20xygen/git-blame.git
cd git-blame
```

3. Установите утилиту:
```bash
go install ./cmd/blame/...
export PATH=$GOPATH/bin:$PATH
```

4. Для сохранение логов (опционально):
```bash
sudo mkdir /var/log/blame
sudo chown $USER:$USER /var/log/blame
```

## Использование
```
blame [flags]

Flags:
  -x, --exclude strings       Exclude glob patterns
  -e, --extensions strings    File extensions filter (comma-separated)
  -f, --format string         Output format (one of 'pretty', 'tabular', 'json', 'json-lines', 'csv')' (default "tabular")
  -h, --help                  help for blame
  -l, --languages strings     Languages filter (comma-separated)
  -o, --order-by strings      Sort key as comma-separated list of 'lines', 'commits', 'names' or 'files' (default [lines,commits,files])
  -r, --repository string     Git repository path (default ".")
  -t, --restrict-to strings   Restrict-to glob patterns
  -R, --revision string       Git revision (default "HEAD")
  -C, --use-committer         Use committer instead of author
```

---

## Примеры использования

Для тестовых запусков распакуем бандлы:

```bash
cd /path/to/sandbox/dir
git clone /path/to/git-blame/test/integration/testdata/bundles/go-cmp.bundle
cd go-cmp
```

#### Простой запуск

```bash
blame
```

```
Name                   Lines Commits Files
Joe Tsai               13818 94      54
colinnewell            130   1       1
A. Ishikawa            92    1       2
Roger Peppe            59    1       2
Tobias Klauser         35    2       3
178inaba               27    2       5
Kyle Lemons            11    1       1
Dmitri Shuralyov       8     1       2
ferhat elmas           7     1       4
Christian Muehlhaeuser 6     3       4
k.nakada               5     1       3
LMMilewski             5     1       2
Ernest Galbrun         3     1       1
Ross Light             2     1       1
Chris Morrow           1     1       1
Fiisio                 1     1       1

```

#### Продвинутое использование

```bash
blame --repository . --extensions .go,.md --order-by files --format pretty
```

```
+------------------------+---------+-------+-------+
| NAME                   | COMMITS | FILES | LINES |
+------------------------+---------+-------+-------+
| Joe Tsai               |      92 |    49 | 12154 |
| 178inaba               |       2 |     4 |    11 |
| ferhat elmas           |       1 |     4 |     7 |
| Christian Muehlhaeuser |       3 |     4 |     6 |
| k.nakada               |       1 |     3 |     5 |
| Roger Peppe            |       1 |     2 |    59 |
| Tobias Klauser         |       1 |     2 |    33 |
| Dmitri Shuralyov       |       1 |     2 |     8 |
| LMMilewski             |       1 |     2 |     5 |
| colinnewell            |       1 |     1 |   130 |
| A. Ishikawa            |       1 |     1 |    36 |
| Kyle Lemons            |       1 |     1 |    11 |
| Ernest Galbrun         |       1 |     1 |     3 |
| Ross Light             |       1 |     1 |     2 |
| Chris Morrow           |       1 |     1 |     1 |
| Fiisio                 |       1 |     1 |     1 |
+------------------------+---------+-------+-------+
```

---

## Структура проекта

#### 1. **cmd**
- **blame**
    - [`main.go`](cmd/blame/main.go) — инициализация и запуск.

#### 2. **configs**
- [`language_extensions.json`](configs/language_extensions.json) — маппинг языков программирования на расширения файлов.

#### 3. **internal** .
- **cli** — обработка командной строки.
    - [`cli.go`](internal/cli/cli.go) — интерфейс команды.
- **format** — форматирование вывода.
    - [`auto.go`](internal/format/auto.go) — автоматическое определение формата.
    - [`format.go`](internal/format/format.go) — реализация форматов вывода.
- **statistics** — сбор статистики.
    - [`params.go`](internal/statistics/params.go) — структуры параметров сбора статистики.
    - [`process.go`](internal/statistics/process.go) — фильтрация и сбор статистики.
    - [`statistics.go`](internal/statistics/statistics.go) — структуры единиц статистики.
- **utils**
    - [`errors.go`](internal/utils/errors.go) — описание ошибок.
    - [`languages.go`](internal/utils/languages.go) — работа с расширениями.
    - [`logger.go`](internal/utils/logger.go) — логирование.
    - [`utils.go`](internal/utils/utils.go) — прочее.

#### 4. **pkg**
- **commands** — работа с системными командами.
    - [`commands.go`](pkg/commands/commands.go) — запуск команд.
    - [`errors.go`](pkg/commands/errors.go) — описание ошибок.
- **files** — работа с файловой системой.
    - [`directories.go`](pkg/files/directories.go) — структуры и методы для взаимодействия с директориями.
    - [`files.go`](pkg/files/files.go) — структуры и методы для взаимодействия с файлами.
    - [`errors.go`](pkg/files/errors.go) — описание ошибок.
- **parsing** — парсинг команды `git blame`.
    - [`output.go`](pkg/parsing/output.go) — структуры единиц вывода.
    - [`parsing.go`](pkg/parsing/parsing.go) — основной процесс обработки.

#### 5. **test** и **tools**
- заимствованные из проекта курса Go в МФТИ файлы для тестирования.
