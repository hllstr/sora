package commands

/*
	Sistem command nya udah plug & play, jadi gak ribet
	Disini kita pake func init() biar gak ribet daftarin command 1 1,
	Jadi semua command udah otomatis ke register...

	Nama file untuk command/fitur-fitur juga bisa ditandai by category
	(misal : downloader.go) nanti isinya fitur2 downloader,
	biar rapih/terorganisir aja sih, untuk subfolder sebenernya bisa
	tapi gw belum nyoba sih, paling nanti gw akan tes di "hoshi", untuk sekarang
	sora gw biarin simpel dulu, biar yang make nanti yang ngembangin sendiri
*/

func init() {
	Plugin(Cmd{
		Name:  "ping",
		Alias: []string{"p"},
		Desc:  "Check bot response.",
		// inline code
		Exec: func(ctx *CommandContext) {
			ctx.Reply("*Pong!*")
		},
	})

	/*
		Ada 2 cara buat masukin kode untuk fiturnya yang pertama pakai anonymous function (inline)
		Yang kedua itu pake function terpisah, tergantung pemakaian, kalo gw misal kodenya pendek
		ya inline ajaa, tapi kalo kodenya panjang mending terpisah sih biar gak berantakan.
	*/

	Plugin(Cmd{
		Name:  "info",
		Alias: []string{"i"},
		Desc:  "Show Information.",
		Exec:  info,
	})
}

// function terpisah
func info(ctx *CommandContext) {
	ingfoText := "Sora is a simple WhatsApp Bot written in GO using whatsmeow library."
	ctx.Reply(ingfoText)
}
