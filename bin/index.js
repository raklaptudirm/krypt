#!/usr/bin/env node

'use strict'

const { pwnedPassword } = require('hibp')
const fs = require('fs')
const crypt = require('../lib/Encrypt.js')
const readlineSync = require('readline-sync')
const chalk = require('chalk')
const clipboardy = require('clipboardy')

//constants

const databaseTemplate = { checksum: { checksum: "", salt: "" }, salt: { key: "", TwoFA: "" }, settings: { TwoFA: { on: false, question: "", answer: { checksum: "", salt: "" } }, hint: { on: false, hint: "" }, alias: {  }, passwordWordy: false }, data: { iv: '', encryptedData: '' } },
_COMMS = [{name:'edit', use: 'Edit; Edit a password', ex: 'edit <pass_id>'}, {name: 'gent', use: 'Get entry; Get a password entry', ex: 'gent <pass_id>'}, {name: 'gpass', use: 'Get password; Get the password of an entry', ex: 'gpass <pass_id>'}, {name: 'npass', use: 'New Password; Create a new Password', ex: 'npass'}, {name: 'sche', use: 'Security Check; Checks your passwords\'s security', ex: 'sche'}, {name: 'cmast', use: 'Change Master; Changes the master password', ex: 'cmast'}, {name: 'dpass', use: 'Delete Password; Delete a password', ex: 'dpass <pass_id>'}, {name: 'mpass', use: 'Make password; Generate a random password', ex: 'mpass'}, {name: 'exit ', use: 'Exit; Exit the process', ex: 'exit'}, {name: 'list', use: 'List; List all passwords', ex: 'list'}, {name: 'search', use: 'Search; Search for a password', ex: 'search <keyword>'}, {name: 'set', use: 'Set; Change a setting', ex: 'set <setting> <args>'}, {name: 'copy', use: 'Copy; Copy a password', ex: 'copy <pass_id>'}],
_BASENAME = /[A-Za-z0-9-_.,]{1,100}/

// Global Vars

let _DATABASE, _PASSWORDS = {}, _KEY, _2F, _NAME, _WORDS

// Main function

async function main() {
	if (fs.existsSync(__dirname + '/../databases/' + _NAME)) {

		if(!loadDatabase()) return
		_KEY = crypt.PBKDF2_HASH(readlineSync.question('PASSWORD: ', { hideEchoBack: true }), _DATABASE.salt.key)

		if (_DATABASE.settings.passwordWordy)
			_WORDS = JSON.parse(fs.readFileSync(__dirname + '/../lib/words.json'))

		if (_DATABASE.settings.TwoFA.on) _2F = crypt.PBKDF2_HASH(readlineSync.question(_DATABASE.settings.TwoFA.question + '? ', { hideEchoBack: true }), _DATABASE.salt.TwoFA)

		if (_DATABASE.checksum.checksum === crypt.PBKDF2_HASH(_KEY, _DATABASE.checksum.salt) && (!_DATABASE.settings.TwoFA.on || _DATABASE.settings.TwoFA.answer.checksum === crypt.PBKDF2_HASH(_2F, _DATABASE.settings.TwoFA.answer.salt))) {

			console.log(chalk.green.bold('Logged in.'))
			loadPasswords()
			while (true) {
				console.log('')
				let input = readlineSync.prompt().toLowerCase()
				if (_DATABASE.settings.alias[input] !== undefined)
					input = _DATABASE.settings.alias[input]
				input = input.split(' ')

				if (input[0] === 'exit') {
					console.clear()
					break
				} else if (input[0] === 'cmast') {
					_KEY = crypt.PBKDF2_HASH(readlineSync.questionNewPassword("Enter new Password: ", { min: 8 }))
					_DATABASE.salt.key = _KEY.salt
					_KEY = _KEY.checksum
					_DATABASE.checksum = crypt.PBKDF2_HASH(_KEY)
					reEncryptData()
				} else if (input[0] === 'npass') {
					const name_ = readlineSync.question('Password Name: ')
					const username_ = readlineSync.question('Username: ')
					const password_ = readlineSync.question('Password (leave empty to generate): ', { hideEchoBack: true }) || generatePassword()
					_PASSWORDS.push(createPass(name_, username_, password_))
					console.log(chalk.green.bold(`Sucessfully added password at ID:${_PASSWORDS.length}.`))
					reEncryptData()
				} else if (input[0] === 'gent') {
					input = parseInt(input[1]) - 1
					if (input === undefined || isNaN(input) || input < 0 || input >= _PASSWORDS.length) {
						console.log(chalk.red.bold('ID out of bounds.'))
					} else {
						printPass(_PASSWORDS[input], input + 1)
					}
				} else if (input[0] === 'gpass') {
					input = parseInt(input[1]) - 1
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
				} else if (input[1] === 'dpass') {
					input = parseInt(input[1]) - 1
					if (input === undefined || isNaN(input) || input < 0 || input >= _PASSWORDS.length) {
						console.log(chalk.red.bold('ID out of bounds.'))
					} else {
						printPass(_PASSWORDS[input], input)
						const sel = readlineSync.question(chalk.red.bold('Delete this entry? (yes): '))
						if (sel === 'yes') {
							_PASSWORDS.splice(input, 1)
							console.log(chalk.green.bold('Password deleted Sucessfully.'))
							reEncryptData()
						} else {
							console.log(chalk.green.bold('Delete aborted.'))
						}
					}
				} else if (input[0] === 'sche') {
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
						if (weak.length !== 0) {
							for (const i of weak) {
								printPass(_PASSWORDS[i.id], parseInt(i.id) + 1)
								if (i.pwned) console.log(chalk.red.bold(`This password has been seen ${i.pwned} times before!`))
								else
									console.log(chalk.green.bold("This password hasn't been seen before."))
									console.log("")
							}
						} else {
							console.log(chalk.green.bold("All your passwords are secure."))
						}
					}
				} else if (input[0] === 'mpass') {
					const newPass = generatePassword()
					console.log(chalk.cyan.bold(newPass))
					console.log(passStrength(newPass) + chalk.red.bold(`[Occurances:${await pwnedPassword(newPass)}]`))
				} else if (input[0] === 'help') {
					console.log(chalk.cyan.bold("Available commands:"))
					for (const comm of _COMMS) console.log(`${chalk.green.bold(comm.name)}: ${chalk.yellow.bold(comm.use)}\nUse: ${chalk.cyan.bold(comm.ex)}\n`)
				} else if (input[0] === 'edit') {
					input = parseInt(input[1]) - 1
					if (input === undefined || isNaN(input) || input < 0 || input >= _PASSWORDS.length) {
						console.log(chalk.red.bold('ID out of bounds.'))
					} else {
						const name_ = readlineSync.question('Password Name (leave empty to keep same): ')
						const username_ = readlineSync.question('Username (leave empty to keep same): ')
						const password_ = readlineSync.question('Password (leave empty to generate): ', { hideEchoBack: true }) || generatePassword()
						_PASSWORDS[input] = createPass(name_ || _PASSWORDS[input].name, username_ || _PASSWORDS[input].username, password_)
						console.log(chalk.green.bold("Sucessfully edited password."))
						reEncryptData()
					}
				} else if (input[0] === 'list') {
					for (const i in _PASSWORDS) {
						console.log("")
						printPass(_PASSWORDS[i], parseInt(i) + 1)
					}
					if (_PASSWORDS.length === 0) console.log(chalk.red.bold("You do not have any stored passwords."))
				} else if (input[0] === 'search') {
					let notFound = true
					for (const i in _PASSWORDS) {
						if (_PASSWORDS[i].name.toLowerCase().includes(input[1])) {
							printPass(_PASSWORDS[i], parseInt(i) + 1)
							notFound = false
							console.log("")
						}
					}
					if (notFound) console.log(chalk.red.bold("No matches found."))
				} else if (input[0] === 'set') {
					if (input[1] === '2fa') {
						if (!_DATABASE.settings.TwoFA.on) {
							const sel = readlineSync.question(chalk.green.bold('Enable 2-Factor Authentication? (yes): '))
							if (sel === 'yes') {
								_DATABASE.settings.TwoFA.on = true
								_DATABASE.settings.TwoFA.question = readlineSync.question('Enter a question: ')
								_2F = crypt.PBKDF2_HASH(readlineSync.question('Enter the answer: '))
								_DATABASE.salt.TwoFA = _2F.salt
								_2F = _2F.checksum
								_DATABASE.settings.TwoFA.answer = crypt.PBKDF2_HASH(_2F)
								console.log(chalk.green.bold("Enabled 2 factor Auth."))
								reEncryptData()					
							} else {
								console.log(chalk.red.bold('Command aborted.'))
							}
						} else {
							if (input[2] === 'dis') {
								const sel = readlineSync.question(chalk.red.bold('Disable 2-Factor Authentication? (yes): '))
								if (sel === 'yes') {
									_DATABASE.settings.TwoFA.on = false
									console.log(chalk.green.bold("Disabled 2 factor Auth."))
									reEncryptData()
								} else {
									console.log(chalk.green.bold('Command aborted.'))
								}
							} else {
								_DATABASE.settings.TwoFA.question = readlineSync.question('Enter new question (Keep empty to keep the same): ') || _DATABASE.settings.TwoFA.question
								_2F = crypt.PBKDF2_HASH(readlineSync.question('Enter new answer: '))
								_DATABASE.salt.TwoFA = _2F.salt
								_2F = _2F.checksum
								_DATABASE.settings.TwoFA.answer = crypt.PBKDF2_HASH(_2F)
								console.log(chalk.green.bold("Changed Auth factors."))
								reEncryptData()
							}
						}
					} else if (input[1] === 'hint') {
						if (!_DATABASE.settings.hint.on) {
							const sel = readlineSync.question(chalk.green.bold('Enable Hint? (yes): '))
							if (sel === 'yes') {
								_DATABASE.settings.hint.on = true
								_DATABASE.settings.hint.hint = readlineSync.question('Enter a hint: ')
								console.log(chalk.green.bold("Enabled Hint."))
								reEncryptData()
								loadDatabase()
								loadPasswords()
							} else {
								console.log(chalk.red.bold('Command aborted.'))
							}
						} else {
							if (input[2] === 'dis') {
								const sel = readlineSync.question(chalk.red.bold('Disable Hint? (yes): '))
								if (sel === 'yes') {
									_DATABASE.settings.hint.on = false
									console.log(chalk.green.bold("Disabled Hint."))
									reEncryptData()
									loadDatabase()
									loadPasswords()
								} else {
									console.log(chalk.green.bold('Command aborted.'))
								}
							} else {
								_DATABASE.settings.hint.hint = readlineSync.question('Enter new hint (Keep empty to keep the same):') || _DATABASE.settings.hint.hint
								console.log(chalk.green.bold("Changed Hint."))
								reEncryptData()
							}
						}
					} else if (input[1] === 'alias') {
						if (input[2] === 'new') {
							const alias = readlineSync.question("Enter alias name: ")
							if (_DATABASE.settings.alias[alias] === undefined && isNotCommand(alias)) {
								_DATABASE.settings.alias[alias] = readlineSync.question("Enter alias command: ")
								console.log(chalk.green.bold("Alias set."))
								reEncryptData()
								loadDatabase()
								loadPasswords()
							} else {
								console.log(chalk.red.bold("Command already exists."))
							}
						} else if (input[2] === 'rename') {
							if (_DATABASE.settings.alias[input[3]] === undefined) {
								console.log(chalk.red.bold("Could not find alias."))
							} else {
								const alias = readlineSync.question("Enter new name: ")
								_DATABASE.settings.alias[alias] = _DATABASE.settings.alias[input[3]]
								delete _DATABASE.settings.alias[input[3]]
								console.log(chalk.green.bold("Alias renamed sucessfully."))
							}
						} else if (input[2] === 'delete') {
							if (_DATABASE.settings.alias[input[3]] === undefined) {
								console.log(chalk.red.bold("Could not find alias."))
							} else {
								delete _DATABASE.settings.alias[input[3]]
								console.log(chalk.green.bold("Alias deleted sucessfully."))
							}
						} else if (input[2] === 'list') {
							console.log("")
							for (const i in _DATABASE.settings.alias) {
								console.log(`${chalk.blue.bold(i)}: ${chalk.yellow.bold(_DATABASE.settings.alias[i])}`)
							}
							if (_DATABASE.settings.alias.length === 0) console.log(chalk.red.bold("You do not have any stored aliases."))
						} else {
							console.log(chalk.red.bold("Undefined argument."))
						}
					} else {
						console.log(chalk.red.bold("Setting not found."))
					}
				} else if (input[0] === 'copy') {
					input = parseInt(input[1]) - 1
					if (input === undefined || isNaN(input) || input < 0 || input >= _PASSWORDS.length) {
						console.log(chalk.red.bold('ID out of bounds.'))
					} else {
						clipboardy.writeSync(_PASSWORDS[input].password)
						console.log(chalk.green.bold("Password copied to clipboard."))
					}
				} else if (input[1] === 'password') {
					_DATABASE.settings.passwordWordy = !_DATABASE.settings.passwordWordy
					if (_DATABASE.settings.passwordWordy)
						console.log(chalk.green.bold("Enabled wordy password."))
				    else
				    	console.log(chalk.green.bold("Disabled wordy password."))
				} else {
					console.log(chalk.red.bold('Invalid command.'))
				}
			}
		} else {
			console.log(chalk.red.bold((_DATABASE.settings.TwoFA.on)?'Wrong Password or 2nd factor.':'Wrong Password.'))
			if (_DATABASE.settings.hint.on)
				console.log(chalk.green.bold(`Hint: ${_DATABASE.settings.hint.hint}`))
		}
	} else {
		_DATABASE = databaseTemplate
		_PASSWORDS = []
		_KEY = crypt.PBKDF2_HASH(readlineSync.questionNewPassword('New Password: ', { min: 8 }))
		_DATABASE.salt.key = _KEY.salt
		_KEY = _KEY.checksum
		_DATABASE.checksum = crypt.PBKDF2_HASH(_KEY)
		reEncryptData()
	}
}
	
// important funcs

function isNotCommand(name) {
	for (const comm of _COMMS) {
		if (comm.name === name) {
			return false
		}
	}
	return true
}

const gatherKeys = obj => {
    const
    isObject = val => typeof val === 'object' && !Array.isArray(val),
    addDelimiter = (a, b) => a ? `${a}.${b}` : b,
    paths = (obj = {}, head = '') => Object.entries(obj)
      .reduce((product, [key, value]) =>
        (fullPath => isObject(value)
          ? product.concat(paths(value, fullPath))
          : product.concat(fullPath))
        (addDelimiter(head, key)), []);
    return new Set(paths(obj));
};

const
  diff = (a, b) => new Set(Array.from(a).filter(item => !b.has(item))),
  equalByKeys = (a, b) => diff(gatherKeys(a), gatherKeys(b)).size === 0;

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
	let password

	if (_DATABASE.settings.passwordWordy) {
		let seperator = _specialChars[crypt.random(_specialChars.length - 1)]
		let len = _WORDS.length - 1, front = '', back = ''
		if (crypt.random(1) === 0)
			front = crypt.random(999)
		else
			back = crypt.random(999)
		password = front + _WORDS[crypt.random(len)] + seperator + _WORDS[crypt.random(len)] + seperator + _WORDS[crypt.random(len)] + back
	} else {
		do {
			password = ""
			const length = 12

			for (let i = 0; i < length; i++) {
				let type = crypt.random(3)
				switch (type) {
					case 0:
					password += _lowerCase[crypt.random(25)]
					break
					case 1:
					password += _upperCase[crypt.random(25)]
					break
					case 2:
					password += _numbers[crypt.random(9)]
					break
					case 3:
					password += _specialChars[crypt.random(_specialChars.length - 1)]
				}
			}
		} while (passStrength(password) !== chalk.green.bold("[STRONG]"))
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
	const data = fs.readFileSync(__dirname + '/../databases/' + _NAME)
	try {
		_DATABASE = JSON.parse(data)
		if (equalByKeys(databaseTemplate, _DATABASE))
			return true
		console.log(chalk.red.bold("[FATAL] The database has been corrupted. Invalid Property List"))
		return false
	} catch (err) {
		console.log(chalk.red.bold("[FATAL] The database has been corrupted. Invalid JSON. ") + err)
		return false
	}
}

function loadPasswords () {
	_PASSWORDS = JSON.parse(decryptPass())
}

function reEncryptData () {
	if (_DATABASE.settings.TwoFA.on)
		_DATABASE.data = crypt.AES_encrypt(JSON.stringify(crypt.AES_encrypt(JSON.stringify(_PASSWORDS), _KEY)), _2F)
	else
		_DATABASE.data = crypt.AES_encrypt(JSON.stringify(_PASSWORDS), _KEY)
	fs.writeFileSync(__dirname + '/../databases/' + _NAME, JSON.stringify(_DATABASE))
}

function decryptPass () {
	if (_DATABASE.settings.TwoFA.on) 
		return crypt.AES_decrypt(JSON.parse(crypt.AES_decrypt(_DATABASE.data, _2F)), _KEY)
	return crypt.AES_decrypt(_DATABASE.data, _KEY)
}

function getDatabases () {
	const data = fs.readFileSync(__dirname + '/../config.json')
	try {
		const config = JSON.parse(data)
		return config
	} catch (err) {
		console.log(chalk.red.bold("[FATAL] The database list has been corrupted. Invalid JSON. ") + err)
		return false
	}
}

// main process

if (process.argv.length === 2) {
	if (!fs.existsSync(__dirname + '/../config.json'))
		fs.writeFileSync(__dirname + '/../config.json', '{"selected": "default","databases": ["default"]}')
	_NAME = getDatabases()
	_NAME = _NAME.selected + ".json"
	main()
} else {
	let args = process.argv.slice(2)
	if (args[0] === 'new') {
		let config = getDatabases()
		const newName = readlineSync.question("Enter database name: ")
		const matches = _BASENAME.exec(newName)
		if (matches && matches[0] === newName && newName.length !== 0) {
			if (!config.databases.includes(newName)) {
				config.databases.push(newName)
				fs.writeFileSync(__dirname + '/../config.json', JSON.stringify(config))
				console.log(chalk.green.bold("Added new database."))
			} else {
				console.log(chalk.red.bold("Database already exists."))
			}
		} else {
			console.log(chalk.red.bold("Illeagal database name. Database names can contain characters A-Z, a-z, 0-9, -, _, ., and ,. They also should have a length between 1-100."))
		}
	} else if (args[0] === 'list') {
		let config = getDatabases()
		for (const databaseName of config.databases) {
			console.log(chalk.blue.bold(databaseName))
		}
	} else if (args[0] === 'switch') {
		let config = getDatabases()
		if (config.databases.includes(args[1])) {
			config.selected = args[1]
			fs.writeFileSync(__dirname + '/../config.json', JSON.stringify(config))
			console.log(chalk.green.bold(`Switched to ${args[1]} database.`))
		} else {
			console.log(chalk.red.bold("Database not found."))
		}
	} else if (args[0] === 'delete') {
		let config = getDatabases()
		if (config.databases.includes(args[1])) {
			if (config.databases.length === 1) {
				console.log(chalk.red.bold("Can't delete last database."))
			} else {
				if (readlineSync.question(chalk.red.bold(`Delete the ${args[1]} database? (yes): `)) === 'yes') {
					if (fs.existsSync(__dirname + '/../databases/' + args[1] + ".json"))
						fs.unlinkSync(__dirname + '/../databases/' + args[1] + ".json")
					config.databases.splice(config.databases.indexOf(args[1]), 1)
					if (config.selected === args[1]) {
						config.selected = config.databases[0]
					}
					console.log(chalk.green.bold(`Deleted ${args[1]} database.`))
					fs.writeFileSync(__dirname + '/../config.json', JSON.stringify(config))
				} else {
					console.log(chalk.green.bold("Delete aborted."))
				}
			}
		} else {
			console.log(chalk.red.bold("Database not found."))
		}
	} else if (args[0] === 'rename') {
		let config = getDatabases()
		if (config.databases.includes(args[1])) {
			const newDBName = readlineSync.question("Enter new name: ")
			const matches = _BASENAME.exec(newDBName)
			if (matches && matches[0] === newDBName && newDBName.length !== 0) {
				if (fs.existsSync(__dirname + '/../databases/' + args[1] + ".json"))
					fs.renameSync(__dirname + `/../databases/${args[1]}.json`, __dirname + `/../databases/${newDBName}.json`)
				config.databases[config.databases.indexOf(args[1])] = newDBName
				console.log(chalk.green.bold(`Renamed ${args[1]} to ${newDBName}.`))
				fs.writeFileSync(__dirname + '/../config.json', JSON.stringify(config))
			} else {
				console.log(chalk.red.bold("Illeagal database name. Database names can contain characters A-Z, a-z, 0-9, -, _, ., and ,. They also should have a length between 1-100."))
			}
		} else {
			console.log(chalk.red.bold("Database not found."))
		}
	} else if (args[0] === 'current') {
		console.log(chalk.blue.bold(getDatabases().selected))
	} else if (args[0] === 'version') {
		const data = fs.readFileSync(__dirname + '/../package.json')
		try {
			console.log('v' + (JSON.parse(data).version || '0.0.0'))
		} catch (err) {
			console.log(chalk.red.bold("[FATAL] The package.json has been corrupted. Invalid JSON. ") + err)
		}
	} else {
		console.log(chalk.red.bold("Invalid argument."))
	}
}