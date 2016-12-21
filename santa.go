package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	gomail "gopkg.in/gomail.v2"
)

var (
	sourceFile   = flag.String("source", "people.csv", "relative path to the source file")
	verbose      = flag.Bool("verbose", false, "whether to print the presents or not")
	dryRun       = flag.Bool("dry", true, "whether or not to actually send the emails")
	smtpHost     = flag.String("host", "smtp.gmail.com", "hostname for the smtp server")
	smtpPort     = flag.Int("port", 587, "port for the smtp server")
	smtpUsername = flag.String("username", "sovanesyan@gmail.com", "username for the smtp server")
	smtpPassword = flag.String("password", "", "password for the smtp server")
)

func main() {
	flag.Parse()

	log.Printf("Welcome to Кинжала's Secret Santa %v\n", time.Now().Year())

	peoplesMap := readPeoplesMap(*sourceFile)
	people := getPeople(peoplesMap)
	log.Printf("Good, we have %v people that want to be Santa\n", len(people))

	presents := makePresents(people)
	if *verbose {
		for _, present := range presents {
			log.Printf("%10v подарява на %v\n", present.giver, present.receiver)
		}
	}

	log.Printf("All good. Looks like it is time to send the presents.")
	sendInvitations(presents, peoplesMap)
	log.Printf("Presents were sent. Now we wait for Christmas to come.")
}

type present struct {
	receiver string
	giver    string
}

func readPeoplesMap(sourceFile string) map[string]string {
	result := make(map[string]string)

	file, err := os.Open(sourceFile)
	check(err)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry := strings.Split(scanner.Text(), ",")
		result[entry[0]] = entry[1]
	}
	err = scanner.Err()
	check(err)

	return result
}

func getPeople(peoplesMap map[string]string) []string {
	people := make([]string, len(peoplesMap))

	i := 0
	for k := range peoplesMap {
		people[i] = k
		i++
	}

	return people
}

func makePresents(people []string) []present {
	givers, receivers := make([]string, len(people)), make([]string, len(people))
	copy(givers, people)
	copy(receivers, people)

	shuffle(receivers)
	for invalidPresents(givers, receivers) {
		log.Print("Somebody got to gift himself. Shuffling again.")
		shuffle(receivers)
	}
	presents := make([]present, len(people))

	for index := 0; index < len(people); index++ {
		presents[index] = present{giver: givers[index], receiver: receivers[index]}
	}

	return presents
}

func invalidPresents(givers, receivers []string) bool {
	for index := 0; index < len(givers); index++ {
		if givers[index] == receivers[index] {
			return true
		}
	}

	return false
}

func shuffle(slice []string) []string {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}

	return slice
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func sendInvitations(presents []present, peoplesMap map[string]string) {
	dialer := gomail.NewDialer(*smtpHost, *smtpPort, *smtpUsername, *smtpPassword)
	template := "Здравей %v,\n\n " +
		"Получаваш този мейл за да разбереш за кой трябва да приготвиш притеснително евтин подарък (до 10 лева).\n\n" +
		"Ти, %v подаряваш подарък на %v. Успех.\n\n" +
		"Поздрави, \nСерж\n\n." +
		"Произведено със SecretSanta Maker 4000."

	for _, present := range presents {
		message := gomail.NewMessage()
		message.SetHeader("To", peoplesMap[present.giver])
		message.SetHeader("From", "sovanesyan@gmail.com")
		message.SetHeader("Subject", "Твой ред е да си дядо Коледа.")
		message.SetBody("text/plain", fmt.Sprintf(template, present.giver, present.giver, present.receiver))
		if *verbose {
			log.Printf("Prepared email for %v\n", present.giver)
		}
		if *dryRun == false {
			dialer.DialAndSend(message)
			log.Printf("Sent email to %v\n", present.giver)
		}

	}
}
