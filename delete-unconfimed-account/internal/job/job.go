package job

import (
	"context"
	"time"
)

// feito de usando o padrão scatter and gather
// wiki https://en.wikipedia.org/wiki/Gather-scatter_(vector_addressing)
func Cron(ctx context.Context, start time.Duration, delay time.Duration, f func()) {
	stream := make(chan time.Time, 1)

	//scatter
	go func() {
		// Feito apenas para executar no inicio da chamada dessa funcao
		t := <-time.After(start)
		stream <- t

		// Ticker para contar de quanto em quanto tempo a funcao do parametro vai ser executada
		ticker := time.NewTicker(delay)
		// fecha ticker no final da execucao do metodo
		// provavelmente so quando o contexto muda de estado para done
		// ou finaliza execcao da go routine main
		defer ticker.Stop()

		// escuta eventos do tick a cada <delay>
		for {
			select {
			case t2 := <-ticker.C:
				stream <- t2
			case <-ctx.Done():
				close(stream)
				return
			}
		}
	}()

	//gather
	go func() {
		//le eventos e executa f() a cada sinal recebido
		// enquanto stream nao for fechado
		for range stream {
			f()
		}
	}()
}
