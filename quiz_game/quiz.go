package main

import(
	"encoding/csv"
	"fmt"
	"flag"
	"os"
	"strings" 
	"time"

)
type problem struct {
	q string
	a string
}

func main(){

	//Specifying the csv file to look for, the first part is the file type, the second part is the probable name, third part is a desc
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	//Creating a time limit
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz")

	flag.Parse()
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))	
		
		}
	
	//Reading in the csv
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil{
		exit("Failed to parse the provided csv file")
	}
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	
	var counter int
	counter = 0
	for i, p :=range problems {
		fmt.Printf("Problem #%d: %s = \n",i+1, p.q)

		ansCh :=make(chan string)
		//Go routine to scan for answer in background
		go func(){
			var answer string
			fmt.Scanf("%s\n",&answer)
			ansCh <- answer
		}()
		
		select{
		case <-timer.C:
			fmt.Printf("You scored %d out of %d \n", counter, len(problems))
			return

		//If something is in that channel it will run
		case answer := <-ansCh: 
			if answer == p.a {
				counter++
			}	
		}
		
	}
	
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}


//Converting the slice that is brought in into a struct so that when data entry changes a majority of the program stays the same
func parseLines(lines [][]string)[]problem{
	ret :=make([]problem, len(lines))
	for i, line := range lines{
		ret[i] = problem{
		q: line[0],
		a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}
