import { request } from '@/utils/request'

type PublicKeyResp = {
  alg: string
  spkiBase64: string
}

function base64ToArrayBuffer(b64: string): ArrayBuffer {
  const binary = atob(b64)
  const bytes = new Uint8Array(binary.length)
  for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i)
  return bytes.buffer
}

function arrayBufferToBase64(buf: ArrayBuffer): string {
  const bytes = new Uint8Array(buf)
  let binary = ''
  for (let i = 0; i < bytes.length; i++) binary += String.fromCharCode(bytes[i]!)
  return btoa(binary)
}

let cachedKeyPromise: Promise<CryptoKey> | null = null

async function getCryptoKey(): Promise<CryptoKey> {
  if (cachedKeyPromise) return cachedKeyPromise
  cachedKeyPromise = (async () => {
    if (!globalThis.crypto?.subtle) {
      throw new Error('WebCrypto not available')
    }

    const pk = await request<PublicKeyResp>({
      url: '/user/public-key',
      method: 'GET',
    })
    if (!pk?.spkiBase64) throw new Error('Failed to load public key')

    const spki = base64ToArrayBuffer(pk.spkiBase64)
    return crypto.subtle.importKey('spki', spki, { name: 'RSA-OAEP', hash: 'SHA-256' }, false, ['encrypt'])
  })()
  return cachedKeyPromise
}

export async function encryptPassword(password: string): Promise<string> {
  const key = await getCryptoKey()
  const encoded = new TextEncoder().encode(password)
  const cipher = await crypto.subtle.encrypt({ name: 'RSA-OAEP' }, key, encoded)
  return arrayBufferToBase64(cipher)
}

