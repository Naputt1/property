import type { JobStatus } from '@/types/job'

export interface ISocketMessage {
  job_update: {
    id: string
    status: JobStatus
    progress: number
    total: number
  }
}

export type SocketMessage = InternalSocketMessage<ISocketMessage>

interface InternalMsg extends ISocketMessage {
  ping: undefined
  pong: undefined
  connect: undefined
}

type InternalSocketMessage<T> = {
  [K in keyof T]: { type: K; data: T[K] }
}[keyof T]

export type SocketCallback = (
  data: InternalSocketMessage<ISocketMessage>,
) => void

export class Ws {
  private ws: WebSocket | null = null
  private reconnectTimeout: any = null
  private heartbeatInterval: any = null
  private readonly HEARTBEAT_INTERVAL = 10000 // 10s
  private readonly PONG_TIMEOUT = 10000
  private readonly RECONNECT_DELAY = 1000
  private readonly MAX_RECONNECT_DELAY = 30000
  private reconnectAttempts = 0
  private pongTimeout: any = undefined

  onmessage: { [key: string]: SocketCallback } = {}
  isConnect = false
  url: string

  constructor(url: string) {
    this.url = url
    this.connect()
  }

  private connect() {
    if (this.ws) {
      this.ws.onclose = null
      this.ws.onmessage = null
      this.ws.close()
    }

    // Convert relative URL to absolute if needed
    let targetUrl = this.url
    if (targetUrl.startsWith('/')) {
      const isSecure = window.location.protocol === 'https:'
      const proto = isSecure ? 'wss:' : 'ws:'
      const host = window.location.host
      targetUrl = `${proto}//${host}${this.url}`
    }

    this.ws = new WebSocket(targetUrl)

    this.ws.onopen = () => {
      this.isConnect = true
      this.reconnectAttempts = 0
      this.startHeartbeat()
    }

    this.ws.onmessage = (event) => {
      let data: InternalSocketMessage<InternalMsg>
      try {
        data = JSON.parse(event.data)
      } catch (e) {
        console.error('Failed to parse WS message', e)
        return
      }

      if (data.type === 'connect') {
        this.isConnect = true
        return
      } else if (data.type === 'ping') {
        this.ws?.send(JSON.stringify({ type: 'pong' }))
        return
      } else if (data.type === 'pong') {
        clearTimeout(this.pongTimeout)
        this.pongTimeout = undefined
        return
      }

      Object.values(this.onmessage).forEach((callback) => callback(data as any))
    }

    this.ws.onclose = () => {
      this.isConnect = false
      this.stopHeartbeat()
      this.scheduleReconnect()
    }

    this.ws.onerror = (err) => {
      console.error('WebSocket error', err)
      this.ws?.close() 
    }
  }

  private scheduleReconnect() {
    if (this.reconnectTimeout) return

    const delay = Math.min(
      this.RECONNECT_DELAY * 2 ** this.reconnectAttempts,
      this.MAX_RECONNECT_DELAY,
    )

    this.reconnectTimeout = setTimeout(() => {
      this.reconnectAttempts++
      this.reconnectTimeout = null
      this.connect()
    }, delay)
  }

  private startHeartbeat() {
    this.stopHeartbeat()

    this.heartbeatInterval = setInterval(() => {
      if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return

      this.ws.send(JSON.stringify({ type: 'ping' }))

      this.pongTimeout = setTimeout(() => {
        console.warn('WebSocket pong timeout, reconnecting...')
        this.ws?.close()
      }, this.PONG_TIMEOUT)
    }, this.HEARTBEAT_INTERVAL)
  }

  private stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
      this.heartbeatInterval = null
    }
    if (this.pongTimeout) {
      clearTimeout(this.pongTimeout)
      this.pongTimeout = undefined
    }
  }

  add(callback: SocketCallback) {
    const uuid = crypto.randomUUID()
    this.onmessage[uuid] = callback
    return uuid
  }

  remove(uuid: string) {
    delete this.onmessage[uuid]
  }

  close() {
    this.stopHeartbeat()
    this.ws?.close()
  }
}

class WsManager {
  wss = new Map<string, Ws>()
  private IdsMap = new Map<string, string>()

  add(cb: SocketCallback, url: string = '/ws') {
    if (this.wss.has(url)) {
      const ws = this.wss.get(url)!
      const id = ws.add(cb)
      this.IdsMap.set(id, url)
      return id
    }

    const ws = new Ws(url)
    this.wss.set(url, ws)
    const id = ws.add(cb)
    this.IdsMap.set(id, url)
    return id
  }

  remove(id: string) {
    const url = this.IdsMap.get(id)
    if (!url) return

    const ws = this.wss.get(url)
    ws?.remove(id)
    this.IdsMap.delete(id)

    if (ws) {
      // Small delay before closing to avoid thrashing if another component re-subscribes
      setTimeout(() => {
        if (Object.keys(ws.onmessage).length === 0) {
          this.wss.delete(url)
          ws.close()
        }
      }, 3000)
    }
  }
}

export const wsManager = new WsManager()
