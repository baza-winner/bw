# Описание

bw -- общая bash-инфраструктура [проектов baza-winner](https://github.com/baza-winner)

# Использование

[Установка команды bw](https://github.com/baza-winner/bw/wiki/%D0%A3%D1%81%D1%82%D0%B0%D0%BD%D0%BE%D0%B2%D0%BA%D0%B0-%D0%BA%D0%BE%D0%BC%D0%B0%D0%BD%D0%B4%D1%8B-bw)

# Разработка

## Настройка окружения

1. [Установить команду bw](https://github.com/baza-winner/bw/wiki)

2. Развернуть рабочее место проекта bwdev

```
bw project bwdev
```

## Сборка

### Только сборка

```
bwdev build -m justBuild
```

### Cборка после тестирования

```
bwdev build
```

### Сборка после полного тестирования

Некоторые тесты (такие как тесты для `_spinner`) при обычном тестировании не запускаются из-за частого ложно отрицательного срабатывания. 

```
bwdev build -m testAll
```

### Просмотр изменений после сборки

```
diff tgz/Имя-файла Имя-файла
```

### Извлечение содержимого архива

#### из bw.bash

```
_mkDir -t tgz && _getBwTar bw.bash | tar xf - -C tgz && _getBwTar bw.bash tests | tar xf - -C tgz
```

#### из old.bw.bash

```
_mkDir -t tgz && _getBwTar old.bw.bash | tar xf - -C tgz && _getBwTar old.bw.bash tests | tar xf - -C tgz
```

## Тестирование

```
bwdev test
bwdev test _rm
bwdev test _rm 0..1
bwdev test _rm 0..1 _mvFile 1
```

### Без прегенерации

```
bwdev test -n 
bwdev test -n _rm
bwdev test -n _rm 0..1
bwdev test -n _rm 0..1 _mvFile 1
```

### Цикл разработки функции

```
bwdev -n test Имя-Функции
```

## Профилирование кода

```
bwdev profile 'bw bt _rm 0'
bwdev profile bw bt _rm 0
```

## Установка для отладки в той же системе

```
curl -O localhost:${_bwdevDockerHttp:-8998}/bw.bash && . bw.bash -u localhost:${_bwdevDockerHttp:-8998} bw bt
```

или

```
curl -O localhost:${_bwdevDockerHttp:-8998}/bw.bash && BW_SELF_UPDATE_SOURCE=localhost:${_bwdevDockerHttp:-8998} . bw.bash bw bt
```

См. также [BW_SELF_UPDATE_SOURCE](#bw_self_update_source)

## Установка для отладки в гостевой системе под Parallels

```
curl -o ~/bw.bash -L 10.211.55.2:${_bwdevDockerHttp:-8998}/bw.bash && . ~/bw.bash -u 10.211.55.2:${_bwdevDockerHttp:-8998} bw bt
```

или

```
curl -o ~/10.211.55.2:${_bwdevDockerHttp:-8998}/bw.bash && BW_SELF_UPDATE_SOURCE=10.211.55.2:${_bwdevDockerHttp:-8998} . ~/bw.bash bw bt
```

См. также [BW_SELF_UPDATE_SOURCE](#bw_self_update_source)

## Проверка корректности удаления bw

```
BW_VERBOSE=true . bw.bash -p - 'bw rm -y && . bw.bash bw bt && bw rm -y'; compgen -v
```

## Отладка самобновления

```
curl -O localhost:${_bwdevDockerHttp:-8998}/bw.bash && . bw.bash -u localhost:${_bwdevDockerHttp:-8998} 'rm -rf .bw; . bw.bash bw bt'
```

или

```
curl -O localhost:${_bwdevDockerHttp:-8998}/bw.bash && BW_SELF_UPDATE_SOURCE=localhost:${_bwdevDockerHttp:-8998} . bw.bash 'rm -rf .bw; . bw.bash bw bt'
```

См. также [BW_SELF_UPDATE_SOURCE](#bw_self_update_source)

## Переменные среды и опции

### BW_PREGEN_ONLY

`BW_PREGEN_ONLY` (опция `-p`) ограничивает прегенерацию вспомогательного кода (codeToParseOptions, codeToParseFuncParams, autoHelp, completion) и приводит к тому, что прегенерация происходит только для заданных функций (например, `BW_PREGEN_ONLY=bw\ bw_install` или `-p "bw bw_install"`). Значение `-` отменяет прегенерацию (но оставляет ранее сгенеренные файлы нетронутыми). К генерации вспомогательного кода для всех функций и формированию списка completion-функций заново (только в режиме разработки) приводит только пустое значение `BW_PREGEN_ONLY` (достаточно просто не упоминать `BW_PREGEN_ONLY` при вызове `. bw.bash`)

### BW_SELF_UPDATE_SOURCE

`BW_SELF_UPDATE_SOURCE` (опция `-u`) задает источник самообновления. По умолчанию `https://raw.githubusercontent.com/baza-winner/bw/master`.
Значение `-` форсирует использование значения по умолчанию (необходимо, если bw.bash был установлен с иным источником самообновления, и надо переключить уже установленный bw.bash на источник по умолчанию, эта задача решается командой `. bw.bash -u -`)

### BW_VERBOSE

`BW_VERBOSE` задает "говорливый" режим выполнения `bw.bash`. Выводит отладочную информацию

### BW_PROFILE

`BW_PROFILE=true` включает обработку `_profileBegin`/`_profileEnd`. Без этого `_profileBegin`/`_profileEnd` игнорируются
