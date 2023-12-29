package main

import "fmt"

const dataLimit int = 20

// club adalah struktur data untuk merepresentasikan informasi klub dalam liga sepak bola.
type club struct {
	name  string // Nama klub
	win   int    // Jumlah win
	lose  int    // Jumlah lose
	draw  int    // Jumlah draw
	point int    // Jumlah poin (3 untuk kemenangan, 1 untuk seri, 0 untuk kekalahan)
	gf    int    // Jumlah gol yang dicetak (goal for)
	ga    int    // Jumlah gol yang diterima (goal against)
}

// match adalah struktur data untuk merepresentasikan informasi pertandingan antara dua klub.
type match struct {
	home      club // Klub tuan rumah
	away      club // Klub tamu
	homeScore int  // Skor klub tuan rumah
	awayScore int  // Skor klub tamu
}

// matchWeek adalah struktur data untuk merepresentasikan satu pekan pertandingan dalam liga.
type matchWeek struct {
	matches [dataLimit / 2]match
}

// league adalah struktur data untuk merepresentasikan seluruh liga sepak bola.
type league struct {
	schedule [dataLimit*2 - 2]matchWeek // Jadwal pertandingan untuk seluruh musim
	clubs    [dataLimit]club            // Informasi klub dalam liga
	nClub    int                        // Jumlah klub dalam liga
}

func main() {
	var premierLeague league
	var input int = 99
	for input != 0 {
		fmt.Println("---------------------------------------")
		printClubNames(premierLeague)
		fmt.Println("---------------------------------------")
		fmt.Println("1. Add club")
		fmt.Println("2. Delete club")
		fmt.Println("3. Edit club")
		fmt.Println("4. Start season!")
		fmt.Println("0. Exit")
		fmt.Println("---------------------------------------")
		fmt.Print("Input : ")
		fmt.Scan(&input)
		if input == 1 {
			processInputClub(&premierLeague)
		} else if input == 2 {
			processDeleteClub(&premierLeague)
		} else if input == 3 {
			processEditClub(&premierLeague)
		} else if input == 4 {
			seasonMenu(&premierLeague)
		} else {
			fmt.Println("Unknown input")
		}
	}
}

// func seasonMenu digunakan untuk menampilkan menu utama selama musim berlangsung.
// Menu ini memungkinkan pengguna untuk mengisi skor pertandingan, melihat papan klasemen, mengedit skor pekan sebelumnya,
// melanjutkan ke pekan berikutnya, atau mengakhiri musim.
func seasonMenu(leagues *league) {
	var input int = 99
	var gameWeek = 0
	fmt.Println("---------------------------------------")
	fmt.Println("Generating double round robin for match schedule . . .")
	generateSchedule(leagues)
	fmt.Println("Schedule generated")
	printAllGameWeek(*leagues)
	sortByName(leagues)
	for input != 0 {
		fmt.Printf("Game week %d\n", gameWeek+1)
		printGameWeek(*leagues, gameWeek)
		fmt.Println("1. Fill score")
		fmt.Println("2. View current leaderboard")
		fmt.Println("3. Edit previous game week")
		fmt.Println("4. Continue next game week")
		fmt.Println("0. End season")
		fmt.Println("---------------------------------------")
		fmt.Print("Input : ")
		fmt.Scan(&input)
		if input == 1 {
			processScoring(leagues, gameWeek)
		} else if input == 2 {
			processViewLeaderboard(*leagues)
		} else if input == 3 {
			processEditScoring(leagues, gameWeek)
		} else if input == 4 {
			gameWeek++
		} else if input == 0 {
			processEndSeason(*leagues)
		} else {
			fmt.Println("Unknown input")
		}
	}
}

// func processInputClub digunakan untuk memproses input dari pengguna
// dan menambahkan klub baru ke dalam struktur data liga jika belum mencapai batas maksimum.
func processInputClub(leagues *league) {
	if leagues.nClub == dataLimit {
		fmt.Println("Tidak dapat menambahkan lebih banyak klub. Batas liga telah tercapai.")
	}
	var newClub club
	fmt.Print("Masukkan nama klub baru: ")
	fmt.Scan(&newClub.name)
	insertClub(leagues, newClub)
}

// func insertClub digunakan untuk menyisipkan klub baru ke dalam struktur data liga.
func insertClub(leagues *league, newClub club) {
	leagues.clubs[leagues.nClub] = newClub
	leagues.nClub++
	fmt.Println("Klub berhasil ditambahkan ke dalam liga.")
}

// func processDeleteClub digunakan untuk memproses penghapusan klub dari struktur data liga.
func processDeleteClub(leagues *league) {
	var clubName string
	fmt.Print("Masukkan nama klub yang ingin dihapus: ")
	fmt.Scan(&clubName)
	index := searchClubName(*leagues, clubName)
	if index != -1 {
		deleteClub(leagues, index)
		fmt.Println("Klub berhasil dihapus dari liga.")
	} else {
		fmt.Printf("Klub dengan nama %s tidak ditemukan\n", clubName)
	}
}

// func deleteClub digunakan untuk menghapus klub dari struktur data liga.
func deleteClub(leagues *league, x int) {
	for i := x; i < leagues.nClub-1; i++ {
		leagues.clubs[i] = leagues.clubs[i+1]
	}
	leagues.nClub--
}

// func processEditClub digunakan untuk memproses pengeditan nama klub dalam struktur data liga.
func processEditClub(leagues *league) {
	var clubName string
	fmt.Print("Masukkan nama klub yang ingin diedit: ")
	fmt.Scan(&clubName)

	index := searchClubName(*leagues, clubName)
	if index != -1 {
		var newClubName string
		fmt.Print("Masukkan nama baru untuk klub: ")
		fmt.Scan(&newClubName)

		editClubName(leagues, index, newClubName)
		fmt.Println("Nama klub berhasil diperbarui.")
	} else {
		fmt.Printf("Klub dengan nama %s tidak ditemukan\n", clubName)
	}
}

// func editClubName digunakan untuk mengubah nama klub dalam struktur data liga.
func editClubName(leagues *league, index int, newName string) {
	leagues.clubs[index].name = newName
}

// func searchClubName digunakan untuk mencari indeks klub berdasarkan nama menggunakan algoritma binary search.
func searchClubName(leagues league, clubName string) int { //binary search ada disini
	low, high := 0, leagues.nClub-1

	for low <= high {
		mid := (low + high) / 2

		if leagues.clubs[mid].name == clubName {
			return mid
		} else if leagues.clubs[mid].name < clubName {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1
}

// func sortByName digunakan untuk mengurutkan klub dalam struktur data liga berdasarkan nama menggunakan algoritma insertion sort.
func sortByName(league *league) {
	n := league.nClub

	for i := 1; i < n; i++ {
		key := league.clubs[i]
		j := i - 1
		for j >= 0 && league.clubs[j].name > key.name {
			league.clubs[j+1] = league.clubs[j]
			j--
		}
		league.clubs[j+1] = key
	}
}

// func sortByPoints digunakan untuk mengurutkan klub dalam struktur data liga berdasarkan poin menggunakan algoritma selection sort.
func sortByPoints(league *league) {
	n := league.nClub
	for i := 0; i < n-1; i++ {
		maxIndex := i
		for j := i + 1; j < n; j++ {
			if league.clubs[j].point > league.clubs[maxIndex].point {
				maxIndex = j
			}
		}
		league.clubs[i], league.clubs[maxIndex] = league.clubs[maxIndex], league.clubs[i]
	}
}

// func generateSchedule digunakan untuk menghasilkan jadwal pertandingan untuk liga.
// Jadwal pertandingan dibuat untuk setengah musim pertama dan setengah musim kedua.
// Fungsi ini memeriksa apakah jumlah klub dalam liga genap, jika tidak, mencetak pesan kesalahan.
func generateSchedule(leagues *league) {
	if leagues.nClub%2 != 0 {
		fmt.Println("Club count must be even!")
		return
	}

	// Menghasilkan jadwal pertandingan untuk setengah musim pertama
	for i := 0; i < leagues.nClub-1; i++ {
		for j := 0; j < leagues.nClub/2; j++ {
			leagues.schedule[i].matches[j].home = leagues.clubs[j]
			leagues.schedule[i].matches[j].away = leagues.clubs[(leagues.nClub-1)-j]
		}
		rotateRobin(leagues)
	}

	rotateRobin(leagues)

	// Menghasilkan jadwal pertandingan untuk setengah musim kedua
	for i := 0; i < leagues.nClub-1; i++ {
		for j := 0; j < leagues.nClub/2; j++ {
			leagues.schedule[i+leagues.nClub-1].matches[j].home = leagues.clubs[(leagues.nClub-1)-j]
			leagues.schedule[i+leagues.nClub-1].matches[j].away = leagues.clubs[j]
		}
		rotateRobin(leagues)
	}

	rotateRobin(leagues)
}

// func rotateRobin digunakan untuk melakukan rotasi klub dalam array klub.
// Rotasi ini diperlukan untuk menghasilkan jadwal pertandingan sesuai dengan format round-robin.
func rotateRobin(leagues *league) {
	var temp club = leagues.clubs[1]
	for j := 1; j < leagues.nClub-1; j++ {
		leagues.clubs[j] = leagues.clubs[j+1]
	}
	leagues.clubs[leagues.nClub-1] = temp
}

// func printAllGameWeek digunakan untuk mencetak semua jadwal pertandingan dalam liga.
// Fungsi ini mencetak setiap pertandingan dalam setiap pekan dengan format yang terstruktur.
func printAllGameWeek(leagues league) {
	fmt.Println("---------------------------------------")
	for i := 0; i < (leagues.nClub-1)*2; i++ {
		fmt.Printf("Game week %d : \n", i+1)
		for j := 0; j < leagues.nClub/2; j++ {
			fmt.Printf("%s - %s\n", leagues.schedule[i].matches[j].home.name, leagues.schedule[i].matches[j].away.name)
		}
	}
	fmt.Println("---------------------------------------")
}

// func printGameWeek digunakan untuk mencetak jadwal pertandingan dalam satu pekan tertentu dalam liga.
func printGameWeek(leagues league, gw int) {
	for i := 0; i < leagues.nClub/2; i++ {
		printMatch(leagues, gw, i)
	}
}

// func printMatch digunakan untuk mencetak detail satu pertandingan dalam liga.
func printMatch(leagues league, gw int, x int) {
	fmt.Printf("Match %d : \n", x+1)
	fmt.Printf("%-3s - %-3s\n", leagues.schedule[gw].matches[x].home.name, leagues.schedule[gw].matches[x].away.name)
	fmt.Printf("%-3d - %-3d\n", leagues.schedule[gw].matches[x].homeScore, leagues.schedule[gw].matches[x].awayScore)
}

// func printLeaderboard digunakan untuk mencetak papan klasemen dalam liga.
func printLeaderboard(l league) {
	fmt.Printf("%-20s%-5s%-5s%-5s%-5s%-5s%-5s%-5s\n", "Club", "Win", "Draw", "Lose", "GF", "GA", "GD", "Pts")
	for i := 0; i < l.nClub; i++ {
		c := l.clubs[i]
		fmt.Printf("%-20s%-5d%-5d%-5d%-5d%-5d%-5d%-5d\n", c.name, c.win, c.draw, c.lose, c.gf, c.ga, c.gf-c.ga, c.point)
	}
}

// func processViewLeaderboard digunakan untuk menampilkan dan mengelola papan klasemen dalam liga.
// Papan klasemen diurutkan berdasarkan poin secara default.
// Pengguna dapat memilih untuk mengurutkan berdasarkan nama atau memilih kembali ke menu pekan pertandingan.
func processViewLeaderboard(leagues league) {
	var input int = 99
	sortByPoints(&leagues)
	for input != 0 {
		printLeaderboard(leagues)
		fmt.Println("1. Sort by name")
		fmt.Println("2. Sort by point")
		fmt.Println("0. return")
		fmt.Print("Input : ")
		fmt.Scan(&input)
		if input == 1 {
			sortByName(&leagues)
		} else if input == 2 {
			sortByPoints(&leagues)
		} else if input == 0 {
			fmt.Println("Returning to game week menu")
		} else {
			fmt.Println("Unknown input")
		}
	}
}

// func printClubNames digunakan untuk mencetak daftar nama klub dalam liga.
func printClubNames(l league) {
	fmt.Println("Club List:")
	for i := 0; i < l.nClub; i++ {
		fmt.Printf("%d. %s\n", i+1, l.clubs[i].name)
	}
}

// func processScoring digunakan untuk memproses input skor pertandingan untuk satu pekan tertentu dalam liga.
// Skor tersebut kemudian digunakan untuk mengupdate statistik klub dan papan klasemen.
func processScoring(leagues *league, gw int) {
	var goalHome, goalAway int
	for i := 0; i < leagues.nClub/2; i++ {
		printMatch(*leagues, gw, i)
		fmt.Printf("%s score : ", leagues.schedule[gw].matches[i].home.name)
		fmt.Scan(&goalHome)
		fmt.Printf("%s score : ", leagues.schedule[gw].matches[i].away.name)
		fmt.Scan(&goalAway)

		homeIndex := searchClubName(*leagues, leagues.schedule[gw].matches[i].home.name)
		awayIndex := searchClubName(*leagues, leagues.schedule[gw].matches[i].away.name)

		leagues.clubs[homeIndex].gf += goalHome
		leagues.clubs[homeIndex].ga += goalAway
		leagues.clubs[awayIndex].gf += goalAway
		leagues.clubs[awayIndex].ga += goalHome

		if goalHome > goalAway {
			leagues.clubs[homeIndex].win++
			leagues.clubs[awayIndex].lose++
			leagues.clubs[homeIndex].point += 3
		} else if goalAway > goalHome {
			leagues.clubs[awayIndex].win++
			leagues.clubs[homeIndex].lose++
			leagues.clubs[awayIndex].point += 3
		} else {
			leagues.clubs[awayIndex].draw++
			leagues.clubs[homeIndex].draw++
			leagues.clubs[awayIndex].point += 1
			leagues.clubs[homeIndex].point += 1
		}

		leagues.schedule[gw].matches[i].homeScore += goalHome
		leagues.schedule[gw].matches[i].awayScore += goalAway
	}
}

// func processEditScoring digunakan untuk memproses pengeditan skor pertandingan pada suatu pekan tertentu dalam liga.
func processEditScoring(league *league, cgw int) {
	var input int
	for i := 0; i < cgw; i++ {
		fmt.Printf("%d. Game week %d\n", i+1, i+1)
	}
	fmt.Printf("Choose game week (1-%d) : ", cgw+1)
	fmt.Scan(&input)

	if input < 1 || input > cgw+1 {
		fmt.Println("Invalid game week")
		return
	}

	minusScoring(league, input-1)
	processScoring(league, input-1)
}

// func minusScoring digunakan untuk mengurangkan statistik klub dan papan klasemen setelah pengeditan skor pertandingan.
func minusScoring(leagues *league, gw int) {
	var goalHome, goalAway int
	for i := 0; i < leagues.nClub/2; i++ {
		goalHome = leagues.schedule[gw].matches[i].homeScore
		goalAway = leagues.schedule[gw].matches[i].awayScore

		homeIndex := searchClubName(*leagues, leagues.schedule[gw].matches[i].home.name)
		awayIndex := searchClubName(*leagues, leagues.schedule[gw].matches[i].away.name)
		leagues.clubs[homeIndex].gf -= goalHome
		leagues.clubs[homeIndex].ga -= goalAway
		leagues.clubs[awayIndex].gf -= goalAway
		leagues.clubs[awayIndex].ga -= goalHome

		if goalHome > goalAway {
			leagues.clubs[homeIndex].win--
			leagues.clubs[awayIndex].lose--
			leagues.clubs[homeIndex].point -= 3
		} else if goalAway > goalHome {
			leagues.clubs[awayIndex].win--
			leagues.clubs[homeIndex].lose--
			leagues.clubs[awayIndex].point -= 3
		} else {
			leagues.clubs[awayIndex].draw--
			leagues.clubs[homeIndex].draw--
			leagues.clubs[awayIndex].point -= 1
			leagues.clubs[homeIndex].point -= 1
		}

		leagues.schedule[gw].matches[i].homeScore -= goalHome
		leagues.schedule[gw].matches[i].awayScore -= goalAway
	}
}

// func processEndSeason digunakan untuk menampilkan papan klasemen akhir dan klub juara setelah musim berakhir.
func processEndSeason(leagues league) {
	sortByPoints(&leagues)
	printLeaderboard(leagues)
	fmt.Printf("%s WINS THE LEAGUE\n", leagues.clubs[0].name)
}
