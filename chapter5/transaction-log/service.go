package main

import (
	"cloud-native-go-study/chapter5/transaction-log/mylogger"
	"fmt"
)

var logger mylogger.TransactionLogger

// initializeTransactionLog 트랜잭션 로그를 바탕으로 store 를 구성한 뒤 파일 로깅 시작
func initializeTransactionLog() error {
	var err error

	logger, err = mylogger.NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event mylogger: %w", err)
	}

	events, errors := logger.ReadEvents()
	e, ok := mylogger.Event{}, true

	for err == nil && ok {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case mylogger.EventPut:
				err = Put(e.Key, e.Value)
			case mylogger.EventDelete:
				err = Delete(e.Key)
			}
		}
	}

	logger.Run()

	return err
}
