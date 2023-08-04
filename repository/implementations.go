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
