package main

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	. "github.com/onsi/gomega"
)

func TestJobsListCmd(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	list := []metalcloud.AFCSearchResult{
		{
			AFCID:               10,
			AFCCreatedTimestamp: "2006-01-02T15:04:05Z",
			AFCFunctionName:     "test_func",
			AFCParamsJSON:       "['param1',10,'param2']",
		},
		{
			AFCID:               11,
			AFCCreatedTimestamp: "2006-01-02T15:04:05Z",
			AFCFunctionName:     "test_func2",
			AFCParamsJSON:       "['p1',55,'param2']",
		},
		{
			AFCID:               13,
			AFCCreatedTimestamp: "2006-01-02T15:04:05Z",
			AFCFunctionName:     "infrastructure_provision",
			AFCParamsJSON:       "[ \"param1\",\"real_func\",[\"real_param1\",\"real_param1\"] ]",
		},
	}

	client.EXPECT().
		AFCSearch(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&list, nil).
		AnyTimes()

	cmd := MakeCommand(map[string]interface{}{})

	ret, err := jobListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("ID"))
	Expect(ret).To(ContainSubstring("test_func(['param1',10,'param2'])"))
	Expect(ret).To(ContainSubstring("test_func2(['p1',55,'param2'])"))
	Expect(ret).To(ContainSubstring("real_func([real_param1 real_param1])"))
	Expect(ret).NotTo(ContainSubstring("infrastructure_provision"))

	client.EXPECT().
		AFCSearch(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&list, nil).
		AnyTimes()

	cmd = MakeCommand(map[string]interface{}{
		"format": "csv",
	})

	ret, err = jobListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("test_func(['param1',10,'param2'])"))
	Expect(ret).To(ContainSubstring("test_func2(['p1',55,'param2'])"))
	Expect(ret).To(ContainSubstring("real_func([real_param1 real_param1])"))
	Expect(ret).To(ContainSubstring("ID,STATUS,DURATION,AFFECTS,RETRIES,REQUEST,RESPONSE"))

	cmd = MakeCommand(map[string]interface{}{
		"filter": "test:test",
	})

	client.EXPECT().
		AFCSearch("test:test", gomock.Any(), gomock.Any()).
		Return(&list, nil).
		AnyTimes()

	ret, err = jobListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("ID"))
	Expect(ret).To(ContainSubstring("test_func(['param1',10,'param2'])"))
	Expect(ret).To(ContainSubstring("test_func2(['p1',55,'param2'])"))
	Expect(ret).To(ContainSubstring("real_func([real_param1 real_param1])"))

}

func TestJobsGetCmd(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	afc := metalcloud.AFC{

		AFCID:               13,
		AFCCreatedTimestamp: "2006-01-02T15:04:05Z",
		AFCFunctionName:     "infrastructure_provision",
		AFCParamsJSON:       "[ \"param1\",\"real_func\",[\"real_param1\",\"real_param1\"] ]",
	}

	client.EXPECT().
		AFCGet(13).
		Return(&afc, nil).
		AnyTimes()

	cmd := MakeCommand(map[string]interface{}{"job_id": "13"})

	ret, err := jobGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("ID"))
	Expect(ret).To(ContainSubstring("real_func([real_param1 real_param1])"))
	Expect(ret).NotTo(ContainSubstring("infrastructure_provision"))

	afc2 := metalcloud.AFC{

		AFCID:               11,
		AFCCreatedTimestamp: "2006-01-02T15:04:05Z",
		AFCFunctionName:     "test_func2",
		AFCParamsJSON:       "['p1',55,'param2']",
	}
	client.EXPECT().
		AFCGet(13).
		Return(&afc2, nil).
		AnyTimes()

	cmd = MakeCommand(map[string]interface{}{
		"format": "csv",
		"job_id": "13",
	})

	ret, err = jobGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("real_func([real_param1 real_param1])"))
	Expect(ret).To(ContainSubstring("ID,STATUS,DURATION,AFFECTS,RETRIES,REQUEST,RESPONSE"))

}
