const readline = require("readline")
const chalk = require("chalk")

function log (query) {
  process.stdout.write(query)
}

const read = {
  visual: undefined,

  prompt: function (query) {
    return new Promise((resolve, reject) => {
      const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout
      })
      const stdin = process.openStdin();
      process.stdin.on('data', char => {
        char = char + ''
        switch (char) {
          case '\n':
          case '\r':
          case '\u0004':
            log("\u001b[0J\u001b[u\n")
            stdin.pause();
            break;
          default:
            log(`\u001b[${query.length + 1}G\u001b[s\u001b[0J\u001b[u`)
            log(rl.line + "\u001b[s");
            if (this.visual !== undefined) {
              this.visual(rl.line)
              log("\u001b[u")
            }
            break;
        }
      })
      rl.question(query, value => {
        rl.history = rl.history.slice(1)
        this.visual = undefined
        resolve(value)
      })
    })
  },
  setVisual: function (func) {
    if (typeof func !== "function" && func !== undefined)
      throw new TypeError("Expected funciton, recieved " + typeof func + ".")
    this.visual = func
  },
  password: function (query) {
    return new Promise((resolve, reject) => {
      const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout
      })
      const stdin = process.openStdin();
      process.stdin.on('data', char => {
        char = char + ''
        switch (char) {
          case '\n':
          case '\r':
          case '\u0004':
            stdin.pause();
            break;
          default:
            process.stdout.clearLine()
            readline.cursorTo(process.stdout, 0)
            log(query + "*".repeat(rl.line.length))
            break;
        }
      })
      rl.question(query, value => {
        rl.history = rl.history.slice(1)
        resolve(value)
      })
    })
  }
}

module.exports = read
