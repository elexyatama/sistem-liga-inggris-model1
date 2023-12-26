package main

import "fmt"

const dataLimit int = 20

type club struct {
	name  string
	win   int
	lose  int
	draw  int
	point int //3 menang, 1 seri, 0 kalah
	gf    int //goal for, gol yang dicetak
	ga    int //goal against, gol yang yang dicetak tim musuh
}

type match struct {
	home      club
	away      club
	homeScore int
	awayScore int
}

type matchWeek struct {
	matches [dataLimit / 2]match
}

type league struct {
	schedule [dataLimit*2 - 2]matchWeek
	clubs    [dataLimit]club
	nClub    int
}

func main() {
	var premierLeague league
	var input int = 99
	for input != 0 {
		sortByName(&premierLeague)
		printClubNames(premierLeague)
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

func seasonMenu(leagues *league) {
	var input int = 99
	var gameWeek = 0
	fmt.Println("Generating double round robin for match schedule . . .")
	generateSchedule(leagues)
	fmt.Println("Schedule generated")
	printAllGameWeek(*leagues)
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

func processInputClub(leagues *league) {
	if leagues.nClub == dataLimit {
		fmt.Println("Cannot insert more clubs. leagues limit reached.")
	}
	var newClub club
	fmt.Print("Enter the name of the new club: ")
	fmt.Scan(&newClub.name)
	insertClub(leagues, newClub)
}

func insertClub(leagues *league, newClub club) {
	leagues.clubs[leagues.nClub] = newClub
	leagues.nClub++
	fmt.Println("Club inserted into the leagues.")
}

func processDeleteClub(leagues *league) {
	var clubName string
	fmt.Print("Enter the name of the desired club to delete: ")
	fmt.Scan(&clubName)
	index := searchClubName(*leagues, clubName)
	if index != -1 {
		deleteClub(leagues, index)
		fmt.Println("Club deleted from the leagues.")
	} else {
		fmt.Printf("Club with name %s not found\n", clubName)
	}
}

func deleteClub(leagues *league, x int) {
	for i := x; i < leagues.nClub-1; i++ {
		leagues.clubs[i] = leagues.clubs[i+1]
	}
	leagues.nClub--
}

func processEditClub(leagues *league) {
	var clubName string
	fmt.Print("Enter the name of the club to edit: ")
	fmt.Scan(&clubName)

	index := searchClubName(*leagues, clubName)
	if index != -1 {
		var newClubName string
		fmt.Print("Enter the new name for the club: ")
		fmt.Scan(&newClubName)

		editClubName(leagues, index, newClubName)
		fmt.Println("Club name updated.")
	} else {
		fmt.Printf("Club with name %s not found\n", clubName)
	}
}

func editClubName(leagues *league, index int, newName string) {
	leagues.clubs[index].name = newName
}

func searchClubName(leagues league, clubName string) int { // binary search disini
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

func sortByName(league *league) { //insertion sort disini
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

func sortByPoints(league *league) { //selection sort disini
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

func generateSchedule(leagues *league) {
	if leagues.nClub%2 != 0 {
		fmt.Println("Club count must be even!")
		return
	}

	//generate week first half season
	for i := 0; i < leagues.nClub-1; i++ {
		for j := 0; j < leagues.nClub/2; j++ {
			leagues.schedule[i].matches[j].home = leagues.clubs[j]
			leagues.schedule[i].matches[j].away = leagues.clubs[(leagues.nClub-1)-j]
		}
		rotateRobin(leagues)
	}

	rotateRobin(leagues)

	//generate week second hald season
	for i := 0; i < leagues.nClub-1; i++ {
		for j := 0; j < leagues.nClub/2; j++ {
			leagues.schedule[i+leagues.nClub-1].matches[j].home = leagues.clubs[(leagues.nClub-1)-j]
			leagues.schedule[i+leagues.nClub-1].matches[j].away = leagues.clubs[j]
		}
		rotateRobin(leagues)
	}

	rotateRobin(leagues)
}

func rotateRobin(leagues *league) {
	var temp club = leagues.clubs[1]
	for j := 1; j < leagues.nClub-1; j++ {
		leagues.clubs[j] = leagues.clubs[j+1]
	}
	leagues.clubs[leagues.nClub-1] = temp
}

func printAllGameWeek(leagues league) {
	for i := 0; i < (leagues.nClub-1)*2; i++ {
		fmt.Printf("Game week %d : ", i+1)
		for j := 0; j < leagues.nClub/2; j++ {
			fmt.Printf("%s - %s\n", leagues.schedule[i].matches[j].home.name, leagues.schedule[i].matches[j].away.name)
		}
	}
}

func printGameWeek(leagues league, gw int) {
	for i := 0; i < leagues.nClub/2; i++ {
		printMatch(leagues, gw, i)
	}
}

func printMatch(leagues league, gw int, x int) {
	fmt.Printf("Match %d\n : ", x+1)
	fmt.Printf("%-5s - %-5s\n", leagues.schedule[gw].matches[x].home.name, leagues.schedule[gw].matches[x].away.name)
	fmt.Printf("%-5d - %-5d\n", leagues.schedule[gw].matches[x].homeScore, leagues.schedule[gw].matches[x].awayScore)
}

func printLeaderboard(l league) {
	fmt.Printf("%-20s%-5s%-5s%-5s%-5s%-5s%-5s%-5s\n", "Club", "Win", "Draw", "Lose", "Pts", "GF", "GA", "GD")
	for i := 0; i < l.nClub; i++ {
		c := l.clubs[i]
		fmt.Printf("%-20s%-5d%-5d%-5d%-5d%-5d%-5d%-5d\n", c.name, c.win, c.draw, c.lose, c.point, c.gf, c.ga, c.gf-c.ga)
	}
}

func processViewLeaderboard(leagues league) {
	var input int = 99
	sortByPoints(&leagues)
	for input != 0 {
		printLeaderboard(leagues)
		fmt.Println("1. Sort by name")
		fmt.Println("2. Sort by point")
		fmt.Println("3. return")
		fmt.Print("Input : ")
		fmt.Scan(&input)
		if input == 1 {
			sortByName(&leagues)
		} else if input == 2 {
			sortByPoints(&leagues)
		} else if input == 3 {
			fmt.Println("Returning to game week menu")
		} else {
			fmt.Println("Unknown input")
		}
	}
}

func printClubNames(l league) {
	fmt.Println("Club List:")
	for i := 0; i < l.nClub; i++ {
		fmt.Printf("%d. %s\n", i+1, l.clubs[i].name)
	}
}

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

func processEndSeason(leagues league) {
	sortByPoints(&leagues)
	printLeaderboard(leagues)
	fmt.Printf("%s WINS THE LEAGUE\n", leagues.clubs[0].name)
}
