package job

type LPJob struct {
	*Job
}

func new(name_ string, contractAddress_ string, hexkey_ string) (*LPJob, error) {
	job, err := New(name_, contractAddress_, hexkey_)
	if err == nil {
		j := LPJob{Job: job}
		j.Job.SendTransaction = j.SendTransaction
		j.MaxAttemptNumber = 0
		return &j, nil
	} else {
		return nil, err
	}
}

func (job *LPJob) SendTransaction() (err error) {
	// callData, _ := contractAbi.Pack("archivePoolInfo", paramTime)
	// err = client.SendData(job.Account, job.To, callData, job.HLog)
	return err
}
