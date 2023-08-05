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
}
type GetUserDataOutput struct {
	Id              int64
	SuccessfulLogin int64
	Password        string
}

type UpdateSuccessfullLoginInput struct {
	GetUserDataOutput
}
type UpdateSuccessfullLoginOutput struct {
	InsertUserDataOutput
}

type GetUserProfileInput struct {
	Id       int64
	Password string
}
type GetUserProfileOutput struct {
	Name        string
	PhoneNumber string
}

type UpdatePhoneNumberInput struct {
	Id          int64
	Password    string
	PhoneNumber string
}
type UpdatePhoneNumberOutput struct {
	Id int64
}

type UpdateNameInput struct {
	Id       int64
	Password string
	Name     string
}
type UpdateNameOutput struct {
	Id int64
}
