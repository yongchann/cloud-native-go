package main

import "fmt"

var logger TransactionLogger

// initializeTransactionLog 트랜잭션 로그를 바탕으로 store 를 구성한 뒤 파일 로깅 시작
func initializeTransactionLog() error {
	var err error

	logger, err = NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := logger.ReadEvents()
	e, ok := Event{}, true

	for err == nil && ok {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventPut:
				err = Put(e.Key, e.Value)
			case EventDelete:
				err = Delete(e.Key)
			}
		}
	}

	logger.Run()

	return err

}
