package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-vfs"
	"github.com/twpayne/go-vfs/vfst"

	"github.com/twpayne/chezmoi/chezmoi2/internal/chezmoitest"
)

func TestApplyCmd(t *testing.T) {
	for _, tc := range []struct {
		name      string
		extraRoot interface{}
		args      []string
		tests     []interface{}
	}{
		{
			name: "all",
			tests: []interface{}{
				vfst.TestPath("/home/user/.create",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o666&^chezmoitest.Umask),
					vfst.TestContentsString("# contents of .create\n"),
				),
				vfst.TestPath("/home/user/.dir",
					vfst.TestIsDir,
					vfst.TestModePerm(0o777&^chezmoitest.Umask),
				),
				vfst.TestPath("/home/user/.dir/file",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o666&^chezmoitest.Umask),
					vfst.TestContentsString("# contents of .dir/file\n"),
				),
				vfst.TestPath("/home/user/.dir/subdir",
					vfst.TestIsDir,
					vfst.TestModePerm(0o777&^chezmoitest.Umask),
				),
				vfst.TestPath("/home/user/.dir/subdir/file",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o666&^chezmoitest.Umask),
					vfst.TestContentsString("# contents of .dir/subdir/file\n"),
				),
				vfst.TestPath("/home/user/.empty",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o666&^chezmoitest.Umask),
					vfst.TestContents(nil),
				),
				vfst.TestPath("/home/user/.executable",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o777&^chezmoitest.Umask),
					vfst.TestContentsString("# contents of .executable\n"),
				),
				vfst.TestPath("/home/user/.file",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o666&^chezmoitest.Umask),
					vfst.TestContentsString("# contents of .file\n"),
				),
				vfst.TestPath("/home/user/.private",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o600&^chezmoitest.Umask),
					vfst.TestContentsString("# contents of .private\n"),
				),
				vfst.TestPath("/home/user/.remove",
					vfst.TestDoesNotExist,
				),
				vfst.TestPath("/home/user/.symlink",
					vfst.TestModeType(os.ModeSymlink),
					vfst.TestSymlinkTarget(filepath.FromSlash(".dir/subdir/file")),
				),
				vfst.TestPath("/home/user/.template",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o666&^chezmoitest.Umask),
					vfst.TestContentsString("key = value\n"),
				),
			},
		},
		{
			name: "all_with_--dry-run",
			args: []string{"--dry-run"},
			tests: []interface{}{
				vfst.TestPath("/home/user/.create",
					vfst.TestDoesNotExist,
				),
				vfst.TestPath("/home/user/.dir",
					vfst.TestDoesNotExist,
				),
				vfst.TestPath("/home/user/.empty",
					vfst.TestDoesNotExist,
				),
				vfst.TestPath("/home/user/.executable",
					vfst.TestDoesNotExist,
				),
				vfst.TestPath("/home/user/.file",
					vfst.TestDoesNotExist,
				),
				vfst.TestPath("/home/user/.private",
					vfst.TestDoesNotExist,
				),
				vfst.TestPath("/home/user/.remove",
					vfst.TestDoesNotExist,
				),
				vfst.TestPath("/home/user/.symlink",
					vfst.TestDoesNotExist,
				),
				vfst.TestPath("/home/user/.template",
					vfst.TestDoesNotExist,
				),
			},
		},
		{
			name: "dir",
			args: []string{"~/.dir"},
			tests: []interface{}{
				vfst.TestPath("/home/user/.dir",
					vfst.TestIsDir,
					vfst.TestModePerm(0o777&^chezmoitest.Umask),
				),
				vfst.TestPath("/home/user/.dir/file",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o666&^chezmoitest.Umask),
					vfst.TestContentsString("# contents of .dir/file\n"),
				),
				vfst.TestPath("/home/user/.dir/subdir",
					vfst.TestIsDir,
					vfst.TestModePerm(0o777&^chezmoitest.Umask),
				),
				vfst.TestPath("/home/user/.dir/subdir/file",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o666&^chezmoitest.Umask),
					vfst.TestContentsString("# contents of .dir/subdir/file\n"),
				),
			},
		},
		{
			name: "dir_with_--recursive=false",
			args: []string{"~/.dir", "--recursive=false"},
			tests: []interface{}{
				vfst.TestPath("/home/user/.dir",
					vfst.TestIsDir,
					vfst.TestModePerm(0o777&^chezmoitest.Umask),
				),
				vfst.TestPath("/home/user/.dir/file",
					vfst.TestDoesNotExist,
				),
				vfst.TestPath("/home/user/.dir/subdir",
					vfst.TestDoesNotExist,
				),
			},
		},
		{
			name: "create",
			args: []string{"~/.create"},
			extraRoot: map[string]interface{}{
				"/home/user/.create": "# existing contents of .create\n",
			},
			tests: []interface{}{
				vfst.TestPath("/home/user/.create",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o666&^chezmoitest.Umask),
					vfst.TestContentsString("# existing contents of .create\n"),
				),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			chezmoitest.WithTestFS(t, map[string]interface{}{
				"/home/user": map[string]interface{}{
					".config": map[string]interface{}{
						"chezmoi": map[string]interface{}{
							"chezmoi.toml": chezmoitest.JoinLines(
								`[data]`,
								`  variable = "value"`,
							),
						},
					},
					".local": map[string]interface{}{
						"share": map[string]interface{}{
							"chezmoi": map[string]interface{}{
								"create_dot_create": "# contents of .create\n",
								"dot_dir": map[string]interface{}{
									"file": "# contents of .dir/file\n",
									"subdir": map[string]interface{}{
										"file": "# contents of .dir/subdir/file\n",
									},
								},
								"dot_file":   "# contents of .file\n",
								"dot_remove": "",
								"dot_template.tmpl": chezmoitest.JoinLines(
									`key = {{ "value" }}`,
								),
								"empty_dot_empty":           "",
								"executable_dot_executable": "# contents of .executable\n",
								"private_dot_private":       "# contents of .private\n",
								"symlink_dot_symlink":       ".dir/subdir/file\n",
							},
						},
					},
				},
			}, func(fs vfs.FS) {
				if tc.extraRoot != nil {
					require.NoError(t, vfst.NewBuilder().Build(fs, tc.extraRoot))
				}
				require.NoError(t, newTestConfig(t, fs).execute(append([]string{"apply"}, tc.args...)))
				vfst.RunTests(t, fs, "", tc.tests)
			})
		})
	}
}
