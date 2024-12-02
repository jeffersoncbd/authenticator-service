package root

import (
	postgresql "authenticator/interfaces"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func Run(pool *pgxpool.Pool, ctx context.Context) error {
	store := postgresql.New(pool)

	applicationID, err := insertApplication(store, ctx)
	if err!= nil {
    return err
  }
	groupID, err := insertRootGroup(store, ctx, applicationID)
	if err!= nil {
    return err
  }

	err = insertRootUser(store, ctx, applicationID, groupID)
	if err!= nil {
    return err
  }

	fmt.Printf(" \033[1;35m~\033[0m Authenticator ID: %s\n\n", applicationID)

	return nil
}

func insertApplication(store *postgresql.Queries, ctx context.Context) (uuid.UUID, error) {
	const appName = "authenticator"

	application, err := store.GetApplicationByName(ctx, appName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			id, err := store.InsertApplication(ctx, appName)
			if err != nil {
				return uuid.Nil, err
			}
			application.ID = id
		} else {
			return uuid.Nil, err
		}
	}
	fmt.Println("\n \033[0;32m✔\033[0m authenticator application inserted")

	return application.ID, nil
}

func insertRootGroup(store *postgresql.Queries, ctx context.Context, applicationID uuid.UUID) (uuid.UUID, error) {
	group, err := store.GetGroupByName(ctx, postgresql.GetGroupByNameParams{
		Name:          "root",
		ApplicationID: applicationID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			id, err := store.InsertGroup(ctx, postgresql.InsertGroupParams{
				Name:          "root",
				ApplicationID: applicationID,
				Permissions:   []byte(`{ "users": 7, "applications": 7, "groups": 7 }`),
			})
			if err != nil {
				return uuid.Nil, err
			}
			group.ID = id
		} else {
			return uuid.Nil, err
		}
	}
	fmt.Println(" \033[0;32m✔\033[0m authenticator root group inserted")
	return group.ID, nil
}

func insertRootUser(store *postgresql.Queries, ctx context.Context, applicationID uuid.UUID, groupID uuid.UUID) error {
	_, err := store.GetUser(ctx, postgresql.GetUserParams{
		ApplicationID: applicationID,
		Email:         os.Getenv("ROOT_MAIL"),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			hash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ROOT_PASS")), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			store.InsertUser(ctx, postgresql.InsertUserParams{
				Email:         os.Getenv("ROOT_MAIL"),
				Name:          "root",
				Password:      string(hash),
				ApplicationID: applicationID,
				GroupID:       groupID,
			})
		} else {
			return err
		}
	}
	fmt.Println(" \033[0;32m✔\033[0m authenticator root user inserted")

	return nil
}
