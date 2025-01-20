package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/pkg/cache"
	"github.com/jiu-u/oai-api/pkg/config"
	"github.com/jiu-u/oai-api/pkg/encrypte"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
	"github.com/jiu-u/oai-api/pkg/sid"
	"os"
	"testing"
)

const configFile = "../config/local.yaml"

var channelSvc service.ChannelService

// setup 初始化测试环境
func setup() {
	fmt.Println("Setting up test environment")
	cfg := config.LoadConfig(configFile)
	logger := log.NewLogger(cfg)
	jwtJWT := jwt.NewJwt(cfg)
	sidSid := sid.NewSid()
	db := repository.NewDB(cfg)
	repositoryRepository := repository.NewRepository(logger, db)
	transaction := repository.NewTransaction(repositoryRepository)
	cacheCache := cache.New()
	serviceService := service.NewService(sidSid, transaction, logger, jwtJWT, cacheCache)
	channelRepository := repository.NewChannelRepository(repositoryRepository)
	channelModelRepository := repository.NewChannelModelRepository(repositoryRepository)
	loadBalanceServiceBeta := service.NewLoadBalanceServiceBeta(serviceService, channelRepository, channelModelRepository)
	channelSvc = service.NewChannelService(serviceService, channelRepository, channelModelRepository, loadBalanceServiceBeta)
}

// teardown 清理测试环境
func teardown() {
	fmt.Println("Tearing down test environment")
	// 执行必要的清理操作，比如删除临时文件、目录等
}

func TestMain(m *testing.M) {
	// 在执行任何测试之前进行设置
	setup()

	// 执行测试
	code := m.Run()

	// 在测试完成后进行清理
	teardown()

	// 使用测试返回的状态码终止程序
	os.Exit(code)

}

// 示例测试函数
func TestExample(t *testing.T) {
	t.Log("Running TestExample")
	// 这里可以写你的测试代码
	if 1+1 != 2 {
		t.Error("TestExample failed")
	}
}

func TestChannelService_GetChannels(t *testing.T) {
	resp, err := channelSvc.GetChannels(context.Background(), &apiV1.ChannelQueryRequest{
		Name:     "",
		Type:     "",
		Status:   1,
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		t.Error(err)
	}
	for _, item := range resp.List {
		fmt.Printf("%+v\n", item)
	}
}

func TestHashId(t *testing.T) {
	const compare = "a53e0dc204795c5b3c380fd4d37f94dc5b9e69d524744f9188f7734f789c1ce6"
	const Type = ""
	const EndPoint = ""
	const APIKey = ""
	hashId := encrypte.Sha256Encode(fmt.Sprintf("%s%s%s", Type, EndPoint, APIKey))
	fmt.Println(hashId)
	if hashId != compare {
		t.Error("hashId error")
	}
}

func TestOption2Print(t *testing.T) {
	type Option struct {
		Key string
		Val string
	}
	jsonBytes := []byte(`{"key":"LINUX_DO_CLIENT_ID","val":"123123123"}`)
	var opt1 Option
	err := json.Unmarshal(jsonBytes, &opt1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(opt1)

	jsonBytes = []byte(`{"key":"LINUX_DO_AUTH_ENABLE","val":true}`)
	var opt2 Option
	err = sonic.Unmarshal(jsonBytes, &opt2)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(opt2)
}
