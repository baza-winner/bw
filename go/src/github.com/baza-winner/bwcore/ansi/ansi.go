// Package ansi предоставляет поддержку "Управляющих последовательностей ANSI".
//
// rus: https://ru.wikipedia.org/wiki/%D0%A3%D0%BF%D1%80%D0%B0%D0%B2%D0%BB%D1%8F%D1%8E%D1%89%D0%B8%D0%B5_%D0%BF%D0%BE%D1%81%D0%BB%D0%B5%D0%B4%D0%BE%D0%B2%D0%B0%D1%82%D0%B5%D0%BB%D1%8C%D0%BD%D0%BE%D1%81%D1%82%D0%B8_ANSI
//
// eng: https://en.wikipedia.org/wiki/ANSI_escape_code
package ansi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/baza-winner/bwcore/bw"
)

// NoColor регулирует обработку <ansi*>разметки функцией String
// Если NoColor = true, то просто убирает <ansi*>-разметку из результата
var NoColor = false

// ============================================================================

type CSI struct {
	parameterBytes    []byte // in the range 0x30–0x3F (ASCII 0–9:;<=>?)
	intermeidateBytes []byte // in the range 0x20–0x2F (ASCII space and !"#$%&'()*+,-./)
	finalByte         byte   // in the range 0x40–0x7E
}

func CSIFrom(parameterBytes []byte, intermeidateBytes []byte, finalByte byte) (result CSI, err error) {
	min := byte(0x30)
	max := byte(0x3F)
	for i, b := range parameterBytes {
		if !(min <= b && b <= max) {
			err = fmt.Errorf("CSIFrom: parameterBytes[%d] (%x) is out of range [%x, %x", i, b, min, max)
			return
		}
	}
	min = byte(0x20)
	max = byte(0x2F)
	for i, b := range intermeidateBytes {
		if !(min <= b && b <= max) {
			err = fmt.Errorf("CSIFrom: intermeidateBytes[%d] (%x) is out of range [%x, %x", i, b, min, max)
			return
		}
	}
	min = byte(0x40)
	max = byte(0x7E)
	if !(min <= finalByte && finalByte <= max) {
		err = fmt.Errorf("CSIFrom: finalByte (%x) is out of range [%x, %x", finalByte, min, max)
		return
	}
	result = CSI{parameterBytes, intermeidateBytes, finalByte}
	return
}

func mustCSIFrom(parameterBytes []byte, intermeidateBytes []byte, finalByte byte) CSI {
	if result, err := CSIFrom(parameterBytes, intermeidateBytes, finalByte); err != nil {
		panic(err.Error())
	} else {
		return result
	}
}

func CSIFromSGRCodes(v ...SGRCode) CSI {
	return CSIFromSGR(v)
}

func CSIFromSGR(v []SGRCode) CSI {
	uint8s := []uint8{}
	for _, i := range v {
		uint8s = append(uint8s, i.codes...)
	}
	ss := make([]string, 0, len(uint8s))
	for _, i := range uint8s {
		ss = append(ss, strconv.FormatUint(uint64(i), 10))
	}
	return mustCSIFrom(
		[]byte(strings.Join(ss, ";")),
		nil,
		'm',
	)
	// return CSI{parameterBytes: []byte(strings.Join(ss, ";")), intermeidateBytes: nil, finalByte: 'm'}
}

func (v CSI) String() string {
	return "\x1b[" + string(v.parameterBytes) + string(v.intermeidateBytes) + string(v.finalByte)
}

// ============================================================================

// SGRCode - Select Graphic Rendition codes
type SGRCode struct {
	codes []uint8
}

// SGRCmd - Select Graphic Rendition commands
type SGRCmd uint8

// SGRCmd
const (
	SGRCmdReset  SGRCmd = iota // all attributes off
	SGRCmdBold                 // increased intensity
	SGRCmdFaint                // decreased intensity
	SGRCmdItalic               // Not widely supported. Sometimes treated as inverse
	SGRCmdUnderline
	SGRCmdSlowBlink    // less than 150 per minute
	SGRCmdRapidBlink   // MS-DOS ANSI.SYS; 150+ per minute; not widely supported
	SGRCmdReverseVideo // swap foreground and background colors
	SGRCmdConceal      // Not widely supported
	SGRCmdCrossedOut   // Characters legible, but marked for deletion
)

// SGRCmd
const (
	SGRCmdFraktur                  SGRCmd = 20 + iota // Rarely supported
	SGRCmdDoublyUnderlineOrBoldOff                    // Double-underline per ECMA-48.[20] See discussion
	SGRCmdNormalColorOrIntensity                      // Neither bold nor faint
	SGRCmdNotItalicNorFraktur
	SGRCmdUnderlineOff // Not singly or doubly underlined
	SGRCmdBlinkOff
)

// SGRCmd
const (
	SGRCmdInverseOff SGRCmd = 27 + iota
	SGRCmdReveal            // conceal off
	SGRNotCrossedOut
)

// SGRCmd
const (
	SGRCmdDefaultForegroundColor SGRCmd = 39
	SGRCmdDefaultBackgroundColor SGRCmd = 49
)

// SGRCmd
const (
	SGRCmdFramed SGRCmd = 51 + iota
	SGRCmdEncircled
	SGRCmdOverlined
	SGRCmdNotFramedOrEncircled
	SGRCmdNotOverlined
)

// SGRCmd, Rarely supported
const (
	SGRCmdIdeogramUnderline       SGRCmd = 60 + iota // or right side line
	SGRCmdIdeogramDoubleUnderline                    // or double line on the right side
	SGRCmdIdeogramOverline                           // or left side line
	SGRCmdIdeogramDoubleOverline                     // or double line on the left side
	SGRCmdIdeogramStressMarking
	SGRCmdIdeogramReset // reset the effects of all of 60–64
)

// SGRCodeOfCmd returns SGRCode for SGRCmd
func SGRCodeOfCmd(a SGRCmd) (result SGRCode, err error) {
	if SGRCmdReset <= a && a <= SGRCmdCrossedOut ||
		SGRCmdFraktur <= a && a <= SGRCmdBlinkOff ||
		SGRCmdInverseOff <= a && a <= SGRNotCrossedOut ||
		a == SGRCmdDefaultForegroundColor ||
		a == SGRCmdDefaultBackgroundColor ||
		SGRCmdFramed <= a && a <= SGRCmdNotOverlined ||
		SGRCmdIdeogramUnderline <= a && a <= SGRCmdIdeogramReset {
		result = SGRCode{[]uint8{uint8(a)}}
	} else {
		err = fmt.Errorf("SGRCmd %d is out of range %d..%d || %d..%d || %d..%d || %d || %d || %d..%d || %d..%d", a,
			SGRCmdReset, SGRCmdCrossedOut,
			SGRCmdFraktur, SGRCmdBlinkOff,
			SGRCmdInverseOff, SGRNotCrossedOut,
			SGRCmdDefaultForegroundColor,
			SGRCmdDefaultBackgroundColor,
			SGRCmdFramed, SGRCmdNotOverlined,
			SGRCmdIdeogramUnderline, SGRCmdIdeogramReset,
		)
	}
	return
}

func MustSGRCodeOfCmd(a SGRCmd) SGRCode {
	if result, err := SGRCodeOfCmd(a); err == nil {
		return result
	} else {
		panic(err.Error())
	}
}

// SGRCodeOfFont returns SGRCode for SGRFont (0..9)
// 0 - Primary(default) font
// 1..9 - alternative font number
func SGRCodeOfFont(a uint8) (result SGRCode, err error) {
	max := uint8(9)
	if a <= max {
		result = SGRCode{[]uint8{a + 10}}
	} else {
		err = fmt.Errorf("SGRFont %d must be less than %d", a, max)
	}
	return
}

func MustSGRCodeOfFont(a uint8) SGRCode {
	if result, err := SGRCodeOfFont(a); err == nil {
		return result
	} else {
		panic(err.Error())
	}
}

type SGRColor8 uint8

const (
	SGRColorBlack SGRColor8 = iota
	SGRColorRed
	SGRColorGreen
	SGRColorYellow
	SGRColorBlue
	SGRColorMagenta
	SGRColorCyan
	SGRColorWhite
)

// Color8 - arg for SGRCodeOfColor8/MustMustSGRCodeOfColor8
type Color8 struct {
	Color      SGRColor8
	Background bool
	Bright     bool
}

// SGRCodeOfColor8 returns SGRCode for Color8 or error
func SGRCodeOfColor8(a Color8) (result SGRCode, err error) {
	max := SGRColorWhite
	if !(a.Color <= max) {
		err = fmt.Errorf("SGRColor8 %d must be less than %d", a.Color, max)
	}
	var base uint8
	if a.Bright {
		base = 90
	} else {
		base = 30
	}
	if a.Background {
		base += 10
	}
	result = SGRCode{[]uint8{base + uint8(a.Color)}}
	return
}

// MustSGRCodeOfColor8 returns SGRCode for Color8 or panics
func MustSGRCodeOfColor8(a Color8) SGRCode {
	if result, err := SGRCodeOfColor8(a); err == nil {
		return result
	} else {
		panic(err.Error())
	}
}

// ColorRGB - arg for SGRCodeOfColorRGB
type ColorRGB struct {
	Red        uint8
	Green      uint8
	Blue       uint8
	Background bool
}

// SGRCodeOfColorRGB returns SGRCode for ColorRGB or error
func SGRCodeOfColorRGB(a ColorRGB) SGRCode {
	var firstCode uint8
	if a.Background {
		firstCode = 48
	} else {
		firstCode = 38
	}
	return SGRCode{[]uint8{firstCode, 2, a.Red, a.Green, a.Blue}}
}

// Color256 - arg for SGRCodeOfColor256
type Color256 struct {
	Code       uint8
	Background bool
}

// SGRCodeOfColor256 returns SGRCode for Color256 or error
func SGRCodeOfColor256(a Color256) SGRCode {
	var firstCode uint8
	if a.Background {
		firstCode = 48
	} else {
		firstCode = 38
	}
	return SGRCode{[]uint8{firstCode, 5, a.Code}}
}

// ============================================================================

func Tag(name string) (result []SGRCode, ok bool) {
	var tag ansiTag
	if tag, ok = ansiTags[name]; ok {
		result = tag.codes
	}
	return
}

const ansiPrefix = "ansi"

func AddTag(name string, codes ...SGRCode) error {
	minLen := len(ansiPrefix) + 1
	if len(name) < minLen {
		return fmt.Errorf("AddTag: len(%q) < minLen (%d)", name, minLen)
	}
	if name[:4] != ansiPrefix {
		return fmt.Errorf("AddTag: name (%q) must start with %q", name, ansiPrefix)
	}
	if _, ok := Tag(name); ok {
		return fmt.Errorf("<%s> already exists", name)
	}
	ansiTags[name] = ansiTag{codes, CSIFromSGR(codes).String()}
	return nil
}

func MustAddTag(name string, codes ...SGRCode) {
	if err := AddTag(name, codes...); err != nil {
		panic(err.Error())
	}
}

func MustTag(name string) []SGRCode {
	if result, ok := Tag(name); !ok {
		panic(bw.Spew.Sprintf("MustTag: <%s> is unknown", name))
	} else {
		return result
	}
}

type ansiTag struct {
	codes     []SGRCode
	CSIString string
}

var ansiTags map[string]ansiTag
var resetCSIString string
var resetCSIStringLen int

func Reset() string {
	return resetCSIString
}

func init() {
	ansiTags = map[string]ansiTag{}
	// resetCode := MustSGRCodeOfCmd(SGRCmdReset)
	MustAddTag("ansiReset", MustSGRCodeOfCmd(SGRCmdReset))
	resetCSIString = ansiTags["ansiReset"].CSIString
	resetCSIStringLen = len(resetCSIString)
}

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
    String
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
    Primary
    Secondary

Пример:

*/

type A struct {
	Default []SGRCode
	S       string
}

func String(s string, ansiDefault ...SGRCode) string {
	return StringA(A{S: s, Default: ansiDefault})
}

func StringA(a A) (result string) {
	ansiCSIString := resetCSIString
	if len(a.Default) > 0 {
		ansiCSIString = CSIFromSGR(a.Default).String()
	}
	if len(a.S) == 0 {
		return
	}
	var didAnsi bool
	result = findAnsiRegexp.ReplaceAllStringFunc(a.S, func(s string) (result string) {
		if NoColor {
			return
		}
		didAnsi = true
		if name := s[1 : len(s)-1]; name == ansiPrefix {
			result = ansiCSIString
		} else if tag, ok := ansiTags[name]; ok {
			result = tag.CSIString
		} else {
			panic(fmt.Sprintf("ansi.From: <%s> is unknown at %q", name, a.S))
		}
		return
	})
	if !NoColor && ansiCSIString != resetCSIString {
		didAnsi = true
		result = ansiCSIString + result
	}
	if !NoColor && didAnsi && !EndsWithReset(result) {
		result = result + resetCSIString
	}
	return result
}

var findAnsiRegexp = regexp.MustCompile("<ansi[^>]*>")

// func Concat(ss ...string) (result string) {
// 	for _, s := range ss {
// 		result += s
// 		// result = ChopReset(result) + s
// 	}
// 	return
// }

func EndsWithReset(s string) bool {
	sLen := len(s)
	return sLen >= resetCSIStringLen && s[sLen-resetCSIStringLen:] == resetCSIString
}

func ChopReset(s string) (result string) {
	result = s
	if EndsWithReset(result) {
		result = result[:len(result)-resetCSIStringLen]
	}
	return
}
