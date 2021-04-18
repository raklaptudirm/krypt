#!/usr/bin/env node

/*
 * krypt
 * https://github.com/Mkorp-Official/Krypt
 *
 * Copyright (c) 2021 Rak Laptudirm
 * Licensed under the MIT license.
 */

"use strict"

const zxcvbn = require("zxcvbn")
const { pwnedPassword } = require("hibp")
const fs = require("fs")
const crypto = require("../lib/crypto.js")
const readlineSync = require("readline-sync")
const chalk = require("chalk")
const clipboardy = require("clipboardy")
const style = require("ansi-styles")

/*
 * Terminal text themes
 */

const WARN = chalk.red.bold,
  CODE = chalk.black.bgGreen,
  OK = chalk.green.bold

/*
 * Constants
 *
 * These variable are global constants
 * and used for content verification,
 * RegExp checking, and help list
 */

const _DATA_TEMPLATE = {
    checksum: { checksum: "", salt: "" },
    salt: { key: "", TwoFA: "" },
    settings: {
      TwoFA: {
        on: false,
        question: "",
        answer: { checksum: "", salt: "" },
      },
      hint: { on: false, hint: "" },
      alias: {},
      passwordWordy: false,
    },
    data: {
      passwords: { iv: "", encryptedData: "" },
      notes: { iv: "", encryptedData: "" },
    },
  },
  // todo: update help
  _COMMS = [
    "new",
    "gent",
    "gpass",
    "secure",
    "change",
    "strength",
    "make",
    "edit",
    "delete",
    "set",
    "help",
    "exit",
    "list",
    "search",
    "copy",
    "archive",
    "notes",
  ],
  _BASENAME = /[A-Za-z0-9-_.,]{1,100}/,
  _HELP = {
    krypt: {
      use: "Global Krypt command.",
      new: {
        format: "krypt new",
        use: "Create a new Krypt database",
      },
      switch: {
        format: "krypt switch <db_name>",
        use: "Switch to the given database.",
      },
      list: {
        format: "krypt list",
        use: "List all the present Krypt databases.",
      },
      delete: {
        format: "krypt delete <db_name>",
        use: "Delete a Krypt database.",
      },
      rename: {
        format: "krypt rename <db_name>",
        use: "Rename a Krypt database.",
      },
      license: {
        format: "krypt license",
        use: "Prints out the Krypt License.",
      },
      version: {
        format: "krypt version",
        use: "Prints out the Krypt version in use.",
      },
      make: {
        format: "krypt make",
        use: "Generate a strong password based on arguments.",
        flags: {
          wordy: {
            use: "Generate a wordy password.",
            alias: "-w",
          },
        },
      },
      strength: {
        format: "krypt strength <password>",
        use: "Gives the strength of the given password.",
      },
      current: {
        format: "krypt current",
        use: "Prints the active database.",
      },
    },
    gpass: {
      format: "gpass <pass_id>",
      use: "Get a password in clear-text.",
    },
    new: {
      format: "new",
      use: "Create a new password.",
    },
    delete: {
      format: "delete <pass_id>",
      use: "Delete an existing password.",
    },
    change: {
      format: "change",
      use: "Change the master password.",
    },
    exit: {
      format: "exit",
      use: "Exit the Krypt session.",
      flags: {
        no_clear: {
          use: "Do not clear the console while exiting.",
          alias: "-ncl",
        },
      },
    },
    gent: {
      format: "gent <entry_id>",
      use: "Get a password entry.",
    },
    secure: {
      format: "secure",
      use: "Get security advise related to your passwords.",
      weak: {
        format: "secure weak",
        use: "Get a list of saved passwords that are weak.",
      },
      leaked: {
        format: "secure leaked",
        use: "Get a list of saved passwords that have been leaked.",
      },
      dups: {
        format: "secure dups",
        use: "Get a list of saved passwords which are duplicates.",
      },
    },
    make: {
      format: "make",
      use: "Generate a strong password based on settings.",
      flags: {
        wordy: {
          use: "Generate a wordy password.",
          alias: "-w",
        },
        no_wordy: {
          use: "Generate a non-wordy password.",
          alias: "-nw",
        },
      },
    },
    strength: {
      format: "strength <password>",
      use: "Gives the strength of the given password.",
    },
    help: {
      format: "help <command>",
      use: "Gives the uses and use example of a command.",
    },
    edit: {
      format: "edit <pass_id>",
      use: "Edit a stored password.",
    },
    list: {
      format: "list",
      use: "Lists all the stored passwords.",
    },
    search: {
      format: "search <keyword>",
      use: "Search a password with a particular name.",
    },
    copy: {
      format: "copy <pass_id>",
      use: "Copy a password to the clipboard.",
    },
    archive: {
      use: "Archive command package to archive files and directories.",
      new: {
        use: "Command package to create a new archive.",
        file: {
          format: "archive new file",
          use: "Archive a single file.",
        },
        dir: {
          format: "archive new dir",
          use: "Archive an entire directory.",
        },
      },
      unarc: {
        format: "archive unarc <archive_name>",
        use: "Un-archive an archive with the given name.",
      },
    },
    set: {
      use: "Command package for changing Krypt settings.",
      tfa: {
        format: "set tfa",
        use: "Enables 2-Factor Authentication, or edits it if already enabled.",
        dis: {
          format: "set tfa dis",
          use: "Disable 2-Factor Authentication.",
        },
      },
      hint: {
        format: "set hint",
        use: "Enables password hint, or edits it if already enabled.",
        dis: {
          format: "set tfa dis",
          use: "Disable password hint.",
        },
      },
      alias: {
        use: "Command package for setting command aliases.",
        new: {
          format: "set alias new",
          use: "Creates a new alias",
        },
        rename: {
          format: "set alias rename",
          use: "Rename a saved alias.",
        },
        delete: {
          format: "set alias delete <alias_name>",
          use: "Deletes the given alias.",
        },
        list: {
          format: "set alias list",
          use: "Lists the stored aliases.",
        },
      },
      password: {
        format: "set password",
        use: "Toggles Wordy-Password",
      },
    },
    notes: {
      use: "Command package for creating, seeing and deleting notes.",
      new: {
        format: "notes new",
        use: "Create a new note.",
      },
      get: {
        use: "Get an existing note from the database.",
        format: "notes get <note_no>",
      },
      delete: {
        use: "Delete an existing note.",
        format: "notes delete <note_no>",
      },
    },
  }

/*
 * Global Variables
 *
 * These variables store globally used
 * data like the database, passwords,
 * master key, 2nd factor key and word
 * list (Webster Dictionary).
 */

let _DATABASE,
  _PASSWORDS = {},
  _KEY,
  _2F,
  _NAME,
  _WORDS,
  _TREE,
  _MAST,
  _NOTES

/*
 * Main function
 *
 * This function definition defines the main process
 * of Krypt. It is invoked if Krypt receives no args.
 */

async function main() {
  if (fs.existsSync(__dirname + "/../databases/" + _NAME + ".json")) {
    if (!loadDatabase()) return
    _MAST = readlineSync.question("PASSWORD: ", { hideEchoBack: true })
    _KEY = crypto.PBKDF2_HASH(_MAST, _DATABASE.salt.key)

    if (_DATABASE.settings.TwoFA.on)
      _2F = crypto.PBKDF2_HASH(
        readlineSync.question(_DATABASE.settings.TwoFA.question + "? ", {
          hideEchoBack: true,
        }),
        _DATABASE.salt.TwoFA
      )

    if (
      _DATABASE.checksum.checksum ===
        crypto.PBKDF2_HASH(_KEY, _DATABASE.checksum.salt) &&
      (!_DATABASE.settings.TwoFA.on ||
        _DATABASE.settings.TwoFA.answer.checksum ===
          crypto.PBKDF2_HASH(_2F, _DATABASE.settings.TwoFA.answer.salt))
    ) {
      console.log("\u001b[2J")
      console.log("")
      LOGO()
      console.log("")
      console.log(`\n${OK(`Database: [ ${_NAME} ]`)}\n`)
      console.log("\n" + OK("Logged in."))
      loadData()
      while (true) {
        console.log("")
        let input = readlineSync.prompt() //.toLowerCase()
        if (_DATABASE.settings.alias[input] !== undefined)
          input = _DATABASE.settings.alias[input]
        input = input.split(" ").filter(item => item !== "")

        if (input[0] === "exit") {
          if (input.length > 2) {
            console.log(
              WARN(`Expected 0-1 arg(s), received ${input.length - 1}`)
            )
            continue
          }
          if (input[1] === "--no-clear" || input[1] === "-ncl") break
          else if (input.length > 1) {
            console.log(WARN(`Invalid argument.`))
            continue
          }
          console.log("\u001b[2J")
          break
        } else if (input[0] === "change") {
          if (input.length > 1) {
            console.log(WARN(`Expected 0 arg(s), received ${input.length - 1}`))
            continue
          }
          _KEY = crypto.PBKDF2_HASH(
            readlineSync.questionNewPassword("Enter new Password: ", { min: 8 })
          )
          _DATABASE.salt.key = _KEY.salt
          _KEY = _KEY.checksum
          _DATABASE.checksum = crypto.PBKDF2_HASH(_KEY)
          reEncryptData()
        } else if (input[0] === "new") {
          if (input.length > 1) {
            console.log(WARN(`Expected 0 arg(s), received ${input.length - 1}`))
            continue
          }
          const name_ = readlineSync.question("Password Name: ")
          const username_ = readlineSync.question("Username: ")
          const password_ =
            readlineSync.question("Password (leave empty to generate): ", {
              hideEchoBack: true,
            }) || generatePassword()
          _PASSWORDS.push(createPass(name_, username_, password_))
          console.log(
            OK(`Sucessfully added password at ID:${_PASSWORDS.length}.`)
          )
          reEncryptData()
        } else if (input[0] === "gent") {
          if (input.length !== 2) {
            console.log(WARN(`Expected 1 arg(s), received ${input.length - 1}`))
            continue
          }
          input = parseInt(input[1]) - 1
          if (
            input === undefined ||
            Number.isNaN(input) ||
            input < 0 ||
            input >= _PASSWORDS.length
          ) {
            console.log(WARN("ID out of bounds."))
            continue
          }
          printPass(_PASSWORDS[input], input + 1)
        } else if (input[0] === "gpass") {
          if (input.length !== 2) {
            console.log(WARN(`Expected 1 arg(s), received ${input.length - 1}`))
            continue
          }
          input = parseInt(input[1]) - 1
          if (
            input === undefined ||
            Number.isNaN(input) ||
            input < 0 ||
            input >= _PASSWORDS.length
          ) {
            console.log(WARN("ID out of bounds."))
          }
          const sel = readlineSync.question(
            WARN(
              "This command will show your password in clear-text. Proceed? (yes): "
            )
          )
          if (sel !== "yes") {
            console.log(OK("Command aborted."))
            continue
          }
          console.log(chalk.cyan.bold(_PASSWORDS[input].password))
        } else if (input[0] === "delete") {
          if (input.length !== 2) {
            console.log(WARN(`Expected 1 arg(s), received ${input.length - 1}`))
            continue
          }
          input = parseInt(input[1]) - 1
          if (
            input === undefined ||
            Number.isNaN(input) ||
            input < 0 ||
            input >= _PASSWORDS.length
          ) {
            console.log(WARN("ID out of bounds."))
            continue
          }
          printPass(_PASSWORDS[input], input)
          const sel = readlineSync.question(WARN("Delete this entry? (yes): "))
          if (sel !== "yes") {
            console.log(OK("Delete aborted."))
            continue
          }
          _PASSWORDS.splice(input, 1)
          console.log(OK("Password deleted Successfully."))
          reEncryptData()
        } else if (input[0] === "secure") {
          if (input.length > 2) {
            console.log(
              WARN(`Expected 0-1 arg(s), received ${input.length - 1}`)
            )
            continue
          }
          const [Sweaks, Pweaks, Duplicates] = await getWeaks()
          if (input[1] === undefined) {
            if (Sweaks.length > 0) {
              console.log(
                WARN(`✗ Some of your passwords are weak. `) +
                  CODE("secure weak") +
                  WARN(` for more information.`)
              )
            } else {
              console.log(OK("✓ All of the passwords are strong."))
            }
            if (Array.isArray(Pweaks)) {
              if (Pweaks.length > 0) {
                console.log(
                  WARN(`✗ Some of your passwords are leaked. `) +
                    CODE("secure leaked") +
                    WARN(` for more information.`)
                )
              } else {
                console.log(OK("✓ None of the passwords are leaked."))
              }
            } else {
              console.log(Pweaks)
            }
            if (Duplicates.length > 0) {
              console.log(
                WARN(`✗ Some of your passwords are duplicates. `) +
                  CODE("secure dups") +
                  WARN(` for more information.`)
              )
            } else {
              console.log(OK("✓ All of the passwords are unique."))
            }
            if (passStrength(_MAST).score !== OK("[VERY STRONG]")) {
              console.log(WARN("✗ Your master password is weak."))
            } else if (await pwnedPassword(_MAST)) {
              console.log(WARN("✗ Your master password has been leaked."))
            } else {
              console.log(OK("✓ Your master password is strong."))
            }
            if (_DATABASE.settings.TwoFA.on) {
              console.log(OK("✓ 2-Factor Auth is enabled."))
            } else {
              console.log(WARN("✗ 2-Factor Auth is disabled."))
            }
          } else if (input[1] === "weak") {
            if (Sweaks.length > 0) {
              for (const i in Sweaks) {
                console.log("")
                printPass(_PASSWORDS[i], parseInt(i))
              }
            } else {
              console.log(OK("All of the passwords are strong."))
            }
          } else if (input[1] === "leaked") {
            if (Pweaks.length > 0) {
              for (const i in Pweaks) {
                console.log("")
                printPass(_PASSWORDS[i], parseInt(i))
                console.log(await timesPwned(_PASSWORDS[i].password))
              }
            } else {
              console.log(OK("None of the passwords are leaked."))
            }
          } else if (input[1] === "dups") {
            if (Duplicates.length > 0) {
              for (const i of Duplicates) {
                console.log(
                  WARN(
                    `Passwords with ids ${i.join(", ")} have the same password.`
                  )
                )
              }
            } else {
              console.log(OK("All of the passwords are unique."))
            }
          } else {
            console.log(WARN("Invalid argument."))
          }
        } else if (input[0] === "make") {
          if (input.length > 2) {
            console.log(
              WARN(`Expected 0-1 arg(s), received ${input.length - 1}`)
            )
            continue
          }
          let type
          if (input[1] === undefined) type = _DATABASE.settings.passwordWordy
          else if (input[1] === "--wordy" || input[1] === "-w") type = true
          else if (input[1] === "--no-wordy" || input[1] === "-nw") type = false
          else {
            console.log(WARN("Invalid argument."))
            continue
          }
          const newPass = generatePassword(type)
          console.log(chalk.cyan.bold(newPass))
          console.log(passStrength(newPass).score + (await timesPwned(newPass)))
        } else if (input[0] === "help") {
          if (input.length < 2) {
            console.log(
              WARN(`Expected multiple arg(s), received ${input.length - 1}`)
            )
            continue
          }
          input.splice(0, 1)
          if (
            input.includes("use") ||
            input.includes("format") ||
            input.includes("flags")
          ) {
            console.log(WARN("Command not found."))
            continue
          }
          let manual = getItem(_HELP, input)
          if (manual === undefined) {
            console.log(WARN("Command not found."))
            continue
          }
          if (manual.format === undefined) {
            console.log(OK(manual.use) + "\n")
            console.log(chalk.bold(`Child commands:`))
            Object.keys(manual)
              .filter(item => item !== "use")
              .forEach(item => {
                console.log(`  ${chalk.bold(item)}: ${manual[item].use}`)
              })
          } else {
            console.log(`${CODE(manual.format)}\n${OK(manual.use)}\n`)
            Object.keys(manual)
              .filter(
                item => item !== "use" && item !== "format" && item !== "flags"
              )
              .forEach(item => {
                console.log(`  ${chalk.bold(item)}: ${manual[item].use}`)
              })
            if (manual.flags !== undefined) {
              console.log(chalk.bold(`Flags:`))
              Object.keys(manual.flags).forEach(item => {
                let id = item
                while (item.includes("_")) item = item.replace("_", "-")
                console.log(
                  `  --${chalk.bold(item)} (${manual.flags[id].alias}): ${
                    manual.flags[id].use
                  }`
                )
              })
            }
          }
        } else if (input[0] === "edit") {
          if (input.length !== 2) {
            console.log(WARN(`Expected 1 arg(s), received ${input.length - 1}`))
            continue
          }
          input = parseInt(input[1]) - 1
          if (
            input === undefined ||
            Number.isNaN(input) ||
            input < 0 ||
            input >= _PASSWORDS.length
          ) {
            console.log(WARN("ID out of bounds."))
            continue
          }
          const name_ = readlineSync.question(
            "Password Name (leave empty to keep same): "
          )
          const username_ = readlineSync.question(
            "Username (leave empty to keep same): "
          )
          const password_ =
            readlineSync.question("Password (leave empty to generate): ", {
              hideEchoBack: true,
            }) || generatePassword()
          _PASSWORDS[input] = createPass(
            name_ || _PASSWORDS[input].name,
            username_ || _PASSWORDS[input].username,
            password_
          )
          console.log(OK("Successfully edited password."))
          reEncryptData()
        } else if (input[0] === "list") {
          if (input.length > 1) {
            console.log(WARN(`Expected 0 arg(s), received ${input.length - 1}`))
            continue
          }
          for (const i in _PASSWORDS) {
            console.log("")
            printPass(_PASSWORDS[i], parseInt(i) + 1)
          }
          if (_PASSWORDS.length === 0)
            console.log(WARN("You do not have any stored passwords."))
        } else if (input[0] === "search") {
          if (input.length !== 2) {
            console.log(WARN(`Expected 1 arg(s), received ${input.length - 1}`))
            continue
          }
          let notFound = true
          for (const i in _PASSWORDS) {
            if (_PASSWORDS[i].name.toLowerCase().includes(input[1])) {
              printPass(_PASSWORDS[i], parseInt(i) + 1)
              notFound = false
              console.log("")
            }
          }
          if (notFound) console.log(WARN("No matches found."))
        } else if (input[0] === "set") {
          if (input.length < 2) {
            console.log(
              WARN(`Expected multiple arg(s), received ${input.length - 1}`)
            )
            continue
          }
          if (input[1] === "tfa") {
            if (!_DATABASE.settings.TwoFA.on) {
              if (input.length > 1) {
                console.log(
                  WARN(`Expected 0 arg(s), received ${input.length - 1}`)
                )
                continue
              }
              const sel = readlineSync.question(
                OK("Enable 2-Factor Authentication? (yes): ")
              )
              if (sel === "yes") {
                _DATABASE.settings.TwoFA.on = true
                _DATABASE.settings.TwoFA.question = readlineSync.question(
                  "Enter a question: "
                )
                _2F = crypto.PBKDF2_HASH(
                  readlineSync.question("Enter the answer: ")
                )
                _DATABASE.salt.TwoFA = _2F.salt
                _2F = _2F.checksum
                _DATABASE.settings.TwoFA.answer = crypto.PBKDF2_HASH(_2F)
                console.log(OK("Enabled 2 factor Auth."))
                reEncryptData()
              } else {
                console.log(WARN("Command aborted."))
              }
            } else {
              if (input[2] === "dis") {
                if (input.length > 3) {
                  console.log(
                    WARN(`Expected 0 arg(s), received ${input.length - 3}`)
                  )
                  continue
                }
                const sel = readlineSync.question(
                  WARN("Disable 2-Factor Authentication? (yes): ")
                )
                if (sel === "yes") {
                  _DATABASE.settings.TwoFA.on = false
                  console.log(OK("Disabled 2 factor Auth."))
                  reEncryptData()
                } else {
                  console.log(OK("Command aborted."))
                }
              } else {
                if (input.length > 2) {
                  console.log(
                    WARN(`Expected 0 arg(s), received ${input.length - 2}`)
                  )
                  continue
                }
                _DATABASE.settings.TwoFA.question =
                  readlineSync.question(
                    "Enter new question (Keep empty to keep the same): "
                  ) || _DATABASE.settings.TwoFA.question
                _2F = crypto.PBKDF2_HASH(
                  readlineSync.question("Enter new answer: ")
                )
                _DATABASE.salt.TwoFA = _2F.salt
                _2F = _2F.checksum
                _DATABASE.settings.TwoFA.answer = crypto.PBKDF2_HASH(_2F)
                console.log(OK("Changed Auth factors."))
                reEncryptData()
              }
            }
          } else if (input[1] === "hint") {
            if (!_DATABASE.settings.hint.on) {
              if (input.length > 2) {
                console.log(
                  WARN(`Expected 0 arg(s), received ${input.length - 2}`)
                )
                continue
              }
              const sel = readlineSync.question(OK("Enable Hint? (yes): "))
              if (sel === "yes") {
                _DATABASE.settings.hint.on = true
                _DATABASE.settings.hint.hint = readlineSync.question(
                  "Enter a hint: "
                )
                console.log(OK("Enabled Hint."))
                reEncryptData()
              } else {
                console.log(WARN("Command aborted."))
              }
            } else {
              if (input[2] === "dis") {
                if (input.length > 3) {
                  console.log(
                    WARN(`Expected 0 arg(s), received ${input.length - 3}`)
                  )
                  continue
                }
                const sel = readlineSync.question(WARN("Disable Hint? (yes): "))
                if (sel === "yes") {
                  _DATABASE.settings.hint.on = false
                  console.log(OK("Disabled Hint."))
                  reEncryptData()
                } else {
                  console.log(OK("Command aborted."))
                }
              } else {
                if (input.length > 2) {
                  console.log(
                    WARN(`Expected 0 arg(s), received ${input.length - 2}`)
                  )
                  continue
                }
                _DATABASE.settings.hint.hint =
                  readlineSync.question(
                    "Enter new hint (Keep empty to keep the same):"
                  ) || _DATABASE.settings.hint.hint
                console.log(OK("Changed Hint."))
                reEncryptData()
              }
            }
          } else if (input[1] === "alias") {
            if (input.length < 3) {
              console.log(
                WARN(`Expected multiple arg(s), received ${input.length - 2}`)
              )
              continue
            }
            if (input[2] === "new") {
              if (input.length > 3) {
                console.log(
                  WARN(`Expected 0 arg(s), received ${input.length - 3}`)
                )
                continue
              }
              const alias = readlineSync.question("Enter alias name: ")
              if (
                _DATABASE.settings.alias[alias] === undefined &&
                isNotCommand(alias)
              ) {
                if (is(alias, _BASENAME)) {
                  _DATABASE.settings.alias[alias] = readlineSync.question(
                    "Enter alias command: "
                  )
                  console.log(OK("Alias set."))
                  reEncryptData()
                } else {
                  console.log(WARN("Illegal alias name."))
                }
              } else {
                console.log(WARN("Command already exists."))
              }
            } else if (input[2] === "rename") {
              if (input.length !== 4) {
                console.log(
                  WARN(`Expected 1 arg(s), received ${input.length - 3}`)
                )
                continue
              }
              if (_DATABASE.settings.alias[input[3]] === undefined) {
                console.log(WARN("Could not find alias."))
              } else {
                const alias = readlineSync.question("Enter new name: ")
                if (is(alias, _BASENAME)) {
                  _DATABASE.settings.alias[alias] =
                    _DATABASE.settings.alias[input[3]]
                  delete _DATABASE.settings.alias[input[3]]
                  console.log(OK("Alias renamed successfully."))
                  reEncryptData()
                } else {
                  console.log(WARN("Illegal alias name."))
                }
              }
            } else if (input[2] === "delete") {
              if (input.length !== 4) {
                console.log(
                  WARN(`Expected 1 arg(s), received ${input.length - 3}`)
                )
                continue
              }
              if (_DATABASE.settings.alias[input[3]] === undefined) {
                console.log(WARN("Could not find alias."))
              } else {
                delete _DATABASE.settings.alias[input[3]]
                console.log(OK("Alias deleted successfully."))
                reEncryptData()
              }
            } else if (input[2] === "list") {
              if (input.length > 3) {
                console.log(
                  WARN(`Expected 0 arg(s), received ${input.length - 3}`)
                )
                continue
              }
              console.log("")
              for (const i in _DATABASE.settings.alias) {
                console.log(
                  `${chalk.blue.bold(i)}: ${chalk.yellow.bold(
                    _DATABASE.settings.alias[i]
                  )}`
                )
              }
              if (Object.keys(_DATABASE.settings.alias).length === 0)
                console.log(WARN("You do not have any stored aliases."))
            } else {
              console.log(WARN("Invalid argument."))
            }
          } else if (input[1] === "password") {
            if (input.length > 2) {
              console.log(
                WARN(`Expected 0 arg(s), received ${input.length - 2}`)
              )
              continue
            }
            _DATABASE.settings.passwordWordy = !_DATABASE.settings.passwordWordy
            if (_DATABASE.settings.passwordWordy)
              console.log(OK("Enabled wordy password."))
            else console.log(OK("Disabled wordy password."))
          } else {
            console.log(WARN("Setting not found."))
          }
        } else if (input[0] === "copy") {
          if (input.length !== 2) {
            console.log(WARN(`Expected 1 arg(s), received ${input.length - 1}`))
            continue
          }
          input = parseInt(input[1]) - 1
          if (
            input === undefined ||
            Number.isNaN(input) ||
            input < 0 ||
            input >= _PASSWORDS.length
          ) {
            console.log(WARN("ID out of bounds."))
          } else {
            clipboardy.writeSync(_PASSWORDS[input].password)
            console.log(OK("Password copied to clipboard."))
          }
        } else if (input[0] === "strength") {
          if (input.length !== 2) {
            console.log(WARN(`Expected 1 arg(s), received ${input.length - 1}`))
            continue
          }
          if (input[1]) {
            const pStrength = passStrength(input[1])
            console.log(
              `${pStrength.score}${await timesPwned(
                input[1]
              )} \nTime required to break password: ${pStrength.time}`
            )
            if (pStrength.feedback.warning)
              console.log(WARN(`Warning: ${pStrength.feedback.warning}`))
            if (pStrength.feedback.suggestions.length !== 0)
              console.log(
                OK(`Suggestions: ${pStrength.feedback.suggestions.join(", ")}`)
              )
          } else {
            console.log(WARN("Please enter a password."))
          }
        } else if (input[0] === "archive") {
          if (input.length !== 3) {
            console.log(WARN(`Expected 2 arg(s), received ${input.length - 1}`))
            continue
          }
          if (!fs.existsSync(__dirname + "/../databases/" + _NAME))
            fs.mkdirSync(__dirname + "/../databases/" + _NAME)
          if (!fs.existsSync(__dirname + "/../databases/" + _NAME + "/.tree"))
            fs.writeFileSync(
              __dirname + "/../databases/" + _NAME + "/.tree",
              "{}"
            )
          _TREE = JSON.parse(
            fs.readFileSync(__dirname + "/../databases/" + _NAME + "/.tree")
          )

          if (input[1] === "new") {
            if (input.length < 2) {
              console.log(
                WARN(`Expected multiple arg(s), received ${input.length - 1}`)
              )
              continue
            }
            if (input[2] === "file") {
              const fPath = readlineSync.questionPath("Enter file path: ", {
                isFile: true,
                exists: true,
              })
              const fName = readlineSync.question("Enter file name: ")
              if (_TREE[fName] === undefined) {
                _TREE[fName] = fPath
                fs.writeFileSync(
                  __dirname + "/../databases/" + _NAME + "/" + fName + ".karc",
                  JSON.stringify(binEncryptFile(fs.readFileSync(fPath)))
                )
                fs.unlinkSync(fPath)
                console.log(OK("Archived file successfully."))
                updateTree()
              } else {
                console.log(WARN("Archive already exists."))
              }
            } else if (input[2] === "dir") {
              const fPath = readlineSync.questionPath(
                "Enter directory path: ",
                {
                  isDirectory: true,
                  exists: true,
                }
              )
              const fName = readlineSync.question("Enter directory name: ")
              if (_TREE[fName] === undefined) {
                fs.mkdirSync(__dirname + "/../databases/" + _NAME + "/" + fName)
                _TREE[fName] = "DIRECTORY"
                let dirTree = { path: fPath + "\\" },
                  num = 0
                let files = getAllFiles(fPath).forEach(item => {
                  dirTree[num.toString()] = item
                  fs.writeFileSync(
                    __dirname +
                      "/../databases/" +
                      _NAME +
                      "/" +
                      fName +
                      "/" +
                      num +
                      ".karc",
                    JSON.stringify(
                      binEncryptFile(fs.readFileSync(fPath + "/" + item))
                    )
                  )
                  fs.unlinkSync(fPath + "/" + item)
                  num++
                })
                fs.writeFileSync(
                  __dirname + "/../databases/" + _NAME + "/" + fName + "/.tree",
                  JSON.stringify(dirTree)
                )
                console.log(OK("Archived directory successfully."))
                updateTree()
              } else {
                console.log(WARN("Archive already exists"))
              }
            } else {
              console.log(WARN("Illegal command."))
            }
          } else if (input[1] === "unarc") {
            if (_TREE[input[2]] === undefined) {
              console.log(
                WARN(`Archived file with name ${input[2]} not found.`)
              )
            } else if (_TREE[input[2]] === "DIRECTORY") {
              let dirTree = JSON.parse(
                fs.readFileSync(
                  __dirname +
                    "/../databases/" +
                    _NAME +
                    "/" +
                    input[2] +
                    "/.tree"
                )
              )
              let treeKeys = Object.keys(dirTree)
              treeKeys.splice(treeKeys.indexOf("path"), 1)
              for (const i of treeKeys) {
                const file = binDecryptFile(
                  JSON.parse(
                    fs.readFileSync(
                      __dirname +
                        "/../databases/" +
                        _NAME +
                        "/" +
                        input[2] +
                        "/" +
                        i +
                        ".karc"
                    )
                  )
                )
                fs.unlinkSync(
                  __dirname +
                    "/../databases/" +
                    _NAME +
                    "/" +
                    input[2] +
                    "/" +
                    i +
                    ".karc"
                )
                fs.writeFileSync(dirTree.path + dirTree[i], file)
              }
              delete _TREE[input[2]]
              console.log(
                OK(
                  `Unarchived file sucessfully. It can be found at ${dirTree.path}`
                )
              )
              fs.unlinkSync(
                __dirname + "/../databases/" + _NAME + "/" + input[2] + "/.tree"
              )
              fs.rmdirSync(
                __dirname + "/../databases/" + _NAME + "/" + input[2]
              )
              updateTree()
            } else {
              const file = binDecryptFile(
                JSON.parse(
                  fs.readFileSync(
                    __dirname +
                      "/../databases/" +
                      _NAME +
                      "/" +
                      input[2] +
                      ".karc"
                  )
                )
              )
              fs.unlinkSync(
                __dirname + "/../databases/" + _NAME + "/" + input[2] + ".karc"
              )
              fs.writeFileSync(_TREE[input[2]], file)
              console.log(
                OK(
                  `Unarchived file sucessfully. It can be found at ${
                    _TREE[input[2]]
                  }`
                )
              )
              delete _TREE[input[2]]
              updateTree()
            }
          } else {
            console.log(WARN("Invalid argument."))
          }
        } else if (input[0] === "notes") {
          if (input.length < 2) {
            console.log(
              WARN(`Expected multiple arg(s), received ${input.length - 1}`)
            )
            continue
          }
          if (input[1] === "new") {
            if (input.length > 2) {
              console.log(
                WARN(`Expected 0 arg(s), received ${input.length - 2}`)
              )
              continue
            }
            let name = readlineSync.question("Enter note name: ")
            let lines = []
            console.log(
              `Enter your note. Enter to go to next line. Type END to end input: \n\n${"-".repeat(
                24
              )}\n`
            )
            readlineSync.promptLoop(
              line => {
                lines.push(line)
                return line === "END"
              },
              { prompt: "" }
            )
            lines.pop()
            lines = lines.join("\n")
            _NOTES.push({ name: name, info: lines, date: new Date() })
            console.log(OK("Added note."))
            reEncryptData()
          } else if (input[1] === "get") {
            if (input.length !== 3) {
              console.log(
                WARN(`Expected 1 arg(s), received ${input.length - 1}`)
              )
              continue
            }
            input = parseInt(input[2]) - 1
            if (
              input === undefined ||
              Number.isNaN(input) ||
              input < 0 ||
              input >= _NOTES.length
            ) {
              console.log(WARN("ID out of bounds."))
            } else {
              printNote(_NOTES[input], input + 1)
            }
          } else if (input[1] === "delete") {
            if (input.length !== 3) {
              console.log(
                WARN(`Expected 1 arg(s), received ${input.length - 1}`)
              )
              continue
            }
            input = parseInt(input[2]) - 1
            if (
              input === undefined ||
              Number.isNaN(input) ||
              input < 0 ||
              input >= _NOTES.length
            ) {
              console.log(WARN("ID out of bounds."))
            } else {
              printNote(_NOTES[input], input + 1)
              const _delete = readlineSync.question(
                WARN("Delete this note (yes)? ")
              )
              if (_delete === "yes") {
                _NOTES.splice(input, 1)
                console.log(OK("Sucessfully deleted note."))
                reEncryptData()
              } else {
                console.log(OK("Delete aborted."))
              }
            }
          } else {
            console.log(WARN("Invalid argument."))
          }
        } else {
          console.log(WARN("Invalid command."))
        }
      }
    } else {
      console.log(
        WARN(
          _DATABASE.settings.TwoFA.on
            ? "Wrong Password or 2nd factor."
            : "Wrong Password."
        )
      )
      if (_DATABASE.settings.hint.on)
        console.log(OK(`Hint: ${_DATABASE.settings.hint.hint}`))
    }
  } else {
    if (!fs.existsSync(__dirname + "/../databases"))
      fs.mkdirSync(__dirname + "/../databases")
    _DATABASE = _DATA_TEMPLATE
    _PASSWORDS = []
    _NOTES = []
    _KEY = crypto.PBKDF2_HASH(
      readlineSync.questionNewPassword("New Password: ", { min: 8 })
    )
    _DATABASE.salt.key = _KEY.salt
    _KEY = _KEY.checksum
    _DATABASE.checksum = crypto.PBKDF2_HASH(_KEY)
    console.log("\n" + OK("Database initialized."))
    reEncryptData()
  }
}

/*
 * Function definitions
 *
 * All the functions used by Krypt are defined here.
 *
 * List:
 * [01] isNotCommand
 *        name -> String
 *      checks if arg:name is a command
 *      returns -> bool
 * [02] gatherKeys
 *        obj -> Object
 *      Gathers the keys of arg:obj
 *      returns -> Set
 * [03] diff
 *        a -> Object
 *        b -> Object
 *      Gathers the key differences between arg:a and arg:b
 *      returns -> Array
 * [04] equalByKeys
 *        a -> Object
 *        b -> Object
 *      Checks if arg:a and arg:b have the same keys
 *      returns -> bool
 * [05] generatePassword
 *      Generate a random password
 *      returns -> String
 * [06] passStrength
 *        pass -> String
 *      Gives the strength of the arg:pass password
 *      returns -> Object > passwordScore
 * [07] printPass
 *        password -> String
 *        id -> Number
 *      Prints the arg:password out with styling
 *      returns -> void
 * [08] createPass
 *        name -> String
 *        username -> String
 *        password -> String
 *      Creates a password Object
 *      returns -> Object > password
 * [09] loadDatabase
 *      Loads the selected database into global:_DATABASE
 *      returns -> void
 * [10] loadData
 *      Loads the passwords from global:_DATABASE to global:_PASSWORDS
 *      returns -> void
 * [11] reEncryptData
 *      Re-encrypts global:_DATABASE into database file
 *      returns -> void
 * [12] decryptData
 *      decrypts the passwords from global:_DATABASE
 *      returns -> Object > [ password, .. ]
 * [13] getDatabases
 *      Loads config.json file
 *      returns -> Object > config
 * [14] is
 *        string -> String
 *        regexp -> RegExp
 *      Checks if arg:string matches arg:regexp
 *      returns -> bool
 * [15] binEncryptFile
 *        file -> Buffer
 *      Encrypt a file
 *      returns -> Object > ciphertext
 * [16] binDecryptFile
 *        ciphertext -> Object > ciphertext
 *      Decrypt a file
 *      returns -> Buffer
 * [17] updateTree
 *      Updates the .tree file of the database archive
 *      returns -> void
 * [18] getAllFiles
 *        dir -> String > path
 *      Get a list of all files in a directory
 *      returns -> Array[path...]
 * [19] LOGO
 *      Print the Krypt logo into the terminal
 *      returns -> void
 * [20] timesPwned
 *        pass -> String
 *      Checks how many times arg:pass has been leaked.
 *      returns -> String
 * [21] getWeaks
 *      Get the weak, leaked and duplicate passwords.
 *      returns -> Array[Array, Array, Array]
 * [22] getItems
 *        ob -> Object
 *        path -> string
 *      Gives the information at arg:path in arg:ob
 *      returns -> Any
 * [23] printNote
 *        note -> Object > note
 *        index -> Number
 *      Prints out the arg:note with formatting.
 */

function isNotCommand(name) {
  return !_COMMS.includes(name)
}

const gatherKeys = obj => {
  const isObject = val => typeof val === "object" && !Array.isArray(val),
    addDelimiter = (a, b) => (a ? `${a}.${b}` : b),
    paths = (obj = {}, head = "") =>
      Object.entries(obj).reduce(
        (product, [key, value]) =>
          (fullPath =>
            isObject(value)
              ? product.concat(paths(value, fullPath))
              : product.concat(fullPath))(addDelimiter(head, key)),
        []
      )
  return new Set(paths(obj))
}

const diff = (a, b) => new Set(Array.from(a).filter(item => !b.has(item))),
  equalByKeys = (a, b) => diff(gatherKeys(a), gatherKeys(b)).size === 0

function generatePassword(wordy) {
  const _lowerCase = "abcdefghijklmnopqrstuvwxyz"
  const _upperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
  const _numbers = "0123456789"
  const _specialChars = ",./;'[]\\=-`<>?\":|}{+_~!@#$%^&*()"
  let password

  if (wordy) {
    let len = _WORDS.length - 1
    password =
      _WORDS[crypto.random(len)] +
      _WORDS[crypto.random(len)] +
      _WORDS[crypto.random(len)] +
      _WORDS[crypto.random(len)]
  } else {
    do {
      password = ""
      const length = 12

      for (let i = 0; i < length; i++) {
        let type = crypto.random(3)
        switch (type) {
          case 0:
            password += _lowerCase[crypto.random(25)]
            break
          case 1:
            password += _upperCase[crypto.random(25)]
            break
          case 2:
            password += _numbers[crypto.random(9)]
            break
          case 3:
            password += _specialChars[crypto.random(_specialChars.length - 1)]
        }
      }
    } while (passStrength(password).score !== OK("[VERY STRONG]"))
  }
  return password
}

function passStrength(passwordS) {
  const measure = [
    WARN("[VERY WEAK]"),
    chalk.yellow.bold("[WEAK]"),
    "[MEDIUM]",
    chalk.blue.bold("[STRONG]"),
    OK("[VERY STRONG]"),
  ]

  const power = zxcvbn(passwordS)
  return {
    score: measure[power.score],
    time: power.crack_times_display.online_no_throttling_10_per_second,
    feedback: power.feedback,
  }
}

function printPass(password, id) {
  console.log(chalk.blue(`[ID:${id}]${passStrength(password.password).score}`))
  console.log(
    "Name: " +
      chalk.yellow.bold(password.name) +
      "\n" +
      "Username: " +
      chalk.yellow.bold(password.username) +
      "\n" +
      "Password: " +
      chalk.yellow.bold("*".repeat(password.password.length))
  )
}

function createPass(name, username, password) {
  return { name: name, username: username, password: password }
}

function loadDatabase() {
  const data = fs.readFileSync(__dirname + "/../databases/" + _NAME + ".json")
  try {
    _DATABASE = JSON.parse(data)
    if (equalByKeys(_DATA_TEMPLATE, _DATABASE)) return true
    console.log(
      WARN("[FATAL] The database has been corrupted. Invalid Property List")
    )
    return false
  } catch (err) {
    console.log(
      WARN("[FATAL] The database has been corrupted. Invalid JSON. ") + err
    )
    return false
  }
}

function loadData() {
  _PASSWORDS = JSON.parse(decryptData(_DATABASE.data.passwords))
  _NOTES = JSON.parse(decryptData(_DATABASE.data.notes))
}

function reEncryptData() {
  if (_DATABASE.settings.TwoFA.on) {
    _DATABASE.data.passwords = crypto.AES_encrypt(
      JSON.stringify(crypto.AES_encrypt(JSON.stringify(_PASSWORDS), _KEY)),
      _2F
    )
    _DATABASE.data.notes = crypto.AES_encrypt(
      JSON.stringify(crypto.AES_encrypt(JSON.stringify(_NOTES), _KEY)),
      _2F
    )
  } else {
    _DATABASE.data.passwords = crypto.AES_encrypt(
      JSON.stringify(_PASSWORDS),
      _KEY
    )
    _DATABASE.data.notes = crypto.AES_encrypt(JSON.stringify(_NOTES), _KEY)
  }
  fs.writeFileSync(
    __dirname + "/../databases/" + _NAME + ".json",
    JSON.stringify(_DATABASE)
  )
}

function decryptData(data) {
  if (_DATABASE.settings.TwoFA.on)
    return crypto.AES_decrypt(JSON.parse(crypto.AES_decrypt(data, _2F)), _KEY)
  return crypto.AES_decrypt(data, _KEY)
}

function getDatabases() {
  const data = fs.readFileSync(__dirname + "/../config.json")
  try {
    const config = JSON.parse(data)
    return config
  } catch (err) {
    console.log(
      WARN("[FATAL] The database list has been corrupted. Invalid JSON. ") + err
    )
    return false
  }
}

const is = (string, regexp) => {
  let match = regexp.exec(string)
  if (string === "") return false
  if (match && match[0] === string) return true
  return false
}

function binEncryptFile(file) {
  if (_DATABASE.settings.TwoFA.on)
    return crypto.Bin_AES_encrypt(crypto.Bin_AES_encrypt(file, _KEY), _2F)
  return crypto.Bin_AES_encrypt(file, _KEY)
}

function binDecryptFile(ciphertext) {
  if (_DATABASE.settings.TwoFA.on)
    return crypto.Bin_AES_decrypt(crypto.Bin_AES_decrypt(ciphertext, _2F), _KEY)
  return crypto.Bin_AES_decrypt(ciphertext, _KEY)
}

function updateTree() {
  fs.writeFileSync(
    __dirname + "/../databases/" + _NAME + "/.tree",
    JSON.stringify(_TREE)
  )
}

function getAllFiles(dir) {
  let files = []
  fs.readdirSync(dir).forEach(file => {
    if (fs.lstatSync(`${dir}\\${file}`).isFile()) files.push(file)
    else
      getAllFiles(`${dir}\\${file}`).forEach(item => {
        files.push(`${file}\\${item}`)
      })
  })
  return files
}

function LOGO() {
  const logo = fs.readFileSync(__dirname + "/../logo").toString()
  console.log(OK(logo))
}

async function timesPwned(pass) {
  try {
    const times = await pwnedPassword(pass)
    if (times === 0) return OK("[No Occurances]")
    return WARN(`[Occurances:${times}]`)
  } catch {
    return WARN("[No Internet]")
  }
}

async function getWeaks() {
  let weakS = [],
    pwned = [],
    list = {},
    duplicates = [],
    Internet = true
  for (const i in _PASSWORDS) {
    if (list[_PASSWORDS[i].password] === undefined)
      list[_PASSWORDS[i].password] = []
    list[_PASSWORDS[i].password].push(parseInt(i) + 1)

    if (passStrength(_PASSWORDS[i].password).score !== OK("[VERY STRONG]"))
      weakS.push(i)
    if (Internet) {
      try {
        if ((await pwnedPassword(_PASSWORDS[i].password)) > 0) pwned.push(i)
      } catch {
        pwned = WARN("[No Internet]")
      }
    }
  }

  for (const i in list) {
    if (list[i].length > 1) duplicates.push(list[i])
  }

  return [weakS, pwned, duplicates]
}

function getItem(ob, path) {
  for (const key of path) {
    ob = ob[key]
    if (ob === undefined) return undefined
  }
  return ob
}

function printNote(note, index) {
  let _FORMAT_OP = /<(reset|bold|dim|italic|underline|overline|inverse|hidden|strikethrough|black|red|green|yellow|blue|magenta|cyan|white|blackBright|redBright|greenBright|yellowBright|blueBright|magentaBright|cyanBright|whiteBright|bgBlack|bgRed|bgGreen|bgYellow|bgBlue|bgMagenta|bgCyan|bgWhite|bgBlackBright|bgRedBright|bgGreenBright|bgYellowBright|bgBlueBright|bgMagentaBright|bgCyanBright|bgWhiteBright)>/,
    _FORMAT_CL = /<\/(reset|bold|dim|italic|underline|overline|inverse|hidden|strikethrough|black|red|green|yellow|blue|magenta|cyan|white|blackBright|redBright|greenBright|yellowBright|blueBright|magentaBright|cyanBright|whiteBright|bgBlack|bgRed|bgGreen|bgYellow|bgBlue|bgMagenta|bgCyan|bgWhite|bgBlackBright|bgRedBright|bgGreenBright|bgYellowBright|bgBlueBright|bgMagentaBright|bgCyanBright|bgWhiteBright)>/,
    str = note.info

  while (_FORMAT_OP.exec(str)) {
    let token = _FORMAT_OP.exec(str)[0]
    token = token.substring(1, token.length - 1)
    str = str.replace(_FORMAT_OP, style[token].open)
  }

  while (_FORMAT_CL.exec(str)) {
    let token = _FORMAT_CL.exec(str)[0]
    token = token.substring(2, token.length - 1)
    str = str.replace(_FORMAT_CL, style[token].close)
  }

  console.log(
    `${chalk.bold(`[${index}] ${note.name}`)}\n${chalk.bold(
      note.date
    )}\n\n${"-".repeat(24)}\n\n${str}\n\n${"-".repeat(24)}`
  )
}

/*
* Filters:
* --name (-n): String
* --username (-u): String
* --leaked (-l): void
* --strength (-s): 0/very-weak | 1/weak | 2/medium | 3/strong | 4/very-strong
* --contains (-c): String -> [,]Array
*
* Logic:
* &
* |
*/

async function filterPass(filters) {
  let length = filters.length, filtered = []
  for (const i in _PASSWORDS) {
    for (const j = 0; j < length; j++) {
      switch (filter[j]) {
        case "--name":
        case "-n":
          if (_PASSWORDS[i].name.toLowerCase().includes(filter[j + 1].toLowerCase()))
            filtered.push(i)
          j++
          break
        case "--username":
        case "-u":
          if (_PASSWORDS[i].username.toLowerCase().includes(filter[j + 1].toLowerCase()))
            filtered.push(i)
          j++
          break
        case "--leaked":
        case "-l":
          if (!(await pwnedPassword(_PASSWORDS[i].password)))
            filtered.push(i)
          break
        case "--strength":
        case "-s":
          break
        default:
        console.log("err")
      }
    }
  }
}

/*
 * Main process
 *
 * The main process checks for command-line args,
 * and works accordingly, invoking main function
 * if no arguments are found.
 */

;(async function () {
  _WORDS = JSON.parse(fs.readFileSync(__dirname + "/../lib/dictionary.json"))
  if (process.argv.length === 2) {
    if (!fs.existsSync(__dirname + "/../config.json"))
      fs.writeFileSync(
        __dirname + "/../config.json",
        '{"selected": "default","databases": ["default"]}'
      )
    _NAME = getDatabases()
    _NAME = _NAME.selected
    console.log("")
    LOGO()
    console.log("")
    console.log(`\n${OK(`Database: [ ${_NAME} ]`)}\n`)
    main()
  } else {
    let args = process.argv.slice(2)
    if (args[0] === "new") {
      if (args.length > 1) {
        console.log(WARN(`Expected 0 arg(s), received ${args.length - 1}`))
        return
      }
      let config = getDatabases()
      const newName = readlineSync.question("Enter database name: ")
      if (is(newName, _BASENAME) && newName.length !== 0) {
        if (!config.databases.includes(newName)) {
          config.databases.push(newName)
          fs.writeFileSync(
            __dirname + "/../config.json",
            JSON.stringify(config)
          )
          console.log(OK("Added new database."))
        } else {
          console.log(WARN("Database already exists."))
        }
      } else {
        console.log(
          WARN(
            "Illegal database name. Database names can contain characters A-Z, a-z, 0-9, -, _, ., and ,. They also should have a length between 1-100."
          )
        )
      }
    } else if (args[0] === "list") {
      if (args.length > 1) {
        console.log(WARN(`Expected 0 arg(s), received ${args.length - 1}`))
        return
      }
      let config = getDatabases()
      for (const databaseName of config.databases) {
        console.log(chalk.blue.bold(databaseName))
      }
    } else if (args[0] === "switch") {
      if (args.length !== 2) {
        console.log(WARN(`Expected 1 arg(s), received ${args.length - 1}`))
        return
      }
      let config = getDatabases()
      if (config.databases.includes(args[1])) {
        config.selected = args[1]
        fs.writeFileSync(__dirname + "/../config.json", JSON.stringify(config))
        console.log(OK(`Switched to ${args[1]} database.`))
      } else {
        console.log(WARN("Database not found."))
      }
    } else if (args[0] === "delete") {
      if (args.length !== 2) {
        console.log(WARN(`Expected 1 arg(s), received ${args.length - 1}`))
        return
      }
      let config = getDatabases()
      if (config.databases.includes(args[1])) {
        if (config.databases.length === 1) {
          console.log(WARN("Can't delete last database."))
        } else {
          if (
            readlineSync.question(
              WARN(`Delete the ${args[1]} database? (yes): `)
            ) === "yes"
          ) {
            if (fs.existsSync(__dirname + "/../databases/" + args[1] + ".json"))
              fs.unlinkSync(__dirname + "/../databases/" + args[1] + ".json")
            config.databases.splice(config.databases.indexOf(args[1]), 1)
            if (config.selected === args[1]) {
              config.selected = config.databases[0]
            }
            console.log(OK(`Deleted ${args[1]} database.`))
            fs.writeFileSync(
              __dirname + "/../config.json",
              JSON.stringify(config)
            )
          } else {
            console.log(OK("Delete aborted."))
          }
        }
      } else {
        console.log(WARN("Database not found."))
      }
    } else if (args[0] === "rename") {
      if (args.length !== 2) {
        console.log(WARN(`Expected 1 arg(s), received ${args.length - 1}`))
        return
      }
      let config = getDatabases()
      if (config.databases.includes(args[1])) {
        const newDBName = readlineSync.question("Enter new name: ")
        if (is(newDBName, _BASENAME) && newDBName.length !== 0) {
          if (fs.existsSync(__dirname + "/../databases/" + args[1] + ".json"))
            fs.renameSync(
              __dirname + `/../databases/${args[1]}.json`,
              __dirname + `/../databases/${newDBName}.json`
            )
          config.databases[config.databases.indexOf(args[1])] = newDBName
          console.log(OK(`Renamed ${args[1]} to ${newDBName}.`))
          fs.writeFileSync(
            __dirname + "/../config.json",
            JSON.stringify(config)
          )
        } else {
          console.log(
            WARN(
              "Illegal database name. Database names can contain characters A-Z, a-z, 0-9, -, _, ., and ,. They also should have a length between 1-100."
            )
          )
        }
      } else {
        console.log(WARN("Database not found."))
      }
    } else if (args[0] === "current") {
      if (args.length > 1) {
        console.log(WARN(`Expected 0 arg(s), received ${args.length - 1}`))
        return
      }
      console.log(chalk.blue.bold(getDatabases().selected))
    } else if (args[0] === "version") {
      if (args.length > 1) {
        console.log(WARN(`Expected 0 arg(s), received ${args.length - 1}`))
        return
      }
      const data = fs.readFileSync(__dirname + "/../package.json")
      try {
        console.log("v" + (JSON.parse(data).version || "0.0.0"))
      } catch (err) {
        console.log(
          WARN("[FATAL] The package.json has been corrupted. Invalid JSON. ") +
            err
        )
      }
    } else if (args[0] === "license") {
      if (args.length > 1) {
        console.log(WARN(`Expected 0 arg(s), received ${args.length - 1}`))
        return
      }
      console.log(
        `\n${chalk.bold("Permissions:")}\n${OK(
          "* Commercial use\n* Distribution\n* Modification\n* Private use"
        )}\n\n${chalk.bold("Conditions:")}\n${chalk.cyan.bold(
          "* License and copyright notice"
        )}\n\n${chalk.bold("Limitations:")}\n${WARN(
          "* Liability\n* Warranty"
        )}\n`
      )
      console.log(
        chalk.whiteBright(fs.readFileSync(`${__dirname}/../LICENSE`).toString())
      )
    } else if (args[0] === "make") {
      if (args.length > 2) {
        console.log(WARN(`Expected 0-1 arg(s), received ${args.length - 1}`))
        return
      }
      let wordy
      if (args[1] === "--wordy" || args[1] === "-w") wordy = true
      else wordy = false
      const newPass = generatePassword(wordy)
      console.log(chalk.cyan.bold(newPass))
      console.log(passStrength(newPass).score + (await timesPwned(newPass)))
    } else if (args[0] === "strength") {
      if (args.length !== 2) {
        console.log(WARN(`Expected 1 arg(s), received ${args.length - 1}`))
        return
      }
      if (args[1]) {
        const pStrength = passStrength(args[1])
        console.log(
          `${pStrength.score}${await timesPwned(
            args[1]
          )} \nTime required to break password: ${pStrength.time}`
        )
        if (pStrength.feedback.warning)
          console.log(WARN(`Warning: ${pStrength.feedback.warning}`))
        if (pStrength.feedback.suggestions.length !== 0)
          console.log(
            OK(`Suggestions: ${pStrength.feedback.suggestions.join(", ")}`)
          )
      } else {
        console.log(WARN("Please enter a password."))
      }
    } else {
      console.log(WARN("Invalid argument."))
    }
  }
})()
