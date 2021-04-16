/*
 * krypt
 * https://github.com/Mkorp-Official/Krypt
 *
 * Copyright (c) 2021 Rak Laptudirm
 * Licensed under the MIT license.
 */

const crypto = require("crypto")

module.exports = {
  AES_encrypt: (text, pkey) => {
    const iv = crypto.randomBytes(16)
    const key = crypto.createHash("sha256").update(pkey).digest()
    const cipher = crypto.createCipheriv("aes-256-cbc", Buffer.from(key), iv)
    let encrypted = cipher.update(text)
    encrypted = Buffer.concat([encrypted, cipher.final()])
    return { iv: iv.toString("hex"), encryptedData: encrypted.toString("hex") }
  },

  AES_decrypt: (text, pkey) => {
    const key = crypto.createHash("sha256").update(pkey).digest()
    const iv = Buffer.from(text.iv, "hex")
    const encryptedText = Buffer.from(text.encryptedData, "hex")
    const decipher = crypto.createDecipheriv(
      "aes-256-cbc",
      Buffer.from(key),
      iv
    )
    let decrypted = decipher.update(encryptedText)
    decrypted = Buffer.concat([decrypted, decipher.final()])
    return decrypted.toString()
  },

  Bin_AES_encrypt: (text, pkey) => {
    const iv = crypto.randomBytes(16)
    const key = crypto.createHash("sha256").update(pkey).digest()
    const cipher = crypto.createCipheriv("aes-256-cbc", Buffer.from(key), iv)
    let encrypted = cipher.update(text)
    encrypted = Buffer.concat([encrypted, cipher.final()])
    return { iv: iv.toString("hex"), encryptedData: encrypted }
  },

  Bin_AES_decrypt: (text, pkey) => {
    const key = crypto.createHash("sha256").update(pkey).digest()
    const iv = Buffer.from(text.iv, "hex")
    const encryptedText = Buffer.from(text.encryptedData, "hex")
    const decipher = crypto.createDecipheriv(
      "aes-256-cbc",
      Buffer.from(key),
      iv
    )
    let decrypted = decipher.update(encryptedText)
    decrypted = Buffer.concat([decrypted, decipher.final()])
    return decrypted
  },

  SHA_hash: string => {
    return crypto.createHash("sha256").update(string).digest("hex")
  },

  random: max => {
    return crypto.randomInt(0, max + 1)
  },

  PBKDF2_HASH: (string, salt) => {
    if (salt) {
      salt = Buffer.from(salt, "hex")
      return crypto
        .pbkdf2Sync(string, salt, 500000, 32, "sha256")
        .toString("hex")
    } else {
      salt = crypto.randomBytes(16)
      return {
        checksum: crypto
          .pbkdf2Sync(string, salt, 500000, 32, "sha256")
          .toString("hex"),
        salt: salt.toString("hex"),
      }
    }
  },
}
