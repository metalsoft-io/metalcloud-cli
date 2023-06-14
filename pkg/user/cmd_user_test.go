package user

import (
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	helper "github.com/metalsoft-io/metalcloud-cli/helpers"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	. "github.com/onsi/gomega"
)

func TestUserListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	user1 := metalcloud.UsersSearchResult{
		UserID:                 1,
		UserEmail:              "user1@email.com",
		UserDisplayName:        "John Oliver",
		UserBlocked:            false,
		UserLastLoginTimestamp: "2022-08-23T08:38:05Z",
		UserCreatedTimestamp:   "2021-02-17T10:35:51Z",
	}

	user2 := metalcloud.UsersSearchResult{
		UserID:                 2,
		UserEmail:              "user2@email.com",
		UserDisplayName:        "Mark Oliver",
		UserBlocked:            false,
		UserLastLoginTimestamp: "2022-08-23T08:38:05Z",
		UserCreatedTimestamp:   "2021-02-17T10:35:51Z",
	}
	user3 := metalcloud.UsersSearchResult{
		UserID:                 3,
		UserEmail:              "user3@email.com",
		UserDisplayName:        "Alexander Oliver",
		UserBlocked:            false,
		UserLastLoginTimestamp: "2020-08-23T08:38:05Z",
		UserCreatedTimestamp:   "2019-02-17T10:35:51Z",
	}

	user4 := metalcloud.UsersSearchResult{
		UserID:          4,
		UserEmail:       "user4@email.com",
		UserDisplayName: "Dan Oliver",
		UserBlocked:     true,
	}

	userList := []metalcloud.UsersSearchResult{
		user1,
		user2,
		user3,
		user4,
	}

	client := helper.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		UserSearch(gomock.Any()).
		Return(&userList, nil).
		AnyTimes()

	//test plaintext return
	format := ""
	cmd := command.Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}

	ret, err := userListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(Not(Equal("")))
	Expect(ret).To(ContainSubstring(user1.UserDisplayName))
	Expect(ret).To(ContainSubstring(user3.UserDisplayName))

	//test json return
	format = "json"
	cmd.Arguments["format"] = &format

	ret, err = userListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(Not(Equal("")))

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)
	Expect(err).To(BeNil())

	//test csv return
	format = "csv"
	cmd.Arguments["format"] = &format

	ret, err = userListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(Not(Equal("")))

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(err).To(BeNil())
	Expect(csv).NotTo(BeNil())
}
