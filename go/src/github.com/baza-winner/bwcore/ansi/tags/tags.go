package tags

import "github.com/baza-winner/bwcore/ansi"

func init() {
	ansi.MustAddTag("ansiErr",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorRed, Bright: true}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansi.MustAddTag("ansiOK",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen, Bright: true}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansi.MustAddTag("ansiWarn",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorYellow, Bright: true}),
	)
	ansi.MustAddTag("ansiFunc",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansi.MustAddTag("ansiVar",
		ansi.SGRCodeOfColor256(ansi.Color256{Code: 201}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansi.MustAddTag("ansiVal",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorCyan, Bright: true}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansi.MustAddTag("ansiType",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorWhite, Bright: true}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansi.MustAddTag("ansiPath",
		ansi.SGRCodeOfColor256(ansi.Color256{Code: 252}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansi.MustAddTag("ansiCmd",
		ansi.SGRCodeOfColor256(ansi.Color256{Code: 253}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
}
