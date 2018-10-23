// Предоставляет функцию Ansi для обработки строк с <ansi*>-разметкой.
package ansi

import (
	"errors"
	"fmt"
	"regexp"
)

/*
Регулирует обработку <ansi*>разметки функцией Ansi
Если NoColor = true, то просто убирает <ansi*>-разметку из результата
*/
var NoColor = false

const escape = "\x1b"

var ansiMap = map[string]string{
	"Reset":     escape + "[0m",
	"Bold":      escape + "[1m",
	"Dim":       escape + "[2m",
	"Italic":    escape + "[3m",
	"Underline": escape + "[4m",
	"Blink":     escape + "[5m",
	"Invert":    escape + "[7m",
	"Hidden":    escape + "[8m",
	"Strike":    escape + "[9m",

	"ResetBold":      escape + "[22m", // именно так, а не 20
	"ResetDim":       escape + "[22m",
	"ResetItalic":    escape + "[23m",
	"ResetUnderline": escape + "[24m",
	"ResetBlink":     escape + "[25m",
	"ResetInvert":    escape + "[27m",
	"ResetHidden":    escape + "[28m",
	"ResetStrike":    escape + "[29m",

	"Black":        escape + "[30m",
	"Red":          escape + "[31m",
	"Green":        escape + "[32m",
	"DarkGreen":    escape + "[38;5;28m",
	"Yellow":       escape + "[33m",
	"Blue":         escape + "[34m",
	"Magenta":      escape + "[38;5;201m", // "Magenta": escape + "[35m",
	"Cyan":         escape + "[36m",
	"LightGray":    escape + "[37m",
	"LightGrey":    "LightGray",
	"Default":      escape + "[39m",
	"DarkGray":     escape + "[90m",
	"DarkGrey":     "DarkGray",
	"LightRed":     escape + "[91m",
	"LightGreen":   escape + "[92m",
	"LightYellow":  escape + "[93m",
	"LightBlue":    escape + "[94m",
	"LightMagenta": escape + "[95m",
	"LightCyan":    escape + "[96m",
	"White":        escape + "[97m",

	"Header":    "Bold;LightGray",
	"Url":       "Blue;Underline",
	"Cmd":       "White;Bold",
	"FileSpec":  "White;Bold",
	"Dir":       "White;Bold",
	"Err":       "Red;Bold",
	"Warn":      "Yellow;Bold",
	"OK":        "Green;Bold",
	"Outline":   "Magenta;Bold",
	"Debug":     "Blue;ResetBold",
	"Primary":   "Cyan;Bold",
	"Secondary": "Cyan;ResetBold",
}

func init() {
	for name, ansiName := range ansiMap {
		if ansiName[0:1] != escape {
			ansiMap[name] = ansi(ansiName)
		}
	}
}

var splitNamesRegexp, _ = regexp.Compile(`[^\w\d]+`)

func ansi(source string) (ansi string) {
	names := splitNamesRegexp.Split(source, -1)
	for _, name := range names {
		var ansiName, err = getAnsiByName(name)
		if err != nil {
			panic(fmt.Sprintf(`ansi: %s`, err.Error()))
		}
		ansi += ansiName
	}
	return
}

var findAnsiRegexp, _ = regexp.Compile(`<ansi[^>]*>`)
var findEnd, _ = regexp.Compile(`\n*$`)

/*
Обрабатывает строку source с <ansi*>-разметкой:

- предваряя строку source ESC-последовательностью, соответствующей тегу <ansi${defaultAnsiName}>,
если defaultAnsiName не пустая строка, или <ansiReset> в противном случае

- заменяя теги <ansi*> на соответствующие ESC-последовательности, тег <ansi> при этом
заменяется на ESC-последовательность, которой была предварена строка

- вставляя ESC-последовательность, соответствующую тегу <ansiReset>, перед замыкающими строку source
символами перевода строки, если такие есть, или в конец строки в противном случае
Возвращает обработанную строку

Список доступных *-значений <ansi*>-тегов (по категориям):

  Форматирование текста:
    Bold
    Dim
    Italic
    Underline
    Blink
    Invert
    Hidden
    Strike

  Отмена форматирования текста:
    Reset - общий сброс
    ResetBold
    ResetDim
    ResetItalic
    ResetUnderline
    ResetBlink
    ResetInvert
    ResetHidden
    ResetStrike

  Цвет текста:
    Black
    Red
    Green
    Yellow
    Blue
    Magenta
    Cyan
    LightGray
    LightGrey
    Default
    DarkGray
    DarkGrey
    LightRed
    LightGreen
    LightYellow
    LightBlue
    LightMagenta
    LightCyan
    White

  Семантическая разметка:
    Header
    Url
    Cmd
    FileSpec
    Dir
    Err
    Warn
    OK
    Outline
    Debug
    PrimaryLiteral
    SecondaryLiteral

Пример:

*/
func Ansi(defaultAnsiName, source string) (result string) {
	var ansiDefault = ansiMap[`Reset`]
	if len(defaultAnsiName) > 0 {
		var err error
		ansiDefault, err = getAnsiByName(defaultAnsiName)
		if err != nil {
			panic(fmt.Sprintf(`Ansi.defaultAnsiName: %s`, err.Error()))
		}
	}
	result = findAnsiRegexp.ReplaceAllStringFunc(source, func(s string) (ansi string) {
		name := s[5 : len(s)-1]
		if len(name) == 0 {
			ansi = ansiDefault
		} else {
			var err error
			ansi, err = getAnsiByName(name)
			if err != nil {
				panic(fmt.Sprintf(`Ansi.source: %s`, err.Error()))
			}
		}
		if NoColor {
			ansi = ``
		}
		return
	})
	result = findEnd.ReplaceAllStringFunc(result, func(s string) (result string) {
		result = s
		if !NoColor {
			result = ansiMap[`Reset`] + result
		}
		return
	})
	if !NoColor {
		result = ansiDefault + result
	}
	return
}

func getAnsiByName(name string) (ansi string, err error) {
	var ok bool
	ansi, ok = ansiMap[name]
	if !ok {
		err = errors.New(fmt.Sprintf(`name'%s' not found`, name))
	} else if ansi[0:1] != escape {
		err = errors.New(fmt.Sprintf(`name'%s' has no escape sequence`, name))
	}
	return
}
