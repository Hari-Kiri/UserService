// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type InsertUserDataInput struct {
	Name        string
	PhoneNumber string
	Password    string
}

type InsertUserDataOutput struct {
	Id int64
}
