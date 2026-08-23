package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallPreCommitHookIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	hookPath := filepath.Join(dir, "pre-commit")
	original := "#!/bin/sh\necho original\n"
	if err := os.WriteFile(hookPath, []byte(original), 0755); err != nil {
		t.Fatal(err)
	}

	if err := installPreCommitHook(hookPath); err != nil {
		t.Fatal(err)
	}
	backupPath := hookPath + preCommitBackupSuffix
	backup, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(backup) != original {
		t.Fatalf("原 Hook 备份内容被改变: %q", backup)
	}

	first, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(first), preCommitMarker) {
		t.Fatal("安装后的 Hook 缺少 git-air 标记")
	}

	if err := installPreCommitHook(hookPath); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(second) != string(first) {
		t.Fatal("重复安装不应改写 git-air Hook")
	}
	backupAfter, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(backupAfter) != original {
		t.Fatal("重复安装不应改变原 Hook 备份")
	}
}

func TestPreCommitScriptRunsOnlyExecutableBackup(t *testing.T) {
	script := preCommitScript()
	if !strings.Contains(script, "if [ -x \"$old_hook\" ]; then") {
		t.Fatalf("Hook 没有检查备份文件的可执行权限:\n%s", script)
	}
}

func TestUninstallPreCommitHookRestoresOriginal(t *testing.T) {
	dir := t.TempDir()
	hookPath := filepath.Join(dir, "pre-commit")
	original := "#!/bin/sh\necho original\n"
	if err := os.WriteFile(hookPath, []byte(original), 0755); err != nil {
		t.Fatal(err)
	}
	if err := installPreCommitHook(hookPath); err != nil {
		t.Fatal(err)
	}
	if err := uninstallPreCommitHook(hookPath); err != nil {
		t.Fatal(err)
	}
	restored, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(restored) != original {
		t.Fatalf("原 Hook 未正确恢复: %q", restored)
	}
	if _, err := os.Stat(hookPath + preCommitBackupSuffix); !os.IsNotExist(err) {
		t.Fatalf("备份文件未清理: %v", err)
	}
}
