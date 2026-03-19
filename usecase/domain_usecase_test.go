package usecase

import (
	interfacemocks "ai-agent/interface/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDomainUsecase_GetInstructions(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repoMock := interfacemocks.NewMockDomainRepo(t)
		u := NewDomainUsecase(repoMock)

		repoMock.EXPECT().GetDomainInstructions("books").Return("You are a book assistant.", nil)

		instructions, err := u.GetInstructions("books")

		require.NoError(t, err)
		require.Equal(t, "You are a book assistant.", instructions)
	})

	t.Run("repo error", func(t *testing.T) {
		repoMock := interfacemocks.NewMockDomainRepo(t)
		u := NewDomainUsecase(repoMock)

		repoMock.EXPECT().GetDomainInstructions("books").Return("", errors.New("domain not found"))

		instructions, err := u.GetInstructions("books")

		require.Error(t, err)
		require.Empty(t, instructions)
		require.ErrorContains(t, err, "domain not found")
	})
}
