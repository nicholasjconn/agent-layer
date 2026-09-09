package wizard

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"github.com/conn-castle/agent-layer/internal/agentoptions"
	"github.com/conn-castle/agent-layer/internal/messages"
)

func TestMain(m *testing.M) {
	wizardOptionDiscoveryRequestFunc = func() agentoptions.DiscoveryRequest { return agentoptions.DiscoveryRequest{} }
	os.Exit(m.Run())
}

func TestModelDiscoveryPrefetchWaitsAndReusesResult(t *testing.T) {
	original := wizardOptionDiscoveryRequestFunc
	t.Cleanup(func() { wizardOptionDiscoveryRequestFunc = original })
	path := filepath.Join(t.TempDir(), "agy")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nprintf 'new-model\\tNew Model\\n'\n"), 0o700); err != nil { // #nosec G306 -- executable harness fixture.
		t.Fatal(err)
	}
	started, release := make(chan struct{}), make(chan struct{})
	var lookups atomic.Int32
	wizardOptionDiscoveryRequestFunc = func() agentoptions.DiscoveryRequest {
		return agentoptions.DiscoveryRequest{Live: true, LookPath: func(string) (string, error) {
			if lookups.Add(1) == 1 {
				close(started)
			}
			<-release
			return path, nil
		}}
	}
	var out bytes.Buffer
	cache := &wizardOptionDiscoveryCache{out: &out}
	cache.prefetch(AgentAntigravity)
	<-started
	// A pending discovery must be awaited rather than replaced by stale options.
	timer := time.AfterFunc(30*time.Millisecond, func() { close(release) })
	defer timer.Stop()
	ui := &MockUI{SelectFunc: func(_ string, options []string, current *string) error {
		if !slices.Contains(options, "New Model") {
			t.Errorf("missing discovered model: %v", options)
		}
		*current = "New Model"
		return nil
	}}
	var selected string
	for range 2 {
		if err := cache.selectModel(ui, AgentAntigravity, messages.WizardAntigravityModelTitle, &selected); err != nil {
			t.Fatal(err)
		}
	}
	if lookups.Load() != 1 {
		t.Fatalf("discovery ran %d times", lookups.Load())
	}
	if !bytes.Contains(out.Bytes(), []byte("Discovering antigravity models")) {
		t.Fatal("missing progress message")
	}
}

func TestModelDiscoveryFailureStopsPickerWithoutChangingConfiguredValue(t *testing.T) {
	original := wizardOptionDiscoveryRequestFunc
	t.Cleanup(func() { wizardOptionDiscoveryRequestFunc = original })
	wizardOptionDiscoveryRequestFunc = func() agentoptions.DiscoveryRequest {
		return agentoptions.DiscoveryRequest{Live: true, LookPath: func(string) (string, error) { return "", errors.New("harness missing") }}
	}
	ui := &MockUI{
		SelectFunc: func(_ string, options []string, current *string) error {
			t.Error("failed discovery opened a picker")
			return nil
		},
	}
	cache := &wizardOptionDiscoveryCache{}
	selected := "my-custom-model"
	for range 2 {
		if err := cache.selectModel(ui, AgentClaude, messages.WizardClaudeModelTitle, &selected); err == nil || !bytes.Contains([]byte(err.Error()), []byte("harness missing")) {
			t.Fatalf("expected actionable discovery failure, got %v", err)
		}
	}
	if selected != "my-custom-model" {
		t.Fatalf("selected=%q", selected)
	}
}

func TestPrefetchEnabledModelsRunsConcurrently(t *testing.T) {
	original := wizardOptionDiscoveryRequestFunc
	t.Cleanup(func() { wizardOptionDiscoveryRequestFunc = original })
	started := make(chan string, 4)
	release := make(chan struct{})
	wizardOptionDiscoveryRequestFunc = func() agentoptions.DiscoveryRequest {
		return agentoptions.DiscoveryRequest{Live: true, LookPath: func(name string) (string, error) { started <- name; <-release; return "", os.ErrNotExist }}
	}
	choices := NewChoices()
	for _, agent := range []string{AgentAntigravity, AgentClaude, AgentCodex, AgentGrok} {
		choices.EnabledAgents[agent] = true
	}
	cache := &wizardOptionDiscoveryCache{}
	cache.prefetchEnabled(choices)
	for range 4 {
		select {
		case <-started:
		case <-time.After(time.Second):
			close(release)
			t.Fatal("discovery did not start concurrently")
		}
	}
	close(release)
	for _, entry := range cache.entries {
		<-entry.done
	}
}

func TestScriptedModelSelectionDoesNotDiscover(t *testing.T) {
	original := wizardOptionDiscoveryRequestFunc
	t.Cleanup(func() { wizardOptionDiscoveryRequestFunc = original })
	wizardOptionDiscoveryRequestFunc = func() agentoptions.DiscoveryRequest {
		t.Fatal("explicit scripted model caused discovery")
		return agentoptions.DiscoveryRequest{}
	}
	ui := &ScriptedUI{answers: scriptedAnswers{Select: map[string]string{messages.WizardAntigravityModelTitle: "future-native-model"}}, used: map[string]struct{}{}}
	var value string
	cache := &wizardOptionDiscoveryCache{}
	if err := cache.selectModel(ui, AgentAntigravity, messages.WizardAntigravityModelTitle, &value); err != nil {
		t.Fatal(err)
	}
	if value != "future-native-model" {
		t.Fatalf("configured value = %q", value)
	}
	if err := ui.AssertComplete(); err != nil {
		t.Fatal(err)
	}
}
