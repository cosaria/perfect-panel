package db_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBaselineMigrationContract(t *testing.T) {
	content := loadBaselineMigration(t)
	sqlText := strings.ToLower(content)

	assertContains(t, sqlText, "create extension if not exists pgcrypto")
	assertContains(t, sqlText, "create table users")
	assertContains(t, sqlText, "id uuid primary key default gen_random_uuid()")

	assertContains(t, sqlText, "create table roles")
	assertContains(t, sqlText, "create table system_settings")
	assertContains(t, sqlText, "scope text not null")
	assertContains(t, sqlText, "key text not null")
	assertContains(t, sqlText, "value_json jsonb not null")
	assertContains(t, sqlText, "unique index idx_system_settings_scope_key on system_settings(scope, key)")

	assertContains(t, sqlText, "create table outbox_events")
	assertContains(t, sqlText, "aggregate_type text not null")
	assertContains(t, sqlText, "aggregate_id uuid not null")
}

func loadBaselineMigration(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("无法定位测试文件路径")
	}

	path := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", "internal", "platform", "db", "migrations", "0001_baseline.sql"))
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读取 baseline migration 失败: %v", err)
	}
	return string(content)
}

func assertContains(t *testing.T, text string, expected string) {
	t.Helper()
	if !strings.Contains(text, expected) {
		t.Fatalf("migration 契约不满足，缺少片段: %s", expected)
	}
}
