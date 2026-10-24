package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var instanceID = os.Getenv("INSTANCE_ID")
var consulAddr = os.Getenv("CONSUL_HTTP_ADDR")
var consulClient = &http.Client{Timeout: 5 * time.Second}

// Realizar el healthcheck
// Service registery
func registerWithConsul(serviceName string, port int) error {
	payload := map[string]any{
		"ID":      instanceID,
		"Name":    serviceName,
		"Address": instanceID,
		"Port":    port,
		"Check": map[string]any{
			"HTTP":                           fmt.Sprintf("http://%s:%d/heartbeat", instanceID, port),
			"Interval":                       "5s",
			"Timeout":                        "2s",
			"DeregisterCriticalServiceAfter": "1m",
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var lastErr error
	// Ciclo de hasta 10 intentos
	for attempt := 1; attempt <= 10; attempt++ {
		// Crear petición PUT para que consul registre
		req, err := http.NewRequest(http.MethodPut, consulAddr+"/v1/agent/service/register", bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		// Crea la petición
		resp, err := consulClient.Do(req)
		if err == nil {
			resp.Body.Close() // Close() libera la conexión
			if resp.StatusCode == http.StatusOK {
				log.Printf("Registrando...")
				return nil
			}
			lastErr = fmt.Errorf("Consul respondió con status %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	return fmt.Errorf("No se pudo registrar en consul: %w", lastErr)
}

func deregisterFromConsul() {
	// Identificar instancia a borrar y enviar petición
	req, err := http.NewRequest(http.MethodPut, consulAddr+"/v1/agent/service/deregister/"+instanceID, nil)
	if err != nil {
		log.Printf("No se pudo armar la petición de baja: %v", err)
		return
	}
	resp, err := consulClient.Do(req)
	if err != nil {
		log.Printf("No se pudo dar de baja %s en consul: %v", instanceID, err)
		return
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Consul respondió con status %d al dar de baja %s", resp.StatusCode, instanceID)
		return
	}
	log.Printf("Dado de baja %s en Consul", instanceID)
}

func serve(handler http.Handler, serviceName string, port int) {
	if err := registerWithConsul(serviceName, port); err != nil {
		log.Printf("WARNING: %v", err)
	}

	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}

	sigCh := make(chan os.Signal, 1)                      // Depositar señales del sistema
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM) // Mandar a este canal lo de SIGINT y SIGTERM
	done := make(chan struct{})                           // Canal sin datos = Terminó de apagarse

	go func() {
		<-sigCh
		log.Printf("Apagando %s", instanceID)
		deregisterFromConsul()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
		close(done)
	}()

	log.Printf("%s (%s) escuchando en: %d", serviceName, instanceID, port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-done
}
