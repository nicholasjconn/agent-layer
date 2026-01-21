package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasjconn/agent-layer/internal/config"
)

func TestBuildVSCodeSettings(t *testing.T) {
	project := &config.ProjectConfig{
		Config: config.Config{
			Approvals: config.ApprovalsConfig{Mode: "commands"},
		},
		CommandsAllow: []string{"git status"},
	}

	settings, err := buildVSCodeSettings(project)
	if err != nil {
		t.Fatalf("buildVSCodeSettings error: %v", err)
	}
	if len(settings.ChatToolsTerminalAutoApprove) != 1 {
		t.Fatalf("expected 1 auto-approve entry")
	}
}

func TestWriteVSCodeSettings(t *testing.T) {
	root := t.TempDir()
	project := &config.ProjectConfig{
		Config: config.Config{
			Approvals: config.ApprovalsConfig{Mode: "commands"},
		},
		CommandsAllow: []string{"git status"},
	}

	if err := WriteVSCodeSettings(root, project); err != nil {
		t.Fatalf("WriteVSCodeSettings error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".vscode", "settings.json")); err != nil {
		t.Fatalf("expected settings.json: %v", err)
	}
}

func TestWriteVSCodeSettingsError(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	project := &config.ProjectConfig{}
	if err := WriteVSCodeSettings(file, project); err == nil {
		t.Fatalf("expected error")
	}
}

func TestWriteVSCodeSettingsWriteError(t *testing.T) {
	root := t.TempDir()
	vscodeDir := filepath.Join(root, ".vscode")
	if err := os.MkdirAll(vscodeDir, 0o500); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	project := &config.ProjectConfig{
		Config: config.Config{
			Approvals: config.ApprovalsConfig{Mode: "none"},
		},
	}
	if err := WriteVSCodeSettings(root, project); err == nil {
		t.Fatalf("expected error")
	}
}

func TestBuildVSCodeMCPConfig(t *testing.T) {
	enabled := true
	project := &config.ProjectConfig{
		Config: config.Config{
			MCP: config.MCPConfig{
				Servers: []config.MCPServer{
					{
						ID:        "example",
						Enabled:   &enabled,
						Transport: "http",
						URL:       "https://example.com?token=${TOKEN}",
					},
				},
			},
		},
		Env: map[string]string{"TOKEN": "abc"},
	}

	cfg, err := buildVSCodeMCPConfig(project)
	if err != nil {
		t.Fatalf("buildVSCodeMCPConfig error: %v", err)
	}
	server, ok := cfg.Servers["example"]
	if !ok {
		t.Fatalf("expected server entry")
	}
	if server.Type != "http" {
		t.Fatalf("unexpected server type: %s", server.Type)
	}
	// VS Code uses ${env:VAR} syntax - VS Code resolves at runtime.
	if server.URL != "https://example.com?token=${env:TOKEN}" {
		t.Fatalf("unexpected url: %s", server.URL)
	}
}

func TestBuildVSCodeMCPConfigHeadersAndEnv(t *testing.T) {
	enabled := true
	project := &config.ProjectConfig{
		Config: config.Config{
			MCP: config.MCPConfig{
				Servers: []config.MCPServer{
					{
						ID:        "http",
						Enabled:   &enabled,
						Transport: "http",
						URL:       "https://example.com",
						Headers:   map[string]string{"X-Token": "${TOKEN}"},
					},
					{
						ID:        "stdio",
						Enabled:   &enabled,
						Transport: "stdio",
						Command:   "tool-${TOKEN}",
						Args:      []string{"--flag", "${KEY}"},
						Env:       map[string]string{"API_KEY": "${KEY}"},
					},
				},
			},
		},
		Env: map[string]string{"TOKEN": "abc", "KEY": "123"},
	}

	cfg, err := buildVSCodeMCPConfig(project)
	if err != nil {
		t.Fatalf("buildVSCodeMCPConfig error: %v", err)
	}
	// VS Code uses ${env:VAR} syntax - VS Code resolves at runtime.
	httpServer, ok := cfg.Servers["http"]
	if !ok {
		t.Fatalf("expected http server entry")
	}
	if httpServer.Headers["X-Token"] != "${env:TOKEN}" {
		t.Fatalf("unexpected header value: %s", httpServer.Headers["X-Token"])
	}

	server, ok := cfg.Servers["stdio"]
	if !ok {
		t.Fatalf("expected stdio server entry")
	}
	if server.Type != "stdio" {
		t.Fatalf("unexpected server type: %s", server.Type)
	}
	if server.Command != "tool-${env:TOKEN}" {
		t.Fatalf("unexpected command: %s", server.Command)
	}
	if len(server.Args) != 2 || server.Args[1] != "${env:KEY}" {
		t.Fatalf("unexpected args: %#v", server.Args)
	}
	if server.Env["API_KEY"] != "${env:KEY}" {
		t.Fatalf("unexpected env value: %s", server.Env["API_KEY"])
	}
}

func TestWriteVSCodeMCPConfig(t *testing.T) {
	root := t.TempDir()
	enabled := true
	project := &config.ProjectConfig{
		Config: config.Config{
			MCP: config.MCPConfig{
				Servers: []config.MCPServer{
					{
						ID:        "example",
						Enabled:   &enabled,
						Transport: "http",
						URL:       "https://example.com?token=${TOKEN}",
					},
				},
			},
		},
		Env: map[string]string{"TOKEN": "abc"},
	}

	if err := WriteVSCodeMCPConfig(root, project); err != nil {
		t.Fatalf("WriteVSCodeMCPConfig error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".vscode", "mcp.json")); err != nil {
		t.Fatalf("expected mcp.json: %v", err)
	}
}

func TestWriteVSCodeMCPConfigError(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	project := &config.ProjectConfig{}
	if err := WriteVSCodeMCPConfig(file, project); err == nil {
		t.Fatalf("expected error")
	}
}

func TestWriteVSCodeMCPConfigMissingEnv(t *testing.T) {
	root := t.TempDir()
	enabled := true
	project := &config.ProjectConfig{
		Config: config.Config{
			MCP: config.MCPConfig{
				Servers: []config.MCPServer{
					{
						ID:        "example",
						Enabled:   &enabled,
						Transport: "http",
						URL:       "https://example.com?token=${TOKEN}",
					},
				},
			},
		},
		Env: map[string]string{},
	}

	if err := WriteVSCodeMCPConfig(root, project); err == nil {
		t.Fatalf("expected error")
	}
}

func TestBuildVSCodeMCPConfigMissingEnv(t *testing.T) {
	enabled := true
	project := &config.ProjectConfig{
		Config: config.Config{
			MCP: config.MCPConfig{
				Servers: []config.MCPServer{
					{
						ID:        "example",
						Enabled:   &enabled,
						Transport: "http",
						URL:       "https://example.com?token=${TOKEN}",
					},
				},
			},
		},
		Env: map[string]string{},
	}

	_, err := buildVSCodeMCPConfig(project)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestWriteVSCodeLaunchers(t *testing.T) {
	root := t.TempDir()

	if err := WriteVSCodeLaunchers(root); err != nil {
		t.Fatalf("WriteVSCodeLaunchers error: %v", err)
	}

	// Verify macOS .command launcher
	shPath := filepath.Join(root, ".agent-layer", "open-vscode.command")
	shInfo, err := os.Stat(shPath)
	if err != nil {
		t.Fatalf("expected open-vscode.command: %v", err)
	}
	if shInfo.Mode().Perm() != 0o755 {
		t.Fatalf("expected 0755 permissions on .command file, got %o", shInfo.Mode().Perm())
	}

	// Verify macOS .app bundle structure
	appDir := filepath.Join(root, ".agent-layer", "open-vscode.app")
	if _, err := os.Stat(appDir); err != nil {
		t.Fatalf("expected open-vscode.app directory: %v", err)
	}

	infoPlistPath := filepath.Join(appDir, "Contents", "Info.plist")
	if _, err := os.Stat(infoPlistPath); err != nil {
		t.Fatalf("expected Info.plist: %v", err)
	}

	execPath := filepath.Join(appDir, "Contents", "MacOS", "open-vscode")
	execInfo, err := os.Stat(execPath)
	if err != nil {
		t.Fatalf("expected open-vscode executable: %v", err)
	}
	if execInfo.Mode().Perm() != 0o755 {
		t.Fatalf("expected 0755 permissions on app executable, got %o", execInfo.Mode().Perm())
	}

	// Verify Windows launcher
	batPath := filepath.Join(root, ".agent-layer", "open-vscode.bat")
	batInfo, err := os.Stat(batPath)
	if err != nil {
		t.Fatalf("expected open-vscode.bat: %v", err)
	}
	if batInfo.Mode().Perm() != 0o755 {
		t.Fatalf("expected 0755 permissions on .bat file, got %o", batInfo.Mode().Perm())
	}
}

func TestWriteVSCodeLaunchersContent(t *testing.T) {
	root := t.TempDir()

	if err := WriteVSCodeLaunchers(root); err != nil {
		t.Fatalf("WriteVSCodeLaunchers error: %v", err)
	}

	// Verify macOS .command launcher content
	shPath := filepath.Join(root, ".agent-layer", "open-vscode.command")
	shContent, err := os.ReadFile(shPath)
	if err != nil {
		t.Fatalf("read .command file: %v", err)
	}
	shStr := string(shContent)

	if len(shStr) == 0 {
		t.Fatal("macOS launcher is empty")
	}
	if shStr[:2] != "#!" {
		t.Fatal("macOS launcher missing shebang")
	}
	if !strings.Contains(shStr, "CODEX_HOME") {
		t.Fatal("macOS launcher missing CODEX_HOME")
	}
	if !strings.Contains(shStr, "code .") {
		t.Fatal("macOS launcher missing 'code .' command")
	}
	if !strings.Contains(shStr, "Shell Command: Install") {
		t.Fatal("macOS launcher missing install instructions")
	}

	// Verify macOS .app bundle content
	appDir := filepath.Join(root, ".agent-layer", "open-vscode.app")

	infoPlistContent, err := os.ReadFile(filepath.Join(appDir, "Contents", "Info.plist"))
	if err != nil {
		t.Fatalf("read Info.plist: %v", err)
	}
	infoPlistStr := string(infoPlistContent)
	if !strings.Contains(infoPlistStr, "CFBundleExecutable") {
		t.Fatal("Info.plist missing CFBundleExecutable")
	}
	if !strings.Contains(infoPlistStr, "com.agent-layer.open-vscode") {
		t.Fatal("Info.plist missing bundle identifier")
	}
	if !strings.Contains(infoPlistStr, "LSUIElement") {
		t.Fatal("Info.plist missing LSUIElement (needed to hide from dock)")
	}

	execContent, err := os.ReadFile(filepath.Join(appDir, "Contents", "MacOS", "open-vscode"))
	if err != nil {
		t.Fatalf("read app executable: %v", err)
	}
	execStr := string(execContent)
	if execStr[:2] != "#!" {
		t.Fatal("app executable missing shebang")
	}
	if !strings.Contains(execStr, "CODEX_HOME") {
		t.Fatal("app executable missing CODEX_HOME")
	}
	if !strings.Contains(execStr, "Contents/Resources/app/bin/code") {
		t.Fatal("app executable missing full path to VS Code CLI")
	}
	if !strings.Contains(execStr, "/Applications/Visual Studio Code.app") {
		t.Fatal("app executable missing VS Code app path")
	}
	if !strings.Contains(execStr, "osascript") {
		t.Fatal("app executable missing osascript error dialog")
	}

	// Verify Windows launcher content
	batPath := filepath.Join(root, ".agent-layer", "open-vscode.bat")
	batContent, err := os.ReadFile(batPath)
	if err != nil {
		t.Fatalf("read .bat file: %v", err)
	}
	batStr := string(batContent)

	if len(batStr) == 0 {
		t.Fatal("Windows launcher is empty")
	}
	if !strings.Contains(batStr, "@echo off") {
		t.Fatal("Windows launcher missing @echo off")
	}
	if !strings.Contains(batStr, "CODEX_HOME") {
		t.Fatal("Windows launcher missing CODEX_HOME")
	}
	if !strings.Contains(batStr, "code .") {
		t.Fatal("Windows launcher missing 'code .' command")
	}
	if !strings.Contains(batStr, "Shell Command: Install") {
		t.Fatal("Windows launcher missing install instructions")
	}
}

func TestWriteVSCodeLaunchersDirectoryError(t *testing.T) {
	root := t.TempDir()
	// Create a file where the directory should be
	file := filepath.Join(root, ".agent-layer")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if err := WriteVSCodeLaunchers(root); err == nil {
		t.Fatalf("expected error when .agent-layer is a file")
	}
}

func TestWriteVSCodeLaunchersWriteError(t *testing.T) {
	root := t.TempDir()
	agentLayerDir := filepath.Join(root, ".agent-layer")
	if err := os.MkdirAll(agentLayerDir, 0o500); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := WriteVSCodeLaunchers(root); err == nil {
		t.Fatalf("expected error when directory is read-only")
	}
}

func TestWriteVSCodeAppBundle(t *testing.T) {
	root := t.TempDir()
	agentLayerDir := filepath.Join(root, ".agent-layer")
	if err := os.MkdirAll(agentLayerDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := writeVSCodeAppBundle(agentLayerDir); err != nil {
		t.Fatalf("writeVSCodeAppBundle error: %v", err)
	}

	// Verify structure
	appDir := filepath.Join(agentLayerDir, "open-vscode.app")
	if _, err := os.Stat(filepath.Join(appDir, "Contents", "Info.plist")); err != nil {
		t.Fatalf("missing Info.plist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(appDir, "Contents", "MacOS", "open-vscode")); err != nil {
		t.Fatalf("missing executable: %v", err)
	}
}

func TestWriteVSCodeAppBundleMkdirError(t *testing.T) {
	root := t.TempDir()
	agentLayerDir := filepath.Join(root, ".agent-layer")
	if err := os.MkdirAll(agentLayerDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create a file where the .app directory should be
	appPath := filepath.Join(agentLayerDir, "open-vscode.app")
	if err := os.WriteFile(appPath, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if err := writeVSCodeAppBundle(agentLayerDir); err == nil {
		t.Fatalf("expected error when .app path is a file")
	}
}
