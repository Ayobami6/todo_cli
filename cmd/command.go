package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Ayobami6/todo_cli/db"
	"github.com/mergestat/timediff"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/rand"
)

func generateUniqueRandomIntegers(n int) string {
	rand.Seed(uint64(time.Now().UnixNano()))

	numbersMap := make(map[int]struct{})
	var numbers []int
	for len(numbers) < n {
		num := rand.Intn(100)
		if _, exists := numbersMap[num]; !exists {
			numbersMap[num] = struct{}{}
			numbers = append(numbers, num)
		}
	}
	numStr := ""
	for _, num := range numbers {
		numStr += strconv.Itoa(num)
	}
	return numStr
}

func configure() {

	byteValue, rErr := os.ReadFile("banner.txt")
	if rErr != nil {
		log.Println(rErr)
	} else {
		fmt.Println(string(byteValue))
	}
	var passcode, generatedPasscode string

	fmt.Println("Press 1 to Enter your passcode")
	fmt.Println("Press 2 to Generate new passcode")
	// get input from the cli
	var input string
	fmt.Print("Enter Choice: ")
	_, err := fmt.Scanln(&input)
	// if input is 1 then prompt for passcode
	if err != nil {
		log.Fatal(err)
	}
	switch input {
	case "1":
		fmt.Print("Enter Passcode: ")
		_, err := fmt.Scanln(&passcode)
		if err != nil {
			log.Fatal(err)
		}
		user, err := db.FetchUser(passcode)
		if err != nil {
			log.Fatal(err)
		}
		viper.Set("passcode", user.Passcode)
		fmt.Println("Configuration successful")
	case "2":
		generatedPasscode = generateUniqueRandomIntegers(4)
		// set passcode to viper
		viper.Set("passcode", generatedPasscode)
		db.SaveUser(generatedPasscode)
		err := viper.WriteConfigAs("config.json") // Writes to a config file
		if err != nil {
			// Error handling if writing fails
			log.Fatalf("Error writing config file: %v", err)
		}
		fmt.Println("Generated Passcode: ", generatedPasscode)
	default:
		log.Fatal("Invalid Choice")
	}

}

var RootCmd = &cobra.Command{
	Use:   "tasks",
	Short: "A CLI To-Do List Application",
}

var AddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		description := args[0]
		// get user passcode from viper
		passcode := viper.GetString("passcode")
		// if not passcode
		if passcode == "" {
			fmt.Println("Please configure the application first")
			return
		}
		db.AddTask(description, passcode)
	},
}
var ListCommand = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement the logic to list all tasks
		tasks, err := db.FindAllUserTasks(viper.GetString("passcode"))
		if err != nil {
			log.Fatal(err)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Task", "Created", "Done"})
		if len(tasks) == 0 {
			log.Fatal("tasks is empty")
		}
		for _, task := range tasks {
			if !task.IsComplete {
				row := []string{
					(task.ID).String(),
					task.Description,
					timediff.TimeDiff(task.CreatedAt),
					strconv.FormatBool(task.IsComplete),
				}
				table.Append(row)
			}
		}
		table.Render()
	},
}
var CompleteCommand = &cobra.Command{
	Use:   "complete",
	Short: "Mark a task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement the logic to mark a task as complete
	},
}

var ConfigureCommand = &cobra.Command{
	Use:   "configure",
	Short: "Configure the application",
	Run: func(cmd *cobra.Command, args []string) {
		configure()
	},
}

func init() {
	RootCmd.AddCommand(AddCommand)
	RootCmd.AddCommand(ListCommand)
	RootCmd.AddCommand(CompleteCommand)
	RootCmd.AddCommand(ConfigureCommand)
}
