package service

import "context"

type TestService interface {
	TestMethod(ctx context.Context, customersIDs []string) ([]interface{}, error)
}
