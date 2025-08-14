package bot

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func Konek() (*whatsmeow.Client, error) {
	store, err := initStore()
	if err != nil {
		return nil, fmt.Errorf("gagal inisialisasi store : %w", err)
	}

	clientLog := waLog.Stdout("CLIENT", "INFO", true)
	device, err := store.GetFirstDevice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("gagal ngambil device dari store : %w", err)
	}
	wangsaf := whatsmeow.NewClient(device, clientLog)
	if wangsaf.Store.ID == nil {
		err = pairing(context.Background(), wangsaf)
	} else {
		wangsaf.Log.Infof("Sesi ditemukan, Connecting....")
		err = wangsaf.Connect()
	}

	if err != nil {
		return nil, err
	}
	return wangsaf, nil
}

func initStore() (*sqlstore.Container, error) {
	dbLog := waLog.Stdout("Database", "INFO", true)
	dbPath := "session/session.db"
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("gagal bikin direktori untuk database : %w", err)
	}
	dbURI := fmt.Sprintf("file:%s?_foreign_keys=on", dbPath)

	rawDB, err := sql.Open("sqlite3", dbURI)
	if err != nil {
		return nil, fmt.Errorf("gagal buka raw connection sqlite : %w", err)
	}
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA cache_size=-16384;", // ngetest cache size di 16MB <== Gacor kink, tapi keknya berlebihan wkwk :v
		// trade off nya ada di memory usage, sedikit lebih besar dibanding gak pake cache_size
		// kalo lu mau yang memory usage nya kecil yaa ilangin aja atau ubah ke -4096 atau -8192
		"PRAGMA busy_timeout=5000;",
	}
	for _, pragma := range pragmas {
		if _, err = rawDB.Exec(pragma); err != nil {
			return nil, fmt.Errorf("gagal make pragma WAL mode (%s) : %w", pragma, err)
		}
	}
	if err = rawDB.Close(); err != nil {
		return nil, fmt.Errorf("gagal close raw db connection : %w", err)
	}

	container, err := sqlstore.New(context.Background(), "sqlite3", dbURI, dbLog)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat container database : %w", err)
	}
	return container, nil
}
func pairing(ctx context.Context, wa *whatsmeow.Client) error {
	wa.Log.Warnf("Session not found, pairing....")
	fmt.Print("Masukkan nomor hp untuk pairing via code (628xxxx)\natau ketik 'qr' untuk pairing via QR : ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	phone := strings.TrimSpace(input)

	if phone == "qr" {
		return pairViaQR(ctx, wa)
	}
	return pairViaCode(ctx, wa, phone)
}

func pairViaQR(ctx context.Context, wa *whatsmeow.Client) error {
	wa.Log.Infof("Requesting QR Code ...")
	qrChan, err := wa.GetQRChannel(ctx)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan kode qr : %w", err)
	}
	err = wa.Connect()
	if err != nil {
		return fmt.Errorf("gagal connect untuk pairing : %w", err)
	}
	for evt := range qrChan {
		switch evt.Event {
		case "code":
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			fmt.Println("Pindai Kode QR di atas dengan WhatsApp Anda.")

		case "success":
			wa.Log.Infof("Pairing via QR berhasil!")
			return nil

		default:
			wa.Log.Infof("Event Login: %s", evt.Event)
		}
	}
	return nil
}

func pairViaCode(ctx context.Context, wa *whatsmeow.Client, phone string) error {
	wa.Log.Infof("Requesting Pairing code for +%s...", phone)
	err := wa.Connect()
	if err != nil {
		return fmt.Errorf("gagal connect untuk pairing : %w", err)
	}
	code, err := wa.PairPhone(ctx, phone, true, whatsmeow.PairClientFirefox, "Firefox (Linux)")
	if err != nil {
		return fmt.Errorf("gagal request pairing code : %w", err)
	}

	fmt.Println("[ PAIRING CODE ]")
	fmt.Printf("Your code : %s\n", code)

	return nil
}
