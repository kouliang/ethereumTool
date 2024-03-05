package job

import (
	"fmt"
	"os"
	"time"

	"github.com/kouliang/ethereumtool/account"
	"github.com/kouliang/ethereumtool/client"
	"github.com/kouliang/ethereumtool/email"
	"github.com/kouliang/ethereumtool/historylog"
)

var EmailNotAddress []string

type IJob interface {
	Run()
	Close()
}

type Job struct {
	Name             string
	ContractAddress  string
	attemptNumber    int64
	MaxAttemptNumber int64

	Account         *account.KLAccount
	HLog            *historylog.HLog
	SendTransaction func() error
}

func New(name_ string, contractAddress_ string, hexkey_ string) (*Job, error) {

	logPath := fmt.Sprintf("./%s.log", name_)
	file, err := historylog.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	hlog := historylog.HLogWithFile(file)

	if !account.IsAvailableAddress(contractAddress_) {
		return nil, fmt.Errorf("unavailable to address")
	}

	account_, err := account.New(hexkey_)
	if err != nil {
		return nil, err
	}

	hlog.Println("========================================================")
	hlog.Println("Job name:", name_)
	hlog.Println("Connected!!! ChainId:", client.ChainID())
	hlog.Println("From address:", account_.AddressStr)
	hlog.Println("Contract address:", contractAddress_)
	hlog.Println("========================================================")

	return &Job{
		Name:            name_,
		ContractAddress: contractAddress_,
		Account:         account_,

		HLog:             hlog,
		MaxAttemptNumber: 3,
	}, nil
}

func (job *Job) Run() {
	job.HLog.Println("========================================================")

	job.HLog.TakeOutHistory()
	job.HLog.Println("Index:", job.attemptNumber)

	err := job.SendTransaction()
	record := job.HLog.TakeOutHistory()
	if err == nil {
		record = fmt.Sprintf("Success!!!\n%s", record)
		job.HLog.Println("Success!!!")

		job.attemptNumber = 0
	} else {
		record = fmt.Sprintf("Fail!!!\n%s", record)
		job.HLog.Println("Fail!!!")

		job.ReSend()
	}

	if len(EmailNotAddress) > 0 {
		emailMsg, err := email.SenEmail(job.Name, record, EmailNotAddress)
		if err != nil {
			job.HLog.Println("SendEmail failed:", err.Error())
		} else {
			job.HLog.Println("EmailMsg:", emailMsg)
		}
	} else {
		job.HLog.Println("Donot need to send email")
	}
}

func (job *Job) Close() {
	job.HLog.File.Close()
}

func (job *Job) ReSend() {
	job.attemptNumber += 1

	if job.attemptNumber > job.MaxAttemptNumber {
		job.attemptNumber = 0
		return
	}

	timer := time.NewTimer(120 * time.Second)
	go func() {
		<-timer.C
		job.Run()
	}()
}
