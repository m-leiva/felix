package scratchpad

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestScratchpadLifecycle(t *testing.T) {
	db, err := sql.Open("sqlite", t.TempDir()+"/scratchpad.db")
	if err != nil {
		t.Fatalf("error abriendo sqlite: %v", err)
	}
	defer db.Close()

	mgr, err := NewManager(db)
	if err != nil {
		t.Fatalf("NewManager fallo: %v", err)
	}

	item, err := mgr.Get("sess-1")
	if err != nil || item != nil {
		t.Fatalf("se esperaba nil en sesion inexistente, obtenido item=%+v err=%v", item, err)
	}

	if err := mgr.Set("sess-1", "Desplegar servicio", "in_progress", map[string]string{"env": "prod"}); err != nil {
		t.Fatalf("Set fallo: %v", err)
	}

	item, err = mgr.Get("sess-1")
	if err != nil || item == nil {
		t.Fatalf("Get fallo tras Set: item=%+v err=%v", item, err)
	}
	if item.Goal != "Desplegar servicio" || item.Status != "in_progress" || item.Variables["env"] != "prod" {
		t.Errorf("datos inesperados en scratchpad: %+v", item)
	}

	if err := mgr.Set("sess-1", "Desplegar servicio", "done", map[string]string{"env": "prod"}); err != nil {
		t.Fatalf("Set update fallo: %v", err)
	}
	item, _ = mgr.Get("sess-1")
	if item.Status != "done" {
		t.Errorf("status esperado 'done', obtenido %q", item.Status)
	}

	if err := mgr.Clear("sess-1"); err != nil {
		t.Fatalf("Clear fallo: %v", err)
	}
	item, err = mgr.Get("sess-1")
	if err != nil || item != nil {
		t.Fatalf("se esperaba nil tras Clear, obtenido item=%+v err=%v", item, err)
	}
}
