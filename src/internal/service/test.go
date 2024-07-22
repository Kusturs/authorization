package service

import "context"

type TestService struct {
}

func NewTestService() *TestService {
	return &TestService{}
}

func (s *TestService) TestMethod(ctx context.Context, customersIDs []string) ([]interface{}, error) {
	return []interface{}{"test page"}, nil
}
