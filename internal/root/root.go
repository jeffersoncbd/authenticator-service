package root

import (
	"authenticator/internal/databases/postgresql"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func Run(store *postgresql.Queries, ctx context.Context) error {
	const appName = "authenticator"

	application, err := store.GetApplicationByName(ctx, appName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			id, err := store.InsertApplication(ctx, postgresql.InsertApplicationParams{
				Name: appName, Keys: []string{"users"},
			})
			if err != nil {
				return err
			}
			application.ID = id
		} else {
			return err
		}
	}
	fmt.Println("\n \033[0;32m✔\033[0m authenticator application inserted")

	group, err := store.GetGroupByName(ctx, postgresql.GetGroupByNameParams{
		Name:          "root",
		ApplicationID: application.ID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			id, err := store.InsertGroup(ctx, postgresql.InsertGroupParams{
				Name:          "root",
				ApplicationID: application.ID,
				Permissions:   []byte(`{ "users": 7 }`),
			})
			if err != nil {
				return err
			}
			group.ID = id
		} else {
			return err
		}
	}
	fmt.Println(" \033[0;32m✔\033[0m authenticator root group inserted")

	_, err = store.GetUser(ctx, os.Getenv("ROOT_MAIL"))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			hash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ROOT_PASS")), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			store.InsertUser(ctx, postgresql.InsertUserParams{
				Email:    os.Getenv("ROOT_MAIL"),
				Name:     "root",
				Password: string(hash),
				Groups:   []byte(fmt.Sprintf("{ \"%v\": \"%v\" }", application.ID, group.ID)),
			})
		} else {
			return err
		}
	}
	fmt.Println(" \033[0;32m✔\033[0m authenticator root user inserted")

	fmt.Printf(" \033[1;35m~\033[0m Authenticator ID: %s\n\n", application.ID)

	return nil
}
