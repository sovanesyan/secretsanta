# Overview 

This is a simple tool to help you out with faily picking who you get to prepare a present for. It works by reading a csv file (name, email per each line) and then sending emails to everybody letting them know who to prepare a present for. 

# Command-line params

* source (default: people.csv) - the csv file
* verbose (default: false) - whether or not to print the actual presents
* dry (default: true) - whether or not to actually send the emails
* host (default: stmp.gmail.com) - smtp host
* port (default: 587) - stmp port
* username (default: sovanesyan@gmail.com) - smtp username
* password (default: string.Empty) - stmp password

# TODO: 

* Create a template file for the email subject and body 
* Switch to net/smtp instead of the gomail package
* Should parse the csv with a csv parser instead of custom code
