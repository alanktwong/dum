package tasks

import (
	i "awong/dotfiles/internal/tasks/installer"

	mock "github.com/stretchr/testify/mock"
)

func NewTestLogger(t interface {
	mock.TestingT
	Cleanup(func())
}) *i.MockLogger {
	l := i.NewMockLogger(t)
	l.On("Debugf", mock.Anything, mock.Anything).Return().Maybe()
	l.On("Infof", mock.Anything, mock.Anything).Return().Maybe()
	l.On("Debug", mock.Anything, mock.Anything).Return().Maybe()
	l.On("Errorf", mock.Anything, mock.Anything).Return().Maybe()
	l.On("Printlnf", mock.Anything, mock.Anything).Return(nil).Maybe()
	l.On("Warnf", mock.Anything, mock.Anything).Return().Maybe()
	return l
}
