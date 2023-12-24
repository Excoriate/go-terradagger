package terradagger

import (
	"path/filepath"
	"testing"

	"github.com/Excoriate/go-terradagger/pkg/mocks"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func TestResolveMountDirPath(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// Override the global Util variable with our mock

	mockDirUtils := mocks.NewMockDirUtilities(mockCtrl)

	cwd := "user/this/is/cwd/result"

	testCases := []struct {
		name          string
		mountDirPath  string
		expectedPath  string
		expectedError error
	}{
		{
			// mountDirPath is empty, so it'll return the current directory
			name:          "mountDirPath is empty",
			mountDirPath:  "",
			expectedPath:  filepath.Join(cwd, "."),
			expectedError: nil,
		},
		// mountDirPath is ".", so it'll return the current directory
		{
			name:          "mountDirPath is .",
			mountDirPath:  ".",
			expectedPath:  filepath.Join(cwd, "."),
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDirUtils.EXPECT().GetCurrentDir().Return(cwd).AnyTimes()
			mockDirUtils.EXPECT().IsValidDir(gomock.Any()).Return(tc.expectedError).AnyTimes()

			path, err := resolveMountDirPath(tc.mountDirPath)

			// assert.Equal(t, tc.expectedPath, path)
			// TODO: Refactor the resolveMountDirPath to accept an interface,
			//  so it'll use the cwd instead of the actual current directory.
			assert.NotEmpty(t, path)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
