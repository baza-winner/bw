# bw

Общая bash-инфраструктура [проектов baza-winner](https://github.com/baza-winner)

## Установка

### macOS и Ubuntu

  curl -o ~/bw.bash -L sob.ru/bw && . ~/bw.bash

или

	curl -o ~/bw.bash -L goo.gl/uh4fBy && . ~/bw.bash

или

  curl -o ~/bw.bash -L https://raw.githubusercontent.com/baza-winner/bw/master/bw.bash && . ~/bw.bash

[sob.ru/bw](https://sob.ru/bw) и (goo.gl/uh4fBy)[https://goo.gl/uh4fBy] - это редиректы на (https://raw.githubusercontent.com/baza-winner/bw/master/bw.bash)[https://raw.githubusercontent.com/baza-winner/bw/master/bw.bash]

### Для Windows 10

То же, что для macOS и Ubuntu, но требуется предварительно [включить подсистему Linux](https://docs.microsoft.com/en-us/windows/wsl/install-win10)

## Использование

TODO

## Разработка

### Тестирование и cборка

  . bw.bash -p - bw rm -y && . bw.bash bw bt

или

  BW_PREGEN_ONLY=- . bw.bash bw rm -y && . bw.bash bw bt

См. также [BW_PREGEN_ONLY](#bw_pregen_only)


#### Полное тестирование

  BW_TEST_ALL=true . bw.bash bw bt

Некоторые тесты (такие как тесты для `_spinner`) при обычном тестировании не запускаются из-за частого ложно отрицательного срабатывания. Чтобы их запустить, используйте `BW_TEST_ALL=true`

#### Больше отладочной информации

  BW_VERBOSE=true . bw.bash -p - bw rm -y; BW_VERBOSE=true  . bw.bash bw bt

См. также [BW_VERBOSE](#bw_verbose)

#### Профилирование кода

  . bw.bash -p _profileInit && BW_PROFILE=true bw bt; _profileResult

  . bw.bash -p _profileInit && BW_PROFILE=true _funcToProfile; _profileResult

См. также [BW_PROFILE](#bw_profile)
ВНИМАНИЕ! Для работы под MacOS профайлер требует установки coreutils (`brew install coreutils`): основан на вызове `$_gdate +%s%3N`

#### Извлечение содержимого архива

	_mkDir -t tgz && _getBwTar old.bw.bash | tar xf - -C tgz && _getBwTar old.bw.bash tests | tar xf - -C tgz

#### Просмотр изменений после сборки

	diff tgz/Имя-файла Имя-файла

### Цикл разработки функции

  . bw.bash -p - bw bt Имя-Функции

или

  BW_PREGEN_ONLY=- . bw.bash bw bt Имя-Функции

Пример:

  . bw.bash -p - bw bt _prepareCodeToParseFuncParams2 _parseFuncParams2 _prepareCodeOfAutoHelp2

См. также [BW_PREGEN_ONLY](#bw_pregen_only)

### Сборка для отладки (минуя тесты)

  . bw.bash -p - _buildBw

или

  BW_PREGEN_ONLY=- . bw.bash _buildBw

См. также [BW_PREGEN_ONLY](#bw_pregen_only)

### Установка для отладки в той же системе

  curl -O localhost:8082/bw.bash && . bw.bash -u localhost:8082 bw bt

или

  curl -O localhost:8082/bw.bash && BW_SELF_UPDATE_SOURCE=localhost:8082 . bw.bash bw bt

См. также [BW_SELF_UPDATE_SOURCE](#bw_self_update_source)

#### Проверка корректности удаления bw

	 BW_VERBOSE=true . bw.bash -p - 'bw rm -y && . bw.bash bw bt && bw rm -y'; compgen -v

### Отладка самобновления

  curl -O localhost:8082/bw.bash && . bw.bash -u localhost:8082 'rm -rf .bw; . bw.bash bw bt'

или

  curl -O localhost:8082/bw.bash && BW_SELF_UPDATE_SOURCE=localhost:8082 . bw.bash 'rm -rf .bw; . bw.bash bw bt'

См. также [BW_SELF_UPDATE_SOURCE](#bw_self_update_source)

### Установка для отладки в гостевой системе под Parallels

  curl -o ~/bw.bash -L 10.211.55.2:8082/bw.bash && . ~/bw.bash -u 10.211.55.2:8082 bw bt

или

  curl -o ~/10.211.55.2:8082/bw.bash && BW_SELF_UPDATE_SOURCE=10.211.55.2:8082 . ~/bw.bash bw bt

См. также [BW_SELF_UPDATE_SOURCE](#bw_self_update_source)

### Переменные среды и опции

#### BW_PREGEN_ONLY

`BW_PREGEN_ONLY` (опция `-p`) ограничивает прегенерацию вспомогательного кода (codeToParseOptions, codeToParseFuncParams, autoHelp, completion) и приводит к тому, что прегенерация происходит только для заданных функций (например, `BW_PREGEN_ONLY=bw\ bw_install` или `-p "bw bw_install"`). Значение `-` отменяет прегенерацию (но оставляет ранее сгенеренные файлы нетронутыми). К генерации вспомогательного кода для всех функций и формированию списка completion-функций заново (только в режиме разработки) приводит только пустое значение `BW_PREGEN_ONLY` (достаточно просто не упоминать `BW_PREGEN_ONLY` при вызове `. bw.bash`)

#### BW_SELF_UPDATE_SOURCE

`BW_SELF_UPDATE_SOURCE` (опция `-u`) задает источник самообновления. По умолчанию `https://raw.githubusercontent.com/baza-winner/bw/master`.
Значение `-` форсирует использование значения по умолчанию (необходимо, если bw.bash был установлен с иным источником самообновления, и надо переключить уже установленный bw.bash на источник по умолчанию, эта задача решается командой `. bw.bash -u -`)

#### BW_VERBOSE

`BW_VERBOSE` задает "говорливый" режим выполнения bw.bash. Выводит отладочную информацию

#### BW_PROFILE

`BW_PROFILE=true` включает отработку _profileBegin/_profileEnd. Без этого _profileBegin/_profileEnd игнорируются
