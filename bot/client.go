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

// kita buwat custom logger biar log nya gak jorok.
type CustomLogger struct {
	oriLogger waLog.Logger
}

func MyCustomLogger(logAseli waLog.Logger) *CustomLogger {
	return &CustomLogger{oriLogger: logAseli}
}

func (cl *CustomLogger) Errorf(msg string, args ...any) {
	log := fmt.Sprintf(msg, args...)
	if strings.Contains(log, "Failed to handle retry receipt ") {
		return
	}
	cl.oriLogger.Errorf(msg, args...)
}

func (cl *CustomLogger) Warnf(msg string, args ...any) {
	log := fmt.Sprintf(msg, args...)
	if strings.Contains(log, "Server returned different participant list hash") {
		return
	}
	cl.oriLogger.Warnf(msg, args...)
}

func (cl *CustomLogger) Infof(msg string, args ...any) {
	cl.oriLogger.Infof(msg, args...)
}

func (cl *CustomLogger) Debugf(msg string, args ...any) {
	cl.oriLogger.Debugf(msg, args...)
}

func (cl *CustomLogger) Sub(prefix string) waLog.Logger {
	return MyCustomLogger(cl.oriLogger.Sub(prefix))
}

// omkeh done.

func Konek() (*whatsmeow.Client, error) {
	store, err := initStore()
	if err != nil {
		return nil, fmt.Errorf("gagal inisialisasi store : %w", err)
	}
	loggerAseli := waLog.Stdout("CLIENT", "INFO", true)
	clientLog := MyCustomLogger(loggerAseli)
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
		"PRAGMA cache_size=-16384;", // setelah gw test ternyata lebih baik pake yang 16M, memory usage nya juga gak terlalu jauh beda sama di set ke 4M ketika idle (5MB~7MB memory usage)
		// -4096 / -8192 / -16384 : 4M / 8M / 16M, sesuaikan aja sesuai kebutuhan
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
