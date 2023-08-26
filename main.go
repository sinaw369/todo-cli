package main

import (
	"bufio"
	"cli/Constant"
	"cli/contract"
	"cli/entity"
	"cli/repository/fileStore"
	"cli/repository/memorystore"
	"cli/service/task"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	userStorage       []entity.User
	authenticatedUser *entity.User
	categoryStorage   []entity.Category
	serializationMode string
)

const userStoragepath = "user.txt"

func main() {
	taskMemoryRepo := memorystore.NewTaskStore()
	taskService := task.NewService(taskMemoryRepo)
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
		runCommand(userFileStore, *command, taskService)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()

	}
	//name:
	// fmt.Printf("userStorage:%+v\n", userStorage)
	/////%+v =fild name
}
func runCommand(store contract.UserWriteStore, command string, taskService *task.Service) {
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
		createTask(taskService)
	case "list-task":
		listTask(taskService)
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
func createTask(taskService *task.Service) {
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
		fmt.Printf("categoryID is not valid ,%v\n", err)

		return
	}
	fmt.Println("Please Enter The Task DueDate")
	scanner.Scan()
	duedate = scanner.Text()
	response, err := taskService.Create(task.CreateRequest{
		Title:               title,
		DueDate:             duedate,
		CategoryID:          categoryID,
		AuthenticatedUserID: authenticatedUser.ID,
	})
	if err != nil {
		fmt.Println("error creating task:", err)

		return
	}

	fmt.Println("created task:", response.Task)
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
	category := entity.Category{
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

func listTask(taskService *task.Service) {
	userTasks, err := taskService.List(task.ListRequest{UserID: authenticatedUser.ID})
	if err != nil {
		fmt.Println("error", err)

		return
	}
	fmt.Println("The tasks", userTasks)
}
func hashThePassword(password string) string {
	hash := md5.Sum([]byte(password))

	return hex.EncodeToString(hash[:])
}
