package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/nickyrolly/ichat/internal/database"
	"github.com/nickyrolly/ichat/internal/repository"
	"github.com/spf13/viper"
)

// TribunData merepresentasikan "table1" dengan ID unik.
type TribunData struct {
	ID   uuid.UUID
	Nama string
}

// KursiData merepresentasikan "table2" dengan ID unik dan TribunID sebagai foreign key.
type KursiData struct {
	ID         uuid.UUID
	Baris      int
	NomorKursi int
	TribunID   uuid.UUID
}

// func main() {
// 	// Nama file CSV yang akan dibaca
// 	// Perbarui nama file sesuai dengan yang Anda unggah, yaitu "tribun_mapping_real.csv"
// 	filename := "cmd/tribun_mapping_real.csv"

// 	// Membaca data dari file CSV
// 	records, err := bacaFileCSV(filename)
// 	if err != nil {
// 		log.Fatalf("Gagal membaca file CSV: %v", err)
// 	}

// 	// Mengolah data CSV menjadi dua grup data
// 	table1, table2 := prosesDataCSV(records)

// 	fmt.Printf("Berhasil memuat %d tribun dan %d kursi dari CSV.\n", len(table1), len(table2))

// 	// Membuat endpoint API
// 	http.HandleFunc("/api/seats", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")

// 		response := map[string]interface{}{
// 			"total_tribun": len(table1),
// 			"total_kursi":  len(table2),
// 			"status":       "success",
// 		}

// 		json.NewEncoder(w).Encode(response)
// 	})

// 	// Menjalankan HTTP server di port 8081 (karena 8080 mungkin sudah dipakai ws-chat)
// 	port := ":8081"
// 	fmt.Printf("Server API berjalan di http://localhost%s\n", port)

// 	if err := http.ListenAndServe(port, nil); err != nil {
// 		log.Fatalf("Gagal menjalankan server: %v", err)
// 	}

// 	// // Mencetak hasil untuk table1 (Daftar Tribun)
// 	// fmt.Println("Grup 1: List Data Tribun (table1)")
// 	// fmt.Println("-------------------------------------------------")
// 	// for _, tribun := range table1 {
// 	// 	fmt.Printf("ID: %s, Nama: %s\n", tribun.ID, tribun.Nama)
// 	// }
// 	// fmt.Println("-------------------------------------------------")
// 	// fmt.Printf("Total tribun unik: %d\n\n", len(table1))

// 	// // Mencetak hasil untuk table2 (Daftar Kursi)
// 	// fmt.Println("Grup 2: List Data Kursi (table2)")
// 	// fmt.Println("-------------------------------------------------")
// 	// for _, kursi := range table2 {
// 	// 	fmt.Printf("ID: %s, Baris: %d, Kursi: %d, TribunID: %s\n", kursi.ID, kursi.Baris, kursi.NomorKursi, kursi.TribunID)
// 	// }
// 	// fmt.Println("-------------------------------------------------")
// 	// fmt.Printf("Total kursi: %d\n", len(table2))
// }

func runMigrations() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("DB_HOST", "127.0.0.1")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "ichat")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: .env file not found for migration")
	}

	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	dbname := viper.GetString("DB_NAME")

	// dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_PORT"),
	// 	os.Getenv("DB_NAME"),
	// )

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	log.Printf("Migration URL: postgres://%s:***@%s:%s/%s?sslmode=disable", user, host, port, dbname)

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Printf("Migration setup error: %v", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Migration error: %v", err)
		return
	}

	log.Println("Migrations completed successfully")
}

func main() {
	// Run migrations first
	runMigrations()

	// Connect to database
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	tribunRepo := repository.NewTribunRepository(db)
	kursiRepo := repository.NewKursiRepository(db)

	// Membaca data dari file CSV
	filename := "cmd/tribun_mapping_real.csv"
	records, err := bacaFileCSV(filename)
	if err != nil {
		log.Printf("Warning: Failed to read CSV file: %v", err)
		log.Println("Using database data instead...")
	} else {
		// Mengolah data CSV menjadi dua grup data
		table1, table2 := prosesDataCSV(records)
		fmt.Printf("Berhasil memuat %d tribun dan %d kursi dari CSV nya.\n", len(table1), len(table2))
	}

	// Membuat endpoint API
	http.HandleFunc("/api/seats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get data from database
		tribunCount, err := tribunRepo.GetCount()
		if err != nil {
			http.Error(w, "Failed to get tribun count", http.StatusInternalServerError)
			return
		}

		kursiCount, err := kursiRepo.GetCount()
		if err != nil {
			http.Error(w, "Failed to get kursi count", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"total_tribun": tribunCount,
			"total_kursi":  kursiCount,
			"status":       "success",
			"source":       "database",
		}

		json.NewEncoder(w).Encode(response)
	})

	// Additional endpoints
	http.HandleFunc("/api/tribun", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tribuns, err := tribunRepo.GetAll()
		if err != nil {
			http.Error(w, "Failed to get tribuns", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(tribuns)
	})

	http.HandleFunc("/api/kursi", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		kursis, err := kursiRepo.GetAll()
		if err != nil {
			http.Error(w, "Failed to get kursis", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(kursis)
	})

	// Menjalankan HTTP server di port 8081
	port := ":8081"
	fmt.Printf("Server API berjalan di http://localhost%s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

func prosesDataCSV(records [][]string) ([]TribunData, []KursiData) {
	var table1 []TribunData
	var table2 []KursiData
	tribunMap := make(map[string]uuid.UUID)

	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}

		tribunNama := strings.TrimSpace(record[0])
		baris, _ := strconv.Atoi(strings.TrimSpace(record[1]))
		nomorKursi, _ := strconv.Atoi(strings.TrimSpace(record[2]))

		// Cek apakah tribun sudah ada di map
		tribunID, exists := tribunMap[tribunNama]
		if !exists {
			tribunID = uuid.New()
			tribunMap[tribunNama] = tribunID
			table1 = append(table1, TribunData{
				ID:   tribunID,
				Nama: tribunNama,
			})
		}

		// Tambahkan kursi
		table2 = append(table2, KursiData{
			ID:         uuid.New(),
			Baris:      baris,
			NomorKursi: nomorKursi,
			TribunID:   tribunID,
		})
	}

	return table1, table2
}

// bacaFileCSV adalah fungsi untuk membaca data dari file CSV.
func bacaFileCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// parseRentangKursi mem-parsing string "1-10" menjadi dua integer,
// yaitu nomor kursi awal dan akhir.
func parseRentangKursi(rentang string) (int, int, error) {
	// Menghapus spasi dan tanda '-' yang tidak perlu
	rentang = strings.TrimSpace(rentang)
	if rentang == "" || rentang == "-" {
		return 0, 0, fmt.Errorf("rentang kursi kosong atau tidak valid")
	}

	parts := strings.Split(rentang, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("format rentang kursi tidak valid: %s", rentang)
	}

	mulai, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("gagal mengkonversi nomor kursi awal: %v", err)
	}

	akhir, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("gagal mengkonversi nomor kursi akhir: %v", err)
	}

	return mulai, akhir, nil
}

// prosesDataCSV mengolah data dari CSV menjadi dua slice baru: TribunData dan KursiData.
// func prosesDataCSV(records [][]string) ([]TribunData, []KursiData) {
// 	var table1 []TribunData
// 	var table2 []KursiData

// 	tribunMap := make(map[string]uuid.UUID)

// 	if len(records) <= 1 {
// 		return table1, table2
// 	}

// 	for _, record := range records[1:] {
// 		// Asumsi format CSV:
// 		// Kolom 0: Nama Tribun
// 		// Kolom 1: Nomor Baris
// 		// Kolom 2: Rentang Nomor Kursi
// 		if len(record) < 3 {
// 			continue
// 		}

// 		tribunNama := record[0]
// 		barisNomorStr := record[1]
// 		nomorKursiStr := record[2]

// 		// Mengabaikan baris yang tidak memiliki informasi yang valid
// 		if tribunNama == "" || barisNomorStr == "" || nomorKursiStr == "" || nomorKursiStr == "-" {
// 			continue
// 		}

// 		// Mendapatkan atau membuat TribunID
// 		var tribunID uuid.UUID
// 		if id, ok := tribunMap[tribunNama]; !ok {
// 			newID := uuid.New()
// 			tribunMap[tribunNama] = newID
// 			table1 = append(table1, TribunData{ID: newID, Nama: tribunNama})
// 			tribunID = newID
// 		} else {
// 			tribunID = id
// 		}

// 		barisNomor, err := strconv.Atoi(barisNomorStr)
// 		if err != nil {
// 			log.Printf("Gagal mengkonversi nomor baris '%s': %v\n", barisNomorStr, err)
// 			continue
// 		}

// 		mulaiKursi, akhirKursi, err := parseRentangKursi(nomorKursiStr)
// 		if err != nil {
// 			log.Printf("Gagal mem-parsing rentang kursi '%s': %v\n", nomorKursiStr, err)
// 			continue
// 		}

// 		// Membuat setiap kursi dalam rentang dan menambahkannya ke table2
// 		for i := mulaiKursi; i <= akhirKursi; i++ {
// 			kursiID := uuid.New()
// 			table2 = append(table2, KursiData{
// 				ID:         kursiID,
// 				Baris:      barisNomor,
// 				NomorKursi: i,
// 				TribunID:   tribunID,
// 			})
// 		}
// 	}

// 	return table1, table2
// }
