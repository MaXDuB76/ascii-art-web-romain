package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type IndexData struct {
	Result string
	Text   string
	Banner string
}

var (
	data IndexData
	err  error
	tmpl *template.Template
)

func main() {
	http.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 Page Not Found")
	})
	/* La méthode http.HandleFunc() est utilisée pour définir une fonction de rappel qui sera appelée lorsque le serveur reçoit une requête HTTP.
	   La fonction de rappel est définie avec deux paramètres, w (de type http.ResponseWriter) et r (de type *http.Request).
	   Elle est appelée à chaque fois qu'une requête est reçue. Si la méthode de la requête est "GET", cela signifie qu'il s'agit d'une nouvelle demande de page.
	   Dans ce cas, la structure data est réinitialisée et la fonction GetAndExecute() est appelée pour traiter la requête.
	   Si la méthode de la requête est "POST", cela signifie qu'un formulaire a été soumis.
	   Dans ce cas, la fonction generateAsciiArt() est appelée pour générer l'art ASCII en fonction des données de la requête,
	   et le résultat est stocké dans la structure data. Ensuite, la fonction GetAndExecute() est appelée pour traiter la requête. */
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			data = IndexData{}
			GetAndExecute(w)

		} else if r.Method == "POST" {
			// Generate ASCII art based on the text and banner selection
			asciiArt := generateAsciiArt(r)
			data = IndexData{
				Result: asciiArt,
			}
			GetAndExecute(w)
		} else {
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}
	})
	LaunchServer()
}

func generateAsciiArt(r *http.Request) string { // Fonction qui génère l'ascii art en fonction de la requête de l'utilisateur sur le site.
	var outputstr string
	entree := r.FormValue("text")                           // On récupère ce que l'utilisateur a écrit dans le champ de saisie de texte du site.
	banner := "templates/" + r.FormValue("banner") + ".txt" // On récupère la bannière choisie par l'utilisateur dans le menu déroulant du site.
	checkvide := false
	checksautligne := false
	mot := Split(entree, "\\n") // Si, dans la requête de l'utilisateur on trouve la chaîne "\\n", alors nous découpons sa requête en 2 requêtes distinctes.
	if entree == "\\n" {
		outputstr += "\n"
	} else {
		for _, arg := range mot {
			checkvide = false
			checksautligne = false // On imprime les mots uns par un si on a trouvé un "\n" dans la requête.
			for j := 0; j < 8; j++ {
				if len(entree) == 0 { // Si le mot est vide, on imprime rien
					checkvide = true
					break
				}
				if len(arg) == 0 {
					checksautligne = true
					break
				}
				if !checkvide {
					for i := 0; i < len(arg); i++ { // Cette boucle nous permet de faire la jonction entre chaque caractère des mots de la requête et
						carac := int(arg[i]) // l'emplacement de ces caractères dans les fichiers .txt où se trouve les ascii arts.
						startLine := 2 + (carac-32)*9 + j
						outputstr = outputstr + printLines(banner, startLine)
					}
				}
				if !checkvide {
					outputstr += "\n"
				}
			}
			if checksautligne {
				outputstr += "\n"
			}
		}
	}
	/* os.WriteFile("output.txt", []byte(outputstr), 0o644) */
	return outputstr
}

func printLines(fileName string, startLine int) string { // Cette fonction nous permet d'imprimer une ligne spécifique d'un fichier tiers.
	var result string
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		if lineNum == startLine {
			result = result + scanner.Text()
			break
		}
		lineNum++
	}
	if scanner.Err() != nil {
		panic(err)
	}
	return result
}

func Split(s, sep string) []string { // Cette fonction nous permet de séparer la requête de l'utilisateur en plusieurs mots
	var answer []string // si on trouve l'index "\n" (fonction suivante), afin de les imprimer les uns après les autres.
	a := 0
	b := 0
	lens := len(sep)
	for i := 0; i < len(s); i++ {
		b = Index(s[a:], sep)
		if b < 0 {
			break
		}
		answer = append(answer, s[a:a+b])
		a = a + b + lens
	}
	answer = append(answer, s[a:])

	return answer
}

func Index(s string, toFind string) int {
	a := len([]rune(s))
	b := len([]rune(toFind))
	for i := 0; i < a-b+1; i++ {
		if s[i:i+b] == toFind {
			return i
		}
	}
	return -1
}

func GetAndExecute(w http.ResponseWriter) { // Fonction qui permet de traiter la requête
	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func LaunchServer() { // Fonction qui permet de démarrer le serveur avec un lien vers la page dans le terminal.
	port := ":8081"
	ip := getLocalIP()
	link := "http://" + ip + port
	log.Println("Starting server on : ", link)
	log.Fatal(http.ListenAndServe(port, nil))
}

func getLocalIP() string { // Fonction qui retourne l'adresse IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, address := range addrs {
		// check if the address is an IP address (IPv4 or IPv6)
		ipNet, ok := address.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				// return the IPv4 address without the subnet mask
				return strings.Split(ipNet.IP.String(), "/")[0]
			}
		}
	}
	return "localhost"
}
