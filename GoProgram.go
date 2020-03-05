/* CS 424, SP 20
    goprogram.go
    Alexandra Noreiga
------------------------------- */
package main

import ("fmt"
        "bufio"
        "os"
        "strconv"
		"strings"
		"sort"
		"io/ioutil"
        )

/**********************
/*Structure/Object to hold Player data
/**********************/
type Player struct {
    firstName string
    lastName string
    atBats uint64
    plateApp uint64
    singles uint64
    doubles uint64
    triples uint64
    homeRuns uint64
    walks uint64
    hitPitch uint64
}

/**********************
/*Structure/Object to hold formatted Player data
/**********************/
type formattedPlayer struct {
    firstName string 
    lastName string
    average float64
    slug float64
    onBasePercent float64
}	

/**********************
/*Function to sort the Players via first and last name
/**********************/	
func sortPlayers(players []Player) []Player {
    sort.Slice(players, func(i, j int) bool {
        if players[i].lastName != players[j].lastName {
            return players[i].lastName < players[j].lastName
        } else {
            return players[i].firstName < players[j].firstName
        }
    })

    return players
}

/**********************
/*Function to calculate various player statistics
/**********************/
func calcStats(players []Player) []formattedPlayer {
    formatPlayers := make([]formattedPlayer, len(players))

    for i := 0; i < len(players); i++ {
	
		//getting the players first name
        formatPlayers[i].firstName = players[i].firstName
		
		//getting players' last name
        formatPlayers[i].lastName = players[i].lastName
		
		//calculating player averages
        formatPlayers[i].average = float64(players[i].singles + players[i].doubles + players[i].triples + players[i].homeRuns) / float64(players[i].atBats)
        
		//calculating player slugging numbers
		formatPlayers[i].slug = float64(players[i].singles + 2 * players[i].doubles + 3 * players[i].triples + 4 * players[i].homeRuns) / float64(players[i].atBats)
        
		//calculating player on base percentage
		formatPlayers[i].onBasePercent = float64(players[i].singles + players[i].doubles + players[i].triples + players[i].homeRuns + players[i].walks + players[i].hitPitch) / float64(players[i].plateApp)

    }
    return formatPlayers
}

/**********************
/*Function to calculate player averages
/**********************/
func calcAverage(players []formattedPlayer) float64 {
    runTotal := float64(0)

    for i := 0; i < len(players); i++ {
        runTotal += players[i].average
    }
    return runTotal / float64(len(players))
}

/**********************
/*Function to return the path of the input file
/*Implemented outside of main to preserve encapsulation and enforce modularity
/**********************/
func getPath() string {
	//Prompt the user
	fmt.Println("Welcome to the player statistics calculator test program! I am going to\n" +
		"read players from an input data file. You will tell me the name of your\n" +
		"input file. I will store all of the players in a list, compute each player's\n" +
		"averages, and then write the resulting team report to your console\n")
		
	fmt.Print("Please enter the name of your input file: ")
	
	//read in file path
	reader := bufio.NewReader(os.Stdin)
	path, _ := reader.ReadString('\n')
	
	//there was an issue with the reader taking in new lines
	//causing an index out of bounds issue
	//this fixes it so imma just leave it here
	path = path[0 : len(path) - 1]
	
	//return file path
	return path
}

/**********************
/*Function to read in input file
/*Implemented outside of main to preserve encapsulation and enforce modularity
/**********************/
	func readLines(path string) (string, error) {
	file, err := os.Open(path)
	
	if err != nil {
		fmt.Println("I was unable to open the file requested.")
		fmt.Println(err)
		return "", err
	} else {
		fmt.Println("\nSuccessfully found file to open!")
		fmt.Println("Now opening requested file...")
	}
	
	filedata, err := ioutil.ReadAll(file)
	
	file.Close()
	
	return string(filedata), nil
}

/**********************
/*Function to parse the input file line by line 
/*loops through each line in the data file
/**********************/	
func parseLines(data string) ([]Player, []string) {
    var err error    

    players := []Player { }
    invalidData := []string { }

	
    data = strings.Replace(data, "\r", "", -1)
    lines := strings.Split(data, "\n")

 //loop to parse file line by line
 Loop:
    for i := 0; i < len(lines); i++ {
        var player Player

        tokens := []string { } 
        
		//seperate strings by whitespace
        spaceDelimit := strings.Split(lines[i], " ")
        
		for j := 0; j < len(spaceDelimit); j++ {
            if spaceDelimit[j] != "" { 
                tokens = append(tokens, spaceDelimit[j])
            }
    }
	//if a line in the input file has less than 10 tokens, prompt an error message
    if len(tokens) != 10 {
        invalidData = append(invalidData, "line " + strconv.Itoa(i) + ") contains not enough data.")
		continue
	}
		//assign firstName tokens to the 0th (first) place in the tokens array
		player.firstName = tokens[0]
		
		//assign lastName tokens to the 1st (second) place in the tokens array
		player.lastName = tokens[1]
	
	//pointer arithmetic because I'm lazy and don't want to copy and paste the same piece of code ten different times
	//because knowing me I'll spend 3 hours trying to find some dumb copy and paste error in those ten different pieces
	//this is an array that holds each of the numeric values related to the Player 
	batterNumericParts := [...]*uint64 { &player.plateApp, &player.atBats, &player.singles, &player.doubles, &player.triples, &player.homeRuns, &player.walks, &player.hitPitch }
		
		//iterating through previously created array
		for j := 0; j < 8; j++ {
			*batterNumericParts[j], err = strconv.ParseUint(tokens[j + 2], 10, 32)
			if err != nil {
				invalidData = append(invalidData, "line " + strconv.Itoa(i) + ") contains an illegal parameter.")
				continue Loop
			}
		}

		players = append(players, player)

    }

    return players, invalidData

}

/**********************
/*Function to format final player report
/**********************/
func formatReport(players []formattedPlayer, errors []string) {
    fmt.Printf("\nBASEBALL TEAM REPORT --- %d PLAYERS FOUND IN FILE\n", len(players))
    fmt.Printf("OVERALL BATTING AVERAGE is %0.3f\n\n", calcAverage(players))

    fmt.Println("   PLAYER NAME     :    AVERAGE    SLUGGING     ONBASE%")
    fmt.Println("------------------------------------------------------------")

    for i := 0; i < len(players); i++ {
        fmt.Printf("%20v :       %0.3f       %0.3f       %0.3f\n", players[i].lastName + ", " + players[i].firstName, players[i].average, players[i].slug, players[i].onBasePercent)
		}
		
        fmt.Printf("\n----- %d ERROR LINES FOUND IN INPUT FILE ----- \n\n", len(errors))

        for i := 0; i < len(errors); i++ {
            fmt.Println(errors[i])
        }
}

/**********************
/*Main() driver function
/**********************/
func main() {

	file := getPath()
	
	data, err := readLines(file)
	
	if err != nil {
		return 
	}
 
	
    players, badlines := parseLines(data)
    players = sortPlayers(players)

    calcData := calcStats(players)

    formatReport(calcData, badlines)

    fmt.Println("\n End Program - goodbye!")
}
