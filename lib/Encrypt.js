const crypto = require('crypto')
const iv = crypto.randomBytes(16)

module.exports = {

  AES_encrypt: (text, pkey) => {
    const key = crypto
      .createHash('sha256')
      .update(pkey)
      .digest()
    const cipher = crypto.createCipheriv('aes-256-cbc', Buffer.from(key), iv)
    let encrypted = cipher.update(text)
    encrypted = Buffer.concat([encrypted, cipher.final()])
    return { iv: iv.toString('hex'), encryptedData: encrypted.toString('hex') }
  },

  AES_decrypt: (text, pkey) => {
    const key = crypto
      .createHash('sha256')
      .update(pkey)
      .digest()
    const iv = Buffer.from(text.iv, 'hex')
    const encryptedText = Buffer.from(text.encryptedData, 'hex')
    const decipher = crypto.createDecipheriv('aes-256-cbc', Buffer.from(key), iv)
    let decrypted = decipher.update(encryptedText)
    decrypted = Buffer.concat([decrypted, decipher.final()])
    return decrypted.toString()
  },

  SHA_hash: (string) => {
    return crypto
      .createHash('sha256')
      .update(string)
      .digest('hex')
  }
}
