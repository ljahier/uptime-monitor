package main

import (
	"fmt"
	"net/http"
	"time"
)

type Resource struct {
	URL       string
	Status    string
	LastCheck time.Time
}

func (r *Resource) Check() error {
	// Code pour vérifier l'état de la ressource ici
	resp, err := http.Get(r.URL)
	// Si une erreur se produit, renvoyer l'erreur
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Si tout se passe bien, mettre à jour r.Status et r.LastCheck
	if resp.StatusCode == http.StatusOK {
		r.Status = "OK"
	} else {
		r.Status = "ERROR"
	}
	r.LastCheck = time.Now()
	return nil
}

func Monitor(resources []*Resource, interval time.Duration) {
	for {
		fmt.Println("hello")
		for _, r := range resources {
			fmt.Println("world")
			err := r.Check()
			if err != nil {
				fmt.Printf("Erreur lors de la vérification de l'état de %s : %s\n", r.URL, err.Error())
			}
		}
		time.Sleep(interval)
	}
}

func main() {
	resources := []*Resource{
		{URL: "http://www.example.com"},
		{URL: "http://www.example.net"},
	}
	Monitor(resources, 1*time.Minute)
}
