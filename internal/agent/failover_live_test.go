//go:build live

package agent

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"felix-go/internal/config"
	"felix-go/internal/tools"
)

func TestLiveFailoverWithToolCalls(t *testing.T) {
	raw, err := os.ReadFile("../../config.json")
	if err != nil {
		t.Fatalf("no se pudo leer config.json: %v", err)
	}
	var cfg config.Config
	if err := json.Unmarshal(raw, &cfg); err != nil {
		t.Fatalf("error parseando config.json: %v", err)
	}
	if cfg.FallbackAPIBase == "" || cfg.FallbackAPIKey == "" || cfg.FallbackModel == "" {
		t.Skip("fallback no configurado en config.json")
	}

	cfg.APIBase = "http://127.0.0.1:1/v1"
	cfg.APIKey = "dead-key"
	cfg.MaxIterations = 5

	reg := tools.NewRegistry()
	reg.Register(&mockEchoTool{name: "status_checker"})

	engine := NewEngine(&cfg, reg, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	reply, err := engine.orchestrator.Run(ctx, "live-failover-test", "Ejecuta status_checker y dime exactamente que devolvio.", "")
	if err != nil {
		t.Fatalf("fallo el turno en failover en vivo: %v", err)
	}
	if reply == "" {
		t.Fatalf("respuesta vacia en failover en vivo")
	}
	st := engine.GetLLMStatus()
	if st.TotalFallbackCalls == 0 {
		t.Errorf("se esperaba al menos 1 llamada al proveedor fallback, obtenido 0")
	}
	t.Logf("Failover en vivo OK (fallback_calls=%d, model=%s): %s", st.TotalFallbackCalls, cfg.FallbackModel, reply)
}
