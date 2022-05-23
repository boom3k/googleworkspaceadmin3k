package googleworkspace3k

import (
	"bytes"
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/groupsmigration/v1"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"strings"
)

type GroupsMigration3k struct {
	Service    *groupsmigration.Service
	AdminEmail string
	Domain     string
}

func BuildGroupsMigration3k(client *http.Client, adminEmail string, ctx context.Context) *GroupsMigration3k {
	groupMigration3k := &GroupsMigration3k{}
	service, err := groupsmigration.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
	groupMigration3k.Service = service
	groupMigration3k.AdminEmail = adminEmail
	groupMigration3k.Domain = strings.Split(adminEmail, "@")[1]
	log.Printf("GroupMigration -->Service: %v,\tCustomerID: %s,\tAdminEmail: %s,\tDomain: %s\n",
		&groupMigration3k.Service, groupMigration3k.AdminEmail, groupMigration3k.Domain)
	return groupMigration3k
}

func BuildGroupsMigration3kOauth2(adminEmail string, scopes []string, clientSecret, authorizationToken []byte, ctx context.Context) *GroupsMigration3k {
	config, err := google.ConfigFromJSON(clientSecret, scopes...)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
	token := &oauth2.Token{}
	err = json.Unmarshal(authorizationToken, token)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
	client := config.Client(context.Background(), token)
	return BuildGroupsMigration3k(client, adminEmail, ctx)
}

func (receiver *GroupsMigration3k) InsertEmail(groupEmail string, emailData []byte) (*groupsmigration.Groups, error) {
	media := bytes.NewReader(emailData)
	mediaOption := googleapi.ContentType("message/rfc822")
	response, err := receiver.Service.Archive.Insert(groupEmail).Media(media, mediaOption).Do()
	if err != nil {
		log.Println(err.Error())
		log.Printf("Mailbox [%s] import of (%d) bytes - FAILED\n", groupEmail, len(emailData))
		return nil, err
	}
	log.Printf("Mailbox [%s] import of (%d) bytes - SUCCESS\n", groupEmail, len(emailData))
	return response, nil
}
