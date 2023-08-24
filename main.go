package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Password string
	Email    string
}
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
type FileStore struct {
	FilePath string
}
type UserWriteStore interface {
	Save(u User)
}
type UserReadStore interface {
	load(serializationMode string) []User
}

func (f FileStore) Save(u User) {
	f.writeUserToFile(u)
}

var (
	userStorage       []User
	authenticatedUser *User
	taskStorage       []Task
	categoryStorage   []Category
	serializationMode string
)

const (
	userStoragepath = "user.txt"
	Jsonsm          = "json" //json serialization mode
	Txtsm           = "txt"
)

var userFileStore = FileStore{
	FilePath: userStoragepath,
}

func main() {

	fmt.Println("Hello to TODO application")
	sm := flag.String("serialize-mode", Jsonsm, "serialization to write data to disk")
	command := flag.String("command", "", "command to execute")
	flag.Parse()
	switch *sm {
	case Txtsm:
		serializationMode = Txtsm
	default:
		serializationMode = Jsonsm
	}
	//var userReadFileStore UserReadStore
	//var userReadStore = FileStore{
	//	FilePath: userStoragepath,
	//}
	//userReadFileStore = userReadStore
	//loudUserFromStorage(userReadFileStore, serializationMode)
	loudUserFromStorage(userFileStore, serializationMode)
	for {
		runCommand(*command)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()

	}
	//name:
	// fmt.Printf("userStorage:%+v\n", userStorage)
	/////%+v =fild name
}
func runCommand(command string) {
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
		registerUser(userFileStore)
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
	if err == nil {
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

func registerUser(store UserWriteStore) {
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
	user := User{
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
func (f FileStore) writeUserToFile(user User) {
	var data []byte
	file, err := os.OpenFile(f.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("error open file")
	}

	if serializationMode == Txtsm {
		data = []byte(fmt.Sprintf("id: %d, name: %s, email: %s, password: %s\n",
			user.ID, user.Name, user.Email, user.Password))

		return
	} else if serializationMode == Jsonsm {
		var jerr error
		data, jerr = json.Marshal(user)
		data = append(data, 10) // append \n
		if jerr != nil {
			fmt.Println("cant marshal user structure", jerr)

			return
		}

	} else {
		fmt.Println("invalid serialization mode ")

		return
	}

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("cant write to the file")
		return
	}
	defer file.Close()
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
			if user.Password == password {
				fmt.Println("login")
				authenticatedUser = &user

				break
			}

		}

	}
	if authenticatedUser == nil {
		fmt.Println("the email or password is incorrect")

	}
}
func (u User) print() {
	fmt.Println("Name: ", u.Name, "Email: ", u.Email, "Password: ", u.Password)
}

func listTask() {
	for _, task := range taskStorage {
		if task.UserID == authenticatedUser.ID {
			fmt.Println(task)
		}
	}
}
func loudUserFromStorage(store UserReadStore, serializationMode string) {
	users := store.load(serializationMode)
	userStorage = append(userStorage, users...)
}
func (f FileStore) load(serializationMode string) []User {
	var uStore []User
	file, err := os.Open(f.FilePath)
	if err != nil {
		fmt.Println("cant read from the file: ", err)

		return nil
	}
	var data = make([]byte, 1024)
	_, err = file.Read(data)
	if err != nil {
		fmt.Println("can't read from the file", err)
	}
	var datastring = string(data)

	//datastring = strings.Trim(datastring, "\n")
	userSlice := strings.Split(datastring, "\n")

	for _, u := range userSlice {
		var userStruct = User{}
		switch serializationMode {
		case Txtsm:
			userStruct, err = deseerializeTxt(u)
			if err != nil {
				fmt.Println("can't deserialize")
				return nil
			}

		case Jsonsm:
			if u[0] != '{' && u[len(u)-1] != '}' {
				continue
			}
			err := json.Unmarshal([]byte(u), &userStruct)
			if err != nil {
				fmt.Println("cant deserialize json")
			}

		}
		fmt.Println("unmarshal:", userStruct)
		uStore = append(uStore, userStruct)
	}
	//fmt.Println(data)
	return uStore
}
func deseerializeTxt(userstr string) (User, error) {
	var user = User{}
	if userstr == "" {
		return User{}, errors.New("user string is empty")
	}
	userField := strings.Split(userstr, ",")
	for _, field := range userField {
		values := strings.Split(field, ": ")
		if len(values) != 2 {
			fmt.Println("field error")
			continue
		}
		fieldName := strings.ReplaceAll(values[0], " ", "")
		fieldValue := values[1]

		switch fieldName {
		case "id":
			id, err := strconv.Atoi(fieldValue)
			if err != nil {
				fmt.Println("strconv error", err)

				return User{}, errors.New("strconv error")
			}
			user.ID = id
		case "name":
			user.Name = fieldValue
		case "email":
			user.Email = fieldValue
		case "password":
			user.Password = fieldValue

		}

	}
	fmt.Printf("user: %+v\n", user)
	return user, nil

}
func hashThePassword(password string) string {
	hash := md5.Sum([]byte(password))

	return hex.EncodeToString(hash[:])
}
