/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stress",
	Short: "Sistema CLI em Go para realizar testes de carga em um serviço web",
	Long: `Sistema CLI em Go para realizar testes de carga em um serviço web.
	
O usuário deverá fornecer a URL do serviço, o número total de requests e a quantidade de chamadas simultâneas.
Serviço irá retornar um relatório com tempo e as responses obtidas.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		url, err := cmd.Flags().GetString("url")
		if err != nil || !isUrl(url) {
			println("Erro ao obter a URL ou URL inválida")
			os.Exit(1)
		}

		requests, err := cmd.Flags().GetInt64("requests")
		if err != nil || requests < 1 {
			println("Erro ao obter o número de requests ou valor inválido")
			os.Exit(1)
		}

		concurrency, err := cmd.Flags().GetInt64("concurrency")
		if err != nil || concurrency < 1 {
			println("Erro ao obter o número de chamadas simultâneas ou valor inválido")
			os.Exit(1)
		}

		startTime := time.Now()
		c := make(chan int)
		go stressTest(url, int(requests), int(concurrency), c)
		counter := 0
		for v := range c {
			report[v]++
			counter++
			if counter == int(requests) {
				break
			}
		}

		println(fmt.Sprintf("Tempo total gasto na execução: %0.2f Sec", time.Since(startTime).Seconds()))
		println(fmt.Sprintf("Quantidade total de requests realizados: %d requests", requests))
		println("Resultados:")
		println(fmt.Sprintf("Quantidade de requests com status HTTP [200]: %d", report[200]))

		success := 0
		for k, v := range report {
			if k == 200 {
				success += v
				continue
			}

			if k > 200 && k < 300 {
				success += v
			}

			if k == 0 {
				println(fmt.Sprintf("Quantidade de requests com erro sem status code definido: %d", v))
				continue
			}

			println(fmt.Sprintf("Quantidade de requests com status HTTP [%d]: %d", k, v))
		}

		close(c)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var report map[int]int

func init() {
	rootCmd.Flags().StringP("url", "u", "", "URL do serviço a ser testado ex. --url=http://google.com")
	rootCmd.Flags().Int64P("requests", "r", 1, "Número total de requests. --requests=1000")
	rootCmd.Flags().Int64P("concurrency", "c", 1, "Número de chamadas simultâneas. --concurrency=10")

	report = make(map[int]int)
}

// isUrl - checks if a string is a valid URL
func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func stressTest(url string, requests int, concurrency int, c chan int) {
	httpClient := &http.Client{}
	for i := 0; i < concurrency; i++ {
		qtdRequests := int(requests / concurrency)
		if (concurrency-1) == i && (requests%concurrency) != 0 {
			qtdRequests += requests % concurrency
		}

		go makeRequests(httpClient, qtdRequests, url, c)
	}
}

func makeRequests(httpClient *http.Client, qtdRequests int, url string, c chan int) {
	for j := 0; j < qtdRequests; j++ {
		resp, err := httpClient.Get(url)
		if err != nil {
			c <- int(0)
			continue
		}

		c <- int(resp.StatusCode)
	}
}
