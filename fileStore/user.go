package fileStore

import (
	"cli/Constant"
	"cli/entity"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type FileStore struct {
	FilePath          string
	SerializationMode string
}

func New(path string, SerializationMode string) FileStore {
	return FileStore{FilePath: path,
		SerializationMode: SerializationMode}
}

func (f FileStore) Save(u entity.User) {
	f.writeUserToFile(u)
}

func (f FileStore) writeUserToFile(user entity.User) {
	var data []byte
	file, err := os.OpenFile(f.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("error open file")
	}

	if f.SerializationMode == Constant.TxtSerializationMode {
		data = []byte(fmt.Sprintf("id: %d, name: %s, email: %s, password: %s\n",
			user.ID, user.Name, user.Email, user.Password))

		return
	} else if f.SerializationMode == Constant.JsonSerializationMode {
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

func (f FileStore) Load() []entity.User {
	var uStore []entity.User
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
		var userStruct = entity.User{}
		switch f.SerializationMode {
		case Constant.TxtSerializationMode:
			userStruct, err = deseerializeTxt(u)
			if err != nil {
				fmt.Println("can't deserialize")
				return nil
			}

		case Constant.JsonSerializationMode:
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
func deseerializeTxt(userstr string) (entity.User, error) {
	var user = entity.User{}
	if userstr == "" {
		return entity.User{}, errors.New("user string is empty")
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

				return entity.User{}, errors.New("strconv error")
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
