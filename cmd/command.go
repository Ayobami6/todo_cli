package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Ayobami6/todo_cli/db"
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

	banner := `	
                                                                                                                                      
                                                 dddddddd                                                                             
TTTTTTTTTTTTTTTTTTTTTTT                          d::::::d                              CCCCCCCCCCCCCLLLLLLLLLLL             IIIIIIIIII
T:::::::::::::::::::::T                          d::::::d                           CCC::::::::::::CL:::::::::L             I::::::::I
T:::::::::::::::::::::T                          d::::::d                         CC:::::::::::::::CL:::::::::L             I::::::::I
T:::::TT:::::::TT:::::T                          d:::::d                         C:::::CCCCCCCC::::CLL:::::::LL             II::::::II
TTTTTT  T:::::T  TTTTTTooooooooooo       ddddddddd:::::d    ooooooooooo         C:::::C       CCCCCC  L:::::L                 I::::I  
        T:::::T      oo:::::::::::oo   dd::::::::::::::d  oo:::::::::::oo      C:::::C                L:::::L                 I::::I  
        T:::::T     o:::::::::::::::o d::::::::::::::::d o:::::::::::::::o     C:::::C                L:::::L                 I::::I  
        T:::::T     o:::::ooooo:::::od:::::::ddddd:::::d o:::::ooooo:::::o     C:::::C                L:::::L                 I::::I  
        T:::::T     o::::o     o::::od::::::d    d:::::d o::::o     o::::o     C:::::C                L:::::L                 I::::I  
        T:::::T     o::::o     o::::od:::::d     d:::::d o::::o     o::::o     C:::::C                L:::::L                 I::::I  
        T:::::T     o::::o     o::::od:::::d     d:::::d o::::o     o::::o     C:::::C                L:::::L                 I::::I  
        T:::::T     o::::o     o::::od:::::d     d:::::d o::::o     o::::o      C:::::C       CCCCCC  L:::::L         LLLLLL  I::::I  
      TT:::::::TT   o:::::ooooo:::::od::::::ddddd::::::ddo:::::ooooo:::::o       C:::::CCCCCCCC::::CLL:::::::LLLLLLLLL:::::LII::::::II
      T:::::::::T   o:::::::::::::::o d:::::::::::::::::do:::::::::::::::o        CC:::::::::::::::CL::::::::::::::::::::::LI::::::::I
      T:::::::::T    oo:::::::::::oo   d:::::::::ddd::::d oo:::::::::::oo           CCC::::::::::::CL::::::::::::::::::::::LI::::::::I
      TTTTTTTTTTT      ooooooooooo      ddddddddd   ddddd   ooooooooooo                CCCCCCCCCCCCCLLLLLLLLLLLLLLLLLLLLLLLLIIIIIIIIII
                                                                                                                                      
                                                                                                                                                                                                                                                            
                                                                                                                                      
                                                                                                                                      
                                                                                                                Sparky Inc. 2024
	`
	fmt.Println(banner)

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
		fmt.Println("Passcode: ", passcode)
	case "2":
		generatedPasscode = generateUniqueRandomIntegers(4)
		// set passcode to viper
		viper.Set("passcode", generatedPasscode)
		db.SaveUser(passcode)
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
		db.AddTask(description)
	},
}
var ListCommand = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement the logic to list all tasks
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
