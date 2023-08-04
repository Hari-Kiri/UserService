package repository

import "context"

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) InsertUserData(ctx context.Context, input InsertUserDataInput) (output InsertUserDataOutput, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		`INSERT INTO database.public.users(name, phone_number, password)
		VALUES ($1, $2, $3)
		RETURNING id;`,
		input.Name,
		input.PhoneNumber,
		input.Password,
	).Scan(&output.Id)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserData(ctx context.Context, input GetUserDataInput) (output GetUserDataOutput, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		`SELECT id, successful_login, password 
		FROM database.public.users 
		WHERE phone_number = $1;`,
		input.PhoneNumber,
	).Scan(&output.Id, &output.SuccessfulLogin, &output.Password)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UpdateSuccessfullLogin(ctx context.Context, input UpdateSuccessfullLoginInput) (output UpdateSuccessfullLoginOutput, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		`UPDATE database.public.users 
		SET successful_login = $1
		WHERE id = $2
		RETURNING id;`,
		input.SuccessfulLogin+1,
		input.Id,
	).Scan(&output.Id)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserProfile(ctx context.Context, input GetUserProfileInput) (output GetUserProfileOutput, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		`SELECT name, phone_number 
		FROM database.public.users 
		WHERE id = $1
		AND password = $2;`,
		input.Id,
		input.Password,
	).Scan(&output.Name, &output.PhoneNumber)
	if err != nil {
		return
	}
	return
}
