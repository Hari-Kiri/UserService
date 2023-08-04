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

type GetUserDataInput struct {
	PhoneNumber string
	Password    string
}
type GetUserDataOutput struct {
	Id              int64
	SuccessfulLogin int64
}

type UpdateSuccessfullLoginInput struct {
	GetUserDataOutput
}
type UpdateSuccessfullLoginOutput struct {
	InsertUserDataOutput
}
