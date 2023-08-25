package main

import (
	"bufio"
	"cli/Constant"
	"cli/contract"
	"cli/entity"
	"cli/fileStore"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Task struct {
	ID         int
	Title      string
	DueDate    string
	CategoryID int
	IsDone     bool
	UserID     int
}
type Category struct {
	ID     int
	Title  string
	Color  string
	UserID int
}

var (
	userStorage       []entity.User
	authenticatedUser *entity.User
	taskStorage       []Task
	categoryStorage   []Category
	serializationMode string
)

const userStoragepath = "user.txt"

func main() {

	fmt.Println("Hello to TODO application")
	sm := flag.String("serialize-mode", Constant.JsonSerializationMode, "serialization to write data to disk")
	command := flag.String("command", "", "command to execute")
	flag.Parse()
	switch *sm {
	case Constant.TxtSerializationMode:
		serializationMode = Constant.TxtSerializationMode
	default:
		serializationMode = Constant.JsonSerializationMode
	}
	var userFileStore = fileStore.New(userStoragepath, serializationMode)

	//var userReadFileStore UserReadStore
	//var userReadStore = FileStore{
	//	FilePath: userStoragepath,
	//}
	//userReadFileStore = userReadStore
	//loudUserFromStorage(userReadFileStore, serializationMode)
	users := userFileStore.Load()
	userStorage = append(userStorage, users...)
	for {
		runCommand(userFileStore, *command)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()

	}
	//name:
	// fmt.Printf("userStorage:%+v\n", userStorage)
	/////%+v =fild name
}
func runCommand(store contract.UserWriteStore, command string) {
	if command != "register-user" && authenticatedUser == nil {
		fmt.Println("To use the program, you must login first")
		login()
		if authenticatedUser == nil {
			return
		}
	}
	/*1 var store UserWriteStore
	var userFileStore = FileStore{
		FilePath: userStoragepath,
	}
	store = userFileStore*/
	switch command {
	case "create-task":
		createTask()
	case "list-task":
		listTask()
	case "create-category":
		createCategory()
	case "register-user":
		//1 registerUser(store)
		registerUser(store)
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("command is not valid", command)
	}

}
func createTask() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, duedate, category string
	fmt.Println("Please Enter The Task Title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("Please Enter The Task Category id")
	scanner.Scan()
	category = scanner.Text()
	categoryID, err := strconv.Atoi(category)
	if err != nil {
		fmt.Println("categoryID is not valid ,%v\n", err)

		return
	}
	isFound := false
	for _, c := range categoryStorage {
		if c.ID == categoryID && c.UserID == authenticatedUser.ID {
			isFound = true

			break
		}
	}

	if !isFound {
		fmt.Println("category id is not found")

		return
	}

	fmt.Println("Please Enter The Task DueDate")
	scanner.Scan()
	duedate = scanner.Text()
	task := Task{
		ID:         len(taskStorage) + 1,
		Title:      title,
		DueDate:    duedate,
		CategoryID: categoryID,
		IsDone:     false,
		UserID:     authenticatedUser.ID,
	}
	taskStorage = append(taskStorage, task)

	fmt.Println("task:", title, category, duedate)

}

func createCategory() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, color string
	fmt.Println("Please Enter The category Title")
	scanner.Scan()
	title = scanner.Text()
	fmt.Println("Please Enter The category color")
	scanner.Scan()
	color = scanner.Text()
	fmt.Println("category:", title, color)
	category := Category{
		ID:     len(categoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}
	categoryStorage = append(categoryStorage, category)
}

func registerUser(store contract.UserWriteStore) {
	scanner := bufio.NewScanner(os.Stdin)
	var id int
	var Name, Email, password string
	fmt.Println("Please Enter The Name")
	scanner.Scan()
	Name = scanner.Text()
	fmt.Println("Please Enter The email")
	scanner.Scan()
	Email = scanner.Text()
	fmt.Println("Please Enter The password")
	scanner.Scan()
	password = scanner.Text()
	id = len(userStorage) + 1
	fmt.Println("user:", Email, password, " ID:", id)
	user := entity.User{
		ID:       id,
		Name:     Name,
		Email:    Email,
		Password: hashThePassword(password),
	}
	userStorage = append(userStorage, user)
	authenticatedUser = &user
	//writeUserToFile(user)
	store.Save(user)
}

func login() {
	scanner := bufio.NewScanner(os.Stdin)
	var Email, password string
	fmt.Println("Please Enter The email")
	scanner.Scan()
	Email = scanner.Text()
	fmt.Println("Please Enter The password")
	scanner.Scan()
	password = scanner.Text()

	for _, user := range userStorage {
		if user.Email == Email && user.Password == hashThePassword(password) {
			fmt.Println("login")
			authenticatedUser = &user

			break
		}

	}
	if authenticatedUser == nil {
		fmt.Println("the email or password is incorrect")

	}
}

func listTask() {
	for _, task := range taskStorage {
		if task.UserID == authenticatedUser.ID {
			fmt.Println(task)
		}
	}
}
func hashThePassword(password string) string {
	hash := md5.Sum([]byte(password))

	return hex.EncodeToString(hash[:])
}
