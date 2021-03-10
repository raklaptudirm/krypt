#!/usr/bin/env node

'use strict'

const { pwnedPassword } = require('hibp')
const fs = require('fs')
const crypt = require('../lib/Encrypt.js')
const readlineSync = require('readline-sync')
const chalk = require('chalk')

// Global Vars

let _DATABASE; let _PASSWORDS = {}; let _KEY
let _COMMS = [{name:'edit ', use: 'Edit; Edit a password', ex: 'edit <pass_id>'}, {name: 'gent ', use: 'Get entry; Get a password entry', ex: 'gent <pass_id>'}, {name: 'gpass', use: 'Get password; Get the password of an entry', ex: 'gpass <pass_id>'}, {name: 'npass', use: 'New Password; Create a new Password', ex: 'npass'}, {name: 'sche ', use: 'Security Check; Checks your passwords\'s security', ex: 'sche'}, {name: 'cmast', use: 'Change Master; Changes the master password', ex: 'cmast'}, {name: 'dpass', use: 'Delete Password; Delete a password', ex: 'dpass <pass_id>'}, {name: 'mpass', use: 'Generate a random password', ex: 'mpass'}, {name: 'exit ', use: 'Exit; Exit the process', ex: 'exit'}]

// Main function

async function main() {
	if (fs.existsSync(__dirname + '/../lib/database.json')) {
	  loadDatabase()
	  _KEY = crypt.SHA_hash(readlineSync.question('PASSWORD: ', { hideEchoBack: true }))
	  if (_DATABASE.checksum === crypt.SHA_hash(_KEY)) {
	    console.log(chalk.green.bold('Logged in.'))
	    loadPasswords()
	    while (true) {
	      console.log('')
	      let input = readlineSync.prompt()
	      if (input === 'exit') {
	        break
	      } else if (input === 'cmast') {
	        _KEY = crypt.SHA_hash(readlineSync.questionNewPassword())
	        _DATABASE.checksum = crypt.SHA_hash(_KEY)
	        reEncryptData()
	      } else if (input === 'npass') {
	        const name_ = readlineSync.question('Password Name: ')
	        const username_ = readlineSync.question('Username: ')
	        const password_ = readlineSync.question('Password (leave empty to generate): ', { hideEchoBack: true }) || generatePassword()
	        _PASSWORDS.push(createPass(name_, username_, password_))
	        console.log(chalk.green.bold(`Sucessfully added password at ID:${_PASSWORDS.length}.`))
	        reEncryptData()
	        loadDatabase()
	        loadPasswords()
	      } else if (input.startsWith('gent')) {
	        input = parseInt(input.slice(4)) - 1
	        if (input === undefined || isNaN(input) || input < 0 || input >= _PASSWORDS.length) {
	          console.log(chalk.red.bold('ID out of bounds.'))
	        } else {
	          printPass(_PASSWORDS[input], input + 1)
	        }
	      } else if (input.startsWith('gpass')) {
	        input = parseInt(input.slice(5)) - 1
	        if (input === undefined || isNaN(input) || input < 0 || input >= _PASSWORDS.length) {
	          console.log(chalk.red.bold('ID out of bounds.'))
	        } else {
	          const sel = readlineSync.question(chalk.red('This command will show your password in cleartext. Proceed? (yes): '))
	          if (sel === 'yes') {
	            console.log(chalk.cyan.bold(_PASSWORDS[input].password))
	          } else {
	            console.log(chalk.green.bold('Command aborted.'))
	          }
	        }
	      } else if (input.startsWith('dpass')) {
	        input = parseInt(input.slice(5)) - 1
	        if (input === undefined || isNaN(input) || input < 0 || input >= _PASSWORDS.length) {
	          console.log(chalk.red.bold('ID out of bounds.'))
	        } else {
	          printPass(_PASSWORDS[input], input)
	          const sel = readlineSync.question(chalk.red.bold('Delete this entry? (yes): '))
	          if (sel === 'yes') {
	            _PASSWORDS.splice(input, 1)
	            console.log(chalk.green.bold('Password deleted Sucessfully.'))
	          } else {
	            console.log(chalk.green.bold('Delete aborted.'))
	          }
	        }
	      } else if (input === 'sche') {
	      		let weak = [], error = false
			  	for (const i in _PASSWORDS) {
			  		let strength = false, pwned
			  		if (passStrength(_PASSWORDS[i].password) !== chalk.green.bold("[STRONG]")) strength = true
			  		try {
						pwned = await pwnedPassword(_PASSWORDS[i].password)
					} catch (err) {
						console.log(chalk.red.bold("You are not connected to the internet. Krypt needs an internet connection to check if your passwords have been leaked or not."))
						error = true
						break
					}
			  		if (!error && (strength || pwned)) weak.push({id: i, pwned: pwned})
			  	}
			    if (!error) {
				    for (const i of weak) {
				    	printPass(_PASSWORDS[i.id], parseInt(i.id) + 1)
				    	if (i.pwned) console.log(chalk.red.bold(`This password has been seen ${i.pwned} times before!`))
				    	else
				    		console.log(chalk.green.bold("This password hasn't been seen before."))
				    	console.log("")
				    }
			    }
	      } else if (input === 'mpass') {
	      	const newPass = generatePassword()
	      	console.log(chalk.cyan.bold(newPass))
	      	console.log(passStrength(newPass) + chalk.red.bold(`[Occurances:${await pwnedPassword(newPass)}]`))
	      } else if (input === 'help') {
	      	console.log(chalk.cyan.bold("Available commands:"))
	      	for (const comm of _COMMS) console.log(`${chalk.green.bold(comm.name)}: ${chalk.yellow.bold(comm.use)}\nUse: ${chalk.cyan.bold(comm.ex)}\n`)
	      } else if (input.startsWith('edit')) {
	        input = parseInt(input.slice(4)) - 1
	        if (input === undefined || isNaN(input) || input < 0 || input >= _PASSWORDS.length) {
	          console.log(chalk.red.bold('ID out of bounds.'))
	        } else {
	            const name_ = readlineSync.question('Password Name (leave empty to keep same): ')
		        const username_ = readlineSync.question('Username (leave empty to keep same): ')
		        const password_ = readlineSync.question('Password (leave empty to generate): ', { hideEchoBack: true }) || generatePassword()
		        _PASSWORDS[input] = createPass(name_ || _PASSWORDS[input].name, username_ || _PASSWORDS[input].username, password_)
		        reEncryptData()
		        loadDatabase()
		        loadPasswords()
	        }
	      } else {
	        console.log(chalk.red.bold('Invalid command.'))
	      }
	    }
	  } else {
	    console.log(chalk.red.bold('Wrong Password.'))
	  }
	} else {
	  _DATABASE = { checksum: '', settings: { TwoFA: false }, data: { iv: '', encryptedData: '' } }
	  _PASSWORDS = []
	  _KEY = crypt.SHA_hash(readlineSync.questionNewPassword())
	  _DATABASE.checksum = crypt.SHA_hash(_KEY)
	  reEncryptData()
	}
}

main()

// important funcs

async function passIsPwned(password) {
	try {
		const appers = await pwnedPassword(password)
		return appers
	} catch (err) {
		console.log(chalk.red.bold("You are not connected to the internet. Krypt needs an internet connection to check if your passwords have been leaked or not."))
	}
}

function generatePassword() {
	const _lowerCase = "abcdefghijklmnopqrstuvwxyz"
	const _upperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const _numbers = "0123456789"
	const _specialChars = ",./;'[]\\=-`<>?\":|}{+_~!@#$%^&*()"

	let password = ""
	const length = 12

	for (let i = 0; i < length; i++) {
		let type = Math.round(Math.random() * 3)
		switch (type) {
			case 0:
				password += _lowerCase[Math.round(Math.random() * 25)]
				break
			case 1:
				password += _upperCase[Math.round(Math.random() * 25)]
				break
			case 2:
				password += _numbers[Math.round(Math.random() * 9)]
				break
			case 3:
				password += _specialChars[Math.round(Math.random() * _specialChars.length)]
		}
	}
	return password
}

function passStrength(pass) {
	const weak = chalk.red.bold("[WEAK]")
	const med = chalk.yellow.bold("[MEDIUM]")
	const strong = chalk.green.bold("[STRONG]")

	if (pass.length < 8) return weak

    var score = 0;
    if (!pass)
        return score;

    var letters = new Object();
    for (var i=0; i<pass.length; i++) {
        letters[pass[i]] = (letters[pass[i]] || 0) + 1;
        score += 5.0 / letters[pass[i]];
    }

    var variations = {
        digits: /\d/.test(pass),
        lower: /[a-z]/.test(pass),
        upper: /[A-Z]/.test(pass),
        nonWords: /\W/.test(pass),
    }

    var variationCount = 0;
    for (var check in variations) {
        variationCount += (variations[check] == true) ? 1 : 0;
    }
    score += (variationCount - 1) * 10;

    let strength = parseInt(score);
    if (score > 80) return strong
    if (score > 60) return med
    return weak
}

function printPass (password, id) {
  console.log(chalk.blue((`[ID:${id}]${passStrength(password.password)}`)))
  console.log('Name: ' + chalk.yellow.bold(password.name) + '\n' + 'Username: ' + chalk.yellow.bold(password.username) + '\n' + 'Password: ' + chalk.yellow.bold(new Array(password.password.length + 1).join('*')))
}

function createPass (name, username, password) {
  return { name: name, username: username, password: password }
}

function loadDatabase () {
  const data = fs.readFileSync(__dirname + '/../lib/database.json')
  _DATABASE = JSON.parse(data)
  return true
}

function loadPasswords () {
  _PASSWORDS = JSON.parse(crypt.AES_decrypt(_DATABASE.data, _KEY))
  return true
}

function reEncryptData () {
  _DATABASE.data = crypt.AES_encrypt(JSON.stringify(_PASSWORDS), _KEY)
  fs.writeFileSync(__dirname + '/../lib/database.json', JSON.stringify(_DATABASE))
  return true
}
